package timeadapter

import (
	stdtime "time"

	"github.com/Vi-72/quest-auth/internal/core/ports"
)

// Clock implements ports.Clock using the system time.
type Clock struct{}

func NewClock() *Clock {
	return &Clock{}
}

// Now returns current time.
func (c *Clock) Now() stdtime.Time {
	return stdtime.Now()
}

var _ ports.Clock = (*Clock)(nil)
