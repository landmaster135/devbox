package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	result, err := service.CraftMarkdown("content", 10, 10, srcJSONFile, srcBodyDir, outDir)
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
	_, err := service.CraftMarkdown("content", 10, 10, srcJSONFile, srcBodyDir, outDir)
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
	_, err := service.CraftMarkdown("content", 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err == nil || !strings.Contains(err.Error(), "## Frequent Tags") {
		t.Fatalf("error = %v", err)
	}
}

func TestService_CraftMarkdown_MissingSourceFileIsSkipped(t *testing.T) {
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
	result, err := service.CraftMarkdown("content", 10, 11, srcJSONFile, srcBodyDir, outDir)
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
		t.Fatalf("missing source should be skipped")
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
			got, err := mapCategoryTag(tt.input)
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
		got, err := mapOwningStatusTag(tt.input)
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
		got, err := mapColorTag(tt.input)
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

	number, err := parseConNumber("CO000012345")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if number != 12345 {
		t.Fatalf("number = %d, want 12345", number)
	}

	_, err = parseConNumber("invalid")
	if err == nil {
		t.Fatalf("expected error")
	}
}
