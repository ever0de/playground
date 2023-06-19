package mvcc

import "sync/atomic"

type DataMap map[string]VersionMap
type VersionMap map[uint64]string

type MVCC struct {
	data   DataMap
	commit uint64
}

func New() *MVCC {
	return &MVCC{
		data:   make(DataMap),
		commit: 0,
	}
}

func (m *MVCC) Write(key string, value string) uint64 {
	var version uint64
	for {
		currentCommit := atomic.LoadUint64(&m.commit)
		newCommit := currentCommit + 1

		if m.data[key] == nil {
			m.data[key] = make(map[uint64]string)
		}
		m.data[key][newCommit] = value

		if atomic.CompareAndSwapUint64(&m.commit, currentCommit, newCommit) {
			version = newCommit
			break
		}
	}

	return version
}

func (m *MVCC) Read(key string, commit uint64) (string, bool) {
	if m.data[key] == nil {
		return "", false
	}

	for i := commit; i > 0; i-- {
		if val, ok := m.data[key][i]; ok {
			return val, true
		}
	}

	return "", false
}

func (m *MVCC) ReadLatestVersion(key string) (string, bool) {
	return m.Read(key, atomic.LoadUint64(&m.commit))
}

func (m *MVCC) Version() uint64 {
	return atomic.LoadUint64(&m.commit)
}
