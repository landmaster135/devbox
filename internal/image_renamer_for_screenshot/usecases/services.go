package usecases

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var xiaomiScreenshotRegexp = regexp.MustCompile(`^Screenshot_(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d+)_([A-Za-z0-9._-]+)$`)

// FileInfo はファイル情報を保持する構造体です
type FileInfo struct {
	Path string
	Name string
}

// Operation はリネーム対象となるスクリーンショットの種別を表します
type Operation string

const (
	OperationUnknown Operation = ""
	OperationVLC     Operation = "vlc"
	OperationWin     Operation = "win"
	OperationPixel   Operation = "pixel"
	OperationXiaomi  Operation = "xiaomi"
	operationAll     Operation = "all"
)

func (o Operation) isValidSelection() bool {
	switch o {
	case OperationVLC, OperationWin, OperationPixel, OperationXiaomi:
		return true
	default:
		return false
	}
}

// Config はプログラムの設定を保持する構造体です
type Config struct {
	SrcDir     string
	Recursive  bool
	Workers    int
	Operation  Operation
	ToDateTime bool
}

// validateConfig は設定の妥当性を検証します
func validateConfig(config Config, stderr io.Writer) error {
	// --to-datetimeが指定されている場合は、他のパターンは不要
	if config.ToDateTime {
		// ディレクトリの存在確認
		_, err := os.Stat(config.SrcDir)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", config.SrcDir, err)
			return err
		}
		return nil
	}

	if !config.Operation.isValidSelection() {
		fmt.Fprintln(stderr, "エラー: -operation には 'vlc'、'win'、'pixel'、または 'xiaomi' のいずれかを指定する必要があります。")
		fmt.Fprintln(stderr, "例: ./image-renamer-for-screenshot -operation=vlc")
		fmt.Fprintln(stderr, "例: ./image-renamer-for-screenshot -to-datetime")
		return fmt.Errorf("無効なoperationが指定されています")
	}

	// ディレクトリの存在確認
	_, err := os.Stat(config.SrcDir)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: ディレクトリ %s へのアクセスエラー: %v\n", config.SrcDir, err)
		return err
	}

	return nil
}

func isImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".avif", ".mp4":
		return true
	default:
		return false
	}
}

// findScreenshotFilesForDateTime は--to-datetimeフラグ用にスクリーンショットファイルを検索します
func findScreenshotFilesForDateTime(srcDir string, recursive bool, stdout, stderr io.Writer) ([]string, error) {
	// 既存の関数を使って基本パターンを検索
	files, err := findScreenshotFiles(srcDir, recursive, operationAll, stdout, stderr)
	if err != nil {
		return nil, err
	}

	// Screenshot_パターンを追加で検索
	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(d.Name()))
			if isImageExt(ext) {
				name := d.Name()
				if strings.HasPrefix(name, "Screenshot_") {
					files = append(files, path)
				}
			}
		}
		return nil
	}

	if recursive {
		err := filepath.WalkDir(srcDir, walkFunc)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の走査中にエラーが発生しました: %v\n", srcDir, err)
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の読み込みに失敗しました: %v\n", srcDir, err)
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				ext := strings.ToLower(filepath.Ext(entry.Name()))
				if isImageExt(ext) {
					name := entry.Name()
					if strings.HasPrefix(name, "Screenshot_") {
						files = append(files, filepath.Join(srcDir, name))
					}
				}
			}
		}
	}

	return files, nil
}

