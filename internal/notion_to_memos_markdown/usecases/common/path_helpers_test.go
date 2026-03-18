package common

import (
	"testing"

	domain "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/domain"
)

func TestBuildJSONConIDSet(t *testing.T) {
	t.Parallel()

	contents := []domain.Content{
		{ConID: " CO0001 "},
		{ConID: ""},
		{ConID: "   "},
		{ConID: "CO0001"},
		{ConID: "CO0002"},
	}

	got := BuildJSONConIDSet(contents)
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if _, exists := got["CO0001"]; !exists {
		t.Fatalf("CO0001 should exist: %#v", got)
	}
	if _, exists := got["CO0002"]; !exists {
		t.Fatalf("CO0002 should exist: %#v", got)
	}
}

func TestCountSrcBodyMetrics(t *testing.T) {
	t.Parallel()

	srcBodyFiles := []string{
		"/tmp/CO0001.md",
		"/tmp/CO0002.md",
		"/tmp/memo.txt",
		"/tmp/ CO0003 .MD",
	}
	jsonConIDSet := map[string]struct{}{
		"CO0001": {},
		"CO0003": {},
	}

	total, mapped, unmapped := CountSrcBodyMetrics(srcBodyFiles, jsonConIDSet)
	if total != 4 {
		t.Fatalf("total = %d, want 4", total)
	}
	if mapped != 2 {
		t.Fatalf("mapped = %d, want 2", mapped)
	}
	if unmapped != 2 {
		t.Fatalf("unmapped = %d, want 2", unmapped)
	}
}

func TestExtractConIDFromPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "normal md", path: "/tmp/CO0001.md", want: "CO0001"},
		{name: "upper extension", path: "/tmp/CO0001.MD", want: "CO0001"},
		{name: "trim filename", path: "/tmp/ CO0002 .md", want: "CO0002"},
		{name: "non md", path: "/tmp/CO0001.txt", want: ""},
		{name: "empty basename", path: "/tmp/.md", want: ""},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ExtractConIDFromPath(tt.path)
			if got != tt.want {
				t.Fatalf("ExtractConIDFromPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestSanitizeCategory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		category string
		want     string
	}{
		{name: "empty", category: "", want: DefaultCategory},
		{name: "trim", category: " software ", want: "software"},
		{name: "slash", category: "tools/dev", want: "tools_dev"},
		{name: "backslash", category: `tools\dev`, want: "tools_dev"},
		{name: "traversal", category: "../tools/./dev//", want: "tools_dev"},
		{name: "only separators", category: " / / ", want: DefaultCategory},
		{name: "null rune", category: "ab\x00cd", want: "ab_cd"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := SanitizeCategory(tt.category)
			if got != tt.want {
				t.Fatalf("SanitizeCategory(%q) = %q, want %q", tt.category, got, tt.want)
			}
		})
	}
}

func TestBuildSplitMarkdownIndex(t *testing.T) {
	t.Parallel()

	index := BuildSplitMarkdownIndex([]string{
		"/tmp/CO0001_02.md",
		"/tmp/CO0001_01.md",
		"/tmp/CO0001.md",
		"/tmp/CO0002_a.md",
		"/tmp/CO0002_.md",
		"/tmp/CO0003_01.txt",
	})

	resolved := index.Resolve("CO0001")
	if len(resolved) != 2 {
		t.Fatalf("len(resolved) = %d, want 2", len(resolved))
	}
	if resolved[0] != "/tmp/CO0001_01.md" || resolved[1] != "/tmp/CO0001_02.md" {
		t.Fatalf("unexpected resolved order: %#v", resolved)
	}

	co2Resolved := index.Resolve(" CO0002 ")
	if len(co2Resolved) != 1 || co2Resolved[0] != "/tmp/CO0002_a.md" {
		t.Fatalf("unexpected CO0002 resolved: %#v", co2Resolved)
	}

	missing := index.Resolve("CO9999")
	if len(missing) != 0 {
		t.Fatalf("missing should be empty: %#v", missing)
	}

	resolved[0] = "mutated"
	resolvedAgain := index.Resolve("CO0001")
	if resolvedAgain[0] != "/tmp/CO0001_01.md" {
		t.Fatalf("Resolve should return a defensive copy: %#v", resolvedAgain)
	}
}
