package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("....")
}

type Resource struct {
	url        string
	polling    bool
	lastPolled int64
}

type Resources struct {
	data []*Resource
	lock *sync.Mutex
}

//PollerMutex 不鼓励
func PollerMutex(res *Resources) {
	for {
		res.lock.Lock()
		var r *Resource
		for _, v := range res.data {
			if v.polling {
				continue
			}
			if r == nil || v.lastPolled < r.lastPolled {
				r = v
			}
		}
		if r != nil {
			r.polling = true
		}
		res.lock.Unlock()
		if r == nil {
			continue
		}

		res.lock.Lock()
		r.polling = false
		r.lastPolled = time.Now().UnixNano()
		res.lock.Unlock()
	}
}

//PollerChan 鼓励
func PollerChan(in, out chan *Resource) {
	for r := range in {
		out <- r
	}
}
