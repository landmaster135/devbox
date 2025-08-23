package infrastructure

import (
	"errors"
	"os"
	"testing"
)

// TestMockFileSystem_MkdirAll はMockFileSystemのMkdirAllメソッドテスト
func TestMockFileSystem_MkdirAll(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockFileSystem)
		path        string
		perm        os.FileMode
		expectError bool
	}{
		{
			name: "WithMockFunc_Normal",
			setupMock: func(mock *MockFileSystem) {
				mock.MkdirAllFunc = func(path string, perm os.FileMode) error {
					return nil
				}
			},
			path:        "/test/path",
			perm:        0755,
			expectError: false,
		},
		{
			name: "WithMockFunc_Error",
			setupMock: func(mock *MockFileSystem) {
				mock.MkdirAllFunc = func(path string, perm os.FileMode) error {
					return errors.New("permission denied")
				}
			},
			path:        "/test/path",
			perm:        0755,
			expectError: true,
		},
		{
			name: "WithoutMockFunc_Normal",
			setupMock: func(mock *MockFileSystem) {
				// MkdirAllFuncを設定しない
			},
			path:        "/test/path",
			perm:        0755,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockFileSystem{}
			tt.setupMock(mock)

			// Act
			err := mock.MkdirAll(tt.path, tt.perm)

			// Assert
			if tt.expectError && err == nil {
				t.Error("エラーが期待されましたが、nilが返されました")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}

// TestMockFileSystem_WriteFile はMockFileSystemのWriteFileメソッドテスト
func TestMockFileSystem_WriteFile(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockFileSystem)
		filename    string
		data        []byte
		perm        os.FileMode
		expectError bool
	}{
		{
			name: "WithMockFunc_Normal",
			setupMock: func(mock *MockFileSystem) {
				mock.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
					return nil
				}
			},
			filename:    "test.txt",
			data:        []byte("test content"),
			perm:        0644,
			expectError: false,
		},
		{
			name: "WithMockFunc_Error",
			setupMock: func(mock *MockFileSystem) {
				mock.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
					return errors.New("write failed")
				}
			},
			filename:    "test.txt",
			data:        []byte("test content"),
			perm:        0644,
			expectError: true,
		},
		{
			name: "WithoutMockFunc_Normal",
			setupMock: func(mock *MockFileSystem) {
				// WriteFileFuncを設定しない
			},
			filename:    "test.txt",
			data:        []byte("test content"),
			perm:        0644,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockFileSystem{}
			tt.setupMock(mock)

			// Act
			err := mock.WriteFile(tt.filename, tt.data, tt.perm)

			// Assert
			if tt.expectError && err == nil {
				t.Error("エラーが期待されましたが、nilが返されました")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
		})
	}
}

// TestMockFileSystem_Join はMockFileSystemのJoinメソッドテスト
func TestMockFileSystem_Join(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(*MockFileSystem)
		elements  []string
		expected  string
	}{
		{
			name: "WithMockFunc_Normal",
			setupMock: func(mock *MockFileSystem) {
				mock.JoinFunc = func(elem ...string) string {
					return "mocked/path"
				}
			},
			elements: []string{"dir", "file.txt"},
			expected: "mocked/path",
		},
		{
			name: "WithoutMockFunc_DefaultImplementation",
			setupMock: func(mock *MockFileSystem) {
				// JoinFuncを設定しない（デフォルト実装を使用）
			},
			elements: []string{"dir", "file.txt"},
			expected: "dir/file.txt",
		},
		{
			name: "WithoutMockFunc_SingleElement",
			setupMock: func(mock *MockFileSystem) {
				// JoinFuncを設定しない（デフォルト実装を使用）
			},
			elements: []string{"file.txt"},
			expected: "file.txt",
		},
		{
			name: "WithoutMockFunc_ThreeElements",
			setupMock: func(mock *MockFileSystem) {
				// JoinFuncを設定しない（デフォルト実装を使用）
			},
			elements: []string{"root", "sub", "file.txt"},
			expected: "root/sub/file.txt",
		},
		{
			name: "WithoutMockFunc_EmptyElements",
			setupMock: func(mock *MockFileSystem) {
				// JoinFuncを設定しない（デフォルト実装を使用）
			},
			elements: []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockFileSystem{}
			tt.setupMock(mock)

			// Act
			result := mock.Join(tt.elements...)

			// Assert
			if result != tt.expected {
				t.Errorf("期待される結果: %s, 実際: %s", tt.expected, result)
			}
		})
	}
}

