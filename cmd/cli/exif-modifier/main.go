package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	usecases "github.com/landmaster135/devbox/internal/exif_modifier/usecases"
)

const version = "1.0.0"

// コマンドライン引数
var (
	dirPath        = flag.String("dir", ".", "画像ファイルがあるディレクトリのパス")
	dateTime       = flag.String("datetime", "", "設定する日時 (yyyyMMddhhmmss形式)")
	fromFilename   = flag.Bool("from-filename", false, "ファイル名から日時を取得してExifに設定する (ファイル名がyyyyMMddhhmmss形式の場合)")
	fromScreenshot = flag.Bool("from-screenshot", false, "スクリーンショットファイル名から日時を取得してExifに設定する (Screenshot_yyyyMMdd-hhmmss形式の場合)")
	extension      = flag.String("ext", "jpg", "対象とする拡張子 (例: jpg, jpeg, png, webp, mp4)")
	recursive      = flag.Bool("recursive", false, "サブフォルダも再帰的に処理する")
	dryRun         = flag.Bool("dry-run", false, "実際には変更せず、処理対象ファイルのみ表示")
	verbose        = flag.Bool("verbose", false, "詳細な出力を表示")
	showHelp       = flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion    = flag.Bool("version", false, "バージョン情報を表示")
	workerCount    = flag.Int("workers", runtime.NumCPU(), "並行処理のワーカー数")
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
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --ext jpg\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --recursive\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --dry-run\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --datetime 20240315143000 --verbose\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --from-filename --ext jpg\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./photos --from-filename --recursive --dry-run\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --dir ./screenshots --from-screenshot --ext png\n")
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

	// ExifModifierServiceを作成
	service := usecases.NewExifModifierService()

	var targetTime time.Time
	var err error

	// 日時のパース
	if *dateTime != "" {
		targetTime, err = service.ParseDateTime(*dateTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 日時の形式が正しくありません: %v\n", err)
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
		WorkerCount:    *workerCount,
	}

	// 実行情報を表示
	printExecutionInfo(config)

	// fromFilename または fromScreenshot モードの場合、ファイル名から日時を抽出してファイルごとに設定
	if *fromFilename || *fromScreenshot {
		var err error

		if *fromFilename {
			err = service.ProcessFilesFromFilename(*dirPath, *extension, *recursive, *dryRun, *verbose, true)
		} else if *fromScreenshot {
			err = service.ProcessFilesFromScreenshot(*dirPath, *extension, *recursive, *dryRun, *verbose, true)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing files: %v\n", err)
			os.Exit(1)
		}

		// 結果を表示
		fmt.Printf("\n処理完了\n")
		return
	}

	// 通常のモード（全ファイルに同じ日時を設定）の場合のみ画像ファイルを検索
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
	fmt.Printf("ワーカー数: %d\n", config.WorkerCount)
	fmt.Println()
}
