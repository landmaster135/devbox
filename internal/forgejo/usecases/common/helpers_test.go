package common

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestFormatDatePtr(t *testing.T) {
	now := time.Date(2026, 5, 9, 12, 34, 56, 0, time.UTC)
	if got := FormatDatePtr(nil); got != "" {
		t.Fatalf("FormatDatePtr(nil) = %q, want %q", got, "")
	}
	if got := FormatDatePtr(&now); got != "2026-05-09T12:34:56Z" {
		t.Fatalf("FormatDatePtr() = %q, want %q", got, "2026-05-09T12:34:56Z")
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
