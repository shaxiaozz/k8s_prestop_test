package main

import (
	"context"
	"fmt"
	"github.com/wonderivan/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 创建一个接收信号的通道
	sigChan := make(chan os.Signal, 1)
	// 监听特定信号，如 SIGINT (Ctrl+C) 和 SIGTERM
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个带有取消功能的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动一个 goroutine 监听信号
	go func() {
		sig := <-sigChan
		logger.Info("Received signal: ", sig)
		cancel() // 取消上下文
	}()

	done := make(chan struct{})

	go func() {
		runTask(ctx)
		close(done)
	}()

	select {
	case <-done:
		logger.Info("Task completed gracefully.")
	case <-time.After(10000 * time.Second):
		logger.Info("Timeout reached, forcing shutdown.")
	}

	logger.Info("Program exiting.")
}

func runTask(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Shutting down task...")
			return
		case t := <-ticker.C:
			fmt.Printf("Working... %v\n", t)
		}
	}
}
