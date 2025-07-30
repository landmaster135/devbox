package usecases

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Base64ExtractorService は画像ファイルのbase64変換を行うサービス
type Base64ExtractorService struct {
	targetPath string
	recursive  bool
}

// ImageResult は画像変換結果を保持する構造体
type ImageResult struct {
	FilePath string `json:"file_path"`
	Base64   string `json:"base64"`
	Error    string `json:"error,omitempty"`
}

// ExtractResult は抽出結果を保持する構造体
type ExtractResult struct {
	Images []ImageResult `json:"images"`
	Total  int           `json:"total"`
}

// NewBase64ExtractorService は新しいBase64ExtractorServiceを作成する
func NewBase64ExtractorService(targetPath string, recursive bool) *Base64ExtractorService {
	return &Base64ExtractorService{
		targetPath: targetPath,
		recursive:  recursive,
	}
}

// ExtractFromPath は指定されたパスから画像ファイルを検索してbase64に変換する
func (s *Base64ExtractorService) ExtractFromPath() (*ExtractResult, error) {
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

	// 画像ファイルをbase64に変換
	results := make([]ImageResult, 0, len(imagePaths))
	for _, imagePath := range imagePaths {
		result := s.convertToBase64(imagePath)
		results = append(results, result)
	}

	return &ExtractResult{
		Images: results,
		Total:  len(results),
	}, nil
}

// findImageFiles はディレクトリ内の画像ファイルを検索する
func (s *Base64ExtractorService) findImageFiles(dirPath string) ([]string, error) {
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
func (s *Base64ExtractorService) isImageFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp"}

	for _, imageExt := range imageExtensions {
		if ext == imageExt {
			return true
		}
	}
	return false
}

// convertToBase64 は画像ファイルをbase64に変換する
func (s *Base64ExtractorService) convertToBase64(filePath string) ImageResult {
	result := ImageResult{
		FilePath: filePath,
	}

	// ファイルを読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		result.Error = fmt.Sprintf("ファイル読み込みエラー: %v", err)
		return result
	}

	// base64エンコード
	result.Base64 = base64.StdEncoding.EncodeToString(data)
	return result
}

// FormatAsText は結果をテキスト形式で出力する
func (r *ExtractResult) FormatAsText() string {
	if r.Total == 0 {
		return "画像ファイルが見つかりませんでした。"
	}

	var output strings.Builder
	output.WriteString(fmt.Sprintf("=== Base64 Extraction Results ===\n"))
	output.WriteString(fmt.Sprintf("Total Images: %d\n\n", r.Total))

	for i, image := range r.Images {
		output.WriteString(fmt.Sprintf("[%d] %s\n", i+1, image.FilePath))
		if image.Error != "" {
			output.WriteString(fmt.Sprintf("Error: %s\n", image.Error))
		} else {
			output.WriteString(fmt.Sprintf("Base64: %s\n", image.Base64))
		}
		output.WriteString("\n")
	}

	return output.String()
}

// FormatAsJSON は結果をJSON形式で出力する
func (r *ExtractResult) FormatAsJSON() (string, error) {
	jsonData, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON変換エラー: %v", err)
	}
	return string(jsonData), nil
}
