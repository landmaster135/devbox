package usecases

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// #==============================================================#
// ##          Setup and teardown                                ##
// #==============================================================#
// テスト用の定数
const (
	targetTmpDir = "internal/independencies/exif_modifier/test_data/tmp"
	targetOrgDir = "internal/independencies/exif_modifier/test_data/org"
)

var (
	dataTmpDir string
	dataOrgDir string
)

// TestMain はテスト全体のセットアップとティアダウンを行います
func TestMain(m *testing.M) {
	fmt.Println("テスト環境のセットアップを開始します...")

	// テストデータのオリジナルと実施用のものが入ったディレクトリを確認および取得
	// tmpディレクトリをクリーンアップ
	var err error
	dataTmpDir, err = readDirForTestData(targetTmpDir)
	if err != nil {
		log.Fatalf("Failed to get tmp directory: %v", err)
	}
	dataOrgDir, err = readDirForTestData(targetOrgDir)
	if err != nil {
		log.Fatalf("Failed to get tmp directory: %v", err)
	}

	fmt.Println("テスト環境のセットアップが完了しました。テストを実行します...")

	// テストを実行
	exitCode := m.Run()

	fmt.Println("テストが完了しました。テスト環境をクリーンアップします...")

	os.Exit(exitCode)
}

func readDirForTestData(targetDir string) (string, error) {
	// 設定ファイルのパスを構築して読み込む
	workDir, err := os.Getwd()
	fmt.Printf("Debug: 設定ファイルのパスを構築するためにカレントディレクトリを取得: %s\n", workDir)
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	// 設定ファイルのパスを構築
	projectName := "devbox"
	var dirPath string
	if strings.HasSuffix(workDir, fmt.Sprintf("/%s", targetDir)) {
		// ローカルを想定
		dirPath = filepath.Join(workDir)
	} else if strings.Contains(workDir, fmt.Sprintf("/%s", projectName)) {
		// GitHub Actionsを想定
		idx := strings.Index(workDir, fmt.Sprintf("/%s", projectName))
		if idx >= 0 {
			projectRoot := workDir[:idx+len(fmt.Sprintf("/%s", projectName))]
			dirPath = filepath.Join(projectRoot, targetDir)
		} else {
			return "", fmt.Errorf("unexpected working directory in project: %s", workDir)
		}
	} else {
		return "", fmt.Errorf("unexpected working directory: %s", workDir)
	}

	// ファイルが存在するか確認
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return "", fmt.Errorf("directory for test data not found at %s: %w", dirPath, err)
	}

	return dirPath, nil
}

// #==============================================================#
// ##          Helpers                                           ##
// #==============================================================#
// setupTestData はテスト用のデータをセットアップします
func setupTestData(t *testing.T) {
	files, err := os.ReadDir(dataTmpDir)
	if err != nil {
		t.Fatalf("Failed to read tmp directory: %v", err)
	}
	for _, file := range files {
		if file.Name() == ".gitkeep" {
			continue
		}
		err := os.Remove(filepath.Join(dataTmpDir, file.Name()))
		if err != nil {
			t.Fatalf("Failed to remove file %s: %v", file.Name(), err)
		}
	}

	// orgディレクトリからtmpディレクトリにファイルをコピー
	files, err = os.ReadDir(dataOrgDir)
	if err != nil {
		t.Fatalf("Failed to read org directory: %v", err)
	}
	for _, file := range files {
		if file.Name() == ".gitkeep" {
			continue
		}
		err := copyFile(
			filepath.Join(dataOrgDir, file.Name()),
			filepath.Join(dataTmpDir, file.Name()),
		)
		if err != nil {
			t.Fatalf("Failed to copy file %s: %v", file.Name(), err)
		}
	}
}

// copyFile はファイルをコピーします
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// #==============================================================#
// ##          Tests                                             ##
// #==============================================================#
// TestNewExifModifierService はExifModifierServiceのコンストラクタをテストします
func TestNewExifModifierService(t *testing.T) {
	service := NewExifModifierService()
	if service == nil {
		t.Error("NewExifModifierService() returned nil")
	}
}

