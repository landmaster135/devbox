package usecases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

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

// GetSupportedExtensions returns supported file extensions
func GetSupportedExtensions() map[string][]string {
	return map[string][]string{
		"input":  {".mp4", ".mkv", ".gif"},
		"output": {".mp4", ".gif"},
	}
}
