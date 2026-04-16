package removetitlehashtags

import (
	"errors"
	"strings"
	"testing"
)

type removeTitleHashTagsMockRepository struct {
	listedFiles     []string
	listErr         error
	readContents    map[string]string
	readErrors      map[string]error
	writtenContents map[string]string
	writeErrors     map[string]error
}

func (r *removeTitleHashTagsMockRepository) ReadFile(filePath string) (string, error) {
	if err, ok := r.readErrors[filePath]; ok {
		return "", err
	}
	return r.readContents[filePath], nil
}

func (r *removeTitleHashTagsMockRepository) WriteFile(filePath string, content string) error {
	r.writtenContents[filePath] = content
	if err, ok := r.writeErrors[filePath]; ok {
		return err
	}
	return nil
}

func (r *removeTitleHashTagsMockRepository) CreateDir(_ string) error {
	return nil
}

func (r *removeTitleHashTagsMockRepository) ListMarkdownFiles(_ string) ([]string, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.listedFiles, nil
}

func (r *removeTitleHashTagsMockRepository) RemoveFile(_ string) error {
	return nil
}

func newRemoveTitleHashTagsMockRepository() *removeTitleHashTagsMockRepository {
	return &removeTitleHashTagsMockRepository{
		readContents:    map[string]string{},
		readErrors:      map[string]error{},
		writtenContents: map[string]string{},
		writeErrors:     map[string]error{},
	}
}

func TestService_RemoveTitleHashTags_Normal(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	repo.listedFiles = []string{"notes/a.md", "notes/b.md"}
	repo.readContents["notes/a.md"] = "## pythonのsetup.pyについてまとめる #Python - Qiita 要約\n- [pythonのsetup.pyについてまとめる #Python - Qiita](https://example.com)\n本文 #Keep\n"
	repo.readContents["notes/b.md"] = "#Python - Tips\nsecond #Tag line\nthird #NoChange\n"

	service := NewService(repo)
	result, err := service.Execute("notes", 1, 2)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	expectedA := "## pythonのsetup.pyについてまとめる - Qiita 要約\n- [pythonのsetup.pyについてまとめる - Qiita](https://example.com)\n本文 #Keep\n"
	if repo.writtenContents["notes/a.md"] != expectedA {
		t.Fatalf("unexpected content for a.md:\nexpected:\n%q\ngot:\n%q", expectedA, repo.writtenContents["notes/a.md"])
	}

	expectedB := "- Tips\nsecond line\nthird #NoChange\n"
	if repo.writtenContents["notes/b.md"] != expectedB {
		t.Fatalf("unexpected content for b.md:\nexpected:\n%q\ngot:\n%q", expectedB, repo.writtenContents["notes/b.md"])
	}

	if !strings.Contains(result, "remove-title-hash-tags: 2 ファイルの 1 行目から 2 行目までのハッシュタグを除去しました") {
		t.Fatalf("unexpected result: %q", result)
	}
	if !strings.Contains(result, "notes/a.md") || !strings.Contains(result, "notes/b.md") {
		t.Fatalf("result does not contain updated file paths: %q", result)
	}
}

func TestService_RemoveTitleHashTags_NoMarkdownFiles_Normal(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	service := NewService(repo)

	result, err := service.Execute("notes", 1, 2)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !strings.Contains(result, "0件") {
		t.Fatalf("unexpected result: %q", result)
	}
}

func TestService_RemoveTitleHashTags_NoChangeFile_Normal(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	repo.listedFiles = []string{"notes/a.md"}
	repo.readContents["notes/a.md"] = "title\nsecond line\nthird #Keep\n"
	service := NewService(repo)

	result, err := service.Execute("notes", 1, 2)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if !strings.Contains(result, "remove-title-hash-tags: 0 ファイルの 1 行目から 2 行目までのハッシュタグを除去しました") {
		t.Fatalf("unexpected result: %q", result)
	}
	if len(repo.writtenContents) != 0 {
		t.Fatalf("unexpected write on unchanged content: %+v", repo.writtenContents)
	}
}

