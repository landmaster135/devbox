package repositories

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/landmaster135/devbox/internal/domain/models"
	domainRepo "github.com/landmaster135/devbox/internal/domain/repositories"
)

// FileRepositoryMock は FileRepository インターフェースのモック実装です
type FileRepositoryMock struct {
	ReadFileFunc       func(path string) (*models.FileContent, error)
	WriteFileFunc      func(path string, content *models.FileContent) error
	FileExistsFunc     func(path string) bool
	FindFilesByExtFunc func(dirPath, ext string) ([]string, error)
	HasFilesWithExtFunc func(dirPath, ext string) (bool, error)
	ReadJSONFileFunc   func(path string) (interface{}, error)
	GetDirectoryPathFunc func(path string) string
	CreateDirectoryFunc func(dirPath string) error
	ReadDirFunc        func(dirPath string) ([]*models.DirEntry, error)
}

// ReadFile はモックの ReadFile 実装です
func (m *FileRepositoryMock) ReadFile(path string) (*models.FileContent, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, errors.New("ReadFile not implemented")
}

// WriteFile はモックの WriteFile 実装です
func (m *FileRepositoryMock) WriteFile(path string, content *models.FileContent) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, content)
	}
	return errors.New("WriteFile not implemented")
}

// FileExists はモックの FileExists 実装です
func (m *FileRepositoryMock) FileExists(path string) bool {
	if m.FileExistsFunc != nil {
		return m.FileExistsFunc(path)
	}
	return false
}

// FindFilesByExt はモックの FindFilesByExt 実装です
func (m *FileRepositoryMock) FindFilesByExt(dirPath, ext string) ([]string, error) {
	if m.FindFilesByExtFunc != nil {
		return m.FindFilesByExtFunc(dirPath, ext)
	}
	return nil, errors.New("FindFilesByExt not implemented")
}

// HasFilesWithExt はモックの HasFilesWithExt 実装です
func (m *FileRepositoryMock) HasFilesWithExt(dirPath, ext string) (bool, error) {
	if m.HasFilesWithExtFunc != nil {
		return m.HasFilesWithExtFunc(dirPath, ext)
	}
	return false, errors.New("HasFilesWithExt not implemented")
}

// ReadJSONFile はモックの ReadJSONFile 実装です
func (m *FileRepositoryMock) ReadJSONFile(path string) (interface{}, error) {
	if m.ReadJSONFileFunc != nil {
		return m.ReadJSONFileFunc(path)
	}
	return nil, errors.New("ReadJSONFile not implemented")
}

// GetDirectoryPath はモックの GetDirectoryPath 実装です
func (m *FileRepositoryMock) GetDirectoryPath(path string) string {
	if m.GetDirectoryPathFunc != nil {
		return m.GetDirectoryPathFunc(path)
	}
	return ""
}

// CreateDirectory はモックの CreateDirectory 実装です
func (m *FileRepositoryMock) CreateDirectory(dirPath string) error {
	if m.CreateDirectoryFunc != nil {
		return m.CreateDirectoryFunc(dirPath)
	}
	return errors.New("CreateDirectory not implemented")
}

// ReadDir はモックの ReadDir 実装です
func (m *FileRepositoryMock) ReadDir(dirPath string) ([]*models.DirEntry, error) {
	if m.ReadDirFunc != nil {
		return m.ReadDirFunc(dirPath)
	}
	return nil, errors.New("ReadDir not implemented")
}

// ISO8601RepositoryMock は ISO8601Repository インターフェースのモック実装です
type ISO8601RepositoryMock struct {
	ParseISO8601Func func(dateStr string) (int64, error)
}

// ParseISO8601 はモックの ParseISO8601 実装です
func (m *ISO8601RepositoryMock) ParseISO8601(dateStr string) (int64, error) {
	if m.ParseISO8601Func != nil {
		return m.ParseISO8601Func(dateStr)
	}
	return 0, errors.New("ParseISO8601 not implemented")
}

// TestNewJSONRepository は NewJSONRepository 関数をテストします
func TestNewJSONRepository(t *testing.T) {
	// モックを作成
	fileRepo := &FileRepositoryMock{}
	iso8601Repo := &ISO8601RepositoryMock{}

	// テスト実行
	repo := NewJSONRepository(fileRepo, iso8601Repo)
	if repo == nil {
		t.Error("NewJSONRepository() がnilを返しました")
	}
}

