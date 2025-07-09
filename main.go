package main

import (
	"context"
	"fmt"
	"time"
)

func worker(ctx context.Context, id int) {
	fmt.Printf("Worker %d 启动\n", id)

	select {
	case <-ctx.Done():
		fmt.Printf("Worker %d 收到取消信号，退出\n", id)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	// 启动多个 goroutine
	for i := 1; i <= 3; i++ {
		go worker(ctx, i)
	}

	// 模拟等待 2 秒
	time.Sleep(2 * time.Second)

	fmt.Println("主线程调用 cancel()")
	cancel()

	// 等待所有 goroutine 打印完
	time.Sleep(1 * time.Second)
	fmt.Println("主线程退出")
}
