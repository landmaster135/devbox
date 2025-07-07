package usecases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// #==============================================================#
// ##          Consts and Types for single process               ##
// #==============================================================#
// ConversionConfig holds the configuration for video conversion
type ConversionConfig struct {
	InputFile   string
	OutputFile  string
	FPS         int
	Width       int     // GIF変換時のみ（0=デフォルト品質）
	Speed       float64 // GIF変換時のみ
	Loop        int     // GIF変換時のみ（0=無限ループ、-1=ループなし）
	UseItsScale bool    // GIF変換時のみ
}

// ConversionResult holds the result of a single file conversion
type ConversionResult struct {
	InputFile  string
	OutputFile string
	Success    bool
	Error      error
}

// MovieConverterService handles video conversion operations
type MovieConverterService struct {
	config ConversionConfig
}

// NewMovieConverterService creates a new MovieConverterService instance
func NewMovieConverterService(config ConversionConfig) *MovieConverterService {
	return &MovieConverterService{config: config}
}

// ConvertMP4ToGIF converts MP4 to GIF using ffmpeg-go
func (s *MovieConverterService) ConvertMP4ToGIF() error {
	log.Printf("MP4からGIFに変換中: %s -> %s", s.config.InputFile, s.config.OutputFile)

	// デフォルト値の設定
	s.setMP4ToGIFDefaults()

	// ファイル名にスペースが含まれている場合の警告
	if strings.Contains(s.config.InputFile, " ") {
		log.Println("警告: ソースファイル名にスペース文字が含まれています。問題が発生する可能性があります。")
	}

	// PowerShellスクリプトと同じフィルター設定
	var vfStr string
	if s.config.Width != 0 {
		// 幅が指定されている場合（PowerShellスクリプトと同じ）
		vfStr = fmt.Sprintf("scale=%d:-1", s.config.Width)
	} else {
		// デフォルト品質（パレット生成とパレット使用）
		vfStr = "split[s0][s1];[s0]palettegen[p];[s1][p]paletteuse"
	}

	// 速度調整の設定
	if s.config.UseItsScale {
		// itsscaleを使用する場合
		itsScale := 1.0 / s.config.Speed
		log.Printf("itsscale値: %f", itsScale)

		// itsscaleをInputレベルで設定（PowerShellスクリプトに近い動作）
		err := ffmpeg.Input(s.config.InputFile, ffmpeg.KwArgs{"itsscale": itsScale}).
			Output(s.config.OutputFile, ffmpeg.KwArgs{
				"lavfi": vfStr,
				"r":     s.config.FPS,
				"loop":  s.config.Loop,
			}).
			OverWriteOutput().
			Run()

		if err != nil {
			return fmt.Errorf("ffmpeg実行エラー: %w", err)
		}
	} else {
		// setptsを使用する場合
		setPts := fmt.Sprintf("PTS/%f", s.config.Speed)
		log.Printf("setpts値: %s", setPts)
		vfStr = fmt.Sprintf("%s,setpts=%s", vfStr, setPts)

		err := ffmpeg.Input(s.config.InputFile).
			Output(s.config.OutputFile, ffmpeg.KwArgs{
				"lavfi": vfStr,
				"r":     s.config.FPS,
				"loop":  s.config.Loop,
			}).
			OverWriteOutput().
			Run()

		if err != nil {
			return fmt.Errorf("ffmpeg実行エラー: %w", err)
		}
	}

	log.Printf("MP4からGIFへの変換が完了しました: %s", s.config.OutputFile)
	return nil
}

// ConvertGIFToMP4 converts GIF to MP4 using ffmpeg-go
func (s *MovieConverterService) ConvertGIFToMP4() error {
	log.Printf("GIFからMP4に変換中: %s -> %s", s.config.InputFile, s.config.OutputFile)

	// デフォルト値の設定
	s.setGIFToMP4Defaults()

	// ファイル名にスペースが含まれている場合の警告
	if strings.Contains(s.config.InputFile, " ") {
		log.Println("警告: ソースファイル名にスペース文字が含まれています。問題が発生する可能性があります。")
	}

	err := ffmpeg.Input(s.config.InputFile).
		Output(s.config.OutputFile, ffmpeg.KwArgs{
			"r":        s.config.FPS,
			"movflags": "faststart",
			"c:v":      "libx264",
			"pix_fmt":  "yuv420p",
		}).
		OverWriteOutput().
		Run()

	if err != nil {
		return fmt.Errorf("ffmpeg実行エラー: %w", err)
	}

	log.Printf("GIFからMP4への変換が完了しました: %s", s.config.OutputFile)
	return nil
}

