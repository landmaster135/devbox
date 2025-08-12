package usecases

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// RealFileTestMockFileOperator は実ファイルテスト用のFileOperatorインターフェースのモック実装です
type RealFileTestMockFileOperator struct {
	ReadFileFunc  func(filename string) ([]byte, error)
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
	MkdirAllFunc  func(path string, perm os.FileMode) error
	WalkDirFunc   func(root string, fn fs.WalkDirFunc) error
	StatFunc      func(name string) (os.FileInfo, error)
	testDataFiles map[string][]byte
	outputFiles   map[string][]byte
	createdDirs   []string
}

// ReadFile はファイルを読み込みます
func (m *RealFileTestMockFileOperator) ReadFile(filename string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(filename)
	}
	if data, exists := m.testDataFiles[filename]; exists {
		return data, nil
	}
	return nil, fmt.Errorf("ファイルが見つかりません: %s", filename)
}

// WriteFile はファイルに書き込みます
func (m *RealFileTestMockFileOperator) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	if m.outputFiles == nil {
		m.outputFiles = make(map[string][]byte)
	}
	m.outputFiles[filename] = data
	return nil
}

// MkdirAll はディレクトリを作成します
func (m *RealFileTestMockFileOperator) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	m.createdDirs = append(m.createdDirs, path)
	return nil
}

