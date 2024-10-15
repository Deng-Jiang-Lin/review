package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	//获取当前时间
	now := time.Now()
	fmt.Println("Current Time:", now)

	//固定格式 一个时间模板
	formattedTime := now.Format("2006-01-02 15:04:05")
	fmt.Println("Formatted Time:", formattedTime)

	//时间解析
	layout := "2006-01-02 15:04:05"
	timeStr := "2024-10-15 13:45:00"
	parsedTime, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
	}
	fmt.Println("Parsed Time:", parsedTime)

	// 加 2 小时
	future := now.Add(2 * time.Hour)
	fmt.Println("Future Time:", future)

	// 计算两个时间差
	duration := future.Sub(now)
	fmt.Println("Duration:", duration)

	//时间间隔
	duration = 2 * time.Hour
	fmt.Println("Duration:", duration)

	// 将 Duration 转换为其他单位
	fmt.Println("Hours:", duration.Hours())
	fmt.Println("Minutes:", duration.Minutes())
	fmt.Println("Seconds:", duration.Seconds())
	fmt.Println("Milliseconds:", duration.Milliseconds())

	//计时器
	ticker := time.NewTicker(1 * time.Microsecond)
	for i := 0; i < 3; i++ {
		//阻塞等待
		<-ticker.C
		fmt.Println("Tick at", time.Now())
	}
	ticker.Stop()

	//	时间比较
	time1 := time.Now()
	time2 := time1.Add(2 * time.Hour)
	fmt.Println("Time1 is before Time2:", time1.Before(time2))
	fmt.Println("Time1 is after Time2:", time1.After(time2))
	fmt.Println("Time1 is equal to Time2:", time1.Equal(time2))

	//超时控制
	select {
	case <-time.After(2 * time.Second):
		fmt.Println("Timeout after 2 seconds")
	}

}
