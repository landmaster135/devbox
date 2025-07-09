package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestImageRenamerService_ProcessImageRename_Normal は全体的な処理フローの正常系テストです
func TestImageRenamerService_ProcessImageRename_Normal(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用ファイルを作成
	testFiles := []string{
		"test1.jpg",
		"test2.jpeg",
		"test3.png",
	}

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()

		// ファイルの更新時刻を設定（異なる時刻にして競合を避ける）
		testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)
		if fileName == "test2.jpeg" {
			testTime = testTime.Add(1 * time.Second)
		} else if fileName == "test3.png" {
			testTime = testTime.Add(2 * time.Second)
		}
		err = os.Chtimes(filePath, testTime, testTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}
	}

	// テストケース1: ドライラン
	config := &Config{
		FolderPath:     tempDir,
		Extension:      "",
		Recursive:      false,
		DryRun:         true,
		Verbose:        true,
		WorkerCount:    2,
		UseFileModTime: true,
	}

	processedCount, errorCount, err := service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	expectedProcessedCount := 3
	if processedCount != expectedProcessedCount {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}

	// ドライランなので元のファイルが存在することを確認
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("ドライランなのにファイルが削除されました: %s", fileName)
		}
	}

	// テストケース2: 実際のリネーム
	config.DryRun = false
	processedCount, errorCount, err = service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	if processedCount != expectedProcessedCount {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}

	// リネーム後のファイルが存在することを確認
	expectedNewFiles := []string{
		"20250629130139.jpg",
		"20250629130140.jpeg",
		"20250629130141.png",
	}

	for _, fileName := range expectedNewFiles {
		filePath := filepath.Join(tempDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("リネーム後のファイルが存在しません: %s", fileName)
		}
	}

	// 元のファイルが存在しないことを確認
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("リネーム後に元のファイルが残っています: %s", fileName)
		}
	}
}

// TestImageRenamerService_ProcessImageRename_NoFiles はファイルが存在しない場合のテストです
func TestImageRenamerService_ProcessImageRename_NoFiles(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成（空のディレクトリ）
	tempDir := t.TempDir()

	config := &Config{
		FolderPath:     tempDir,
		Extension:      "",
		Recursive:      false,
		DryRun:         false,
		Verbose:        true,
		WorkerCount:    2,
		UseFileModTime: true,
	}

	processedCount, errorCount, err := service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	if processedCount != 0 {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 0, 実際の値 %d", processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}
}

// TestImageRenamerService_ProcessImageRename_WithConflicts は競合があるファイルの処理テストです
func TestImageRenamerService_ProcessImageRename_WithConflicts(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 同じ時刻のファイルを複数作成（競合を発生させる）
	testFiles := []string{
		"conflict1.jpg",
		"conflict2.jpg",
		"conflict3.jpg",
	}

	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)

	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()

		// 全て同じ時刻に設定（競合を発生させる）
		err = os.Chtimes(filePath, testTime, testTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}
	}

	config := &Config{
		FolderPath:     tempDir,
		Extension:      "",
		Recursive:      false,
		DryRun:         false,
		Verbose:        true,
		WorkerCount:    2,
		UseFileModTime: true,
	}

	processedCount, errorCount, err := service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	expectedProcessedCount := 3
	if processedCount != expectedProcessedCount {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}

	// 競合解決後のファイルが存在することを確認
	expectedNewFiles := []string{
		"20250629130139.jpg", // 最初のファイル
		"20250629130140.jpg", // 1秒後
		"20250629130141.jpg", // 2秒後
	}

	for _, fileName := range expectedNewFiles {
		filePath := filepath.Join(tempDir, fileName)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("競合解決後のファイルが存在しません: %s", fileName)
		}
	}

	// 元のファイルが存在しないことを確認
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("リネーム後に元のファイルが残っています: %s", fileName)
		}
	}
}

