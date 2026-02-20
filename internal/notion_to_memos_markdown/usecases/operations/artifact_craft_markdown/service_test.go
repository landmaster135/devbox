package artifactcraftmarkdown

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	filesystem "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/filesystem"
)

func TestService_Execute_Normal(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "artifacts.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	jsonContent := `[
		{
			"con_id": "AF0186",
			"page_title": "my devbox release note",
			"category": "system",
			"output_url": "https://example.com/devbox",
			"tags": [
				{"page_title": "Golang"},
				{"page_title": "CustomTag"}
			]
		},
		{
			"con_id": "AF0187",
			"page_title": "out of range",
			"category": "system",
			"output_url": "https://example.com/skip",
			"tags": []
		}
	]`
	if err := os.WriteFile(srcJSONFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}

	tagsMarkdown := `## Artifact
#06-af/life #06-af/movie #06-af/product
#06-af/system/devbox #06-af/system/others

## Frequent Tags
#31-programming/language/golang #openai
`
	if err := os.WriteFile(tagsFile, []byte(tagsMarkdown), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "AF0186.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("artifact", "", false, 186, 186, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	craftedData, err := os.ReadFile(filepath.Join(outDir, "AF0186.md"))
	if err != nil {
		t.Fatalf("read crafted file failed: %v", err)
	}
	crafted := string(craftedData)
	if !strings.Contains(crafted, "# my devbox release note") {
		t.Fatalf("heading missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#31-programming/language/golang") {
		t.Fatalf("frequent tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#customtag") {
		t.Fatalf("fallback tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#06-af/system/devbox") {
		t.Fatalf("system category tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#91-backup/tool-migration/202602-notion") {
		t.Fatalf("required backup tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "https://example.com/devbox") {
		t.Fatalf("output_url missing: %s", crafted)
	}
	if _, err := os.Stat(filepath.Join(outDir, "AF0187.md")); !os.IsNotExist(err) {
		t.Fatalf("out-of-range file should not be created")
	}
}

func TestService_Execute_CategoryMappings(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "artifacts.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"AF0100",
			"page_title":"note活動（投稿部）",
			"category":"article",
			"output_url":"https://example.com/note",
			"tags":[]
		},
		{
			"con_id":"AF0101",
			"page_title":"movie summary",
			"category":"movie",
			"output_url":"https://example.com/movie",
			"tags":[]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(tagsFile, []byte(`## Artifact
#06-af/life #06-af/movie #06-af/product

## Frequent Tags
#31-programming/language/golang
`), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "AF0100.md"), []byte("body"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "AF0101.md"), []byte("body"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	_, err := service.Execute("artifact", "", false, 100, 101, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	articleData, err := os.ReadFile(filepath.Join(outDir, "AF0100.md"))
	if err != nil {
		t.Fatalf("read article file failed: %v", err)
	}
	if !strings.Contains(string(articleData), "#06-af/article/note") {
		t.Fatalf("article tag missing: %s", string(articleData))
	}

	otherData, err := os.ReadFile(filepath.Join(outDir, "AF0101.md"))
	if err != nil {
		t.Fatalf("read other category file failed: %v", err)
	}
	if !strings.Contains(string(otherData), "#06-af/movie") {
		t.Fatalf("artifact category tag missing: %s", string(otherData))
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	t.Run("unsupported page type", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("content", "", false, 1, 1, "/tmp/artifacts.json", "/tmp/body", "/tmp/out")
		if err == nil || err.Error() != "未対応のpage-typeです: content" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte(`{invalid`), nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("artifact", "", false, 1, 1, "/tmp/artifacts.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "Artifact JSONの解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("artifact section not found", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/artifacts.json":
				return []byte(`[]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#golang\n"), nil
			default:
				return nil, errors.New("not found")
			}
		}

		service := NewService(mockRepo)
		_, err := service.Execute("artifact", "", false, 1, 1, "/tmp/artifacts.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "## Artifact") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("empty page title", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/artifacts.json":
				return []byte(`[
					{
						"con_id":"AF0001",
						"page_title":" ",
						"category":"system",
						"output_url":"https://example.com",
						"tags":[]
					}
				]`), nil
			case "/tmp/tags.md":
				return []byte("## Artifact\n#06-af/system/devbox\n\n## Frequent Tags\n#golang\n"), nil
			default:
				return []byte("source"), nil
			}
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}

		service := NewService(mockRepo)
		_, err := service.Execute("artifact", "", false, 1, 10, "/tmp/artifacts.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "page_title が空です") {
			t.Fatalf("error = %v", err)
		}
	})
}
