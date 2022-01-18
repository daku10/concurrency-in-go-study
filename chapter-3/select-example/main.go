package main

import (
	"fmt"
	"time"
)

func main() {
	
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}

	c1 := make(chan interface{}); close(c1)
	c2 := make(chan interface{}); close(c2)

	var c1Count, c2Count int
	for i := 1000; i >= 0; i-- {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)

	var cnil <-chan int
	select {
	case <-cnil:
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}

	begin := time.Now()
	var cdefault1, cdefault2 <-chan int
	select {
	case <-cdefault1:
	case <-cdefault2:
	default:
		fmt.Printf("In default after %v\n\n", time.Since(begin))
	}

	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCount := 0
	loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}

		// Simulate work
		workCount++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achived %v cycles of work before signalled to stop.\n", workCount)
}
