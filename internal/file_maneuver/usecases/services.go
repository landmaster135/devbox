package usecases

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Config はfile-maneuverツールの設定を保持する構造体です
type Config struct {
	SrcDirs    []string // ソースディレクトリのリスト
	Extensions []string // 対象拡張子のリスト
	DestDir    string   // 宛先ディレクトリ
	Recursive  bool     // 再帰的検索フラグ
	Workers    int      // ワーカー数
	DryRun     bool     // ドライランフラグ
	CopyMode   bool     // コピーモードフラグ
}

// FileManeuverService はファイル移動サービスを提供する構造体です
type FileManeuverService struct {
	config Config
}

// NewConfig は設定を作成し、全てのバリデーションを実行します
func NewConfig(srcDirs []string, extensions []string, destDir string, recursive bool, workers int, dryRun bool, copyMode bool) (*Config, error) {
	config := &Config{
		SrcDirs:    srcDirs,
		Extensions: extensions,
		DestDir:    destDir,
		Recursive:  recursive,
		Workers:    workers,
		DryRun:     dryRun,
		CopyMode:   copyMode,
	}

	// 構造体作成時に全てのバリデーションを実行
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// NewFileManeuverService はサービスを作成し、設定の妥当性を保証します
func NewFileManeuverService(config *Config) *FileManeuverService {
	// configは既にバリデーション済みなので、そのまま使用可能
	return &FileManeuverService{
		config: *config,
	}
}

// validate は全てのフィールドの妥当性を検証します
func (c *Config) validate() error {
	// 1. ソースディレクトリの検証
	if err := c.validateSrcDirs(); err != nil {
		return fmt.Errorf("ソースディレクトリの検証エラー: %w", err)
	}

	// 2. 拡張子の検証と正規化
	if err := c.validateAndNormalizeExtensions(); err != nil {
		return fmt.Errorf("拡張子の検証エラー: %w", err)
	}

	// 3. 宛先ディレクトリの検証
	if err := c.validateDestDir(); err != nil {
		return fmt.Errorf("宛先ディレクトリの検証エラー: %w", err)
	}

	// 4. ワーカー数の検証と調整
	c.normalizeWorkers()

	return nil
}

// validateSrcDirs はソースディレクトリの妥当性を検証します
func (c *Config) validateSrcDirs() error {
	if len(c.SrcDirs) == 0 {
		return fmt.Errorf("ソースディレクトリが指定されていません")
	}

	for _, dir := range c.SrcDirs {
		if dir == "" {
			return fmt.Errorf("空のディレクトリパスが含まれています")
		}

		info, err := os.Stat(dir)
		if err != nil {
			return fmt.Errorf("ディレクトリ %s にアクセスできません: %w", dir, err)
		}

		if !info.IsDir() {
			return fmt.Errorf("%s はディレクトリではありません", dir)
		}
	}

	return nil
}

// validateAndNormalizeExtensions は拡張子の検証と正規化を行います
func (c *Config) validateAndNormalizeExtensions() error {
	if len(c.Extensions) == 0 {
		return fmt.Errorf("拡張子が指定されていません")
	}

	normalized := make([]string, len(c.Extensions))
	for i, ext := range c.Extensions {
		if ext == "" {
			return fmt.Errorf("空の拡張子が含まれています")
		}

		// 拡張子の正規化（先頭にドットを追加、小文字に変換）
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		normalized[i] = strings.ToLower(ext)
	}

	c.Extensions = normalized
	return nil
}

// validateDestDir は宛先ディレクトリの妥当性を検証します
func (c *Config) validateDestDir() error {
	if c.DestDir == "" {
		return fmt.Errorf("宛先ディレクトリが指定されていません")
	}

	info, err := os.Stat(c.DestDir)
	if err != nil {
		return fmt.Errorf("宛先ディレクトリ %s にアクセスできません: %w", c.DestDir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s はディレクトリではありません", c.DestDir)
	}

	// 書き込み権限の確認
	testFile := filepath.Join(c.DestDir, ".write_test_"+fmt.Sprintf("%d", time.Now().UnixNano()))
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("宛先ディレクトリ %s に書き込み権限がありません: %w", c.DestDir, err)
	}
	os.Remove(testFile) // テストファイルを削除

	return nil
}

// normalizeWorkers はワーカー数を正規化します
func (c *Config) normalizeWorkers() {
	if c.Workers <= 0 {
		c.Workers = runtime.NumCPU()
	}

	// 最大値の制限（過度な並行処理を防ぐ）
	maxWorkers := runtime.NumCPU() * 2
	if c.Workers > maxWorkers {
		c.Workers = maxWorkers
	}
}

// FindTargetFiles は対象ファイルを検索します
func (s *FileManeuverService) FindTargetFiles(stdout, stderr io.Writer) ([]string, error) {
	var allFiles []string
	var mu sync.Mutex

	// 各ソースディレクトリから並行してファイルを検索
	var wg sync.WaitGroup
	errChan := make(chan error, len(s.config.SrcDirs))

	for _, srcDir := range s.config.SrcDirs {
		wg.Add(1)
		go func(dir string) {
			defer wg.Done()

			files, err := s.findFilesInDirectory(dir)
			if err != nil {
				errChan <- fmt.Errorf("ディレクトリ %s の検索エラー: %w", dir, err)
				return
			}

			mu.Lock()
			allFiles = append(allFiles, files...)
			mu.Unlock()

			fmt.Fprintf(stdout, "ディレクトリ %s から %d ファイルを発見しました\n", dir, len(files))
		}(srcDir)
	}

	wg.Wait()
	close(errChan)

	// エラーチェック
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	fmt.Fprintf(stdout, "合計 %d ファイルが見つかりました\n", len(allFiles))
	return allFiles, nil
}

// findFilesInDirectory は指定されたディレクトリからファイルを検索します
func (s *FileManeuverService) findFilesInDirectory(dir string) ([]string, error) {
	var files []string

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// 再帰的検索が無効で、ルートディレクトリでない場合はスキップ
			if !s.config.Recursive && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		// 拡張子チェック
		ext := strings.ToLower(filepath.Ext(path))
		for _, targetExt := range s.config.Extensions {
			if ext == targetExt {
				files = append(files, path)
				break
			}
		}

		return nil
	}

	err := filepath.Walk(dir, walkFunc)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ProcessFiles はファイルを移動またはコピーします
func (s *FileManeuverService) ProcessFiles(files []string, stdout, stderr io.Writer) (int, int, error) {
	if len(files) == 0 {
		if s.config.CopyMode {
			fmt.Fprintf(stdout, "コピー対象のファイルがありません\n")
		} else {
			fmt.Fprintf(stdout, "移動対象のファイルがありません\n")
		}
		return 0, 0, nil
	}

	// 重複チェック
	conflicts, err := s.checkFileConflicts(files)
	if err != nil {
		return 0, 0, fmt.Errorf("ファイル衝突チェックエラー: %w", err)
	}

	if len(conflicts) > 0 {
		fmt.Fprintf(stderr, "警告: 以下のファイルは宛先に同名ファイルが存在するためスキップされます:\n")
		for _, conflict := range conflicts {
			fmt.Fprintf(stderr, "  %s\n", conflict)
		}
	}

	// 衝突するファイルを除外
	validFiles := s.excludeConflictFiles(files, conflicts)
	if s.config.CopyMode {
		fmt.Fprintf(stdout, "%d ファイルをコピーします（%d ファイルをスキップ）\n", len(validFiles), len(conflicts))
	} else {
		fmt.Fprintf(stdout, "%d ファイルを移動します（%d ファイルをスキップ）\n", len(validFiles), len(conflicts))
	}

	if s.config.DryRun {
		if s.config.CopyMode {
			fmt.Fprintf(stdout, "ドライランモード: 実際のコピーは行いません\n")
			for _, file := range validFiles {
				destPath := filepath.Join(s.config.DestDir, filepath.Base(file))
				fmt.Fprintf(stdout, "コピー予定: %s -> %s\n", file, destPath)
			}
		} else {
			fmt.Fprintf(stdout, "ドライランモード: 実際の移動は行いません\n")
			for _, file := range validFiles {
				destPath := filepath.Join(s.config.DestDir, filepath.Base(file))
				fmt.Fprintf(stdout, "移動予定: %s -> %s\n", file, destPath)
			}
		}
		return len(validFiles), 0, nil
	}

	// 並行処理でファイル処理
	return s.processFilesParallel(validFiles, stdout, stderr)
}

// checkFileConflicts はファイル名の衝突をチェックします
func (s *FileManeuverService) checkFileConflicts(files []string) ([]string, error) {
	var conflicts []string

	for _, file := range files {
		destPath := filepath.Join(s.config.DestDir, filepath.Base(file))
		if _, err := os.Stat(destPath); err == nil {
			conflicts = append(conflicts, file)
		}
	}

	return conflicts, nil
}

// excludeConflictFiles は衝突するファイルを除外します
func (s *FileManeuverService) excludeConflictFiles(files []string, conflicts []string) []string {
	conflictMap := make(map[string]bool)
	for _, conflict := range conflicts {
		conflictMap[conflict] = true
	}

	var validFiles []string
	for _, file := range files {
		if !conflictMap[file] {
			validFiles = append(validFiles, file)
		}
	}

	return validFiles
}

// processFilesParallel は並行処理でファイルを移動またはコピーします
func (s *FileManeuverService) processFilesParallel(files []string, stdout, stderr io.Writer) (int, int, error) {
	workerCount := s.config.Workers
	if workerCount > len(files) {
		workerCount = len(files)
	}

	if s.config.CopyMode {
		fmt.Fprintf(stdout, "ファイルコピーに %d ワーカーを使用します\n", workerCount)
	} else {
		fmt.Fprintf(stdout, "ファイル移動に %d ワーカーを使用します\n", workerCount)
	}

	// カウンターとワーカーの同期用
	var mu sync.Mutex
	var wg sync.WaitGroup
	successCount := 0
	errorCount := 0

	// ジョブチャネル
	jobChan := make(chan string, len(files))

	// ワーカーの起動
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobChan {
				if err := s.processFile(file, stdout, stderr); err != nil {
					if s.config.CopyMode {
						fmt.Fprintf(stderr, "エラー: %s のコピーに失敗しました: %v\n", file, err)
					} else {
						fmt.Fprintf(stderr, "エラー: %s の移動に失敗しました: %v\n", file, err)
					}
					mu.Lock()
					errorCount++
					mu.Unlock()
				} else {
					mu.Lock()
					successCount++
					mu.Unlock()
				}
			}
		}()
	}

	// ジョブの送信
	for _, file := range files {
		jobChan <- file
	}
	close(jobChan)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	return successCount, errorCount, nil
}

