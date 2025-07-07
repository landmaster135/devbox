package usecases

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolveConflicts_MultipleConflictGroups(t *testing.T) {
	service := NewImageRenamerService()

	// テストケース: 複数の競合グループが存在する場合
	// グループ1: 20250629130139.jpg に3つのファイルが競合
	// グループ2: 20250629130138.jpg に3つのファイルが競合

	// 基準時刻を設定
	baseTime1, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:39")
	baseTime2, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:38")

	renameInfos := []FileRenameInfo{
		// グループ1: 20250629130139.jpg
		{
			OriginalPath: "1-1_image_renamer_with_exif/PXL_20250629_040139164.jpg",
			NewFileName:  "20250629130139.jpg",
			CreateDate:   baseTime1,
			Directory:    "1-1_image_renamer_with_exif",
		},
		{
			OriginalPath: "1-1_image_renamer_with_exif/PXL_20250629_040139378.jpg",
			NewFileName:  "20250629130139.jpg",
			CreateDate:   baseTime1,
			Directory:    "1-1_image_renamer_with_exif",
		},
		{
			OriginalPath: "1-1_image_renamer_with_exif/PXL_20250629_040139801.jpg",
			NewFileName:  "20250629130139.jpg",
			CreateDate:   baseTime1,
			Directory:    "1-1_image_renamer_with_exif",
		},
		// グループ2: 20250629130138.jpg
		{
			OriginalPath: "1-1_image_renamer_with_exif/PXL_20250629_040138106.jpg",
			NewFileName:  "20250629130138.jpg",
			CreateDate:   baseTime2,
			Directory:    "1-1_image_renamer_with_exif",
		},
		{
			OriginalPath: "1-1_image_renamer_with_exif/PXL_20250629_040138353.jpg",
			NewFileName:  "20250629130138.jpg",
			CreateDate:   baseTime2,
			Directory:    "1-1_image_renamer_with_exif",
		},
		{
			OriginalPath: "1-1_image_renamer_with_exif/PXL_20250629_040138704.jpg",
			NewFileName:  "20250629130138.jpg",
			CreateDate:   baseTime2,
			Directory:    "1-1_image_renamer_with_exif",
		},
	}

	config := &Config{
		Verbose: true,
	}

	// 競合解決を実行
	err := service.resolveConflicts(renameInfos, config)
	if err != nil {
		t.Fatalf("resolveConflicts failed: %v", err)
	}

	// 結果を検証
	t.Logf("=== 競合解決後の結果 ===")
	for i, info := range renameInfos {
		t.Logf("[%d] %s → %s", i, info.OriginalPath, info.NewFileName)
	}

	// 具体的なNewFileNameの値を検証
	// 決定的な処理順序: キーをソートして処理するため、20250629130138.jpgが先に処理される
	expectedResults := map[string]string{
		"1-1_image_renamer_with_exif/PXL_20250629_040138106.jpg": "20250629130138.jpg", // グループ1: 最初のファイルは元の時刻維持
		"1-1_image_renamer_with_exif/PXL_20250629_040138353.jpg": "20250629130140.jpg", // グループ1: 2秒後
		"1-1_image_renamer_with_exif/PXL_20250629_040138704.jpg": "20250629130141.jpg", // グループ1: 3秒後
		"1-1_image_renamer_with_exif/PXL_20250629_040139164.jpg": "20250629130139.jpg", // グループ2: 最初のファイルは元の時刻維持
		"1-1_image_renamer_with_exif/PXL_20250629_040139378.jpg": "20250629130142.jpg", // グループ2: 3秒後（130140,130141は使用済み）
		"1-1_image_renamer_with_exif/PXL_20250629_040139801.jpg": "20250629130143.jpg", // グループ2: 4秒後
	}

	// NewFileNameの値を検証
	for _, info := range renameInfos {
		expected, exists := expectedResults[info.OriginalPath]
		if !exists {
			t.Errorf("予期しないファイル: %s", info.OriginalPath)
			continue
		}

		if info.NewFileName != expected {
			t.Errorf("ファイル %s: 期待値 %s, 実際の値 %s",
				info.OriginalPath, expected, info.NewFileName)
		}
	}

	// 全てのファイル名がユニークかチェック
	fileNameMap := make(map[string][]string)
	for _, info := range renameInfos {
		fullPath := info.Directory + "/" + info.NewFileName
		fileNameMap[fullPath] = append(fileNameMap[fullPath], info.OriginalPath)
	}

	// 重複チェック
	hasConflict := false
	for newPath, originalPaths := range fileNameMap {
		if len(originalPaths) > 1 {
			t.Errorf("競合が解決されていません: %s に以下のファイルが競合:", newPath)
			for _, originalPath := range originalPaths {
				t.Errorf("  - %s", originalPath)
			}
			hasConflict = true
		}
	}

	if !hasConflict {
		t.Log("✅ 全ての競合が正常に解決されました")
	}
}

