package common

import (
	"strings"
	"testing"
)

func TestNormalizeNewlines_Normal(t *testing.T) {
	t.Parallel()

	got := NormalizeNewlines("a\r\nb\r\n")
	if got != "a\nb\n" {
		t.Fatalf("unexpected normalized content: %q", got)
	}
}

func TestSplitFrontMatterBlock_NoFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	hasFrontMatter, block, body, err := SplitFrontMatterBlock("hello")
	if err != nil {
		t.Fatalf("SplitFrontMatterBlock returned error: %v", err)
	}
	if hasFrontMatter {
		t.Fatal("expected hasFrontMatter=false")
	}
	if block != "" {
		t.Fatalf("expected empty block, got %q", block)
	}
	if body != "hello" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestSplitFrontMatterBlock_WithBody_Normal(t *testing.T) {
	t.Parallel()

	content := "---\ntitle: sample\n---\n\nbody\n"
	hasFrontMatter, block, body, err := SplitFrontMatterBlock(content)
	if err != nil {
		t.Fatalf("SplitFrontMatterBlock returned error: %v", err)
	}
	if !hasFrontMatter {
		t.Fatal("expected hasFrontMatter=true")
	}
	expectedBlock := "---\ntitle: sample\n---\n"
	if block != expectedBlock {
		t.Fatalf("unexpected block: %q", block)
	}
	if body != "\nbody\n" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestSplitFrontMatterBlock_OnlyFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	content := "---\ntitle: sample\n---"
	hasFrontMatter, block, body, err := SplitFrontMatterBlock(content)
	if err != nil {
		t.Fatalf("SplitFrontMatterBlock returned error: %v", err)
	}
	if !hasFrontMatter {
		t.Fatal("expected hasFrontMatter=true")
	}
	expectedBlock := "---\ntitle: sample\n---\n"
	if block != expectedBlock {
		t.Fatalf("unexpected block: %q", block)
	}
	if body != "" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestSplitFrontMatterBlock_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, _, _, err := SplitFrontMatterBlock("---\ntitle: sample\n")
	if err == nil {
		t.Fatal("expected error for missing front matter end delimiter")
	}
}

func TestParseFrontMatterMap_Empty_Normal(t *testing.T) {
	t.Parallel()

	keys, values, err := ParseFrontMatterMap("---\n\n---\n")
	if err != nil {
		t.Fatalf("ParseFrontMatterMap returned error: %v", err)
	}
	if len(keys) != 0 || len(values) != 0 {
		t.Fatalf("expected empty map, got keys=%v values=%v", keys, values)
	}
}

func TestParseFrontMatterMap_Valid_Normal(t *testing.T) {
	t.Parallel()

	keys, values, err := ParseFrontMatterMap("---\ntitle: sample\nstatus: draft\n---\n")
	if err != nil {
		t.Fatalf("ParseFrontMatterMap returned error: %v", err)
	}
	if strings.Join(keys, ",") != "title,status" {
		t.Fatalf("unexpected keys: %v", keys)
	}
	if values["title"] != "sample" || values["status"] != "draft" {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestParseFrontMatterMap_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, _, err := ParseFrontMatterMap("---\ninvalid-line\n---\n")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseFrontMatterMap_EmptyKey(t *testing.T) {
	t.Parallel()

	_, _, err := ParseFrontMatterMap("---\n: value\n---\n")
	if err == nil {
		t.Fatal("expected empty key error")
	}
}

func TestParseKVPairs_Valid_Normal(t *testing.T) {
	t.Parallel()

	keys, values, err := ParseKVPairs([]string{"title=sample", "status=draft", "title=updated"})
	if err != nil {
		t.Fatalf("ParseKVPairs returned error: %v", err)
	}
	if strings.Join(keys, ",") != "title,status" {
		t.Fatalf("unexpected keys: %v", keys)
	}
	if values["title"] != "updated" || values["status"] != "draft" {
		t.Fatalf("unexpected values: %v", values)
	}
}

func TestParseKVPairs_InvalidFormat(t *testing.T) {
	t.Parallel()

	_, _, err := ParseKVPairs([]string{"invalid"})
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestParseKVPairs_EmptyKey(t *testing.T) {
	t.Parallel()

	_, _, err := ParseKVPairs([]string{"=value"})
	if err == nil {
		t.Fatal("expected empty key error")
	}
}

func TestBuildFrontMatter_Normal(t *testing.T) {
	t.Parallel()

	got := BuildFrontMatter([]string{"title", "status"}, map[string]string{
		"title":  "sample",
		"status": "draft",
	})
	expected := "---\ntitle: sample\nstatus: draft\n---\n"
	if got != expected {
		t.Fatalf("unexpected front matter:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestUniqueTrimmedTags_Normal(t *testing.T) {
	t.Parallel()

	tags := UniqueTrimmedTags(" go, #go,markdown , ,#memo ")
	if strings.Join(tags, ",") != "go,markdown,memo" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestBuildTagLine_Normal(t *testing.T) {
	t.Parallel()

	got, err := BuildTagLine([]string{"go", "markdown"})
	if err != nil {
		t.Fatalf("BuildTagLine returned error: %v", err)
	}
	if got != "#go #markdown" {
		t.Fatalf("unexpected tag line: %q", got)
	}
}

func TestBuildTagLine_Invalid(t *testing.T) {
	t.Parallel()

	_, err := BuildTagLine(nil)
	if err == nil {
		t.Fatal("expected error for empty tags")
	}
}
