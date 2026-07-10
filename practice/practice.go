package main

import (
	"fmt"
)

func main() {
	// var a=10
	var b = 13
	var c = &b
	var d = *c
	fmt.Println(c)
	fmt.Println(d)

}