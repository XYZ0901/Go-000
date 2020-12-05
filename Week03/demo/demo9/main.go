package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	Wait    sync.WaitGroup
	Counter int = 0
)

func main() {
	for routine := 1; routine <= 2; routine++ {
		Wait.Add(1)
		go Routine(Counter)
	}
	Wait.Wait()
	fmt.Printf("Final Counter: %d\n", Counter)
}

func Routine(id int) {
	for count := 0; count < 2; count++ {
		value := Counter
		time.Sleep(1 * time.Nanosecond)
		value++
		Counter = value
	}
	Wait.Done()
}
