package usecases

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var xiaomiScreenshotRegexp = regexp.MustCompile(`^Screenshot_(\d{4})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d{2})-(\d+)_([A-Za-z0-9._-]+)$`)

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
	switch {
	case strings.HasPrefix(baseName, "vlcsnap-"):
		newName, err = renameVlcToDateTime(baseName, ext)
	case strings.HasPrefix(baseName, "スクリーンショット "):
		newName, err = renameWindowsToDateTime(baseName, ext)
	case strings.HasPrefix(baseName, "screen-"):
		newName, err = renamePixelToDateTime(baseName, ext)
	case xiaomiScreenshotRegexp.MatchString(baseName):
		newName, err = renameXiaomiToDateTime(baseName, ext)
	case strings.HasPrefix(baseName, "Screenshot_"):
		newName, err = renameScreenshotToDateTime(baseName, ext)
	default:
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

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
		return
	}

	mu.Lock()
	*successCount++
	mu.Unlock()
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

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		mu.Lock()
		*errorCount++
		mu.Unlock()
		return
	}

	mu.Lock()
	*successCount++
	mu.Unlock()
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
