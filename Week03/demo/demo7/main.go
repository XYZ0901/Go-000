package main

import (
	"fmt"
)

func main() {
	a := make(map[int]int)
	b := a
	fmt.Printf("%p %p\n", a, b)
}
