package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"caiyun/internal/bootstrap"
	"caiyun/internal/models"
	"caiyun/internal/services"
)

func main() {
	bootstrap.LoadEnvFile()
	closeLogger := bootstrap.ConfigureStandardLogger("export")
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

	exportService := sharedServices.Export
	if exportService == nil {
		log.Fatal("ExportService 未初始化")
	}

	pollInterval := bootstrap.GetDurationEnv("EXPORT_WORKER_POLL_INTERVAL", 5*time.Second)
	batchSize := bootstrap.GetIntEnv("EXPORT_WORKER_BATCH_SIZE", 10)

	log.Printf("Export Worker 启动，轮询间隔: %v，批次大小: %d", pollInterval, batchSize)

	runExportWorkerLoop(exportService, pollInterval, batchSize)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Export Worker 已启动，按 Ctrl+C 停止")
	sig := <-sigChan
	log.Printf("Export Worker 收到停止信号: %s", sig)
}

func runExportWorkerLoop(exportService *services.ExportService, pollInterval time.Duration, batchSize int) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// 立即执行一次
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
			// 简单解析，生产环境可使用 json.Unmarshal
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
