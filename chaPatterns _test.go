package main

import (
	"fmt"
	"testing"
	"time"
)

//通过约定和设计模式来实现数据在并发上下文中的临时限制（ad hoc confinement），即确保只有一个并发进程可以访问共享数据。

func TestConfinement(t *testing.T) {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}

}

func TestCancellation(t *testing.T) {
	doWork := func(
		done <-chan interface{},
		strings <-chan string,
	) <-chan interface{} {

		terminated := make(chan interface{})

		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
		}()

		for {
			select {
			case s := <-strings:
				fmt.Println(s)
			case <-done:
				return nil
			}
		}
	}
	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		// Cancel the operation after 1 second.
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}

// 一个或多个已完成通道合并为单个已完成通道
func TestOrCh(t *testing.T) {
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				// 递归调用 or 函数，以便将所有通道合并为一个。
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()

		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
