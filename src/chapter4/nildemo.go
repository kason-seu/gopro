package main

import "fmt"

func main() {

	var a *int
	fmt.Println(a == nil)

	var b chan int
	fmt.Println(b == nil)

	var c func()
	fmt.Println(c == nil)

	var d interface{}
	fmt.Println(d == nil)

	var e map[string]int
	fmt.Println(e == nil)

	var f []int
	fmt.Println(f == nil)

	//fmt.Println(a == b)

}
