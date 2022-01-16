package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	increment := func() {
		fmt.Println("increment!")
		count++
	}
	decrement := func() {
		fmt.Println("decrement")
		count--
	}

	var once sync.Once
	once.Do(increment)
	once.Do(decrement)
	fmt.Printf("count: %d\n", count)

	var onceA, onceB sync.Once
	var initB func()
	initA := func() {onceB.Do(initB)}
	initB = func(){fmt.Println("B!");onceA.Do(initA)}
	onceA.Do(initA)
}
