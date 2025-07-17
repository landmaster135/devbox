package usecases

import (
	"fmt"
	"image"
	"image/draw"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/anthonynsimon/bild/blur"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gen2brain/webp"
)

// FilterMode は適用するフィルターの種類を表します
type FilterMode string

const (
	// BlurMode はぼかしフィルターを表します
	BlurMode FilterMode = "blur"
	// GrayscaleMode はグレースケールフィルターを表します
	GrayscaleMode FilterMode = "grayscale"
)

// ProcessingConfig は画像処理の設定を表します
type ProcessingConfig struct {
	SrcDir         string
	OutDir         string
	ArcDir         string
	X1, Y1, X2, Y2 int
	Suffix         string
	Move           bool
	Recursive      bool
	Workers        int
}

// FilterConfig はフィルター設定を表します
type FilterConfig struct {
	Mode    FilterMode
	Radius  float64
	RWeight float64
	GWeight float64
	BWeight float64
}

// ProcessingResult は処理結果を表します
type ProcessingResult struct {
	SuccessCount int
	ErrorCount   int
}

// ImageFilterService は画像フィルタリングサービスを表します
type ImageFilterService struct {
	ProcessingConfig *ProcessingConfig
	FilterConfig     *FilterConfig
}

// NewImageFilterService は新しいImageFilterServiceを作成します
func NewImageFilterService(procConfig *ProcessingConfig, filterConfig *FilterConfig) (*ImageFilterService, error) {
	service := &ImageFilterService{
		ProcessingConfig: procConfig,
		FilterConfig:     filterConfig,
	}

	// バリデーション実行
	if err := service.validate(); err != nil {
		return nil, err
	}

	return service, nil
}

// validate はサービスの設定をバリデーションします
func (s *ImageFilterService) validate() error {
	if err := s.validateProcessingConfig(); err != nil {
		return fmt.Errorf("processing config validation failed: %w", err)
	}

	if err := s.validateFilterConfig(); err != nil {
		return fmt.Errorf("filter config validation failed: %w", err)
	}

	return nil
}

// validateProcessingConfig は処理設定をバリデーションします
func (s *ImageFilterService) validateProcessingConfig() error {
	config := s.ProcessingConfig

	// 座標バリデーション
	if !(config.X1 == 0 && config.Y1 == 0 && config.X2 == 0 && config.Y2 == 0) {
		if config.X2 <= config.X1 || config.Y2 <= config.Y1 {
			return fmt.Errorf("invalid coordinates: x2 > x1, y2 > y1 required")
		}
	}

	// ワーカー数バリデーション
	if config.Workers <= 0 {
		return fmt.Errorf("workers must be positive: %d", config.Workers)
	}

	return nil
}

// validateFilterConfig はフィルター設定をバリデーションします
func (s *ImageFilterService) validateFilterConfig() error {
	config := s.FilterConfig

	// フィルターモードバリデーション
	if config.Mode != BlurMode && config.Mode != GrayscaleMode {
		return fmt.Errorf("unsupported filter mode: %s", config.Mode)
	}

	// RGB重みバリデーション
	if config.RWeight < 0.0 || config.RWeight > 1.0 {
		return fmt.Errorf("r-weight must be 0.0-1.0: %.2f", config.RWeight)
	}
	if config.GWeight < 0.0 || config.GWeight > 1.0 {
		return fmt.Errorf("g-weight must be 0.0-1.0: %.2f", config.GWeight)
	}
	if config.BWeight < 0.0 || config.BWeight > 1.0 {
		return fmt.Errorf("b-weight must be 0.0-1.0: %.2f", config.BWeight)
	}

	return nil
}

// ProcessImages は画像処理のメイン処理を実行します
func (s *ImageFilterService) ProcessImages() (*ProcessingResult, error) {
	// ファイルパス収集
	paths, err := s.CollectImagePaths()
	if err != nil {
		return nil, err
	}

	// ワーカープールで並行処理
	return s.ProcessImagesWithWorkers(paths)
}