// Convert performs conversion based on file extensions
func (s *MovieConverterService) Convert() error {
	// 入力ファイルの存在確認
	if _, err := os.Stat(s.config.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("入力ファイルが見つかりません: %s", s.config.InputFile)
	}

	// 拡張子の取得
	inputExt := strings.ToLower(filepath.Ext(s.config.InputFile))
	outputExt := strings.ToLower(filepath.Ext(s.config.OutputFile))

	// 拡張子の検証
	if !strings.Contains(s.config.InputFile, ".") {
		return fmt.Errorf("入力ファイル名に拡張子が含まれていません: %s", s.config.InputFile)
	}

	// 変換方向の判定と実行
	if (inputExt == ".mp4" || inputExt == ".mkv") && outputExt == ".gif" {
		return s.ConvertMP4ToGIF()
	} else if inputExt == ".gif" && outputExt == ".mp4" {
		return s.ConvertGIFToMP4()
	} else {
		return fmt.Errorf("サポートされていない変換: %s -> %s", inputExt, outputExt)
	}
}

// setMP4ToGIFDefaults sets default values for MP4 to GIF conversion
func (s *MovieConverterService) setMP4ToGIFDefaults() {
	if s.config.FPS == 0 {
		s.config.FPS = 60 // バッチファイルのデフォルト
	}
	if s.config.Speed == 0 {
		s.config.Speed = 2.0 // バッチファイルのデフォルト
	}
	// Width: 0はそのまま（デフォルト品質）
	// Loop: 0はそのまま（無限ループ）
	// UseItsScale: trueがデフォルト
}

// setGIFToMP4Defaults sets default values for GIF to MP4 conversion
func (s *MovieConverterService) setGIFToMP4Defaults() {
	if s.config.FPS == 0 {
		s.config.FPS = 15 // PowerShellスクリプトのデフォルト
	}
}

// GenerateOutputFile generates output filename if not provided
func GenerateOutputFile(inputFile string) string {
	ext := strings.ToLower(filepath.Ext(inputFile))
	base := strings.TrimSuffix(inputFile, ext)

	switch ext {
	case ".mp4", ".mkv":
		return base + ".gif"
	case ".gif":
		return base + ".mp4"
	default:
		return base + "_converted"
	}
}

// ValidateConfig validates the conversion configuration
func ValidateConfig(config ConversionConfig) error {
	if config.InputFile == "" {
		return fmt.Errorf("入力ファイルが指定されていません")
	}

	if !strings.Contains(config.InputFile, ".") {
		return fmt.Errorf("入力ファイル名に拡張子が含まれていません: %s", config.InputFile)
	}

	if _, err := os.Stat(config.InputFile); os.IsNotExist(err) {
		return fmt.Errorf("入力ファイルが見つかりません: %s", config.InputFile)
	}

	return nil
}

// #==============================================================#
// ##          Consts and Types for batch process                ##
// #==============================================================#
// BatchConversionConfig holds the configuration for batch video conversion
type BatchConversionConfig struct {
	InputDir    string
	InputExt    string
	OutputDir   string
	OutputExt   string
	Recursive   bool
	FPS         int
	Width       int     // GIF変換時のみ（0=デフォルト品質）
	Speed       float64 // GIF変換時のみ
	Loop        int     // GIF変換時のみ（0=無限ループ、-1=ループなし）
	UseItsScale bool    // GIF変換時のみ
}

// BatchConversionResult holds the result of batch conversion
type BatchConversionResult struct {
	TotalFiles   int
	SuccessCount int
	FailureCount int
	Results      []ConversionResult
	FailedFiles  []string
}

// BatchMovieConverterService handles batch video conversion operations
type BatchMovieConverterService struct {
	config BatchConversionConfig
}

// NewBatchMovieConverterService creates a new BatchMovieConverterService instance
func NewBatchMovieConverterService(config BatchConversionConfig) *BatchMovieConverterService {
	return &BatchMovieConverterService{config: config}
}

