package main

import (
	"fmt"
	"math/rand"
	"time"
)

func DoWork(done <-chan interface{}, nums ...int) (<-chan interface{}, <-chan int) {
	heartbeat := make(chan interface{}, 1)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		for _, n := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}

			select {
			case <-done:
				return
			case intStream <- n:				
			}
		}
	}()
	return heartbeat, intStream
}

func main() {

	doWork := func(done <-chan interface{}) (<-chan interface{}, <-chan int) {
		heartbeatStream := make(chan interface{}, 1)
		workStream := make(chan int)

		go func() {
			defer close(heartbeatStream)
			defer close(workStream)
			
			for i := 0; i < 10; i++ {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}

				select {
				case <-done:
					return
				case workStream <-rand.Intn(10):
				}
			}
		}()

		return heartbeatStream, workStream
	}

	done := make(chan interface{})
	defer close(done)

	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r, ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)				
			} else {
				return
			}
		}
	}

	// doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
	// 	heartbeat := make(chan interface{})
	// 	results := make(chan time.Time)
	// 	go func() {
	// 		// closeを忘れたgoroutine
	// 		// defer close(heartbeat)
	// 		// defer close(results)

	// 		pulse := time.Tick(pulseInterval)
	// 		workGen := time.Tick(2 * pulseInterval)

	// 		sendPulse := func() {
	// 			select { 
	// 			case heartbeat <-struct{}{}:
	// 			default:
	// 			}
	// 		}

	// 		sendResult := func(r time.Time) {
	// 			for {
	// 				select {
	// 				case <-done:
	// 					return
	// 				case <-pulse:
	// 					sendPulse()
	// 				case results <- r:
	// 					return
	// 				}
	// 			}
	// 		}

	// 		for i := 0; i < 2; i++ {
	// 			select {
	// 			case <-done:
	// 				return
	// 			case <-pulse:
	// 				sendPulse()
	// 			case r := <-workGen:
	// 				sendResult(r)
	// 			}
	// 		}
	// 	}()
	// 	return heartbeat, results
	// }

	// done := make(chan interface{})
	// time.AfterFunc(10 * time.Second, func() {close(done)})

	// const timeout = 2 * time.Second
	// heartbeat, results := doWork(done, timeout / 2)
	// for {
	// 	select {
	// 	case _, ok := <-heartbeat:
	// 		if ok == false {
	// 			fmt.Println("heart down")
	// 			return
	// 		}
	// 		fmt.Println("pulse")
	// 	case r, ok := <-results:
	// 		if ok == false {
	// 			fmt.Println("results down")
	// 			return
	// 		}
	// 		fmt.Printf("results %v\n", r.Second())
	// 		// ここで宣言すると、caseの実行時に毎回timeoutのchannelが作成される、しかしheartbeatはtimeout/2でやってくるので、closeされるまでは毎回作り直されるので期待した動作をする
	// 	case <-time.After(timeout):
	// 		fmt.Println("worker goroutine is not healthy")
	// 		return
	// 	}
	// 	fmt.Println("finish for")
	// }
}
