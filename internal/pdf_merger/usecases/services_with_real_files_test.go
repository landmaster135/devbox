package usecases

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// #==============================================================#
// ##         Tests                                              ##
// #==============================================================#
// extractPageNumber のテスト
func TestImageExtractionService_extractPageNumber(t *testing.T) {
	service := NewImageExtractionService()

	tests := []struct {
		name     string
		leftPart string
		want     int
	}{
		{
			name:     "document_001 パターン",
			leftPart: "document_001",
			want:     1,
		},
		{
			name:     "document_page1 パターン",
			leftPart: "document_page1",
			want:     1,
		},
		{
			name:     "document_1 パターン",
			leftPart: "document_1",
			want:     1,
		},
		{
			name:     "document_010 パターン",
			leftPart: "document_010",
			want:     10,
		},
		{
			name:     "document_page123 パターン",
			leftPart: "document_page123",
			want:     123,
		},
		{
			name:     "末尾アンダースコア付き",
			leftPart: "document_001_",
			want:     1,
		},
		{
			name:     "無効なパターン（数値なし）",
			leftPart: "document_abc",
			want:     1, // デフォルト値
		},
		{
			name:     "アンダースコアなし",
			leftPart: "document",
			want:     1, // デフォルト値
		},
		{
			name:     "空文字列",
			leftPart: "",
			want:     1, // デフォルト値
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.extractPageNumber(tt.leftPart)
			if got != tt.want {
				t.Errorf("extractPageNumber(%q) = %v, want %v", tt.leftPart, got, tt.want)
			}
		})
	}
}

