package clock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSystemClock_Now(t *testing.T) {
	const tolerance = 10 * time.Millisecond
	clock := NewSystemClock()

	timeNow := time.Now()
	clockTimeNow := clock.Now()

	diff := timeNow.Sub(clockTimeNow).Abs()

	assert.Less(t, diff, tolerance, "The clock time and the actual system time differed too much. ")
}