// CollectImagePaths は処理対象の画像ファイルパスを収集します
func (s *ImageFilterService) CollectImagePaths() ([]string, error) {
	var paths []string
	config := s.ProcessingConfig

	walkFunc := func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			paths = append(paths, path)
		}
		return nil
	}

	if config.Recursive {
		// 再帰的にディレクトリを走査
		err := filepath.WalkDir(config.SrcDir, walkFunc)
		if err != nil {
			return nil, fmt.Errorf("directory walk failed: %w", err)
		}
	} else {
		// 単一ディレクトリのみ処理
		entries, err := os.ReadDir(config.SrcDir)
		if err != nil {
			return nil, fmt.Errorf("directory read failed: %w", err)
		}

		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(e.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
				paths = append(paths, filepath.Join(config.SrcDir, e.Name()))
			}
		}
	}

	return paths, nil
}

// ProcessImagesWithWorkers はワーカープールを使用して画像を並行処理します
func (s *ImageFilterService) ProcessImagesWithWorkers(paths []string) (*ProcessingResult, error) {
	config := s.ProcessingConfig
	pathsChan := make(chan string, 512)

	// パスをチャンネルに送信
	go func() {
		defer close(pathsChan)
		for _, path := range paths {
			pathsChan <- path
		}
	}()

	// ワーカープールの作成と実行
	var wg sync.WaitGroup
	result := &ProcessingResult{}
	var countMutex sync.Mutex

	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range pathsChan {
				// 画像にフィルターを適用して保存
				if err := s.ApplyFilterAndSave(path); err != nil {
					fmt.Printf("警告: %v\n", err)
					countMutex.Lock()
					result.ErrorCount++
					countMutex.Unlock()
					continue
				}

				// 元ファイルを移動（オプションが有効な場合）
				if config.Move {
					if err := s.MoveOriginal(path); err != nil {
						fmt.Printf("警告: %v\n", err)
						countMutex.Lock()
						result.ErrorCount++
						countMutex.Unlock()
					}
				}

				countMutex.Lock()
				result.SuccessCount++
				countMutex.Unlock()
			}
		}()
	}

	// すべてのワーカーの完了を待機
	wg.Wait()

	return result, nil
}

