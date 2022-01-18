package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	chanOwner := func() chan int {
		results := make(chan int, 5)
		go func() {
			defer close(results)
			for i := 0; i <= 5; i++ {
				results <- i
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for result := range results {
			fmt.Printf("Received: %v\n", result)
		}
		fmt.Println("Done receiving!")
	}

	results := chanOwner()
	consumer(results)

	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buff, "%c", b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	// 本書中ではlexical、つまり構文として問題ないようにする(この場合はアクセスするデータを分けている)
	go printData(&wg, data[:3])
	go printData(&wg, data[3:])

	wg.Wait()
}
