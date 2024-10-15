package main

import (
	"context"
	"fmt"
	"testing"
)

/*
在 Go 中，ctx.Value("traceId").(string) 使用了类型断言（Type Assertion）来获取上下文中存储的值并将其转换为字符串类型。
类型断言的语法为 value.(type)，其中 value 是要断言的值，type 是要转换到的目标类型。如果断言成功，它会返回转换后的值；如果断言失败，它会触发一个 panic。

assert：如果你有一个接口值，并且想要根据其实际类型执行不同的操作，可以使用断言。
convert：如果你想要在兼容的类型之间进行转换，可以使用类型转换。

==================================
可以相互转换的类型
1、数值类型之间的转换

	整数类型（int、int8、int16、int32、int64、uint、uint8、uint16、uint32、uint64）之间可以相互转换。
	浮点类型（float32、float64）之间可以相互转换。
	整数类型和浮点类型之间可以相互转换。

2、字符串和字节切片之间的转换
3、指针类型之间的转换
4、结构体类型和其他类型之间不能直接转换。
5、接口类型和具体类型之间的转换需要使用类型断言或类型 switch。
*/
type traceIDKey struct{}

func TestAssertion(t *testing.T) {
	// 创建一个上下文
	ctx := context.WithValue(context.Background(), "traceId", "123456")
	ctx = context.WithValue(ctx, traceIDKey{}, "12345") //需要返回一个新的context 不能直接修改原来的context
	// 从上下文中获取 traceID 对应的值

	traceID, ok := ctx.Value("traceId").(string)
	if !ok {
		// 断言失败，执行相应的处理逻辑
		t.Error("Failed to assert traceId as string")
		return
	}
	// 断言成功，可以使用 traceID 进行后续操作
	t.Logf("Trace ID: %s", traceID)

	traceID2, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		t.Error("Failed to assert traceId as string")
		return
	}
	t.Logf("Trace ID: %+v", traceID2)
}

func TestConvert(t *testing.T) {
	var iface interface{} = "hello"
	str, ok := iface.(string)
	if ok {
		t.Log(str)
	} else {
		t.Error("not a string")
	}

	floatVal, ok := iface.(float64)
	if ok {
		t.Log(floatVal)
	} else {
		t.Error("not a float64")
	}
	//	========================转换

	var intVal int = 42
	var floatVal2 float64 = float64(intVal)
	t.Log(floatVal2) // 输出：42.0

	// var strVal string = string(intVal) // 编译错误：cannot convert i (type int) to type string
	// fmt.Println(strVal)

}

type afterFuncer interface {
	AfterFunc()
}

type Parent struct{}

func (p *Parent) AfterFunc() {
	fmt.Println("Parent: AfterFunc")
}

type Child struct {
	Parent
}

// 可以断言是否实现了某个接口
func TestAssertIface(t *testing.T) {
	var parent interface{} = &Parent{}
	var child interface{} = &Child{}
	// 检查 parent 是否实现了 afterFuncer 接口
	if _, ok := parent.(afterFuncer); ok {
		t.Log("Parent implements afterFuncer")
	}

	// 检查 child 是否实现了 afterFuncer 接口
	if _, ok := child.(afterFuncer); ok {
		t.Log("Child implements afterFuncer")
	}

	// 使用类型断言调用方法
	if af, ok := parent.(afterFuncer); ok {
		af.AfterFunc()
	}
}
