package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

//https://mp.weixin.qq.com/s/uNajVcWr4mZpof1eNemfmQ 资料地址
//type slice struct {
//	array unsafe.Pointer // 元素指针
//	len   int // 长度
//	cap   int // 容量
//}

// Slice作为参数其实传递的是拷贝的副本，但是邮由于其底任然为一个数组 ，因此函数内部的修改 本质上会影响到原始变量
func TestMe(t *testing.T) {
	ints := make([]int, 8) //默认将元素以0填充,容量不够先为两倍扩容
	ints = append(ints, 9)
	t.Logf("s:%v,len:%d,cap:%d", ints, len(ints), cap(ints)) //[10.....] ,1,10

	ints = make([]int, 8, 10) //虽然内存空间已经分配了，但是逻辑意义上不存在元素
	//i := ints[9] 此时的len 为8 如果访问 9仍然 index out of range 因此 只能访问 0-len 之间的元素
	t.Logf("s:%v,len:%d,cap:%d", ints, len(ints), cap(ints)) //[10.....] ,1,10

	//	打印s每个元素的地址
	s := []int{1, 2, 3}
	for i := 0; i < len(s); i++ {
		fmt.Printf("s[%d]: %p, size: %d\n", i, &s[i], unsafe.Sizeof(s[i]))
	}
	i := s[0:1] // 这是切片i的地址，其底层指向s的第一个元素 共享底层数组

	fmt.Printf("Address of s's underlying array: %p\n", &s[0])
	fmt.Printf("Address of i's underlying array: %p\n", &i[0])

	t.Logf("s:%v,%p,len:%d,cap:%d", i, &i, len(i), cap(i)) //这里打印的地址是切片的地址，而不是底层数组的地址

	t.Logf("before:s:%p", &s)

	s = append(s, 4) //在原来的数组末尾追加元素 ，末尾指的是len的位置,而非cap的位置
	//s和原来s的底层数组是同一个，但是s的地址已经改变了
	t.Logf("after:s:%p", &s)
}

/*
	 切片本身地址不会变，变动的是底层数组
		当切片的容量（cap）足够时：
		append 会直接在原切片的底层数组中追加元素，切片仍然会引用原来的底层数组，不会分配新的底层数组。因此，底层数组的地址不会改变。
		但即使底层数组不变，切片变量本身的地址（即 &s）在 Go 中可能仍然会发生变化，特别是当你打印 &s（切片变量的内存地址）时，因为 s 是一个结构体，它包含指向底层数组的指针、长度和容量字段。Go 的编译器可能会对这些结构体重新分配内存。

当切片的容量（cap）不足时：
append 操作会创建一个新的底层数组，复制旧的元素到新的数组中，并将新的元素追加到新数组中。此时，切片会指向新的底层数组。因此，切片的底层数组地址将会改变，切片变量的地址也可能会发生变化。
*/
func TestAppend(t *testing.T) {
	s := make([]int, 2, 4) // 长度为2，容量为4的切片
	s[0] = 1
	s[1] = 2

	fmt.Printf("before: s:%p, cap:%d, len:%d\n", &s, cap(s), len(s))
	fmt.Printf("before: s[0]: %p\n", &s[0]) // 底层数组的地址

	s = append(s, 3) // 此时容量够，底层数组不会改变
	fmt.Printf("after append (no reallocation): s:%p, cap:%d, len:%d\n", &s, cap(s), len(s))
	fmt.Printf("after append (no reallocation): s[0]: %p\n", &s[0]) // 底层数组的地址仍然相同

	s = append(s, 4) // 追加元素，此时容量依然够
	fmt.Printf("after append (no reallocation): s:%p, cap:%d, len:%d\n", &s, cap(s), len(s))
	fmt.Printf("after append (no reallocation): s[0]: %p\n", &s[0]) // 底层数组地址依然相同

	s = append(s, 5) // 容量不够了，会触发底层数组的重新分配
	fmt.Printf("after append (with reallocation): s:%p, cap:%d, len:%d\n", &s, cap(s), len(s))
	fmt.Printf("after append (with reallocation): s[0]: %p\n", &s[0]) // 底层数组地址发生改变
}

