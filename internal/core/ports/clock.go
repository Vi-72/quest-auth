package ports

import "time"

// Clock provides current time.
type Clock interface {
	Now() time.Time
}
