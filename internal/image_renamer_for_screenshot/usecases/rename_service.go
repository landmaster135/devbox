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

type renameTask struct {
	oldPath string
	newPath string
}

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

func resolveScreenshotRenameTarget(file FileInfo, config Config) (string, bool, error) {
	oldPath := file.Path
	oldName := file.Name
	if oldName == "" {
		oldName = filepath.Base(oldPath)
	}
	ext := filepath.Ext(oldName)
	baseName := strings.TrimSuffix(oldName, ext)

	var newName string
	var err error

	if config.ToDateTime {
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
			return "", false, nil
		}
	} else {
		switch config.Operation {
		case OperationVLC:
			if !strings.HasPrefix(baseName, "vlcsnap-") {
				return "", false, nil
			}
			newName, err = renameVlcScreenshot(baseName, ext)
		case OperationWin:
			if !strings.HasPrefix(baseName, "スクリーンショット ") {
				return "", false, nil
			}
			newName, err = renameWindowsScreenshot(baseName, ext)
		case OperationPixel:
			if !strings.HasPrefix(baseName, "screen-") {
				return "", false, nil
			}
			newName, err = renamePixelScreenshot(baseName, ext)
		case OperationXiaomi:
			if !xiaomiScreenshotRegexp.MatchString(baseName) {
				return "", false, nil
			}
			newName, err = renameXiaomiScreenshot(baseName, ext)
		default:
			return "", false, fmt.Errorf("未対応のoperation: %v", config.Operation)
		}
	}

	if err != nil {
		return "", false, err
	}

	return filepath.Join(filepath.Dir(oldPath), newName), true, nil
}

func performScreenshotRename(oldPath, newPath string, stdout, stderr io.Writer) error {
	if oldPath == newPath {
		fmt.Fprintf(stdout, "スキップ: %s は既に目的のファイル名です。\n", oldPath)
		return nil
	}

	fmt.Fprintf(stdout, "処理中: %s -> %s\n", oldPath, newPath)
	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(stderr, "エラー: %s のリネームに失敗しました: %v\n", oldPath, err)
		return err
	}

	return nil
}

// renameScreenshotFiles はスクリーンショットファイルをリネームします
func renameScreenshotFiles(fileInfos []FileInfo, config Config, stdout, stderr io.Writer) (int, int) {
	if len(fileInfos) == 0 {
		return 0, 0
	}

	renamedTasks, conflicts, precheckErrors := buildScreenshotRenamePlan(fileInfos, config, stderr)
	if len(conflicts) > 0 {
		fmt.Fprintln(stderr, "エラー: リネーム予定のファイル名が既存のファイル名と衝突しています。以下を確認してください:")
		for _, msg := range conflicts {
			fmt.Fprintf(stderr, "  - %s\n", msg)
		}
		return 0, precheckErrors + len(conflicts)
	}

	if len(renamedTasks) == 0 {
		return 0, precheckErrors
	}

	// ワーカープールの設定
	workerCount := config.Workers
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > len(fileInfos) {
		workerCount = len(fileInfos)
	}

	fmt.Fprintf(stdout, "リネーム操作に %d ワーカーを使用します。\n", workerCount)

	jobChan := make(chan renameTask, len(renamedTasks))
	var mu sync.Mutex
	var wg sync.WaitGroup
	successCount := 0
	errorCount := precheckErrors

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobChan {
				if err := performScreenshotRename(task.oldPath, task.newPath, stdout, stderr); err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}

				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	for _, task := range renamedTasks {
		jobChan <- task
	}
	close(jobChan)
	wg.Wait()

	return successCount, errorCount
}

func buildScreenshotRenamePlan(fileInfos []FileInfo, config Config, stderr io.Writer) ([]renameTask, []string, int) {
	plannedPaths := make(map[string]string, len(fileInfos))
	renamedTasks := make([]renameTask, 0, len(fileInfos))
	conflicts := make([]string, 0)
	precheckErrors := 0

	for _, file := range fileInfos {
		newPath, shouldRename, err := resolveScreenshotRenameTarget(file, config)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %s の解析に失敗しました: %v\n", file.Path, err)
			precheckErrors++
			continue
		}

		if !shouldRename {
			continue
		}

		if newPath != file.Path {
			if _, err := os.Stat(newPath); err == nil {
				conflicts = append(conflicts, fmt.Sprintf("%s -> %s (既存ファイルが存在します)", file.Path, newPath))
				continue
			} else if !os.IsNotExist(err) {
				fmt.Fprintf(stderr, "エラー: %s の存在確認に失敗しました: %v\n", newPath, err)
				precheckErrors++
				continue
			}
		}

		if owner, exists := plannedPaths[newPath]; exists && owner != file.Path {
			conflicts = append(conflicts, fmt.Sprintf("%s と %s が同じリネーム先 %s を要求しています。", file.Path, owner, newPath))
			continue
		}

		plannedPaths[newPath] = file.Path
		renamedTasks = append(renamedTasks, renameTask{oldPath: file.Path, newPath: newPath})
	}

	return renamedTasks, conflicts, precheckErrors
}
