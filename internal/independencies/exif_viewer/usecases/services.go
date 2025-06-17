package usecases

import (
	"encoding/binary"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

// Config はEXIF表示の設定を保持します
type Config struct {
	Directory      string
	Extensions     []string
	Properties     []string
	MaxProps       int
	Verbose        bool
	Recursive      bool
	ShowProperties bool // プロパティ一覧を表示するフラグ
	ShowDataTypes  bool // データ型を表示するフラグ
}

// ExifData は単一ファイルのEXIF情報を保持します
type ExifData struct {
	FilePath   string
	Properties map[string]string
}

// PropertyInfo はプロパティの詳細情報を保持します
type PropertyInfo struct {
	Name     string
	DataType string
	Count    int // このプロパティが見つかったファイル数
	Examples []string // 値の例（最大3つ）
}

// ExifViewerService はEXIF表示サービスです
type ExifViewerService struct{}

// NewExifViewerService は新しいExifViewerServiceを作成します
func NewExifViewerService() *ExifViewerService {
	return &ExifViewerService{}
}

// ensureUTF8String は文字列がUTF-8として有効かチェックし、無効な場合は修正する
func (s *ExifViewerService) ensureUTF8String(str string) string {
	if utf8.ValidString(str) {
		return str
	}
	// 無効なUTF-8文字を置換
	return strings.ToValidUTF8(str, "�")
}

// FindImageFiles は指定された設定に基づいて画像ファイルを検索します
func (s *ExifViewerService) FindImageFiles(config *Config) ([]string, error) {
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

// ExtractExifData は複数のファイルからEXIF情報を抽出します
func (s *ExifViewerService) ExtractExifData(imageFiles []string, config *Config) ([]ExifData, error) {
	var exifDataList []ExifData

	for _, filePath := range imageFiles {
		data, err := s.ExtractSingleFileExif(filePath, config)
		if err != nil {
			// エラーがあってもスキップして続行
			continue
		}

		exifDataList = append(exifDataList, data)
	}

	return exifDataList, nil
}

// ExtractSingleFileExif は単一ファイルからEXIF情報を抽出します
func (s *ExifViewerService) ExtractSingleFileExif(filePath string, config *Config) (ExifData, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// JPEGファイルの場合
	if ext == ".jpg" || ext == ".jpeg" {
		return s.extractJpegExif(filePath, config)
	}

	// PNGファイルの場合
	if ext == ".png" {
		return s.extractPngMetadata(filePath, config)
	}

	// その他の画像形式（TIFF等）の場合
	return s.extractGenericExif(filePath, config)
}

// extractJpegExif はJPEGファイルからEXIF情報を抽出します
func (s *ExifViewerService) extractJpegExif(filePath string, config *Config) (ExifData, error) {
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

		if config.Properties == nil || s.contains(config.Properties, tagName) {
			data.Properties[tagName] = s.ensureUTF8String(fmt.Sprintf("%v", value))
		}

		return nil
	})

	return data, err
}

// extractPngMetadata はPNGファイルからメタデータを抽出します
func (s *ExifViewerService) extractPngMetadata(filePath string, config *Config) (ExifData, error) {
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
		data.Properties["File Size"] = s.formatFileSize(stat.Size())
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
	pngInfo, err := s.readPngChunks(file)
	if err == nil {
		for key, value := range pngInfo {
			if config.Properties == nil || s.contains(config.Properties, key) {
				data.Properties[key] = value
			}
		}
	}

	return data, nil
}

// extractGenericExif は汎用画像ファイルからEXIF情報を抽出します
func (s *ExifViewerService) extractGenericExif(filePath string, config *Config) (ExifData, error) {
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
		data.Properties["File Size"] = s.formatFileSize(stat.Size())
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

			if config.Properties == nil || s.contains(config.Properties, tagName) {
				data.Properties[tagName] = s.ensureUTF8String(fmt.Sprintf("%v", value))
			}
		}
	}

	return data, nil
}

