package main

import (
	"fmt"
	"strconv"
)

func main() {

	itoa := strconv.Itoa(5)

	fmt.Println(itoa)

	s := ""
	arr := []byte(s)
	fmt.Println(arr)
	fmt.Println(len(arr))

	var ss string
	arr = []byte(ss)
	fmt.Println(arr)
	fmt.Println(len(arr))

}
