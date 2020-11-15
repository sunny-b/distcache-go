package main

import (
	"fmt"

	"github.com/sunny-b/distcache-go/list"
)

func main() {
	val := 5

	l := list.New().InsertRoot(val)

	var err error
	l, err = l.Insert(6, 1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(l)
	fmt.Println(l.Length())

}
