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

	// ディレクトリの存在確認
	if _, err := os.Stat(*dirPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: ディレクトリが存在しません: %s\n", *dirPath)
		os.Exit(1)
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

	return time.ParseInLocation("20060102150405", dateTimeStr, time.Local)
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