// parseImageFileName のテスト
func TestImageExtractionService_parseImageFileName(t *testing.T) {
	service := NewImageExtractionService()
	outputDir := "/test/output"

	tests := []struct {
		name     string
		filename string
		want     fileInfo
		wantErr  bool
	}{
		{
			name:     "正常なファイル名 - document_001_Im0.jpg",
			filename: "document_001_Im0.jpg",
			want: fileInfo{
				originalPath: filepath.Join(outputDir, "document_001_Im0.jpg"),
				pageNum:      1,
				imageNum:     0,
				ext:          ".jpg",
			},
			wantErr: false,
		},
		{
			name:     "正常なファイル名 - document_page1_Im2.png",
			filename: "document_page1_Im2.png",
			want: fileInfo{
				originalPath: filepath.Join(outputDir, "document_page1_Im2.png"),
				pageNum:      1,
				imageNum:     2,
				ext:          ".png",
			},
			wantErr: false,
		},
		{
			name:     "正常なファイル名 - document_010_Im5.jpg",
			filename: "document_010_Im5.jpg",
			want: fileInfo{
				originalPath: filepath.Join(outputDir, "document_010_Im5.jpg"),
				pageNum:      10,
				imageNum:     5,
				ext:          ".jpg",
			},
			wantErr: false,
		},
		{
			name:     "_Im が含まれていない",
			filename: "document_001.jpg",
			want:     fileInfo{},
			wantErr:  true,
		},
		{
			name:     "無効な形式 - _Im で分割できない",
			filename: "document_Im.jpg",
			want: fileInfo{
				originalPath: filepath.Join(outputDir, "document_Im.jpg"),
				pageNum:      1, // デフォルト値
				imageNum:     0, // デフォルト値
				ext:          ".jpg",
			},
			wantErr: false,
		},
		{
			name:     "画像番号が数値でない",
			filename: "document_001_Imabc.jpg",
			want: fileInfo{
				originalPath: filepath.Join(outputDir, "document_001_Imabc.jpg"),
				pageNum:      1,
				imageNum:     0, // デフォルト値
				ext:          ".jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.parseImageFileName(tt.filename, outputDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseImageFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseImageFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

// isSupportedFormat のテスト
func TestImageExtractionService_isSupportedFormat(t *testing.T) {
	service := NewImageExtractionService()

	tests := []struct {
		name       string
		format     string
		wantResult bool
		wantMsg    string
		wantErr    bool
	}{
		{
			name:       "サポート形式 - jpg",
			format:     "jpg",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "サポート形式 - jpeg",
			format:     "jpeg",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "サポート形式 - png",
			format:     "png",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "サポート形式 - tiff",
			format:     "tiff",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "サポート形式 - webp",
			format:     "webp",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "大文字小文字混在 - JPG",
			format:     "JPG",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "大文字小文字混在 - PNG",
			format:     "PNG",
			wantResult: true,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    false,
		},
		{
			name:       "非サポート形式 - gif",
			format:     "gif",
			wantResult: false,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    true,
		},
		{
			name:       "非サポート形式 - bmp",
			format:     "bmp",
			wantResult: false,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    true,
		},
		{
			name:       "非サポート形式 - svg",
			format:     "svg",
			wantResult: false,
			wantMsg:    "(サポート形式: jpg, jpeg, png, tiff, webp)",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotMsg, err := service.isSupportedFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("isSupportedFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("isSupportedFormat() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("isSupportedFormat() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}

// ValidatePageRange のテスト
func TestImageExtractionService_ValidatePageRange(t *testing.T) {
	service := NewImageExtractionService()

	tests := []struct {
		name       string
		startPage  int
		endPage    int
		totalPages int
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "正常範囲 - 両方指定",
			startPage:  3,
			endPage:    7,
			totalPages: 10,
			wantErr:    false,
		},
		{
			name:       "正常範囲 - 開始ページのみ",
			startPage:  5,
			endPage:    0,
			totalPages: 10,
			wantErr:    false,
		},
		{
			name:       "正常範囲 - 終了ページのみ",
			startPage:  0,
			endPage:    8,
			totalPages: 10,
			wantErr:    false,
		},
		{
			name:       "正常範囲 - 全ページ",
			startPage:  0,
			endPage:    0,
			totalPages: 10,
			wantErr:    false,
		},
		{
			name:       "正常範囲 - 単一ページ",
			startPage:  5,
			endPage:    5,
			totalPages: 10,
			wantErr:    false,
		},
		{
			name:       "負の開始ページ",
			startPage:  -1,
			endPage:    5,
			totalPages: 10,
			wantErr:    true,
			errMsg:     "開始ページは0以上の値を指定してください",
		},
		{
			name:       "負の終了ページ",
			startPage:  3,
			endPage:    -1,
			totalPages: 10,
			wantErr:    true,
			errMsg:     "終了ページは0以上の値を指定してください",
		},
		{
			name:       "開始ページが総ページ数を超過",
			startPage:  15,
			endPage:    0,
			totalPages: 10,
			wantErr:    true,
			errMsg:     "開始ページがPDFの総ページ数を超えています",
		},
		{
			name:       "終了ページが総ページ数を超過",
			startPage:  0,
			endPage:    15,
			totalPages: 10,
			wantErr:    true,
			errMsg:     "終了ページがPDFの総ページ数を超えています",
		},
		{
			name:       "開始ページが終了ページより大きい",
			startPage:  8,
			endPage:    3,
			totalPages: 10,
			wantErr:    true,
			errMsg:     "開始ページが終了ページより後になっています",
		},
		{
			name:       "境界値 - 最初のページ",
			startPage:  1,
			endPage:    1,
			totalPages: 10,
			wantErr:    false,
		},
		{
			name:       "境界値 - 最後のページ",
			startPage:  10,
			endPage:    10,
			totalPages: 10,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidatePageRange(tt.startPage, tt.endPage, tt.totalPages)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePageRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !reflect.DeepEqual(err.Error()[:len(tt.errMsg)], tt.errMsg) {
					t.Errorf("ValidatePageRange() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// sortImageFiles のテスト
func TestImageExtractionService_sortImageFiles(t *testing.T) {
	service := NewImageExtractionService()

	tests := []struct {
		name  string
		input []fileInfo
		want  []fileInfo
	}{
		{
			name: "ページ番号順でソート",
			input: []fileInfo{
				{originalPath: "doc_003_Im0.jpg", pageNum: 3, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_002_Im0.jpg", pageNum: 2, imageNum: 0, ext: ".jpg"},
			},
			want: []fileInfo{
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_002_Im0.jpg", pageNum: 2, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_003_Im0.jpg", pageNum: 3, imageNum: 0, ext: ".jpg"},
			},
		},
		{
			name: "同じページ番号内で画像番号順でソート",
			input: []fileInfo{
				{originalPath: "doc_001_Im2.jpg", pageNum: 1, imageNum: 2, ext: ".jpg"},
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_001_Im1.jpg", pageNum: 1, imageNum: 1, ext: ".jpg"},
			},
			want: []fileInfo{
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_001_Im1.jpg", pageNum: 1, imageNum: 1, ext: ".jpg"},
				{originalPath: "doc_001_Im2.jpg", pageNum: 1, imageNum: 2, ext: ".jpg"},
			},
		},
		{
			name: "ページ番号と画像番号の混在",
			input: []fileInfo{
				{originalPath: "doc_002_Im1.jpg", pageNum: 2, imageNum: 1, ext: ".jpg"},
				{originalPath: "doc_001_Im2.jpg", pageNum: 1, imageNum: 2, ext: ".jpg"},
				{originalPath: "doc_002_Im0.jpg", pageNum: 2, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
			},
			want: []fileInfo{
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_001_Im2.jpg", pageNum: 1, imageNum: 2, ext: ".jpg"},
				{originalPath: "doc_002_Im0.jpg", pageNum: 2, imageNum: 0, ext: ".jpg"},
				{originalPath: "doc_002_Im1.jpg", pageNum: 2, imageNum: 1, ext: ".jpg"},
			},
		},
		{
			name:  "空のスライス",
			input: []fileInfo{},
			want:  []fileInfo{},
		},
		{
			name: "単一要素",
			input: []fileInfo{
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
			},
			want: []fileInfo{
				{originalPath: "doc_001_Im0.jpg", pageNum: 1, imageNum: 0, ext: ".jpg"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.sortImageFiles(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortImageFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

// calculateActualPageNumber のテスト
func TestImageExtractionService_calculateActualPageNumber(t *testing.T) {
	service := NewImageExtractionService()

	tests := []struct {
		name            string
		fileInfo        fileInfo
		index           int
		startPageOffset int
		want            int
	}{
		{
			name: "startPageOffset = 0, ページ番号あり",
			fileInfo: fileInfo{
				pageNum: 5,
			},
			index:           2,
			startPageOffset: 0,
			want:            5,
		},
		{
			name: "startPageOffset = 0, ページ番号なし",
			fileInfo: fileInfo{
				pageNum: 0,
			},
			index:           2,
			startPageOffset: 0,
			want:            3, // index + 1
		},
		{
			name: "startPageOffset > 0",
			fileInfo: fileInfo{
				pageNum: 5, // この値は無視される
			},
			index:           2,
			startPageOffset: 10,
			want:            12, // startPageOffset + index
		},
		{
			name: "startPageOffset = 1, index = 0",
			fileInfo: fileInfo{
				pageNum: 3,
			},
			index:           0,
			startPageOffset: 1,
			want:            1,
		},
		{
			name: "startPageOffset = 0, index = 0, ページ番号あり",
			fileInfo: fileInfo{
				pageNum: 1,
			},
			index:           0,
			startPageOffset: 0,
			want:            1,
		},
		{
			name: "startPageOffset = 0, index = 0, ページ番号なし",
			fileInfo: fileInfo{
				pageNum: 0,
			},
			index:           0,
			startPageOffset: 0,
			want:            1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.calculateActualPageNumber(tt.fileInfo, tt.index, tt.startPageOffset)
			if got != tt.want {
				t.Errorf("calculateActualPageNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

// createTemporaryPDF のテスト
func TestPDFCreationService_createTemporaryPDF(t *testing.T) {
	service := NewPDFCreationService()

	t.Run("出力ディレクトリ配下に一時ファイルを生成", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "pdf-temp-*")
		if err != nil {
			t.Fatalf("一時ディレクトリの作成に失敗: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		output := filepath.Join(tmpDir, "result.pdf")
		path, err := service.createTemporaryPDF(output)
		if err != nil {
			t.Fatalf("createTemporaryPDF() でエラーが発生: %v", err)
		}
		defer os.Remove(path)

		if filepath.Dir(path) != tmpDir {
			t.Errorf("一時ファイルのディレクトリが一致しません: got %s, want %s", filepath.Dir(path), tmpDir)
		}

		if _, err := os.Stat(path); err == nil {
			t.Errorf("一時PDFが事前に存在しています: %s", path)
		} else if !os.IsNotExist(err) {
			t.Fatalf("一時PDFの存在確認で想定外のエラー: %v", err)
		}
	})

	t.Run("連続呼び出しで異なるパスを返す", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "pdf-temp-*")
		if err != nil {
			t.Fatalf("一時ディレクトリの作成に失敗: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		output := filepath.Join(tmpDir, "result.pdf")
		path1, err := service.createTemporaryPDF(output)
		if err != nil {
			t.Fatalf("1回目の createTemporaryPDF() でエラーが発生: %v", err)
		}
		defer os.Remove(path1)

		path2, err := service.createTemporaryPDF(output)
		if err != nil {
			t.Fatalf("2回目の createTemporaryPDF() でエラーが発生: %v", err)
		}
		defer os.Remove(path2)

		if path1 == path2 {
			t.Errorf("一時PDFのパスが重複しました: %s", path1)
		}
	})
}