// readPngChunks はPNGチャンクを読み取ります
func (s *ExifViewerService) readPngChunks(file *os.File) (map[string]string, error) {
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
				info["Color Type"] = s.getColorTypeName(colorType)
				info["Compression"] = s.getCompressionName(compression)
				info["Filter"] = s.getFilterName(filter)
				info["Interlace"] = s.getInterlaceName(interlace)

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
				info["SRGB Rendering"] = s.getSRGBRenderingIntent(renderingIntent)
			}
		case "IEND":
			return info, nil
		}
	}

	return info, nil
}

// FormatExifTable はEXIF情報をテーブル形式で表示します
func (s *ExifViewerService) FormatExifTable(exifDataList []ExifData, config *Config) string {
	if len(exifDataList) == 0 {
		return "No EXIF data found."
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
			if s.contains(config.Properties, key) {
				filteredKeys = append(filteredKeys, key)
			}
		}
		sortedKeys = filteredKeys
	}

	// プロパティ数制限の適用
	if config.MaxProps > 0 && len(sortedKeys) > config.MaxProps {
		sortedKeys = sortedKeys[:config.MaxProps]
	}

	// 文字列ビルダーでテーブルを構築
	var result strings.Builder
	w := tabwriter.NewWriter(&result, 0, 0, 2, ' ', 0)

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

	// サマリーを追加
	totalProperties := len(propertyKeys)
	displayedProperties := len(sortedKeys)

	if config.MaxProps > 0 && totalProperties > displayedProperties {
		fmt.Fprintf(&result, "\nSummary: %d files processed, %d of %d properties displayed (limited by -max %d)\n",
			len(exifDataList), displayedProperties, totalProperties, config.MaxProps)
	} else {
		fmt.Fprintf(&result, "\nSummary: %d files processed, %d unique properties found\n",
			len(exifDataList), displayedProperties)
	}

	return s.ensureUTF8String(result.String())
}

// AnalyzeProperties は全プロパティの詳細情報を分析します
func (s *ExifViewerService) AnalyzeProperties(exifDataList []ExifData) []PropertyInfo {
	propertyStats := make(map[string]*PropertyInfo)

	// 全ファイルのプロパティを分析
	for _, data := range exifDataList {
		for propName, propValue := range data.Properties {
			if info, exists := propertyStats[propName]; exists {
				info.Count++
				// 例を追加（最大3つまで）
				if len(info.Examples) < 3 && !s.containsString(info.Examples, propValue) {
					info.Examples = append(info.Examples, propValue)
				}
			} else {
				propertyStats[propName] = &PropertyInfo{
					Name:     propName,
					DataType: s.inferDataType(propValue),
					Count:    1,
					Examples: []string{propValue},
				}
			}
		}
	}

	// PropertyInfoのスライスに変換してソート
	var result []PropertyInfo
	for _, info := range propertyStats {
		result = append(result, *info)
	}

	sort.Slice(result, func(i, j int) bool {
		// 使用頻度でソート（降順）、同じ場合は名前でソート（昇順）
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Name < result[j].Name
	})

	return result
}

// FormatPropertyList はプロパティ一覧をフォーマットします
func (s *ExifViewerService) FormatPropertyList(propertyInfos []PropertyInfo, totalFiles int) string {
	if len(propertyInfos) == 0 {
		return "No properties found."
	}

	var result strings.Builder
	w := tabwriter.NewWriter(&result, 0, 0, 2, ' ', 0)

	// ヘッダーを出力
	fmt.Fprint(w, "Property Name\tData Type\tFrequency\tUsage %\tExamples\n")
	fmt.Fprint(w, strings.Repeat("-", 50)+"\t"+strings.Repeat("-", 15)+"\t"+strings.Repeat("-", 10)+"\t"+strings.Repeat("-", 8)+"\t"+strings.Repeat("-", 40)+"\n")

	// プロパティ情報を出力
	for _, info := range propertyInfos {
		usagePercent := float64(info.Count) / float64(totalFiles) * 100
		examplesStr := strings.Join(info.Examples, ", ")
		if len(examplesStr) > 40 {
			examplesStr = examplesStr[:37] + "..."
		}

		fmt.Fprintf(w, "%s\t%s\t%d/%d\t%.1f%%\t%s\n",
			info.Name,
			info.DataType,
			info.Count,
			totalFiles,
			usagePercent,
			examplesStr,
		)
	}

	w.Flush()

	fmt.Fprintf(&result, "\nSummary: %d unique properties found across %d files\n",
		len(propertyInfos), totalFiles)

	return s.ensureUTF8String(result.String())
}

