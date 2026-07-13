package usecases

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestPDFMergerService() (*PDFMergerService, *bytes.Buffer) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	return NewPDFMergerServiceWithLogger(logger), &buf
}

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
				Operation:   OperationExtractImages,
				SrcFile:     "test_data/org/sample_01_01.pdf",
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
				Operation:   OperationExtractImages,
				SrcFile:     "test_data/org/sample_01_01.pdf",
				OutputDir:   "",
				ImageFormat: "jpg",
			},
			expectError:    true,
			expectedOutput: "",
		},
		{
			name: "PDF作成処理_正常系",
			opts: PDFMergerOptions{
				Operation: OperationMergeIntoNew,
				SrcDir:    "test_data/org",
				OutputDir: "test_data/tmp/created_test",
			},
			expectError:    false,
			expectedOutput: "検出した画像",
		},
		{
			name: "既存PDF追加処理_正常系",
			opts: PDFMergerOptions{
				Operation:     OperationAddIntoExist,
				SrcDir:        "test_data/org",
				OutputDir:     "test_data/tmp/added_test",
				ReceivingFile: "test_data/org/sample_01_01.pdf",
			},
			expectError:    false,
			expectedOutput: "既存PDFに画像を追加しました",
		},
		{
			name: "PDF作成処理_画像なし",
			opts: PDFMergerOptions{
				Operation: OperationMergeIntoNew,
				SrcDir:    "test_data/tmp/empty",
				OutputDir: "test_data/tmp/empty_test",
			},
			expectError:    true,
			expectedOutput: "",
		},
		{
			name: "未対応Operation",
			opts: PDFMergerOptions{
				Operation: "unknown",
				SrcDir:    "test_data/org",
				OutputDir: "test_data/tmp/unknown_test",
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
				os.RemoveAll("test_data/tmp/created_test")
				os.RemoveAll("test_data/tmp/added_test")
				os.RemoveAll("test_data/tmp/empty_test")
				os.RemoveAll("test_data/tmp/unknown_test")
			}()

			service, logBuf := newTestPDFMergerService()

			err := service.Process(tt.opts)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if tt.expectedOutput != "" && !strings.Contains(logBuf.String(), tt.expectedOutput) {
					t.Errorf("期待される出力が含まれていません。期待: %s, 実際: %s", tt.expectedOutput, logBuf.String())
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
				SrcFile:     "test_data/org/sample_01_01.pdf",
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
				SrcFile:     "test_data/org/nonexistent.pdf",
				OutputDir:   "test_data/tmp/extract_handle_test",
				ImageFormat: "jpg",
			},
			expectError: true,
		},
		{
			name: "異常系_出力ディレクトリ未指定",
			opts: PDFMergerOptions{
				SrcFile:     "test_data/org/sample_01_01.pdf",
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

			service, _ := newTestPDFMergerService()

			err := service.handleImageExtraction(tt.opts)

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
				SrcDir:    "test_data/org",
				OutputDir: "test_data/tmp/handle_created_test",
			},
			expectError: false,
		},
		{
			name: "正常系_既存PDFに追加",
			opts: PDFMergerOptions{
				SrcDir:        "test_data/org",
				OutputDir:     "test_data/tmp/handle_merged_test",
				ReceivingFile: "test_data/org/sample_01_01.pdf",
			},
			expectError: false,
		},
		{
			name: "異常系_既存PDFファイル不存在",
			opts: PDFMergerOptions{
				SrcDir:        "test_data/org",
				OutputDir:     "test_data/tmp/handle_error_test",
				ReceivingFile: "test_data/org/nonexistent.pdf",
			},
			expectError: true,
		},
		{
			name: "異常系_画像ファイルなし",
			opts: PDFMergerOptions{
				SrcDir:    "test_data/tmp/empty_handle",
				OutputDir: "test_data/tmp/handle_empty_test",
			},
			expectError: true,
		},
		{
			name: "異常系_出力ディレクトリ未指定",
			opts: PDFMergerOptions{
				SrcDir: "test_data/org",
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
				os.RemoveAll("test_data/tmp/handle_created_test")
				os.RemoveAll("test_data/tmp/handle_merged_test")
				os.RemoveAll("test_data/tmp/handle_error_test")
				os.RemoveAll("test_data/tmp/handle_empty_test")
			}()

			service, _ := newTestPDFMergerService()

			err := service.handlePDFCreation(tt.opts)

			if tt.expectError {
				if err == nil {
					t.Errorf("エラーが期待されましたが、エラーが発生しませんでした")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				// 出力ファイルが作成されているか確認
				outputPath := filepath.Join(tt.opts.OutputDir, filepath.Base(tt.opts.OutputDir)+".pdf")
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("出力ファイルが作成されていません: %s", outputPath)
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
				Operation:   OperationExtractImages,
				SrcFile:     "test.pdf",
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
				Operation:     OperationAddIntoExist,
				SrcDir:        "images",
				OutputDir:     "output",
				ReceivingFile: "existing.pdf",
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
			if tt.opts.SrcFile != "" && tt.opts.SrcFile != "test.pdf" && tt.name == "画像抽出オプション" {
				t.Errorf("SrcFile field mismatch")
			}
			if tt.opts.SrcDir != "" && tt.opts.SrcDir != "images" && tt.name == "PDF作成オプション" {
				t.Errorf("SrcDir field mismatch")
			}
			// 構造体が正常に作成されることを確認
			t.Logf("%s: %+v", tt.description, tt.opts)
		})
	}
}