func TestResolveConflicts_SingleConflictGroup(t *testing.T) {
	service := NewImageRenamerService()

	// シンプルなケース: 1つの競合グループのみ
	baseTime, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:39")

	// NewFileNameの具体的な値を検証
	expectedFile1 := "20250629130139.jpg" // 最初のファイルは元の時刻維持
	expectedFile2 := "20250629130140.jpg" // 1秒後

	renameInfos := []FileRenameInfo{
		{
			OriginalPath: "test/file1.jpg",
			NewFileName:  expectedFile1,
			CreateDate:   baseTime,
			Directory:    "test",
		},
		{
			OriginalPath: "test/file2.jpg",
			NewFileName:  expectedFile1,
			CreateDate:   baseTime,
			Directory:    "test",
		},
	}

	config := &Config{
		Verbose: true,
	}

	err := service.resolveConflicts(renameInfos, config)
	if err != nil {
		t.Fatalf("resolveConflicts failed: %v", err)
	}

	if renameInfos[0].NewFileName != expectedFile1 {
		t.Errorf("file1.jpg: 期待値 %s, 実際の値 %s", expectedFile1, renameInfos[0].NewFileName)
	}

	if renameInfos[1].NewFileName != expectedFile2 {
		t.Errorf("file2.jpg: 期待値 %s, 実際の値 %s", expectedFile2, renameInfos[1].NewFileName)
	}

	// 競合が解決されているかチェック
	if renameInfos[0].NewFileName == renameInfos[1].NewFileName {
		t.Errorf("競合が解決されていません: 両方とも %s", renameInfos[0].NewFileName)
	}

	t.Logf("file1.jpg → %s", renameInfos[0].NewFileName)
	t.Logf("file2.jpg → %s", renameInfos[1].NewFileName)
}

func TestResolveConflicts_NoConflicts(t *testing.T) {
	service := NewImageRenamerService()

	// 競合がないケース
	baseTime1, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:39")
	baseTime2, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:40")

	// NewFileNameの具体的な値を検証
	expectedFile1 := "20250629130139.jpg" // 最初のファイルは元の時刻維持
	expectedFile2 := "20250629130140.jpg" // 1秒後

	renameInfos := []FileRenameInfo{
		{
			OriginalPath: "test/file1.jpg",
			NewFileName:  expectedFile1,
			CreateDate:   baseTime1,
			Directory:    "test",
		},
		{
			OriginalPath: "test/file2.jpg",
			NewFileName:  expectedFile2,
			CreateDate:   baseTime2,
			Directory:    "test",
		},
	}

	config := &Config{
		Verbose: true,
	}

	err := service.resolveConflicts(renameInfos, config)
	if err != nil {
		t.Fatalf("resolveConflicts failed: %v", err)
	}

	// NewFileNameが変更されていないことを確認
	if renameInfos[0].NewFileName != expectedFile1 {
		t.Errorf("file1.jpg: NewFileNameが変更されました: %s", renameInfos[0].NewFileName)
	}

	if renameInfos[1].NewFileName != expectedFile2 {
		t.Errorf("file2.jpg: NewFileNameが変更されました: %s", renameInfos[1].NewFileName)
	}

	t.Log("✅ 競合がない場合、NewFileNameは変更されませんでした")
}

