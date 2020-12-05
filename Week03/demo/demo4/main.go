package main

import (
	"sync"
	"time"
)

type t struct {
	sync.Mutex
}

var wg sync.WaitGroup

func main() {
	mem := make(map[int]t)
	for i := 0; i < 100; i++ {
		mem[i] = t{}
	}

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, v := range mem {
				v.Lock()
				for _, v2 := range mem {
					if v2 == v {
						continue
					}
					v2.Lock()
					time.Sleep(10 * time.Millisecond)
					v2.Unlock()
				}
				v.Unlock()
			}
		}()
	}
	wg.Wait()
}