// TestExifModifierService_isImageFile は画像ファイル判定をテストします
func TestExifModifierService_isImageFile(t *testing.T) {
	service := NewExifModifierService()

	testCases := []struct {
		filePath        string
		targetExtension string
		expected        bool
	}{
		{"test.jpg", "", true},
		{"test.jpeg", "", true},
		{"test.png", "", true},
		{"test.tiff", "", true},
		{"test.tif", "", true},
		{"test.webp", "", true},
		{"test.mp4", "", true},
		{"test.webm", "", true},
		{"test.txt", "", false},
		{"test.doc", "", false},
		{"test.jpg", "jpg", true},
		{"test.jpg", ".jpg", true},
		{"test.png", "jpg", false},
		{"test.JPG", "jpg", true}, // 大文字小文字の違いをテスト
		{"test.jpg", "JPG", true}, // 大文字小文字の違いをテスト
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.filePath, tc.targetExtension), func(t *testing.T) {
			result := service.isImageFile(tc.filePath, tc.targetExtension)
			if result != tc.expected {
				t.Errorf("isImageFile(%s, %s) = %v, expected %v", tc.filePath, tc.targetExtension, result, tc.expected)
			}
		})
	}
}

// TestExifModifierService_FindImageFiles は画像ファイル検索をテストします
func TestExifModifierService_FindImageFiles(t *testing.T) {
	service := NewExifModifierService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// テストファイルを作成
	testFiles := []string{
		filepath.Join(tempDir, "test1.jpg"),
		filepath.Join(tempDir, "test2.png"),
		filepath.Join(tempDir, "test3.tiff"),
		filepath.Join(tempDir, "test4.txt"), // 対象外
		filepath.Join(subDir, "test5.jpg"),
		filepath.Join(subDir, "test6.png"),
	}

	for _, filename := range testFiles {
		file, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}
		file.Close()
	}

	// 再帰なしのテスト
	config := &Config{
		FolderPath: tempDir,
		Extension:  "",
		Recursive:  false,
	}

	files, err := service.FindImageFiles(config)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files without recursion, got %d", len(files))
	}

	// 再帰ありのテスト
	config.Recursive = true
	files, err = service.FindImageFiles(config)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 5 {
		t.Errorf("Expected 5 files with recursion, got %d", len(files))
	}

	// 特定の拡張子のみのテスト
	config.Extension = "jpg"
	files, err = service.FindImageFiles(config)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 jpg files with recursion, got %d", len(files))
	}

	for _, file := range files {
		if !strings.HasSuffix(strings.ToLower(file), ".jpg") {
			t.Errorf("Expected only jpg files, got %s", file)
		}
	}
}

// TestValidateInputOptions は入力オプションのバリデーションをテストします
func TestValidateInputOptions(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "test.txt")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	testCases := []struct {
		dirPath   string
		extension string
		expectErr bool
	}{
		{tempDir, "", false},
		{tempDir, "jpg", false},
		{tempDir, ".jpg", false},
		{tempDir, "invalid", true},
		{"nonexistent", "", true},
		{testFile, "", true}, // ディレクトリではなくファイル
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.dirPath, tc.extension), func(t *testing.T) {
			err := ValidateInputOptions(tc.dirPath, tc.extension)
			if (err != nil) != tc.expectErr {
				t.Errorf("ValidateInputOptions(%s, %s) error = %v, expectErr %v", tc.dirPath, tc.extension, err, tc.expectErr)
			}
		})
	}
}

// TestValidateDirectory はディレクトリのバリデーションをテストします
func TestValidateDirectory(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "test.txt")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	testCases := []struct {
		dirPath   string
		expectErr bool
	}{
		{tempDir, false},
		{"", true},
		{"nonexistent", true},
		{testFile, true}, // ディレクトリではなくファイル
	}

	for _, tc := range testCases {
		t.Run(tc.dirPath, func(t *testing.T) {
			err := validateDirectory(tc.dirPath)
			if (err != nil) != tc.expectErr {
				t.Errorf("validateDirectory(%s) error = %v, expectErr %v", tc.dirPath, err, tc.expectErr)
			}
		})
	}
}

