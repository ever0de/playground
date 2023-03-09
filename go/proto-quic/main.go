package main

import "sync"

func main() {
	addr := "localhost:4242"

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		NewClient(addr)
		wg.Done()
	}()

	go func() {
		NewServer(addr)
		wg.Done()
	}()

	wg.Wait()
}
