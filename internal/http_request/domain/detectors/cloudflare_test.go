package detectors

import (
	"strings"
	"testing"
)

func TestIsCloudflareChallenge(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		headers  map[string]string
		body     []byte
		expected bool
	}{
		{
			name:   "Cloudflare403",
			status: 403,
			headers: map[string]string{
				"Server":       "cloudflare",
				"Cf-Mitigated": "challenge",
			},
			body:     []byte("<title>Just a moment...</title>"),
			expected: true,
		},
		{
			name:   "NonCloudflare403",
			status: 403,
			headers: map[string]string{
				"Server": "nginx",
			},
			body:     []byte("forbidden"),
			expected: false,
		},
		{
			name:   "Retryable503",
			status: 503,
			headers: map[string]string{
				"Server": "cloudflare",
			},
			body:     []byte("Checking your browser before accessing"),
			expected: true,
		},
		{
			name:   "BenignResponse",
			status: 200,
			headers: map[string]string{
				"Server": "cloudflare",
			},
			body:     []byte("Hello"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCloudflareChallenge(tt.status, tt.headers, tt.body)
			if got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestBuildCloudflareWarning(t *testing.T) {
	headers := map[string]string{
		"Server":       "cloudflare",
		"Cf-Ray":       "12345",
		"Cf-Mitigated": "challenge",
	}

	warning := BuildCloudflareWarning(403, headers)
	if warning == "" {
		t.Fatal("warning should not be empty")
	}
	if !contains(warning, "Ray ID: 12345") {
		t.Fatalf("warning should contain Ray ID, got %s", warning)
	}
}

func contains(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
