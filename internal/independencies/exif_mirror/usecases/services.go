package usecases

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// サポートする画像拡張子
var supportedExtensions = []string{".jpg", ".jpeg", ".tiff", ".tif", ".png", ".webp", ".mp4", ".webm"}

// Config はEXIFミラーリングの設定を保持します
type Config struct {
	SourceFolderPath string // ソースフォルダパス
	TargetFolderPath string // ターゲットフォルダパス
	SourceExtension  string // ソース拡張子
	TargetExtension  string // ターゲット拡張子
	Recursive        bool   // 再帰処理
	DryRun           bool   // ドライラン
	Verbose          bool   // 詳細モード
	WorkerCount      int    // 並行処理のワーカー数
}

// ExifMirrorService はEXIFミラーリングサービスです
type ExifMirrorService struct{}

// NewExifMirrorService は新しいExifMirrorServiceを作成します
func NewExifMirrorService() *ExifMirrorService {
	return &ExifMirrorService{}
}

// MirrorExifData は指定された設定に基づいてEXIFデータをミラーリングします
func (s *ExifMirrorService) MirrorExifData(config *Config) (int, int, error) {
	// ターゲットファイルを検索
	targetFiles, err := s.findImageFiles(config.TargetFolderPath, config.TargetExtension, config.Recursive)
	if err != nil {
		return 0, 0, fmt.Errorf("ターゲットファイルの検索に失敗: %v", err)
	}

	if len(targetFiles) == 0 {
		return 0, 0, fmt.Errorf("ターゲットファイルが見つかりません")
	}

	if config.Verbose {
		log.Printf("Found %d target files", len(targetFiles))
	}

	// ワーカー数のデフォルト設定
	workerCount := config.WorkerCount
	if workerCount <= 0 {
		workerCount = 4 // デフォルトは4並列
	}

	// ファイル数がワーカー数より少ない場合、ワーカー数をファイル数に調整
	if len(targetFiles) < workerCount {
		workerCount = len(targetFiles)
	}

	// チャネルの作成
	jobs := make(chan string, len(targetFiles))
	results := make(chan MirrorResult, len(targetFiles))

	// ワーカーゴルーチンの起動
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for targetFilePath := range jobs {
				result := MirrorResult{
					TargetFilePath: targetFilePath,
					Success:        false,
					Error:          nil,
				}

				// 対応するソースファイルを検索
				sourceFilePath := s.findCorrespondingSourceFile(targetFilePath, config)
				if sourceFilePath == "" {
					result.Error = fmt.Errorf("対応するソースファイルが見つかりません")
					results <- result
					continue
				}

				result.SourceFilePath = sourceFilePath

				if config.DryRun {
					fmt.Printf("処理対象: %s <- %s\n", targetFilePath, sourceFilePath)
					result.Success = true
				} else {
					fmt.Printf("処理中: %s <- %s\n", targetFilePath, sourceFilePath)
					err := s.copyExifData(sourceFilePath, targetFilePath, config)
					if err != nil {
						fmt.Printf("  ⚠️  エラー: %v\n", err)
						result.Error = err
						if config.Verbose {
							log.Printf("Error processing file %s: %v", targetFilePath, err)
						}
					} else {
						fmt.Printf("  ✅ EXIFデータをコピーしました\n")
						result.Success = true
					}
				}
				results <- result
			}
		}()
	}

	// ジョブをチャネルに送信
	for _, filePath := range targetFiles {
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

// MirrorResult はファイル処理の結果を表します
type MirrorResult struct {
	TargetFilePath string
	SourceFilePath string
	Success        bool
	Error          error
}

// findImageFiles は指定された設定に基づいて画像ファイルを検索します
func (s *ExifMirrorService) findImageFiles(folderPath, extension string, recursive bool) ([]string, error) {
	var imageFiles []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 再帰フラグが設定されていない場合、サブディレクトリをスキップ
		if !recursive && info.IsDir() && path != folderPath {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			if s.isTargetFile(path, extension) {
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

// isTargetFile は対象ファイルかどうかをチェック
func (s *ExifMirrorService) isTargetFile(filePath, targetExtension string) bool {
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

// findCorrespondingSourceFile は対応するソースファイルを検索します
func (s *ExifMirrorService) findCorrespondingSourceFile(targetFilePath string, config *Config) string {
	// ターゲットファイルの基本情報を取得
	targetDir := filepath.Dir(targetFilePath)
	targetBaseName := removeExtension(filepath.Base(targetFilePath))

	// ソースフォルダでの相対パスを計算
	relPath, err := filepath.Rel(config.TargetFolderPath, targetDir)
	if err != nil {
		// 相対パス計算に失敗した場合、同じディレクトリで検索
		relPath = "."
	}

	// ソースディレクトリを決定
	sourceDir := config.SourceFolderPath
	if relPath != "." {
		sourceDir = filepath.Join(config.SourceFolderPath, relPath)
	}

	// ソース拡張子を正規化
	sourceExtension := config.SourceExtension
	if !strings.HasPrefix(sourceExtension, ".") {
		sourceExtension = "." + sourceExtension
	}

	// ソースファイルパスを構築
	sourceFilePath := filepath.Join(sourceDir, targetBaseName+sourceExtension)

	// ファイルの存在確認
	if _, err := os.Stat(sourceFilePath); err == nil {
		return sourceFilePath
	}

	return ""
}

// copyExifData はソースファイルからターゲットファイルにEXIFデータをコピーします
// exiftoolがある場合は使用し、なければ基本的なファイル操作で代替
func (s *ExifMirrorService) copyExifData(sourceFilePath, targetFilePath string, config *Config) error {
	if config.Verbose {
		log.Printf("Copying EXIF from %s to %s", sourceFilePath, targetFilePath)
	}

	// exiftoolを使用してEXIFデータをコピーする方法を試す
	if s.hasExifTool() {
		return s.copyExifWithExifTool(sourceFilePath, targetFilePath, config)
	}

	// exiftoolが利用できない場合は基本的な処理のみ
	if config.Verbose {
		log.Printf("exiftool not available, skipping EXIF copy for %s", targetFilePath)
	}
	
	return fmt.Errorf("EXIF copying requires exiftool to be installed")
}

// hasExifTool はexiftoolが利用可能かチェックします
func (s *ExifMirrorService) hasExifTool() bool {
	_, err := exec.LookPath("exiftool")
	return err == nil
}

// copyExifWithExifTool はexiftoolを使用してEXIFデータをコピーします
func (s *ExifMirrorService) copyExifWithExifTool(sourceFilePath, targetFilePath string, config *Config) error {
	// exiftoolコマンドを構築
	args := []string{
		"-tagsFromFile",
		sourceFilePath,
		"-exif:all",
		"-overwrite_original",
		targetFilePath,
	}

	if config.Verbose {
		log.Printf("Running: exiftool %s", strings.Join(args, " "))
	}

	// exiftoolを実行
	cmd := exec.Command("exiftool", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("exiftool failed: %v, stderr: %s", err, stderr.String())
	}

	if config.Verbose {
		log.Printf("exiftool output: %s", stdout.String())
	}

	return nil
}

// removeExtension はファイル名から拡張子を除去します
func removeExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext != "" {
		return fileName[:len(fileName)-len(ext)]
	}
	return fileName
}

// ValidateDirectory はディレクトリの存在と権限をバリデーションします
func ValidateDirectory(dirPath string) error {
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

	return nil
}

// ValidateExtension はファイル拡張子をバリデーションします
func ValidateExtension(ext string) error {
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

// CopyFileExifSimple は単純なファイルベースのEXIFコピー（fallback用）
func (s *ExifMirrorService) CopyFileExifSimple(sourceFilePath, targetFilePath string) error {
	// この実装は、ライブラリを使わない基本的なアプローチです
	// 実際のEXIF操作は複雑なため、exiftoolの使用を推奨します
	
	sourceInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return fmt.Errorf("ソースファイルの情報取得に失敗: %v", err)
	}
	
	targetInfo, err := os.Stat(targetFilePath)
	if err != nil {
		return fmt.Errorf("ターゲットファイルの情報取得に失敗: %v", err)
	}
	
	// 少なくともファイルの変更時刻をコピー
	err = os.Chtimes(targetFilePath, sourceInfo.ModTime(), sourceInfo.ModTime())
	if err != nil {
		return fmt.Errorf("ファイル時刻の更新に失敗: %v", err)
	}
	
	// ターゲットファイルのサイズが変わっていないことを確認
	newInfo, err := os.Stat(targetFilePath)
	if err != nil {
		return fmt.Errorf("処理後のファイル確認に失敗: %v", err)
	}
	
	if newInfo.Size() != targetInfo.Size() {
		return fmt.Errorf("ファイルサイズが変わりました")
	}
	
	return nil
}

// BackupFile はファイルのバックアップを作成します
func (s *ExifMirrorService) BackupFile(filePath string) (string, error) {
	backupPath := filePath + ".backup"
	
	source, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("ソースファイルを開けません: %v", err)
	}
	defer source.Close()
	
	backup, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("バックアップファイルを作成できません: %v", err)
	}
	defer backup.Close()
	
	_, err = io.Copy(backup, source)
	if err != nil {
		os.Remove(backupPath) // 失敗時は作成したファイルを削除
		return "", fmt.Errorf("ファイルのコピーに失敗: %v", err)
	}
	
	return backupPath, nil
}
