package clock

import (
	"time"

	clockpkg "quest-auth/internal/pkg/clock"
)

// SystemClock implements clock.Clock using the system time.
type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

// Now returns current system time.
func (c *SystemClock) Now() time.Time {
	return time.Now()
}

// Compile-time check that SystemClock implements clock.Clock
var _ clockpkg.Clock = (*SystemClock)(nil)
