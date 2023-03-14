package main

import (
	"sync"
)

func main() {
	addr := "localhost:4242"

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		NewClientStream(addr)
		wg.Done()
		defer println("done client")
	}()

	go func() {
		NewServerStream(addr)
	}()

	wg.Wait()
}
