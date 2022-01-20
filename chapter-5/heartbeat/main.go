package main

import (
	"fmt"
	"time"
)

func main() {
	doWork := func(done <-chan interface{}, pulseInterval time.Duration) (<-chan interface{}, <-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			// closeを忘れたgoroutine
			// defer close(heartbeat)
			// defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select { 
				case heartbeat <-struct{}{}:
				default:
				}
			}

			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}

			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}

	done := make(chan interface{})
	time.AfterFunc(10 * time.Second, func() {close(done)})

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout / 2)
	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				fmt.Println("heart down")
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				fmt.Println("results down")
				return
			}
			fmt.Printf("results %v\n", r.Second())
			// ここで宣言すると、caseの実行時に毎回timeoutのchannelが作成される、しかしheartbeatはtimeout/2でやってくるので、closeされるまでは毎回作り直されるので期待した動作をする
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy")
			return
		}
		fmt.Println("finish for")
	}
}
