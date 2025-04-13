package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int, 1)
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer wg.Done()
			if id > 0 {
				<-ch
			}
			fmt.Println(id)
			ch <- 1
		}(i)
	}
	ch <- 1
	wg.Wait()
}