// findScreenshotFiles は指定されたディレクトリからスクリーンショットファイルを検索します
func findScreenshotFiles(srcDir string, recursive bool, operation Operation, stdout, stderr io.Writer) ([]string, error) {
	var files []string

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			name := d.Name()
			origExt := filepath.Ext(name)
			ext := strings.ToLower(origExt)
			if isImageExt(ext) {
				baseName := strings.TrimSuffix(name, origExt)
				matchesXiaomi := false
				if operation == OperationXiaomi || operation == operationAll {
					matchesXiaomi = xiaomiScreenshotRegexp.MatchString(baseName)
				}
				if matchesOperation(name, matchesXiaomi, operation) {
					files = append(files, path)
				}
			}
		}
		return nil
	}

	if recursive {
		err := filepath.WalkDir(srcDir, walkFunc)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の走査中にエラーが発生しました: %v\n", srcDir, err)
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(srcDir)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ %s の読み込みに失敗しました: %v\n", srcDir, err)
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				name := entry.Name()
				origExt := filepath.Ext(name)
				ext := strings.ToLower(origExt)
				if isImageExt(ext) {
					baseName := strings.TrimSuffix(name, origExt)
					matchesXiaomi := false
					if operation == OperationXiaomi || operation == operationAll {
						matchesXiaomi = xiaomiScreenshotRegexp.MatchString(baseName)
					}
					if matchesOperation(name, matchesXiaomi, operation) {
						files = append(files, filepath.Join(srcDir, name))
					}
				}
			}
		}
	}

	return files, nil
}

func matchesOperation(name string, matchesXiaomi bool, operation Operation) bool {
	switch operation {
	case OperationVLC:
		return strings.HasPrefix(name, "vlcsnap-")
	case OperationWin:
		return strings.HasPrefix(name, "スクリーンショット ")
	case OperationPixel:
		return strings.HasPrefix(name, "screen-")
	case OperationXiaomi:
		return matchesXiaomi
	case operationAll:
		return strings.HasPrefix(name, "vlcsnap-") ||
			strings.HasPrefix(name, "スクリーンショット ") ||
			strings.HasPrefix(name, "screen-") ||
			matchesXiaomi
	default:
		return false
	}
}

// getFileInfos はファイルパスのリストからファイル情報を取得します
func getFileInfos(files []string, stderr io.Writer) ([]FileInfo, error) {
	fileInfos := make([]FileInfo, len(files))
	var hasError bool

	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ファイル %s の情報取得に失敗しました: %v\n", file, err)
			hasError = true
			continue
		}
		fileInfos[i] = FileInfo{
			Path: file,
			Name: info.Name(),
		}
	}

	if hasError {
		return fileInfos, fmt.Errorf("一部のファイル情報の取得に失敗しました")
	}
	return fileInfos, nil
}

// renameVlcToDateTime はVLCスクリーンショットファイルをYYYYMMDDHHMMSS形式にリネームします
func renameVlcToDateTime(baseName, ext string) (string, error) {
	// パターン1: vlcsnap-YYYY-MM-DD-HH-MM-SS
	re1 := regexp.MustCompile(`vlcsnap-(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})`)
	matches1 := re1.FindStringSubmatch(baseName)
	if len(matches1) == 7 {
		year, month, day := matches1[1], matches1[2], matches1[3]
		hour, minute, second := matches1[4], matches1[5], matches1[6]
		return fmt.Sprintf("%s%s%s%s%s%s%s", year, month, day, hour, minute, second, ext), nil
	}

	// パターン2: vlcsnap-YYYY-MM-DD-HHhMMmSSsNNN
	re2 := regexp.MustCompile(`vlcsnap-(\d{4})-(\d{2})-(\d{2})-(\d{2})h(\d{2})m(\d{2})s\d+`)
	matches2 := re2.FindStringSubmatch(baseName)
	if len(matches2) == 7 {
		year, month, day := matches2[1], matches2[2], matches2[3]
		hour, minute, second := matches2[4], matches2[5], matches2[6]
		return fmt.Sprintf("%s%s%s%s%s%s%s", year, month, day, hour, minute, second, ext), nil
	}

	return "", fmt.Errorf("VLCスクリーンショットのパターンに一致しません: %s", baseName)
}

