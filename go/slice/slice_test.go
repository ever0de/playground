package slice_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	length := 5
	capacity := 10

	slice := make([]int, length, capacity)

	t.Log(len(slice), cap(slice))

	assert.Equal(t, length, len(slice))
	assert.Equal(t, capacity, cap(slice))
}

func TestSliceInStruct(t *testing.T) {
	type Temp struct {
		slice []int
	}

	a := Temp{
		slice: make([]int, 5, 10),
	}
	b := Temp{
		slice: make([]int, 5, 10),
	}

	assert.Equal(t, a, b)

	c := Temp{
		slice: make([]int, 3),
	}
	assert.NotEqual(t, a, c)

	a.slice[0] = 1
	num := (*int)(unsafe.Pointer(&a.slice[0]))
	assert.Equal(t, 1, *num)
}
