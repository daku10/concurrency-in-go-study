package main

import (
	"fmt"
	"sync"
)

func main() {

	hello := func(wg *sync.WaitGroup, i int) {
		defer wg.Done()
		fmt.Printf("hello %d\n", i)
	}

	var wg sync.WaitGroup
	num := 5
	wg.Add(5)
	for i := 0; i < num; i++ {
		go hello(&wg, i)
	}

	wg.Wait()
	fmt.Println("all goroutine are finished")
}
