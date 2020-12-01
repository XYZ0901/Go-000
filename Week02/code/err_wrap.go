package code

import (
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	err := rap2()
	fmt.Printf("%+v", err)
}

func rap1() error {
	return errors.New("test err")
}

func rap2() error {
	return rap1()
}

// **print out**
//test err
//main.rap1
//        .../main.go:14
//main.rap2
//        .../main.go:19
//main.main
//        .../main.go:9
//runtime.main
//        /usr/local/go/src/runtime/proc.go:204
//runtime.goexit
//        /usr/local/go/src/runtime/asm_amd64.s:1374
