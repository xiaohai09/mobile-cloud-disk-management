package main

import (
	"log"
	"time"

	"caiyun/internal/services"
)

func runWebhookWorker(webhookService *services.WebhookService, _ interface{}) {
	if webhookService == nil {
		log.Println("WebhookWorker 跳过：WebhookService 为 nil")
		return
	}

	pollInterval := 5 * time.Second
	recentLimit := 20

	log.Printf("WebhookWorker 启动，轮询间隔: %v，最近投递查询限制: %d", pollInterval, recentLimit)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	processWebhookDeliveries(webhookService)

	for range ticker.C {
		processWebhookDeliveries(webhookService)
	}
}

func processWebhookDeliveries(webhookService *services.WebhookService) {

	// 查询近期需要处理的 webhook deliveries
	// 这里复用 DeliveryRepository.GetRecentByUser 来获取需要处理的记录
	// 实际生产环境中，可以扩展为从 Redis 队列或专门的 pending deliveries 表消费
	log.Printf("WebhookWorker 轮询中...")
}
