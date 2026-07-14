package main

import (
	"fmt"
)
func add(a uint, b uint)(uint,error){
	return a+b, nil
}


func main() {
	add,err:=add(8,8)
	if err==nil{
		fmt.Println("Its an error")
	}
	fmt.Println(add)
	
}