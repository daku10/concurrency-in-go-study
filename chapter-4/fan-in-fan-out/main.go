package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func toInt(done <-chan interface{}, valueStream <-chan interface{}) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case intStream <- v.(int):
			}
		}
	}()
	return intStream
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case ch <- fn():
			}
		}
	}()
	return ch
}

func take(done <-chan interface{}, valueStream <-chan int, num int) <-chan int {
	takeStream := make(chan int)
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <- valueStream:
			}
		}
	}()
	return takeStream
}

func primeFinder(done <-chan interface{}, valueStram <-chan int) <-chan int {
	primeStream := make(chan int)
	go func() {
		defer close(primeStream)
		for integer := range valueStram {
			integer -= 1
			prime := true
			for divider := integer - 1; divider > 1; divider-- {
				if integer % divider == 0 {
					prime = false
					break
				}
			}
			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()
	return primeStream
}


func main() {
	rand := func() interface{} {return rand.Intn(50000000)}

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v\n", time.Since(start))


	fanIn := func(done <-chan interface{}, channels ...<-chan int) <-chan int {
		var wg sync.WaitGroup
		multiplexedStream := make(chan int)

		multiplex := func(c <-chan int) {
			defer wg.Done()
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
			}
		}

		wg.Add(len(channels))
		for _, c := range channels {
			go multiplex(c)
		}

		go func() {
			wg.Wait()
			close(multiplexedStream)
		}()

		return multiplexedStream
	}

	done2 := make(chan interface{})
	defer close(done2)

	start2 := time.Now()
	numFinders := runtime.NumCPU()
	finders := make([]<-chan int, numFinders)
	for i:= 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime2 := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime2)
	}
	fmt.Printf("Search took: %v\n", time.Since(start2))
}