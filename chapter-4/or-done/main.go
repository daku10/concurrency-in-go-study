package main

import "fmt"

func main() {
	orDone := func(done, c <-chan interface{}) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-c:
					if ok == false {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}

	done := make(chan interface{})
	defer close(done)
	myChan := make(chan interface{})
	go func() {
		defer close(myChan)	
		myChan <- 10
	}()
	for val := range orDone(done, myChan) {
		fmt.Printf("%d\n", val)
	}
}

