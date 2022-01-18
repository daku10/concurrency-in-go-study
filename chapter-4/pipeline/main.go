package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	multiply := func(values []int, multiplier int) []int {
		multipliedValues := make([]int, len(values))
		for i, v := range values {
			multipliedValues[i] = v * multiplier
		}
		return multipliedValues
	}

	add := func(values []int, additive int) []int {
		addedValues := make([]int, len(values))
		for i, v := range values {
			addedValues[i] = v + additive
		}
		return addedValues
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}

	generator := func(done <-chan interface{}, integers ...int) <-chan int {
		intStream := make(chan int, len(integers))
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiplyStream := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipliedStream := make(chan int)
		go func() {
			defer close(multipliedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipliedStream <- i * multiplier:
				}
			}
		}()
		return multipliedStream
	}

	addStream := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i + additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4)
	pipeline := multiplyStream(done, addStream(done, multiplyStream(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}

	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valuesStream := make(chan interface{})
		go func() {
			defer close(valuesStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valuesStream <- v:						
					}
				}
			}
		}()
		return valuesStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
					// <- が1つだけだとvalueStreamのアドレスっぽいものを入れることになってそう(interface{}だと型が何でもできるのがしんどい気がする)
				case takeStream <- <- valueStream:
				}
			}
		}()
		return takeStream
	}

	done2 := make(chan interface{})
	defer close(done2)

	for num := range take(done2, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}

	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		repeatFnStream := make(chan interface{})
		go func() {
			defer close(repeatFnStream)
			for {
				select {
				case <-done:
					return
				case repeatFnStream <- fn():
				}
			}
		}()
		return repeatFnStream
	}

	done3 := make(chan interface{})
	defer close(done3)

	rand := func() interface{} {return rand.Int()}

	for num := range take(done3, repeatFn(done3, rand), 10) {
		fmt.Println(num)
	}

	toString := func(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- v.(string):
				}
			}
		}()
		return stringStream
	}

	done4 := make(chan interface{})
	defer close(done4)
	var message string
	for token := range toString(done4, take(done4, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)

	test := func(done <-chan interface{}, values <-chan int) {
		tCh := make(chan int, 20)
		go func() {
			defer close(tCh)

			// ここでブロックは起きるので、case <-done: はそんなに行かないような...
			for v := range values {
				select {
				case <-done:
					fmt.Println("done test")
					return
				case tCh <- v:
					fmt.Printf("v:\n")
					<-tCh
				}
			}
		}()
	}

	doneTest := make(chan interface{})
	testCh := make(chan int)
	test(doneTest, testCh)
	go func() {
		for i := 0; i < 10; i++ {
			testCh <- i
			fmt.Printf("push i: %d\n", i)
			time.Sleep(200 * time.Second)
			// testCh <- i
		}
	}()
	time.Sleep(1 * time.Second)
	close(doneTest)
	fmt.Println("finish!")
	time.Sleep(10 * time.Second)
}
