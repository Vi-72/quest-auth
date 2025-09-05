package clock

import "time"

// Clock provides current time abstraction.
type Clock interface {
	Now() time.Time
}
