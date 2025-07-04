package finder

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

// TestFindFiles_Normal はFindFiles関数の正常系テストです
func TestFindFiles_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "finder_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイル構造を作成
	testFiles := []struct {
		path    string
		content string
	}{
		{filepath.Join(tempDir, "file1.go"), "package main\n\nfunc main() {}\n"},
		{filepath.Join(tempDir, "file2.go"), "package test\n\nfunc test() {}\n"},
		{filepath.Join(tempDir, "file.txt"), "This is a text file"},
		{filepath.Join(tempDir, "subdir", "file3.go"), "package subdir\n\nfunc sub() {}\n"},
		{filepath.Join(tempDir, "subdir", "file4.py"), "def main():\n    pass\n"},
	}

	for _, tf := range testFiles {
		dir := filepath.Dir(tf.path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("ディレクトリの作成に失敗しました: %v", err)
		}
		if err := os.WriteFile(tf.path, []byte(tf.content), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	// 非再帰的検索のテスト
	t.Run("NonRecursive", func(t *testing.T) {
		files, err := FindFiles(tempDir, false, []string{".go"})
		if err != nil {
			t.Fatalf("FindFiles関数の実行に失敗しました: %v", err)
		}

		// 期待されるファイル数
		if len(files) != 2 {
			t.Errorf("期待するファイル数は2ですが、実際は%dでした: %v", len(files), files)
		}

		// ファイル名を抽出して比較
		var fileNames []string
		for _, f := range files {
			fileNames = append(fileNames, filepath.Base(f))
		}
		sort.Strings(fileNames)

		expectedNames := []string{"file1.go", "file2.go"}
		if !reflect.DeepEqual(fileNames, expectedNames) {
			t.Errorf("期待するファイル名: %v, 実際のファイル名: %v", expectedNames, fileNames)
		}
	})

	// 再帰的検索のテスト
	t.Run("Recursive", func(t *testing.T) {
		files, err := FindFiles(tempDir, true, []string{".go"})
		if err != nil {
			t.Fatalf("FindFiles関数の実行に失敗しました: %v", err)
		}

		// 期待されるファイル数
		if len(files) != 3 {
			t.Errorf("期待するファイル数は3ですが、実際は%dでした: %v", len(files), files)
		}

		// ファイル名を抽出して比較
		var fileNames []string
		for _, f := range files {
			fileNames = append(fileNames, filepath.Base(f))
		}
		sort.Strings(fileNames)

		expectedNames := []string{"file1.go", "file2.go", "file3.go"}
		if !reflect.DeepEqual(fileNames, expectedNames) {
			t.Errorf("期待するファイル名: %v, 実際のファイル名: %v", expectedNames, fileNames)
		}
	})

	// 複数拡張子のテスト
	t.Run("MultipleExtensions", func(t *testing.T) {
		files, err := FindFiles(tempDir, true, []string{".go", ".py"})
		if err != nil {
			t.Fatalf("FindFiles関数の実行に失敗しました: %v", err)
		}

		// 期待されるファイル数
		if len(files) != 4 {
			t.Errorf("期待するファイル数は4ですが、実際は%dでした: %v", len(files), files)
		}

		// ファイル名を抽出して比較
		var fileNames []string
		for _, f := range files {
			fileNames = append(fileNames, filepath.Base(f))
		}
		sort.Strings(fileNames)

		expectedNames := []string{"file1.go", "file2.go", "file3.go", "file4.py"}
		if !reflect.DeepEqual(fileNames, expectedNames) {
			t.Errorf("期待するファイル名: %v, 実際のファイル名: %v", expectedNames, fileNames)
		}
	})

	// 空の拡張子リストのテスト
	t.Run("EmptyExtensions", func(t *testing.T) {
		files, err := FindFiles(tempDir, false, []string{})
		if err != nil {
			t.Fatalf("FindFiles関数の実行に失敗しました: %v", err)
		}

		// 期待されるファイル数（空の拡張子リストでは何も見つからない）
		if len(files) != 0 {
			t.Errorf("期待するファイル数は0ですが、実際は%dでした: %v", len(files), files)
		}
	})
}

// TestFindFiles_DirectoryNotFound はFindFiles関数のディレクトリ不在テストです
func TestFindFiles_DirectoryNotFound(t *testing.T) {
	// 存在しないディレクトリパス
	nonExistentPath := "/path/to/nonexistent/directory"

	// テスト実行
	_, err := FindFiles(nonExistentPath, false, []string{".go"})
	if err == nil {
		t.Errorf("存在しないディレクトリでもエラーが発生しませんでした")
	}
}

// TestFindFiles_InvalidDirectory はFindFiles関数の無効なディレクトリテストです
func TestFindFiles_InvalidDirectory(t *testing.T) {
	// テスト用の一時ファイルを作成
	tempFile, err := os.CreateTemp("", "finder_test_file")
	if err != nil {
		t.Fatalf("テスト用ファイルの作成に失敗しました: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// ファイルをディレクトリとして指定
	_, err = FindFiles(tempFile.Name(), false, []string{".go"})
	if err == nil {
		t.Errorf("ファイルをディレクトリとして指定してもエラーが発生しませんでした")
	}
}