// moveFile は単一ファイルを移動します
func (s *FileManeuverService) moveFile(srcPath string, stdout, stderr io.Writer) error {
	destPath := filepath.Join(s.config.DestDir, filepath.Base(srcPath))

	fmt.Fprintf(stdout, "移動中: %s -> %s\n", srcPath, destPath)

	err := os.Rename(srcPath, destPath)
	if err != nil {
		return fmt.Errorf("ファイル移動エラー: %w", err)
	}

	return nil
}

// copyFile は単一ファイルをコピーします
func (s *FileManeuverService) copyFile(srcPath string, stdout, stderr io.Writer) error {
	destPath := filepath.Join(s.config.DestDir, filepath.Base(srcPath))

	fmt.Fprintf(stdout, "コピー中: %s -> %s\n", srcPath, destPath)

	// ソースファイルを開く
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("ソースファイルオープンエラー: %w", err)
	}
	defer srcFile.Close()

	// ソースファイルの情報を取得
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("ソースファイル情報取得エラー: %w", err)
	}

	// 宛先ファイルを作成
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("宛先ファイル作成エラー: %w", err)
	}
	defer destFile.Close()

	// ファイル内容をコピー
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("ファイルコピーエラー: %w", err)
	}

	// 権限とタイムスタンプを保持
	err = os.Chmod(destPath, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("権限設定エラー: %w", err)
	}

	err = os.Chtimes(destPath, srcInfo.ModTime(), srcInfo.ModTime())
	if err != nil {
		return fmt.Errorf("タイムスタンプ設定エラー: %w", err)
	}

	return nil
}