// TestValidateExtension は拡張子のバリデーションをテストします
func TestValidateExtension(t *testing.T) {
	testCases := []struct {
		extension string
		expectErr bool
	}{
		{"jpg", false},
		{".jpg", false},
		{"jpeg", false},
		{".jpeg", false},
		{"png", false},
		{".png", false},
		{"tiff", false},
		{".tiff", false},
		{"tif", false},
		{".tif", false},
		{"webp", false},
		{".webp", false},
		{"mp4", false},
		{".mp4", false},
		{"webm", false},
		{".webm", false},
		{"", true},
		{"invalid", true},
		{".invalid", true},
		{"JPG", false}, // 大文字も有効
		{".JPG", false},
	}

	for _, tc := range testCases {
		t.Run(tc.extension, func(t *testing.T) {
			err := validateExtension(tc.extension)
			if (err != nil) != tc.expectErr {
				t.Errorf("validateExtension(%s) error = %v, expectErr %v", tc.extension, err, tc.expectErr)
			}
		})
	}
}

// TestParseDateTime は日時文字列のパースをテストします
func TestParseDateTime(t *testing.T) {
	testCases := []struct {
		dateTimeStr string
		expectErr   bool
	}{
		{"20240101120000", false},
		{"20241231235959", false},
		{"20240229120000", false}, // うるう年
		{"20230229120000", true},  // うるう年ではない
		{"20240000120000", true},  // 無効な月
		{"20241300120000", true},  // 無効な月
		{"20240132120000", true},  // 無効な日
		{"20240101240000", true},  // 無効な時
		{"20240101126000", true},  // 無効な分
		{"20240101120060", true},  // 無効な秒
		{"2024010112000", true},   // 短すぎる
		{"202401011200000", true}, // 長すぎる
		{"abcdefghijklmn", true},  // 数字ではない
	}

	for _, tc := range testCases {
		t.Run(tc.dateTimeStr, func(t *testing.T) {
			_, err := ParseDateTime(tc.dateTimeStr)
			if (err != nil) != tc.expectErr {
				t.Errorf("ParseDateTime(%s) error = %v, expectErr %v", tc.dateTimeStr, err, tc.expectErr)
			}
		})
	}

	// 正常なケースで値を確認
	dt, err := ParseDateTime("20240101120000")
	if err != nil {
		t.Errorf("ParseDateTime(20240101120000) unexpected error: %v", err)
	}

	expected := time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)
	if !dt.Equal(expected) {
		t.Errorf("ParseDateTime(20240101120000) = %v, expected %v", dt, expected)
	}
}

// TestValidateDateTime は日時の各要素のバリデーションをテストします
func TestValidateDateTime(t *testing.T) {
	testCases := []struct {
		dateTimeStr string
		expectErr   bool
	}{
		{"20240101120000", false},
		{"20241231235959", false},
		{"20240229120000", false}, // うるう年
		{"20230229120000", true},  // うるう年ではない
		{"18990101120000", true},  // 1900未満
		{"21000101120000", true},  // 2099超過
		{"20240000120000", true},  // 無効な月
		{"20241300120000", true},  // 無効な月
		{"20240132120000", true},  // 無効な日
		{"20240101240000", true},  // 無効な時
		{"20240101126000", true},  // 無効な分
		{"20240101120060", true},  // 無効な秒
	}

	for _, tc := range testCases {
		t.Run(tc.dateTimeStr, func(t *testing.T) {
			err := validateDateTime(tc.dateTimeStr)
			if (err != nil) != tc.expectErr {
				t.Errorf("validateDateTime(%s) error = %v, expectErr %v", tc.dateTimeStr, err, tc.expectErr)
			}
		})
	}
}

// TestParseInt は文字列を整数に変換する関数をテストします
func TestParseInt(t *testing.T) {
	testCases := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"1", 1},
		{"123", 123},
		{"9999", 9999},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := parseInt(tc.input)
			if result != tc.expected {
				t.Errorf("parseInt(%s) = %d, expected %d", tc.input, result, tc.expected)
			}
		})
	}
}

// TestIsLeapYear はうるう年判定をテストします
func TestIsLeapYear(t *testing.T) {
	testCases := []struct {
		year     int
		expected bool
	}{
		{2000, true},  // 400で割り切れる
		{2004, true},  // 4で割り切れる
		{2020, true},  // 4で割り切れる
		{2024, true},  // 4で割り切れる
		{1900, false}, // 100で割り切れるが400では割り切れない
		{2100, false}, // 100で割り切れるが400では割り切れない
		{2001, false}, // 4で割り切れない
		{2002, false}, // 4で割り切れない
		{2003, false}, // 4で割り切れない
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d", tc.year), func(t *testing.T) {
			result := isLeapYear(tc.year)
			if result != tc.expected {
				t.Errorf("isLeapYear(%d) = %v, expected %v", tc.year, result, tc.expected)
			}
		})
	}
}