// inferDataType は値からデータ型を推測します
func (s *ExifViewerService) inferDataType(value string) string {
	if value == "" {
		return "string"
	}

	// 整数判定
	if _, err := strconv.Atoi(value); err == nil {
		return "integer"
	}

	// 浮動小数点数判定
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return "float"
	}

	// 日付/時刻判定
	if s.isDateTime(value) {
		return "datetime"
	}

	// ファイルサイズ判定
	if s.isFileSize(value) {
		return "filesize"
	}

	// 比率判定（座標判定より先に）
	if s.isRatio(value) {
		return "ratio"
	}

	// 座標判定（より厳密にチェック）
	if s.isCoordinate(value) {
		return "coordinate"
	}

	return "string"
}

// containsString は文字列スライスに特定の文字列が含まれているかチェックします
func (s *ExifViewerService) containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isDateTime は日付・時刻形式かどうかを判定します
func (s *ExifViewerService) isDateTime(value string) bool {
	patterns := []string{
		"2006:01:02 15:04:05",
		"2006:01:02 15:04:05-07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
	}

	for _, pattern := range patterns {
		if len(value) == len(pattern) {
			if strings.Contains(value, ":") && (strings.Contains(value, " ") || strings.Contains(value, "T")) {
				return true
			}
		}
	}
	return false
}

// isFileSize はファイルサイズ形式かどうかを判定します
func (s *ExifViewerService) isFileSize(value string) bool {
	return strings.HasSuffix(value, " B") || strings.HasSuffix(value, "B") ||
		strings.HasSuffix(value, "kB") || strings.HasSuffix(value, "MB") ||
		strings.HasSuffix(value, "GB") || strings.HasSuffix(value, "TB")
}

// isCoordinate は座標形式かどうかを判定します
func (s *ExifViewerService) isCoordinate(value string) bool {
	// 度記号(°)が含まれている場合
	if strings.Contains(value, "°") {
		return true
	}
	
	// 分記号(')が含まれており、かつ数字と組み合わされている場合
	if strings.Contains(value, "'") && strings.ContainsAny(value, "0123456789") {
		return true
	}
	
	// N, S, E, W が文字列の最後にあり、かつ数字が含まれている場合のみ
	if strings.ContainsAny(value, "0123456789") {
		if strings.HasSuffix(value, "N") || strings.HasSuffix(value, "S") ||
		   strings.HasSuffix(value, "E") || strings.HasSuffix(value, "W") {
			return true
		}
	}
	
	return false
}

// isRatio は比率形式かどうかを判定します
func (s *ExifViewerService) isRatio(value string) bool {
	return strings.Contains(value, "/") && len(strings.Split(value, "/")) == 2
}

// ヘルパーメソッド

func (s *ExifViewerService) formatFileSize(size int64) string {
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

func (s *ExifViewerService) getColorTypeName(colorType byte) string {
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

func (s *ExifViewerService) getCompressionName(compression byte) string {
	if compression == 0 {
		return "Deflate/Inflate"
	}
	return "Unknown"
}

func (s *ExifViewerService) getFilterName(filter byte) string {
	if filter == 0 {
		return "Adaptive"
	}
	return "Unknown"
}

func (s *ExifViewerService) getInterlaceName(interlace byte) string {
	switch interlace {
	case 0:
		return "Noninterlaced"
	case 1:
		return "Adam7 Interlaced"
	default:
		return "Unknown"
	}
}

func (s *ExifViewerService) getSRGBRenderingIntent(intent byte) string {
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

func (s *ExifViewerService) contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(strings.TrimSpace(s), item) {
			return true
		}
	}
	return false
}
