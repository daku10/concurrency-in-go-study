package main

import "fmt"

func main() {
	bridge := func(done <-chan interface{}, chanStream <-chan <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <- chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for val := range stream {
					select {
					case valStream <-val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{}, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()
		return chanStream
	}

	for v := range bridge(nil, genVals()) {
		fmt.Printf("%v ", v)
	}
	fmt.Println()

	ch1 := make(chan interface{})
	ch2 := make(chan interface{})
	ch3 := make(chan interface{})

	go func() {
		fmt.Println("input ch1")
		defer close(ch1)
		ch1 <- 1
		fmt.Println("finish ch1")
	}()

	go func() {
		fmt.Println("input ch3")
		defer close(ch3)
		ch3 <- 3
		fmt.Println("finish ch3")
	}()

	go func() {
		fmt.Println("input ch2")
		// 閉じない場合bridgeのところで読み込みがブロックされる
		ch2 <- 2
		fmt.Println("finish ch2")
	}()

	genAgg := func(ch1, ch2, ch3 chan interface{}) <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			chanStream <- ch1
			chanStream <- ch2
			chanStream <- ch3
		}()
		return chanStream
	}

	agg := genAgg(ch1, ch2, ch3)

	for v := range bridge(nil, agg) {
		fmt.Printf("%v\n", v)
	}
}
