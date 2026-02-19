package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	common"github.com/landmaster135/devbox/internal/notion_to_memos_markdown/usecases/common"
)

func TestService_CraftMarkdown_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "Google Apps Script",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"url": "https://example.com/gas",
			"tags": [
				{"page_title": "JavaScript"},
				{"page_title": "CustomTag"}
			]
		},
		{
			"con_id": "CO000000011",
			"page_title": "Out of range",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 0,
			"price": 0,
			"tags": []
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}

	tagsMarkdown := `# tags

## Frequent Tags
#31-programming/language/javascript #31-programming/language/golang
`
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}

	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000011.md"), []byte("範囲外です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	result, err := service.CraftMarkdown("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	craftedPath := filepath.Join(outDir, "CO000000010.md")
	craftedData, err := os.ReadFile(craftedPath)
	if err != nil {
		t.Fatalf("read crafted file failed: %v", err)
	}
	crafted := string(craftedData)
	if !strings.Contains(crafted, "bought_at: 2024-01-01T09:00:00+09:00") {
		t.Fatalf("front matter missing: %s", crafted)
	}
	if !strings.Contains(crafted, "score_of_100: 88") || !strings.Contains(crafted, "price_yen: 1200") {
		t.Fatalf("front matter values missing: %s", crafted)
	}
	if !strings.Contains(crafted, "url: https://example.com/gas") {
		t.Fatalf("front matter url missing: %s", crafted)
	}
	if !strings.Contains(crafted, "# Google Apps Script") {
		t.Fatalf("heading missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#31-programming/language/javascript") {
		t.Fatalf("frequent tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#customtag") {
		t.Fatalf("fallback tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#91-backup/tool-migration/202602-notion") {
		t.Fatalf("backup tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#0a-content/software") ||
		!strings.Contains(crafted, "#01-p/own-status/2-already") ||
		!strings.Contains(crafted, "#01-p/color/gray") {
		t.Fatalf("classification tags missing: %s", crafted)
	}

	idxCategory := strings.Index(crafted, "#0a-content/software")
	idxOwningStatus := strings.Index(crafted, "#01-p/own-status/2-already")
	idxColor := strings.Index(crafted, "#01-p/color/gray")
	idxBackup := strings.Index(crafted, "#91-backup/tool-migration/202602-notion")
	idxPageTag1 := strings.Index(crafted, "#31-programming/language/javascript")
	idxPageTag2 := strings.Index(crafted, "#customtag")
	if !(idxCategory < idxOwningStatus &&
		idxOwningStatus < idxColor &&
		idxBackup < idxCategory &&
		idxColor < idxPageTag1 &&
		idxPageTag1 < idxPageTag2) {
		t.Fatalf("tag order is invalid: %s", crafted)
	}

	if _, err := os.Stat(filepath.Join(outDir, "CO000000011.md")); !os.IsNotExist(err) {
		t.Fatalf("out-of-range file should not be created")
	}
}

func TestService_CraftMarkdown_InvalidCategory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "Invalid category sample",
			"category": "unknown-category",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}

	tagsMarkdown := "## Frequent Tags\n#31-programming/language/javascript\n"
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	_, err := service.CraftMarkdown("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err == nil || !strings.Contains(err.Error(), "未対応のcategoryです") {
		t.Fatalf("error = %v", err)
	}
}

func TestService_CraftMarkdown_FrequentTagsNotFound(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "sample",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(tagsFile, []byte("# no section\n"), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	_, err := service.CraftMarkdown("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err == nil || !strings.Contains(err.Error(), "## Frequent Tags") {
		t.Fatalf("error = %v", err)
	}
}

func TestService_CraftMarkdown_MissingSourceFileCreatesEmptyMarkdown(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "Exists",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		},
		{
			"con_id": "CO000000011",
			"page_title": "Missing",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	tagsMarkdown := "## Frequent Tags\n#31-programming/language/javascript\n"
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	result, err := service.CraftMarkdown("content", "", false, 10, 11, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=2") || !strings.Contains(result, "加工成功=2") {
		t.Fatalf("unexpected result: %s", result)
	}

	if _, err := os.Stat(filepath.Join(outDir, "CO000000010.md")); err != nil {
		t.Fatalf("existing source should be crafted: %v", err)
	}
	missingPath := filepath.Join(outDir, "CO000000011.md")
	missingData, err := os.ReadFile(missingPath)
	if err != nil {
		t.Fatalf("missing source output should exist: %v", err)
	}
	missingContent := string(missingData)
	if !strings.Contains(missingContent, "con_id: CO000000011") {
		t.Fatalf("front matter should be applied to empty markdown: %s", missingContent)
	}
	if !strings.Contains(missingContent, "# Missing") {
		t.Fatalf("heading should be applied to empty markdown: %s", missingContent)
	}
	if !strings.Contains(missingContent, "#91-backup/tool-migration/202602-notion") {
		t.Fatalf("tags should be applied to empty markdown: %s", missingContent)
	}
}

func TestService_CraftMarkdown_MissingSourceFileCanSkipByFlag(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "Exists",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		},
		{
			"con_id": "CO000000011",
			"page_title": "Missing",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	tagsMarkdown := "## Frequent Tags\n#31-programming/language/javascript\n"
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	result, err := service.CraftMarkdown("content", "", true, 10, 11, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=2") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	if _, err := os.Stat(filepath.Join(outDir, "CO000000010.md")); err != nil {
		t.Fatalf("existing source should be crafted: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "CO000000011.md")); !os.IsNotExist(err) {
		t.Fatalf("missing source should be skipped when --skips-no-src-body=true")
	}
}

func TestService_CraftMarkdown_FilterByCategory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "Software item",
			"category": "software",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		},
		{
			"con_id": "CO000000011",
			"page_title": "Book item",
			"category": "book",
			"owning_status": "already",
			"color": "gray",
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	tagsMarkdown := "## Frequent Tags\n#31-programming/language/javascript\n"
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000011.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	result, err := service.CraftMarkdown("content", "software", false, 10, 11, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	if _, err := os.Stat(filepath.Join(outDir, "CO000000010.md")); err != nil {
		t.Fatalf("software file should be output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "CO000000011.md")); !os.IsNotExist(err) {
		t.Fatalf("book file should be filtered out")
	}
}

func TestService_CraftMarkdown_NullColorSkipsColorTag(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "CO000000010",
			"page_title": "No color content",
			"category": "software",
			"owning_status": "already",
			"color": null,
			"bought_at": "2024-01-01T09:00:00+09:00",
			"score": 88,
			"price": 1200,
			"tags": [{"page_title": "JavaScript"}]
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	tagsMarkdown := "## Frequent Tags\n#31-programming/language/javascript\n"
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(nil)
	result, err := service.CraftMarkdown("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	outPath := filepath.Join(outDir, "CO000000010.md")
	outData, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output failed: %v", err)
	}
	outText := string(outData)
	if strings.Contains(outText, "#01-p/color/") {
		t.Fatalf("color tag should not be added: %s", outText)
	}
	if !strings.Contains(outText, `url: ""`) {
		t.Fatalf("empty url should be explicit empty string: %s", outText)
	}
	if !strings.Contains(outText, "#0a-content/software") || !strings.Contains(outText, "#01-p/own-status/2-already") {
		t.Fatalf("required tags missing: %s", outText)
	}
}

func TestMapCategoryTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		want     string
		wantFail bool
	}{
		{name: "web-clip", input: "web-clip", want: "0a-content/web-clip"},
		{name: "subscription_year", input: "subscription_year", want: "0a-content/subscription-year"},
		{name: "subscription month", input: "subscription month", want: "0a-content/subscription-month"},
		{name: "unknown", input: "unknown", wantFail: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := common.MapCategoryTag(tt.input)
			if tt.wantFail {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMapOwningStatusTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		want     string
		wantFail bool
	}{
		{input: "yet", want: "01-p/own-status/1-yet"},
		{input: "already", want: "01-p/own-status/2-already"},
		{input: "gone", want: "01-p/own-status/3-gone"},
		{input: "bad", wantFail: true},
	}

	for _, tt := range tests {
		got, err := common.MapOwningStatusTag(tt.input)
		if tt.wantFail {
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("got %q, want %q", got, tt.want)
		}
	}
}

func TestMapColorTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		want     string
		wantFail bool
	}{
		{input: "gray", want: "01-p/color/gray"},
		{input: "white", want: "01-p/color/white"},
		{input: "purple", want: "01-p/color/purple"},
		{input: "unknown", wantFail: true},
	}

	for _, tt := range tests {
		got, err := common.MapColorTag(tt.input)
		if tt.wantFail {
			if err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tt.input, err)
		}
		if got != tt.want {
			t.Fatalf("got %q, want %q", got, tt.want)
		}
	}
}

func TestParseConNumber(t *testing.T) {
	t.Parallel()

	number, err := common.ParseConNumber("CO000012345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if number != 12345 {
		t.Fatalf("number = %d, want 12345", number)
	}

	_, err = common.ParseConNumber("invalid")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestResolveFrequentTag(t *testing.T) {
	t.Parallel()

	frequentTags := []string{
		"chromeextension",
		"31-programming/language/javascript",
		"31-programming/language/typescript",
	}

	tests := []struct {
		name string
		tag  string
		want string
	}{
		{
			name: "no slash exact match",
			tag:  "chromeextension",
			want: "chromeextension",
		},
		{
			name: "no slash partial not match",
			tag:  "chrome",
			want: "chrome",
		},
		{
			name: "with slash match by suffix exact",
			tag:  "javascript",
			want: "31-programming/language/javascript",
		},
		{
			name: "with slash partial not match",
			tag:  "java",
			want: "java",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := common.ResolveFrequentTag(tt.tag, frequentTags)
			if got != tt.want {
				t.Fatalf("resolveFrequentTag(%q) = %q, want %q", tt.tag, got, tt.want)
			}
		})
	}
}

func TestNormalizeFrontMatterURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "normal", input: "https://example.com", want: "https://example.com"},
		{name: "empty", input: "", want: `""`},
		{name: "spaces", input: "   ", want: `""`},
		{name: "trim", input: "  https://example.com/path  ", want: "https://example.com/path"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := common.NormalizeFrontMatterURL(tt.input)
			if got != tt.want {
				t.Fatalf("normalizeFrontMatterURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
