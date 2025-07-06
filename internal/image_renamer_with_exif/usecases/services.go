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


// ConflictResolver は競合解決を行います
type ConflictResolver struct {
	usedFileNames map[string]bool // 使用済みファイル名を追跡
}

// NewConflictResolver は新しいConflictResolverを作成します
func NewConflictResolver(existingFiles map[string]bool) *ConflictResolver {
	return &ConflictResolver{
		usedFileNames: existingFiles,
	}
}

// findNextAvailableTime は利用可能な次の時刻を検索します
func (r *ConflictResolver) findNextAvailableTime(baseTime time.Time, ext string, directory string) time.Time {
	candidate := baseTime
	for {
		fileName := candidate.Format("20060102150405") + ext
		fullPath := filepath.Join(directory, fileName)

		if !r.usedFileNames[fullPath] {
			r.usedFileNames[fullPath] = true // 使用済みとしてマーク
			return candidate
		}
		candidate = candidate.Add(1 * time.Second) // 1秒進める
	}
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

// RenameResult はファイルリネームの結果を表します
type RenameResult struct {
	OriginalPath string
	NewPath      string
	Success      bool
	Error        error
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

// FileRenameInfo はリネーム前の情報を保持します
type FileRenameInfo struct {
	OriginalPath string
	NewFileName  string
	CreateDate   time.Time
	Directory    string
}

// prepareRenameInfo は全ファイルのリネーム情報を事前に準備します
func (s *ImageRenamerService) prepareRenameInfo(imageFiles []string, config *Config) ([]FileRenameInfo, error) {
	var renameInfos []FileRenameInfo

	for _, filePath := range imageFiles {
		// CreateDateを抽出
		var createDate time.Time
		var err error

		if config.UseFileModTime {
			// ファイルの更新時刻を使用
			info, err := os.Stat(filePath)
			if err != nil {
				return nil, fmt.Errorf("ファイル情報の取得に失敗: %s - %v", filePath, err)
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
					return nil, fmt.Errorf("ファイル情報の取得に失敗: %s - %v", filePath, err)
				}
				createDate = info.ModTime()

				if config.Verbose {
					log.Printf("フォールバック: ファイル更新時刻を使用 %s", filePath)
				}
			}
		}

		// 新しいファイル名を生成
		newFileName := s.generateNewFileName(createDate, filePath)
		directory := filepath.Dir(filePath)

		renameInfo := FileRenameInfo{
			OriginalPath: filePath,
			NewFileName:  newFileName,
			CreateDate:   createDate,
			Directory:    directory,
		}

		renameInfos = append(renameInfos, renameInfo)
	}

	return renameInfos, nil
}

// collectExistingFileNames は既存ファイル名を収集します
func (s *ImageRenamerService) collectExistingFileNames(renameInfos []FileRenameInfo) map[string]bool {
	existingFiles := make(map[string]bool)

	// 各ディレクトリの既存ファイルを収集
	processedDirs := make(map[string]bool)
	for _, info := range renameInfos {
		if processedDirs[info.Directory] {
			continue
		}
		processedDirs[info.Directory] = true

		// ディレクトリ内の既存ファイルを取得
		entries, err := os.ReadDir(info.Directory)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				fullPath := filepath.Join(info.Directory, entry.Name())
				existingFiles[fullPath] = true
			}
		}
	}

	return existingFiles
}

// resolveConflicts は競合を自動解決します
func (s *ImageRenamerService) resolveConflicts(renameInfos []FileRenameInfo, config *Config) error {
	// 既存ファイル名を収集
	existingFiles := s.collectExistingFileNames(renameInfos)

	// ConflictResolverを作成
	resolver := NewConflictResolver(existingFiles)

	// ディレクトリ別に競合をグループ化
	dirFileMap := make(map[string]map[string][]*FileRenameInfo)
	for i := range renameInfos {
		info := &renameInfos[i]
		if dirFileMap[info.Directory] == nil {
			dirFileMap[info.Directory] = make(map[string][]*FileRenameInfo)
		}
		dirFileMap[info.Directory][info.NewFileName] = append(
			dirFileMap[info.Directory][info.NewFileName],
			info,
		)
	}

	// 競合を解決
	var resolvedConflicts []string
	for dir, fileMap := range dirFileMap {
		for newFileName, conflictInfos := range fileMap {
			if len(conflictInfos) > 1 {
				// 競合が検出された場合
				resolvedConflicts = append(resolvedConflicts, fmt.Sprintf(
					"競合解決: '%s' (ディレクトリ: %s)",
					newFileName, dir,
				))

				// 元の時刻順にソート
				sort.Slice(conflictInfos, func(i, j int) bool {
					return conflictInfos[i].CreateDate.Before(conflictInfos[j].CreateDate)
				})

				// 最初のファイルは元の時刻を維持、残りは1秒ずつ後の時刻を割り当て
				for i, info := range conflictInfos {
					if i == 0 {
						// 最初のファイルも既存ファイルとの競合をチェック
						ext := filepath.Ext(info.OriginalPath)
						newTime := resolver.findNextAvailableTime(info.CreateDate, ext, info.Directory)
						if !newTime.Equal(info.CreateDate) {
							info.CreateDate = newTime
							info.NewFileName = s.generateNewFileName(newTime, info.OriginalPath)
							resolvedConflicts = append(resolvedConflicts, fmt.Sprintf(
								"  - %s → %s (既存ファイルとの競合回避)",
								filepath.Base(info.OriginalPath), info.NewFileName,
							))
						} else {
							resolvedConflicts = append(resolvedConflicts, fmt.Sprintf(
								"  - %s → %s (元の時刻維持)",
								filepath.Base(info.OriginalPath), info.NewFileName,
							))
						}
					} else {
						// 2番目以降のファイルは安全な時刻を検索
						ext := filepath.Ext(info.OriginalPath)
						newTime := resolver.findNextAvailableTime(info.CreateDate.Add(time.Duration(i)*time.Second), ext, info.Directory)
						info.CreateDate = newTime
						info.NewFileName = s.generateNewFileName(newTime, info.OriginalPath)
						resolvedConflicts = append(resolvedConflicts, fmt.Sprintf(
							"  - %s → %s (%d秒後に調整)",
							filepath.Base(info.OriginalPath), info.NewFileName, int(newTime.Sub(conflictInfos[0].CreateDate).Seconds()),
						))
					}
				}
				resolvedConflicts = append(resolvedConflicts, "")
			}
		}
	}

	// 競合解決の結果を表示
	if len(resolvedConflicts) > 0 {
		fmt.Println("ファイル名の競合が検出されました。自動解決を実行します:")
		fmt.Println()
		for _, msg := range resolvedConflicts {
			fmt.Println(msg)
		}
	}

	return nil
}

// renameImageFiles は複数のファイルを並行処理でリネームします
func (s *ImageRenamerService) renameImageFiles(imageFiles []string, config *Config) (int, int, error) {
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

	// 全ファイルのリネーム情報を事前に準備
	renameInfos, err := s.prepareRenameInfo(imageFiles, config)
	if err != nil {
		return 0, 0, fmt.Errorf("リネーム情報の準備に失敗: %w", err)
	}

	// 競合を自動解決
	if err := s.resolveConflicts(renameInfos, config); err != nil {
		return 0, 0, err
	}

	// リネーム実行
	return s.renameImageFiles(imageFiles, config)
}
