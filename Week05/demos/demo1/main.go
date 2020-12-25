package main

import (
	"errors"
	"log"
)

func main() {
	a()
}

func a() {
	err := errors.New("a")
	defer func(err2 error) {
		if err != nil {
			log.Fatal(err)
		}
	}(err)
	err = errors.New("b")
	err = errors.New("c")
}