// TestImageRenamerService_ProcessImageRename_RecursiveSearch は再帰検索のテストです
func TestImageRenamerService_ProcessImageRename_RecursiveSearch(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// サブディレクトリを作成
	subDir1 := filepath.Join(tempDir, "subdir1")
	subDir2 := filepath.Join(tempDir, "subdir2")
	err := os.Mkdir(subDir1, 0755)
	if err != nil {
		t.Fatalf("サブディレクトリの作成に失敗: %v", err)
	}
	err = os.Mkdir(subDir2, 0755)
	if err != nil {
		t.Fatalf("サブディレクトリの作成に失敗: %v", err)
	}

	// 各ディレクトリにファイルを作成
	testFiles := map[string]string{
		"root.jpg":           tempDir,
		"subdir1/sub1.jpg":   tempDir,
		"subdir2/sub2.jpeg":  tempDir,
	}

	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)
	fileIndex := 0

	for relativePath, baseDir := range testFiles {
		filePath := filepath.Join(baseDir, relativePath)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", relativePath, err)
		}
		file.Close()

		// 異なる時刻を設定
		fileTime := testTime.Add(time.Duration(fileIndex) * time.Second)
		err = os.Chtimes(filePath, fileTime, fileTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}
		fileIndex++
	}

	// テストケース1: 再帰検索なし
	config := &Config{
		FolderPath:     tempDir,
		Extension:      "",
		Recursive:      false,
		DryRun:         true,
		Verbose:        true,
		WorkerCount:    2,
		UseFileModTime: true,
	}

	processedCount, _, err := service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	// 再帰検索なしの場合、ルートディレクトリのファイルのみ処理される
	expectedProcessedCount := 1
	if processedCount != expectedProcessedCount {
		t.Errorf("再帰検索なしの処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	// テストケース2: 再帰検索あり
	config.Recursive = true
	var errorCount int
	processedCount, errorCount, err = service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	// 再帰検索ありの場合、全てのファイルが処理される
	expectedProcessedCount = 3
	if processedCount != expectedProcessedCount {
		t.Errorf("再帰検索ありの処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}
}

// TestImageRenamerService_ProcessImageRename_SpecificExtension は特定拡張子のテストです
func TestImageRenamerService_ProcessImageRename_SpecificExtension(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 異なる拡張子のファイルを作成
	testFiles := []string{
		"test1.jpg",
		"test2.jpeg",
		"test3.png",
		"test4.tiff",
	}

	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)

	for i, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()

		// 異なる時刻を設定
		fileTime := testTime.Add(time.Duration(i) * time.Second)
		err = os.Chtimes(filePath, fileTime, fileTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}
	}

	// テストケース: jpg拡張子のみ処理
	config := &Config{
		FolderPath:     tempDir,
		Extension:      "jpg",
		Recursive:      false,
		DryRun:         true,
		Verbose:        true,
		WorkerCount:    2,
		UseFileModTime: true,
	}

	processedCount, errorCount, err := service.ProcessImageRename(config)
	if err != nil {
		t.Fatalf("ProcessImageRename failed: %v", err)
	}

	// jpg拡張子のファイルのみ処理される
	expectedProcessedCount := 1
	if processedCount != expectedProcessedCount {
		t.Errorf("jpg拡張子指定の処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}
}

// TestImageRenamerService_RenameImageFilesWithInfo_ParallelProcessing は並行処理のテストです
func TestImageRenamerService_RenameImageFilesWithInfo_ParallelProcessing(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 多数のファイルを作成（並行処理をテストするため）
	fileCount := 10
	var renameInfos []FileRenameInfo
	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)

	for i := 0; i < fileCount; i++ {
		fileName := fmt.Sprintf("test%d.jpg", i)
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()

		// 異なる時刻を設定（競合を避ける）
		fileTime := testTime.Add(time.Duration(i) * time.Second)
		err = os.Chtimes(filePath, fileTime, fileTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}

		// RenameInfoを作成
		renameInfo := FileRenameInfo{
			OriginalPath: filePath,
			NewFileName:  fileTime.Format("20060102150405") + ".jpg",
			CreateDate:   fileTime,
			Directory:    tempDir,
		}
		renameInfos = append(renameInfos, renameInfo)
	}

	// テストケース1: 並行処理数を指定
	config := &Config{
		DryRun:      false,
		Verbose:     true,
		WorkerCount: 3, // 3並列で処理
	}

	processedCount, errorCount, err := service.renameImageFilesWithInfo(renameInfos, config)
	if err != nil {
		t.Fatalf("renameImageFilesWithInfo failed: %v", err)
	}

	if processedCount != fileCount {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", fileCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}

	// リネーム後のファイルが存在することを確認
	for _, info := range renameInfos {
		newPath := filepath.Join(info.Directory, info.NewFileName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Errorf("リネーム後のファイルが存在しません: %s", newPath)
		}
	}
}

// TestImageRenamerService_RenameImageFilesWithInfo_WorkerCountAdjustment はワーカー数調整のテストです
func TestImageRenamerService_RenameImageFilesWithInfo_WorkerCountAdjustment(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 少数のファイルを作成
	fileCount := 2
	var renameInfos []FileRenameInfo
	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)

	for i := 0; i < fileCount; i++ {
		fileName := fmt.Sprintf("test%d.jpg", i)
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()

		// 異なる時刻を設定
		fileTime := testTime.Add(time.Duration(i) * time.Second)
		err = os.Chtimes(filePath, fileTime, fileTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}

		// RenameInfoを作成
		renameInfo := FileRenameInfo{
			OriginalPath: filePath,
			NewFileName:  fileTime.Format("20060102150405") + ".jpg",
			CreateDate:   fileTime,
			Directory:    tempDir,
		}
		renameInfos = append(renameInfos, renameInfo)
	}

	// テストケース1: ワーカー数がファイル数より多い場合
	config := &Config{
		DryRun:      false,
		Verbose:     true,
		WorkerCount: 5, // ファイル数(2)より多い
	}

	processedCount, errorCount, err := service.renameImageFilesWithInfo(renameInfos, config)
	if err != nil {
		t.Fatalf("renameImageFilesWithInfo failed: %v", err)
	}

	if processedCount != fileCount {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", fileCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}

	// テストケース2: ワーカー数が0以下の場合（デフォルト値が使用される）
	// 新しいファイルを作成（新しいrenameInfosスライスを作成）
	var newRenameInfos []FileRenameInfo
	for i := 0; i < fileCount; i++ {
		fileName := fmt.Sprintf("test_default%d.jpg", i)
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()

		// 異なる時刻を設定
		fileTime := testTime.Add(time.Duration(i+10) * time.Second)
		err = os.Chtimes(filePath, fileTime, fileTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}

		// RenameInfoを作成
		renameInfo := FileRenameInfo{
			OriginalPath: filePath,
			NewFileName:  fileTime.Format("20060102150405") + ".jpg",
			CreateDate:   fileTime,
			Directory:    tempDir,
		}
		newRenameInfos = append(newRenameInfos, renameInfo)
	}

	config.WorkerCount = 0 // デフォルト値が使用される

	processedCount, errorCount, err = service.renameImageFilesWithInfo(newRenameInfos, config)
	if err != nil {
		t.Fatalf("renameImageFilesWithInfo failed: %v", err)
	}

	expectedProcessedCount := len(newRenameInfos)
	if processedCount != expectedProcessedCount {
		t.Errorf("処理されたファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedProcessedCount, processedCount)
	}

	if errorCount != 0 {
		t.Errorf("エラー数が期待値と異なります: 期待値 0, 実際の値 %d", errorCount)
	}
}
