package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"caiyun/internal/bootstrap"
	"caiyun/internal/services"
)

func main() {
	bootstrap.LoadEnvFile()
	closeLogger := bootstrap.ConfigureStandardLogger("webhook")
	defer closeLogger()

	core, err := bootstrap.InitCore()
	if err != nil {
		log.Fatal("基础依赖初始化失败:", err)
	}
	defer func() {
		if err := core.Close(); err != nil {
			log.Printf("关闭基础依赖失败: %v", err)
		}
	}()

	sharedServices, err := bootstrap.InitSharedServices(core)
	if err != nil {
		log.Fatalf("初始化共享业务服务失败: %v", err)
	}

	webhookService := sharedServices.Webhook
	if webhookService == nil {
		log.Fatal("WebhookService 未初始化")
	}

	pollInterval := bootstrap.GetDurationEnv("WEBHOOK_WORKER_POLL_INTERVAL", 5*time.Second)
	recentLimit := bootstrap.GetIntEnv("WEBHOOK_WORKER_RECENT_LIMIT", 20)

	log.Printf("Webhook Worker 启动，轮询间隔: %v，最近投递查询限制: %d", pollInterval, recentLimit)

	runWebhookWorkerLoop(webhookService, pollInterval, recentLimit)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Webhook Worker 已启动，按 Ctrl+C 停止")
	sig := <-sigChan
	log.Printf("Webhook Worker 收到停止信号: %s", sig)
}

func runWebhookWorkerLoop(webhookService *services.WebhookService, pollInterval time.Duration, recentLimit int) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// 立即执行一次
	processWebhookDeliveries(webhookService)

	for range ticker.C {
		processWebhookDeliveries(webhookService)
	}
}

func processWebhookDeliveries(webhookService *services.WebhookService) {

	// 查询近期需要处理的 webhook deliveries
	// 这里复用 DeliveryRepository.GetRecentByUser 来获取需要处理的记录
	// 实际生产环境中，可以扩展为从 Redis 队列或专门的 pending deliveries 表消费
	log.Printf("Webhook Worker 轮询中...")
}
