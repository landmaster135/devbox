package usecases

import "testing"

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