// renameWindowsToDateTime はWindowsスクリーンショットファイルをYYYYMMDDHHMMSS形式にリネームします
func renameWindowsToDateTime(baseName, ext string) (string, error) {
	// スクリーンショット YYYY-MM-DD HHMMSS
	re := regexp.MustCompile(`スクリーンショット (\d{4})-(\d{2})-(\d{2}) (\d{2})(\d{2})(\d{2})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 7 {
		return "", fmt.Errorf("windowsスクリーンショットのパターンに一致しません: %s", baseName)
	}

	year, month, day := matches[1], matches[2], matches[3]
	hour, minute, second := matches[4], matches[5], matches[6]

	return fmt.Sprintf("%s%s%s%s%s%s%s", year, month, day, hour, minute, second, ext), nil
}

// renamePixelToDateTime はPixelスクリーンレコードファイルをYYYYMMDDHHMMSS形式にリネームします
func renamePixelToDateTime(baseName, ext string) (string, error) {
	// screen-YYYYMMDD-HHMMSS
	re := regexp.MustCompile(`screen-(\d{8})-(\d{6})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return "", fmt.Errorf("Pixelスクリーンレコードのパターンに一致しません: %s", baseName)
	}

	dateStr := matches[1]
	timeStr := matches[2]

	if len(dateStr) != 8 || len(timeStr) != 6 {
		return "", fmt.Errorf("Pixelスクリーンレコードの日時形式が不正です: %s", baseName)
	}

	return fmt.Sprintf("%s%s%s", dateStr, timeStr, ext), nil
}

// renameScreenshotToDateTime はScreenshot_ファイルをYYYYMMDDHHMMSS形式にリネームします
func renameScreenshotToDateTime(baseName, ext string) (string, error) {
	// Screenshot_YYYYMMDD-HHMMSS
	re := regexp.MustCompile(`Screenshot_(\d{8})-(\d{6})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return "", fmt.Errorf("Screenshot_パターンに一致しません: %s", baseName)
	}

	dateStr := matches[1]
	timeStr := matches[2]

	if len(dateStr) != 8 || len(timeStr) != 6 {
		return "", fmt.Errorf("Screenshot_の日時形式が不正です: %s", baseName)
	}

	return fmt.Sprintf("%s%s%s", dateStr, timeStr, ext), nil
}

// processScreenshotRenameToDateTime は--to-datetimeフラグ用のリネーム処理を行います
func processScreenshotRenameToDateTime(file FileInfo, mu *sync.Mutex, successCount, errorCount *int, stdout, stderr io.Writer) {
	oldPath := file.Path
	dir := filepath.Dir(oldPath)
	oldName := filepath.Base(oldPath)
	ext := filepath.Ext(oldName)
	baseName := strings.TrimSuffix(oldName, ext)

	var newName string
	var err error

	// 各パターンを試行してYYYYMMDDHHMMSS形式に変換
	if strings.HasPrefix(baseName, "vlcsnap-") {
		newName, err = renameVlcToDateTime(baseName, ext)
	} else if strings.HasPrefix(baseName, "スクリーンショット ") {
		newName, err = renameWindowsToDateTime(baseName, ext)
	} else if strings.HasPrefix(baseName, "screen-") {
		newName, err = renamePixelToDateTime(baseName, ext)
	} else if xiaomiScreenshotRegexp.MatchString(baseName) {
		newName, err = renameXiaomiToDateTime(baseName, ext)
	} else if strings.HasPrefix(baseName, "Screenshot_") {
		newName, err = renameScreenshotToDateTime(baseName, ext)
	} else {
		// パターンに一致しないファイルはスキップ
		return
	}

	if err != nil {
		fmt.Fprintf(stderr, "エラー: %s の解析に失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
		return
	}

	newPath := filepath.Join(dir, newName)
	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
	} else {
		mu.Lock()
		*successCount++
		mu.Unlock()
	}
}

// renameVlcScreenshot はVLCスクリーンショットファイルをリネームします
func renameVlcScreenshot(baseName, ext string) (string, error) {
	// パターン1: vlcsnap-YYYY-MM-DD-HH-MM-SS
	re1 := regexp.MustCompile(`vlcsnap-(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})`)
	matches1 := re1.FindStringSubmatch(baseName)
	if len(matches1) == 7 {
		year, month, day := matches1[1], matches1[2], matches1[3]
		hour, minute, second := matches1[4], matches1[5], matches1[6]
		return fmt.Sprintf("Screenshot_%s%s%s-%s%s%s%s", year, month, day, hour, minute, second, ext), nil
	}

	// パターン2: vlcsnap-YYYY-MM-DD-HHhMMmSSsNNN
	re2 := regexp.MustCompile(`vlcsnap-(\d{4})-(\d{2})-(\d{2})-(\d{2})h(\d{2})m(\d{2})s\d+`)
	matches2 := re2.FindStringSubmatch(baseName)
	if len(matches2) == 7 {
		year, month, day := matches2[1], matches2[2], matches2[3]
		hour, minute, second := matches2[4], matches2[5], matches2[6]
		return fmt.Sprintf("Screenshot_%s%s%s-%s%s%s%s", year, month, day, hour, minute, second, ext), nil
	}

	return "", fmt.Errorf("VLCスクリーンショットのパターンに一致しません: %s", baseName)
}

// renameWindowsScreenshot はWindowsスクリーンショットファイルをリネームします
func renameWindowsScreenshot(baseName, ext string) (string, error) {
	// スクリーンショット YYYY-MM-DD HH-MM-SS
	re := regexp.MustCompile(`スクリーンショット (\d{4})-(\d{2})-(\d{2}) (\d{2})(\d{2})(\d{2})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 7 {
		return "", fmt.Errorf("[Error] Windowsスクリーンショットのパターンに一致しません: %s", baseName)
	}

	year, month, day := matches[1], matches[2], matches[3]
	hour, minute, second := matches[4], matches[5], matches[6]

	return fmt.Sprintf("Screenshot_%s%s%s-%s%s%s%s", year, month, day, hour, minute, second, ext), nil
}

