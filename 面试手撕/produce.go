package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(ch chan<- int, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		data := i + id*10
		ch <- data
		fmt.Printf("生产者 %d 生产数据: %d", id, data)
		time.Sleep(time.Millisecond * 200)
	}
}

func consumer(ch <-chan int, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for data := range ch {
		fmt.Printf("消费者 %d 消费数据: %d", id, data)
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	dataChannel := make(chan int, 10)
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go producer(dataChannel, &wg, i)
	}
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go consumer(dataChannel, &wg, i)
	}
	go func() {
		wg.Wait()
		close(dataChannel)
	}()
	wg.Wait()
	fmt.Println("所有工作完成")
}
