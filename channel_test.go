package main

//最后的 https://mp.weixin.qq.com/s/QgNndPgN1kqxWh-ijSofkw

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type intMap map[int]int

type ConcurrentMap struct {
	mu     sync.RWMutex
	data   intMap
	readCh chan []int
}

func (cm *ConcurrentMap) Put(k, v int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	ints := []int{k, v}
	cm.readCh <- ints
}

func (cm *ConcurrentMap) Get(k int, maxWaitingTime time.Duration) (int, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	v, ok := cm.data[k]
	if ok {
		return v, nil
	}
	select {
	case data := <-cm.readCh:
		// 类型检查
		cm.data[data[0]] = data[1]
		if k == data[0] {
			return data[1], nil
		}
	case <-time.After(maxWaitingTime):
		return 0, fmt.Errorf("超时错误。。。")
	}
	return 0, nil
}

func TestConcurrentMap(t *testing.T) {
	concurrentMap := ConcurrentMap{
		data:   make(intMap),
		readCh: make(chan []int),
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		v, _ := concurrentMap.Get(1, time.Second)
		t.Log("v", v)
	}()

	go func() {
		defer wg.Done()
		concurrentMap.Put(1, 2)
	}()

	wg.Wait()
	close(concurrentMap.readCh)
}

func TestDataRace(t *testing.T) {
	var data int
	go func() {
		data++
	}()

	if data == 0 {
		fmt.Printf("the value is %d", data)
	}
}

// 按钮
type Button struct {
	Clicked *sync.Cond
}

// Cond 和 Broadcast 是通知在 Wait 调用上阻塞的 goroutine 条件已被触发的方法。
func TestCond(t *testing.T) {

	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	// running on goroutine every function that passed/registered
	// and wait, not exit until that goroutine is confirmed to be running
	subscribe := func(c *sync.Cond, param string, fn func(s string)) {
		//保每个订阅者 goroutine 确实启动
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func(p string) {
			goroutineRunning.Done()
			c.L.Lock() // 获取条件变量关联的锁
			defer c.L.Unlock()

			fmt.Println("Registered and wait ... ")
			c.Wait() // 等待条件触发

			fn(p) //// 条件触发后执行回调函数
		}(param)

		goroutineRunning.Wait()
	}
	// 确保所有回调都执行完成
	var clickRegistered sync.WaitGroup

	for _, v := range []string{
		"Maximizing window.",
		"Displaying annoying dialog box!",
		"Mouse clicked."} {

		clickRegistered.Add(1) // 记录需要等待的回调数

		subscribe(button.Clicked, v, func(s string) {
			fmt.Println(s)
			clickRegistered.Done() // 标记回调完成
		})
	}

	// 循环处理完，触发所有等待的goroutine
	button.Clicked.Broadcast()
	// 等待所有回调执行完成
	clickRegistered.Wait()
}

// once 保证 Do 函数只执行一次
func TestOnce(t *testing.T) {
	var count int

	increment := func() {
		count++
	}

	var once sync.Once

	var increments sync.WaitGroup
	increments.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}

	increments.Wait()
	fmt.Printf("Count is %d\n", count)
}

// 连接池
func TestPool(t *testing.T) {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating new instance.")

			return struct{}{}
		},
	}
	//Get call 如果未启动实例，则在 pool 中定义的新函数
	myPool.Get()
	instance := myPool.Get()
	fmt.Println("instance", instance)
	//在这里，我们将之前检索到的实例放回池中。这会增加一个可用的实例数
	myPool.Put(instance)
	myPool.Get()

	var numCalcsCreated int

	calcPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("new calc pool")

			numCalcsCreated += 1
			mem := make([]byte, 1024)

			return &mem
		},
	}

	fmt.Println("calcPool.New", calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	calcPool.Get()

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup

	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)

			// Assume something interesting, but quick is being done with
			// this memory.
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}

// 死锁 所有并发进程都相互等待的程序。

type value struct {
	mu    sync.Mutex
	value int
}

func TestDeadLock(t *testing.T) {
	var wg sync.WaitGroup
	// 创建两个 value 实例
	printSum := func(v1, v2 *value) {
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		// deadlock
		time.Sleep(2 * time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%v\n", v1.value+v2.value)
	}
	var a, b value
	wg.Add(2)
	// 两个 goroutine 互相等待对方释放锁
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}

func TestRangeChannel(t *testing.T) {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()
	//通道关闭会退出循环
	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
}

func TestBufferCh(t *testing.T) {
	//带缓冲
	dataStream := make(chan interface{}, 4)
	//<-dataStream 读取空通道会死锁
	//dataStream <- struct{}{} //有空间缓冲区写入，不会死锁或者阻塞
	dataStream <- struct{}{}
	dataStream <- struct{}{}
	dataStream <- struct{}{}
	dataStream <- struct{}{}
	//dataStream <- struct{}{} //缓冲区满了，写入会死锁

	close(dataStream)
	//dataStream <- struct{}{} //panic: send on closed channel [recovered] 往已关闭的通道写入数据会引发 panic
	v, ok := <-dataStream //只有当缓冲区里没有数据(且)通道已关闭时，ok 才会返回 false。
	fmt.Println(v, ok)

	//无缓冲
	//var dataStream2 chan struct{}
	//dataStream2 := make(chan interface{})
	//dataStream2 <- struct{}{} //无缓冲区，写入会死锁,
	// <-dataStream2 //无缓冲区，读取会死锁

}

// 1 - 当多个频道有内容可读时会发生什么？
func TestMulCh(t *testing.T) {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		//随机从一个通道里面读取数据,50%
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func TestNoPare(t *testing.T) {
	var c <-chan int //未初始化的只读通道
	select {
	case <-c:
	//	使用默认操作来解决这个问题
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}

}