// WalkDir はディレクトリを走査します
func (m *RealFileTestMockFileOperator) WalkDir(root string, fn fs.WalkDirFunc) error {
	if m.WalkDirFunc != nil {
		return m.WalkDirFunc(root, fn)
	}

	// テストデータファイルを走査
	for filePath := range m.testDataFiles {
		if strings.HasPrefix(filePath, root) && strings.HasSuffix(strings.ToLower(filePath), ".md") {
			// ファイル情報を作成
			info := &realFileTestMockFileInfo{
				name:  filepath.Base(filePath),
				isDir: false,
			}
			dirEntry := &realFileTestMockDirEntry{
				name:  filepath.Base(filePath),
				isDir: false,
				info:  info,
			}

			if err := fn(filePath, dirEntry, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

// Stat はファイル情報を取得します
func (m *RealFileTestMockFileOperator) Stat(name string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}

	// ディレクトリの場合
	if name == "/test/src" || name == "/test/dest" {
		return &realFileTestMockFileInfo{
			name:  filepath.Base(name),
			isDir: true,
		}, nil
	}

	// ファイルの場合
	if _, exists := m.testDataFiles[name]; exists {
		return &realFileTestMockFileInfo{
			name:  filepath.Base(name),
			isDir: false,
		}, nil
	}

	return nil, os.ErrNotExist
}

// realFileTestMockFileInfo はos.FileInfoのモック実装です
type realFileTestMockFileInfo struct {
	name  string
	isDir bool
}

func (m *realFileTestMockFileInfo) Name() string       { return m.name }
func (m *realFileTestMockFileInfo) Size() int64        { return 0 }
func (m *realFileTestMockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *realFileTestMockFileInfo) ModTime() time.Time { return time.Now() }
func (m *realFileTestMockFileInfo) IsDir() bool        { return m.isDir }
func (m *realFileTestMockFileInfo) Sys() interface{}   { return nil }

// realFileTestMockDirEntry はfs.DirEntryのモック実装です
type realFileTestMockDirEntry struct {
	name  string
	isDir bool
	info  os.FileInfo
}

func (m *realFileTestMockDirEntry) Name() string               { return m.name }
func (m *realFileTestMockDirEntry) IsDir() bool                { return m.isDir }
func (m *realFileTestMockDirEntry) Type() fs.FileMode          { return 0644 }
func (m *realFileTestMockDirEntry) Info() (os.FileInfo, error) { return m.info, nil }

func TestService_ExtractBlogContent_WithRealData(t *testing.T) {
	// テストデータファイルのパス
	testDataPaths := []string{
		"./test_data/org/sample_python_01.md",
		"./test_data/org/sample_speaker_01.md",
		"./test_data/org/sample_youtube_01.md",
	}

	// 実際のテストデータファイルを読み込み
	testDataFiles := make(map[string][]byte)
	for _, path := range testDataPaths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("テストデータファイルの読み込みに失敗しました: %s (%v)", path, err)
		}

		// テスト用のパスに変換
		testPath := filepath.Join("/test/src", filepath.Base(path))
		testDataFiles[testPath] = content
	}

	// モックFileOperatorを作成
	mockFileOperator := &RealFileTestMockFileOperator{
		testDataFiles: testDataFiles,
	}

	// テスト対象のサービスを作成
	service := NewServiceWithFileOperator(mockFileOperator)

	tests := []struct {
		name            string
		srcDir          string
		destDir         string
		wantExtracted   int
		wantErr         bool
		wantResultMatch string
	}{
		{
			name:            "正常系: 3つの実ファイルからコンテンツを抽出",
			srcDir:          "/test/src",
			destDir:         "/test/dest",
			wantExtracted:   3,
			wantErr:         false,
			wantResultMatch: "処理完了: 3件のファイルからコンテンツを抽出しました。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			result, err := service.ExtractBlogContent(tt.srcDir, tt.destDir)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractBlogContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 結果メッセージの検証
			if !tt.wantErr && !strings.Contains(result, tt.wantResultMatch) {
				t.Errorf("ExtractBlogContent() result = %v, want to contain %v", result, tt.wantResultMatch)
			}

			// 出力ファイルの検証
			if !tt.wantErr {
				// 出力ディレクトリが作成されたか確認
				dirCreated := false
				for _, dir := range mockFileOperator.createdDirs {
					if dir == tt.destDir {
						dirCreated = true
						break
					}
				}
				if !dirCreated {
					t.Errorf("出力ディレクトリが作成されませんでした: %s", tt.destDir)
				}

				// 各テストファイルに対応する出力ファイルが作成されたか確認
				expectedOutputFiles := []string{
					filepath.Join(tt.destDir, "sample_python_01.md"),
					filepath.Join(tt.destDir, "sample_speaker_01.md"),
					filepath.Join(tt.destDir, "sample_youtube_01.md"),
				}

				for _, expectedFile := range expectedOutputFiles {
					if _, exists := mockFileOperator.outputFiles[expectedFile]; !exists {
						t.Errorf("期待される出力ファイルが作成されませんでした: %s", expectedFile)
					} else {
						// 出力ファイルの内容を検証
						outputContent := string(mockFileOperator.outputFiles[expectedFile])

						// "# Content" マーカーが含まれているか確認
						if !strings.Contains(outputContent, "# Content") {
							t.Errorf("出力ファイルに '# Content' マーカーが含まれていません: %s", expectedFile)
						}

						// "## はじまり" マーカーが含まれているか確認
						if !strings.Contains(outputContent, "## はじまり") {
							t.Errorf("出力ファイルに '## はじまり' マーカーが含まれていません: %s", expectedFile)
						}

						// ファイル固有の内容が含まれているか確認
						fileName := filepath.Base(expectedFile)
						switch fileName {
						case "sample_python_01.md":
							containsSample := strings.Contains(outputContent, "サンプルちゃん")
							if !containsSample {
								maxLen := 500
								if len(outputContent) < maxLen {
									maxLen = len(outputContent)
								}
								t.Errorf("sample_python_01.mdの固有内容「サンプルちゃん」が含まれていません: %s\n実際の内容: %s", expectedFile, outputContent[:maxLen])
							}

							containsTabSize := strings.Contains(outputContent, "Tab Size")
							if !containsTabSize {
								maxLen := 500
								if len(outputContent) < maxLen {
									maxLen = len(outputContent)
								}
								t.Errorf("sample_python_01.mdの固有内容「Tab Size」が含まれていません: %s\n実際の内容: %s", expectedFile, outputContent[:maxLen])
							}
						case "sample_speaker_01.md":
							if !strings.Contains(outputContent, "スピーカーが欲しいぞ～") {
								t.Errorf("sample_speaker_01.mdの固有内容が含まれていません: %s", expectedFile)
							}
							if !strings.Contains(outputContent, "一体型スピーカー") {
								t.Errorf("sample_speaker_01.mdの固有内容が含まれていません: %s", expectedFile)
							}
						case "sample_youtube_01.md":
							if !strings.Contains(outputContent, "先月の動画の更新状況です。") {
								t.Errorf("sample_youtube_01.mdの固有内容が含まれていません: %s", expectedFile)
							}
							if !strings.Contains(outputContent, "投稿した動画の数：8") {
								t.Errorf("sample_youtube_01.mdの固有内容が含まれていません: %s", expectedFile)
							}
						}
					}
				}
			}
		})
	}
}

