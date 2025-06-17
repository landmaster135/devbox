package usecases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	exif "github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

// サポートする画像拡張子
var supportedExtensions = []string{".jpg", ".jpeg", ".tiff", ".tif"}

// Config はEXIF修正の設定を保持します
type Config struct {
	FolderPath string
	DateTime   time.Time
	Extension  string
	Recursive  bool
	DryRun     bool
	Verbose    bool
}

// ExifModifierService はEXIF修正サービスです
type ExifModifierService struct{}

// NewExifModifierService は新しいExifModifierServiceを作成します
func NewExifModifierService() *ExifModifierService {
	return &ExifModifierService{}
}

// ensureUTF8String は文字列がUTF-8として有効かチェックし、無効な場合は修正する
func (s *ExifModifierService) ensureUTF8String(str string) string {
	if utf8.ValidString(str) {
		return str
	}
	return strings.ToValidUTF8(str, "�")
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

	return imageFiles, err
}

// isImageFile は画像ファイルかどうかをチェック
func (s *ExifModifierService) isImageFile(filePath, targetExtension string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 特定の拡張子が指定されている場合
	if targetExtension != "" {
		return strings.ToLower(targetExtension) == ext
	}

	// サポートされている拡張子かチェック
	for _, supportedExt := range supportedExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// ModifyExifData は複数のファイルのEXIF情報を修正します
func (s *ExifModifierService) ModifyExifData(imageFiles []string, config *Config) (int, int, error) {
	processedCount := 0
	errorCount := 0

	for _, filePath := range imageFiles {
		fmt.Printf("処理中: %s\n", filePath)

		if config.DryRun {
			fmt.Printf("  → (ドライラン) File Modification Date/Time を %s に設定\n",
				config.DateTime.Format("2006-01-02 15:04:05"))
			processedCount++
			continue
		}

		err := s.ModifySingleFileExif(filePath, config)
		if err != nil {
			fmt.Printf("  ⚠️  エラー: %v\n", err)
			errorCount++
			if config.Verbose {
				log.Printf("Error processing file %s: %v", filePath, err)
			}
		} else {
			fmt.Printf("  ✅ File Modification Date/Time を %s に設定しました\n",
				config.DateTime.Format("2006-01-02 15:04:05"))
			processedCount++
		}
	}

	return processedCount, errorCount, nil
}

// ModifySingleFileExif は単一ファイルのEXIF情報を修正します
func (s *ExifModifierService) ModifySingleFileExif(filePath string, config *Config) error {
	ext := strings.ToLower(filepath.Ext(filePath))

	// JPEGファイルの場合
	if ext == ".jpg" || ext == ".jpeg" {
		return s.modifyJpegExif(filePath, config)
	}

	// TIFF等の汎用ファイルの場合
	return s.modifyGenericExif(filePath, config)
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

// updateExistingExif は既存のExifデータを更新します
func (s *ExifModifierService) updateExistingExif(filePath string, config *Config, originalData []byte, rootIfd *exif.Ifd, rawExif []byte) error {
	// 現在のgo-exifライブラリでは、Exifデータの直接的な書き込みが複雑なため、
	// ファイルシステムレベルでの時刻更新のみを行います
	if config.Verbose {
		log.Printf("Exifデータが見つかりました。ファイル時刻を更新: %s", filePath)
	}

	return s.updateFileTime(filePath, config.DateTime)
}

// updateFileTime はファイルの更新時刻を変更します
func (s *ExifModifierService) updateFileTime(filePath string, targetTime time.Time) error {
	return os.Chtimes(filePath, targetTime, targetTime)
}
