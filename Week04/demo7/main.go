package main

import (
	"time"
)

type Option func(o *options) Option

type options struct {
	stopTimeout time.Duration
}

func StopTimeout(d time.Duration) Option {
	return func(o *options) Option {
		prev := o.stopTimeout
		o.stopTimeout = d
		return StopTimeout(prev)
	}
}

func main() {
	o := options{}
	f := StopTimeout(6)
	pre := f(&o)
	defer pre(&o)
}
