package bloomfilter

import (
	"github.com/spaolacci/murmur3"
)

type CounterBloomFilter struct {
	bitset  []int
	hashNum uint32
}

func NewCounterBloomFilter(size int, hashNum uint32) *CounterBloomFilter {
	return &CounterBloomFilter{
		bitset:  make([]int, size),
		hashNum: hashNum,
	}
}

func (cbf *CounterBloomFilter) Add(item []byte) {
	for i := uint32(0); i < cbf.hashNum; i++ {
		hasher := murmur3.New32WithSeed(i)
		hasher.Write(item)
		hash := hasher.Sum32() % uint32(len(cbf.bitset))

		cbf.bitset[hash]++
	}
}

func (cbf *CounterBloomFilter) Contains(item []byte) bool {
	for i := uint32(0); i < cbf.hashNum; i++ {
		hasher := murmur3.New32WithSeed(i)
		hasher.Write(item)
		hash := hasher.Sum32() % uint32(len(cbf.bitset))

		if cbf.bitset[hash] == 0 {
			return false
		}
	}

	return true
}

func (cbf *CounterBloomFilter) Remove(item []byte) {
	for i := uint32(0); i < cbf.hashNum; i++ {
		hasher := murmur3.New32WithSeed(i)
		hasher.Write(item)
		hash := hasher.Sum32() % uint32(len(cbf.bitset))

		if cbf.bitset[hash] > 0 {
			cbf.bitset[hash]--
		}
	}
}
