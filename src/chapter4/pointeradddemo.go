package main

import (
	"fmt"
	"unsafe"
)

// 指针的长度都是8
func main() {
	var c *int
	fmt.Println(unsafe.Sizeof(c))
	var d *map[string]int
	fmt.Println(unsafe.Sizeof(d))
}