func TestService_RemoveTitleHashTags_ListError(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	repo.listErr = errors.New("list error")
	service := NewService(repo)

	if _, err := service.Execute("notes", 1, 2); err == nil {
		t.Fatal("expected list error")
	}
}

func TestService_RemoveTitleHashTags_ReadError(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	repo.listedFiles = []string{"notes/a.md"}
	repo.readErrors["notes/a.md"] = errors.New("read error")
	service := NewService(repo)

	if _, err := service.Execute("notes", 1, 2); err == nil {
		t.Fatal("expected read error")
	}
}

func TestService_RemoveTitleHashTags_WriteError(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	repo.listedFiles = []string{"notes/a.md"}
	repo.readContents["notes/a.md"] = "## title #Tag\nline #Tag\n"
	repo.writeErrors["notes/a.md"] = errors.New("write error")
	service := NewService(repo)

	if _, err := service.Execute("notes", 1, 2); err == nil {
		t.Fatal("expected write error")
	}
}

func TestService_RemoveTitleHashTags_InvalidLineRange(t *testing.T) {
	t.Parallel()

	repo := newRemoveTitleHashTagsMockRepository()
	service := NewService(repo)

	if _, err := service.Execute("notes", 0, 2); err == nil {
		t.Fatal("expected invalid line range error")
	}
	if _, err := service.Execute("notes", 2, 1); err == nil {
		t.Fatal("expected invalid line range error")
	}
}

func TestRemoveTitleHashTagsFromContent_OnlyFirstTwoLines_Normal(t *testing.T) {
	t.Parallel()

	content := "line1 #A\nline2 #B\nline3 #C\n"
	got := removeTitleHashTagsFromContent(content, 1, 2)
	expected := "line1\nline2\nline3 #C\n"
	if got != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestRemoveTitleHashTagsFromContent_ShortContent_Normal(t *testing.T) {
	t.Parallel()

	content := "#Python"
	got := removeTitleHashTagsFromContent(content, 1, 2)
	expected := ""
	if got != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestRemoveTitleHashTagsFromContent_LineRange_Normal(t *testing.T) {
	t.Parallel()

	content := "line1 #A\nline2 #B\nline3 #C\n"
	got := removeTitleHashTagsFromContent(content, 2, 3)
	expected := "line1 #A\nline2\nline3\n"
	if got != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestRemoveTitleHashTagsFromContent_StartLineOverFileLength_Normal(t *testing.T) {
	t.Parallel()

	content := "line1 #A\nline2 #B\n"
	got := removeTitleHashTagsFromContent(content, 4, 5)
	if got != content {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", content, got)
	}
}

func TestRemoveHashTagsFromLine_HeadingPrefix_Normal(t *testing.T) {
	t.Parallel()

	line := "## heading #Tag - title\n"
	got := removeHashTagsFromLine(line)
	expected := "## heading - title\n"
	if got != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestRemoveHashTagsFromLine_EmptyLine_Normal(t *testing.T) {
	t.Parallel()

	line := "\n"
	got := removeHashTagsFromLine(line)
	if got != line {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", line, got)
	}
}

func TestRemoveHashTagsFromLine_WithoutTrailingNewline_Normal(t *testing.T) {
	t.Parallel()

	line := "title #Tag"
	got := removeHashTagsFromLine(line)
	expected := "title"
	if got != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}

func TestRemoveTitleHashTagsFromContent_FullWidthSpaceAndMultipleTags_Normal(t *testing.T) {
	t.Parallel()

	content := "## Goの初心者が見ると幸せになれる場所 #golang #Go 要約\n- [Goの初心者が見ると幸せになれる場所　#golang #Go - Qiita](https://qiita.com/tenntenn/items/0e33a4959250d1a55045)\n本文\n"
	got := removeTitleHashTagsFromContent(content, 1, 2)
	expected := "## Goの初心者が見ると幸せになれる場所 要約\n- [Goの初心者が見ると幸せになれる場所 - Qiita](https://qiita.com/tenntenn/items/0e33a4959250d1a55045)\n本文\n"
	if got != expected {
		t.Fatalf("unexpected content:\nexpected:\n%q\ngot:\n%q", expected, got)
	}
}
