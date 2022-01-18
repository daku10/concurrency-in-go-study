package main

import "testing"



func multiplyStream(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
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

func addStream(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
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


func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
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

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
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


func toString(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {
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

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range toString(done, take(done, repeat(done, "a"), b.N)) {

	}
}

func BenchmarkTyped(b *testing.B) {
	repeat := func(done <-chan interface{}, values ...string) <-chan string {
		valueStream := make(chan string)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <-v:						
					}
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan string, num int) <-chan string {
		takeStream := make(chan string)
		go func() {
			defer close(takeStream)
			for i := num; i > 0 || i == -1; {
				// なんでこんな書き方？
				if i != -1 {
					i--
				}
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)
	b.ResetTimer()
	for range take(done, repeat(done, "a"), b.N) {
		
	}
}