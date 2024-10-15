package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func recordLog(ctx context.Context) {
	select {
	case <-ctx.Done():
		// 操作被取消
		return
	default:
		// 模拟一些工作
		time.Sleep(2 * time.Second)
		fmt.Println("work done")
	}
}

// withCancel 带取消的上下文 用于取消操作
// 启动一个长时间运行的任务：例如，启动一个 goroutine 来执行某些工作。
// 在某个条件下取消任务：例如，超时、用户请求取消等。
// 清理资源：在取消任务时，确保释放相关资源。
func TestWithCancel(t *testing.T) {
	//取消机制
	ctx, cancel := context.WithCancel(context.Background())
	go recordLog(ctx)
	time.Sleep(3 * time.Second)
	cancel()
}

// 设置截止时间
func readBook(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("读书时间到")
	default:
		time.Sleep(4 * time.Second)
		fmt.Println("读书中")
	}
}

/*
超时控制：

	这种模式非常适用于需要设置操作超时的场景。例如：
	API 请求超时控制
	数据库查询超时限制
	文件读写操作的时间限制
*/
func TestWithDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	deadline, ok := ctx.Deadline()
	if ok {
		t.Log("deadline:", deadline)
	}
	defer cancel() // 手动取消

	readBook(ctx)
	select {
	case <-ctx.Done():
		fmt.Println("操作超时")
	default:
		fmt.Println("操作完成")
	}
}

func TestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // 手动取消

	readBook(ctx)
	select {
	case <-ctx.Done():
		fmt.Println("操作超时")
	default:
		fmt.Println("操作完成")
	}

}

func processRequest(ctx context.Context, userIDKey interface{}) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		// 处理错误
		fmt.Println("userID not found in context")
		return
	}
	// 使用 userID
	fmt.Println("userID:", userID)
}

// 传递请求的元数据：在处理 HTTP 请求时，可以将请求的 ID、用户 ID、认证信息等元数据存储在上下文中，以便在请求的整个生命周期内访问这些数据。
// 数据库事务：在数据库事务中，可以将事务对象存储在上下文中，以便在请求的不同部分中访问和管理事务。
// 日志记录：可以将请求的唯一标识符存储在上下文中，以便在日志中关联特定请求的日志条目。

func TestWithValue(t *testing.T) {
	type key int
	const userIDKey key = 0

	ctx := context.WithValue(context.Background(), userIDKey, "12345")

	// 在其他函数中
	processRequest(ctx, userIDKey)
}
