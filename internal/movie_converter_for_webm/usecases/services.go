package usecases

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// #==============================================================#
// ##          Consts and Types for single process               ##
// #==============================================================#
// ConversionConfig holds the configuration for video conversion
type ConversionConfig struct {
	InputFile      string
	OutputFile     string
	VideoBitrate   string // ビデオビットレート（例: "1M", "1500k"）
	AudioBitrate   string // オーディオビットレート（例: "128k"）
	CRF            int    // Constant Rate Factor（0-63、低いほど高品質）
	AudioCodec     string // オーディオコーデック（opus/vorbis）
	ConversionMode string // 変換モード（crf/cbr）
	VideoQuality   int    // ビデオ品質（CBRモード時、デフォルト75）
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

// setMP4ToWEBMDefaults sets default values for MP4 to WEBM conversion
func (s *MovieConverterService) setMP4ToWEBMDefaults() {
	if s.config.AudioBitrate == "" {
		s.config.AudioBitrate = "128k"
	}
	if s.config.AudioCodec == "" {
		s.config.AudioCodec = "opus"
	}
	if s.config.ConversionMode == "" {
		s.config.ConversionMode = "crf"
	}
	if s.config.CRF == 0 && s.config.ConversionMode == "crf" {
		s.config.CRF = 30 // デフォルトCRF値
	}
	if s.config.VideoQuality == 0 {
		s.config.VideoQuality = 75
	}
}

// getSourceVideoBitrate retrieves the source video bitrate using ffprobe
func (s *MovieConverterService) getSourceVideoBitrate() (string, error) {
	// ffprobe でソース動画のビットレートを取得
	probeData, err := ffmpeg.Probe(s.config.InputFile)
	if err != nil {
		log.Printf("ffprobeでビットレート取得に失敗: %v", err)
		return "1M", nil // デフォルト値を返す
	}

	// JSONデータの構造体定義
	type probeStream struct {
		CodecType string `json:"codec_type"`
		BitRate   string `json:"bit_rate"`
	}

	type probeResult struct {
		Streams []probeStream `json:"streams"`
	}

	// JSONをパース
	var result probeResult
	if err := json.Unmarshal([]byte(probeData), &result); err != nil {
		log.Printf("ffprobe結果のパースに失敗: %v", err)
		return "1M", nil
	}

	// ビデオストリームを検索
	for _, stream := range result.Streams {
		if stream.CodecType == "video" && stream.BitRate != "" {
			if bitRate, err := strconv.ParseFloat(stream.BitRate, 64); err == nil {
				// bit/s から kbit/s に変換し、品質調整を適用
				bitrateKbps := int(bitRate / 1000 * float64(s.config.VideoQuality) / 100)
				return fmt.Sprintf("%dk", bitrateKbps), nil
			}
		}
	}

	log.Println("ソース動画のビットレートを取得できませんでした。デフォルト値1Mを使用します。")
	return "1M", nil
}

// convertMP4ToWEBM converts MP4 to WEBM using ffmpeg-go
func (s *MovieConverterService) convertMP4ToWEBM() error {
	log.Printf("MP4からWEBMに変換中: %s -> %s", s.config.InputFile, s.config.OutputFile)

	// デフォルト値の設定
	s.setMP4ToWEBMDefaults()

	// ファイル名にスペースが含まれている場合の警告
	if strings.Contains(s.config.InputFile, " ") {
		log.Println("警告: ソースファイル名にスペース文字が含まれています。問題が発生する可能性があります。")
	}

	// オーディオコーデックの設定
	audioCodec := "libopus"
	if s.config.AudioCodec == "vorbis" {
		audioCodec = "libvorbis"
	}

	var err error
	if s.config.ConversionMode == "crf" {
		// CRFモード: 固定品質
		log.Printf("CRFモードを使用: CRF=%d", s.config.CRF)
		err = ffmpeg.Input(s.config.InputFile).
			Output(s.config.OutputFile, ffmpeg.KwArgs{
				"c:v": "libvpx-vp9",
				"crf": s.config.CRF,
				"b:v": "0", // CRFモードでは0に設定
				"c:a": audioCodec,
				"b:a": s.config.AudioBitrate,
			}).
			OverWriteOutput().
			Run()
	} else {
		// CBRモード: 固定ビットレート
		videoBitrate := s.config.VideoBitrate
		if videoBitrate == "" {
			// ソース動画のビットレートを取得
			videoBitrate, _ = s.getSourceVideoBitrate()
		}
		log.Printf("CBRモードを使用: ビットレート=%s", videoBitrate)
		err = ffmpeg.Input(s.config.InputFile).
			Output(s.config.OutputFile, ffmpeg.KwArgs{
				"c:v": "libvpx-vp9",
				"b:v": videoBitrate,
				"c:a": audioCodec,
				"b:a": s.config.AudioBitrate,
			}).
			OverWriteOutput().
			Run()
	}

	if err != nil {
		return fmt.Errorf("ffmpeg実行エラー: %w", err)
	}

	log.Printf("MP4からWEBMへの変換が完了しました: %s", s.config.OutputFile)
	return nil
}

// setWEBMToMP4Defaults sets default values for WEBM to MP4 conversion
func (s *MovieConverterService) setWEBMToMP4Defaults() {
	// WEBM to MP4変換では特別なデフォルト設定は不要
}

// convertWEBMToMP4 converts WEBM to MP4 using ffmpeg-go
func (s *MovieConverterService) convertWEBMToMP4() error {
	log.Printf("WEBMからMP4に変換中: %s -> %s", s.config.InputFile, s.config.OutputFile)

	// デフォルト値の設定
	s.setWEBMToMP4Defaults()

	// ファイル名にスペースが含まれている場合の警告
	if strings.Contains(s.config.InputFile, " ") {
		log.Println("警告: ソースファイル名にスペース文字が含まれています。問題が発生する可能性があります。")
	}

	err := ffmpeg.Input(s.config.InputFile).
		Output(s.config.OutputFile, ffmpeg.KwArgs{
			"c:v":      "libx264",
			"c:a":      "aac",
			"movflags": "faststart",
		}).
		OverWriteOutput().
		Run()

	if err != nil {
		return fmt.Errorf("ffmpeg実行エラー: %w", err)
	}

	log.Printf("WEBMからMP4への変換が完了しました: %s", s.config.OutputFile)
	return nil
}

// convert performs conversion based on file extensions
func (s *MovieConverterService) convert() error {
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
	supportedInputExts := []string{".mp4", ".mkv", ".avi", ".mov", ".flv"}
	if slices.Contains(supportedInputExts, inputExt) && outputExt == ".webm" {
		return s.convertMP4ToWEBM()
	} else if inputExt == ".webm" && outputExt == ".mp4" {
		return s.convertWEBMToMP4()
	} else {
		return fmt.Errorf("サポートされていない変換: %s -> %s", inputExt, outputExt)
	}
}

// #==============================================================#
// ##          Consts and Types for batch process                ##
// #==============================================================#
// BatchConversionConfig holds the configuration for batch video conversion
type BatchConversionConfig struct {
	InputDir       string
	InputExt       string
	OutputDir      string
	OutputExt      string
	Recursive      bool
	VideoBitrate   string // ビデオビットレート（例: "1M", "1500k"）
	AudioBitrate   string // オーディオビットレート（例: "128k"）
	CRF            int    // Constant Rate Factor（0-63、低いほど高品質）
	AudioCodec     string // オーディオコーデック（opus/vorbis）
	ConversionMode string // 変換モード（crf/cbr）
	VideoQuality   int    // ビデオ品質（CBRモード時、デフォルト75）
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

// batchConvert performs batch conversion of multiple files
func (bs *BatchMovieConverterService) batchConvert() (*BatchConversionResult, error) {
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
		InputFile:      inputFile,
		OutputFile:     outputFile,
		VideoBitrate:   bs.config.VideoBitrate,
		AudioBitrate:   bs.config.AudioBitrate,
		CRF:            bs.config.CRF,
		AudioCodec:     bs.config.AudioCodec,
		ConversionMode: bs.config.ConversionMode,
		VideoQuality:   bs.config.VideoQuality,
	}

	// 変換の実行
	service := NewMovieConverterService(config)
	err = service.convert()

	return ConversionResult{
		InputFile:  inputFile,
		OutputFile: outputFile,
		Success:    err == nil,
		Error:      err,
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

// GetSupportedExtensions returns supported file extensions
func GetSupportedExtensions() map[string][]string {
	return map[string][]string{
		"input":  {".mp4", ".mkv", ".avi", ".mov", ".flv", ".webm"},
		"output": {".mp4", ".webm"},
	}
}

// #==============================================================#
// ##          Unified Service for both single and batch         ##
// #==============================================================#

// ProcessingMode represents the type of processing to be performed
type ProcessingMode int

const (
	SingleFileMode ProcessingMode = iota
	BatchMode
)

// UnifiedConversionConfig holds configuration for both single and batch processing
type UnifiedConversionConfig struct {
	// Single file processing config
	SingleConfig *ConversionConfig
	// Batch processing config
	BatchConfig *BatchConversionConfig
	// Processing mode
	Mode ProcessingMode
}

// UnifiedConversionResult holds the result of unified conversion
type UnifiedConversionResult struct {
	Mode    ProcessingMode
	Success bool
	Error   error
	// Single file result
	SingleResult *ConversionResult
	// Batch result
	BatchResult *BatchConversionResult
}

// UnifiedMovieConverterService handles both single file and batch video conversion operations
type UnifiedMovieConverterService struct {
	config UnifiedConversionConfig
}

// NewUnifiedMovieConverterService creates a new UnifiedMovieConverterService instance
func NewUnifiedMovieConverterService(singleConfig *ConversionConfig, batchConfig *BatchConversionConfig) *UnifiedMovieConverterService {
	var mode ProcessingMode

	// Determine processing mode based on provided configurations
	if batchConfig != nil && (batchConfig.InputDir != "" || batchConfig.InputExt != "" || batchConfig.OutputDir != "" || batchConfig.OutputExt != "") {
		mode = BatchMode
	} else {
		mode = SingleFileMode
	}

	return &UnifiedMovieConverterService{
		config: UnifiedConversionConfig{
			SingleConfig: singleConfig,
			BatchConfig:  batchConfig,
			Mode:         mode,
		},
	}
}

// #==============================================================#
// ##          Methods of Unified Service                        ##
// #==============================================================#
// validateSingleConfig validates the conversion configuration
func validateSingleConfig(config ConversionConfig) error {
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

// generateOutputFile generates output filename if not provided
func generateOutputFile(inputFile string) string {
	ext := strings.ToLower(filepath.Ext(inputFile))
	originalExt := filepath.Ext(inputFile)
	base := strings.TrimSuffix(inputFile, originalExt)

	supportedInputExts := []string{".mp4", ".mkv", ".avi", ".mov", ".flv"}
	if slices.Contains(supportedInputExts, ext) {
		return base + ".webm"
	} else if ext == ".webm" {
		return base + ".mp4"
	} else {
		return base + "_converted"
	}
}

// processSingleFile handles single file conversion
func (us *UnifiedMovieConverterService) processSingleFile(result *UnifiedConversionResult) (*UnifiedConversionResult, error) {
	if us.config.SingleConfig == nil {
		result.Success = false
		result.Error = fmt.Errorf("単一ファイル処理の設定が指定されていません")
		return result, result.Error
	}

	// 設定の検証
	if err := validateSingleConfig(*us.config.SingleConfig); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("設定エラー: %w", err)
		return result, result.Error
	}

	// 出力ファイルの自動生成
	if us.config.SingleConfig.OutputFile == "" {
		us.config.SingleConfig.OutputFile = generateOutputFile(us.config.SingleConfig.InputFile)
		log.Printf("出力ファイル名を自動生成しました: %s", us.config.SingleConfig.OutputFile)
	}

	// 変換サービスの作成と実行
	service := NewMovieConverterService(*us.config.SingleConfig)
	err := service.convert()

	singleResult := &ConversionResult{
		InputFile:  us.config.SingleConfig.InputFile,
		OutputFile: us.config.SingleConfig.OutputFile,
		Success:    err == nil,
		Error:      err,
	}

	result.SingleResult = singleResult
	result.Success = err == nil
	result.Error = err

	if err == nil {
		log.Printf("変換が正常に完了しました: %s -> %s", us.config.SingleConfig.InputFile, us.config.SingleConfig.OutputFile)
	}

	return result, err
}

// validateBatchConfig validates the batch conversion configuration
func validateBatchConfig(config *BatchConversionConfig) error {
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
	if !slices.Contains(supportedExts["input"], strings.ToLower(config.InputExt)) {
		return fmt.Errorf("サポートされていない入力拡張子: %s", config.InputExt)
	}
	if !slices.Contains(supportedExts["output"], strings.ToLower(config.OutputExt)) {
		return fmt.Errorf("サポートされていない出力拡張子: %s", config.OutputExt)
	}

	// 入力ディレクトリの存在確認
	if _, err := os.Stat(config.InputDir); os.IsNotExist(err) {
		return fmt.Errorf("入力ディレクトリが見つかりません: %s", config.InputDir)
	}

	return nil
}

// processBatchFiles handles batch file conversion
func (us *UnifiedMovieConverterService) processBatchFiles(result *UnifiedConversionResult) (*UnifiedConversionResult, error) {
	if us.config.BatchConfig == nil {
		result.Success = false
		result.Error = fmt.Errorf("バッチ処理の設定が指定されていません")
		return result, result.Error
	}

	// バッチ設定の検証
	if err := validateBatchConfig(us.config.BatchConfig); err != nil {
		result.Success = false
		result.Error = fmt.Errorf("バッチ設定エラー: %w", err)
		return result, result.Error
	}

	// バッチ変換サービスの作成と実行
	batchService := NewBatchMovieConverterService(*us.config.BatchConfig)
	batchResult, err := batchService.batchConvert()

	result.BatchResult = batchResult
	result.Success = err == nil
	result.Error = err

	if err == nil {
		log.Printf("バッチ変換が完了しました: %s (%s) -> %s (%s)",
			us.config.BatchConfig.InputDir, us.config.BatchConfig.InputExt,
			us.config.BatchConfig.OutputDir, us.config.BatchConfig.OutputExt)
	}

	return result, err
}

// ProcessConversion performs conversion based on the configuration mode
func (us *UnifiedMovieConverterService) ProcessConversion() (*UnifiedConversionResult, error) {
	result := &UnifiedConversionResult{
		Mode: us.config.Mode,
	}

	switch us.config.Mode {
	case SingleFileMode:
		return us.processSingleFile(result)
	case BatchMode:
		return us.processBatchFiles(result)
	default:
		result.Success = false
		result.Error = fmt.Errorf("不明な処理モード: %d", us.config.Mode)
		return result, result.Error
	}
}
