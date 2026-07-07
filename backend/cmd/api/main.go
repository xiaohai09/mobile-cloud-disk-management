package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"caiyun/internal/bootstrap"
)

func main() {
	bootstrap.LoadEnvFile()
	closeLogger := bootstrap.ConfigureStandardLogger("api")
	defer closeLogger()

	result := bootstrap.BootstrapAPI()
	defer result.Close()

	go func() {
		log.Printf("API 服务启动在端口 %s", bootstrap.GetEnv("PORT", "8080"))
		if err := result.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("服务运行失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = result.Server.Shutdown(ctx)
}
