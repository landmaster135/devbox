package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

const version = "1.0.0"

type Config struct {
	Directory   string
	Extensions  []string
	Properties  []string
	ShowHelp    bool
	ShowVersion bool
	Verbose     bool
	Recursive   bool
}

type ExifData struct {
	FilePath   string
	Properties map[string]string
}

func main() {
	// フラグを定義
	directory := flag.String("dir", ".", "画像ファイルを検索するディレクトリ")
	extensionsStr := flag.String("ext", "jpg,jpeg,tiff,tif", "対象の画像拡張子（カンマ区切り）")
	propertiesStr := flag.String("props", "", "表示するExifプロパティ（カンマ区切り、空の場合は全て表示）")
	showHelp := flag.Bool("help", false, "ヘルプメッセージを表示")
	showVersion := flag.Bool("version", false, "バージョン情報を表示")
	verbose := flag.Bool("v", false, "詳細なログを出力")
	recursive := flag.Bool("r", false, "サブディレクトリも再帰的に検索")

	// カスタムUsage関数を設定
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "exif-modifier - EXIF Property Viewer Tool\n\n")
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier [OPTIONS]\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier -dir ./photos\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier -dir ./photos -ext jpg,png -props DateTime,Camera\n")
		fmt.Fprintf(os.Stderr, "    exif-modifier -dir ./photos -r -v\n")
	}

	// フラグを解析
	flag.Parse()

	// ログ設定
	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	} else {
		log.SetFlags(log.Ldate | log.Ltime)
	}

	config := &Config{
		Directory:   *directory,
		Extensions:  strings.Split(strings.ToLower(*extensionsStr), ","),
		ShowHelp:    *showHelp,
		ShowVersion: *showVersion,
		Verbose:     *verbose,
		Recursive:   *recursive,
	}

	if *propertiesStr != "" {
		config.Properties = strings.Split(*propertiesStr, ",")
	}

	// ヘルプまたはバージョン表示
	if config.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

	if config.ShowVersion {
		fmt.Printf("exif-modifier version %s\n", version)
		os.Exit(0)
	}

	// ディレクトリの存在確認
	if _, err := os.Stat(config.Directory); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist\n", config.Directory)
		os.Exit(1)
	}

	// 画像ファイルを検索
	imageFiles, err := findImageFiles(config)
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
	exifDataList, err := extractExifData(imageFiles, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting EXIF data: %v\n", err)
		os.Exit(1)
	}

	// 結果をテーブル形式で表示
	displayExifTable(exifDataList, config)
}

func findImageFiles(config *Config) ([]string, error) {
	var imageFiles []string

	err := filepath.Walk(config.Directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 再帰フラグが設定されていない場合、サブディレクトリをスキップ
		if !config.Recursive && info.IsDir() && path != config.Directory {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			ext = strings.TrimPrefix(ext, ".")

			for _, targetExt := range config.Extensions {
				if ext == strings.TrimSpace(targetExt) {
					imageFiles = append(imageFiles, path)
					break
				}
			}
		}

		return nil
	})

	return imageFiles, err
}

func extractExifData(imageFiles []string, config *Config) ([]ExifData, error) {
	var exifDataList []ExifData

	for _, filePath := range imageFiles {
		if config.Verbose {
			log.Printf("Processing: %s\n", filePath)
		}

		data, err := extractSingleFileExif(filePath, config)
		if err != nil {
			if config.Verbose {
				log.Printf("Warning: Could not extract EXIF from %s: %v\n", filePath, err)
			}
			// エラーがあってもスキップして続行
			continue
		}

		exifDataList = append(exifDataList, data)
	}

	return exifDataList, nil
}