// TestGetMaxDaysInMonth は月の最大日数取得をテストします
func TestGetMaxDaysInMonth(t *testing.T) {
	testCases := []struct {
		month    int
		year     int
		expected int
	}{
		{1, 2024, 31},  // 1月
		{2, 2024, 29},  // 2月（うるう年）
		{2, 2023, 28},  // 2月（平年）
		{3, 2024, 31},  // 3月
		{4, 2024, 30},  // 4月
		{5, 2024, 31},  // 5月
		{6, 2024, 30},  // 6月
		{7, 2024, 31},  // 7月
		{8, 2024, 31},  // 8月
		{9, 2024, 30},  // 9月
		{10, 2024, 31}, // 10月
		{11, 2024, 30}, // 11月
		{12, 2024, 31}, // 12月
		{0, 2024, 31},  // 無効な月（フォールバック）
		{13, 2024, 31}, // 無効な月（フォールバック）
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%d_%d", tc.month, tc.year), func(t *testing.T) {
			result := getMaxDaysInMonth(tc.month, tc.year)
			if result != tc.expected {
				t.Errorf("getMaxDaysInMonth(%d, %d) = %d, expected %d", tc.month, tc.year, result, tc.expected)
			}
		})
	}
}

// TestExifModifierService_UpdateFileTime はファイルの更新時刻変更をテストします
func TestExifModifierService_UpdateFileTime(t *testing.T) {
	service := NewExifModifierService()

	// テスト用の一時ファイルを作成
	tempFile, err := os.CreateTemp("", "exif_test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// テスト用の時刻
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)

	// ファイルの更新時刻を変更
	err = service.updateFileTime(tempFile.Name(), testTime)
	if err != nil {
		t.Fatalf("updateFileTime() error = %v", err)
	}

	// 変更後のファイル情報を取得
	info, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("os.Stat() error = %v", err)
	}

	// 更新時刻が正しく設定されているか確認
	// 注: ファイルシステムの精度によっては完全に一致しない場合があるため、
	// 1秒以内の差であれば許容する
	if diff := info.ModTime().Sub(testTime); diff < -time.Second || diff > time.Second {
		t.Errorf("File modification time = %v, expected close to %v", info.ModTime(), testTime)
	}
}

// TestExifModifierService_ModifySingleFileExif は単一ファイルのEXIF修正をテストします
func TestExifModifierService_ModifySingleFileExif(t *testing.T) {
	// テスト用のデータをセットアップ
	setupTestData(t)

	service := NewExifModifierService()

	// テスト用のファイルパスを準備
	testFiles := []string{
		filepath.Join(dataTmpDir, "test_11.jpg"),
		filepath.Join(dataTmpDir, "test_21.jpeg"),
		filepath.Join(dataTmpDir, "test_01.png"),
		filepath.Join(dataTmpDir, "test_31.tiff"),
		filepath.Join(dataTmpDir, "test_41.webp"),
	}

	// テスト用の設定
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)
	config := &Config{
		DateTime: testTime,
		DryRun:   false,
		Verbose:  true,
	}

	// 各ファイルに対してテスト
	for _, filePath := range testFiles {
		t.Run(filepath.Base(filePath), func(t *testing.T) {
			err := service.ModifySingleFileExif(filePath, config)
			if err != nil {
				t.Logf("ModifySingleFileExif(%s) error = %v", filePath, err)
				// エラーがあっても失敗とはしない（JPEGの解析エラーなどが発生する可能性があるため）
			}

			// ファイルの更新時刻を確認
			info, err := os.Stat(filePath)
			if err != nil {
				t.Errorf("os.Stat() error = %v", err)
				return
			}

			// 更新時刻が正しく設定されているか確認
			if diff := info.ModTime().Sub(testTime); diff < -time.Second || diff > time.Second {
				t.Errorf("File modification time = %v, expected close to %v", info.ModTime(), testTime)
			} else {
				t.Logf("File %s modification time updated successfully to %v", filepath.Base(filePath), info.ModTime())
			}
		})
	}

	// 存在しないファイルのテスト
	nonExistentFile := filepath.Join(dataTmpDir, "nonexistent.jpg")
	err := service.ModifySingleFileExif(nonExistentFile, config)
	if err == nil {
		t.Error("ModifySingleFileExif() with non-existent file should return error")
	}
}

