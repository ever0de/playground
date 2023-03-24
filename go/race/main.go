package main

import (
	"fmt"
	"sync"
)

// NOTE: Run `go run --race .`

type Conflict struct {
	stub []byte
}

var GlobalFunctions = func(haha *Conflict) {
	haha.stub = []byte("haha")
	fmt.Printf("haha: %s", haha.stub)
}

func main() {
	conflict := &Conflict{}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		// NOTE: Need function copy
		GlobalFunctions(conflict)
		wg.Done()
	}()
	go func() {
		// NOTE: Need function copy
		GlobalFunctions(conflict)
		conflict.stub = []byte("baba")
		fmt.Printf("baba: %s", conflict.stub)
		wg.Done()
	}()

	wg.Wait()
}
