package main

import (
	"fmt"
	"time"
)

func SayHello() {
	fmt.Println("你好")
}

func main() {
	go SayHello()
	fmt.Println("hello from main")
	time.Sleep(100 * time.Millisecond)
}