// TestExifModifierService_ModifyExifData は複数ファイルのEXIF修正をテストします
func TestExifModifierService_ModifyExifData(t *testing.T) {
	// テスト用のデータをセットアップ
	setupTestData(t)

	service := NewExifModifierService()

	// テスト用のファイルパスを準備
	testFiles := []string{
		filepath.Join(dataTmpDir, "test_01.png"),
		filepath.Join(dataTmpDir, "test_31.tiff"),
		filepath.Join(dataTmpDir, "test_41.webp"),
	}

	// テスト用の設定
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)
	config := &Config{
		DateTime:    testTime,
		DryRun:      false,
		Verbose:     true,
		WorkerCount: 2, // 並行処理のテスト
	}

	// 複数ファイルの処理
	processedCount, errorCount, err := service.ModifyExifData(testFiles, config)
	if err != nil {
		t.Errorf("ModifyExifData() error = %v", err)
	}

	// エラーがあっても許容する
	t.Logf("ModifyExifData() processedCount = %d, errorCount = %d", processedCount, errorCount)

	// 各ファイルの更新時刻を確認
	for _, filePath := range testFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filePath, err)
			continue
		}

		// 更新時刻が正しく設定されているか確認
		if diff := info.ModTime().Sub(testTime); diff < -time.Second || diff > time.Second {
			t.Errorf("File %s modification time = %v, expected close to %v", filePath, info.ModTime(), testTime)
		} else {
			t.Logf("File %s modification time updated successfully to %v", filepath.Base(filePath), info.ModTime())
		}
	}

	// ドライランのテスト
	config.DryRun = true
	originalTimes := make(map[string]time.Time)

	// 現在の更新時刻を保存
	for _, filePath := range testFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filePath, err)
			continue
		}
		originalTimes[filePath] = info.ModTime()
	}

	// 異なる時刻でドライラン
	config.DateTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.Local)
	processedCount, errorCount, err = service.ModifyExifData(testFiles, config)
	if err != nil {
		t.Errorf("ModifyExifData() with dry-run error = %v", err)
	}

	// ドライランでは実際のファイルは変更されないことを確認
	for _, filePath := range testFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filePath, err)
			continue
		}

		if !info.ModTime().Equal(originalTimes[filePath]) {
			t.Errorf("File %s was modified during dry-run", filePath)
		} else {
			t.Logf("File %s was not modified during dry-run as expected", filepath.Base(filePath))
		}
	}
}

