package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/landmaster135/devbox/internal/independencies/exif_modifier/usecases"
)

const version = "1.0.0"

// コマンドライン引数
var (
	dirPath        = flag.String("dir", ".", "画像ファイルがあるディレクトリのパス")
	dateTime       = flag.String("datetime", "", "設定する日時 (yyyyMMddhhmmss形式)")
	fromFilename   = flag.Bool("from-filename", false, "ファイル名から日時を取得してExifに設定する (ファイル名がyyyyMMddhhmmss形式の場合)")
	fromScreenshot = flag.Bool("from-screenshot", false, "スクリーンショットファイル名から日時を取得してExifに設定する (Screenshot_yyyyMMdd-hhmmss形式の場合)")
	extension      = flag.String("ext", "", "対象とする拡張子 (例: .jpg, .jpeg, .png, .webp, .mp4)")
	recursive      = flag.Bool("recursive", false, "サブフォルダも再帰的に処理する")
	dryRun         = flag.Bool("dry-run", false, "実際には変更せず、処理対象ファイルのみ表示")
	verbose        = flag.Bool("verbose", false, "詳細な出力を表示")
	showHelp       = flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion    = flag.Bool("version", false, "バージョン情報を表示")
)

func main() {
	// カスタムUsage関数を設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "exif-modifier - EXIF Property Modifier Tool\n\n")
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --datetime 20240315143000\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --ext .jpg\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --recursive\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --dry-run\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --verbose\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --from-filename --ext .jpg\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --from-filename --recursive --dry-run\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./screenshots --from-screenshot --ext .png\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./screenshots --from-screenshot --recursive --dry-run\n")
	}

	flag.Parse()

	// ヘルプまたはバージョン表示
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("exif-modifier version %s\n", version)
		os.Exit(0)
	}

	// 引数の検証
	activeOptionsCount := 0
	if *dateTime != "" {
		activeOptionsCount++
	}
	if *fromFilename {
		activeOptionsCount++
	}
	if *fromScreenshot {
		activeOptionsCount++
	}

	if activeOptionsCount == 0 {
		fmt.Fprintf(os.Stderr, "Error: --datetime, --from-filename, --from-screenshot のいずれかのパラメータが必要です\n")
		flag.Usage()
		os.Exit(1)
	}

	if activeOptionsCount > 1 {
		fmt.Fprintf(os.Stderr, "Error: --datetime, --from-filename, --from-screenshot は同時に指定できません\n")
		flag.Usage()
		os.Exit(1)
	}

	var targetTime time.Time
	var err error

	// 日時のパース
	if *dateTime != "" {
		targetTime, err = parseDateTime(*dateTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 日時の形式が正しくありません: %v\n", err)
			os.Exit(1)
		}
	}

	// ディレクトリの存在確認とバリデーション
	if err := validateDirectory(*dirPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 拡張子のバリデーション
	if *extension != "" {
		if err := validateExtension(*extension); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// ログ設定
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	// 設定を作成
	config := &usecases.Config{
		FolderPath:     *dirPath,
		DateTime:       targetTime,
		Extension:      *extension,
		Recursive:      *recursive,
		DryRun:         *dryRun,
		Verbose:        *verbose,
		FromFilename:   *fromFilename,
		FromScreenshot: *fromScreenshot,
	}

	// 実行情報を表示
	printExecutionInfo(config)

	// ExifModifierServiceを作成
	service := usecases.NewExifModifierService()

	// 画像ファイルを検索
	imageFiles, err := service.FindImageFiles(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding image files: %v\n", err)
		os.Exit(1)
	}

	if len(imageFiles) == 0 {
		fmt.Printf("対象の画像ファイルが見つかりませんでした: %s\n", *dirPath)
		os.Exit(0)
	}

	if config.Verbose {
		log.Printf("Found %d image files\n", len(imageFiles))
	}

	// fromFilename または fromScreenshot モードの場合、ファイル名から日時を抽出してファイルごとに設定
	if *fromFilename || *fromScreenshot {
		var processedCount, errorCount int
		var err error
		
		if *fromFilename {
			processedCount, errorCount, err = processFilesFromFilename(service, imageFiles, config)
		} else if *fromScreenshot {
			processedCount, errorCount, err = processFilesFromScreenshot(service, imageFiles, config)
		}
		
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
			os.Exit(1)
		}

		// 結果を表示
		fmt.Printf("\n処理完了: %d個のファイルを処理しました", processedCount)
		if errorCount > 0 {
			fmt.Printf(" (%d個のエラー)", errorCount)
		}
		fmt.Println()
	} else {
		// 通常のモード（全ファイルに同じ日時を設定）
		processedCount, errorCount, err := service.ModifyExifData(imageFiles, config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error modifying EXIF data: %v\n", err)
			os.Exit(1)
		}

		// 結果を表示
		fmt.Printf("\n処理完了: %d個のファイルを処理しました", processedCount)
		if errorCount > 0 {
			fmt.Printf(" (%d個のエラー)", errorCount)
		}
		fmt.Println()
	}
}

