package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	s := "慕课网test0001"

	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))

	fmt.Println(sh.Len)
	fmt.Println(sh.Data)
	fmt.Printf("%p", &s)

	// 这种循环获取的是底层数组的每个字节，并非字符
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}

	// 这种可以获取对应utf8 下的rune的字符
	for _, value := range s {
		fmt.Println(value)
		fmt.Printf("%c\n", value)
	}

	// 切分字符
	split_str := string(([]rune(s))[:3])
	fmt.Println(split_str)
}
