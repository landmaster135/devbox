package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/landmaster135/devbox/internal/independencies/exif_modifier/usecases"
)

const version = "1.0.0"

// コマンドライン引数
var (
	folderPath  = flag.String("folder", ".", "画像ファイルがあるフォルダのパス")
	dateTime    = flag.String("datetime", "", "設定する日時 (yyyyMMddhhmmss形式)")
	extension   = flag.String("ext", "", "対象とする拡張子 (例: .jpg, .jpeg, .tiff)")
	recursive   = flag.Bool("recursive", false, "サブフォルダも再帰的に処理する")
	dryRun      = flag.Bool("dry-run", false, "実際には変更せず、処理対象ファイルのみ表示")
	verbose     = flag.Bool("verbose", false, "詳細な出力を表示")
	showHelp    = flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion = flag.Bool("version", false, "バージョン情報を表示")
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
		fmt.Fprintf(os.Stderr, "    exif-modifier --folder ./photos --datetime 20240315143000 --ext .jpg\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --folder ./photos --datetime 20240315143000 --recursive\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --folder ./photos --datetime 20240315143000 --dry-run\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier --folder ./photos --datetime 20240315143000 --verbose\n")
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
	if *dateTime == "" {
		fmt.Fprintf(os.Stderr, "Error: --datetime パラメータは必須です (yyyyMMddhhmmss形式)\n")
		flag.Usage()
		os.Exit(1)
	}

	// 日時のパース
	targetTime, err := parseDateTime(*dateTime)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 日時の形式が正しくありません: %v\n", err)
		os.Exit(1)
	}

	// フォルダの存在確認
	if _, err := os.Stat(*folderPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: フォルダが存在しません: %s\n", *folderPath)
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
		FolderPath: *folderPath,
		DateTime:   targetTime,
		Extension:  *extension,
		Recursive:  *recursive,
		DryRun:     *dryRun,
		Verbose:    *verbose,
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
		fmt.Printf("対象の画像ファイルが見つかりませんでした: %s\n", *folderPath)
		os.Exit(0)
	}

	if config.Verbose {
		log.Printf("Found %d image files\n", len(imageFiles))
	}

	// Exif情報を更新
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

// 日時文字列をtime.Timeに変換
func parseDateTime(dateTimeStr string) (time.Time, error) {
	if len(dateTimeStr) != 14 {
		return time.Time{}, fmt.Errorf("日時は14文字である必要があります (yyyyMMddhhmmss)")
	}

	return time.ParseInLocation("20060102150405", dateTimeStr, time.Local)
}

// 実行情報を表示
func printExecutionInfo(config *usecases.Config) {
	fmt.Printf("フォルダ: %s\n", config.FolderPath)
	fmt.Printf("設定する日時: %s\n", config.DateTime.Format("2006-01-02 15:04:05"))
	if config.Extension != "" {
		fmt.Printf("対象拡張子: %s\n", config.Extension)
	}
	fmt.Printf("再帰処理: %t\n", config.Recursive)
	fmt.Printf("ドライラン: %t\n", config.DryRun)
	fmt.Printf("詳細モード: %t\n", config.Verbose)
	fmt.Println()
}
