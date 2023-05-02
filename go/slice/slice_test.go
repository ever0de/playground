package slice_test

import (
	"testing"

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
