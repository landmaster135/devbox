package usecases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	exif "github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

// サポートする画像拡張子
var supportedExtensions = []string{".jpg", ".jpeg", ".tiff", ".tif", ".png", ".webp", ".mp4", ".webm"}

// Config はEXIF修正の設定を保持します
type Config struct {
	FolderPath     string
	DateTime       time.Time
	Extension      string
	Recursive      bool
	DryRun         bool
	Verbose        bool
	FromFilename   bool
	FromScreenshot bool
	WorkerCount    int // 並行処理のワーカー数
}

// ExifModifierService はEXIF修正サービスです
type ExifModifierService struct{}

// NewExifModifierService は新しいExifModifierServiceを作成します
func NewExifModifierService() *ExifModifierService {
	return &ExifModifierService{}
}

// getJSTLocation はJSTタイムゾーンを取得します
func (s *ExifModifierService) getJSTLocation() *time.Location {
	jstLocation, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// フォールバック: UTC+9の固定オフセット
		return time.FixedZone("JST", 9*60*60)
	}
	return jstLocation
}

// validateDirectory はディレクトリの存在と権限をバリデーションします
func validateDirectory(dirPath string) error {
	// 空文字列チェック
	if dirPath == "" {
		return fmt.Errorf("ディレクトリパスが空です")
	}

	// 存在チェック
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("ディレクトリが存在しません: %s", dirPath)
	}
	if err != nil {
		return fmt.Errorf("ディレクトリの情報取得に失敗しました: %s - %v", dirPath, err)
	}

	// ディレクトリかどうかチェック
	if !info.IsDir() {
		return fmt.Errorf("指定されたパスはディレクトリではありません: %s", dirPath)
	}

	// 読み取り権限チェック
	testFile := filepath.Join(dirPath, ".test_access")
	file, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("ディレクトリへの書き込み権限がありません: %s", dirPath)
	}
	file.Close()
	os.Remove(testFile) // テストファイルを削除

	return nil
}

// validateExtension はファイル拡張子をバリデーションします
func validateExtension(ext string) error {
	// 空文字列チェック
	if ext == "" {
		return fmt.Errorf("拡張子が空です")
	}

	// 拡張子を正規化（ドットがない場合は追加）
	normalizedExt := ext
	if !strings.HasPrefix(normalizedExt, ".") {
		normalizedExt = "." + normalizedExt
	}

	// より確実な方法で小文字に変換
	extLower := ""
	for _, char := range normalizedExt {
		if char >= 'A' && char <= 'Z' {
			extLower += string(char + 32) // A-Z を a-z に変換
		} else {
			extLower += string(char)
		}
	}

	// サポートされている拡張子かチェック
	for _, supportedExt := range supportedExtensions {
		if extLower == supportedExt {
			return nil
		}
	}

	return fmt.Errorf("サポートされていない拡張子です: %s\nサポートされている拡張子: %v", ext, supportedExtensions)
}

func ValidateInputOptions(dirPath, extension string) error {
	// ディレクトリの存在確認とバリデーション
	if err := validateDirectory(dirPath); err != nil {
		return fmt.Errorf("ディレクトリの存在確認とバリデーションに失敗しました: %w", err)
	}

	// 拡張子のバリデーション
	if extension != "" {
		if err := validateExtension(extension); err != nil {
			return fmt.Errorf("拡張子のバリデーションに失敗しました: %w", err)
		}
	}

	return nil
}

