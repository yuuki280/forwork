package main

import (
	"fmt"
	"sync"
)

func main() {
	letter, number := make(chan bool), make(chan bool)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 26; i += 2 {
			<-number
			fmt.Print(i)
			i++
			fmt.Print(i)
			letter <- true
		}
	}()
	go func() {
		defer wg.Done()
		for i := 'a'; i <= 'z'; i += 2 {
			<-letter
			fmt.Print(string(i))
			i++
			fmt.Print(string(i))
			if i < 'z' {
				number <- true
			}
		}
	}()
	number <- true
	wg.Wait()
}