func TestResolveConflicts_MixedConflictAndNonConflict(t *testing.T) {
	service := NewImageRenamerService()

	// 混合ケース: 競合ファイル + 非競合ファイル（ただし非競合ファイルが競合解決後のファイル名と重複）
	baseTime1, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:39")
	baseTime2, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:43") // 143秒 = 競合解決後に使用される時刻

	renameInfos := []FileRenameInfo{
		// 競合グループ: 20250629130139.jpg
		{
			OriginalPath: "test/file1.jpg",
			NewFileName:  "20250629130139.jpg",
			CreateDate:   baseTime1,
			Directory:    "test",
		},
		{
			OriginalPath: "test/file2.jpg",
			NewFileName:  "20250629130139.jpg",
			CreateDate:   baseTime1,
			Directory:    "test",
		},
		{
			OriginalPath: "test/file3.jpg",
			NewFileName:  "20250629130139.jpg",
			CreateDate:   baseTime1,
			Directory:    "test",
		},
		// 非競合ファイル（競合解決で使用される時刻と重複）
		{
			OriginalPath: "test/file4.jpg",
			NewFileName:  "20250629130141.jpg", // 競合解決で使用される時刻と重複
			CreateDate:   baseTime2,
			Directory:    "test",
		},
	}

	config := &Config{
		Verbose: true,
	}

	err := service.resolveConflicts(renameInfos, config)
	if err != nil {
		t.Fatalf("resolveConflicts failed: %v", err)
	}

	// 結果を検証
	t.Logf("=== 混合ケースの競合解決後の結果 ===")
	for i, info := range renameInfos {
		t.Logf("[%d] %s → %s", i, info.OriginalPath, info.NewFileName)
	}

	// 期待される結果
	expectedResults := map[string]string{
		"test/file1.jpg": "20250629130139.jpg", // 最初のファイルは元の時刻維持
		"test/file2.jpg": "20250629130140.jpg", // 1秒後
		"test/file3.jpg": "20250629130142.jpg", // 3秒後（130141は非競合ファイルで使用済み）
		"test/file4.jpg": "20250629130141.jpg", // 非競合ファイルは変更されない
	}

	// NewFileNameの値を検証
	for _, info := range renameInfos {
		expected, exists := expectedResults[info.OriginalPath]
		if !exists {
			t.Errorf("予期しないファイル: %s", info.OriginalPath)
			continue
		}

		if info.NewFileName != expected {
			t.Errorf("ファイル %s: 期待値 %s, 実際の値 %s",
				info.OriginalPath, expected, info.NewFileName)
		}
	}

	// 全てのファイル名がユニークかチェック
	fileNameMap := make(map[string][]string)
	for _, info := range renameInfos {
		fullPath := info.Directory + "/" + info.NewFileName
		fileNameMap[fullPath] = append(fileNameMap[fullPath], info.OriginalPath)
	}

	// 重複チェック
	hasConflict := false
	for newPath, originalPaths := range fileNameMap {
		if len(originalPaths) > 1 {
			t.Errorf("競合が解決されていません: %s に以下のファイルが競合:", newPath)
			for _, originalPath := range originalPaths {
				t.Errorf("  - %s", originalPath)
			}
			hasConflict = true
		}
	}

	if !hasConflict {
		t.Log("✅ 混合ケースでも全ての競合が正常に解決されました")
	}
}

// TestValidateInputOptions_Normal はValidateInputOptionsの正常系テストです
func TestValidateInputOptions_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// 正常ケース: 有効なディレクトリと拡張子
	err := ValidateInputOptions(tempDir, "jpg")
	if err != nil {
		t.Errorf("正常なケースでエラーが発生: %v", err)
	}

	// 正常ケース: 拡張子が空文字列
	err = ValidateInputOptions(tempDir, "")
	if err != nil {
		t.Errorf("拡張子が空文字列の正常ケースでエラーが発生: %v", err)
	}

	// 正常ケース: ドット付き拡張子
	err = ValidateInputOptions(tempDir, ".jpeg")
	if err != nil {
		t.Errorf("ドット付き拡張子の正常ケースでエラーが発生: %v", err)
	}
}

// TestValidateInputOptions_InvalidDirectory はValidateInputOptionsのディレクトリ異常系テストです
func TestValidateInputOptions_InvalidDirectory(t *testing.T) {
	// 異常ケース: 空文字列のディレクトリ
	err := ValidateInputOptions("", "jpg")
	if err == nil {
		t.Error("空文字列のディレクトリでエラーが発生しませんでした")
	}

	// 異常ケース: 存在しないディレクトリ
	err = ValidateInputOptions("/nonexistent/directory", "jpg")
	if err == nil {
		t.Error("存在しないディレクトリでエラーが発生しませんでした")
	}

	// 異常ケース: ファイルをディレクトリとして指定
	tempFile := filepath.Join(t.TempDir(), "testfile.txt")
	file, _ := os.Create(tempFile)
	file.Close()

	err = ValidateInputOptions(tempFile, "jpg")
	if err == nil {
		t.Error("ファイルをディレクトリとして指定してもエラーが発生しませんでした")
	}
}