// applyFilterAndSave は画像を読み込み、指定した領域にフィルターを適用して outDir に保存します
func applyFilterAndSave(inPath, outDir string, x1, y1, x2, y2 int, suffix string, mode FilterMode, radius float64, rWeight, gWeight, bWeight float64) error {
	// ログ出力用のフォーマット文字列
	logFormat := "処理情報: ファイル=%s, 範囲=(%d,%d)-(%d,%d), モード=%s, 半径=%.1f"
	fmt.Printf(logFormat+"\n", filepath.Base(inPath), x1, y1, x2, y2, mode, radius)

	// ── 読み込み ──
	img, err := imgio.Open(inPath) // 拡張子を気にせずデコード
	if err != nil {
		return fmt.Errorf("open %s: %w", inPath, err)
	}

	// ── 矩形チェック ──
	bounds := img.Bounds()
	// 画像の実際の境界を取得
	minX, minY := bounds.Min.X, bounds.Min.Y
	maxX, maxY := bounds.Max.X, bounds.Max.Y

	fmt.Printf("画像情報: サイズ=%dx%d, 境界=(%d,%d)-(%d,%d)\n",
		bounds.Dx(), bounds.Dy(), minX, minY, maxX, maxY)

	// 全て0の場合は画像全体を対象にする
	if x1 == 0 && y1 == 0 && x2 == 0 && y2 == 0 {
		x1, y1, x2, y2 = minX, minY, maxX, maxY
		fmt.Printf("画像全体を処理対象に設定: (%d,%d)-(%d,%d)\n", x1, y1, x2, y2)
	}

	// 座標が画像の範囲内かチェック
	if x1 < minX || y1 < minY || x2 > maxX || y2 > maxY || x2 <= x1 || y2 <= y1 {
		return fmt.Errorf("invalid rectangle (%d,%d)-(%d,%d) for %s with bounds (%d,%d)-(%d,%d)",
			x1, y1, x2, y2, inPath, minX, minY, maxX, maxY)
	}

	// ── 元画像のクローンを作成 ──
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	// ── 指定領域を切り出し ──
	// x1, y1, x2, y2 = multiplyAndRound(x1, 1.5), multiplyAndRound(y1, 1.5), multiplyAndRound(x2, 1.5), multiplyAndRound(y2, 1.5)
	cropRect := image.Rect(x1, y1, x2, y2)
	fmt.Printf("切り出し範囲: (%d,%d)-(%d,%d)\n", cropRect.Min.X, cropRect.Min.Y, cropRect.Max.X, cropRect.Max.Y)
	subImg := transform.Crop(img, cropRect)

	// ── フィルター適用 ──
	var filtered *image.RGBA
	switch mode {
	case BlurMode:
		filtered = blur.Gaussian(subImg, radius)
	case GrayscaleMode:
		filtered = effect.GrayscaleWithWeights(subImg, rWeight, gWeight, bWeight)
	default:
		return fmt.Errorf("unsupported filter mode: %s", mode)
	}

	// ── フィルター適用した領域を元の画像に合成 ──
	targetRect := image.Rect(x1, y1, x2, y2)
	fmt.Printf("合成先範囲: (%d,%d)-(%d,%d)\n", targetRect.Min.X, targetRect.Min.Y, targetRect.Max.X, targetRect.Max.Y)
	fmt.Printf("フィルター適用画像の境界: (%d,%d)-(%d,%d)\n",
		filtered.Bounds().Min.X, filtered.Bounds().Min.Y,
		filtered.Bounds().Max.X, filtered.Bounds().Max.Y)

	draw.Draw(dst, targetRect, filtered, filtered.Bounds().Min, draw.Src)

	// ── 保存パス準備 ──
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}
	base := strings.TrimSuffix(filepath.Base(inPath), filepath.Ext(inPath))
	outPath := filepath.Join(outDir,
		fmt.Sprintf("%s_%s%s", base, suffix, filepath.Ext(inPath)))

	// ── 形式に応じて保存 ──
	ext := strings.ToLower(filepath.Ext(outPath))
	fmt.Printf("保存形式: %s, 保存先: %s\n", ext, outPath)

	switch ext {
	case ".jpg", ".jpeg":
		err = imgio.Save(outPath, dst, imgio.JPEGEncoder(95))
	case ".png":
		err = imgio.Save(outPath, dst, imgio.PNGEncoder())
	case ".webp":
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()
		err = webp.Encode(f, dst, webp.Options{Quality: 90})
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("unsupported extension: %s", ext)
	}
	return err
}

// ApplyFilterAndSave は画像にフィルターを適用して保存します（サービスメソッド版）
func (s *ImageFilterService) ApplyFilterAndSave(inPath string) error {
	config := s.ProcessingConfig
	filterConfig := s.FilterConfig

	return applyFilterAndSave(
		inPath,
		config.OutDir,
		config.X1, config.Y1, config.X2, config.Y2,
		config.Suffix,
		filterConfig.Mode,
		filterConfig.Radius,
		filterConfig.RWeight, filterConfig.GWeight, filterConfig.BWeight,
	)
}

// moveOriginal は元画像を arcDir に移動します
func moveOriginal(src, arcDir string) error {
	if err := os.MkdirAll(arcDir, 0o755); err != nil {
		return err
	}
	dst := filepath.Join(arcDir, filepath.Base(src))
	return os.Rename(src, dst)
}

// MoveOriginal は元画像を移動します（サービスメソッド版）
func (s *ImageFilterService) MoveOriginal(src string) error {
	return moveOriginal(src, s.ProcessingConfig.ArcDir)
}
