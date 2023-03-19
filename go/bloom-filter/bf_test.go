package bloomfilter_test

import (
	"testing"

	bloomfilter "github.com/ever0de/playground/bloom-filter"
	"github.com/stretchr/testify/assert"
)

func TestBloomFilter(t *testing.T) {
	bf := bloomfilter.NewBloomFilter(100, 3)

	bf.Add([]byte("apple"))
	bf.Add([]byte("banana"))
	bf.Add([]byte("orange"))

	assert.True(t, bf.Contains([]byte("apple")))
	// XXX: size = 10 일때 True (:thinking_face:?????)
	assert.False(t, bf.Contains([]byte("grape")))

	bf.Delete([]byte("banana"))

	assert.False(t, bf.Contains([]byte("banana")))
}

func TestCounterBloomFilter(t *testing.T) {
	cbf := bloomfilter.NewCounterBloomFilter(10, 3)

	cbf.Add([]byte("apple"))
	cbf.Add([]byte("banana"))
	cbf.Add([]byte("orange"))

	assert.True(t, cbf.Contains([]byte("apple")))
	assert.False(t, cbf.Contains([]byte("grape")))

	cbf.Remove([]byte("banana"))

	assert.False(t, cbf.Contains([]byte("banana")))
}
