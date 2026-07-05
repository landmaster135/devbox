package dedupimages

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

const (
	grayscaleDiffTolerance  = 20
	cannyHighThresholdRatio = 0.20
	cannyLowThresholdRatio  = 0.10
)

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
	height := bounds1.Dy()
	width := bounds1.Dx()
	matchCount := 0
	totalPixelCount := width * height
	for y := range height {
		for x := range width {
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
	for y := range height {
		row := make([]float64, width)
		for x := range width {
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
	for y := range height {
		row := make([]float64, width)
		for x := range width {
			sum := 0.0
			for index := -2; index <= 2; index++ {
				sum += matrix[y][clampIndex(x+index, width)] * weights[index+2]
			}
			row[x] = sum / normalizer
		}
		horizontal[y] = row
	}

	vertical := make([][]float64, height)
	for y := range height {
		row := make([]float64, width)
		for x := range width {
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

	for y := range height {
		magRow := make([]float64, width)
		dirRow := make([]float64, width)
		for x := range width {
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
	for y := range height {
		for x := range width {
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
	for y := range height {
		row := make([]uint8, width)
		for x := range width {
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
	for y := range height {
		for x := range width {
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
	for y := range height {
		row := make([]bool, width)
		for x := range width {
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