// TestJSONRepositoryImpl_FindJSONFiles は FindJSONFiles メソッドをテストします
func TestJSONRepositoryImpl_FindJSONFiles(t *testing.T) {
	// テスト用のディレクトリ構造を作成
	tempDir, err := os.MkdirTemp("", "json_repository_test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("サブディレクトリの作成に失敗しました: %v", err)
	}

	// テスト用のJSONファイルを作成
	jsonFile1 := filepath.Join(tempDir, "file1.json")
	if err := os.WriteFile(jsonFile1, []byte("{}"), 0644); err != nil {
		t.Fatalf("JSONファイルの作成に失敗しました: %v", err)
	}

	jsonFile2 := filepath.Join(subDir, "file2.json")
	if err := os.WriteFile(jsonFile2, []byte("{}"), 0644); err != nil {
		t.Fatalf("JSONファイルの作成に失敗しました: %v", err)
	}

	// テキストファイルも作成
	txtFile := filepath.Join(tempDir, "file.txt")
	if err := os.WriteFile(txtFile, []byte("text"), 0644); err != nil {
		t.Fatalf("テキストファイルの作成に失敗しました: %v", err)
	}

	// テストケース
	tests := []struct {
		name      string
		dirPath   string
		recursive bool
		setup     func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository)
		expected  []string
		wantErr   bool
	}{
		{
			name:      "非再帰的な検索",
			dirPath:   tempDir,
			recursive: false,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					FindFilesByExtFunc: func(dirPath, ext string) ([]string, error) {
						if dirPath == tempDir && ext == ".json" {
							return []string{jsonFile1}, nil
						}
						return nil, errors.New("unexpected arguments")
					},
				}
				iso8601Repo := &ISO8601RepositoryMock{}
				return fileRepo, iso8601Repo
			},
			expected: []string{jsonFile1},
			wantErr:  false,
		},
		{
			name:      "再帰的な検索",
			dirPath:   tempDir,
			recursive: true,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					ReadDirFunc: func(dirPath string) ([]*models.DirEntry, error) {
						if dirPath == tempDir {
							return []*models.DirEntry{
								{Name: "file1.json", Path: jsonFile1, IsDir: false},
								{Name: "file.txt", Path: txtFile, IsDir: false},
								{Name: "subdir", Path: subDir, IsDir: true},
							}, nil
						} else if dirPath == subDir {
							return []*models.DirEntry{
								{Name: "file2.json", Path: jsonFile2, IsDir: false},
							}, nil
						}
						return nil, errors.New("unexpected directory")
					},
				}
				iso8601Repo := &ISO8601RepositoryMock{}
				return fileRepo, iso8601Repo
			},
			expected: []string{jsonFile1, jsonFile2},
			wantErr:  false,
		},
		{
			name:      "ディレクトリ読み込みエラー",
			dirPath:   "invalid_dir",
			recursive: true,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					ReadDirFunc: func(dirPath string) ([]*models.DirEntry, error) {
						return nil, errors.New("directory read error")
					},
				}
				iso8601Repo := &ISO8601RepositoryMock{}
				return fileRepo, iso8601Repo
			},
			expected: nil,
			wantErr:  true,
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックをセットアップ
			fileRepo, iso8601Repo := tt.setup(t)

			// テスト対象のインスタンスを作成
			repo := NewJSONRepository(fileRepo, iso8601Repo)

			// テスト実行
			got, err := repo.FindJSONFiles(tt.dirPath, tt.recursive)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("FindJSONFiles() エラー = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// エラーが期待される場合は、ここで終了
			if tt.wantErr {
				return
			}

			// 結果の検証
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("FindJSONFiles() = %v, want = %v", got, tt.expected)
			}
		})
	}
}

