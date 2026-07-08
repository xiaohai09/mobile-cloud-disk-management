package bootstrap

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"caiyun/internal/handlers"
	"caiyun/internal/middleware"
	"caiyun/internal/monitor"
	"caiyun/internal/repository"
	"caiyun/internal/services"
	"caiyun/internal/ws"
	"caiyun/pkg/jwt"
)

// RouteHandlers 汇总所有 HTTP handler 实例，供路由注册函数按需使用。
type RouteHandlers struct {
	Auth         *handlers.AuthHandler
	Account      *handlers.AccountHandler
	Task         *handlers.TaskHandler
	QueueStatus  *handlers.QueueStatusHandler
	Admin        *handlers.AdminHandler
	Exchange     *handlers.ExchangeHandler
	Announcement *handlers.AnnouncementHandler
	ExportSvc    *services.ExportService
	WebhookSvc   *services.WebhookService
}

// RouteDependencies 汇总路由注册所需的全部依赖。
type RouteDependencies struct {
	JWTManager       *jwt.Manager
	Repos            Repositories
	RateLimitConfig  *middleware.RateLimitConfig
	PostAuthRateMw   *middleware.RateLimitMiddlewareInstance
	AuditFilter      *middleware.AuditLogFilter
	AuditRepo        *repository.AuditLogRepository
	MetricsCollector *monitor.Metrics
	Handlers         RouteHandlers
}

// BootstrapResult 携带 API 服务启动后的核心对象与关闭逻辑。
type BootstrapResult struct {
	Engine  *gin.Engine
	Server  *http.Server
	Repos   Repositories
	Close   func()
}

func newJWTManagerFromEnv() *jwt.Manager {
	algorithm := strings.ToUpper(strings.TrimSpace(GetEnv("JWT_ALGORITHM", "")))
	if algorithm == "" {
		log.Fatal("JWT_ALGORITHM 必须设置为 HS256 或 RS256")
	}
	switch algorithm {
	case "RS256":
		manager, err := jwt.NewRS256Manager(
			GetSecretEnv("JWT_PRIVATE_KEY", ""),
			GetSecretEnv("JWT_PUBLIC_KEY", ""),
		)
		if err != nil {
			log.Fatalf("初始化 RS256 JWT 失败: %v", err)
		}
		return manager
	case "HS256":
		secret := GetSecretEnv("JWT_SECRET", "")
		if secret == "" {
			log.Fatal("JWT_SECRET 必须设置")
		}
		return jwt.NewManager(secret)
	default:
		log.Fatalf("不支持的 JWT_ALGORITHM: %s，仅支持 HS256 或 RS256", algorithm)
		return nil
	}
}

