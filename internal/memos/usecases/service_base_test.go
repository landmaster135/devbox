package usecases

import "testing"

func TestNewService_DefaultClient_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		BaseURL:  "https://memos.example.com",
		APIToken: "token",
		Timeout:  0,
	})
	if service.client == nil {
		t.Fatal("client is nil")
	}
	if service.baseURL != "https://memos.example.com/api/v1" {
		t.Fatalf("baseURL = %s, want https://memos.example.com/api/v1", service.baseURL)
	}
	if service.fileSystem == nil {
		t.Fatal("fileSystem is nil")
	}
}

func TestNormalizeBaseURL_Normal(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "https://example.com", want: "https://example.com/api/v1"},
		{name: "with trailing slash", in: "https://example.com/", want: "https://example.com/api/v1"},
		{name: "with api v1", in: "https://example.com/api/v1", want: "https://example.com/api/v1"},
		{name: "empty", in: "", want: "/api/v1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeBaseURL(tc.in)
			if got != tc.want {
				t.Fatalf("normalizeBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestNormalizeMemoIdentifier_Normal(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "plain", in: "memo-1", want: "memo-1"},
		{name: "with memos prefix", in: "memos/memo-1", want: "memo-1"},
		{name: "with api prefix", in: "api/v1/memos/memo-1", want: "memo-1"},
		{name: "with full url", in: "https://example.com/api/v1/memos/memo-1", want: "memo-1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeMemoIdentifier(tc.in)
			if got != tc.want {
				t.Fatalf("normalizeMemoIdentifier(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestBuildMemoResourceName_Normal(t *testing.T) {
	if got := buildMemoResourceName("https://example.com/api/v1/memos/memo-1"); got != "memos/memo-1" {
		t.Fatalf("buildMemoResourceName() = %s, want memos/memo-1", got)
	}
	if got := buildMemoResourceName("  "); got != "" {
		t.Fatalf("buildMemoResourceName() = %q, want empty", got)
	}
}
