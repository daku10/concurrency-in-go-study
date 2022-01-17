package main

import (
	"fmt"
	"sync"
)

func main() {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels"
	}()
	fmt.Println(<-stringStream)

	// 一方向チャネルとコンパイルエラーを起こせる例
	// writeStream := make(chan<- interface{})
	// readStream := make(<-chan interface{})

	// <-writeStream
	// readStream <- struct{}{}

	go func() {
		stringStream <- "Hello channels"
	}()
	salutation, ok := <- stringStream
	fmt.Printf("(%v): %v\n", ok, salutation)

	intStream := make(chan int)
	close(intStream)
	integer, ok := <-intStream
	fmt.Printf("(%v): %v\n", ok, integer)

	intStream = make(chan int)
	go func() {
		defer close(intStream)
		for i := 0; i < 5; i++ {
			intStream <- i
		}
	}()
	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
	fmt.Println()

	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v begin\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}