// TestValidateInputOptions_InvalidExtension はValidateInputOptionsの拡張子異常系テストです
func TestValidateInputOptions_InvalidExtension(t *testing.T) {
	tempDir := t.TempDir()

	// 異常ケース: サポートされていない拡張子
	err := ValidateInputOptions(tempDir, "txt")
	if err == nil {
		t.Error("サポートされていない拡張子でエラーが発生しませんでした")
	}

	// 異常ケース: サポートされていない拡張子（ドット付き）
	err = ValidateInputOptions(tempDir, ".doc")
	if err == nil {
		t.Error("サポートされていない拡張子（ドット付き）でエラーが発生しませんでした")
	}
}

// TestImageRenamerService_IsImageFile はisImageFileメソッドのテストです
func TestImageRenamerService_IsImageFile(t *testing.T) {
	service := NewImageRenamerService()

	// 正常ケース: サポートされている拡張子
	testCases := []struct {
		filePath        string
		targetExtension string
		expected        bool
		description     string
	}{
		{"test.jpg", "", true, "JPGファイル（拡張子指定なし）"},
		{"test.jpeg", "", true, "JPEGファイル（拡張子指定なし）"},
		{"test.png", "", true, "PNGファイル（拡張子指定なし）"},
		{"test.tiff", "", true, "TIFFファイル（拡張子指定なし）"},
		{"test.webp", "", true, "WebPファイル（拡張子指定なし）"},
		{"test.mp4", "", true, "MP4ファイル（拡張子指定なし）"},
		{"test.txt", "", false, "テキストファイル（拡張子指定なし）"},
		{"test.jpg", "jpg", true, "JPGファイル（jpg指定）"},
		{"test.jpeg", "jpg", false, "JPEGファイル（jpg指定）"},
		{"test.jpg", ".jpg", true, "JPGファイル（.jpg指定）"},
		{"test.JPG", "", true, "大文字JPGファイル"},
		{"test.JPG", "jpg", true, "大文字JPGファイル（jpg指定）"},
	}

	for _, tc := range testCases {
		result := service.isImageFile(tc.filePath, tc.targetExtension)
		if result != tc.expected {
			t.Errorf("%s: 期待値 %v, 実際の値 %v", tc.description, tc.expected, result)
		}
	}
}

// TestImageRenamerService_GenerateNewFileName はgenerateNewFileNameメソッドのテストです
func TestImageRenamerService_GenerateNewFileName(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の時刻
	testTime, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:39")

	testCases := []struct {
		originalPath string
		expected     string
		description  string
	}{
		{"test.jpg", "20250629130139.jpg", "JPGファイル"},
		{"test.jpeg", "20250629130139.jpeg", "JPEGファイル"},
		{"test.png", "20250629130139.png", "PNGファイル"},
		{"path/to/test.tiff", "20250629130139.tiff", "TIFFファイル（パス付き）"},
	}

	for _, tc := range testCases {
		result := service.generateNewFileName(testTime, tc.originalPath)
		if result != tc.expected {
			t.Errorf("%s: 期待値 %s, 実際の値 %s", tc.description, tc.expected, result)
		}
	}
}

// TestNewImageRenamerService はNewImageRenamerServiceのテストです
func TestNewImageRenamerService(t *testing.T) {
	service := NewImageRenamerService()
	if service == nil {
		t.Error("NewImageRenamerServiceがnilを返しました")
	}
}

// TestNewConflictResolver はNewConflictResolverのテストです
func TestNewConflictResolver(t *testing.T) {
	existingFiles := map[string]bool{
		"test1.jpg": true,
		"test2.jpg": true,
	}

	resolver := NewConflictResolver(existingFiles)

	if len(resolver.usedFileNames) != 2 {
		t.Errorf("usedFileNamesの長さが期待値と異なります: 期待値 2, 実際の値 %d", len(resolver.usedFileNames))
	}
}

