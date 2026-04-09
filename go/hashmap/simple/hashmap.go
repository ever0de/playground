package simple

// https://en.wikipedia.org/wiki/Poisson_distribution
const loadFactor = 0.75

type (
	// non thread-safe
	Map[K comparable, V any] struct {
		hashFunc hashFunc[K]
		seed     uintptr

		size    int
		buckets []*bucket[K, V]
	}

	bucket[K comparable, V any] struct {
		key   K
		value V
		next  *bucket[K, V]
	}

	hashFunc[K comparable] func(key *K, seed uintptr) uintptr
)

func NewMap[K comparable, V any](hashFunc hashFunc[K], seed uintptr) *Map[K, V] {
	return &Map[K, V]{
		hashFunc: hashFunc,
		seed:     seed,
		size:     0,
		buckets:  make([]*bucket[K, V], 16),
	}
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	hash := m.hashFunc(&key, m.seed)
	index := hash % uintptr(len(m.buckets))

	for b := m.buckets[index]; b != nil; b = b.next {
		if b.key == key {
			return b.value, true
		}
	}

	var zero V
	return zero, false
}

func (m *Map[K, V]) Put(key K, value V) {
	hash := m.hashFunc(&key, m.seed)
	index := hash % uintptr(len(m.buckets))

	for b := m.buckets[index]; b != nil; b = b.next {
		if b.key == key {
			b.value = value
			return
		}
	}

	m.buckets[index] = &bucket[K, V]{key: key, value: value, next: m.buckets[index]}
	m.size++

	if float64(m.size)/float64(len(m.buckets)) > loadFactor {
		m.resize()
	}
}

func (m *Map[K, V]) Size() int {
	return m.size
}

func (m *Map[K, V]) Delete(key K) {
	hash := m.hashFunc(&key, m.seed)
	index := hash % uintptr(len(m.buckets))

	var prev *bucket[K, V]
	for b := m.buckets[index]; b != nil; b = b.next {
		if b.key == key {
			if prev == nil {
				m.buckets[index] = b.next
			} else {
				prev.next = b.next
			}
			m.size--
			return
		}
		prev = b
	}
}

func (m *Map[K, V]) resize() {
	newBuckets := make([]*bucket[K, V], len(m.buckets)*2)
	for _, b := range m.buckets {
		for b != nil {
			hash := m.hashFunc(&b.key, m.seed)
			index := hash % uintptr(len(newBuckets))
			newBuckets[index] = &bucket[K, V]{key: b.key, value: b.value, next: newBuckets[index]}
			b = b.next
		}
	}
	m.buckets = newBuckets
}
