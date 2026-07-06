package main

import (
	"net/http"

	"caiyun/internal/handlers"
	httphandlers "caiyun/internal/interfaces/http"
	"caiyun/internal/middleware"
	"caiyun/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func registerRoutes(r *gin.Engine, deps routeDependencies) {
	registerHealthAndMetricsRoutes(r, deps)
	registerAuthRoutes(r, deps.handlers.auth)
	registerProtectedRoutes(r, deps)
	registerAdminRoutes(r, deps)
	registerWebSocketRoute(r, deps)
	registerExportRoutes(r, deps)
	registerWebhookRoutes(r, deps)
	registerPlatformRoutes(r, deps)
}

func registerHealthAndMetricsRoutes(r *gin.Engine, deps routeDependencies) {
	healthHandler := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "caiyun-api"})
	}
	r.GET("/", healthHandler)
	r.GET("/health", healthHandler)
	r.GET(
		"/metrics",
		middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User),
		middleware.AdminMiddleware(),
		gin.WrapH(promhttp.HandlerFor(deps.metricsCollector.Registry(), promhttp.HandlerOpts{})),
	)
}

func registerAuthRoutes(r *gin.Engine, authHandler *handlers.AuthHandler) {
	public := r.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/password/reset-code/send", authHandler.SendPasswordResetCode)
		public.POST("/password/reset", authHandler.ResetPassword)
	}
}

func registerProtectedRoutes(r *gin.Engine, deps routeDependencies) {
	h := deps.handlers
	api := r.Group("/api")
	api.Use(
		middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User),
		middleware.CSRFMiddleware(),
		deps.postAuthRateMw.HandlerFunc(),
	)
	{
		api.GET("/accounts", h.account.ListAccounts)
		api.POST("/accounts", h.account.CreateAccount)
		api.GET("/accounts/:id", h.account.GetAccount)
		api.PUT("/accounts/:id", h.account.UpdateAccount)
		api.DELETE("/accounts/:id", h.account.DeleteAccount)
		api.POST("/accounts/:id/task", h.account.TriggerTask)

		api.GET("/tasks", h.task.GetTaskLogs)
		api.GET("/tasks/:id", h.task.GetTaskLogs)
		api.GET("/tasks/status", h.task.GetTaskStatus)
		api.GET("/queue/status", h.task.GetQueueStatus)

		api.GET("/announcements", h.announcement.GetPublishedAnnouncements)

		api.GET("/exchange/products", h.exchange.SearchProducts)
		api.GET("/exchange/records", h.exchange.GetExchangeTasks)
		api.POST("/exchange/records/batch", h.exchange.BatchExecuteExchangeTasks)
	}
}

func registerAdminRoutes(r *gin.Engine, deps routeDependencies) {
	h := deps.handlers
	admin := r.Group("/api/admin")
	admin.Use(
		middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User),
		middleware.AdminMiddleware(),
		middleware.CSRFMiddleware(),
		deps.postAuthRateMw.HandlerFunc(),
	)
	{
		admin.GET("/users", h.admin.GetAllUsers)
		admin.POST("/users", h.admin.GetAllUsers)
		admin.PUT("/users/:id", h.admin.UpdateUserRole)
		admin.DELETE("/users/:id", h.admin.DeleteUser)
		admin.GET("/stats", h.admin.GetStatsOverview)
		admin.GET("/audit-logs", h.admin.GetStatsOverview)
	}
}

func registerWebSocketRoute(r *gin.Engine, deps routeDependencies) {
	r.GET("/ws", middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User), func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid token"})
			return
		}
		userID, ok := userIDVal.(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "invalid user"})
			return
		}
		ws.GetHub().HandleWebSocket(c.Writer, c.Request, userID)
	})
}

func registerExportRoutes(r *gin.Engine, deps routeDependencies) {
	exportHandler := httphandlers.NewExportHandler(deps.handlers.exportSvc)
	exportRoutes := r.Group("/api")
	exportRoutes.Use(
		middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User),
		middleware.CSRFMiddleware(),
		deps.postAuthRateMw.HandlerFunc(),
	)
	{
		exportRoutes.GET("/export", exportHandler.ExportData)
		exportRoutes.GET("/export/history", exportHandler.GetExportHistory)
	}
}

func registerWebhookRoutes(r *gin.Engine, deps routeDependencies) {
	webhookHandler := httphandlers.NewWebhookHandler(deps.handlers.webhookSvc)
	webhookRoutes := r.Group("/api")
	webhookRoutes.Use(
		middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User),
		middleware.CSRFMiddleware(),
		deps.postAuthRateMw.HandlerFunc(),
	)
	{
		webhookRoutes.GET("/webhooks", webhookHandler.ListEndpoints)
		webhookRoutes.POST("/webhooks", webhookHandler.CreateEndpoint)
		webhookRoutes.PUT("/webhooks/:id", webhookHandler.UpdateEndpoint)
		webhookRoutes.DELETE("/webhooks/:id", webhookHandler.DeleteEndpoint)
		webhookRoutes.POST("/webhooks/:id/test", webhookHandler.TestWebhook)
	}
}

func registerPlatformRoutes(r *gin.Engine, deps routeDependencies) {
	platformRoutes := r.Group("/api/platforms")
	platformRoutes.Use(
		middleware.AuthMiddlewareWithUser(deps.jwtManager, deps.repos.User),
		middleware.CSRFMiddleware(),
		deps.postAuthRateMw.HandlerFunc(),
	)
	{
		platformRoutes.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"platforms": []string{"caiyun"}})
		})
	}
}
