package dedupimages

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/landmaster135/devbox/internal/movie_extractor/usecases/common"
)

var supportedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
}

type Input struct {
	SrcDir    string
	MatchRate float64
	Log       bool
	LogWriter io.Writer
	OutDir    string
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Handle(input Input) (string, error) {
	if strings.TrimSpace(input.SrcDir) == "" {
		return "", fmt.Errorf("src-dir は必須です")
	}
	if input.MatchRate < 0 || input.MatchRate > 100 {
		return "", fmt.Errorf("match-rate は0から100の範囲で指定してください")
	}
	if strings.TrimSpace(input.OutDir) == "" {
		return "", fmt.Errorf("out-dir は必須です")
	}

	srcInfo, err := os.Stat(input.SrcDir)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("入力ディレクトリが存在しません: %s", input.SrcDir)
		}
		return "", fmt.Errorf("入力ディレクトリの確認に失敗しました: %w", err)
	}
	if !srcInfo.IsDir() {
		return "", fmt.Errorf("src-dir にファイルは指定できません: %s", input.SrcDir)
	}

	imagePaths, err := collectImageFiles(input.SrcDir)
	if err != nil {
		return "", err
	}
	if len(imagePaths) == 0 {
		return "", fmt.Errorf("src-dir に対象画像が存在しません: %s", input.SrcDir)
	}

	absOutDir, err := common.PrepareOutputDir(input.OutDir)
	if err != nil {
		return "", err
	}

	selectedImages := make([]string, 0, len(imagePaths))

	for _, candidate := range imagePaths {
		if _, decodeErr := decodeImage(candidate); decodeErr != nil {
			return "", fmt.Errorf("画像デコードに失敗しました (%s): %w", candidate, decodeErr)
		}

		isDuplicate := false
		if len(selectedImages) > 0 {
			recentSelected := selectedImages[len(selectedImages)-1]
			rate, matchErr := calculatePixelMatchRate(candidate, recentSelected)
			if matchErr != nil {
				return "", fmt.Errorf("画像比較に失敗しました (%s, %s): %w", candidate, recentSelected, matchErr)
			}
			if input.Log && input.LogWriter != nil {
				fmt.Fprintf(
					input.LogWriter,
					`照合率: "%s" vs "%s": %.2f%%`+"\n",
					filepath.Base(candidate),
					filepath.Base(recentSelected),
					rate*100,
				)
			}
			if rate*100 >= input.MatchRate {
				isDuplicate = true
			}
		}

		if isDuplicate {
			continue
		}

		destPath := filepath.Join(absOutDir, filepath.Base(candidate))
		if copyErr := copyFile(candidate, destPath); copyErr != nil {
			return "", fmt.Errorf("画像コピーに失敗しました (%s -> %s): %w", candidate, destPath, copyErr)
		}
		selectedImages = append(selectedImages, candidate)
	}

	var result strings.Builder
	result.WriteString("重複除外が完了しました。\n")
	fmt.Fprintf(&result, "入力画像数: %d\n", len(imagePaths))
	fmt.Fprintf(&result, "出力画像数: %d\n", len(selectedImages))
	fmt.Fprintf(&result, "出力ディレクトリ: %s\n", absOutDir)
	if len(selectedImages) > 0 {
		result.WriteString("出力ファイル:\n")
		for _, selectedImage := range selectedImages {
			fmt.Fprintf(&result, "- %s\n", filepath.Base(selectedImage))
		}
	}

	return result.String(), nil
}

func collectImageFiles(srcDir string) ([]string, error) {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, fmt.Errorf("入力ディレクトリの読み取りに失敗しました: %w", err)
	}

	paths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if _, ok := supportedImageExtensions[ext]; !ok {
			continue
		}
		paths = append(paths, filepath.Join(srcDir, entry.Name()))
	}
	sort.Strings(paths)
	return paths, nil
}

func copyFile(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}
	return destFile.Sync()
}
