package usecases

import (
	"errors"
	"testing"
)

func TestTruncateForLog(t *testing.T) {
	short := []byte("hello")
	if got := truncateForLog(short); got != "hello" {
		t.Fatalf("unexpected short result: %s", got)
	}

	long := make([]byte, 600)
	for i := range long {
		long[i] = 'a'
	}
	got := truncateForLog(long)
	if len(got) != 515 || got[len(got)-3:] != "..." {
		t.Fatalf("unexpected long result: %q", got[len(got)-10:])
	}
}

func TestSetBaseURL(t *testing.T) {
	svc := NewHealthService(0)
	original := svc.baseURL
	svc.SetBaseURL(" https://custom.example/api/ ")
	if svc.baseURL != "https://custom.example/api" {
		t.Fatalf("baseURL not trimmed: %s", svc.baseURL)
	}
	svc.SetBaseURL("")
	if svc.baseURL != "https://custom.example/api" {
		t.Fatalf("baseURL should remain unchanged when empty")
	}

	// ensure default was set initially
	if original == "" {
		t.Fatalf("default baseURL should not be empty")
	}
}

func TestHealthServiceShouldRetryDailySummaryWithRefresh(t *testing.T) {
	svc := NewHealthService(0)

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "401", err: errors.New("status=401 unauthorized"), want: true},
		{name: "403", err: errors.New("Status=403"), want: true},
		{name: "unauthorized", err: errors.New("UNAUTHORIZED"), want: true},
		{name: "invalid token", err: errors.New("INVALID_TOKEN"), want: true},
		{name: "other", err: errors.New("timeout"), want: false},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			got := svc.ShouldRetryDailySummaryWithRefresh(tc.err)
			if got != tc.want {
				t.Fatalf("ShouldRetryDailySummaryWithRefresh(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestLabelForMeasureType(t *testing.T) {
	if got := labelForMeasureType(1); got != "weight_kg" {
		t.Fatalf("unexpected known label: %s", got)
	}
	if got := labelForMeasureType(999); got != "measure_999" {
		t.Fatalf("unexpected fallback label: %s", got)
	}
}

func TestFlexibleBoolUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		want  bool
		valid bool
	}{
		{name: "true", data: "true", want: true, valid: true},
		{name: "false", data: "false", want: false, valid: true},
		{name: "one", data: "1", want: true, valid: true},
		{name: "zero", data: "0", want: false, valid: true},
		{name: "whitespace number", data: " 5 ", want: true, valid: true},
		{name: "empty", data: "", want: false, valid: true},
		{name: "invalid", data: "\"yes\"", valid: false},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			var fb flexibleBool
			err := fb.UnmarshalJSON([]byte(tc.data))
			if !tc.valid {
				if err == nil {
					t.Fatalf("expected error for %s", tc.data)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fb.Bool() != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, fb.Bool())
			}
		})
	}
}
