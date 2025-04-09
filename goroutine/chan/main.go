package main

import "fmt"

func main() {
	ch := make(chan string)
	go func() {
		ch <- "Hello from channel"
	}()
	msg := <-ch
	fmt.Println(msg)
}