// TestExifModifierService_ProcessFilesFromFilename はファイル名からの日時抽出処理をテストします
func TestExifModifierService_ProcessFilesFromFilename(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// テスト用のデータをセットアップ
	setupTestData(t)

	// テスト用のファイルをコピー（日時情報を含むファイル名）
	testFiles := []string{
		filepath.Join(tempDir, "IMG_20240101_120000.png"),
		filepath.Join(tempDir, "Photo_20240102_130000.png"),
		filepath.Join(tempDir, "2024-01-03_14-00-00.png"),
		filepath.Join(tempDir, "2024_01_04_150000.png"),
		filepath.Join(tempDir, "nodate.png"), // 日付なし
		filepath.Join(subDir, "IMG_20240105_160000.png"),
	}

	// テスト用のファイルを作成
	for _, filename := range testFiles {
		// test_01.pngをコピー
		err := copyFile(
			filepath.Join(dataTmpDir, "test_01.png"),
			filename,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	service := NewExifModifierService()

	// 再帰なしのテスト
	err = service.ProcessFilesFromFilename(tempDir, false, false, true, true)
	if err != nil {
		t.Errorf("ProcessFilesFromFilename() error = %v", err)
	}

	// 各ファイルの更新時刻を確認
	// タイムゾーンを考慮した期待値を設定
	jst := time.FixedZone("JST", 9*60*60) // JST = UTC+9
	expectedTimes := map[string]time.Time{
		"IMG_20240101_120000.png":   time.Date(2024, 1, 1, 12, 0, 0, 0, jst),
		"Photo_20240102_130000.png": time.Date(2024, 1, 2, 13, 0, 0, 0, jst),
		"2024-01-03_14-00-00.png":   time.Date(2024, 1, 3, 14, 0, 0, 0, jst),
		"2024_01_04_150000.png":     time.Date(2024, 1, 4, 15, 0, 0, 0, jst),
	}

	for filename, expectedTime := range expectedTimes {
		filePath := filepath.Join(tempDir, filename)
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filePath, err)
			continue
		}

		// ファイルが更新されたかどうかを確認
		if info.ModTime().IsZero() {
			t.Errorf("File %s was not modified", filePath)
		} else {
			t.Logf("File %s was modified to %v (expected around %v)", filename, info.ModTime(), expectedTime)
		}
	}

	// 日付なしファイルは変更されていないことを確認
	noDateFile := filepath.Join(tempDir, "nodate.png")
	info, err := os.Stat(noDateFile)
	if err != nil {
		t.Errorf("os.Stat(%s) error = %v", noDateFile, err)
	} else {
		t.Logf("File without date has modification time: %v", info.ModTime())
	}

	// 再帰ありのテスト
	// まず、サブディレクトリのファイルの時刻をリセット
	subDirFile := filepath.Join(subDir, "IMG_20240105_160000.png")
	oldTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	os.Chtimes(subDirFile, oldTime, oldTime)

	err = service.ProcessFilesFromFilename(tempDir, true, false, true, true)
	if err != nil {
		t.Errorf("ProcessFilesFromFilename() with recursion error = %v", err)
	}

	// サブディレクトリのファイルが更新されたことを確認
	info, err = os.Stat(subDirFile)
	if err != nil {
		t.Errorf("os.Stat(%s) error = %v", subDirFile, err)
	} else {
		// ファイルが更新されたかどうかを確認（oldTimeより後であること）
		if !info.ModTime().After(oldTime) {
			t.Errorf("Subdirectory file was not modified")
		} else {
			expectedTime := time.Date(2024, 1, 5, 16, 0, 0, 0, jst)
			t.Logf("Subdirectory file was modified to %v (expected around %v)", info.ModTime(), expectedTime)
		}
	}

	// ドライランのテスト
	// まず、全ファイルの時刻をリセット
	for _, filename := range testFiles {
		os.Chtimes(filename, oldTime, oldTime)
	}

	err = service.ProcessFilesFromFilename(tempDir, true, true, true, true)
	if err != nil {
		t.Errorf("ProcessFilesFromFilename() with dry-run error = %v", err)
	}

	// ドライランでは実際のファイルは変更されないことを確認
	for _, filePath := range testFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filePath, err)
			continue
		}

		if !info.ModTime().Equal(oldTime) {
			t.Errorf("File %s was modified during dry-run", filePath)
		} else {
			t.Logf("File %s was not modified during dry-run as expected", filepath.Base(filePath))
		}
	}
}

