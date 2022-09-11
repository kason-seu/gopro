package main

import "fmt"

func main() {

	// 一开始申明空接口，零值就是nil
	var a interface{}
	fmt.Println(a == nil)

	var b *int

	// 一旦空接口里面有了类型即*int后，即使value还是nil，但是该接口已经不是nil了
	a = b
	fmt.Println(a == nil)

}
