package main

import "sync"

func main() {
	type Map struct {
		inner sync.Map
	}

	type ReferenceOther struct {
		inner *sync.Map
	}

	syncmap := Map{
		inner: sync.Map{},
	}

	syncmap.inner.Store("a", 1)

	ref := ReferenceOther{
		inner: &syncmap.inner,
	}

	value, ok := ref.inner.Load("a")

	if !ok {
		panic("not ok")
	}

	println(value.(int))
}