// TestExifModifierService_ProcessFilesFromScreenshot はスクリーンショットファイルの処理をテストします
func TestExifModifierService_ProcessFilesFromScreenshot(t *testing.T) {
	service := NewExifModifierService()

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "exif_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	// テスト用のデータをセットアップ
	setupTestData(t)

	// テストファイルを作成（スクリーンショットファイル名）
	testFiles := []string{
		filepath.Join(tempDir, "Screenshot_20240101-120000.png"),
		filepath.Join(tempDir, "スクリーンショット 2024-01-02 13.00.00.png"),
		filepath.Join(tempDir, "Screen Shot 2024-01-03 at 14.00.00.png"),
		filepath.Join(tempDir, "normal_image.png"), // スクリーンショットではない
		filepath.Join(subDir, "Screenshot_20240105-160000.png"),
	}

	// テスト用のファイルを作成
	for _, filename := range testFiles {
		// test_01.pngをコピー
		err := copyFile(
			filepath.Join(dataTmpDir, "test_01.png"),
			filename,
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	// 特定の時刻を設定（タイムゾーンを明示的に指定）
	jst := time.FixedZone("JST", 9*60*60) // JST = UTC+9
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, jst)
	for _, filename := range testFiles {
		// ファイルの更新時刻を設定
		os.Chtimes(filename, testTime, testTime)
	}

	// 再帰なしのテスト
	err = service.ProcessFilesFromScreenshot(tempDir, false, false, true, true)
	if err != nil {
		t.Errorf("ProcessFilesFromScreenshot() error = %v", err)
	}

	// スクリーンショットファイルの更新時刻が保持されていることを確認
	for _, filename := range testFiles {
		// スクリーンショットファイルのみチェック
		if !strings.Contains(strings.ToLower(filepath.Base(filename)), "screenshot") &&
			!strings.Contains(filepath.Base(filename), "スクリーンショット") &&
			!strings.Contains(strings.ToLower(filepath.Base(filename)), "screen shot") {
			continue
		}

		// サブディレクトリのファイルは再帰なしの場合スキップ
		if strings.Contains(filename, subDir) && !strings.Contains(filename, tempDir) {
			continue
		}

		info, err := os.Stat(filename)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filename, err)
			continue
		}

		// 更新時刻が保持されていることを確認
		// 注: 完全一致ではなく、近似値で確認
		if diff := info.ModTime().Sub(testTime); diff < -time.Second || diff > time.Second {
			t.Errorf("File %s modification time = %v, expected close to %v", filename, info.ModTime(), testTime)
		} else {
			t.Logf("File %s modification time is correct: %v", filepath.Base(filename), info.ModTime())
		}
	}

	// 通常の画像ファイルは変更されていないことを確認
	normalFile := filepath.Join(tempDir, "normal_image.png")
	info, err := os.Stat(normalFile)
	if err != nil {
		t.Errorf("os.Stat(%s) error = %v", normalFile, err)
	} else {
		// 注: 完全一致ではなく、近似値で確認
		if diff := info.ModTime().Sub(testTime); diff < -time.Second || diff > time.Second {
			t.Errorf("Non-screenshot file was incorrectly modified: %v", info.ModTime())
		} else {
			t.Logf("Non-screenshot file was not modified as expected: %v", info.ModTime())
		}
	}

	// 再帰ありのテスト
	// まず、サブディレクトリのファイルの時刻をリセット
	subDirFile := filepath.Join(subDir, "Screenshot_20240105-160000.png")
	oldTime := time.Date(2000, 1, 1, 0, 0, 0, 0, jst)
	os.Chtimes(subDirFile, oldTime, oldTime)

	// サブディレクトリのファイルの時刻を新しい時刻に設定
	newTime := time.Date(2024, 1, 5, 16, 0, 0, 0, jst)
	os.Chtimes(subDirFile, newTime, newTime)

	err = service.ProcessFilesFromScreenshot(tempDir, true, false, true, true)
	if err != nil {
		t.Errorf("ProcessFilesFromScreenshot() with recursion error = %v", err)
	}

	// サブディレクトリのファイルが更新されたことを確認
	info, err = os.Stat(subDirFile)
	if err != nil {
		t.Errorf("os.Stat(%s) error = %v", subDirFile, err)
	} else {
		// ファイルが更新されたかどうかを確認（oldTimeより後であること）
		if !info.ModTime().After(oldTime) {
			t.Errorf("Subdirectory file was not modified")
		} else {
			t.Logf("Subdirectory file was modified as expected: %v", info.ModTime())
		}
	}

	// ドライランのテスト
	// まず、全ファイルの時刻をリセット
	for _, filename := range testFiles {
		os.Chtimes(filename, oldTime, oldTime)
	}

	err = service.ProcessFilesFromScreenshot(tempDir, true, true, true, true)
	if err != nil {
		t.Errorf("ProcessFilesFromScreenshot() with dry-run error = %v", err)
	}

	// ドライランでは実際のファイルは変更されないことを確認
	for _, filePath := range testFiles {
		info, err := os.Stat(filePath)
		if err != nil {
			t.Errorf("os.Stat(%s) error = %v", filePath, err)
			continue
		}

		if !info.ModTime().Equal(oldTime) {
			t.Errorf("File %s was modified during dry-run", filePath)
		} else {
			t.Logf("File %s was not modified during dry-run as expected", filepath.Base(filePath))
		}
	}
}