func TestService_ExtractContentFromFile_WithRealData(t *testing.T) {
	// テストデータファイルのパス
	testDataPath := "./test_data/org/sample_python_01.md"

	// 実際のテストデータファイルを読み込み
	content, err := os.ReadFile(testDataPath)
	if err != nil {
		t.Fatalf("テストデータファイルの読み込みに失敗しました: %v", err)
	}

	// モックFileOperatorを作成
	mockFileOperator := &RealFileTestMockFileOperator{
		testDataFiles: map[string][]byte{
			"/test/sample_python_01.md": content,
		},
	}

	// テスト対象のサービスを作成
	service := NewServiceWithFileOperator(mockFileOperator)

	tests := []struct {
		name     string
		filePath string
		destDir  string
		wantErr  bool
	}{
		{
			name:     "正常系: sample_python_01.mdからコンテンツを抽出",
			filePath: "/test/sample_python_01.md",
			destDir:  "/test/dest",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト実行
			err := service.extractContentFromFile(tt.filePath, tt.destDir)

			// エラーの検証
			if (err != nil) != tt.wantErr {
				t.Errorf("extractContentFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 出力ファイルの検証
			if !tt.wantErr {
				expectedOutputFile := filepath.Join(tt.destDir, "sample_python_01.md")
				if _, exists := mockFileOperator.outputFiles[expectedOutputFile]; !exists {
					t.Errorf("期待される出力ファイルが作成されませんでした: %s", expectedOutputFile)
				} else {
					// 出力ファイルの内容を検証
					outputContent := string(mockFileOperator.outputFiles[expectedOutputFile])

					// "# Content" マーカーが含まれているか確認
					if !strings.Contains(outputContent, "# Content") {
						t.Errorf("出力ファイルに '# Content' マーカーが含まれていません")
					}

					// "## はじまり" マーカーが含まれているか確認
					if !strings.Contains(outputContent, "## はじまり") {
						t.Errorf("出力ファイルに '## はじまり' マーカーが含まれていません")
					}

					// 元のファイルの "# Content" より前の部分が含まれていないか確認
					if strings.Contains(outputContent, "priority: 5") {
						t.Errorf("出力ファイルに '# Content' より前の内容が含まれています")
					}

					// 抽出されるべき内容が含まれているか確認
					if !strings.Contains(outputContent, "サンプルちゃん") {
						t.Errorf("出力ファイルに期待される内容が含まれていません")
					}
				}
			}
		})
	}
}

func TestService_ExtractContentAfterMarker_WithRealData(t *testing.T) {
	// テストデータファイルのパス
	testDataPaths := []string{
		"./test_data/org/sample_python_01.md",
		"./test_data/org/sample_speaker_01.md",
		"./test_data/org/sample_youtube_01.md",
	}

	// モックFileOperatorを作成（この場合は使用しない）
	mockFileOperator := &RealFileTestMockFileOperator{}

	// テスト対象のサービスを作成
	service := NewServiceWithFileOperator(mockFileOperator)

	for _, testDataPath := range testDataPaths {
		t.Run(fmt.Sprintf("正常系: %sからマーカー以降のコンテンツを抽出", filepath.Base(testDataPath)), func(t *testing.T) {
			// 実際のテストデータファイルを読み込み
			content, err := os.ReadFile(testDataPath)
			if err != nil {
				t.Fatalf("テストデータファイルの読み込みに失敗しました: %v", err)
			}

			// テスト実行
			extractedContent, err := service.extractContentAfterMarker(string(content))

			// エラーの検証
			if err != nil {
				t.Errorf("extractContentAfterMarker() error = %v", err)
				return
			}

			// 抽出されたコンテンツの検証
			if !strings.Contains(extractedContent, "# Content") {
				t.Errorf("抽出されたコンテンツに '# Content' マーカーが含まれていません")
			}

			if !strings.Contains(extractedContent, "## はじまり") {
				t.Errorf("抽出されたコンテンツに '## はじまり' マーカーが含まれていません")
			}

			// マーカーより前の内容が含まれていないか確認
			if strings.Contains(extractedContent, "priority: 5") {
				t.Errorf("抽出されたコンテンツに '# Content' より前の内容が含まれています")
			}

			// ファイル固有の内容が含まれているか確認
			fileName := filepath.Base(testDataPath)
			switch fileName {
			case "sample_python_01.md":
				if !strings.Contains(extractedContent, "サンプルちゃん") {
					t.Errorf("sample_python_01.mdの固有内容が含まれていません")
				}
				if !strings.Contains(extractedContent, "Tab Size") {
					t.Errorf("sample_python_01.mdの固有内容が含まれていません")
				}
			case "sample_speaker_01.md":
				if !strings.Contains(extractedContent, "スピーカーが欲しいぞ～") {
					t.Errorf("sample_speaker_01.mdの固有内容が含まれていません")
				}
				if !strings.Contains(extractedContent, "一体型スピーカー") {
					t.Errorf("sample_speaker_01.mdの固有内容が含まれていません")
				}
			case "sample_youtube_01.md":
				if !strings.Contains(extractedContent, "先月の動画の更新状況です。") {
					t.Errorf("sample_youtube_01.mdの固有内容が含まれていません")
				}
				if !strings.Contains(extractedContent, "投稿した動画の数：8") {
					t.Errorf("sample_youtube_01.mdの固有内容が含まれていません")
				}
			}
		})
	}
}
