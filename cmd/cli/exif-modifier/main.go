package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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
	extensionsStr := flag.String("ext", "jpg,jpeg,tiff,tif,png", "対象の画像拡張子（カンマ区切り）")
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
	ext := strings.ToLower(filepath.Ext(filePath))

	// JPEGファイルの場合
	if ext == ".jpg" || ext == ".jpeg" {
		return extractJpegExif(filePath, config)
	}

	// PNGファイルの場合
	if ext == ".png" {
		return extractPngMetadata(filePath, config)
	}

	// その他の画像形式（TIFF等）の場合
	return extractGenericExif(filePath, config)
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

func extractPngMetadata(filePath string, config *Config) (ExifData, error) {
	data := ExifData{
		FilePath:   filePath,
		Properties: make(map[string]string),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// ファイル情報を追加
	stat, err := file.Stat()
	if err == nil {
		data.Properties["File Size"] = formatFileSize(stat.Size())
		data.Properties["File Modification Date/Time"] = stat.ModTime().Format("2006:01:02 15:04:05-07:00")
		data.Properties["File Name"] = filepath.Base(filePath)
		data.Properties["Directory"] = filepath.Dir(filePath)
		data.Properties["File Type"] = "PNG"
		data.Properties["File Type Extension"] = "png"
		data.Properties["MIME Type"] = "image/png"
	}

	// 画像メタデータを読み取り
	img, err := png.Decode(file)
	if err != nil {
		return data, err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	data.Properties["Image Width"] = strconv.Itoa(width)
	data.Properties["Image Height"] = strconv.Itoa(height)
	data.Properties["Image Size"] = fmt.Sprintf("%dx%d", width, height)
	data.Properties["Megapixels"] = fmt.Sprintf("%.1f", float64(width*height)/1000000.0)

	// PNG固有の情報を読み取り（再度ファイルを開く）
	file.Seek(0, 0)
	pngInfo, err := readPngChunks(file)
	if err == nil {
		for key, value := range pngInfo {
			if config.Properties == nil || contains(config.Properties, key) {
				data.Properties[key] = value
			}
		}
	}

	return data, nil
}

func extractGenericExif(filePath string, config *Config) (ExifData, error) {
	data := ExifData{
		FilePath:   filePath,
		Properties: make(map[string]string),
	}

	file, err := os.Open(filePath)
	if err != nil {
		return data, err
	}
	defer file.Close()

	// ファイル情報を追加
	stat, err := file.Stat()
	if err == nil {
		data.Properties["File Size"] = formatFileSize(stat.Size())
		data.Properties["File Modification Date/Time"] = stat.ModTime().Format("2006:01:02 15:04:05-07:00")
		data.Properties["File Name"] = filepath.Base(filePath)
		data.Properties["Directory"] = filepath.Dir(filePath)
		ext := strings.ToUpper(strings.TrimPrefix(filepath.Ext(filePath), "."))
		data.Properties["File Type"] = ext
		data.Properties["File Type Extension"] = strings.ToLower(ext)
	}

	// EXIFデータを読み取り試行
	fileData := make([]byte, stat.Size())
	_, err = file.Read(fileData)
	if err != nil {
		return data, err
	}

	rawExif, err := exif.SearchAndExtractExif(fileData)
	if err != nil {
		return data, nil // EXIFがない場合はファイル情報のみ返す
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

func readPngChunks(file *os.File) (map[string]string, error) {
	info := make(map[string]string)

	// PNG signature check
	signature := make([]byte, 8)
	_, err := file.Read(signature)
	if err != nil {
		return info, err
	}

	for {
		// Read chunk length
		lengthBytes := make([]byte, 4)
		n, err := file.Read(lengthBytes)
		if err != nil || n != 4 {
			break
		}
		length := binary.BigEndian.Uint32(lengthBytes)

		// Read chunk type
		chunkType := make([]byte, 4)
		_, err = file.Read(chunkType)
		if err != nil {
			break
		}

		chunkTypeName := string(chunkType)

		// Read chunk data
		chunkData := make([]byte, length)
		if length > 0 {
			_, err = file.Read(chunkData)
			if err != nil {
				break
			}
		}

		// Skip CRC
		file.Seek(4, 1)

		// Process specific chunks
		switch chunkTypeName {
		case "IHDR":
			if len(chunkData) >= 13 {
				width := binary.BigEndian.Uint32(chunkData[0:4])
				height := binary.BigEndian.Uint32(chunkData[4:8])
				bitDepth := chunkData[8]
				colorType := chunkData[9]
				compression := chunkData[10]
				filter := chunkData[11]
				interlace := chunkData[12]

				info["Bit Depth"] = strconv.Itoa(int(bitDepth))
				info["Color Type"] = getColorTypeName(colorType)
				info["Compression"] = getCompressionName(compression)
				info["Filter"] = getFilterName(filter)
				info["Interlace"] = getInterlaceName(interlace)

				// 既に設定されていなければ幅と高さも設定
				if _, exists := info["Image Width"]; !exists {
					info["Image Width"] = strconv.Itoa(int(width))
					info["Image Height"] = strconv.Itoa(int(height))
				}
			}
		case "pHYs":
			if len(chunkData) >= 9 {
				pixelsPerUnitX := binary.BigEndian.Uint32(chunkData[0:4])
				pixelsPerUnitY := binary.BigEndian.Uint32(chunkData[4:8])
				unitSpecifier := chunkData[8]

				info["Pixels Per Unit X"] = strconv.Itoa(int(pixelsPerUnitX))
				info["Pixels Per Unit Y"] = strconv.Itoa(int(pixelsPerUnitY))
				if unitSpecifier == 1 {
					info["Pixel Units"] = "meters"
				} else {
					info["Pixel Units"] = "unknown"
				}
			}
		case "gAMA":
			if len(chunkData) >= 4 {
				gamma := binary.BigEndian.Uint32(chunkData[0:4])
				gammaValue := float64(gamma) / 100000.0
				info["Gamma"] = fmt.Sprintf("%.1f", gammaValue)
			}
		case "sRGB":
			if len(chunkData) >= 1 {
				renderingIntent := chunkData[0]
				info["SRGB Rendering"] = getSRGBRenderingIntent(renderingIntent)
			}
		case "IEND":
			return info, nil
		}
	}

	return info, nil
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.0f %cB", float64(size)/float64(div), "kMGTPE"[exp])
}

func getColorTypeName(colorType byte) string {
	switch colorType {
	case 0:
		return "Grayscale"
	case 2:
		return "RGB"
	case 3:
		return "Palette"
	case 4:
		return "Grayscale with Alpha"
	case 6:
		return "RGB with Alpha"
	default:
		return "Unknown"
	}
}

func getCompressionName(compression byte) string {
	if compression == 0 {
		return "Deflate/Inflate"
	}
	return "Unknown"
}

func getFilterName(filter byte) string {
	if filter == 0 {
		return "Adaptive"
	}
	return "Unknown"
}

func getInterlaceName(interlace byte) string {
	switch interlace {
	case 0:
		return "Noninterlaced"
	case 1:
		return "Adam7 Interlaced"
	default:
		return "Unknown"
	}
}

func getSRGBRenderingIntent(intent byte) string {
	switch intent {
	case 0:
		return "Perceptual"
	case 1:
		return "Relative Colorimetric"
	case 2:
		return "Saturation"
	case 3:
		return "Absolute Colorimetric"
	default:
		return "Unknown"
	}
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
