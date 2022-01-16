package main

import (
	"fmt"
	"sync"
)

func main() {

	salutation := []string{"hello", "greetings", "good day"}
	var wg sync.WaitGroup
	ch := make(chan string)
	for _, word := range salutation {
		wg.Add(1)
		// 先に入れると受信されるまでブロックされる、かつ受信するgoroutineはこの後の行なので、deadlockと判定される
		// ch <- word
		go func() {
			word2 := <-ch
			fmt.Println(word2)
			wg.Done()
		}()		
		ch <- word
	}
	wg.Wait()
}
