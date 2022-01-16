package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var o sync.Once

	m := new(sync.Mutex)
	c := sync.NewCond(m)

	for i := 0; i < 10; i++ {
		go func(i int) {
			o.Do(func() {
				fmt.Println("once!")
			})
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

	p := sync.Pool{
		New: func() interface{} {
			return "定時作業"
		},
	}
	
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			p.Put("割込作業")
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			fmt.Println(p.Get())
			time.Sleep(500 * time.Millisecond)
		}
	}()

	wg.Wait()
}