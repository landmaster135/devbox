package craftmarkdown

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
	if err := os.WriteFile(tagsFile, []byte("## Frequent Tags\n#31-programming/language/javascript\n"), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("本文です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000011.md"), []byte("範囲外です。\n"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
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
	if !strings.Contains(crafted, "# Google Apps Script") {
		t.Fatalf("heading missing: %s", crafted)
	}
	if !strings.Contains(crafted, "bought_at: 2024-01-01T09:00:00+09:00") {
		t.Fatalf("front matter missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#91-backup/tool-migration/202602-notion") {
		t.Fatalf("required tag missing: %s", crafted)
	}
	if !strings.Contains(crafted, "#31-programming/language/javascript") {
		t.Fatalf("frequent tag missing: %s", crafted)
	}
}

func TestService_Execute_SplitSourceFiles(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"CO000001946",
			"page_title":"split body sample",
			"category":"subscription_year",
			"owning_status":"already",
			"color":"gray",
			"bought_at":"2024-01-01T09:00:00+09:00",
			"score":70,
			"price":1000,
			"tags":[{"page_title":"JavaScript"}]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(tagsFile, []byte("## Frequent Tags\n#31-programming/language/javascript\n"), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000001946_02.md"), []byte("second\n"), 0644); err != nil {
		t.Fatalf("write split body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000001946_01.md"), []byte("first\n"), 0644); err != nil {
		t.Fatalf("write split body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("content", "", false, 1946, 1946, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=2") || !strings.Contains(result, "加工成功=2") {
		t.Fatalf("unexpected result: %s", result)
	}

	for _, fileName := range []string{"CO000001946_01.md", "CO000001946_02.md"} {
		outPath := filepath.Join(outDir, fileName)
		outData, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatalf("read split output failed (%s): %v", fileName, err)
		}
		outText := string(outData)
		if !strings.Contains(outText, "# split body sample") {
			t.Fatalf("heading missing (%s): %s", fileName, outText)
		}
		if !strings.Contains(outText, "con_id: CO000001946") {
			t.Fatalf("front matter missing (%s): %s", fileName, outText)
		}
		if !strings.Contains(outText, "#0a-content/subscription-year") {
			t.Fatalf("category tag missing (%s): %s", fileName, outText)
		}
	}

	if _, err := os.Stat(filepath.Join(outDir, "CO000001946.md")); !os.IsNotExist(err) {
		t.Fatalf("base output should not be created when split files exist")
	}
}

func TestService_Execute_ExactSourcePreferredOverSplit(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}

	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"CO000000010",
			"page_title":"exact preferred",
			"category":"software",
			"owning_status":"already",
			"color":"gray",
			"tags":[{"page_title":"JavaScript"}]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(tagsFile, []byte("## Frequent Tags\n#31-programming/language/javascript\n"), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("exact\n"), 0644); err != nil {
		t.Fatalf("write exact body failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010_01.md"), []byte("split\n"), 0644); err != nil {
		t.Fatalf("write split body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	result, err := service.Execute("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=1") {
		t.Fatalf("unexpected result: %s", result)
	}

	if _, err := os.Stat(filepath.Join(outDir, "CO000000010.md")); err != nil {
		t.Fatalf("exact output should exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outDir, "CO000000010_01.md")); !os.IsNotExist(err) {
		t.Fatalf("split output should not be generated when exact source exists")
	}
}

func TestService_Execute_MissingSourceFileCanSkipByFlag(t *testing.T) {
	t.Parallel()

	mockRepo := &filesystem.MockRepository{}
	mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
		switch path {
		case "/tmp/contents.json":
			return []byte(`[
				{
					"con_id": "CO000000010",
					"page_title": "Missing body",
					"category": "software",
					"owning_status": "already",
					"color": "gray",
					"bought_at": "2024-01-01T09:00:00+09:00",
					"score": 88,
					"price": 1200,
					"tags": [{"page_title": "JavaScript"}]
				}
			]`), nil
		case "/tmp/tags.md":
			return []byte("## Frequent Tags\n#31-programming/language/javascript\n"), nil
		default:
			return nil, errors.New("not found")
		}
	}
	mockRepo.FileExistsFunc = func(path string) (bool, error) {
		return false, nil
	}

	service := NewService(mockRepo)
	result, err := service.Execute("content", "", true, 10, 10, "/tmp/contents.json", "/tmp/body", "/tmp/out")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "対象件数=1") || !strings.Contains(result, "加工成功=0") {
		t.Fatalf("unexpected result: %s", result)
	}
	if len(mockRepo.WriteFileCalls) != 0 {
		t.Fatalf("WriteFile should not be called when skipsNoSrcBody=true")
	}
}

func TestService_Execute_Error(t *testing.T) {
	t.Parallel()

	t.Run("unsupported page type", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("artifact", "", false, 1, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || err.Error() != "未対応のpage-typeです: artifact" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid range value", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("content", "", false, 0, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || err.Error() != "con_number_start と con_number_end は1以上で指定してください" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("start greater than end", func(t *testing.T) {
		t.Parallel()
		service := NewService(&filesystem.MockRepository{})
		_, err := service.Execute("content", "", false, 2, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || err.Error() != "con_number_start は con_number_end 以下である必要があります" {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("read json error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return nil, errors.New("read failed")
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "src-json-file の読み込みに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			return []byte("{invalid"), nil
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "Content JSONの解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("load frequent tags error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			if path == "/tmp/contents.json" {
				return []byte(`[]`), nil
			}
			return nil, errors.New("read tags failed")
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "tags.md の読み込みに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("mkdir all error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/contents.json":
				return []byte(`[]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#foo\n"), nil
			default:
				return nil, errors.New("not found")
			}
		}
		mockRepo.MkdirAllFunc = func(path string, perm os.FileMode) error {
			return errors.New("mkdir failed")
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 1, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "出力ディレクトリの作成に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("parse con number error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/contents.json":
				return []byte(`[
					{
						"con_id": "CO-invalid",
						"page_title": "sample",
						"category": "software",
						"owning_status": "already",
						"color": "gray",
						"tags": []
					}
				]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#foo\n"), nil
			default:
				return nil, errors.New("not found")
			}
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 10, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "con_id の解析に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("file exists error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/contents.json":
				return []byte(`[
					{
						"con_id": "CO0001",
						"page_title": "sample",
						"category": "software",
						"owning_status": "already",
						"color": "gray",
						"tags": []
					}
				]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#foo\n"), nil
			default:
				return nil, errors.New("not found")
			}
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return false, errors.New("stat failed")
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 10, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "コピー元ファイルの確認に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("copy file error", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/contents.json":
				return []byte(`[
					{
						"con_id": "CO0001",
						"page_title": "sample",
						"category": "software",
						"owning_status": "already",
						"color": "gray",
						"tags": []
					}
				]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#foo\n"), nil
			default:
				return nil, errors.New("not found")
			}
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		mockRepo.CopyFileFunc = func(srcPath, dstPath string) error {
			return errors.New("copy failed")
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 10, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "Markdownのコピーに失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("write file error on empty markdown create", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/contents.json":
				return []byte(`[
					{
						"con_id": "CO0001",
						"page_title": "sample",
						"category": "software",
						"owning_status": "already",
						"color": "gray",
						"tags": []
					}
				]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#foo\n"), nil
			default:
				return nil, errors.New("not found")
			}
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return false, nil
		}
		mockRepo.WriteFileFunc = func(path string, data []byte) error {
			return errors.New("write failed")
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 10, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "空Markdownの作成に失敗しました") {
			t.Fatalf("error = %v", err)
		}
	})

	t.Run("empty page title", func(t *testing.T) {
		t.Parallel()
		mockRepo := &filesystem.MockRepository{}
		mockRepo.ReadFileFunc = func(path string) ([]byte, error) {
			switch path {
			case "/tmp/contents.json":
				return []byte(`[
					{
						"con_id": "CO0001",
						"page_title": " ",
						"category": "software",
						"owning_status": "already",
						"color": "gray",
						"tags": []
					}
				]`), nil
			case "/tmp/tags.md":
				return []byte("## Frequent Tags\n#foo\n"), nil
			default:
				return []byte("source"), nil
			}
		}
		mockRepo.FileExistsFunc = func(path string) (bool, error) {
			return true, nil
		}
		service := NewService(mockRepo)
		_, err := service.Execute("content", "", false, 1, 10, "/tmp/contents.json", "/tmp/body", "/tmp/out")
		if err == nil || !strings.Contains(err.Error(), "page_title が空です") {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestService_Execute_InvalidCategory(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	srcJSONFile := filepath.Join(tmpDir, "contents.json")
	tagsFile := filepath.Join(tmpDir, "tags.md")
	srcBodyDir := filepath.Join(tmpDir, "body")
	outDir := filepath.Join(tmpDir, "out")

	if err := os.MkdirAll(srcBodyDir, 0755); err != nil {
		t.Fatalf("mkdir failed: %v", err)
	}
	if err := os.WriteFile(srcJSONFile, []byte(`[
		{
			"con_id":"CO000000010",
			"page_title":"Invalid category",
			"category":"unknown",
			"owning_status":"already",
			"color":"gray",
			"tags":[{"page_title":"JavaScript"}]
		}
	]`), 0644); err != nil {
		t.Fatalf("write json failed: %v", err)
	}
	if err := os.WriteFile(tagsFile, []byte("## Frequent Tags\n#31-programming/language/javascript\n"), 0644); err != nil {
		t.Fatalf("write tags failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(srcBodyDir, "CO000000010.md"), []byte("body"), 0644); err != nil {
		t.Fatalf("write body failed: %v", err)
	}

	service := NewService(filesystem.NewRepository())
	_, err := service.Execute("content", "", false, 10, 10, srcJSONFile, srcBodyDir, outDir)
	if err == nil || !strings.Contains(err.Error(), "タグ生成に失敗しました") {
		t.Fatalf("error = %v", err)
	}
}
