package time

import (
	stdtime "time"
)

// Repository exposes time helpers backed by the Go time package.
type Repository interface {
	// Now returns the current time in the specified timezone.
	Now(timezone string) (stdtime.Time, error)
}