func extractSingleFileExif(filePath string, config *Config) (ExifData, error) {
	data := ExifData{
		FilePath:   filePath,
		Properties: make(map[string]string),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// JPEGファイルの場合の処理
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".jpg" || ext == ".jpeg" {
		return extractJpegExif(filePath, config)
	}

	// その他の画像形式の場合（基本的なEXIF読み取り）
	stat, err := file.Stat()
	if err != nil {
		return data, err
	}

	fileData := make([]byte, stat.Size())
	_, err = file.Read(fileData)
	if err != nil {
		return data, err
	}

	rawExif, err := exif.SearchAndExtractExif(fileData)
	if err != nil {
		return data, err
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return data, err
	}

	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		return data, err
	}

	// EXIFデータを収集
	children := index.RootIfd.Children()
	for _, ifd := range children {
		entries := ifd.Entries()
		for _, entry := range entries {
			tagName := entry.TagName()
			value, err := entry.FormatFirst()
			if err != nil {
				continue
			}

			if config.Properties == nil || contains(config.Properties, tagName) {
				data.Properties[tagName] = fmt.Sprintf("%v", value)
			}
		}
	}

	return data, nil
}

func extractJpegExif(filePath string, config *Config) (ExifData, error) {
	data := ExifData{
		FilePath:   filePath,
		Properties: make(map[string]string),
	}

	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseFile(filePath)
	if err != nil {
		return data, err
	}

	sl := intfc.(*jpegstructure.SegmentList)
	rootIfd, _, err := sl.Exif()
	if err != nil {
		return data, err
	}

	// EXIFデータを収集
	err = rootIfd.EnumerateTagsRecursively(func(ifd *exif.Ifd, ite *exif.IfdTagEntry) error {
		tagName := ite.TagName()
		value, err := ite.FormatFirst()
		if err != nil {
			return nil // エラーがあってもスキップ
		}

		if config.Properties == nil || contains(config.Properties, tagName) {
			data.Properties[tagName] = fmt.Sprintf("%v", value)
		}

		return nil
	})

	return data, err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(strings.TrimSpace(s), item) {
			return true
		}
	}
	return false
}

func displayExifTable(exifDataList []ExifData, config *Config) {
	if len(exifDataList) == 0 {
		fmt.Println("No EXIF data found.")
		return
	}

	// 全てのプロパティキーを収集
	propertyKeys := make(map[string]bool)
	for _, data := range exifDataList {
		for key := range data.Properties {
			propertyKeys[key] = true
		}
	}

	// プロパティキーをソート
	var sortedKeys []string
	for key := range propertyKeys {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// フィルタリング（指定されたプロパティがある場合）
	if config.Properties != nil {
		var filteredKeys []string
		for _, key := range sortedKeys {
			if contains(config.Properties, key) {
				filteredKeys = append(filteredKeys, key)
			}
		}
		sortedKeys = filteredKeys
	}

	// テーブルライターを作成
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// ヘッダーを出力
	fmt.Fprint(w, "File Path")
	for _, key := range sortedKeys {
		fmt.Fprintf(w, "\t%s", key)
	}
	fmt.Fprintln(w)

	// 区切り線
	fmt.Fprint(w, strings.Repeat("-", 50))
	for range sortedKeys {
		fmt.Fprint(w, "\t", strings.Repeat("-", 20))
	}
	fmt.Fprintln(w)

	// データ行を出力
	for _, data := range exifDataList {
		// ファイルパスを相対パスで表示
		relPath, err := filepath.Rel(".", data.FilePath)
		if err != nil {
			relPath = data.FilePath
		}
		fmt.Fprint(w, relPath)

		for _, key := range sortedKeys {
			value := data.Properties[key]
			if value == "" {
				value = "-"
			}
			// 長すぎる値は短縮
			if len(value) > 30 {
				value = value[:27] + "..."
			}
			fmt.Fprintf(w, "\t%s", value)
		}
		fmt.Fprintln(w)
	}

	// テーブルをフラッシュ
	w.Flush()

	// サマリーを表示
	fmt.Printf("\nSummary: %d files processed, %d unique EXIF properties found\n",
		len(exifDataList), len(sortedKeys))
}
