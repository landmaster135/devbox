package common

import (
	"fmt"
	"strings"
	"testing"
)

type mockContentReader struct {
	readFileFunc func(filePath string) ([]byte, error)
}

func (m *mockContentReader) ReadFile(filePath string) ([]byte, error) {
	if m.readFileFunc != nil {
		return m.readFileFunc(filePath)
	}
	return nil, nil
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
			got := NormalizeBaseURL(tc.in)
			if got != tc.want {
				t.Fatalf("NormalizeBaseURL(%q) = %q, want %q", tc.in, got, tc.want)
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
			got := NormalizeMemoIdentifier(tc.in)
			if got != tc.want {
				t.Fatalf("NormalizeMemoIdentifier(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestBuildMemoResourceName_Normal(t *testing.T) {
	if got := BuildMemoResourceName("https://example.com/api/v1/memos/memo-1"); got != "memos/memo-1" {
		t.Fatalf("BuildMemoResourceName() = %s, want memos/memo-1", got)
	}
	if got := BuildMemoResourceName("  "); got != "" {
		t.Fatalf("BuildMemoResourceName() = %q, want empty", got)
	}
}

func TestResolveContent_ContentFileReadError_Error(t *testing.T) {
	reader := &mockContentReader{
		readFileFunc: func(filePath string) ([]byte, error) {
			return nil, fmt.Errorf("read failed")
		},
	}

	_, err := ResolveContent("", "./memo.md", reader)
	if err == nil {
		t.Fatal("ResolveContent() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "content-file の読み込みに失敗しました") {
		t.Fatalf("error = %v, want content-file の読み込みに失敗しました", err)
	}
}
