package bootstrap

import (
	"net/http"

	"caiyun/internal/handlers"
	htthandlers "caiyun/internal/interfaces/http"
	"caiyun/internal/middleware"
	"caiyun/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerRoutes(r *gin.Engine, deps RouteDependencies) {
	registerHealthAndMetricsRoutes(r, deps)
	registerAuthRoutes(r, deps.Handlers.Auth)
	registerProtectedRoutes(r, deps)
	registerAdminRoutes(r, deps)
	registerWebSocketRoute(r, deps)
	registerExportRoutes(r, deps)
	registerWebhookRoutes(r, deps)
	registerPlatformRoutes(r, deps)
}

func registerHealthAndMetricsRoutes(r *gin.Engine, deps RouteDependencies) {
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "caiyun-api"})
	}
	r.GET("/", healthHandler)
	r.GET("/health", healthHandler)
	r.GET(
		"/metrics",
		middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User),
		middleware.AdminMiddleware(),
		gin.WrapH(promhttp.HandlerFor(deps.MetricsCollector.Registry(), promhttp.HandlerOpts{})),
	)
}

func registerAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	public := r.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
		public.POST("/logout", authHandler.Logout)
		public.GET("/me", authHandler.GetCurrentUser)
		public.POST("/password/reset-code/send", authHandler.SendPasswordResetCode)
		public.POST("/password/reset", authHandler.ResetPassword)
	}
}

func registerProtectedRoutes(r *gin.Engine, deps RouteDependencies) {
	h := deps.Handlers
	api := r.Group("/api")
	api.Use(
		middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User),
		middleware.CSRFMiddleware(),
		deps.PostAuthRateMw.HandlerFunc(),
	)
	{
		api.GET("/accounts", h.Account.ListAccounts)
		api.POST("/accounts", h.Account.CreateAccount)
		api.GET("/accounts/:id", h.Account.GetAccount)
		api.PUT("/accounts/:id", h.Account.UpdateAccount)
		api.DELETE("/accounts/:id", h.Account.DeleteAccount)
		api.POST("/accounts/:id/task", h.Account.TriggerTask)
		api.PUT("/accounts/:id/status", h.Account.SetAccountStatus)
		api.POST("/accounts/:id/refresh", h.Account.RefreshToken)
		api.POST("/accounts/:id/trigger", h.Account.TriggerTask)
		api.POST("/accounts/sms/send", h.Account.SendSmsCode)
		api.GET("/accounts/sms/status/:phone", h.Account.GetSmsStatus)
		api.POST("/accounts/sms/verify", h.Account.SmsLogin)

		api.GET("/tasks/logs", h.Task.GetTaskLogs)
		api.GET("/tasks", h.Task.GetTaskLogs)
		api.GET("/tasks/:id", h.Task.GetTaskLogs)
		api.POST("/tasks/trigger-all", h.Task.TriggerAllTasks)
		api.GET("/tasks/status", h.Task.GetTaskStatus)
		api.GET("/queue/status", h.Task.GetQueueStatus)

		api.GET("/stats/dashboard", h.Task.GetDashboard)
		api.GET("/stats/cloud", h.Task.GetCloudStats)
		api.GET("/stats/trend", h.Task.GetTrendData)
		api.POST("/stats/calculate", h.Task.CalculateStats)
		api.GET("/stats/total-cloud", h.Task.GetTotalCloudCount)

		api.GET("/announcements", h.Announcement.GetPublishedAnnouncements)

		api.GET("/exchange/products", h.Exchange.SearchProducts)
		api.GET("/exchange/records", h.Exchange.GetExchangeTasks)
		api.POST("/exchange/records/batch", h.Exchange.BatchExecuteExchangeTasks)
	}
}

func registerAdminRoutes(r *gin.Engine, deps RouteDependencies) {
	h := deps.Handlers
	admin := r.Group("/api/admin")
	admin.Use(
		middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User),
		middleware.AdminMiddleware(),
		middleware.CSRFMiddleware(),
		deps.PostAuthRateMw.HandlerFunc(),
	)
	{
		admin.GET("/users", h.Admin.GetAllUsers)
		admin.PUT("/users/:id", h.Admin.UpdateUserRole)
		admin.DELETE("/users/:id", h.Admin.DeleteUser)
		admin.GET("/accounts", h.Admin.GetAllAccounts)
		admin.GET("/accounts/search", h.Admin.SearchAllAccounts)
		admin.GET("/accounts/summaries", h.Admin.GetAccountSummaries)
		admin.PUT("/accounts/:id/status", h.Admin.UpdateAccountStatus)
		admin.DELETE("/accounts/:id", h.Admin.DeleteAccount)
		admin.GET("/dashboard", h.Admin.GetAdminDashboard)
		admin.GET("/stats/overview", h.Admin.GetStatsOverview)
		admin.GET("/audit-logs", h.Admin.GetStatsOverview)
		admin.GET("/task-definitions", h.Admin.GetTaskConfigs)
		admin.PATCH("/task-definitions/:id", h.Admin.UpdateTaskConfig)
	}
}

func registerWebSocketRoute(r *gin.Engine, deps RouteDependencies) {
	r.GET("/ws", middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User), func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": "invalid token",
				"trace_id": c.GetString("X-Request-ID"),
			})
			return
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized, "message": "invalid user",
				"trace_id": c.GetString("X-Request-ID"),
			})
			return
		}
		ws.GetHub().HandleWebSocket(c.Writer, c.Request, userID)
	})
}

func registerExportRoutes(r *gin.Engine, deps RouteDependencies) {
	exportHandler := htthandlers.NewExportHandler(deps.Handlers.ExportSvc)
	exportRoutes := r.Group("/api")
	exportRoutes.Use(
		middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User),
		middleware.CSRFMiddleware(),
		deps.PostAuthRateMw.HandlerFunc(),
	)
	{
		exportRoutes.GET("/export", exportHandler.ExportData)
		exportRoutes.GET("/export/history", exportHandler.GetExportHistory)
	}
}

func registerWebhookRoutes(r *gin.Engine, deps RouteDependencies) {
	webhookHandler := htthandlers.NewWebhookHandler(deps.Handlers.WebhookSvc)
	webhookRoutes := r.Group("/api")
	webhookRoutes.Use(
		middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User),
		middleware.CSRFMiddleware(),
		deps.PostAuthRateMw.HandlerFunc(),
	)
	{
		webhookRoutes.GET("/webhooks", webhookHandler.ListEndpoints)
		webhookRoutes.POST("/webhooks", webhookHandler.CreateEndpoint)
		webhookRoutes.PUT("/webhooks/:id", webhookHandler.UpdateEndpoint)
		webhookRoutes.DELETE("/webhooks/:id", webhookHandler.DeleteEndpoint)
		webhookRoutes.POST("/webhooks/:id/test", webhookHandler.TestWebhook)
	}
}

func registerPlatformRoutes(r *gin.Engine, deps RouteDependencies) {
	platformRoutes := r.Group("/api/platforms")
	platformRoutes.Use(
		middleware.AuthMiddlewareWithUser(deps.JWTManager, deps.Repos.User),
		middleware.CSRFMiddleware(),
		deps.PostAuthRateMw.HandlerFunc(),
	)
	{
		platformRoutes.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"platforms": []string{"caiyun"}})
		})
	}
}
