package main

import (
	"fmt"
	"log"
	"os"

	"github.com/daku10/concurrency-in-go-study/chapter-5/error-handling/intermediate"
)

func handleError(key int, err error, message string) {
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ", key))
	log.Printf("%#v", err)
	fmt.Printf("[%v] %v", key, message)
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	err := intermediate.RunJob("1")
	if err != nil {
		msg := "There was an unexpected issue: please report this as a bug."
		if _, ok := err.(intermediate.IntermediateErr); ok {
			msg = err.Error()
		}
		handleError(1, err, msg)
	}
}
