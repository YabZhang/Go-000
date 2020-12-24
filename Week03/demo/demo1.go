package main

import (
	"fmt"
	"time"
)

func demo(a int) {
	time.Sleep(time.Second)
	a++
	fmt.Println(a)
}
