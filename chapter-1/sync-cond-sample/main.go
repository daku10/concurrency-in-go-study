package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	m := new(sync.Mutex)
	c := sync.NewCond(m)

	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("wait %d\n", i)
			c.L.Lock()
			defer c.L.Unlock()
			fmt.Printf("locking... %d\n", i)
			c.Wait()
			fmt.Printf("go %d\n", i)
		}(i)
	}

	// for i := 0; i < 10; i++ {
	// 	time.Sleep(1 * time.Second)
	// 	c.Signal()
	// }
	// for i := 3; i > 0; i-- {
	// 	fmt.Printf("%d\n", i)
	// 	time.Sleep(1 * time.Second)
	// }

	// goroutineがwaitになる前にBroadcastしても意味がないためSleepが入っている(もちろんアンチパターン)
	time.Sleep(1 * time.Second)
	c.Broadcast()
	time.Sleep(3 * time.Second)
}