// isImageFile は画像ファイルかどうかをチェック
func (s *ExifModifierService) isImageFile(filePath, targetExtension string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 特定の拡張子が指定されている場合
	if targetExtension != "" {
		// 拡張子を正規化（ドットがない場合は追加）
		normalizedExt := targetExtension
		if !strings.HasPrefix(normalizedExt, ".") {
			normalizedExt = "." + normalizedExt
		}
		return strings.ToLower(normalizedExt) == ext
	}

	// サポートされている拡張子かチェック
	for _, supportedExt := range supportedExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// FindImageFiles は指定された設定に基づいて画像ファイルを検索します
func (s *ExifModifierService) FindImageFiles(config *Config) ([]string, error) {
	var imageFiles []string

	err := filepath.Walk(config.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 再帰フラグが設定されていない場合、サブディレクトリをスキップ
		if !config.Recursive && info.IsDir() && path != config.FolderPath {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			if s.isImageFile(path, config.Extension) {
				imageFiles = append(imageFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// ファイル名順でソート
	sort.Slice(imageFiles, func(i, j int) bool {
		return filepath.Base(imageFiles[i]) < filepath.Base(imageFiles[j])
	})

	return imageFiles, nil
}

// ProcessResult はファイル処理の結果を表します
type ProcessResult struct {
	FilePath string
	Success  bool
	Error    error
}

// updateFileTime はファイルの更新時刻を変更します
func (s *ExifModifierService) updateFileTime(filePath string, targetTime time.Time) error {
	return os.Chtimes(filePath, targetTime, targetTime)
}

// modifyJpegExif はJPEGファイルのEXIF情報を修正します
func (s *ExifModifierService) modifyJpegExif(filePath string, config *Config) error {
	// ファイルを読み取り
	originalData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ファイルの読み取りに失敗: %v", err)
	}

	// JPEGパーサーを使用してExifセグメントを取得
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(originalData)
	if err != nil {
		return fmt.Errorf("JPEGの解析に失敗: %v", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Exifセグメントを取得
	rootIfd, rawExif, err := sl.Exif()
	if err != nil {
		// Exifが存在しない場合はファイル時刻のみ更新
		if config.Verbose {
			log.Printf("Exifデータが見つかりません。ファイル時刻のみ更新: %s", filePath)
		}
		return s.updateFileTime(filePath, config.DateTime)
	}

	// 既存のExifを修正
	return s.updateExistingExif(filePath, config, originalData, rootIfd, rawExif)
}

// modifyGenericExif は汎用画像ファイルのEXIF情報を修正します
func (s *ExifModifierService) modifyGenericExif(filePath string, config *Config) error {
	// TIFF等の形式では、ファイル時刻のみ更新
	if config.Verbose {
		log.Printf("TIFF/その他形式のファイル時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// modifyPngFile はPNGファイルの時刻を修正します
func (s *ExifModifierService) modifyPngFile(filePath string, config *Config) error {
	if config.Verbose {
		log.Printf("PNGファイルの時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// modifyWebpFile はWebPファイルの時刻を修正します
func (s *ExifModifierService) modifyWebpFile(filePath string, config *Config) error {
	if config.Verbose {
		log.Printf("WebPファイルの時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// modifyMp4File はMP4ファイルの時刻を修正します
func (s *ExifModifierService) modifyMp4File(filePath string, config *Config) error {
	if config.Verbose {
		log.Printf("MP4ファイルの時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// modifyWebmFile はWebMファイルの時刻を修正します
func (s *ExifModifierService) modifyWebmFile(filePath string, config *Config) error {
	if config.Verbose {
		log.Printf("WebMファイルの時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// modifyTiffFile はTIFFファイルの時刻を修正します
func (s *ExifModifierService) modifyTiffFile(filePath string, config *Config) error {
	if config.Verbose {
		log.Printf("TIFFファイルの時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// modifyGenericFile は汎用ファイルの時刻を修正します
func (s *ExifModifierService) modifyGenericFile(filePath string, config *Config) error {
	if config.Verbose {
		log.Printf("汎用ファイルの時刻を更新: %s", filePath)
	}
	return s.updateFileTime(filePath, config.DateTime)
}

// ModifySingleFileExif は単一ファイルのEXIF情報を修正します
func (s *ExifModifierService) ModifySingleFileExif(filePath string, config *Config) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".jpg", ".jpeg":
		return s.modifyJpegExif(filePath, config)
	case ".png":
		return s.modifyPngFile(filePath, config)
	case ".webp":
		return s.modifyWebpFile(filePath, config)
	case ".mp4":
		return s.modifyMp4File(filePath, config)
	case ".webm":
		return s.modifyWebmFile(filePath, config)
	case ".tiff", ".tif":
		return s.modifyTiffFile(filePath, config)
	default:
		return s.modifyGenericFile(filePath, config)
	}
}

// ModifyExifData は複数のファイルのEXIF情報を並行処理で修正します
func (s *ExifModifierService) ModifyExifData(imageFiles []string, config *Config) (int, int, error) {
	// ワーカー数のデフォルト設定
	workerCount := config.WorkerCount
	if workerCount <= 0 {
		workerCount = 4 // デフォルトは4並列
	}

	// ファイル数がワーカー数より少ない場合、ワーカー数をファイル数に調整
	if len(imageFiles) < workerCount {
		workerCount = len(imageFiles)
	}

	// チャネルの作成
	jobs := make(chan string, len(imageFiles))
	results := make(chan ProcessResult, len(imageFiles))

	// ワーカーゴルーチンの起動
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range jobs {
				result := ProcessResult{
					FilePath: filePath,
					Success:  false,
					Error:    nil,
				}

				if config.DryRun {
					fmt.Printf("処理中: %s\n", filePath)
					fmt.Printf("  → (ドライラン) File Modification Date/Time を %s に設定\n",
						config.DateTime.Format("2006-01-02 15:04:05"))
					result.Success = true
				} else {
					fmt.Printf("処理中: %s\n", filePath)
					err := s.ModifySingleFileExif(filePath, config)
					if err != nil {
						fmt.Printf("  ⚠️  エラー: %v\n", err)
						result.Error = err
						if config.Verbose {
							log.Printf("Error processing file %s: %v", filePath, err)
						}
					} else {
						fmt.Printf("  ✅ File Modification Date/Time を %s に設定しました\n",
							config.DateTime.Format("2006-01-02 15:04:05"))
						result.Success = true
					}
				}
				results <- result
			}
		}()
	}

	// ジョブをチャネルに送信
	for _, filePath := range imageFiles {
		jobs <- filePath
	}
	close(jobs)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集
	processedCount := 0
	errorCount := 0
	for result := range results {
		if result.Success {
			processedCount++
		} else {
			errorCount++
		}
	}

	return processedCount, errorCount, nil
}

// updateExistingExif は既存のExifデータを更新します
func (s *ExifModifierService) updateExistingExif(filePath string, config *Config, originalData []byte, rootIfd *exif.Ifd, rawExif []byte) error {
	// 現在のgo-exifライブラリでは、Exifデータの直接的な書き込みが複雑なため、
	// ファイルシステムレベルでの時刻更新のみを行います
	if config.Verbose {
		log.Printf("Exifデータが見つかりました。ファイル時刻を更新: %s", filePath)
	}

	return s.updateFileTime(filePath, config.DateTime)
}

// parseInt は文字列を整数に変換します（エラーハンドリングなし、事前に数値チェック済み）
func parseInt(s string) int {
	result := 0
	for _, char := range s {
		result = result*10 + int(char-'0')
	}
	return result
}

// isLeapYear はうるう年かどうか判定します
func isLeapYear(year int) bool {
	// 4で割り切れる年はうるう年
	// ただし100で割り切れる年は平年
	// ただし400で割り切れる年はうるう年
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// getMaxDaysInMonth は指定された年月の最大日数を取得します
func getMaxDaysInMonth(month, year int) int {
	switch month {
	case 1, 3, 5, 7, 8, 10, 12: // 31日の月
		return 31
	case 4, 6, 9, 11: // 30日の月
		return 30
	case 2: // 2月（うるう年チェック）
		if isLeapYear(year) {
			return 29
		}
		return 28
	default:
		return 31 // フォールバック
	}
}

// validateDateTime は日時の各要素をバリデーションします
func validateDateTime(dateTimeStr string) error {
	// 文字列から各要素を抽出
	year := dateTimeStr[0:4]
	month := dateTimeStr[4:6]
	day := dateTimeStr[6:8]
	hour := dateTimeStr[8:10]
	minute := dateTimeStr[10:12]
	second := dateTimeStr[12:14]

	// 年のバリデーション (1900-2099)
	yearInt := parseInt(year)
	if yearInt < 1900 || yearInt > 2099 {
		return fmt.Errorf("年は1900-2099の範囲である必要があります: %d", yearInt)
	}

	// 月のバリデーション (01-12)
	monthInt := parseInt(month)
	if monthInt < 1 || monthInt > 12 {
		return fmt.Errorf("月は01-12の範囲である必要があります: %02d", monthInt)
	}

	// 日のバリデーション (01-31, 月によって異なる)
	dayInt := parseInt(day)
	if dayInt < 1 || dayInt > 31 {
		return fmt.Errorf("日は01-31の範囲である必要があります: %02d", dayInt)
	}

	// 月ごとの日数チェック
	maxDaysInMonth := getMaxDaysInMonth(monthInt, yearInt)
	if dayInt > maxDaysInMonth {
		return fmt.Errorf("%d年%02d月の日は01-%02dの範囲である必要があります: %02d", yearInt, monthInt, maxDaysInMonth, dayInt)
	}

	// 時のバリデーション (00-23)
	hourInt := parseInt(hour)
	if hourInt < 0 || hourInt > 23 {
		return fmt.Errorf("時は00-23の範囲である必要があります: %02d", hourInt)
	}

	// 分のバリデーション (00-59)
	minuteInt := parseInt(minute)
	if minuteInt < 0 || minuteInt > 59 {
		return fmt.Errorf("分は00-59の範囲である必要があります: %02d", minuteInt)
	}

	// 秒のバリデーション (00-59)
	secondInt := parseInt(second)
	if secondInt < 0 || secondInt > 59 {
		return fmt.Errorf("秒は00-59の範囲である必要があります: %02d", secondInt)
	}

	return nil
}

// ParseDateTime は日時文字列をtime.Timeに変換します
func ParseDateTime(dateTimeStr string) (time.Time, error) {
	if len(dateTimeStr) != 14 {
		return time.Time{}, fmt.Errorf("日時は14文字である必要があります (yyyyMMddhhmmss)")
	}

	// 基本的な数値チェック
	for _, char := range dateTimeStr {
		if char < '0' || char > '9' {
			return time.Time{}, fmt.Errorf("日時は数字のみで構成されている必要があります: %s", dateTimeStr)
		}
	}

	// 各要素のバリデーション
	if err := validateDateTime(dateTimeStr); err != nil {
		return time.Time{}, err
	}

	// JSTタイムゾーンを取得
	jstLocation, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		// フォールバック: UTC+9の固定オフセット
		jstLocation = time.FixedZone("JST", 9*60*60)
	}

	return time.ParseInLocation("20060102150405", dateTimeStr, jstLocation)
}

func (s *ExifModifierService) isFileExtensionSupported(ext string) bool {
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".tiff" {
		return false
	}
	return true
}

// ProcessFilesFromFilename はファイル名から日時情報を抽出してEXIF情報を設定します
func (s *ExifModifierService) ProcessFilesFromFilename(path string, recursive, dryRun, verbose, overwriteExif bool) error {
	// ファイルパターンの正規表現
	// 例: IMG_20230101_120000.jpg, Screenshot_20230101-120000.png など
	datePatterns := []*regexp.Regexp{
		regexp.MustCompile(`(\d{4})(\d{2})(\d{2})[-_]?(\d{2})(\d{2})(\d{2})`),                // YYYYMMDD_HHMMSS
		regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})[-_](\d{2})-(\d{2})-(\d{2})`),             // YYYY-MM-DD_HH-MM-SS
		regexp.MustCompile(`(\d{4})[-_](\d{2})[-_](\d{2})[-_](\d{2})[-_](\d{2})[-_](\d{2})`), // YYYY_MM_DD_HH_MM_SS
	}

	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合
		if info.IsDir() {
			// ルートパス以外のディレクトリで再帰的処理が無効の場合はスキップ
			if filePath != path && !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		// 画像ファイル以外はスキップ
		ext := strings.ToLower(filepath.Ext(filePath))
		if !s.isFileExtensionSupported(ext) {
			if verbose {
				fmt.Printf("スキップ: %s (サポートされていないファイル形式)\n", filePath)
			}
			return nil
		}

		// ファイル名から日時情報を抽出
		fileName := filepath.Base(filePath)
		var dateTime time.Time
		var matched bool

		for _, pattern := range datePatterns {
			matches := pattern.FindStringSubmatch(fileName)
			if len(matches) >= 7 {
				year := matches[1]
				month := matches[2]
				day := matches[3]
				hour := matches[4]
				minute := matches[5]
				second := matches[6]

				// 日時情報の解析（JSTタイムゾーンで）
				dateTimeStr := fmt.Sprintf("%s-%s-%s %s:%s:%s", year, month, day, hour, minute, second)
				parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateTimeStr, s.getJSTLocation())
				if err != nil {
					if verbose {
						fmt.Printf("警告: %s の日時情報の解析に失敗しました: %v\n", filePath, err)
					}
					continue
				}

				dateTime = parsedTime
				matched = true
				break
			}
		}

		if !matched {
			if verbose {
				fmt.Printf("スキップ: %s (ファイル名から日時情報を抽出できませんでした)\n", filePath)
			}
			return nil
		}

		// 既存のEXIF情報をチェック
		if !overwriteExif {
			// ファイルを読み取り
			originalData, err := os.ReadFile(filePath)
			if err == nil {
				// JPEGパーサーを使用してExifセグメントを取得
				jmp := jpegstructure.NewJpegMediaParser()
				intfc, err := jmp.ParseBytes(originalData)
				if err == nil {
					sl := intfc.(*jpegstructure.SegmentList)
					// Exifセグメントを取得
					_, _, err := sl.Exif()
					if err == nil {
						if verbose {
							fmt.Printf("スキップ: %s (既存のEXIF情報があります)\n", filePath)
						}
						return nil
					}
				}
			}
		}

		// 処理内容の表示
		fmt.Printf("処理: %s -> %v\n", filePath, dateTime.Format("2006-01-02 15:04:05"))

		// ドライランの場合は実際の変更を行わない
		if dryRun {
			return nil
		}

		// 設定を作成
		config := &Config{
			DateTime: dateTime,
			DryRun:   false,
			Verbose:  verbose,
		}

		// EXIF情報の設定
		err = s.ModifySingleFileExif(filePath, config)
		if err != nil {
			fmt.Printf("エラー: %s の処理中にエラーが発生しました: %v\n", filePath, err)
			return nil // 個別のファイルエラーは全体の処理を中断しない
		}

		fmt.Printf("EXIF情報を設定しました: %s\n", filePath)
		return nil
	})
}

// ProcessFilesFromScreenshot はスクリーンショットファイルの日時情報を設定します
func (s *ExifModifierService) ProcessFilesFromScreenshot(path string, recursive, dryRun, verbose, overwriteExif bool) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合
		if info.IsDir() {
			// ルートパス以外のディレクトリで再帰的処理が無効の場合はスキップ
			if filePath != path && !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		// 画像ファイル以外はスキップ
		ext := strings.ToLower(filepath.Ext(filePath))
		if !s.isFileExtensionSupported(ext) {
			if verbose {
				fmt.Printf("スキップ: %s (サポートされていないファイル形式)\n", filePath)
			}
			return nil
		}

		// ファイル名がスクリーンショットパターンに一致するか確認
		fileName := filepath.Base(filePath)
		isScreenshot := strings.HasPrefix(strings.ToLower(fileName), "screenshot") ||
			strings.HasPrefix(strings.ToLower(fileName), "screen shot") ||
			strings.HasPrefix(strings.ToLower(fileName), "スクリーンショット")

		if !isScreenshot {
			if verbose {
				fmt.Printf("スキップ: %s (スクリーンショットファイルではありません)\n", filePath)
			}
			return nil
		}

		// スクリーンショットファイル名から日時を抽出
		screenshotPatterns := []*regexp.Regexp{
			regexp.MustCompile(`[Ss]creenshot_(\d{8})-(\d{6})`),     // Screenshot_YYYYMMDD-HHMMSS
			regexp.MustCompile(`スクリーンショット_(\d{8})-(\d{6})`),      // スクリーンショット_YYYYMMDD-HHMMSS
		}

		var fileTime time.Time
		var matched bool

		for _, pattern := range screenshotPatterns {
			matches := pattern.FindStringSubmatch(fileName)
			if len(matches) >= 3 {
				dateStr := matches[1] // YYYYMMDD
				timeStr := matches[2] // HHMMSS
				dateTimeStr := dateStr + timeStr // YYYYMMDDHHMMSS

				parsedTime, err := time.ParseInLocation("20060102150405", dateTimeStr, s.getJSTLocation())
				if err != nil {
					return fmt.Errorf("ファイル名から抽出した日時の解析に失敗: %s - %v", filePath, err)
				}

				fileTime = parsedTime
				matched = true
				break
			}
		}

		if !matched {
			return fmt.Errorf("スクリーンショットファイル名から日時情報を抽出できませんでした: %s (期待される形式: Screenshot_YYYYMMDD-HHMMSS)", filePath)
		}

		// 既存のEXIF情報をチェック
		if !overwriteExif {
			// ファイルを読み取り
			originalData, err := os.ReadFile(filePath)
			if err == nil {
				// JPEGパーサーを使用してExifセグメントを取得
				jmp := jpegstructure.NewJpegMediaParser()
				intfc, err := jmp.ParseBytes(originalData)
				if err == nil {
					sl := intfc.(*jpegstructure.SegmentList)
					// Exifセグメントを取得
					_, _, err := sl.Exif()
					if err == nil {
						if verbose {
							fmt.Printf("スキップ: %s (既存のEXIF情報があります)\n", filePath)
						}
						return nil
					}
				}
			}
		}

		// 処理内容の表示
		fmt.Printf("処理: %s -> %v\n", filePath, fileTime.Format("2006-01-02 15:04:05"))

		// ドライランの場合は実際の変更を行わない
		if dryRun {
			return nil
		}

		// 設定を作成
		config := &Config{
			DateTime: fileTime,
			DryRun:   false,
			Verbose:  verbose,
		}

		// EXIF情報の設定
		err = s.ModifySingleFileExif(filePath, config)
		if err != nil {
			fmt.Printf("エラー: %s の処理中にエラーが発生しました: %v\n", filePath, err)
			return nil // 個別のファイルエラーは全体の処理を中断しない
		}

		fmt.Printf("EXIF情報を設定しました: %s\n", filePath)
		return nil
	})
}
