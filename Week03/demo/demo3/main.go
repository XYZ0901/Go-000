package main

import "fmt"

func main() {
	var ben = &Ben{
		id:   0,
		name: "Ben",
	}
	var jerry = &Jerry{name: "Jerry"}
	var maker IceCreamMaker

	var loop0, loop1 func()

	loop0 = func() {
		maker = ben
		go loop1()
	}

	loop1 = func() {
		maker = jerry
		go loop0()
	}

	go loop0()

	for {
		maker.Hello()
	}

}

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	id   int // 如果这里注释掉 Ben的结构与Jerry的结构一致 将不会发生data race
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", b.name)
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	fmt.Printf("Ben says, \"Hello my name is %s\"\n", j.name)
}