// 日時文字列をtime.Timeに変換
func parseDateTime(dateTimeStr string) (time.Time, error) {
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

	return time.ParseInLocation("20060102150405", dateTimeStr, time.Local)
}

// 日時の各要素をバリデーション
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

// 文字列を整数に変換（エラーハンドリングなし、事前に数値チェック済み）
func parseInt(s string) int {
	result := 0
	for _, char := range s {
		result = result*10 + int(char-'0')
	}
	return result
}

// 指定された年月の最大日数を取得
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

// うるう年判定
func isLeapYear(year int) bool {
	// 4で割り切れる年はうるう年
	// ただし100で割り切れる年は平年
	// ただし400で割り切れる年はうるう年
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

// ディレクトリの存在と権限をバリデーション
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

// ファイル拡張子をバリデーション
func validateExtension(ext string) error {
	// 空文字列チェック
	if ext == "" {
		return fmt.Errorf("拡張子が空です")
	}

	// ドットで始まるかチェック
	if !filepath.HasPrefix(ext, ".") {
		return fmt.Errorf("拡張子はドット（.）で始まる必要があります: %s", ext)
	}

	// サポートされている拡張子のリスト
	supportedExtensions := []string{".jpg", ".jpeg", ".tiff", ".tif", ".png", ".webp", ".mp4", ".webm"}
	
	// 大文字小文字を無視して比較
	extLower := filepath.ToSlash(ext) // パスを正規化
	extLower = filepath.Ext(extLower + "dummy") // 拡張子として認識させる
	if extLower == "" {
		extLower = ext
	}
	extLower = filepath.ToSlash(extLower)
	
	// より確実な方法で小文字に変換
	extLower = ""
	for _, char := range ext {
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

// 実行情報を表示
func printExecutionInfo(config *usecases.Config) {
	fmt.Printf("ディレクトリ: %s\n", config.FolderPath)
	if config.FromFilename {
		fmt.Println("モード: ファイル名から日時を取得")
	} else if config.FromScreenshot {
		fmt.Println("モード: スクリーンショットファイル名から日時を取得")
	} else {
		fmt.Printf("設定する日時: %s\n", config.DateTime.Format("2006-01-02 15:04:05"))
	}
	if config.Extension != "" {
		fmt.Printf("対象拡張子: %s\n", config.Extension)
	}
	fmt.Printf("再帰処理: %t\n", config.Recursive)
	fmt.Printf("ドライラン: %t\n", config.DryRun)
	fmt.Printf("詳細モード: %t\n", config.Verbose)
	fmt.Println()
}

// ファイル名から日時を抽出してファイルごとに処理
func processFilesFromFilename(service *usecases.ExifModifierService, imageFiles []string, config *usecases.Config) (int, int, error) {
	processedCount := 0
	errorCount := 0
	
	// yyyyMMddhhmmss形式の正規表現（ファイル名の先頭に14桁の数字がある場合）
	dateTimeRegex := regexp.MustCompile(`^(\d{14})`)
	
	for _, filePath := range imageFiles {
		fileName := filepath.Base(filePath)
		fileName = removeExtension(fileName)
		
		// ファイル名から日時を抽出
		matches := dateTimeRegex.FindStringSubmatch(fileName)
		if len(matches) < 2 {
			if config.Verbose {
				log.Printf("ファイル名が日時形式ではありません（スキップ）: %s", fileName)
			}
			continue
		}
		
		// 日時をパース
		dateTimeStr := matches[1]
		fileDateTime, err := parseDateTime(dateTimeStr)
		if err != nil {
			if config.Verbose {
				log.Printf("日時のパースに失敗しました（スキップ）: %s - %v", fileName, err)
			}
			errorCount++
			continue
		}
		
		if config.Verbose {
			log.Printf("ファイル: %s -> 日時: %s", fileName, fileDateTime.Format("2006-01-02 15:04:05"))
		}
		
		// 一時的にconfigの日時を変更
		tempConfig := *config
		tempConfig.DateTime = fileDateTime
		
		// ドライランの場合は処理をスキップ
		if config.DryRun {
			fmt.Printf("[DRY RUN] %s -> %s\n", filePath, fileDateTime.Format("2006-01-02 15:04:05"))
			processedCount++
			continue
		}
		
		// 単一ファイルの処理
		_, fileErrorCount, err := service.ModifyExifData([]string{filePath}, &tempConfig)
		if err != nil {
			if config.Verbose {
				log.Printf("ファイル処理中にエラーが発生しました: %s - %v", filePath, err)
			}
			errorCount++
			continue
		}
		
		if fileErrorCount > 0 {
			errorCount += fileErrorCount
		} else {
			processedCount++
		}
	}
	
	return processedCount, errorCount, nil
}

// スクリーンショットファイル名から日時を抽出してファイルごとに処理
func processFilesFromScreenshot(service *usecases.ExifModifierService, imageFiles []string, config *usecases.Config) (int, int, error) {
	processedCount := 0
	errorCount := 0
	
	// Screenshot_yyyyMMdd-hhmmss形式の正規表現
	screenshotRegex := regexp.MustCompile(`^Screenshot_(\d{8})-(\d{6})`)
	
	for _, filePath := range imageFiles {
		fileName := filepath.Base(filePath)
		fileName = removeExtension(fileName)
		
		// ファイル名から日時を抽出
		matches := screenshotRegex.FindStringSubmatch(fileName)
		if len(matches) < 3 {
			if config.Verbose {
				log.Printf("ファイル名がスクリーンショット形式ではありません（スキップ）: %s", fileName)
			}
			continue
		}
		
		// 日付と時刻を結合してyyyyMMddhhmmss形式にする
		dateStr := matches[1]  // yyyyMMdd
		timeStr := matches[2]  // hhmmss
		dateTimeStr := dateStr + timeStr
		
		// 日時をパース
		fileDateTime, err := parseDateTime(dateTimeStr)
		if err != nil {
			if config.Verbose {
				log.Printf("日時のパースに失敗しました（スキップ）: %s - %v", fileName, err)
			}
			errorCount++
			continue
		}
		
		if config.Verbose {
			log.Printf("ファイル: %s -> 日時: %s", fileName, fileDateTime.Format("2006-01-02 15:04:05"))
		}
		
		// 一時的にconfigの日時を変更
		tempConfig := *config
		tempConfig.DateTime = fileDateTime
		
		// ドライランの場合は処理をスキップ
		if config.DryRun {
			fmt.Printf("[DRY RUN] %s -> %s\n", filePath, fileDateTime.Format("2006-01-02 15:04:05"))
			processedCount++
			continue
		}
		
		// 単一ファイルの処理
		_, fileErrorCount, err := service.ModifyExifData([]string{filePath}, &tempConfig)
		if err != nil {
			if config.Verbose {
				log.Printf("ファイル処理中にエラーが発生しました: %s - %v", filePath, err)
			}
			errorCount++
			continue
		}
		
		if fileErrorCount > 0 {
			errorCount += fileErrorCount
		} else {
			processedCount++
		}
	}
	
	return processedCount, errorCount, nil
}

// ファイル名から拡張子を除去
func removeExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return fileName[:len(fileName)-len(ext)]
	}
	return fileName
}
