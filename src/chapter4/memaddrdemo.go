package main

import (
	"fmt"
	"unsafe"
)

type User struct {
	A int32
	B []int32
	C string
	D bool
	E struct{}
}

func main() {

	fmt.Printf("bool 长度： %d, 内存系数: %d\n", unsafe.Sizeof(bool(true)), unsafe.Alignof(bool(true)))
	fmt.Printf("byte 长度： %d, 内存系数: %d\n", unsafe.Sizeof(byte(0)), unsafe.Alignof(byte(0)))
	fmt.Printf("int 长度： %d, 内存系数: %d\n", unsafe.Sizeof(int(0)), unsafe.Alignof(int(0)))
	fmt.Printf("int8 长度： %d, 内存系数: %d\n", unsafe.Sizeof(int8(0)), unsafe.Alignof(int8(0)))
	fmt.Printf("int16 长度： %d, 内存系数: %d\n", unsafe.Sizeof(int16(0)), unsafe.Alignof(int16(0)))
	fmt.Printf("int32 长度： %d, 内存系数: %d\n", unsafe.Sizeof(int32(0)), unsafe.Alignof(int32(0)))
	fmt.Printf("int64 长度： %d, 内存系数: %d\n", unsafe.Sizeof(int64(0)), unsafe.Alignof(int64(0)))
	fmt.Printf("string 长度： %d, 内存系数: %d\n", unsafe.Sizeof("abc"), unsafe.Alignof("abc"))
	fmt.Printf("float32 长度： %d, 内存系数: %d\n", unsafe.Sizeof(float32(1.0)), unsafe.Alignof(float32(1.0)))
	fmt.Printf("float64 长度： %d, 内存系数: %d\n", unsafe.Sizeof(float64(1.0)), unsafe.Alignof(float64(1.0)))
	i := 10.0
	fmt.Printf("指针 长度： %d, 内存系数: %d\n", unsafe.Sizeof(&i), unsafe.Alignof(&i))
	ii := map[string]int{"a": 1}
	fmt.Printf("指针 长度： %d, 内存系数: %d\n", unsafe.Sizeof(&ii), unsafe.Alignof(&ii))

	var arr []int32
	fmt.Printf("[]int32 长度： %d, 内存系数: %d\n", unsafe.Sizeof(arr), unsafe.Alignof(arr))

	var u User
	fmt.Println(unsafe.Sizeof(u), unsafe.Alignof(u))
}
