package main

import "sync"

func main() {
	tmp := sync.Map{}

	tmp.Store("a", 1)
	value, ok := tmp.Load("a")

	if !ok {
		panic("not ok")
	}

	println(value.(int))
}
