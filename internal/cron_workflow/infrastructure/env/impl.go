package env

import (
	"fmt"
	"os"
	"strings"
)

// NewRepository creates the default environment repository implementation.
func NewRepository() Repository {
	return &osRepository{}
}

type osRepository struct{}

func (r *osRepository) GetEnv(envKey string) (string, error) {
	key := strings.TrimSpace(envKey)
	if key == "" {
		return "", fmt.Errorf("environment variable key: %q must not be empty", envKey)
	}

	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return "", MissingEnvError{Key: key}
	}

	return value, nil
}
