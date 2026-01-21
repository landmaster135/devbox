package time

import (
	"fmt"
	"strings"
	stdtime "time"
)

// NewRepository creates a repository backed by the standard library time package.
func NewRepository() Repository {
	return &stdRepository{}
}

type stdRepository struct{}

func (r *stdRepository) Now(timezone string) (stdtime.Time, error) {
	tz := strings.TrimSpace(timezone)
	if tz == "" {
		return stdtime.Time{}, fmt.Errorf("timezone must not be empty")
	}

	location, err := stdtime.LoadLocation(tz)
	if err != nil {
		return stdtime.Time{}, fmt.Errorf("load timezone %q: %w", tz, err)
	}

	return stdtime.Now().In(location), nil
}
