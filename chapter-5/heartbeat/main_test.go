package main

import (
	"testing"
)

func TestDoWork_GeneratesAllNumbers(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := DoWork(done, intSlice...)

	// これいる？
	// <-heartbeat

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected: %v, but received :%v, ", i, expected, r)
		}
		i++
	}


	// for i, expected := range intSlice {
	// 	select {
	// 	case r := <-results:
	// 		if r != expected {
	// 			t.Errorf("index %v: expected %v, but received: %v, ", i, expected, r)
	// 		}
	// 	case <-time.After(1 * time.Second):
	// 		t.Fatal("test timed out")
	// 	}
	// }
}