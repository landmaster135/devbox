package env

import "fmt"

// Repository abstracts environment variable lookups needed by the workflow scheduler.
type Repository interface {
	GetEnv(envKey string) (string, error)
}

// MissingEnvError indicates that a required environment variable is not set.
type MissingEnvError struct {
	Key string
}

func (e MissingEnvError) Error() string {
	return fmt.Sprintf("environment variable %q is required", e.Key)
}
