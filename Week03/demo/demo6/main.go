package main

import (
	"sync"
	"time"
)

func main() {
	done := make(chan bool, 1)
	var mu sync.Mutex

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				mu.Lock()
				time.Sleep(100 * time.Microsecond)
				mu.Unlock()
			}
		}
	}()
	for i := 0; i < 10; i++ {
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		mu.Unlock()
	}
	done <- true
}
