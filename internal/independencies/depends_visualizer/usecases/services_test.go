package usecases

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/analyzer"
	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/config"
)

// TestApp_Run_Normal はAppのRunメソッドの正常系テストです
func TestApp_Run_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "depends_visualizer_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のGoファイルを作成
	testFilePath := filepath.Join(tempDir, "test.go")
	testContent := `package test

func main() {
	hello()
	world()
}

func hello() {
	fmt.Println("Hello")
}

func world() {
	fmt.Println("World")
}
`
	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// テスト用の設定
	cfg := &config.AppConfig{
		SourceFile: testFilePath,
		Extension:  ".go",
		Format:     "mermaid",
	}

	// テスト対象のインスタンスを作成
	app := NewApp(cfg)

	// 出力先のバッファを用意
	var stdout, stderr bytes.Buffer

	// テスト実行
	exitCode := app.Run(&stdout, &stderr)

	// 検証
	if exitCode != ExitCodeOK {
		t.Errorf("期待する終了コードは %d ですが、実際は %d でした", ExitCodeOK, exitCode)
	}

	// 標準出力の内容を検証
	output := stdout.String()
	if !strings.Contains(output, "```mermaid") {
		t.Errorf("出力にmermaidヘッダーが含まれていません: %s", output)
	}
	if !strings.Contains(output, "main") {
		t.Errorf("出力にmain関数が含まれていません: %s", output)
	}
	if !strings.Contains(output, "hello") {
		t.Errorf("出力にhello関数が含まれていません: %s", output)
	}
	if !strings.Contains(output, "world") {
		t.Errorf("出力にworld関数が含まれていません: %s", output)
	}
}

// TestApp_CollectFiles_Normal はcollectFilesメソッドの正常系テストです
func TestApp_CollectFiles_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "depends_visualizer_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	testFiles := []struct {
		name    string
		content string
	}{
		{"test1.go", "package test\n\nfunc main() {}\n"},
		{"test2.go", "package test\n\nfunc hello() {}\n"},
		{"test.txt", "This is not a Go file"},
		{"subdir/test3.go", "package test\n\nfunc world() {}\n"},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tempDir, tf.name)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("ディレクトリの作成に失敗しました: %v", err)
		}
		if err := os.WriteFile(path, []byte(tf.content), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	// 単一ファイルモードのテスト
	t.Run("SingleFileMode", func(t *testing.T) {
		cfg := &config.AppConfig{
			SourceFile: filepath.Join(tempDir, "test1.go"),
		}
		app := NewApp(cfg)

		files, err := app.collectFiles()
		if err != nil {
			t.Fatalf("ファイル収集に失敗しました: %v", err)
		}

		if len(files) != 1 {
			t.Errorf("期待するファイル数は1ですが、実際は%dでした", len(files))
		}

		if !strings.Contains(files[0], "test1.go") {
			t.Errorf("期待するファイル名はtest1.goですが、実際は%sでした", files[0])
		}

		if app.Config.Extension != ".go" {
			t.Errorf("期待する拡張子は.goですが、実際は%sでした", app.Config.Extension)
		}
	})

	// ディレクトリモード（非再帰的）のテスト
	t.Run("DirectoryModeNonRecursive", func(t *testing.T) {
		cfg := &config.AppConfig{
			Directory: tempDir,
			Extension: ".go",
			Recursive: false,
		}
		app := NewApp(cfg)

		files, err := app.collectFiles()
		if err != nil {
			t.Fatalf("ファイル収集に失敗しました: %v", err)
		}

		if len(files) != 2 { // test1.go と test2.go
			t.Errorf("期待するファイル数は2ですが、実際は%dでした", len(files))
		}
	})

	// ディレクトリモード（再帰的）のテスト
	t.Run("DirectoryModeRecursive", func(t *testing.T) {
		cfg := &config.AppConfig{
			Directory: tempDir,
			Extension: ".go",
			Recursive: true,
		}
		app := NewApp(cfg)

		files, err := app.collectFiles()
		if err != nil {
			t.Fatalf("ファイル収集に失敗しました: %v", err)
		}

		if len(files) != 3 { // test1.go, test2.go, subdir/test3.go
			t.Errorf("期待するファイル数は3ですが、実際は%dでした: %v", len(files), files)
		}
	})
}

