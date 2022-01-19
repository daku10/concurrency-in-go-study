package main

import (
	"fmt"
	"time"
)

func take(done <-chan interface{}, num int, inputStream <-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case ch <- <- inputStream:
			}
		}
	}()
	return ch
}

func repeat(done <-chan interface{}, val interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case ch <- val:				
			}
		}
	}()
	return ch
}

func sleep(done <-chan interface{}, delay time.Duration, inputStream <-chan interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		startSleep := time.Now()
		defer close(ch)
		for i := range inputStream {
			time.Sleep(delay)
			select {
			case <-done:
				return
			case ch <- i:
			}
		}
		fmt.Printf("finished: %v\n", time.Since(startSleep))
	}()
	return ch
}

func buffer(done <-chan interface{}, num int, inputCh <-chan interface{}) <-chan interface{} {
	ch := make(chan interface{}, num)
	go func() {
		defer close(ch)
		for i := range inputCh {
			select {
			case <-done:
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

func main() {
	done := make(chan interface{})
	defer close(done)

	zeros := take(done, 3, repeat(done, 0))
	short := sleep(done, 1 * time.Second, zeros)
	buff := buffer(done, 2, short)
	long := sleep(done, 4 * time.Second, buff)
	pipeline := long

	start := time.Now()
	for v := range pipeline {
		fmt.Println(v)
	}
	fmt.Printf("elapsed: %v\n", time.Since(start))
}