func BootstrapAPI() *BootstrapResult {
	core, err := InitCore()
	if err != nil {
		log.Printf("基础依赖初始化失败: %v", err)
		os.Exit(1)
	}

	// 初始化认证与仓储依赖。
	jwtExpiry := 7 * 24 * time.Hour
	jwtManager := newJWTManagerFromEnv()
	repos := core.Repository

	// 初始化服务层。
	passwordResetConfig := services.PasswordResetConfig{
		SMTP: services.SMTPConfig{
			Host:     GetEnv("SMTP_HOST", ""),
			Port:     GetEnv("SMTP_PORT", "587"),
			Username: GetEnv("SMTP_USERNAME", ""),
			Password: GetEnv("SMTP_PASSWORD", ""),
			From:     GetEnv("SMTP_FROM", ""),
			FromName: GetEnv("SMTP_FROM_NAME", "移动云盘"),
			UseTLS:   GetBoolEnv("SMTP_USE_TLS", false),
		},
	}
	authService := services.NewAuthServiceWithPasswordResetCache(repos.User, jwtManager, jwtExpiry, passwordResetConfig, core.Redis)

	sharedServices, err := InitSharedServices(core)
	if err != nil {
		log.Printf("初始化共享业务服务失败: %v", err)
		os.Exit(1)
	}
	accountService := sharedServices.Account
	taskService := sharedServices.Task
	cloudService := sharedServices.Cloud
	tokenManager := sharedServices.TokenManager
	exchangeService := sharedServices.Exchange
	taskQueue := sharedServices.TaskQueue
	log.Printf("任务队列后端: %s", TaskQueueBackendName())

	adminService := services.NewAdminService(repos.User, repos.Account, repos.TaskLog, repos.TaskConfig)

	// 初始化默认管理员账号（仅当不存在时创建）
	if err := SeedDefaultAdmin(core.Repository); err != nil {
		log.Printf("[admin] 初始化默认管理员失败: %v", err)
	}

	productService := services.NewProductService(repos.Product, repos.Account)
	announcementService := services.NewAnnouncementService(repos.Announcement)
	exportService := services.NewExportService(
		repos.ExportHistory,
		repos.Account,
		repos.TaskLog,
		repos.ExchangeRecord,
	)
	webhookService := services.NewWebhookService(
		repos.WebhookEndpoint,
		repos.WebhookDelivery,
		core.DB,
	)

	// 初始化任务监控器，并注册为 API 进程可见的全局实例。
	taskMonitor := monitor.NewTaskMonitor(monitor.Config{Logger: log.Default(), MaxHistory: 1000})
	taskMonitor.StartCleanupJob(5*time.Minute, 30*time.Minute)
	monitor.SetGlobalTaskMonitor(taskMonitor)
	metricsCollector := monitor.NewMetrics()

	// 定时同步基础监控指标到 Prometheus。
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		subCh := taskMonitor.Subscribe()
		defer taskMonitor.Unsubscribe(subCh)

		updateMetrics := func() {
			taskStats := taskMonitor.GetStats()
			running := ToInt(taskStats["active_tasks"])
			completed := ToInt(taskStats["completed_tasks"])
			total := ToInt(taskStats["total_tasks"])
			metricsCollector.SetTaskStats(total, 0, running, completed)

			tokenStats := tokenManager.GetTokenStats()
			metricsCollector.SetTokenStats(
				ToInt(tokenStats["total"]),
				ToInt(tokenStats["healthy"]),
				ToInt(tokenStats["error"]),
			)
			metricsCollector.SetAuditDropped(middleware.AuditDroppedCount())
		}

		updateMetrics()
		for {
			select {
			case <-ticker.C:
				updateMetrics()
			case _, ok := <-subCh:
				if !ok {
					return
				}
				// 收到任务状态变更后立即刷新一次指标。
				updateMetrics()
			}
		}
	}()

	// 初始化处理器。
	authHandler := handlers.NewAuthHandler(authService, jwtManager)
	accountHandler := handlers.NewAccountHandler(accountService, taskService, core.Redis)
	taskHandler := handlers.NewTaskHandler(taskService, cloudService, accountService)
	taskHandler.SetRedisCache(core.Redis)
	queueStatusHandler := handlers.NewQueueStatusHandler(taskQueue)
	adminHandler := handlers.NewAdminHandler(adminService)
	exchangeHandler := handlers.NewExchangeHandler(exchangeService, productService)
	announcementHandler := handlers.NewAnnouncementHandler(announcementService)
	auditFilter := middleware.NewAuditLogFilter()
	// 初始化全局共享的审计 writer（避免每个请求新建 worker goroutine）。
	middleware.InitGlobalAuditWriter(repos.AuditLog)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Use(middleware.RecoveryWithLogger())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.BodySizeLimitMiddleware(10 << 20)) // 10 MiB
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	// 使用高级限流中间件保护接口；保留实例引用以便优雅关闭。
	rateLimitConfig := middleware.DefaultRateLimitConfig()
	preAuthLimiter := middleware.NewAdvancedRateLimitMiddleware(rateLimitConfig)
	postAuthLimiter := middleware.NewAuthenticatedRateLimitMiddleware(rateLimitConfig)
	preAuthLimiter.SetRedisStore(core.Redis)
	postAuthLimiter.SetRedisStore(core.Redis)
	r.Use(preAuthLimiter.HandlerFunc())
	// 释放 SMS 限流器后台协程。
	defer accountHandler.Close()

	registerRoutes(r, RouteDependencies{
		JWTManager:       jwtManager,
		Repos:            repos,
		RateLimitConfig:  rateLimitConfig,
		PostAuthRateMw:   postAuthLimiter,
		AuditFilter:      auditFilter,
		AuditRepo:        repos.AuditLog,
		MetricsCollector: metricsCollector,
		Handlers: RouteHandlers{
			Auth:         authHandler,
			Account:      accountHandler,
			Task:         taskHandler,
			QueueStatus:  queueStatusHandler,
			Admin:        adminHandler,
			Exchange:     exchangeHandler,
			Announcement: announcementHandler,
			ExportSvc:    exportService,
			WebhookSvc:   webhookService,
		},
	})

	port := GetEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API 服务启动在端口 %s", port)

	return &BootstrapResult{
		Engine: r,
		Server: srv,
		Repos:  repos,
		Close: func() {
			preAuthLimiter.Stop()
			postAuthLimiter.Stop()
			middleware.StopGlobalAuditWriter()
			tokenManager.Stop()
			ws.GetHub().Stop()
			if err := core.Close(); err != nil {
				log.Printf("关闭基础依赖失败: %v", err)
			}
			taskMonitor.Stop()
		},
	}
}