// TestMockJSONProcessor_MarshalIndent はMockJSONProcessorのMarshalIndentメソッドテスト
func TestMockJSONProcessor_MarshalIndent(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockJSONProcessor)
		input       interface{}
		prefix      string
		indent      string
		expected    []byte
		expectError bool
	}{
		{
			name: "WithMockFunc_Normal",
			setupMock: func(mock *MockJSONProcessor) {
				mock.MarshalIndentFunc = func(v interface{}, prefix, indent string) ([]byte, error) {
					return []byte(`{"mocked": "response"}`), nil
				}
			},
			input:       map[string]string{"key": "value"},
			prefix:      "",
			indent:      "  ",
			expected:    []byte(`{"mocked": "response"}`),
			expectError: false,
		},
		{
			name: "WithMockFunc_Error",
			setupMock: func(mock *MockJSONProcessor) {
				mock.MarshalIndentFunc = func(v interface{}, prefix, indent string) ([]byte, error) {
					return nil, errors.New("marshal failed")
				}
			},
			input:       map[string]string{"key": "value"},
			prefix:      "",
			indent:      "  ",
			expected:    nil,
			expectError: true,
		},
		{
			name: "WithoutMockFunc_DefaultImplementation",
			setupMock: func(mock *MockJSONProcessor) {
				// MarshalIndentFuncを設定しない（デフォルト実装を使用）
			},
			input:       map[string]string{"key": "value"},
			prefix:      "",
			indent:      "  ",
			expected:    []byte("{}"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mock := &MockJSONProcessor{}
			tt.setupMock(mock)

			// Act
			result, err := mock.MarshalIndent(tt.input, tt.prefix, tt.indent)

			// Assert
			if tt.expectError && err == nil {
				t.Error("エラーが期待されましたが、nilが返されました")
			}
			if !tt.expectError && err != nil {
				t.Errorf("エラーが発生しました: %v", err)
			}
			if string(result) != string(tt.expected) {
				t.Errorf("期待される結果: %s, 実際: %s", string(tt.expected), string(result))
			}
		})
	}
}

// TestMockFileSystem_Join_WithEmptyString はMockFileSystemのJoinメソッドの空文字列テスト
func TestMockFileSystem_Join_WithEmptyString(t *testing.T) {
	// Arrange
	mock := &MockFileSystem{}
	// JoinFuncを設定しない（デフォルト実装を使用）

	// Act
	result := mock.Join("dir", "", "file.txt")

	// Assert
	expected := "dir//file.txt"
	if result != expected {
		t.Errorf("期待される結果: %s, 実際: %s", expected, result)
	}
}

// TestMockFileSystem_Integration は複数のMockFileSystemメソッドの統合テスト
func TestMockFileSystem_Integration(t *testing.T) {
	// Arrange
	mock := &MockFileSystem{
		MkdirAllFunc: func(path string, perm os.FileMode) error {
			if path == "/invalid" {
				return errors.New("invalid path")
			}
			return nil
		},
		WriteFileFunc: func(filename string, data []byte, perm os.FileMode) error {
			if len(data) == 0 {
				return errors.New("empty data")
			}
			return nil
		},
		JoinFunc: func(elem ...string) string {
			return "custom/joined/path"
		},
	}

	// Act & Assert - MkdirAll正常系
	err := mock.MkdirAll("/valid", 0755)
	if err != nil {
		t.Errorf("MkdirAll正常系でエラーが発生しました: %v", err)
	}

	// Act & Assert - MkdirAllエラー系
	err = mock.MkdirAll("/invalid", 0755)
	if err == nil {
		t.Error("MkdirAllエラー系でエラーが期待されましたが、nilが返されました")
	}

	// Act & Assert - WriteFile正常系
	err = mock.WriteFile("test.txt", []byte("content"), 0644)
	if err != nil {
		t.Errorf("WriteFile正常系でエラーが発生しました: %v", err)
	}

	// Act & Assert - WriteFileエラー系
	err = mock.WriteFile("test.txt", []byte(""), 0644)
	if err == nil {
		t.Error("WriteFileエラー系でエラーが期待されましたが、nilが返されました")
	}

	// Act & Assert - Join
	result := mock.Join("dir", "file.txt")
	expected := "custom/joined/path"
	if result != expected {
		t.Errorf("Join結果が期待と異なります。期待: %s, 実際: %s", expected, result)
	}
}

// TestMockJSONProcessor_Integration はMockJSONProcessorの統合テスト
func TestMockJSONProcessor_Integration(t *testing.T) {
	// Arrange
	mock := &MockJSONProcessor{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			// 入力に応じて異なる結果を返す
			if str, ok := v.(string); ok && str == "error" {
				return nil, errors.New("marshal error")
			}
			return []byte(`{"processed": true}`), nil
		},
	}

	// Act & Assert - 正常系
	result, err := mock.MarshalIndent("valid", "", "  ")
	if err != nil {
		t.Errorf("正常系でエラーが発生しました: %v", err)
	}
	expected := `{"processed": true}`
	if string(result) != expected {
		t.Errorf("期待される結果: %s, 実際: %s", expected, string(result))
	}

	// Act & Assert - エラー系
	result, err = mock.MarshalIndent("error", "", "  ")
	if err == nil {
		t.Error("エラー系でエラーが期待されましたが、nilが返されました")
	}
	if result != nil {
		t.Error("エラー時は結果がnilである必要があります")
	}
}
