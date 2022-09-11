package main

import (
	"fmt"
)

type Car interface {
	Drive()
}

type TrafficTool interface {
	Drive()
	Blow()
}

type Vehicle struct {
	Name string
}

func (v *Vehicle) Drive() {

	fmt.Println("drive")

}

func main() {

	var d Car = &Vehicle{}
	fmt.Println(d)

	switch d.(type) {

	case TrafficTool:
		fmt.Println("vehicle is a car and traffictool interface instance")
		//case TrafficTool:
		//	fmt.Println("vehicle ia a traffictool interfance instance")

	}
}
