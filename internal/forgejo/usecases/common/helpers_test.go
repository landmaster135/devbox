package common

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestDecodeProjects(t *testing.T) {
	plain := `[{"name":"P1","created_at":"2020-01-01T00:00:00Z","updated_at":"2020-01-02T00:00:00Z"}]`
	decoded, err := DecodeProjects([]byte(plain))
	if err != nil {
		t.Fatalf("DecodeProjects() error = %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("len(decoded) = %d, want 1", len(decoded))
	}

	wrapped := `{"data":[{"name":"P2","title":"Title","created_at":"2020-01-03T00:00:00Z","updated_at":"2020-01-04T00:00:00Z"}], "projects":[]}`
	decoded, err = DecodeProjects([]byte(wrapped))
	if err != nil {
		t.Fatalf("DecodeProjects() wrapped error = %v", err)
	}
	if len(decoded) != 1 {
		t.Fatalf("len(decoded) = %d, want 1", len(decoded))
	}
	if decoded[0].Name != "P2" {
		t.Fatalf("Name = %q, want %q", decoded[0].Name, "P2")
	}
}

func TestIsNotFoundError(t *testing.T) {
	if ok := IsNotFoundError(&RequestError{Status: http.StatusNotFound, Body: "not found"}); !ok {
		t.Fatal("IsNotFoundError() should return true")
	}
	if ok := IsNotFoundError(nil); ok {
		t.Fatal("IsNotFoundError(nil) should return false")
	}
	if ok := IsNotFoundError(fmt.Errorf("other")); ok {
		t.Fatal("IsNotFoundError(other) should return false")
	}
}

func TestFormatDate(t *testing.T) {
	if got := FormatDate(time.Time{}); got != "" {
		t.Fatalf("FormatDate(zero) = %q, want %q", got, "")
	}
	if got := FormatDate(time.Date(2026, 5, 9, 12, 34, 56, 0, time.UTC)); got != "2026-05-09T12:34:56Z" {
		t.Fatalf("FormatDate() = %q, want %q", got, "2026-05-09T12:34:56Z")
	}
}

func TestPrimaryLanguage(t *testing.T) {
	if got := PrimaryLanguage(map[string]float64{"A": 1, "B": 3, "C": 2}); got != "B" {
		t.Fatalf("PrimaryLanguage() = %q, want %q", got, "B")
	}
	if got := PrimaryLanguage(map[string]float64{}); got != "" {
		t.Fatalf("PrimaryLanguage(empty) = %q, want %q", got, "")
	}
}
