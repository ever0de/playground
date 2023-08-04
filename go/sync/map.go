package map_test

import (
	"sync"
	"testing"
)

func TestSyncMap(t *testing.T) {
	tmp := sync.Map{}

	tmp.Store("a", 1)
	value, ok := tmp.Load("a")

	if !ok {
		panic("not ok")
	}

	println(value.(int))
}
