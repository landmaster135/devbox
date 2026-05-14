package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestSplitByComma_Normal(t *testing.T) {
	got := splitByComma(" content, visibility ,state ")
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	if got[0] != "content" || got[1] != "visibility" || got[2] != "state" {
		t.Fatalf("split result = %v, want [content visibility state]", got)
	}
}

func TestSplitByComma_Empty_Normal(t *testing.T) {
	got := splitByComma("  ")
	if got != nil {
		t.Fatalf("got = %v, want nil", got)
	}
}

func TestBuildAnyTagsFilter_Normal(t *testing.T) {
	got := buildAnyTagsFilter([]string{"health", "book"})
	want := "tag in ['health','book']"
	if got != want {
		t.Fatalf("buildAnyTagsFilter() = %q, want %q", got, want)
	}
}

func TestBuildAnyTagsFilter_Escape_Normal(t *testing.T) {
	got := buildAnyTagsFilter([]string{`a'b`, `path\tag`})
	want := `tag in ['a\'b','path\\tag']`
	if got != want {
		t.Fatalf("buildAnyTagsFilter() = %q, want %q", got, want)
	}
}

func TestMergeFilters_Normal(t *testing.T) {
	got := mergeFilters(`visibility == "PUBLIC"`, "tag in ['health','book']")
	want := `(visibility == "PUBLIC") && (tag in ['health','book'])`
	if got != want {
		t.Fatalf("mergeFilters() = %q, want %q", got, want)
	}
}

func TestMergeFilters_ExtraEmpty_Normal(t *testing.T) {
	got := mergeFilters(`visibility == "PUBLIC"`, "")
	if got != `visibility == "PUBLIC"` {
		t.Fatalf("mergeFilters() = %q, want base filter", got)
	}
}

func TestBoolPointer_Normal(t *testing.T) {
	if got := boolPointer(true, false); got != nil {
		t.Fatalf("got = %v, want nil", got)
	}
	got := boolPointer(false, true)
	if got == nil || *got != false {
		t.Fatalf("got = %v, want pointer to false", got)
	}
}

func TestPrintJSON_Normal(t *testing.T) {
	var out bytes.Buffer
	if err := printJSON(&out, map[string]string{"hello": "world"}); err != nil {
		t.Fatalf("printJSON() error = %v", err)
	}
	if !strings.Contains(out.String(), "\"hello\": \"world\"") {
		t.Fatalf("output = %s, want JSON content", out.String())
	}
}

func TestPrintJSON_MarshalError_Error(t *testing.T) {
	var out bytes.Buffer
	err := printJSON(&out, map[string]any{"invalid": make(chan int)})
	if err == nil {
		t.Fatal("printJSON() error = nil, want error")
	}
}
