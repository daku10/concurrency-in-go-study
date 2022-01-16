package main

import (
	"fmt"
	"sync"
)

func main() {

	salutation := []string{"hello", "greetings", "good day"}
	var wg sync.WaitGroup
	for _, word := range salutation {
		wg.Add(1)
		go func(word string) {
			fmt.Println(word)
			wg.Done()
		}(word)		
	}
	wg.Wait()
}