// processFile はコピーモードに応じてファイルを処理します
func (s *FileManeuverService) processFile(srcPath string, stdout, stderr io.Writer) error {
	if s.config.CopyMode {
		return s.copyFile(srcPath, stdout, stderr)
	}
	return s.moveFile(srcPath, stdout, stderr)
}

// ExecuteFileManeuver はファイル移動またはコピー処理を一括実行します
func (s *FileManeuverService) ExecuteFileManeuver(stdout, stderr io.Writer) (int, int, error) {
	// 対象ファイルの検索
	fmt.Fprintf(stdout, "🔍 対象ファイルを検索中...\n")
	files, err := s.FindTargetFiles(stdout, stderr)
	if err != nil {
		return 0, 0, fmt.Errorf("ファイル検索エラー: %w", err)
	}

	isCopyMode := s.config.CopyMode

	if len(files) == 0 {
		if isCopyMode {
			fmt.Fprintf(stdout, "✅ コピー対象のファイルが見つかりませんでした\n")
		} else {
			fmt.Fprintf(stdout, "✅ 移動対象のファイルが見つかりませんでした\n")
		}
		return 0, 0, nil
	}

	fmt.Fprintln(stdout)

	// ファイルの処理
	if isCopyMode {
		fmt.Fprintf(stdout, "📦 ファイルコピーを開始します...\n")
	} else {
		fmt.Fprintf(stdout, "📦 ファイル移動を開始します...\n")
	}
	successCount, errorCount, err := s.ProcessFiles(files, stdout, stderr)
	if err != nil {
		if isCopyMode {
			return 0, 0, fmt.Errorf("ファイルコピーエラー: %w", err)
		} else {
			return 0, 0, fmt.Errorf("ファイル移動エラー: %w", err)
		}
	}

	return successCount, errorCount, nil
}