func TestCopy(t *testing.T) {
	s := []int{0, 1, 2, 3, 4}
	s1 := s     //直接拷贝 slice 地址都一样 简单拷贝
	s2 := s[:1] //直接拷贝 slice 地址都一样 简单拷贝
	t.Logf("address of s: %p, address of s1: %p", s, s1)
	t.Logf("address of s: %p, address of s2: %p", s, s2)
}

func TestCopy2(t *testing.T) {
	s := []int{0, 1, 2, 3, 4}
	s1 := make([]int, len(s))
	copy(s1, s)
	t.Logf("s: %v, s1: %v", s, s1)
	t.Logf("address of s: %p, address of s1: %p", s, s1)
	t.Logf("address of s[0]: %p, address of s1[0]: %p", &s[0], &s1[0]) //底层数组地址不同
}
func TestCompareArray(t *testing.T) {
	var a [4]int
	var b [3]int

	// a和b类型不同，不是。因为数组的长度是类型的一部分，这是与 slice 不同的一点。
	t.Log("a和b是否相等：", reflect.TypeOf(a) == reflect.TypeOf(b))
}

func TestPickTest(t *testing.T) {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 := slice[2:5] //[2,3,4] 指向原始slice的第三个元素
	fmt.Println("s1第一个指针的地址和slice的地址相等：", &slice[2] == &s1[0])
	fmt.Println("S1 :len", len(s1), "cap:", cap(s1))

	s2 := s1[2:6:7] //取2-6，相对于s1的第三个元素，第三个是容量索引
	fmt.Println(s2)

	s2 = append(s2, 100) //末尾追加
	s2 = append(s2, 200) //超出容量，重新分配内存

	s1[2] = 20 //这次修改不影响s2 因为s2和s1的底层数组不同

	fmt.Println(s1)    //[2,3,20]
	fmt.Println(s2)    //
	fmt.Println(slice) //slice 不变
}

func TestIncreaseCap(t *testing.T) {

	//初始增长：
	//当切片的容量不足以容纳新的元素时，Go 会创建一个新的底层数组。新数组的大小通常是当前容量的 2 倍。
	//大容量增长：
	//当切片容量大于或等于 1024 时，增长策略变为每次增加当前容量的 25%。

}

func f(s []int) {
	// i只是一个副本，不能改变s中元素的值
	/*for _, i := range s {
		i++
	}
	*/

	for i := range s {
		s[i] += 1
	}
}

func modifySlice(s []int) {
	// Modify the slice
	s[0] = 100
}

// 所有的参数传递在 Go 中技术上都是值传递。
// 对于引用类型（如切片、映射、通道等），传递的是包含指针的结构体的副本。

func TestBeParams(t *testing.T) {
	s := []int{1, 1, 1}
	//传递的是切片的副本，但是切片的底层数组是同一个
	f(s)
	//比较两者是否为同一个地址
	fmt.Println(&s[0], &s[1], &s[2])
	fmt.Println(s)

	// Create a slice
	slice := []int{1, 2, 3}

	// Print the original slice
	fmt.Println("Before:", slice)

	// Pass the slice to the function
	modifySlice(slice)

	// Print the modified slice
	fmt.Println("After:", slice)
}

func myAppend(s []int) []int {
	// 这里 s 虽然改变了，但并不会影响外层函数的 s
	s = append(s, 100) //这是一个新的切片
	return s
}

func myAppendPtr(s *[]int) {
	// 会改变外层 s 本身
	*s = append(*s, 100) //使用旧的切片
	return
}

func TestChangeSlice(t *testing.T) {
	s := []int{1, 1, 1}
	newS := myAppend(s)

	fmt.Println(s)
	fmt.Println(newS)

	s = newS

	myAppendPtr(&s)
	fmt.Println(s)
}

