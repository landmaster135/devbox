package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/exif_viewer/usecases"
)

const version = "1.0.0"

func main() {
	// フラグを定義
	directory := flag.String("dir", ".", "画像ファイルを検索するディレクトリ")
	extensionsStr := flag.String("ext", "jpg,jpeg,tiff,tif,png", "対象の画像拡張子（カンマ区切り）")
	propertiesStr := flag.String("props", "", "表示するExifプロパティ（カンマ区切り、空の場合は全て表示）")
	maxProps := flag.Int("max", 0, "表示するプロパティの最大数（0の場合は制限なし）")
	showHelp := flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion := flag.Bool("version", false, "バージョン情報を表示")
	verbose := flag.Bool("v", false, "詳細なログを出力")
	recursive := flag.Bool("r", false, "サブディレクトリも再帰的に検索")

	// カスタムUsage関数を設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "exif-viewer - EXIF Property Viewer Tool\n\n")
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "    exif-viewer [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    exif-viewer -dir ./photos\n")
		fmt.Fprintf(os.Stderr, "    exif-viewer -dir ./photos -ext jpg,png -props DateTime,Camera\n")
		fmt.Fprintf(os.Stderr, "    exif-viewer -dir ./photos -max 5\n")
		fmt.Fprintf(os.Stderr, "    exif-viewer -dir ./photos -r -v\n")
	}

	// フラグを解析
	flag.Parse()

	// ログ設定
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	config := &usecases.Config{
		Directory:  *directory,
		Extensions: strings.Split(strings.ToLower(*extensionsStr), ","),
		MaxProps:   *maxProps,
		Verbose:    *verbose,
		Recursive:  *recursive,
	}

	if *propertiesStr != "" {
		config.Properties = strings.Split(*propertiesStr, ",")
	}

	// ヘルプまたはバージョン表示
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("exif-viewer version %s\n", version)
		os.Exit(0)
	}

	// ディレクトリの存在確認
	if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist\n", config.Directory)
		os.Exit(1)
	}

	// ExifViewerServiceを作成
	service := usecases.NewExifViewerService()

	// 画像ファイルを検索
	imageFiles, err := service.FindImageFiles(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding image files: %v\n", err)
		os.Exit(1)
	}

	if len(imageFiles) == 0 {
		fmt.Printf("No image files found in directory: %s\n", config.Directory)
		os.Exit(0)
	}

	if config.Verbose {
		log.Printf("Found %d image files\n", len(imageFiles))
	}

	// Exif情報を抽出
	exifDataList, err := service.ExtractExifData(imageFiles, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting EXIF data: %v\n", err)
		os.Exit(1)
	}

	// 結果をテーブル形式で表示
	output := service.FormatExifTable(exifDataList, config)
	fmt.Print(output)
}
