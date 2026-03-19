package usecases

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/markdown_crafter/config"
)

type operationDispatchRepository struct {
	readFileContents map[string]string
	listedFiles      map[string][]string
	createdDirs      []string
	writtenFiles     map[string]string
	removedFiles     []string
}

func newOperationDispatchRepository() *operationDispatchRepository {
	return &operationDispatchRepository{
		readFileContents: map[string]string{},
		listedFiles:      map[string][]string{},
		writtenFiles:     map[string]string{},
	}
}

func (r *operationDispatchRepository) ReadFile(filePath string) (string, error) {
	return r.readFileContents[filePath], nil
}

func (r *operationDispatchRepository) WriteFile(filePath string, content string) error {
	r.writtenFiles[filePath] = content
	return nil
}

func (r *operationDispatchRepository) CreateDir(dirPath string) error {
	r.createdDirs = append(r.createdDirs, dirPath)
	return nil
}

func (r *operationDispatchRepository) ListMarkdownFiles(dirPath string) ([]string, error) {
	return r.listedFiles[dirPath], nil
}

func (r *operationDispatchRepository) RemoveFile(filePath string) error {
	r.removedFiles = append(r.removedFiles, filePath)
	return nil
}

func TestService_ExecuteByConfig_Normal(t *testing.T) {
	tests := []struct {
		name   string
		cfg    *config.Config
		setup  func(repo *operationDispatchRepository)
		assert func(t *testing.T, result string, repo *operationDispatchRepository)
	}{
		{
			name: "SplitHeadings",
			cfg: &config.Config{
				Operation:    config.OperationSplitHeadings,
				FilePath:     "note.md",
				HeadingLevel: 2,
				OutputDir:    "out",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.readFileContents["note.md"] = "## 見出し\n本文\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if len(repo.createdDirs) != 1 || repo.createdDirs[0] != "out" {
					t.Fatalf("unexpected created dirs: %v", repo.createdDirs)
				}
				if _, ok := repo.writtenFiles[filepath.Join("out", "001.md")]; !ok {
					t.Fatalf("expected output file was not written: %v", repo.writtenFiles)
				}
				if !strings.Contains(result, "split-headings: 1 ファイルを出力しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "AddFrontMatter",
			cfg: &config.Config{
				Operation: config.OperationAddFrontMatter,
				FilePath:  "note.md",
				KVPairs:   []string{"title=記事"},
			},
			setup: func(repo *operationDispatchRepository) {
				repo.readFileContents["note.md"] = "本文\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if !strings.Contains(repo.writtenFiles["note.md"], "title: 記事") {
					t.Fatalf("front matter was not added: %q", repo.writtenFiles["note.md"])
				}
				if !strings.Contains(result, "add-front-matter: note.md") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "AddTagsByFile",
			cfg: &config.Config{
				Operation: config.OperationAddTags,
				FilePath:  "note.md",
				Tags:      "go,markdown",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.readFileContents["note.md"] = "本文\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if !strings.Contains(repo.writtenFiles["note.md"], "#go #markdown") {
					t.Fatalf("tags were not added: %q", repo.writtenFiles["note.md"])
				}
				if !strings.Contains(result, "add-tags: note.md にタグを追加しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "AddTagsByDir",
			cfg: &config.Config{
				Operation: config.OperationAddTags,
				DirPath:   "notes",
				Tags:      "go",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.listedFiles["notes"] = []string{"notes/a.md"}
				repo.readFileContents["notes/a.md"] = "本文\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if !strings.Contains(repo.writtenFiles["notes/a.md"], "#go") {
					t.Fatalf("tags were not added: %q", repo.writtenFiles["notes/a.md"])
				}
				if !strings.Contains(result, "add-tags: 1 ファイルにタグを追加しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "DeleteEmptyFiles",
			cfg: &config.Config{
				Operation: config.OperationDeleteEmptyFiles,
				DirPath:   "notes",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.listedFiles["notes"] = []string{"notes/a.md"}
				repo.readFileContents["notes/a.md"] = ""
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if len(repo.removedFiles) != 1 || repo.removedFiles[0] != "notes/a.md" {
					t.Fatalf("unexpected removed files: %v", repo.removedFiles)
				}
				if !strings.Contains(result, "delete-empty-files: 1 ファイルを削除しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "AddHeading1",
			cfg: &config.Config{
				Operation:       config.OperationAddHeading1,
				FilePath:        "note.md",
				HeadingText:     "概要",
				HeadingPosition: config.HeadingPositionHead,
			},
			setup: func(repo *operationDispatchRepository) {
				repo.readFileContents["note.md"] = "本文\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if !strings.Contains(repo.writtenFiles["note.md"], "# 概要") {
					t.Fatalf("heading was not added: %q", repo.writtenFiles["note.md"])
				}
				if !strings.Contains(result, "add-heading1: note.md に見出しを追加しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "ReplaceImages",
			cfg: &config.Config{
				Operation:       config.OperationReplaceImages,
				FilePath:        "note.md",
				ReplacementText: "(添付画像)",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.readFileContents["note.md"] = "![alt](img.png)\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if !strings.Contains(repo.writtenFiles["note.md"], "(添付画像)") {
					t.Fatalf("images were not replaced: %q", repo.writtenFiles["note.md"])
				}
				if !strings.Contains(result, "replace-images: note.md の画像記法 1 件を置換しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "RemoveHeadingAnnotations",
			cfg: &config.Config{
				Operation:    config.OperationRemoveHeadingAnnotations,
				FilePath:     "note.md",
				HeadingLevel: 3,
			},
			setup: func(repo *operationDispatchRepository) {
				repo.readFileContents["note.md"] = "### **注釈付き見出し**\n本文\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if !strings.Contains(repo.writtenFiles["note.md"], "### 注釈付き見出し") {
					t.Fatalf("heading annotations were not removed: %q", repo.writtenFiles["note.md"])
				}
				if !strings.Contains(result, "remove-heading-annotations: note.md の見出し注釈 1 件を除去しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "RemoveTitleHashTags",
			cfg: &config.Config{
				Operation: config.OperationRemoveTitleHashTags,
				DirPath:   "notes",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.listedFiles["notes"] = []string{"notes/a.md"}
				repo.readFileContents["notes/a.md"] = "## タイトル #Go - 記事\n- [リンク #Go](https://example.com)\n本文 #Keep\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				got := repo.writtenFiles["notes/a.md"]
				if !strings.Contains(got, "## タイトル - 記事\n- [リンク](https://example.com)\n本文 #Keep\n") {
					t.Fatalf("title hash tags were not removed: %q", got)
				}
				if !strings.Contains(result, "remove-title-hash-tags: 1 ファイルの先頭2行からハッシュタグを除去しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
		{
			name: "RemoveTitleHashTagsNoChange",
			cfg: &config.Config{
				Operation: config.OperationRemoveTitleHashTags,
				DirPath:   "notes",
			},
			setup: func(repo *operationDispatchRepository) {
				repo.listedFiles["notes"] = []string{"notes/a.md"}
				repo.readFileContents["notes/a.md"] = "title\nsecond line\nthird #Keep\n"
			},
			assert: func(t *testing.T, result string, repo *operationDispatchRepository) {
				t.Helper()
				if _, ok := repo.writtenFiles["notes/a.md"]; ok {
					t.Fatalf("unchanged file should not be written: %q", repo.writtenFiles["notes/a.md"])
				}
				if !strings.Contains(result, "remove-title-hash-tags: 0 ファイルの先頭2行からハッシュタグを除去しました") {
					t.Fatalf("unexpected result: %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_Normal", func(t *testing.T) {
			repo := newOperationDispatchRepository()
			tt.setup(repo)

			service := NewService(repo)
			result, err := service.ExecuteByConfig(tt.cfg)
			if err != nil {
				t.Fatalf("ExecuteByConfig returned error: %v", err)
			}

			tt.assert(t, result, repo)
		})
	}
}

func TestService_ExecuteByConfig_UnsupportedOperation(t *testing.T) {
	service := NewService(newOperationDispatchRepository())

	_, err := service.ExecuteByConfig(&config.Config{
		Operation: "unknown-op",
	})
	if err == nil {
		t.Fatal("expected unsupported operation error")
	}
	if !strings.Contains(err.Error(), "未サポートのoperationです") {
		t.Fatalf("unexpected error: %v", err)
	}
}