func TestSlice(t *testing.T) {
	ints := make([]int, 10) //[0*10]
	i := append(ints, 10)   //[0*10，10]
	t.Logf("s:%v,len:%d,cap:%d", i, len(i), cap(i))
	//[19*0,10] ,11,20
}

func TestSlice2(t *testing.T) {
	ints := make([]int, 0, 10)
	i := append(ints, 10)
	t.Logf("s:%v,len:%d,cap:%d", i, len(i), cap(i))
	//[10	],1,10
}

func TestSlice3(t *testing.T) {
	s := make([]int, 0, 20)
	i := s[8:]
	t.Logf("s:%v,len:%d,cap:%d", i, len(i), cap(i))
	//index out of range 因为此时的len==0
}

func TestSlice5(t *testing.T) {
	s := make([]int, 10, 20)
	i := s[8:9]
	t.Logf("s:%v,len:%d,cap:%d", i, len(i), cap(i))
	// [0] ,1,12 切片的容量未 start-原始cap
}

func TestSlice6(t *testing.T) {
	s := make([]int, 10, 20)
	i := s[8:] //[0,0]
	i[0] = -1  //[-1,0]
	t.Logf("s:%v", s)
	//[0....,-1,0]
}

func TestSlice7(t *testing.T) {
	s := make([]int, 10, 20)
	i := s[10]
	t.Log(i)
	//index out of range
}

func TestSlice8(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:]                              //[0,0]
	ints := append(s1, []int{10, 11, 12}...) //容量超出，切片底层数组会重新复制，重新申请容量
	v := s[10]
	t.Log(v)
	t.Log(ints)
}

func TestSlice9(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:] //[ 8,0]
	changeSlice(s1)
	t.Logf("s: %v", s) //[-1,0]
}

func changeSlice(s1 []int) {
	s1[0] = -1
}

func printUnderlyingArray(s []int, t *testing.T) {
	// 获取切片的头部信息
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&s))

	// 创建一个指向底层数组的指针
	basePtr := unsafe.Pointer(sliceHeader.Data)

	// 打印底层数组的所有元素（包括超出切片长度的部分）
	underlyingArray := make([]int, sliceHeader.Cap)
	for i := 0; i < sliceHeader.Cap; i++ {
		underlyingArray[i] = *(*int)(unsafe.Pointer(uintptr(basePtr) + uintptr(i)*unsafe.Sizeof(int(0))))
	}

	t.Logf("Underlying array: %v", underlyingArray)

}

func TestSlice10(t *testing.T) {
	s := make([]int, 10, 12)
	s1 := s[8:] //[0,0] ,cap=4 在切片原容量小于 256 的情况下，扩容时会采用原容量的2倍作为新的容量
	changeSlice1(s1)
	t.Logf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))       //[0*10] 10,12
	t.Logf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1)) //[0,0,10]   3,4
	//	打印底层数组
	// 打印底层数组
	printUnderlyingArray(s, t) // [0 0 0 0 0 0 0 0 0 0 10 0] 底层数组已经改变
}

func changeSlice1(s1 []int) {
	s1 = append(s1, 10) //未扩容 还是一样，但s1已经改变
}

func TestSlice12(t *testing.T) {
	s := []int{0, 1, 2, 3, 4}
	s = append(s[:2], s[3:]...)
	t.Logf("s: %v, len: %d, cap: %d", s, len(s), cap(s))
	v := s[4]
	// 是否会数组访问越界 不会//
	t.Log(v) //[]
}

func TestSlice13(t *testing.T) {
	s := make([]int, 512)
	s = append(s, 1) // 512 * (512 + 3*256)/4 = 832   832*sizeOf(type 8) = 然后执行内存对对齐 补齐 超过256 就会扩容
	t.Logf("len of s: %d, cap of s: %d", len(s), cap(s))
}