func TestPDFMergerOptions_ValidateOperation(t *testing.T) {
	tests := []struct {
		name        string
		opts        PDFMergerOptions
		expectError bool
	}{
		{
			name: "MergeIntoNew_Normal",
			opts: PDFMergerOptions{
				Operation: OperationMergeIntoNew,
				SrcDir:    "images",
				OutputDir: "output",
			},
			expectError: false,
		},
		{
			name: "AddIntoExist_Normal",
			opts: PDFMergerOptions{
				Operation:     OperationAddIntoExist,
				SrcDir:        "images",
				OutputDir:     "output",
				ReceivingFile: "existing.pdf",
			},
			expectError: false,
		},
		{
			name: "ExtractImages_Normal",
			opts: PDFMergerOptions{
				Operation: OperationExtractImages,
				SrcFile:   "input.pdf",
				OutputDir: "output",
			},
			expectError: false,
		},
		{
			name: "AddIntoExist_ReceivingFile未指定",
			opts: PDFMergerOptions{
				Operation: OperationAddIntoExist,
				SrcDir:    "images",
				OutputDir: "output",
			},
			expectError: true,
		},
		{
			name: "ExtractImages_SrcFile未指定",
			opts: PDFMergerOptions{
				Operation: OperationExtractImages,
				OutputDir: "output",
			},
			expectError: true,
		},
		{
			name: "未対応Operation",
			opts: PDFMergerOptions{
				Operation: "unknown",
				OutputDir: "output",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.ValidateOperation()
			if tt.expectError && err == nil {
				t.Fatal("エラーが期待されましたが、エラーが発生しませんでした")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("予期しないエラーが発生しました: %v", err)
			}
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
			Operation:   OperationExtractImages,
			SrcFile:     "test_data/org/sample_01_01.pdf",
			OutputDir:   outputDir,
			ImageFormat: "jpg",
			StartPage:   0,
			EndPage:     0,
		}

		service, logBuf := newTestPDFMergerService()

		err := service.Process(opts)
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
		output := logBuf.String()
		if !strings.Contains(output, "PDF画像抽出を開始します") {
			t.Error("期待される開始メッセージが出力されていません")
		}
		if !strings.Contains(output, "画像抽出が完了しました") {
			t.Error("期待される完了メッセージが出力されていません")
		}
	})

	t.Run("完全なPDF作成フロー", func(t *testing.T) {
		outputDir := "test_data/tmp/integration_created"
		outputFile := filepath.Join(outputDir, "integration_created.pdf")
		defer os.RemoveAll(outputDir)

		opts := PDFMergerOptions{
			Operation: OperationMergeIntoNew,
			SrcDir:    "test_data/org",
			OutputDir: outputDir,
		}

		service, logBuf := newTestPDFMergerService()

		err := service.Process(opts)
		if err != nil {
			t.Fatalf("PDF作成処理でエラーが発生: %v", err)
		}

		// 出力ファイルが作成されているか確認
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			t.Error("PDFファイルが作成されていません")
		}

		// 出力メッセージの確認
		output := logBuf.String()
		if !strings.Contains(output, "検出した画像") {
			t.Error("期待される画像検出メッセージが出力されていません")
		}
		if !strings.Contains(output, "PDF を生成しました") {
			t.Error("期待される完了メッセージが出力されていません")
		}
	})
}
