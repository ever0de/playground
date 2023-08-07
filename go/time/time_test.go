package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAfterFunc(t *testing.T) {
	duration := time.Microsecond

	called := false
	time.AfterFunc(duration, func() {
		called = true
	})

	// Wait for the callback to be called
	time.Sleep(duration * 2)

	assert.True(t, called)
}
