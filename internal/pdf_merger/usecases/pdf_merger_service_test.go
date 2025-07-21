package usecases

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestPDFMergerService_Process(t *testing.T) {
	tests := []struct {
		name           string
		opts           PDFMergerOptions
		expectError    bool
		expectedOutput string
	}{
		{
			name: "画像抽出処理_正常系",
			opts: PDFMergerOptions{
				Extract:     "test_data/org/sample_01_01.pdf",
				OutputDir:   "test_data/tmp/extract_test",
				ImageFormat: "jpg",
				StartPage:   0,
				EndPage:     0,
			},
			expectError:    false,
			expectedOutput: "PDF画像抽出を開始します",
		},
		{
			name: "画像抽出処理_出力ディレクトリ未指定",
			opts: PDFMergerOptions{
				Extract:     "test_data/org/sample_01_01.pdf",
				OutputDir:   "",
				ImageFormat: "jpg",
			},
			expectError:    true,
			expectedOutput: "",
		},
		{
			name: "PDF作成処理_正常系",
			opts: PDFMergerOptions{
				Dir: "test_data/org",
				Out: "test_data/tmp/created_test.pdf",
			},
			expectError:    false,
			expectedOutput: "検出した画像",
		},
		{
			name: "PDF作成処理_画像なし",
			opts: PDFMergerOptions{
				Dir: "test_data/tmp/empty",
				Out: "test_data/tmp/empty_test.pdf",
			},
			expectError:    true,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用ディレクトリの準備
			if err := os.MkdirAll("test_data/tmp/extract_test", 0755); err != nil {
				t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
			}
			if err := os.MkdirAll("test_data/tmp/empty", 0755); err != nil {
				t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
			}
			defer func() {
				os.RemoveAll("test_data/tmp/extract_test")
				os.RemoveAll("test_data/tmp/empty")
				os.Remove("test_data/tmp/created_test.pdf")
				os.Remove("test_data/tmp/empty_test.pdf")
			}()

			service := NewPDFMergerService()
			var stdout, stderr bytes.Buffer

			err := service.Process(tt.opts, &stdout, &stderr)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if tt.expectedOutput != "" && !strings.Contains(stdout.String(), tt.expectedOutput) {
					t.Errorf("期待される出力が含まれていません。期待: %s, 実際: %s", tt.expectedOutput, stdout.String())
				}
			}
		})
	}
}

func TestPDFMergerService_handleImageExtraction(t *testing.T) {
	tests := []struct {
		name        string
		opts        PDFMergerOptions
		expectError bool
	}{
		{
			name: "正常系_全ページ抽出",
			opts: PDFMergerOptions{
				Extract:     "test_data/org/sample_01_01.pdf",
				OutputDir:   "test_data/tmp/extract_handle_test",
				ImageFormat: "jpg",
				StartPage:   0,
				EndPage:     0,
			},
			expectError: false,
		},
		{
			name: "異常系_PDFファイル不存在",
			opts: PDFMergerOptions{
				Extract:     "test_data/org/nonexistent.pdf",
				OutputDir:   "test_data/tmp/extract_handle_test",
				ImageFormat: "jpg",
			},
			expectError: true,
		},
		{
			name: "異常系_出力ディレクトリ未指定",
			opts: PDFMergerOptions{
				Extract:     "test_data/org/sample_01_01.pdf",
				OutputDir:   "",
				ImageFormat: "jpg",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用ディレクトリの準備
			if err := os.MkdirAll("test_data/tmp/extract_handle_test", 0755); err != nil {
				t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
			}
			defer os.RemoveAll("test_data/tmp/extract_handle_test")

			service := NewPDFMergerService()
			var stdout, stderr bytes.Buffer

			err := service.handleImageExtraction(tt.opts, &stdout, &stderr)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
			}
		})
	}
}

func TestPDFMergerService_handlePDFCreation(t *testing.T) {
	tests := []struct {
		name        string
		opts        PDFMergerOptions
		expectError bool
	}{
		{
			name: "正常系_新規PDF作成",
			opts: PDFMergerOptions{
				Dir: "test_data/org",
				Out: "test_data/tmp/handle_created_test.pdf",
			},
			expectError: false,
		},
		{
			name: "正常系_既存PDFに追加",
			opts: PDFMergerOptions{
				Dir: "test_data/org",
				Out: "test_data/tmp/handle_merged_test.pdf",
				Add: "test_data/org/sample_01_01.pdf",
			},
			expectError: false,
		},
		{
			name: "異常系_既存PDFファイル不存在",
			opts: PDFMergerOptions{
				Dir: "test_data/org",
				Out: "test_data/tmp/handle_error_test.pdf",
				Add: "test_data/org/nonexistent.pdf",
			},
			expectError: true,
		},
		{
			name: "異常系_画像ファイルなし",
			opts: PDFMergerOptions{
				Dir: "test_data/tmp/empty_handle",
				Out: "test_data/tmp/handle_empty_test.pdf",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// テスト用ディレクトリの準備
			if err := os.MkdirAll("test_data/tmp/empty_handle", 0755); err != nil {
				t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
			}
			defer func() {
				os.RemoveAll("test_data/tmp/empty_handle")
				os.Remove("test_data/tmp/handle_created_test.pdf")
				os.Remove("test_data/tmp/handle_merged_test.pdf")
				os.Remove("test_data/tmp/handle_error_test.pdf")
				os.Remove("test_data/tmp/handle_empty_test.pdf")
			}()

			service := NewPDFMergerService()
			var stdout, stderr bytes.Buffer

			err := service.handlePDFCreation(tt.opts, &stdout, &stderr)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				// 出力ファイルが作成されているか確認
				if _, err := os.Stat(tt.opts.Out); os.IsNotExist(err) {
					t.Errorf("出力ファイルが作成されていません: %s", tt.opts.Out)
				}
			}
		})
	}
}

