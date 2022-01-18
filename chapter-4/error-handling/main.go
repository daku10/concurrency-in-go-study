package main

import (
	"fmt"
	"net/http"
)

func main() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()
		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}

	type Result struct {
		Error error
		Response *http.Response
	}

	checkStatus2 := func(done <-chan interface{}, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
				case results <- result:					
				}
			}
		}()
		return results
	}

	done2 := make(chan interface{})
	defer close(done2)

	errorCount := 0
	urls = []string{"a", "b", "https://www.google.com", "c"}
	for result := range checkStatus2(done2, urls...) {
		if result.Error != nil {
			fmt.Printf("error :%v\n", result.Error)
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}
