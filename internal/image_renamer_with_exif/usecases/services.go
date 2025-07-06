package usecases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
)

// サポートする画像拡張子
var supportedExtensions = []string{".jpg", ".jpeg", ".tiff", ".tif", ".png", ".webp", ".mp4", ".webm"}

// Config は画像リネームの設定を保持します
type Config struct {
	FolderPath     string
	Extension      string
	Recursive      bool
	DryRun         bool
	Verbose        bool
	WorkerCount    int  // 並行処理のワーカー数
	UseCreateDate  bool // CreateDateを使用するかどうか
	UseModifyDate  bool // ModifyDateを使用するかどうか
	UseFileModTime bool // ファイルの更新時刻を使用するかどうか
}

// ImageRenamerService は画像リネームサービスです
type ImageRenamerService struct{}

// NewImageRenamerService は新しいImageRenamerServiceを作成します
func NewImageRenamerService() *ImageRenamerService {
	return &ImageRenamerService{}
}

// validateDirectory はディレクトリの存在と権限をバリデーションします
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

// validateExtension はファイル拡張子をバリデーションします
func validateExtension(ext string) error {
	// 空文字列チェック
	if ext == "" {
		return fmt.Errorf("拡張子が空です")
	}

	// 拡張子を正規化（ドットがない場合は追加）
	normalizedExt := ext
	if !strings.HasPrefix(normalizedExt, ".") {
		normalizedExt = "." + normalizedExt
	}

	// より確実な方法で小文字に変換
	extLower := ""
	for _, char := range normalizedExt {
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

// ValidateInputOptions は入力オプションをバリデーションします
func ValidateInputOptions(dirPath, extension string) error {
	// ディレクトリの存在確認とバリデーション
	if err := validateDirectory(dirPath); err != nil {
		return fmt.Errorf("ディレクトリの存在確認とバリデーションに失敗しました: %w", err)
	}

	// 拡張子のバリデーション
	if extension != "" {
		if err := validateExtension(extension); err != nil {
			return fmt.Errorf("拡張子のバリデーションに失敗しました: %w", err)
		}
	}

	return nil
}

// isImageFile は画像ファイルかどうかをチェック
func (s *ImageRenamerService) isImageFile(filePath, targetExtension string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 特定の拡張子が指定されている場合
	if targetExtension != "" {
		// 拡張子を正規化（ドットがない場合は追加）
		normalizedExt := targetExtension
		if !strings.HasPrefix(normalizedExt, ".") {
			normalizedExt = "." + normalizedExt
		}
		return strings.ToLower(normalizedExt) == ext
	}

	// サポートされている拡張子かチェック
	for _, supportedExt := range supportedExtensions {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// FindImageFiles は指定された設定に基づいて画像ファイルを検索します
func (s *ImageRenamerService) FindImageFiles(config *Config) ([]string, error) {
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

	if err != nil {
		return nil, err
	}

	// ファイル名順でソート
	sort.Slice(imageFiles, func(i, j int) bool {
		return filepath.Base(imageFiles[i]) < filepath.Base(imageFiles[j])
	})

	return imageFiles, nil
}

// RenameResult はファイルリネームの結果を表します
type RenameResult struct {
	OriginalPath string
	NewPath      string
	Success      bool
	Error        error
}

// extractCreateDateFromJpeg はJPEGファイルからCreateDateを抽出します
func (s *ImageRenamerService) extractCreateDateFromJpeg(filePath string) (time.Time, error) {
	// ファイルを読み取り
	data, err := os.ReadFile(filePath)
	if err != nil {
		return time.Time{}, fmt.Errorf("ファイルの読み取りに失敗: %v", err)
	}

	// JPEGパーサーを使用してExifセグメントを取得
	jmp := jpegstructure.NewJpegMediaParser()
	intfc, err := jmp.ParseBytes(data)
	if err != nil {
		return time.Time{}, fmt.Errorf("JPEGの解析に失敗: %v", err)
	}

	sl := intfc.(*jpegstructure.SegmentList)

	// Exifセグメントを取得
	rootIfd, _, err := sl.Exif()
	if err != nil {
		return time.Time{}, fmt.Errorf("Exifデータが見つかりません: %v", err)
	}

	// CreateDateまたはDateTimeOriginalを取得
	if err != nil {
		return time.Time{}, fmt.Errorf("IFDマッピングの作成に失敗: %v", err)
	}

	// まずCreateDateを試す
	results, err := rootIfd.FindTagWithName("DateTime")
	if err == nil && len(results) > 0 {
		value, err := results[0].Value()
		if err == nil {
			if dateStr, ok := value.(string); ok {
				parsedTime, err := time.Parse("2006:01:02 15:04:05", dateStr)
				if err == nil {
					return parsedTime, nil
				}
			}
		}
	}

	// DateTimeOriginalを試す
	results, err = rootIfd.FindTagWithName("DateTimeOriginal")
	if err == nil && len(results) > 0 {
		value, err := results[0].Value()
		if err == nil {
			if dateStr, ok := value.(string); ok {
				parsedTime, err := time.Parse("2006:01:02 15:04:05", dateStr)
				if err == nil {
					return parsedTime, nil
				}
			}
		}
	}

	return time.Time{}, fmt.Errorf("CreateDateまたはDateTimeOriginalが見つかりません")
}

// extractCreateDate はファイル形式に応じてCreateDateを抽出します
func (s *ImageRenamerService) extractCreateDate(filePath string) (time.Time, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".jpg", ".jpeg":
		return s.extractCreateDateFromJpeg(filePath)
	case ".png", ".tiff", ".tif", ".webp", ".mp4", ".webm":
		// PNG, TIFF, WebPは依存関係を減らすためファイルのModTimeを使用
		info, err := os.Stat(filePath)
		if err != nil {
			return time.Time{}, fmt.Errorf("ファイル情報の取得に失敗: %v", err)
		}
		return info.ModTime(), nil
	default:
		return time.Time{}, fmt.Errorf("サポートされていないファイル形式: %s", ext)
	}
}

// generateNewFileName はCreateDateから新しいファイル名を生成します
func (s *ImageRenamerService) generateNewFileName(createDate time.Time, originalPath string) string {
	ext := filepath.Ext(originalPath)
	return createDate.Format("20060102150405") + ext
}

// renameSingleFile は単一ファイルをリネームします
func (s *ImageRenamerService) renameSingleFile(filePath string, config *Config) RenameResult {
	result := RenameResult{
		OriginalPath: filePath,
		Success:      false,
		Error:        nil,
	}

	// CreateDateを抽出
	var createDate time.Time
	var err error

	if config.UseFileModTime {
		// ファイルの更新時刻を使用
		info, err := os.Stat(filePath)
		if err != nil {
			result.Error = fmt.Errorf("ファイル情報の取得に失敗: %v", err)
			return result
		}
		createDate = info.ModTime()
	} else {
		// ExifのCreateDateを使用
		createDate, err = s.extractCreateDate(filePath)
		if err != nil {
			if config.Verbose {
				log.Printf("Exif CreateDateの抽出に失敗: %s - %v", filePath, err)
			}

			// フォールバックでファイルの更新時刻を使用
			info, err := os.Stat(filePath)
			if err != nil {
				result.Error = fmt.Errorf("ファイル情報の取得に失敗: %v", err)
				return result
			}
			createDate = info.ModTime()

			if config.Verbose {
				log.Printf("フォールバック: ファイル更新時刻を使用 %s", filePath)
			}
		}
	}

	// 新しいファイル名を生成
	dir := filepath.Dir(filePath)
	newFileName := s.generateNewFileName(createDate, filePath)
	newPath := filepath.Join(dir, newFileName)

	result.NewPath = newPath

	// 同じ名前のファイルが既に存在する場合の処理
	if filePath != newPath {
		counter := 1
		originalNewPath := newPath
		for {
			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				break
			}
			// 連番を追加
			ext := filepath.Ext(originalNewPath)
			nameWithoutExt := strings.TrimSuffix(originalNewPath, ext)
			newPath = fmt.Sprintf("%s_%02d%s", nameWithoutExt, counter, ext)
			counter++
		}
		result.NewPath = newPath
	}

	// ドライランの場合は実際のリネームを行わない
	if config.DryRun {
		fmt.Printf("  → %s (ドライラン)\n", filepath.Base(newPath))
		result.Success = true
		return result
	}

	// ファイルをリネーム
	if filePath != newPath {
		err = os.Rename(filePath, newPath)
		if err != nil {
			result.Error = fmt.Errorf("ファイルのリネームに失敗: %v", err)
			return result
		}
		fmt.Printf("  → %s\n", filepath.Base(newPath))
	} else {
		fmt.Printf("  → (変更なし)\n")
	}

	result.Success = true
	return result
}

// RenameImageFiles は複数のファイルを並行処理でリネームします
func (s *ImageRenamerService) RenameImageFiles(imageFiles []string, config *Config) (int, int, error) {
	// ワーカー数のデフォルト設定
	workerCount := config.WorkerCount
	if workerCount <= 0 {
		workerCount = 4 // デフォルトは4並列
	}

	// ファイル数がワーカー数より少ない場合、ワーカー数をファイル数に調整
	if len(imageFiles) < workerCount {
		workerCount = len(imageFiles)
	}

	// チャネルの作成
	jobs := make(chan string, len(imageFiles))
	results := make(chan RenameResult, len(imageFiles))

	// ワーカーゴルーチンの起動
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range jobs {
				fmt.Printf("処理中: %s\n", filePath)
				result := s.renameSingleFile(filePath, config)

				if result.Error != nil {
					fmt.Printf("  ⚠️  エラー: %v\n", result.Error)
					if config.Verbose {
						log.Printf("Error processing file %s: %v", filePath, result.Error)
					}
				}

				results <- result
			}
		}()
	}

	// ジョブをチャネルに送信
	for _, filePath := range imageFiles {
		jobs <- filePath
	}
	close(jobs)

	// ワーカーの完了を待つ
	go func() {
		wg.Wait()
		close(results)
	}()

	// 結果を収集
	processedCount := 0
	errorCount := 0
	for result := range results {
		if result.Success {
			processedCount++
		} else {
			errorCount++
		}
	}

	return processedCount, errorCount, nil
}

// ProcessImageRename は画像ファイルの検索とリネームを一括で処理します
func (s *ImageRenamerService) ProcessImageRename(config *Config) (int, int, error) {
	// 画像ファイルを検索
	imageFiles, err := s.FindImageFiles(config)
	if err != nil {
		return 0, 0, fmt.Errorf("画像ファイルの検索に失敗: %w", err)
	}

	if len(imageFiles) == 0 {
		return 0, 0, nil
	}

	if config.Verbose {
		log.Printf("Found %d image files", len(imageFiles))
	}

	// ファイルをリネーム
	return s.RenameImageFiles(imageFiles, config)
}
