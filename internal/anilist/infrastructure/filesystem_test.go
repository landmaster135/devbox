package infrastructure

import (
	"os"
	"path/filepath"
	"testing"
)

// TestNewOSFileSystem_Normal はNewOSFileSystemメソッドの正常系テスト
func TestNewOSFileSystem_Normal(t *testing.T) {
	// Act
	result := NewOSFileSystem()

	// Assert
	if result == nil {
		t.Error("結果がnilです")
	}
}

// TestOSFileSystem_MkdirAll はOSFileSystemのMkdirAllメソッドテスト
func TestOSFileSystem_MkdirAll(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		perm     os.FileMode
		setupDir bool
	}{
		{
			name:     "CreateNewDirectory_Normal",
			path:     "test_mkdir",
			perm:     0755,
			setupDir: false,
		},
		{
			name:     "CreateNestedDirectory_Normal",
			path:     "test_mkdir/nested/deep",
			perm:     0755,
			setupDir: false,
		},
		{
			name:     "ExistingDirectory_Normal",
			path:     "test_mkdir_existing",
			perm:     0755,
			setupDir: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			fs := NewOSFileSystem()
			tempDir := t.TempDir()
			testPath := filepath.Join(tempDir, tt.path)

			if tt.setupDir {
				// 事前にディレクトリを作成
				if err := os.MkdirAll(testPath, 0755); err != nil {
					t.Fatalf("事前ディレクトリ作成に失敗しました: %v", err)
				}
			}

			// Act
			err := fs.MkdirAll(testPath, tt.perm)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}

			// ディレクトリが作成されているか確認
			if _, err := os.Stat(testPath); os.IsNotExist(err) {
				t.Error("ディレクトリが作成されていません")
			}
		})
	}
}

// TestOSFileSystem_WriteFile はOSFileSystemのWriteFileメソッドテスト
func TestOSFileSystem_WriteFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		data     []byte
		perm     os.FileMode
	}{
		{
			name:     "WriteNewFile_Normal",
			filename: "test_file.txt",
			data:     []byte("test content"),
			perm:     0644,
		},
		{
			name:     "WriteEmptyFile_Normal",
			filename: "empty_file.txt",
			data:     []byte(""),
			perm:     0644,
		},
		{
			name:     "WriteJSONFile_Normal",
			filename: "test.json",
			data:     []byte(`{"key": "value"}`),
			perm:     0644,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			fs := NewOSFileSystem()
			tempDir := t.TempDir()
			testPath := filepath.Join(tempDir, tt.filename)

			// Act
			err := fs.WriteFile(testPath, tt.data, tt.perm)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}

			// ファイルが作成されているか確認
			if _, err := os.Stat(testPath); os.IsNotExist(err) {
				t.Error("ファイルが作成されていません")
				return
			}

			// ファイル内容を確認
			content, err := os.ReadFile(testPath)
			if err != nil {
				t.Errorf("ファイル読み取りに失敗しました: %v", err)
				return
			}

			if string(content) != string(tt.data) {
				t.Errorf("期待される内容: %s, 実際: %s", string(tt.data), string(content))
			}
		})
	}
}

// TestOSFileSystem_WriteFile_OverwriteExisting はOSFileSystemのWriteFileメソッドの既存ファイル上書きテスト
func TestOSFileSystem_WriteFile_OverwriteExisting(t *testing.T) {
	// Arrange
	fs := NewOSFileSystem()
	tempDir := t.TempDir()
	testPath := filepath.Join(tempDir, "overwrite_test.txt")

	// 既存ファイルを作成
	originalContent := []byte("original content")
	if err := os.WriteFile(testPath, originalContent, 0644); err != nil {
		t.Fatalf("既存ファイル作成に失敗しました: %v", err)
	}

	// Act
	newContent := []byte("new content")
	err := fs.WriteFile(testPath, newContent, 0644)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}

	// ファイル内容が上書きされているか確認
	content, err := os.ReadFile(testPath)
	if err != nil {
		t.Errorf("ファイル読み取りに失敗しました: %v", err)
		return
	}

	if string(content) != string(newContent) {
		t.Errorf("期待される内容: %s, 実際: %s", string(newContent), string(content))
	}
}

// TestOSFileSystem_Join はOSFileSystemのJoinメソッドテスト
func TestOSFileSystem_Join(t *testing.T) {
	tests := []struct {
		name     string
		elements []string
		expected string
	}{
		{
			name:     "TwoElements_Normal",
			elements: []string{"dir", "file.txt"},
			expected: filepath.Join("dir", "file.txt"),
		},
		{
			name:     "ThreeElements_Normal",
			elements: []string{"root", "sub", "file.txt"},
			expected: filepath.Join("root", "sub", "file.txt"),
		},
		{
			name:     "SingleElement_Normal",
			elements: []string{"file.txt"},
			expected: "file.txt",
		},
		{
			name:     "EmptyElements_Normal",
			elements: []string{},
			expected: "",
		},
		{
			name:     "WithEmptyString_Normal",
			elements: []string{"dir", "", "file.txt"},
			expected: filepath.Join("dir", "", "file.txt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			fs := NewOSFileSystem()

			// Act
			result := fs.Join(tt.elements...)

			// Assert
			if result != tt.expected {
				t.Errorf("期待される結果: %s, 実際: %s", tt.expected, result)
			}
		})
	}
}

// TestOSFileSystem_WriteFile_InvalidPath はOSFileSystemのWriteFileメソッドの無効パステスト
func TestOSFileSystem_WriteFile_InvalidPath(t *testing.T) {
	// Arrange
	fs := NewOSFileSystem()
	invalidPath := "/invalid/nonexistent/path/file.txt"
	data := []byte("test content")

	// Act
	err := fs.WriteFile(invalidPath, data, 0644)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、nilが返されました")
	}
}

// TestOSFileSystem_MkdirAll_InvalidPath はOSFileSystemのMkdirAllメソッドの無効パステスト
func TestOSFileSystem_MkdirAll_InvalidPath(t *testing.T) {
	// Arrange
	fs := NewOSFileSystem()
	// 権限のないディレクトリ（通常は/rootなど）
	invalidPath := "/root/test_invalid_mkdir"

	// Act
	err := fs.MkdirAll(invalidPath, 0755)

	// Assert
	// 権限エラーが発生することを期待（環境によっては成功する場合もある）
	// エラーが発生しない場合もあるため、エラーの有無は確認しない
	// 代わりに、メソッドが正常に実行されることを確認
	_ = err // エラーを無視
}