// TestApp_RenderResults_Normal はrenderResultsメソッドの正常系テストです
func TestApp_RenderResults_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "depends_visualizer_test")
	if err != nil {
		t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用の出力ファイルパス
	outputPath := filepath.Join(tempDir, "output.md")

	// テスト用の解析結果
	results := []analyzer.AnalysisResult{
		{
			FilePath: "/path/to/test.go",
			Dependencies: map[string][]string{
				"main":  {"hello", "world"},
				"hello": {},
				"world": {"hello"},
			},
		},
	}

	testCases := []struct {
		name       string
		format     string
		outputPath string
		wantErr    bool
		contains   []string
	}{
		{
			name:     "Mermaid Format",
			format:   "mermaid",
			wantErr:  false,
			contains: []string{"```mermaid", "class", "main", "hello", "world"},
		},
		{
			name:     "MermaidFlowchart Format",
			format:   "mermaid-flowchart",
			wantErr:  false,
			contains: []string{"```mermaid", "flowchart", "main", "hello", "world"},
		},
		{
			name:     "PlantUML Format",
			format:   "plantuml",
			wantErr:  false,
			contains: []string{"@startuml", "main", "hello", "world"},
		},
		{
			name:     "DOT Format",
			format:   "dot",
			wantErr:  false,
			contains: []string{"digraph", "main", "hello", "world"},
		},
		{
			name:       "File Output",
			format:     "mermaid",
			outputPath: outputPath,
			wantErr:    false,
			contains:   []string{},
		},
		{
			name:     "Invalid Format",
			format:   "invalid",
			wantErr:  true,
			contains: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の設定
			cfg := &config.AppConfig{
				Format:     tc.format,
				OutputPath: tc.outputPath,
			}
			app := NewApp(cfg)

			// 出力先のバッファを用意
			var stdout bytes.Buffer

			// テスト実行
			err := app.renderResults(results, &stdout)

			// エラー検証
			if (err != nil) != tc.wantErr {
				t.Errorf("renderResults() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			// 出力ファイルの検証
			if tc.outputPath != "" && !tc.wantErr {
				if _, err := os.Stat(tc.outputPath); os.IsNotExist(err) {
					t.Errorf("出力ファイルが作成されていません: %s", tc.outputPath)
				}
				content, err := os.ReadFile(tc.outputPath)
				if err != nil {
					t.Fatalf("出力ファイルの読み込みに失敗しました: %v", err)
				}
				output := string(content)
				if !strings.Contains(output, "mermaid") {
					t.Errorf("出力ファイルに期待する内容が含まれていません")
				}
			} else if !tc.wantErr {
				// 標準出力の検証
				output := stdout.String()
				for _, s := range tc.contains {
					if !strings.Contains(output, s) {
						t.Errorf("出力に %s が含まれていません: %s", s, output)
					}
				}
			}
		})
	}
}

// TestIsSupportedExtension_Normal はisSupportedExtension関数の正常系テストです
func TestIsSupportedExtension_Normal(t *testing.T) {
	testCases := []struct {
		name     string
		ext      string
		expected bool
	}{
		{"Go", ".go", true},
		{"Python", ".py", true},
		{"JavaScript", ".js", true},
		{"Unsupported", ".txt", false},
		{"Empty", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isSupportedExtension(tc.ext)
			if result != tc.expected {
				t.Errorf("%sの場合、期待する結果は%vですが、実際は%vでした", tc.ext, tc.expected, result)
			}
		})
	}
}

// TestGetSupportedExtensions_Normal はgetSupportedExtensions関数の正常系テストです
func TestGetSupportedExtensions_Normal(t *testing.T) {
	extensions := getSupportedExtensions()

	// 期待される拡張子が含まれているか確認
	expectedExts := []string{".go", ".py", ".js"}
	for _, ext := range expectedExts {
		found := false
		for _, e := range extensions {
			if e == ext {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期待する拡張子 %s が結果に含まれていません", ext)
		}
	}
}
