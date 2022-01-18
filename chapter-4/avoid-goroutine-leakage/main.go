package main

import (
	"fmt"
	"math/rand"
	"time"
)


func main() {
	doWork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strings {
				// 何か面白い処理
				fmt.Println(s)
			}
		}()
		return completed
	}

	// nilを渡しているが、doWorkの内部的にはclose可能なchannelが来ることを期待しており、永久のブロックされている
	doWork(nil)
	fmt.Println("Done")

	doWork2 := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork2(done, nil)

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("Canceling doWork2 goroutine")
		close(done)
	}()

	<-terminated
	fmt.Println("Done")

	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randStream
	}

	done2 := make(chan interface{})
	randStream := newRandStream(done2)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done2)
	// 処理の実行中のシミュレート
	time.Sleep(1 * time.Second)
}
