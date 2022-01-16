package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	// プールに4KB確保する
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	time.Sleep(1 * time.Second)

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func(i int) {
			defer wg.Done()

			mem := calcPool.Get()
			defer calcPool.Put(mem)
		}(i)
	}

	wg.Wait()
	fmt.Printf("%d calculators were created.", numCalcsCreated)
}
