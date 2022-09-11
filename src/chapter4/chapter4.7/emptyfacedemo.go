package main

import "fmt"

func empty(e interface{}) {
	fmt.Println(e)
}
func main() {

	empty("abc")

	empty(12345)
}
