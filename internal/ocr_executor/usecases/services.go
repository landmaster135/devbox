package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gosseract "github.com/otiai10/gosseract/v2"
)

// OCRExecutorService は画像ファイルのOCR処理を行うサービス
type OCRExecutorService struct {
	targetPath   string
	recursive    bool
	languages    string
	outputDir    string
	outputFormat string
}

// OCRResult は個別ファイルのOCR結果を保持する構造体
type OCRResult struct {
	FilePath string `json:"file_path"`
	Text     string `json:"text"`
	Error    string `json:"error,omitempty"`
}

// ExecuteResult はOCR実行結果を保持する構造体
type ExecuteResult struct {
	Results []OCRResult `json:"results"`
	Total   int         `json:"total"`
}

// NewOCRExecutorService は新しいOCRExecutorServiceを作成する
func NewOCRExecutorService(targetPath string, recursive bool, languages string, outputDir string, outputFormat string) *OCRExecutorService {
	return &OCRExecutorService{
		targetPath:   targetPath,
		recursive:    recursive,
		languages:    languages,
		outputDir:    outputDir,
		outputFormat: outputFormat,
	}
}

// ExecuteFromPath は指定されたパスから画像ファイルを検索してOCRを実行する
func (s *OCRExecutorService) ExecuteFromPath() (*ExecuteResult, error) {
	// パスの存在確認
	info, err := os.Stat(s.targetPath)
	if err != nil {
		return nil, fmt.Errorf("パスが存在しません: %v", err)
	}

	var imagePaths []string

	if info.IsDir() {
		// ディレクトリの場合、画像ファイルを検索
		imagePaths, err = s.findImageFiles(s.targetPath)
		if err != nil {
			return nil, fmt.Errorf("画像ファイルの検索に失敗しました: %v", err)
		}
	} else {
		// ファイルの場合、画像ファイルかどうかを確認
		if s.isImageFile(s.targetPath) {
			imagePaths = []string{s.targetPath}
		}
	}

	// 画像ファイルをOCR処理
	results := make([]OCRResult, 0, len(imagePaths))
	for _, imagePath := range imagePaths {
		result := s.performOCR(imagePath)
		results = append(results, result)
	}

	executeResult := &ExecuteResult{
		Results: results,
		Total:   len(results),
	}

	// ファイル出力（指定されている場合）
	if s.outputDir != "" {
		if err := s.saveToFile(executeResult); err != nil {
			return nil, fmt.Errorf("ファイル出力に失敗しました: %v", err)
		}
	}

	return executeResult, nil
}

// findImageFiles はディレクトリ内の画像ファイルを検索する
func (s *OCRExecutorService) findImageFiles(dirPath string) ([]string, error) {
	var imagePaths []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合
		if info.IsDir() {
			// 再帰検索が無効で、ルートディレクトリでない場合はスキップ
			if !s.recursive && path != dirPath {
				return filepath.SkipDir
			}
			return nil
		}

		// 画像ファイルの場合
		if s.isImageFile(path) {
			imagePaths = append(imagePaths, path)
		}

		return nil
	})

	return imagePaths, err
}

// isImageFile はファイルが画像ファイルかどうかを判定する
func (s *OCRExecutorService) isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff", ".tif"}

	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// performOCR は画像ファイルにOCR処理を実行する
func (s *OCRExecutorService) performOCR(filePath string) OCRResult {
	result := OCRResult{
		FilePath: filePath,
	}

	// Tesseractクライアントを作成
	client := gosseract.NewClient()
	defer client.Close()

	// 言語設定
	if err := client.SetLanguage(s.languages); err != nil {
		result.Error = fmt.Sprintf("言語設定エラー: %v", err)
		return result
	}

	// 画像ファイルを設定
	if err := client.SetImage(filePath); err != nil {
		result.Error = fmt.Sprintf("画像読み込みエラー: %v", err)
		return result
	}

	// OCR実行
	text, err := client.Text()
	if err != nil {
		result.Error = fmt.Sprintf("OCR実行エラー: %v", err)
		return result
	}

	result.Text = strings.TrimSpace(text)
	return result
}

// saveToFile は結果をファイルに保存する
func (s *OCRExecutorService) saveToFile(result *ExecuteResult) error {
	// 出力ディレクトリが存在することを確認
	if err := os.MkdirAll(s.outputDir, 0755); err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// ファイル名を決定
	var fileName string
	var content string
	var err error

	if s.outputFormat == "json" {
		// JSON形式
		fileName = "ocr_results.json"
		content, err = result.FormatAsJSON()
	} else {
		// テキスト形式
		fileName = "ocr_results.txt"
		content = result.FormatAsText()
	}

	if err != nil {
		return fmt.Errorf("フォーマット変換エラー: %v", err)
	}

	// ファイルパスを構築
	filePath := filepath.Join(s.outputDir, fileName)

	// ファイルに書き込み
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %v", err)
	}

	return nil
}

// FormatAsText は結果をテキスト形式で出力する
func (r *ExecuteResult) FormatAsText() string {
	if r.Total == 0 {
		return "画像ファイルが見つかりませんでした。"
	}

	var output strings.Builder
	output.WriteString("=== OCR Results ===\n")
	output.WriteString(fmt.Sprintf("Total Images: %d\n\n", r.Total))

	for i, result := range r.Results {
		output.WriteString(fmt.Sprintf("[%d] %s\n", i+1, result.FilePath))
		if result.Error != "" {
			output.WriteString(fmt.Sprintf("Error: %s\n", result.Error))
		} else {
			output.WriteString(fmt.Sprintf("Text: %s\n", result.Text))
		}
		output.WriteString("\n")
	}

	return output.String()
}

// FormatAsJSON は結果をJSON形式で出力する
func (r *ExecuteResult) FormatAsJSON() (string, error) {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %v", err)
	}
	return string(jsonData), nil
}
