package main

import (
	"log"
	"time"

	"caiyun/internal/models"
	"caiyun/internal/services"
)

func runExportWorker(exportService *services.ExportService, _ interface{}) {
	if exportService == nil {
		log.Println("ExportWorker 跳过：ExportService 为 nil")
		return
	}

	pollInterval := 5 * time.Second
	batchSize := 10

	log.Printf("ExportWorker 启动，轮询间隔: %v，批次大小: %d", pollInterval, batchSize)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	processPendingExports(exportService, batchSize)

	for range ticker.C {
		processPendingExports(exportService, batchSize)
	}
}

func processPendingExports(exportService *services.ExportService, batchSize int) {
	if batchSize <= 0 {
		batchSize = 10
	}

	histories, _, err := exportService.GetExportHistory(0, 1, batchSize)
	if err != nil {
		log.Printf("查询导出历史失败: %v", err)
		return
	}

	var pending []*models.ExportHistory
	for _, h := range histories {
		if h.Status == "pending" {
			pending = append(pending, h)
		}
	}

	if len(pending) == 0 {
		return
	}

	log.Printf("发现 %d 条待处理导出任务", len(pending))
	for _, history := range pending {
		req := &services.ExportRequest{
			ExportType: history.ExportType,
			Format:     history.Format,
			UserID:     history.UserID,
		}
		if history.Filters != "" {
			req.Filters = map[string]interface{}{"filters": history.Filters}
		}

		log.Printf("开始处理导出任务 ID=%d, 用户=%d, 类型=%s, 格式=%s",
			history.ID, history.UserID, history.ExportType, history.Format)
		if err := exportService.ExecuteExportJob(history.ID, req); err != nil {
			log.Printf("导出任务失败 ID=%d: %v", history.ID, err)
		} else {
			log.Printf("导出任务完成 ID=%d", history.ID)
		}
	}
}
