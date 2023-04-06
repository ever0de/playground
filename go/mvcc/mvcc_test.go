package mvcc_test

import (
	"testing"

	"github.com/ever0de/playground/mvcc"
	"github.com/stretchr/testify/assert"
)

func TestMVCC(t *testing.T) {
	mvcc := mvcc.New()

	// write values for key "name"
	mvcc.Write("name", "Bob")
	mvcc.Write("name", "Charlie")

	// read the latest value for key "name"
	if val, ok := mvcc.Read("name", mvcc.Version()); ok {
		assert.Equal(t, "Charlie", val)
	}

	// read the previous value for key "name"
	if val, ok := mvcc.Read("name", 1); ok {
		assert.Equal(t, "Bob", val)
	}
}
