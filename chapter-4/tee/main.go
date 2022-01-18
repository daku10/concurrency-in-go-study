package main

import (
	"fmt"
	"time"
)

func main() {
	tee := func(done <-chan interface{}, in <-chan interface{}) (_, _ <-chan interface{}) {
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range in {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case out1<-val:
						out1 = nil
					case out2<-val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}

	done := make(chan interface{})
	defer close(done)
	ch := make(chan interface{})
	go func() {
		defer close(ch)
		for i := 0; i < 5; i++ {
			ch <- i
		}
	}()
	out1, out2 := tee(done, ch)

	go func() {		
		for {
			fmt.Printf("out2: %v\n", <-out2)
			time.Sleep(1 * time.Second)
		}
	}()

	for val1 := range out1 {
		fmt.Printf("out1: %v\n", val1)
	}
}