// TestJSONRepositoryImpl_ProcessJSONData は ProcessJSONData メソッドをテストします
func TestJSONRepositoryImpl_ProcessJSONData(t *testing.T) {
	// テストケース
	tests := []struct {
		name       string
		data       interface{}
		targetKey  string
		setup      func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository)
		expected   interface{}
		wantChange bool
	}{
		{
			name: "ISO8601文字列を含むオブジェクト",
			data: map[string]interface{}{
				"id":        1,
				"name":      "test",
				"createdAt": "2023-04-01T12:34:56Z",
			},
			targetKey: "createdAt",
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{}
				iso8601Repo := &ISO8601RepositoryMock{
					ParseISO8601Func: func(dateStr string) (int64, error) {
						if dateStr == "2023-04-01T12:34:56Z" {
							return 1680351296, nil // 2023-04-01T12:34:56Z のUNIXタイムスタンプ
						}
						return 0, errors.New("unexpected date string")
					},
				}
				return fileRepo, iso8601Repo
			},
			expected: map[string]interface{}{
				"id":        1,
				"name":      "test",
				"createdAt": int64(1680351296),
			},
			wantChange: true,
		},
		{
			name: "ネストされたオブジェクト内のISO8601文字列",
			data: map[string]interface{}{
				"id":   1,
				"name": "test",
				"metadata": map[string]interface{}{
					"createdAt": "2023-04-01T12:34:56Z",
					"updatedAt": "2023-05-01T12:34:56Z",
				},
			},
			targetKey: "createdAt",
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{}
				iso8601Repo := &ISO8601RepositoryMock{
					ParseISO8601Func: func(dateStr string) (int64, error) {
						if dateStr == "2023-04-01T12:34:56Z" {
							return 1680351296, nil // 2023-04-01T12:34:56Z のUNIXタイムスタンプ
						}
						return 0, errors.New("unexpected date string")
					},
				}
				return fileRepo, iso8601Repo
			},
			expected: map[string]interface{}{
				"id":   1,
				"name": "test",
				"metadata": map[string]interface{}{
					"createdAt": int64(1680351296),
					"updatedAt": "2023-05-01T12:34:56Z",
				},
			},
			wantChange: true,
		},
		{
			name: "配列内のオブジェクトのISO8601文字列",
			data: map[string]interface{}{
				"id":   1,
				"name": "test",
				"items": []interface{}{
					map[string]interface{}{
						"id":        101,
						"createdAt": "2023-04-01T12:34:56Z",
					},
					map[string]interface{}{
						"id":        102,
						"createdAt": "2023-05-01T12:34:56Z",
					},
				},
			},
			targetKey: "createdAt",
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{}
				iso8601Repo := &ISO8601RepositoryMock{
					ParseISO8601Func: func(dateStr string) (int64, error) {
						if dateStr == "2023-04-01T12:34:56Z" {
							return 1680351296, nil // 2023-04-01T12:34:56Z のUNIXタイムスタンプ
						} else if dateStr == "2023-05-01T12:34:56Z" {
							return 1682943296, nil // 2023-05-01T12:34:56Z のUNIXタイムスタンプ
						}
						return 0, errors.New("unexpected date string")
					},
				}
				return fileRepo, iso8601Repo
			},
			expected: map[string]interface{}{
				"id":   1,
				"name": "test",
				"items": []interface{}{
					map[string]interface{}{
						"id":        101,
						"createdAt": int64(1680351296),
					},
					map[string]interface{}{
						"id":        102,
						"createdAt": int64(1682943296),
					},
				},
			},
			wantChange: true,
		},
		{
			name: "ISO8601文字列が含まれていない場合",
			data: map[string]interface{}{
				"id":   1,
				"name": "test",
				"date": "2023/04/01", // ISO8601形式ではない
			},
			targetKey: "date",
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{}
				iso8601Repo := &ISO8601RepositoryMock{
					ParseISO8601Func: func(dateStr string) (int64, error) {
						return 0, errors.New("invalid date format")
					},
				}
				return fileRepo, iso8601Repo
			},
			expected: map[string]interface{}{
				"id":   1,
				"name": "test",
				"date": "2023/04/01",
			},
			wantChange: false,
		},
		{
			name:      "プリミティブ型の場合",
			data:      "not an object",
			targetKey: "createdAt",
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{}
				iso8601Repo := &ISO8601RepositoryMock{}
				return fileRepo, iso8601Repo
			},
			expected:   "not an object",
			wantChange: false,
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックをセットアップ
			fileRepo, iso8601Repo := tt.setup(t)

			// テスト対象のインスタンスを作成
			repo := NewJSONRepository(fileRepo, iso8601Repo)

			// テスト実行
			got, changed := repo.ProcessJSONData(tt.data, tt.targetKey)

			// 変更フラグの検証
			if changed != tt.wantChange {
				t.Errorf("ProcessJSONData() changed = %v, wantChange = %v", changed, tt.wantChange)
			}

			// 結果の検証
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ProcessJSONData() = %v, want = %v", got, tt.expected)
			}
		})
	}
}