// TestConflictResolver_FindNextAvailableTime はfindNextAvailableTimeメソッドのテストです
func TestConflictResolver_FindNextAvailableTime(t *testing.T) {
	// 既存ファイルを設定
	existingFiles := map[string]bool{
		"/test/20250629130139.jpg": true,
		"/test/20250629130140.jpg": true,
	}

	resolver := NewConflictResolver(existingFiles)
	baseTime, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:39")

	// 利用可能な次の時刻を検索
	nextTime := resolver.findNextAvailableTime(baseTime, ".jpg", "/test")

	// 期待される時刻（20250629130141）
	expectedTime, _ := time.Parse("2006:01:02 15:04:05", "2025:06:29 13:01:41")

	if !nextTime.Equal(expectedTime) {
		t.Errorf("期待される時刻と異なります: 期待値 %v, 実際の値 %v", expectedTime, nextTime)
	}

	// 使用済みファイル名に追加されているかチェック
	expectedPath := "/test/20250629130141.jpg"
	if !resolver.usedFileNames[expectedPath] {
		t.Errorf("使用済みファイル名に追加されていません: %s", expectedPath)
	}
}

// TestImageRenamerService_FindImageFiles はfindImageFilesメソッドのテストです
func TestImageRenamerService_FindImageFiles(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用ファイルを作成
	testFiles := []string{
		"test1.jpg",
		"test2.jpeg",
		"test3.png",
		"test4.txt", // サポートされていない拡張子
		"subdir/test5.jpg",
	}

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("サブディレクトリの作成に失敗: %v", err)
	}

	// ファイルを作成
	for _, fileName := range testFiles {
		filePath := filepath.Join(tempDir, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %s - %v", fileName, err)
		}
		file.Close()
	}

	// テストケース1: 再帰検索なし、拡張子指定なし
	config := &Config{
		FolderPath: tempDir,
		Extension:  "",
		Recursive:  false,
	}

	imageFiles, err := service.findImageFiles(config)
	if err != nil {
		t.Fatalf("findImageFiles failed: %v", err)
	}

	// 期待される結果: test1.jpg, test2.jpeg, test3.png（サブディレクトリは除外）
	expectedCount := 3
	if len(imageFiles) != expectedCount {
		t.Errorf("再帰検索なしの場合のファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedCount, len(imageFiles))
	}

	// テストケース2: 再帰検索あり、拡張子指定なし
	config.Recursive = true
	imageFiles, err = service.findImageFiles(config)
	if err != nil {
		t.Fatalf("findImageFiles failed: %v", err)
	}

	// 期待される結果: test1.jpg, test2.jpeg, test3.png, subdir/test5.jpg
	expectedCount = 4
	if len(imageFiles) != expectedCount {
		t.Errorf("再帰検索ありの場合のファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedCount, len(imageFiles))
	}

	// テストケース3: 特定の拡張子のみ
	config.Extension = "jpg"
	imageFiles, err = service.findImageFiles(config)
	if err != nil {
		t.Fatalf("findImageFiles failed: %v", err)
	}

	// 期待される結果: test1.jpg, subdir/test5.jpg
	expectedCount = 2
	if len(imageFiles) != expectedCount {
		t.Errorf("jpg拡張子指定の場合のファイル数が期待値と異なります: 期待値 %d, 実際の値 %d", expectedCount, len(imageFiles))
	}

	// ファイルがソートされているかチェック
	for i := 1; i < len(imageFiles); i++ {
		if filepath.Base(imageFiles[i-1]) > filepath.Base(imageFiles[i]) {
			t.Errorf("ファイルがソートされていません: %s > %s", filepath.Base(imageFiles[i-1]), filepath.Base(imageFiles[i]))
		}
	}
}

// TestImageRenamerService_PrepareRenameInfo はprepareRenameInfoメソッドのテストです
func TestImageRenamerService_PrepareRenameInfo(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用ファイルを作成
	testFile := filepath.Join(tempDir, "test.jpg")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	file.Close()

	// ファイルの更新時刻を設定
	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)
	err = os.Chtimes(testFile, testTime, testTime)
	if err != nil {
		t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
	}

	imageFiles := []string{testFile}

	// テストケース1: ファイルの更新時刻を使用
	config := &Config{
		UseFileModTime: true,
		Verbose:        true,
	}

	renameInfos, err := service.prepareRenameInfo(imageFiles, config)
	if err != nil {
		t.Fatalf("prepareRenameInfo failed: %v", err)
	}

	if len(renameInfos) != 1 {
		t.Fatalf("renameInfosの長さが期待値と異なります: 期待値 1, 実際の値 %d", len(renameInfos))
	}

	info := renameInfos[0]
	if info.OriginalPath != testFile {
		t.Errorf("OriginalPathが期待値と異なります: 期待値 %s, 実際の値 %s", testFile, info.OriginalPath)
	}

	expectedFileName := "20250629130139.jpg"
	if info.NewFileName != expectedFileName {
		t.Errorf("NewFileNameが期待値と異なります: 期待値 %s, 実際の値 %s", expectedFileName, info.NewFileName)
	}

	if info.Directory != tempDir {
		t.Errorf("Directoryが期待値と異なります: 期待値 %s, 実際の値 %s", tempDir, info.Directory)
	}

	// CreateDateがファイルの更新時刻と一致するかチェック
	if !info.CreateDate.Equal(testTime) {
		t.Errorf("CreateDateが期待値と異なります: 期待値 %v, 実際の値 %v", testTime, info.CreateDate)
	}
}

