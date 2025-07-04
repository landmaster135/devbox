package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	usecases "github.com/landmaster135/devbox/internal/image_renamer_with_exif/usecases"
)

const version = "1.0.0"

// コマンドライン引数
var (
	dirPath        = flag.String("dir", ".", "画像ファイルがあるディレクトリのパス")
	extension      = flag.String("ext", "", "対象とする拡張子 (例: jpg, jpeg, png, webp, tiff)")
	recursive      = flag.Bool("recursive", false, "サブフォルダも再帰的に処理する")
	dryRun         = flag.Bool("dry-run", false, "実際には変更せず、処理対象ファイルのみ表示")
	verbose        = flag.Bool("verbose", false, "詳細な出力を表示")
	showHelp       = flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion    = flag.Bool("version", false, "バージョン情報を表示")
	workerCount    = flag.Int("workers", runtime.NumCPU(), "並行処理のワーカー数")
	useFileModTime = flag.Bool("use-file-modtime", false, "ExifのCreateDateではなくファイルの更新時刻を使用")
)

func main() {
	// カスタムUsage関数を設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "image-renamer-with-exif - EXIF CreateDateからファイルをリネームするツール\n\n")
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "DESCRIPTION:\n")
		fmt.Fprintf(os.Stderr, "    画像ファイルのEXIF CreateDateプロパティまたはファイルの更新時刻を使用して、\n")
		fmt.Fprintf(os.Stderr, "    ファイル名を年月日時分秒の形式（YYYYMMDDHHMMSS.拡張子）にリネームします。\n")
		fmt.Fprintf(os.Stderr, "    PowerShellの X0-01_rename_by_exif.ps1 と同等の機能を提供します。\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    # 現在のディレクトリの画像ファイルをリネーム\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif\n\n")
		fmt.Fprintf(os.Stderr, "    # 特定のディレクトリのJPEGファイルのみリネーム\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif --dir ./photos --ext jpg\n\n")
		fmt.Fprintf(os.Stderr, "    # サブディレクトリも含めて再帰的にリネーム\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif --dir ./photos --recursive\n\n")
		fmt.Fprintf(os.Stderr, "    # ドライランで処理対象ファイルを確認\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif --dir ./photos --dry-run\n\n")
		fmt.Fprintf(os.Stderr, "    # ファイルの更新時刻を使用してリネーム\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif --dir ./photos --use-file-modtime\n\n")
		fmt.Fprintf(os.Stderr, "    # 詳細出力でリネーム\n")
		fmt.Fprintf(os.Stderr, "    image-renamer-with-exif --dir ./photos --verbose\n\n")
	}

	flag.Parse()

	// ヘルプまたはバージョン表示
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("image-renamer-with-exif version %s\n", version)
		os.Exit(0)
	}

	// バリデーション
	if err := usecases.ValidateInputOptions(*dirPath, *extension); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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
		Extension:      *extension,
		Recursive:      *recursive,
		DryRun:         *dryRun,
		Verbose:        *verbose,
		WorkerCount:    *workerCount,
		UseFileModTime: *useFileModTime,
	}

	// 実行情報を表示
	printExecutionInfo(config)

	// ImageRenamerServiceを作成
	service := usecases.NewImageRenamerService()

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

	// ファイルをリネーム
	processedCount, errorCount, err := service.RenameImageFiles(imageFiles, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming image files: %v\n", err)
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
	if config.UseFileModTime {
		fmt.Println("モード: ファイルの更新時刻を使用")
	} else {
		fmt.Println("モード: EXIF CreateDateを使用（フォールバック: ファイル更新時刻）")
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
