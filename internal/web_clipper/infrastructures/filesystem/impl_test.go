package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOSRepositoryFileOperations_Normal(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "input.txt")
	renamedPath := filepath.Join(dir, "renamed.txt")

	if err := repo.WriteFile(filePath, []byte("hello")); err != nil {
		t.Fatalf("WriteFileがエラーを返しました: %v", err)
	}

	data, err := repo.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFileがエラーを返しました: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("読み込み内容が期待と異なります: got=%q want=%q", string(data), "hello")
	}

	exists, err := repo.Exists(filePath)
	if err != nil {
		t.Fatalf("Existsがエラーを返しました: %v", err)
	}
	if !exists {
		t.Fatal("作成済みファイルが存在しない判定になりました")
	}

	if err := repo.Rename(filePath, renamedPath); err != nil {
		t.Fatalf("Renameがエラーを返しました: %v", err)
	}

	exists, err = repo.Exists(renamedPath)
	if err != nil {
		t.Fatalf("Existsがエラーを返しました: %v", err)
	}
	if !exists {
		t.Fatal("リネーム後ファイルが存在しない判定になりました")
	}
}

func TestOSRepositoryListDirectory_Normal(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644); err != nil {
		t.Fatalf("テストファイル作成に失敗しました: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "nested"), 0755); err != nil {
		t.Fatalf("テストディレクトリ作成に失敗しました: %v", err)
	}

	files, err := repo.ListDirectory(dir)
	if err != nil {
		t.Fatalf("ListDirectoryがエラーを返しました: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("取得件数が期待と異なります: got=%d want=2", len(files))
	}

	foundFile := false
	foundDir := false
	for _, file := range files {
		switch file.Name {
		case "b.txt":
			foundFile = true
			if file.IsDir {
				t.Fatal("通常ファイルがディレクトリ判定になりました")
			}
		case "nested":
			foundDir = true
			if !file.IsDir {
				t.Fatal("ディレクトリが通常ファイル判定になりました")
			}
		}
	}
	if !foundFile || !foundDir {
		t.Fatalf("期待するエントリが見つかりません: foundFile=%t foundDir=%t", foundFile, foundDir)
	}
}

func TestOSRepositoryErrors(t *testing.T) {
	t.Parallel()

	repo := NewRepository()
	dir := t.TempDir()
	filePath := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("テストファイル作成に失敗しました: %v", err)
	}

	if _, err := repo.ListDirectory(filePath); err == nil {
		t.Fatal("ファイルをListDirectoryへ渡した場合はエラーが期待されます")
	}

	exists, err := repo.Exists(filepath.Join(dir, "missing.txt"))
	if err != nil {
		t.Fatalf("存在しないファイルのExistsがエラーを返しました: %v", err)
	}
	if exists {
		t.Fatal("存在しないファイルが存在する判定になりました")
	}

	if _, err := repo.ReadFile(" "); err == nil {
		t.Fatal("空パスのReadFileはエラーが期待されます")
	}
}