func TestNewPDFMergerService(t *testing.T) {
	service := NewPDFMergerService()
	if service == nil {
		t.Error("NewPDFMergerService() returned nil")
	}
}

func TestPDFMergerOptions_Validation(t *testing.T) {
	tests := []struct {
		name        string
		opts        PDFMergerOptions
		description string
	}{
		{
			name: "画像抽出オプション",
			opts: PDFMergerOptions{
				Extract:     "test.pdf",
				OutputDir:   "output",
				ImageFormat: "jpg",
				StartPage:   1,
				EndPage:     5,
			},
			description: "画像抽出用のオプションが正しく設定される",
		},
		{
			name: "PDF作成オプション",
			opts: PDFMergerOptions{
				Dir: "images",
				Out: "output.pdf",
				Add: "existing.pdf",
			},
			description: "PDF作成用のオプションが正しく設定される",
		},
		{
			name:        "空のオプション",
			opts:        PDFMergerOptions{},
			description: "空のオプションでも構造体が作成される",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// オプション構造体の各フィールドが期待通りに設定されているか確認
			if tt.opts.Extract != "" && tt.opts.Extract != "test.pdf" && tt.name == "画像抽出オプション" {
				t.Errorf("Extract field mismatch")
			}
			if tt.opts.Dir != "" && tt.opts.Dir != "images" && tt.name == "PDF作成オプション" {
				t.Errorf("Dir field mismatch")
			}
			// 構造体が正常に作成されることを確認
			t.Logf("%s: %+v", tt.description, tt.opts)
		})
	}
}

func TestPDFMergerService_Integration(t *testing.T) {
	t.Run("完全な画像抽出フロー", func(t *testing.T) {
		// テスト用ディレクトリの準備
		outputDir := "test_data/tmp/integration_extract"
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			t.Fatalf("テスト用ディレクトリの作成に失敗: %v", err)
		}
		defer os.RemoveAll(outputDir)

		opts := PDFMergerOptions{
			Extract:     "test_data/org/sample_01_01.pdf",
			OutputDir:   outputDir,
			ImageFormat: "jpg",
			StartPage:   0,
			EndPage:     0,
		}

		service := NewPDFMergerService()
		var stdout, stderr bytes.Buffer

		err := service.Process(opts, &stdout, &stderr)
		if err != nil {
			t.Fatalf("画像抽出処理でエラーが発生: %v", err)
		}

		// 抽出された画像ファイルが存在するか確認
		files, err := os.ReadDir(outputDir)
		if err != nil {
			t.Fatalf("出力ディレクトリの読み取りに失敗: %v", err)
		}

		if len(files) == 0 {
			t.Error("画像ファイルが抽出されていません")
		}

		// 出力メッセージの確認
		output := stdout.String()
		if !strings.Contains(output, "PDF画像抽出を開始します") {
			t.Error("期待される開始メッセージが出力されていません")
		}
		if !strings.Contains(output, "画像抽出が完了しました") {
			t.Error("期待される完了メッセージが出力されていません")
		}
	})

	t.Run("完全なPDF作成フロー", func(t *testing.T) {
		outputFile := "test_data/tmp/integration_created.pdf"
		defer os.Remove(outputFile)

		opts := PDFMergerOptions{
			Dir: "test_data/org",
			Out: outputFile,
		}

		service := NewPDFMergerService()
		var stdout, stderr bytes.Buffer

		err := service.Process(opts, &stdout, &stderr)
		if err != nil {
			t.Fatalf("PDF作成処理でエラーが発生: %v", err)
		}

		// 出力ファイルが作成されているか確認
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Error("PDFファイルが作成されていません")
		}

		// 出力メッセージの確認
		output := stdout.String()
		if !strings.Contains(output, "検出した画像") {
			t.Error("期待される画像検出メッセージが出力されていません")
		}
		if !strings.Contains(output, "PDF を生成しました") {
			t.Error("期待される完了メッセージが出力されていません")
		}
	})
}