// renamePixelScreenshot はPixelスクリーンレコードファイルをリネームします
func renamePixelScreenshot(baseName, ext string) (string, error) {
	// screen-YYYYMMDD-HHMMSS
	re := regexp.MustCompile(`screen-(\d{8})-(\d{6})`)
	matches := re.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return "", fmt.Errorf("[Error] Pixelスクリーンレコードのパターンに一致しません: %s", baseName)
	}

	dateStr := matches[1]
	timeStr := matches[2]

	if len(dateStr) != 8 || len(timeStr) != 6 {
		return "", fmt.Errorf("[Error] Pixelスクリーンレコードの日時形式が不正です: %s", baseName)
	}

	year := dateStr[0:4]
	month := dateStr[4:6]
	day := dateStr[6:8]

	hour := timeStr[0:2]
	minute := timeStr[2:4]
	second := timeStr[4:6]

	return fmt.Sprintf("Screenshot_%s%s%s-%s%s%s%s", year, month, day, hour, minute, second, ext), nil
}

func parseXiaomiDateTime(baseName string) (string, string, error) {
	matches := xiaomiScreenshotRegexp.FindStringSubmatch(baseName)
	if len(matches) != 9 {
		return "", "", fmt.Errorf("Xiaomiスクリーンショットのパターンに一致しません: %s", baseName)
	}

	dateStr := matches[1] + matches[2] + matches[3]
	timeStr := matches[4] + matches[5] + matches[6]
	return dateStr, timeStr, nil
}

// renameXiaomiScreenshot はXiaomi端末のスクリーンショットファイルをリネームします
func renameXiaomiScreenshot(baseName, ext string) (string, error) {
	dateStr, timeStr, err := parseXiaomiDateTime(baseName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Screenshot_%s-%s%s", dateStr, timeStr, ext), nil
}

func renameXiaomiToDateTime(baseName, ext string) (string, error) {
	dateStr, timeStr, err := parseXiaomiDateTime(baseName)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%s", dateStr, timeStr, ext), nil
}

// processScreenshotRename は1つのスクリーンショットファイルをリネームします
func processScreenshotRename(file FileInfo, operation Operation, mu *sync.Mutex, successCount, errorCount *int, stdout, stderr io.Writer) {
	oldPath := file.Path
	dir := filepath.Dir(oldPath)
	oldName := filepath.Base(oldPath)
	ext := filepath.Ext(oldName)
	baseName := strings.TrimSuffix(oldName, ext)

	var newName string
	var err error

	switch operation {
	case OperationVLC:
		if !strings.HasPrefix(baseName, "vlcsnap-") {
			return
		}
		newName, err = renameVlcScreenshot(baseName, ext)
	case OperationWin:
		if !strings.HasPrefix(baseName, "スクリーンショット ") {
			return
		}
		newName, err = renameWindowsScreenshot(baseName, ext)
	case OperationPixel:
		if !strings.HasPrefix(baseName, "screen-") {
			return
		}
		newName, err = renamePixelScreenshot(baseName, ext)
	case OperationXiaomi:
		if !xiaomiScreenshotRegexp.MatchString(baseName) {
			return
		}
		newName, err = renameXiaomiScreenshot(baseName, ext)
	default:
		return
	}

	if err != nil {
		fmt.Fprintf(stderr, "エラー: %s の解析に失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
		return
	}

	newPath := filepath.Join(dir, newName)
	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)

	err = os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
	} else {
		mu.Lock()
		*successCount++
		mu.Unlock()
	}
}

