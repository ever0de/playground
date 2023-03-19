package bloomfilter

import (
	"hash"
	"hash/fnv"

	"github.com/bits-and-blooms/bitset"
)

type BloomFilter struct {
	bits        *bitset.BitSet
	hashFuncs   []hash.Hash64
	size        uint
	numElements uint
}

// BloomFilter 생성
func NewBloomFilter(size uint, numHashFuncs uint) *BloomFilter {
	bf := &BloomFilter{
		bits:      bitset.New(size),
		hashFuncs: make([]hash.Hash64, numHashFuncs),
		size:      size,
	}
	for i := range bf.hashFuncs {
		bf.hashFuncs[i] = fnv.New64()
	}
	return bf
}

// BloomFilter에 요소 추가
func (bf *BloomFilter) Add(element []byte) {
	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		hashFunc.Write(element)
		index := hashFunc.Sum64() % uint64(bf.size)
		bf.bits.Set(uint(index))
	}
	bf.numElements++
}

// BloomFilter에서 요소 삭제
func (bf *BloomFilter) Delete(element []byte) {
	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		hashFunc.Write(element)
		index := hashFunc.Sum64() % uint64(bf.size)
		if bf.bits.Test(uint(index)) {
			bf.bits.Clear(uint(index))
			bf.numElements--
		}
	}
}

// BloomFilter에서 요소 존재 여부 확인
func (bf *BloomFilter) Contains(element []byte) bool {
	for _, hashFunc := range bf.hashFuncs {
		hashFunc.Reset()
		hashFunc.Write(element)
		index := hashFunc.Sum64() % uint64(bf.size)
		if !bf.bits.Test(uint(index)) {
			return false
		}
	}
	return true
}

// func main() {
// 	bf := NewBloomFilter(100, 3)
// 	bf.Add([]byte("hello"))
// 	bf.Add([]byte("world"))
// 	fmt.Println(bf.Contains([]byte("hello"))) // true
// 	fmt.Println(bf.Contains([]byte("world"))) // true
// 	bf.Delete([]byte("hello"))
// 	fmt.Println(bf.Contains([]byte("hello"))) // false
// 	fmt.Println(bf.Contains([]byte("world"))) // true
// }
