package time_test

import (
	"testing"
	"time"
	_ "unsafe"

	"github.com/stretchr/testify/assert"
)

var (
	now = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
)

//go:linkname timeNow time.Now
func timeNow() time.Time {
	return now
}

func sleep(d time.Duration) {
	now = now.Add(d)
}

func TestTimeNow(t *testing.T) {
	assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())

	sleep(time.Second)

	assert.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), time.Now())
}