// BatchConvert performs batch conversion of multiple files
func (bs *BatchMovieConverterService) BatchConvert() (*BatchConversionResult, error) {
	log.Printf("バッチ変換を開始: %s (%s) -> %s (%s)", bs.config.InputDir, bs.config.InputExt, bs.config.OutputDir, bs.config.OutputExt)

	// 入力ディレクトリの存在確認
	if _, err := os.Stat(bs.config.InputDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("入力ディレクトリが見つかりません: %s", bs.config.InputDir)
	}

	// 出力ディレクトリの作成
	if err := os.MkdirAll(bs.config.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	// ファイルリストの取得
	files, err := bs.scanFiles()
	if err != nil {
		return nil, fmt.Errorf("ファイルスキャンエラー: %w", err)
	}

	if len(files) == 0 {
		return &BatchConversionResult{
			TotalFiles:   0,
			SuccessCount: 0,
			FailureCount: 0,
			Results:      []ConversionResult{},
			FailedFiles:  []string{},
		}, nil
	}

	log.Printf("変換対象ファイル数: %d", len(files))

	// バッチ変換の実行
	result := &BatchConversionResult{
		TotalFiles:  len(files),
		Results:     make([]ConversionResult, 0, len(files)),
		FailedFiles: make([]string, 0),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	// 並列処理のためのセマフォ（最大4並列）
	semaphore := make(chan struct{}, 4)

	for _, inputFile := range files {
		wg.Add(1)
		go func(file string) {
			defer wg.Done()
			semaphore <- struct{}{}        // セマフォ取得
			defer func() { <-semaphore }() // セマフォ解放

			convResult := bs.convertSingleFile(file)

			mu.Lock()
			result.Results = append(result.Results, convResult)
			if convResult.Success {
				result.SuccessCount++
			} else {
				result.FailureCount++
				result.FailedFiles = append(result.FailedFiles, convResult.InputFile)
			}
			mu.Unlock()
		}(inputFile)
	}

	wg.Wait()

	log.Printf("バッチ変換完了: 成功 %d/%d, 失敗 %d", result.SuccessCount, result.TotalFiles, result.FailureCount)
	return result, nil
}

// scanFiles scans the input directory for files with the specified extension
func (bs *BatchMovieConverterService) scanFiles() ([]string, error) {
	var files []string
	inputExt := strings.ToLower(bs.config.InputExt)

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 再帰処理が無効で、サブディレクトリの場合はスキップ
			if !bs.config.Recursive && path != bs.config.InputDir {
				return filepath.SkipDir
			}
			return nil
		}

		// 拡張子のチェック
		if strings.ToLower(filepath.Ext(path)) == inputExt {
			files = append(files, path)
		}

		return nil
	}

	err := filepath.Walk(bs.config.InputDir, walkFunc)
	return files, err
}

// convertSingleFile converts a single file
func (bs *BatchMovieConverterService) convertSingleFile(inputFile string) ConversionResult {
	// 出力ファイル名の生成
	relPath, err := filepath.Rel(bs.config.InputDir, inputFile)
	if err != nil {
		return ConversionResult{
			InputFile:  inputFile,
			OutputFile: "",
			Success:    false,
			Error:      fmt.Errorf("相対パス取得エラー: %w", err),
		}
	}

	// 拡張子を変更
	outputFile := filepath.Join(bs.config.OutputDir, relPath)
	outputFile = strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + bs.config.OutputExt

	// 出力ディレクトリの作成
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return ConversionResult{
			InputFile:  inputFile,
			OutputFile: outputFile,
			Success:    false,
			Error:      fmt.Errorf("出力ディレクトリ作成エラー: %w", err),
		}
	}

	// 変換設定の作成
	config := ConversionConfig{
		InputFile:   inputFile,
		OutputFile:  outputFile,
		FPS:         bs.config.FPS,
		Width:       bs.config.Width,
		Speed:       bs.config.Speed,
		Loop:        bs.config.Loop,
		UseItsScale: bs.config.UseItsScale,
	}

	// 変換の実行
	service := NewMovieConverterService(config)
	err = service.Convert()

	return ConversionResult{
		InputFile:  inputFile,
		OutputFile: outputFile,
		Success:    err == nil,
		Error:      err,
	}
}

// GetSupportedExtensions returns supported file extensions
func GetSupportedExtensions() map[string][]string {
	return map[string][]string{
		"input":  {".mp4", ".mkv", ".gif"},
		"output": {".mp4", ".gif"},
	}
}

// normalizeExtension normalizes file extension by adding dot if missing
func normalizeExtension(ext string) string {
	if ext == "" {
		return ext
	}
	if !strings.HasPrefix(ext, ".") {
		return "." + ext
	}
	return ext
}

// ValidateBatchConfig validates the batch conversion configuration
func ValidateBatchConfig(config *BatchConversionConfig) error {
	if config.InputDir == "" {
		return fmt.Errorf("入力ディレクトリが指定されていません")
	}

	if config.InputExt == "" {
		return fmt.Errorf("入力拡張子が指定されていません")
	}

	if config.OutputDir == "" {
		return fmt.Errorf("出力ディレクトリが指定されていません")
	}

	if config.OutputExt == "" {
		return fmt.Errorf("出力拡張子が指定されていません")
	}

	// 拡張子の正規化（ドットを自動追加）
	config.InputExt = normalizeExtension(config.InputExt)
	config.OutputExt = normalizeExtension(config.OutputExt)

	// サポートされている拡張子の確認
	supportedExts := GetSupportedExtensions()
	inputSupported := false
	for _, ext := range supportedExts["input"] {
		if strings.ToLower(config.InputExt) == ext {
			inputSupported = true
			break
		}
	}
	if !inputSupported {
		return fmt.Errorf("サポートされていない入力拡張子: %s", config.InputExt)
	}

	outputSupported := false
	for _, ext := range supportedExts["output"] {
		if strings.ToLower(config.OutputExt) == ext {
			outputSupported = true
			break
		}
	}
	if !outputSupported {
		return fmt.Errorf("サポートされていない出力拡張子: %s", config.OutputExt)
	}

	// 入力ディレクトリの存在確認
	if _, err := os.Stat(config.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("入力ディレクトリが見つかりません: %s", config.InputDir)
	}

	return nil
}
