package main

import (
	"fmt"
	"testing"
	"unsafe"
)

//unsafe.Pointer 是一种特殊的指针类型,它可以指向任何类型的数据。它的主要用途是在不同的指针类型之间进行转换,或者在指针和整数之间进行转换。
//可以绕过 Go 的类型系统,实现任意类型之间的转换。

func TestPointerConvert(t *testing.T) {
	i := 42
	pointer := unsafe.Pointer(&i)
	f := (*float32)(pointer) // unsafe.Pointer 转换为 *float32
	fmt.Printf("i = %d\n", i)
	fmt.Printf("*f = %f\n", *f)
}

// 指针读写运算
func TestCal(t *testing.T) {
	arr := [3]int{1, 2, 3}
	p := unsafe.Pointer(&arr[0])
	for i := 0; i < len(arr); i++ {
		element := *(*int)(unsafe.Pointer(uintptr(p) + uintptr(i)*unsafe.Sizeof(arr[0])))
		fmt.Printf("Element %d: %d\n", i, element)
	}

}

type secretStruct struct {
	a int
	b string
}

// 访问结构体的未导出(私有)字段
func TestPrivateAttr(t *testing.T) {
	s := secretStruct{1, "hello"}
	p := unsafe.Pointer(&s)

	bPtr := (*string)(unsafe.Pointer(uintptr(p) + unsafe.Offsetof(s.b)))
	*bPtr = "world"

	fmt.Printf("%+v\n", s)
}

func TestZeroCopy(t *testing.T) {

}