// renameScreenshotFiles はスクリーンショットファイルをリネームします
func renameScreenshotFiles(fileInfos []FileInfo, config Config, stdout, stderr io.Writer) (int, int) {
	// ワーカープールの設定
	workerCount := config.Workers
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > len(fileInfos) {
		workerCount = len(fileInfos)
	}

	fmt.Fprintf(stdout, "リネーム操作に %d ワーカーを使用します。\n", workerCount)

	// カウンターとワーカーの同期用
	var mu sync.Mutex
	var wg sync.WaitGroup
	successCount := 0
	errorCount := 0

	// ジョブチャネル
	jobChan := make(chan FileInfo, len(fileInfos))

	// ワーカーの起動
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobChan {
				if config.ToDateTime {
					processScreenshotRenameToDateTime(file, &mu, &successCount, &errorCount, stdout, stderr)
				} else {
					processScreenshotRename(file, config.Operation, &mu, &successCount, &errorCount, stdout, stderr)
				}
			}
		}()
	}

	// ジョブの送信
	for _, file := range fileInfos {
		jobChan <- file
	}
	close(jobChan)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	return successCount, errorCount
}

// ProcessScreenshotRename はスクリーンショットファイルのリネーム処理を統合して実行します
func ProcessScreenshotRename(config Config, stdout, stderr io.Writer) (int, int, error) {
	// 設定の妥当性を検証
	if err := validateConfig(config, stderr); err != nil {
		return 0, 0, err
	}

	// スクリーンショットファイルの検索
	var files []string
	var err error
	if config.ToDateTime {
		files, err = findScreenshotFilesForDateTime(config.SrcDir, config.Recursive, stdout, stderr)
	} else {
		files, err = findScreenshotFiles(config.SrcDir, config.Recursive, config.Operation, stdout, stderr)
	}
	if err != nil {
		return 0, 0, err
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "スクリーンショットファイルが見つかりませんでした。")
		return 0, 0, nil
	}

	fmt.Fprintf(stdout, "スクリーンショットファイルが %d 件見つかりました。\n", len(files))

	if config.ToDateTime {
		fmt.Fprintln(stdout, "YYYYMMDDHHMMSS形式でのリネームを実行します。")
	} else {
		switch config.Operation {
		case OperationVLC:
			fmt.Fprintln(stdout, "VLCスナップショットパターンを使用します。")
		case OperationWin:
			fmt.Fprintln(stdout, "Windowsスクリーンショットパターンを使用します。")
		case OperationPixel:
			fmt.Fprintln(stdout, "Pixelスクリーンレコードパターンを使用します。")
		case OperationXiaomi:
			fmt.Fprintln(stdout, "Xiaomiスクリーンショットパターンを使用します。")
		}
	}

	// ファイル情報の取得
	fileInfos, err := getFileInfos(files, stderr)
	if err != nil {
		// エラーがあっても続行するため、ここではエラーコードを返さない
	}

	// リネーム処理の実行
	successCount, errorCount := renameScreenshotFiles(fileInfos, config, stdout, stderr)

	// 処理結果の出力
	fmt.Fprintf(stdout, "✔ ファイルリネームが完了しました\n")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", successCount)

	if errorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", errorCount)
		return successCount, errorCount, fmt.Errorf("一部のファイルのリネームに失敗しました")
	}

	return successCount, errorCount, nil
}
