package env

import (
	"errors"
	"testing"
)

func TestOsRepositoryGetEnv(t *testing.T) {
	t.Run("empty key returns error", func(t *testing.T) {
		repo := NewRepository()

		const whitespaceKey = "   "
		_, err := repo.GetEnv(whitespaceKey)
		if err == nil {
			t.Fatalf("expected error for whitespace-only key")
		}

		expected := "environment variable key: \"   \" must not be empty"
		if err.Error() != expected {
			t.Fatalf("unexpected error message. want %q, got %q", expected, err.Error())
		}
	})

	t.Run("missing environment variable returns MissingEnvError", func(t *testing.T) {
		repo := NewRepository()

		const key = "CRON_WORKFLOW_TEST_MISSING"
		t.Setenv(key, "")

		_, err := repo.GetEnv(key)
		if err == nil {
			t.Fatalf("expected MissingEnvError when env is unset")
		}

		var missing MissingEnvError
		if !errors.As(err, &missing) {
			t.Fatalf("expected MissingEnvError, got %T", err)
		}

		if missing.Key != key {
			t.Fatalf("MissingEnvError.Key mismatch. want %q, got %q", key, missing.Key)
		}
	})

	t.Run("returns trimmed value when key contains surrounding spaces", func(t *testing.T) {
		repo := NewRepository()

		const (
			actualKey   = "CRON_WORKFLOW_TEST_VALUE"
			providedKey = "  CRON_WORKFLOW_TEST_VALUE  "
			rawValue    = "   ready-to-run   "
			wantValue   = "ready-to-run"
		)

		t.Setenv(actualKey, rawValue)

		got, err := repo.GetEnv(providedKey)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != wantValue {
			t.Fatalf("unexpected value. want %q, got %q", wantValue, got)
		}
	})
}
