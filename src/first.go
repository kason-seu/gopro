package main

import "fmt"

func main() {

	fmt.Println("hello, world")

	// 空结构体的用法. 充当hashset，节省内存
	m := map[string]struct{}{}

	m["s"] = struct {
	}{}

	// chan 如果只需要传递一个信号，而不需要任何有意义的值，。可以传空结构体，不占内存. a : make(chan struct{})

}
