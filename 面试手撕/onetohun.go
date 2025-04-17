package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1, ch2 := make(chan bool), make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 99; i += 2 {
			<-ch1
			fmt.Println(i)
			ch2 <- true
		}
	}()
	go func() {
		defer wg.Done()
		for i := 2; i <= 100; i += 2 {
			<-ch2
			fmt.Println(i)
			if i < 100 {
				ch1 <- true
			}
		}
	}()
	ch1 <- true
	wg.Wait()
}
