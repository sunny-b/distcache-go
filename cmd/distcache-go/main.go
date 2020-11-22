package main

import (
	"fmt"

	"distcache-go/list"
)

func main() {
	val := 5

	l := list.New()
	err := l.Unshift(val)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(l)
	fmt.Println(l.Length())

}
