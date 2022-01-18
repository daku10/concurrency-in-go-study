package main

import (
	"fmt"
	"time"
)

func main() {

	done := make(chan interface{})
	go func() {
		defer close(done)
		time.Sleep(2 * time.Second);
	}()

	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 0; i < 5; i++ {
			intStream <- i
			time.Sleep(1 * time.Second)
		}
	}()

	loop:
	for {
		select {
		case <-done: {
			fmt.Println("Done!")
			break loop
		}
		case i := <-intStream:
			fmt.Printf("%d\n", i)
		default:
		}
		// 上記のdefaultを消すと、caseで参照しているチャネルがブロックされているので、高々2回ほど(i := <-intStreamが読まれるタイミング)しか実行されない
		// defaultが存在している場合ループで呼ばれる
		fmt.Println("work!!")
	}

	fmt.Println("finished!")
}