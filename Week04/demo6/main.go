package main

import "time"

type Option func(o *options)

type options struct {
	stopTimeout time.Duration
}

func StopTimeout(d time.Duration) Option {
	return func(o *options) {
		o.stopTimeout = d
	}
}

func main() {
	o := options{}
	f := StopTimeout(6)
	f(&o)
}