// TestJSONRepositoryImpl_ConvertFile は ConvertFile メソッドをテストします
func TestJSONRepositoryImpl_ConvertFile(t *testing.T) {
	// テスト用のJSONデータ
	testJSON := map[string]interface{}{
		"id":        1,
		"name":      "test",
		"created_at": "2023-04-01T12:34:56Z",
	}

	// テストケース
	tests := []struct {
		name     string
		filePath string
		key      string
		dryRun   bool
		setup    func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository)
		expected bool
		wantErr  bool
	}{
		{
			name:     "正常な変換（実際に書き込む）",
			filePath: "test.json",
			key:      "created_at",
			dryRun:   false,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					ReadJSONFileFunc: func(path string) (interface{}, error) {
						if path == "test.json" {
							return testJSON, nil
						}
						return nil, errors.New("unexpected file path")
					},
					WriteFileFunc: func(path string, content *models.FileContent) error {
						if path != "test.json" {
							return errors.New("unexpected file path")
						}
						// 書き込み成功を返す（内容の検証はスキップ）
						return nil
					},
				}
				iso8601Repo := &ISO8601RepositoryMock{
					ParseISO8601Func: func(dateStr string) (int64, error) {
						if dateStr == "2023-04-01T12:34:56Z" {
							return 1680351296, nil
						}
						return 0, errors.New("unexpected date string")
					},
				}
				return fileRepo, iso8601Repo
			},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "正常な変換（ドライラン）",
			filePath: "test.json",
			key:      "created_at",
			dryRun:   true,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					ReadJSONFileFunc: func(path string) (interface{}, error) {
						if path == "test.json" {
							return testJSON, nil
						}
						return nil, errors.New("unexpected file path")
					},
					// ドライランでは書き込みは行われないので、WriteFileFuncは設定しない
				}
				iso8601Repo := &ISO8601RepositoryMock{
					ParseISO8601Func: func(dateStr string) (int64, error) {
						if dateStr == "2023-04-01T12:34:56Z" {
							return 1680351296, nil
						}
						return 0, errors.New("unexpected date string")
					},
				}
				return fileRepo, iso8601Repo
			},
			expected: false, // ドライランでは変換は行われない
			wantErr:  false,
		},
		{
			name:     "JSONファイル読み込みエラー",
			filePath: "invalid.json",
			key:      "created_at",
			dryRun:   false,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					ReadJSONFileFunc: func(path string) (interface{}, error) {
						return nil, errors.New("file read error")
					},
				}
				iso8601Repo := &ISO8601RepositoryMock{}
				return fileRepo, iso8601Repo
			},
			expected: false,
			wantErr:  true,
		},
		{
			name:     "JSONデータが期待した型でない",
			filePath: "invalid_type.json",
			key:      "created_at",
			dryRun:   false,
			setup: func(t *testing.T) (domainRepo.FileRepository, domainRepo.ISO8601Repository) {
				fileRepo := &FileRepositoryMock{
					ReadJSONFileFunc: func(path string) (interface{}, error) {
						return "not an object", nil
					},
				}
				iso8601Repo := &ISO8601RepositoryMock{}
				return fileRepo, iso8601Repo
			},
			expected: false,
			wantErr:  true,
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// モックをセットアップ
			fileRepo, iso8601Repo := tt.setup(t)

			// テスト対象のインスタンスを作成
			repo := NewJSONRepository(fileRepo, iso8601Repo)

			// テスト実行
			got, err := repo.ConvertFile(tt.filePath, tt.key, tt.dryRun)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertFile() エラー = %v, wantErr = %v", err, tt.wantErr)
				return
			}

			// 結果の検証
			if got != tt.expected {
				t.Errorf("ConvertFile() = %v, want = %v", got, tt.expected)
			}
		})
	}
}
