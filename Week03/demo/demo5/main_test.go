package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

type Config struct {
	a []int
}

func (c *Config) T() {}

func BenchmarkAtomic(b *testing.B) {
	var v atomic.Value
	v.Store(&Config{})

	go func() {
		i := 0
		for {
			i++
			cfg := &Config{[]int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			v.Store(cfg)
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				cfg := v.Load().(*Config)
				cfg.T()
				//fmt.Printf("%v\n", cfg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkRWMutex(b *testing.B) {
	l := sync.RWMutex{}
	cfg := &Config{}

	go func() {
		i := 0
		for {
			i++
			l.RLock()
			cfg = &Config{[]int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			l.RUnlock()
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				l.RLock()
				cfg.T()
				//fmt.Printf("%v\n", cfg)
				l.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkMutex(b *testing.B) {
	l := sync.Mutex{}
	cfg := &Config{}

	go func() {
		i := 0
		for {
			i++
			l.Lock()
			cfg = &Config{[]int{i, i + 1, i + 2, i + 3, i + 4, i + 5}}
			l.Unlock()
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 4; n++ {
		wg.Add(1)
		go func() {
			for n := 0; n < b.N; n++ {
				l.Lock()
				cfg.T()
				//fmt.Printf("%v\n", cfg)
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
