package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// NewRepository returns the default filesystem-backed repository.
func NewRepository() Repository {
	return &osRepository{}
}

type osRepository struct{}

func (r *osRepository) Write(path string, overwrites bool, content string) error {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return fmt.Errorf("path must not be empty")
	}

	cleanPath := filepath.Clean(trimmed)
	dir := filepath.Dir(cleanPath)
	if dir != "." && dir != "" {
		if err := r.EnsureDir(dir); err != nil {
			return fmt.Errorf("ensure parent directory for %s: %w", cleanPath, err)
		}
	}

	flags := os.O_CREATE | os.O_WRONLY
	if overwrites {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_APPEND
	}

	f, err := os.OpenFile(cleanPath, flags, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", cleanPath, err)
	}
	defer f.Close()

	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("write %s: %w", cleanPath, err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync %s: %w", cleanPath, err)
	}

	return nil
}

func (r *osRepository) EnsureDir(path string) error {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return fmt.Errorf("directory path must not be empty")
	}

	cleanPath := filepath.Clean(trimmed)
	if err := os.MkdirAll(cleanPath, 0o755); err != nil {
		return fmt.Errorf("create directory %s: %w", cleanPath, err)
	}

	return nil
}

var defaultRepository = NewRepository()

// WriteFile is a convenience wrapper around the default Repository for callers
// that do not need dependency injection.
func WriteFile(path string, overwrites bool, content string) error {
	return defaultRepository.Write(path, overwrites, content)
}

// EnsureDir ensures that the specified directory exists using the default
// repository implementation.
func EnsureDir(path string) error {
	return defaultRepository.EnsureDir(path)
}