// TestImageRenamerService_RenameSingleFileWithInfo はrenameSingleFileWithInfoメソッドのテストです
func TestImageRenamerService_RenameSingleFileWithInfo(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// テスト用ファイルを作成
	originalFile := filepath.Join(tempDir, "original.jpg")
	file, err := os.Create(originalFile)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	file.Close()

	// テストケース1: ドライラン
	testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)
	renameInfo := FileRenameInfo{
		OriginalPath: originalFile,
		NewFileName:  "20250629130139.jpg",
		CreateDate:   testTime,
		Directory:    tempDir,
	}

	config := &Config{
		DryRun: true,
	}

	result := service.renameSingleFileWithInfo(renameInfo, config)

	if !result.Success {
		t.Errorf("ドライランが失敗しました: %v", result.Error)
	}

	if result.OriginalPath != originalFile {
		t.Errorf("OriginalPathが期待値と異なります: 期待値 %s, 実際の値 %s", originalFile, result.OriginalPath)
	}

	expectedNewPath := filepath.Join(tempDir, "20250629130139.jpg")
	if result.NewPath != expectedNewPath {
		t.Errorf("NewPathが期待値と異なります: 期待値 %s, 実際の値 %s", expectedNewPath, result.NewPath)
	}

	// ドライランなので元のファイルが存在することを確認
	if _, err := os.Stat(originalFile); os.IsNotExist(err) {
		t.Error("ドライランなのに元のファイルが削除されました")
	}

	// テストケース2: 実際のリネーム
	config.DryRun = false
	result = service.renameSingleFileWithInfo(renameInfo, config)

	if !result.Success {
		t.Errorf("実際のリネームが失敗しました: %v", result.Error)
	}

	// 元のファイルが存在しないことを確認
	if _, err := os.Stat(originalFile); !os.IsNotExist(err) {
		t.Error("リネーム後に元のファイルが残っています")
	}

	// 新しいファイルが存在することを確認
	if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
		t.Error("リネーム後に新しいファイルが存在しません")
	}
}

// TestImageRenamerService_ExtractCreateDate はextractCreateDateメソッドのテストです
func TestImageRenamerService_ExtractCreateDate(t *testing.T) {
	service := NewImageRenamerService()

	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()

	// PNG/TIFF/WebPファイルのテスト（ファイルのModTimeを使用）
	testCases := []struct {
		extension   string
		description string
	}{
		{".png", "PNGファイル"},
		{".tiff", "TIFFファイル"},
		{".webp", "WebPファイル"},
		{".mp4", "MP4ファイル"},
	}

	for _, tc := range testCases {
		testFile := filepath.Join(tempDir, "test"+tc.extension)
		file, err := os.Create(testFile)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗: %v", err)
		}
		file.Close()

		// ファイルの更新時刻を設定
		testTime := time.Date(2025, 6, 29, 13, 1, 39, 0, time.Local)
		err = os.Chtimes(testFile, testTime, testTime)
		if err != nil {
			t.Fatalf("ファイルの更新時刻設定に失敗: %v", err)
		}

		// CreateDateを抽出
		createDate, err := service.extractCreateDate(testFile)
		if err != nil {
			t.Errorf("%s: extractCreateDate failed: %v", tc.description, err)
			continue
		}

		// ModTimeと一致するかチェック
		if !createDate.Equal(testTime) {
			t.Errorf("%s: CreateDateが期待値と異なります: 期待値 %v, 実際の値 %v", tc.description, testTime, createDate)
		}
	}

	// サポートされていない拡張子のテスト
	unsupportedFile := filepath.Join(tempDir, "test.txt")
	file, err := os.Create(unsupportedFile)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗: %v", err)
	}
	file.Close()

	_, err = service.extractCreateDate(unsupportedFile)
	if err == nil {
		t.Error("サポートされていない拡張子でエラーが発生しませんでした")
	}
}
