package usecases

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
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

const (
	grayscaleDiffTolerance  = 20
	cannyHighThresholdRatio = 0.20
	cannyLowThresholdRatio  = 0.10
)

// DedupImagesInput は dedup-images 操作の入力です。
type DedupImagesInput struct {
	SrcDir    string
	MatchRate float64
	Log       bool
	LogWriter io.Writer
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
	result.WriteString(fmt.Sprintf("入力画像数: %d\n", len(imagePaths)))
	result.WriteString(fmt.Sprintf("出力画像数: %d\n", len(selectedImages)))
	result.WriteString(fmt.Sprintf("出力ディレクトリ: %s\n", absOutDir))
	if len(selectedImages) > 0 {
		result.WriteString("出力ファイル:\n")
		for _, selectedImage := range selectedImages {
			result.WriteString(fmt.Sprintf("- %s\n", filepath.Base(selectedImage)))
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

	edges1 := detectCannyEdges(img1)
	edges2 := detectCannyEdges(img2)
	edgeRate, ok := calculateEdgeIoURate(edges1, edges2)
	if ok {
		return edgeRate, nil
	}

	// エッジが取得できないベタ画像同士は、既存の階調差比較へフォールバックする。
	return calculateGrayscaleMatchRate(img1, img2), nil
}

func calculateGrayscaleMatchRate(img1 image.Image, img2 image.Image) float64 {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	matchCount := 0
	totalPixelCount := bounds1.Dx() * bounds1.Dy()
	for y := 0; y < bounds1.Dy(); y++ {
		for x := 0; x < bounds1.Dx(); x++ {
			pixel1 := color.GrayModel.Convert(img1.At(bounds1.Min.X+x, bounds1.Min.Y+y)).(color.Gray)
			pixel2 := color.GrayModel.Convert(img2.At(bounds2.Min.X+x, bounds2.Min.Y+y)).(color.Gray)
			diff := int(pixel1.Y) - int(pixel2.Y)
			if diff < 0 {
				diff = -diff
			}
			if diff <= grayscaleDiffTolerance {
				matchCount++
			}
		}
	}

	return float64(matchCount) / float64(totalPixelCount)
}

func detectCannyEdges(img image.Image) [][]bool {
	gray := toGrayMatrix(img)
	blurred := applyGaussianBlur(gray)
	magnitude, direction := applySobel(blurred)
	suppressed := applyNonMaximumSuppression(magnitude, direction)
	return applyHysteresisThreshold(suppressed)
}

func toGrayMatrix(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	matrix := make([][]float64, height)
	for y := 0; y < height; y++ {
		row := make([]float64, width)
		for x := 0; x < width; x++ {
			pixel := color.GrayModel.Convert(img.At(bounds.Min.X+x, bounds.Min.Y+y)).(color.Gray)
			row[x] = float64(pixel.Y)
		}
		matrix[y] = row
	}
	return matrix
}

func applyGaussianBlur(matrix [][]float64) [][]float64 {
	height := len(matrix)
	width := len(matrix[0])
	weights := []float64{1, 4, 6, 4, 1}
	normalizer := 16.0

	horizontal := make([][]float64, height)
	for y := 0; y < height; y++ {
		row := make([]float64, width)
		for x := 0; x < width; x++ {
			sum := 0.0
			for index := -2; index <= 2; index++ {
				sum += matrix[y][clampIndex(x+index, width)] * weights[index+2]
			}
			row[x] = sum / normalizer
		}
		horizontal[y] = row
	}

	vertical := make([][]float64, height)
	for y := 0; y < height; y++ {
		row := make([]float64, width)
		for x := 0; x < width; x++ {
			sum := 0.0
			for index := -2; index <= 2; index++ {
				sum += horizontal[clampIndex(y+index, height)][x] * weights[index+2]
			}
			row[x] = sum / normalizer
		}
		vertical[y] = row
	}

	return vertical
}

func applySobel(matrix [][]float64) ([][]float64, [][]float64) {
	height := len(matrix)
	width := len(matrix[0])
	magnitude := make([][]float64, height)
	direction := make([][]float64, height)

	for y := 0; y < height; y++ {
		magRow := make([]float64, width)
		dirRow := make([]float64, width)
		for x := 0; x < width; x++ {
			xm1 := clampIndex(x-1, width)
			xp1 := clampIndex(x+1, width)
			ym1 := clampIndex(y-1, height)
			yp1 := clampIndex(y+1, height)

			gx := -matrix[ym1][xm1] + matrix[ym1][xp1] +
				-2*matrix[y][xm1] + 2*matrix[y][xp1] +
				-matrix[yp1][xm1] + matrix[yp1][xp1]

			gy := matrix[ym1][xm1] + 2*matrix[ym1][x] + matrix[ym1][xp1] +
				-matrix[yp1][xm1] - 2*matrix[yp1][x] - matrix[yp1][xp1]

			magRow[x] = math.Hypot(gx, gy)
			dirRow[x] = math.Atan2(gy, gx)
		}
		magnitude[y] = magRow
		direction[y] = dirRow
	}

	return magnitude, direction
}

func applyNonMaximumSuppression(magnitude [][]float64, direction [][]float64) [][]float64 {
	height := len(magnitude)
	width := len(magnitude[0])
	result := make([][]float64, height)
	for y := range result {
		result[y] = make([]float64, width)
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			angle := direction[y][x] * 180.0 / math.Pi
			if angle < 0 {
				angle += 180
			}

			current := magnitude[y][x]
			var before float64
			var after float64

			switch {
			case angle < 22.5 || angle >= 157.5:
				before = magnitude[y][x-1]
				after = magnitude[y][x+1]
			case angle < 67.5:
				before = magnitude[y-1][x+1]
				after = magnitude[y+1][x-1]
			case angle < 112.5:
				before = magnitude[y-1][x]
				after = magnitude[y+1][x]
			default:
				before = magnitude[y-1][x-1]
				after = magnitude[y+1][x+1]
			}

			if current >= before && current >= after {
				result[y][x] = current
			}
		}
	}

	return result
}

func applyHysteresisThreshold(suppressed [][]float64) [][]bool {
	height := len(suppressed)
	width := len(suppressed[0])
	maxValue := 0.0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if suppressed[y][x] > maxValue {
				maxValue = suppressed[y][x]
			}
		}
	}
	if maxValue == 0 {
		return make([][]bool, height)
	}

	high := maxValue * cannyHighThresholdRatio
	low := maxValue * cannyLowThresholdRatio

	state := make([][]uint8, height)
	for y := 0; y < height; y++ {
		row := make([]uint8, width)
		for x := 0; x < width; x++ {
			value := suppressed[y][x]
			if value >= high {
				row[x] = 2
			} else if value >= low {
				row[x] = 1
			}
		}
		state[y] = row
	}

	type point struct {
		y int
		x int
	}
	queue := make([]point, 0, width*height/8)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if state[y][x] == 2 {
				queue = append(queue, point{y: y, x: x})
			}
		}
	}

	head := 0
	for head < len(queue) {
		current := queue[head]
		head++

		for offsetY := -1; offsetY <= 1; offsetY++ {
			for offsetX := -1; offsetX <= 1; offsetX++ {
				if offsetY == 0 && offsetX == 0 {
					continue
				}
				nextY := current.y + offsetY
				nextX := current.x + offsetX
				if nextY < 0 || nextY >= height || nextX < 0 || nextX >= width {
					continue
				}
				if state[nextY][nextX] == 1 {
					state[nextY][nextX] = 2
					queue = append(queue, point{y: nextY, x: nextX})
				}
			}
		}
	}

	edges := make([][]bool, height)
	for y := 0; y < height; y++ {
		row := make([]bool, width)
		for x := 0; x < width; x++ {
			row[x] = state[y][x] == 2
		}
		edges[y] = row
	}
	return edges
}

func calculateEdgeIoURate(edges1 [][]bool, edges2 [][]bool) (float64, bool) {
	if len(edges1) == 0 || len(edges2) == 0 {
		return 0, false
	}
	if len(edges1) != len(edges2) || len(edges1[0]) != len(edges2[0]) {
		return 0, false
	}

	intersection := 0
	union := 0
	for y := range edges1 {
		for x := range edges1[y] {
			a := edges1[y][x]
			b := edges2[y][x]
			if a || b {
				union++
			}
			if a && b {
				intersection++
			}
		}
	}
	if union == 0 {
		return 0, false
	}
	return float64(intersection) / float64(union), true
}

func clampIndex(value int, length int) int {
	if value < 0 {
		return 0
	}
	if value >= length {
		return length - 1
	}
	return value
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
