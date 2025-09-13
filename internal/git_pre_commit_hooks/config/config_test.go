package config

import (
	"errors"
	"os"
	"testing"
	"time"
)

// MockFileInfo はテスト用のFileInfo実装
type MockFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (m MockFileInfo) Name() string       { return m.name }
func (m MockFileInfo) Size() int64        { return m.size }
func (m MockFileInfo) Mode() os.FileMode  { return m.mode }
func (m MockFileInfo) ModTime() time.Time { return m.modTime }
func (m MockFileInfo) IsDir() bool        { return m.isDir }
func (m MockFileInfo) Sys() interface{}   { return nil }

// TestLoadConfigWithFileSystem_Normal はLoadConfigWithFileSystemの正常系テスト
func TestLoadConfigWithFileSystem_Normal(t *testing.T) {
	testCases := []struct {
		name           string
		filename       string
		mockStatFunc   func(filename string) (os.FileInfo, error)
		mockReadFunc   func(filename string) ([]byte, error)
		expectedError  bool
		expectedConfig bool
	}{
		{
			name:     "Valid JSON config",
			filename: "test.json",
			mockStatFunc: func(filename string) (os.FileInfo, error) {
				return MockFileInfo{name: "test.json", size: 100}, nil
			},
			mockReadFunc: func(filename string) ([]byte, error) {
				return []byte(`{
					"mcpServers": {
						"test-server": {
							"command": "test-command",
							"args": ["arg1", "arg2"],
							"env": {
								"API_KEY": "YOUR_API_KEY",
								"SECRET": "test-secret"
							}
						}
					}
				}`), nil
			},
			expectedError:  false,
			expectedConfig: true,
		},
		{
			name:     "File does not exist",
			filename: "nonexistent.json",
			mockStatFunc: func(filename string) (os.FileInfo, error) {
				return nil, os.ErrNotExist
			},
			mockReadFunc: func(filename string) ([]byte, error) {
				return nil, nil
			},
			expectedError:  true,
			expectedConfig: false,
		},
		{
			name:     "File read error",
			filename: "test.json",
			mockStatFunc: func(filename string) (os.FileInfo, error) {
				return MockFileInfo{name: "test.json", size: 100}, nil
			},
			mockReadFunc: func(filename string) ([]byte, error) {
				return nil, errors.New("permission denied")
			},
			expectedError:  true,
			expectedConfig: false,
		},
		{
			name:     "Invalid JSON format",
			filename: "invalid.json",
			mockStatFunc: func(filename string) (os.FileInfo, error) {
				return MockFileInfo{name: "invalid.json", size: 50}, nil
			},
			mockReadFunc: func(filename string) ([]byte, error) {
				return []byte(`{invalid json`), nil
			},
			expectedError:  true,
			expectedConfig: false,
		},
		{
			name:     "Empty JSON",
			filename: "empty.json",
			mockStatFunc: func(filename string) (os.FileInfo, error) {
				return MockFileInfo{name: "empty.json", size: 2}, nil
			},
			mockReadFunc: func(filename string) ([]byte, error) {
				return []byte(`{}`), nil
			},
			expectedError:  false,
			expectedConfig: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockFS := &MockFileSystem{
				StatFunc:     tc.mockStatFunc,
				ReadFileFunc: tc.mockReadFunc,
			}

			config, err := LoadConfigWithFileSystem(tc.filename, mockFS)

			if tc.expectedError {
				if err == nil {
					t.Errorf("LoadConfigWithFileSystem() expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadConfigWithFileSystem() unexpected error: %v", err)
				return
			}

			if tc.expectedConfig {
				if config == nil {
					t.Errorf("LoadConfigWithFileSystem() expected config, but got nil")
				}
			}
		})
	}
}

// TestLoadConfig_Normal はLoadConfigの正常系テスト（実際のファイルシステムを使用）
func TestLoadConfig_Normal(t *testing.T) {
	// このテストは実際のファイルシステムに依存するため、
	// 通常は統合テストとして別途実装することを推奨
	// ここでは関数が存在することのみを確認
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Errorf("LoadConfig() with nonexistent file should return error")
	}
}
