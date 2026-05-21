package usecases

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var supportedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
}

// DedupImagesInput は dedup-images 操作の入力です。
type DedupImagesInput struct {
	SrcDir    string
	MatchRate float64
	OutDir    string
}

// HandleDedupImages は画像ディレクトリから重複画像を除去し、代表画像のみを出力します。
func (s *Service) HandleDedupImages(input DedupImagesInput) (string, error) {
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

	if err := os.MkdirAll(input.OutDir, 0755); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	absOutDir, err := filepath.Abs(input.OutDir)
	if err != nil {
		return "", fmt.Errorf("出力ディレクトリの絶対パス変換に失敗しました: %w", err)
	}

	selectedImages := make([]string, 0, len(imagePaths))
	writtenNames := make([]string, 0, len(imagePaths))

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
		writtenNames = append(writtenNames, filepath.Base(candidate))
	}

	var result strings.Builder
	result.WriteString("重複除外が完了しました。\n")
	result.WriteString(fmt.Sprintf("入力画像数: %d\n", len(imagePaths)))
	result.WriteString(fmt.Sprintf("出力画像数: %d\n", len(writtenNames)))
	result.WriteString(fmt.Sprintf("出力ディレクトリ: %s\n", absOutDir))
	if len(writtenNames) > 0 {
		result.WriteString("出力ファイル:\n")
		for _, name := range writtenNames {
			result.WriteString(fmt.Sprintf("- %s\n", name))
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

func calculatePixelMatchRate(filePath1 string, filePath2 string) (float64, error) {
	img1, err := decodeImage(filePath1)
	if err != nil {
		return 0, err
	}
	img2, err := decodeImage(filePath2)
	if err != nil {
		return 0, err
	}

	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return 0, nil
	}
	if bounds1.Dx() == 0 || bounds1.Dy() == 0 {
		return 0, nil
	}

	matchCount := 0
	totalPixelCount := bounds1.Dx() * bounds1.Dy()
	for y := 0; y < bounds1.Dy(); y++ {
		for x := 0; x < bounds1.Dx(); x++ {
			pixel1 := color.GrayModel.Convert(img1.At(bounds1.Min.X+x, bounds1.Min.Y+y)).(color.Gray)
			pixel2 := color.GrayModel.Convert(img2.At(bounds2.Min.X+x, bounds2.Min.Y+y)).(color.Gray)
			if pixel1.Y == pixel2.Y {
				matchCount++
			}
		}
	}

	return float64(matchCount) / float64(totalPixelCount), nil
}

func decodeImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルを開けません: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("画像デコードに失敗しました: %w", err)
	}
	return img, nil
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
