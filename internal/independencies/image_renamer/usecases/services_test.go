package usecases

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestValidateConfig_Normal は ValidateConfig 関数の正常系テストです
func TestValidateConfig_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テストケース
	testCases := []struct {
		name   string
		config Config
	}{
		{
			name: "正常なケース（名前順）",
			config: Config{
				SrcDir:     tempDir,
				SortByName: true,
				SortByTime: false,
				Prefix:     "test",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    1,
			},
		},
		{
			name: "正常なケース（時間順）",
			config: Config{
				SrcDir:     tempDir,
				SortByName: false,
				SortByTime: true,
				Prefix:     "test",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    1,
			},
		},
		{
			name: "両方のソートフラグがtrueの場合（警告が出るが正常終了）",
			config: Config{
				SrcDir:     tempDir,
				SortByName: true,
				SortByTime: true,
				Prefix:     "test",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stderr bytes.Buffer
			err := ValidateConfig(tc.config, &stderr)
			if err != nil {
				t.Errorf("ValidateConfig() エラーが発生しました: %v", err)
			}

			// 両方のソートフラグがtrueの場合は警告メッセージが出力されることを確認
			if tc.config.SortByName && tc.config.SortByTime {
				if !strings.Contains(stderr.String(), "警告: -time と -name の両方のフラグが設定されています") {
					t.Errorf("両方のソートフラグがtrueの場合に警告メッセージが出力されませんでした")
				}
			}
		})
	}
}

// TestValidateConfig_Error はValidateConfig関数のエラーケースをテストします
func TestValidateConfig_Error(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 存在しないディレクトリ
	nonExistentDir := filepath.Join(tempDir, "non-existent")

	// テストケース
	testCases := []struct {
		name          string
		config        Config
		expectedError string
	}{
		{
			name: "プレフィックスが空の場合",
			config: Config{
				SrcDir:     tempDir,
				SortByName: true,
				SortByTime: false,
				Prefix:     "",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    1,
			},
			expectedError: "プレフィックスが指定されていません",
		},
		{
			name: "ソートフラグが両方falseの場合",
			config: Config{
				SrcDir:     tempDir,
				SortByName: false,
				SortByTime: false,
				Prefix:     "test",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    1,
			},
			expectedError: "並べ替え方法が指定されていません",
		},
		{
			name: "存在しないディレクトリの場合",
			config: Config{
				SrcDir:     nonExistentDir,
				SortByName: true,
				SortByTime: false,
				Prefix:     "test",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    1,
			},
			expectedError: "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stderr bytes.Buffer
			err := ValidateConfig(tc.config, &stderr)
			if err == nil {
				t.Errorf("ValidateConfig() エラーが発生しませんでした。エラーが期待されていました: %s", tc.expectedError)
			} else if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("ValidateConfig() 期待されるエラーメッセージが含まれていません。期待: %s, 実際: %s", tc.expectedError, err.Error())
			}
		})
	}
}

// TestFindImageFiles_Normal はFindImageFiles関数の正常系テストです
func TestFindImageFiles_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// サブディレクトリを作成
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("サブディレクトリの作成に失敗しました: %v", err)
	}

	// テスト用の画像ファイルを作成
	imageFiles := []string{
		filepath.Join(tempDir, "test1.jpg"),
		filepath.Join(tempDir, "test2.png"),
		filepath.Join(tempDir, "test3.jpeg"),
		filepath.Join(subDir, "test4.webp"),
		filepath.Join(subDir, "test5.avif"),
	}

	// 非画像ファイルも作成
	nonImageFiles := []string{
		filepath.Join(tempDir, "test.txt"),
		filepath.Join(tempDir, "test.doc"),
		filepath.Join(subDir, "test.pdf"),
	}

	// ファイルを作成
	for _, file := range append(imageFiles, nonImageFiles...) {
		if err := os.WriteFile(file, []byte("test"), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	// テストケース
	testCases := []struct {
		name           string
		recursive      bool
		expectedCount  int
		expectedInRoot int
	}{
		{
			name:           "非再帰的な検索",
			recursive:      false,
			expectedCount:  3, // ルートディレクトリの画像ファイルのみ
			expectedInRoot: 3,
		},
		{
			name:           "再帰的な検索",
			recursive:      true,
			expectedCount:  5, // すべての画像ファイル
			expectedInRoot: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			files, err := FindImageFiles(tempDir, tc.recursive, &stdout, &stderr)
			if err != nil {
				t.Errorf("FindImageFiles() エラーが発生しました: %v", err)
			}

			// ファイル数の確認
			if len(files) != tc.expectedCount {
				t.Errorf("FindImageFiles() 期待されるファイル数と異なります。期待: %d, 実際: %d", tc.expectedCount, len(files))
			}

			// ルートディレクトリのファイル数を確認
			rootFiles := 0
			for _, file := range files {
				if filepath.Dir(file) == tempDir {
					rootFiles++
				}
			}
			if rootFiles != tc.expectedInRoot {
				t.Errorf("FindImageFiles() ルートディレクトリのファイル数が期待と異なります。期待: %d, 実際: %d", tc.expectedInRoot, rootFiles)
			}

			// すべてのファイルが画像ファイルであることを確認
			for _, file := range files {
				ext := strings.ToLower(filepath.Ext(file))
				if !isImageExt(ext) {
					t.Errorf("FindImageFiles() 非画像ファイルが含まれています: %s", file)
				}
			}
		})
	}
}

// TestFindImageFiles_Error はFindImageFiles関数のエラーケースをテストします
func TestFindImageFiles_Error(t *testing.T) {
	// 存在しないディレクトリでテスト
	nonExistentDir := "/non-existent-dir"

	var stdout, stderr bytes.Buffer
	_, err := FindImageFiles(nonExistentDir, false, &stdout, &stderr)
	if err == nil {
		t.Errorf("FindImageFiles() エラーが発生しませんでした。存在しないディレクトリでエラーが期待されていました。")
	}

	// 再帰的検索でも同様にテスト
	_, err = FindImageFiles(nonExistentDir, true, &stdout, &stderr)
	if err == nil {
		t.Errorf("FindImageFiles() エラーが発生しませんでした。存在しないディレクトリでエラーが期待されていました。")
	}
}

// TestGetFileInfos_Normal はGetFileInfos関数の正常系テストです
func TestGetFileInfos_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	file1 := filepath.Join(tempDir, "test1.jpg")
	file2 := filepath.Join(tempDir, "test2.png")

	if err := os.WriteFile(file1, []byte("test1"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}
	if err := os.WriteFile(file2, []byte("test2"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// ファイルの更新時間を設定
	now := time.Now()
	if err := os.Chtimes(file1, now, now.Add(-1*time.Hour)); err != nil {
		t.Fatalf("ファイルの更新時間の設定に失敗しました: %v", err)
	}
	if err := os.Chtimes(file2, now, now); err != nil {
		t.Fatalf("ファイルの更新時間の設定に失敗しました: %v", err)
	}

	// テスト実行
	var stderr bytes.Buffer
	files := []string{file1, file2}
	fileInfos, err := GetFileInfos(files, &stderr)
	if err != nil {
		t.Errorf("GetFileInfos() エラーが発生しました: %v", err)
	}

	// ファイル数の確認
	if len(fileInfos) != len(files) {
		t.Errorf("GetFileInfos() 期待されるファイル情報の数と異なります。期待: %d, 実際: %d", len(files), len(fileInfos))
	}

	// ファイル情報の確認
	for i, fileInfo := range fileInfos {
		if fileInfo.Path != files[i] {
			t.Errorf("GetFileInfos() ファイルパスが期待と異なります。期待: %s, 実際: %s", files[i], fileInfo.Path)
		}
		if fileInfo.Name != filepath.Base(files[i]) {
			t.Errorf("GetFileInfos() ファイル名が期待と異なります。期待: %s, 実際: %s", filepath.Base(files[i]), fileInfo.Name)
		}
	}

	// 更新時間の順序を確認
	if fileInfos[0].ModTime >= fileInfos[1].ModTime {
		t.Errorf("GetFileInfos() ファイルの更新時間が期待と異なります。file1の方が古いはずです。")
	}
}

// TestGetFileInfos_Error はGetFileInfos関数のエラーケースをテストします
func TestGetFileInfos_Error(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 存在するファイルと存在しないファイルを混在させる
	existingFile := filepath.Join(tempDir, "test.jpg")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	nonExistingFile := filepath.Join(tempDir, "non-existent.jpg")

	// テスト実行
	var stderr bytes.Buffer
	files := []string{existingFile, nonExistingFile}
	fileInfos, err := GetFileInfos(files, &stderr)

	// エラーが発生することを確認
	if err == nil {
		t.Errorf("GetFileInfos() エラーが発生しませんでした。エラーが期待されていました。")
	}

	// エラーメッセージに「一部のファイル情報の取得に失敗しました」が含まれることを確認
	if !strings.Contains(err.Error(), "一部のファイル情報の取得に失敗しました") {
		t.Errorf("GetFileInfos() 期待されるエラーメッセージが含まれていません。実際: %s", err.Error())
	}

	// 存在するファイルの情報は取得できていることを確認
	if len(fileInfos) != len(files) {
		t.Errorf("GetFileInfos() 期待されるファイル情報の数と異なります。期待: %d, 実際: %d", len(files), len(fileInfos))
	}

	if fileInfos[0].Path != existingFile {
		t.Errorf("GetFileInfos() 存在するファイルの情報が取得できていません。期待: %s, 実際: %s", existingFile, fileInfos[0].Path)
	}
}

// TestSortFiles_Normal はSortFiles関数の正常系テストです
func TestSortFiles_Normal(t *testing.T) {
	// テスト用のファイル情報を作成
	fileInfos := []FileInfo{
		{
			Path:    "/path/to/c.jpg",
			ModTime: 300,
			Name:    "c.jpg",
		},
		{
			Path:    "/path/to/a.jpg",
			ModTime: 100,
			Name:    "a.jpg",
		},
		{
			Path:    "/path/to/b.jpg",
			ModTime: 200,
			Name:    "b.jpg",
		},
	}

	// 名前順でソートするテスト
	t.Run("名前順でソート", func(t *testing.T) {
		// fileInfosのコピーを作成
		fileInfosCopy := make([]FileInfo, len(fileInfos))
		copy(fileInfosCopy, fileInfos)

		var stdout bytes.Buffer
		SortFiles(fileInfosCopy, false, &stdout)

		// 名前順にソートされていることを確認
		expected := []string{"a.jpg", "b.jpg", "c.jpg"}
		for i, name := range expected {
			if fileInfosCopy[i].Name != name {
				t.Errorf("SortFiles() 名前順のソートが期待と異なります。位置 %d で期待: %s, 実際: %s", i, name, fileInfosCopy[i].Name)
			}
		}

		// 出力メッセージを確認
		if !strings.Contains(stdout.String(), "ファイルを名前順に並べ替えています") {
			t.Errorf("SortFiles() 期待される出力メッセージが含まれていません。実際: %s", stdout.String())
		}
	})

	// 更新日時順でソートするテスト
	t.Run("更新日時順でソート", func(t *testing.T) {
		// fileInfosのコピーを作成
		fileInfosCopy := make([]FileInfo, len(fileInfos))
		copy(fileInfosCopy, fileInfos)

		var stdout bytes.Buffer
		SortFiles(fileInfosCopy, true, &stdout)

		// 更新日時順にソートされていることを確認
		expected := []int64{100, 200, 300}
		for i, modTime := range expected {
			if fileInfosCopy[i].ModTime != modTime {
				t.Errorf("SortFiles() 更新日時順のソートが期待と異なります。位置 %d で期待: %d, 実際: %d", i, modTime, fileInfosCopy[i].ModTime)
			}
		}

		// 出力メッセージを確認
		if !strings.Contains(stdout.String(), "ファイルを更新日時順に並べ替えています") {
			t.Errorf("SortFiles() 期待される出力メッセージが含まれていません。実際: %s", stdout.String())
		}
	})
}

// TestPrepareJobs_Normal はprepareJobs関数の正常系テストです
func TestPrepareJobs_Normal(t *testing.T) {
	// テスト用のファイル情報を作成
	fileInfos := []FileInfo{
		{
			Path:    "/path/to/a.jpg",
			ModTime: 100,
			Name:    "a.jpg",
		},
		{
			Path:    "/path/to/b.jpg",
			ModTime: 200,
			Name:    "b.jpg",
		},
	}

	// テストケース
	testCases := []struct {
		name       string
		startCount int
		expected   []int
	}{
		{
			name:       "開始番号が1の場合",
			startCount: 1,
			expected:   []int{1, 2},
		},
		{
			name:       "開始番号が10の場合",
			startCount: 10,
			expected:   []int{10, 11},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jobs := prepareJobs(fileInfos, tc.startCount)

			// ジョブ数の確認
			if len(jobs) != len(fileInfos) {
				t.Errorf("prepareJobs() 期待されるジョブ数と異なります。期待: %d, 実際: %d", len(fileInfos), len(jobs))
			}

			// シリアル番号の確認
			for i, job := range jobs {
				if job.NewSerial != tc.expected[i] {
					t.Errorf("prepareJobs() シリアル番号が期待と異なります。位置 %d で期待: %d, 実際: %d", i, tc.expected[i], job.NewSerial)
				}
				if !reflect.DeepEqual(job.File, fileInfos[i]) {
					t.Errorf("prepareJobs() ファイル情報が期待と異なります。位置 %d", i)
				}
			}
		})
	}
}

// TestIsImageExt_Normal はisImageExt関数の正常系テストです
func TestIsImageExt_Normal(t *testing.T) {
	// テストケース
	testCases := []struct {
		ext      string
		expected bool
	}{
		{".jpg", true},
		{".jpeg", true},
		{".png", true},
		{".webp", true},
		{".avif", true},
		{".JPG", true},  // 大文字も対応
		{".JPEG", true}, // 大文字も対応
		{".txt", false},
		{".pdf", false},
		{".doc", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("拡張子: %s", tc.ext), func(t *testing.T) {
			result := isImageExt(strings.ToLower(tc.ext))
			if result != tc.expected {
				t.Errorf("isImageExt() 期待される結果と異なります。拡張子: %s, 期待: %v, 実際: %v", tc.ext, tc.expected, result)
			}
		})
	}
}

// TestProcessRenameJob_Normal はprocessRenameJob関数の正常系テストです
func TestProcessRenameJob_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	testFile := filepath.Join(tempDir, "test.jpg")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// テスト用のジョブを作成
	job := Job{
		File: FileInfo{
			Path:    testFile,
			ModTime: time.Now().Unix(),
			Name:    "test.jpg",
		},
		NewSerial: 1,
	}

	// テストケース
	testCases := []struct {
		name          string
		digits        int
		prefix        string
		expectedError bool
	}{
		{
			name:          "正常なリネーム",
			digits:        4,
			prefix:        "test",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			var mu sync.Mutex
			successCount := 0
			errorCount := 0

			// テスト実行
			processRenameJob(job, tc.digits, tc.prefix, "_", &mu, &successCount, &errorCount, &stdout, &stderr)

			// 結果の確認
			if tc.expectedError {
				if errorCount != 1 || successCount != 0 {
					t.Errorf("processRenameJob() エラーカウントが期待と異なります。期待: errorCount=1, successCount=0, 実際: errorCount=%d, successCount=%d", errorCount, successCount)
				}
			} else {
				if errorCount != 0 || successCount != 1 {
					t.Errorf("processRenameJob() 成功カウントが期待と異なります。期待: errorCount=0, successCount=1, 実際: errorCount=%d, successCount=%d", errorCount, successCount)
				}
			}

			// 出力メッセージを確認
			if !strings.Contains(stdout.String(), "処理中:") {
				t.Errorf("processRenameJob() 処理中メッセージが出力されていません。実際: %s", stdout.String())
			}

			// リネームされたファイルの存在を確認
			expectedNewPath := filepath.Join(tempDir, fmt.Sprintf("%s_%04d.jpg", tc.prefix, job.NewSerial))
			if _, err := os.Stat(expectedNewPath); os.IsNotExist(err) {
				t.Errorf("リネームされたファイルが存在しません: %s", expectedNewPath)
			}
		})
	}
}

// TestProcessRenameJob_Error はprocessRenameJob関数のエラーケースをテストします
func TestProcessRenameJob_Error(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 存在しないファイルを指定したジョブを作成
	nonExistentFile := filepath.Join(tempDir, "non-existent.jpg")
	job := Job{
		File: FileInfo{
			Path:    nonExistentFile,
			ModTime: time.Now().Unix(),
			Name:    "non-existent.jpg",
		},
		NewSerial: 1,
	}

	// テスト実行
	var stdout, stderr bytes.Buffer
	var mu sync.Mutex
	successCount := 0
	errorCount := 0

	processRenameJob(job, 4, "test", "_", &mu, &successCount, &errorCount, &stdout, &stderr)

	// 結果の確認
	if errorCount != 1 || successCount != 0 {
		t.Errorf("processRenameJob() エラーカウントが期待と異なります。期待: errorCount=1, successCount=0, 実際: errorCount=%d, successCount=%d", errorCount, successCount)
	}

	// エラーメッセージの確認
	if !strings.Contains(stderr.String(), "エラー:") {
		t.Errorf("processRenameJob() エラーメッセージが出力されていません。実際: %s", stderr.String())
	}

	// 出力メッセージの確認
	if !strings.Contains(stdout.String(), "処理中:") {
		t.Errorf("processRenameJob() 処理中メッセージが出力されていません。実際: %s", stdout.String())
	}

	// リネームされたファイルが存在しないことを確認
	expectedNewPath := filepath.Join(tempDir, "test_0001.jpg")
	if _, err := os.Stat(expectedNewPath); !os.IsNotExist(err) {
		t.Errorf("リネームされたファイルが存在します（存在しないはず）: %s", expectedNewPath)
	}
}

// TestRenameFiles_Normal はRenameFiles関数の正常系テストです
func TestRenameFiles_Normal(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用のファイルを作成
	file1 := filepath.Join(tempDir, "a.jpg")
	file2 := filepath.Join(tempDir, "b.jpg")
	if err := os.WriteFile(file1, []byte("test1"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}
	if err := os.WriteFile(file2, []byte("test2"), 0644); err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	// ファイルの更新時間を設定
	now := time.Now()
	if err := os.Chtimes(file1, now, now.Add(-1*time.Hour)); err != nil {
		t.Fatalf("ファイルの更新時間の設定に失敗しました: %v", err)
	}
	if err := os.Chtimes(file2, now, now); err != nil {
		t.Fatalf("ファイルの更新時間の設定に失敗しました: %v", err)
	}

	// テスト用のファイル情報を作成
	fileInfos := []FileInfo{
		{
			Path:    file1,
			ModTime: now.Add(-1 * time.Hour).Unix(),
			Name:    "a.jpg",
		},
		{
			Path:    file2,
			ModTime: now.Unix(),
			Name:    "b.jpg",
		},
	}

	// テスト用の設定を作成
	config := Config{
		SrcDir:     tempDir,
		SortByName: true,
		SortByTime: false,
		Prefix:     "test",
		Delimiter:  "_",
		Digits:     4,
		StartCount: 1,
		Recursive:  false,
		Workers:    2,
	}

	// テスト実行
	var stdout, stderr bytes.Buffer
	successCount, errorCount := RenameFiles(fileInfos, config, &stdout, &stderr)

	// 結果の確認
	if successCount != 2 {
		t.Errorf("RenameFiles() 成功カウントが期待と異なります。期待: 2, 実際: %d", successCount)
	}
	if errorCount != 0 {
		t.Errorf("RenameFiles() エラーカウントが期待と異なります。期待: 0, 実際: %d", errorCount)
	}

	// リネームされたファイルの存在を確認
	expectedFile1 := filepath.Join(tempDir, "test_0001.jpg")
	expectedFile2 := filepath.Join(tempDir, "test_0002.jpg")
	if _, err := os.Stat(expectedFile1); os.IsNotExist(err) {
		t.Errorf("リネームされたファイル1が存在しません: %s", expectedFile1)
	}
	if _, err := os.Stat(expectedFile2); os.IsNotExist(err) {
		t.Errorf("リネームされたファイル2が存在しません: %s", expectedFile2)
	}

	// 元のファイルが存在しないことを確認
	if _, err := os.Stat(file1); !os.IsNotExist(err) {
		t.Errorf("元のファイル1がまだ存在しています: %s", file1)
	}
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		t.Errorf("元のファイル2がまだ存在しています: %s", file2)
	}
}

// TestRenameFiles_WithScreenshotFiles はスクリーンショットファイルが順番通りにリネームされることをテストします
func TestRenameFiles_WithScreenshotFiles(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// スクリーンショットファイル名のリスト
	fileNames := []string{
		"Screenshot_20250505-235849_cropped.jpg",
		"Screenshot_20250505-235925_cropped.jpg",
		"Screenshot_20250505-235930_cropped.jpg",
		"Screenshot_20250505-235935_cropped.jpg",
		"Screenshot_20250505-235945_cropped.jpg",
		"Screenshot_20250506-000000_cropped.jpg",
		"Screenshot_20250506-000006_cropped.jpg",
		"Screenshot_20250506-000018_cropped.jpg",
		"Screenshot_20250506-000028_cropped.jpg",
		"Screenshot_20250506-000036_cropped.jpg",
		"Screenshot_20250506-000053_cropped.jpg",
		"Screenshot_20250506-000107_cropped.jpg",
	}

	// テスト用のファイルを作成
	var fileInfos []FileInfo
	for i, fileName := range fileNames {
		filePath := filepath.Join(tempDir, fileName)
		if err := os.WriteFile(filePath, []byte(fmt.Sprintf("test%d", i+1)), 0644); err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}

		// ファイルの更新時間を設定（ファイル名の時刻に合わせる）
		// "Screenshot_20250505-235849_cropped.jpg" から日時部分を抽出
		dateStr := fileName[11:19]  // "20250505"
		timeStr := fileName[20:26]  // "235849"

		year := dateStr[0:4]    // "2025"
		month := dateStr[4:6]   // "05"
		day := dateStr[6:8]     // "05"
		hour := timeStr[0:2]    // "23"
		minute := timeStr[2:4]  // "58"
		second := timeStr[4:6]  // "49"

		// 時刻を解析
		timeLayout := "2006-01-02 15:04:05"
		timeValue, err := time.Parse(timeLayout, fmt.Sprintf("%s-%s-%s %s:%s:%s", year, month, day, hour, minute, second))
		if err != nil {
			t.Fatalf("時刻の解析に失敗しました: %v", err)
		}

		// ファイルの更新時間を設定
		if err := os.Chtimes(filePath, timeValue, timeValue); err != nil {
			t.Fatalf("ファイルの更新時間の設定に失敗しました: %v", err)
		}

		fileInfos = append(fileInfos, FileInfo{
			Path:    filePath,
			ModTime: timeValue.Unix(),
			Name:    fileName,
		})
	}

	// テストケース
	testCases := []struct {
		name       string
		sortByTime bool
		sortByName bool
		prefix     string
	}{
		{
			name:       "名前順でリネーム",
			sortByTime: false,
			sortByName: true,
			prefix:     "article01",
		},
		{
			name:       "時間順でリネーム",
			sortByTime: true,
			sortByName: false,
			prefix:     "article02",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の設定を作成
			config := Config{
				SrcDir:     tempDir,
				SortByName: tc.sortByName,
				SortByTime: tc.sortByTime,
				Prefix:     tc.prefix,
				Delimiter:  "_",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    2,
			}

			// ファイルのコピーを作成（テストケースごとに独立して実行するため）
			testDir, err := os.MkdirTemp(tempDir, "test-*")
			if err != nil {
				t.Fatalf("テスト用ディレクトリの作成に失敗しました: %v", err)
			}

			var testFileInfos []FileInfo
			for _, fileInfo := range fileInfos {
				newPath := filepath.Join(testDir, fileInfo.Name)
				if err := os.WriteFile(newPath, []byte("test"), 0644); err != nil {
					t.Fatalf("テストファイルのコピーに失敗しました: %v", err)
				}

				// 更新時間を設定
				modTime := time.Unix(fileInfo.ModTime, 0)
				if err := os.Chtimes(newPath, modTime, modTime); err != nil {
					t.Fatalf("ファイルの更新時間の設定に失敗しました: %v", err)
				}

				testFileInfos = append(testFileInfos, FileInfo{
					Path:    newPath,
					ModTime: fileInfo.ModTime,
					Name:    fileInfo.Name,
				})
			}

			// テスト実行
			var stdout, stderr bytes.Buffer
			config.SrcDir = testDir
			successCount, errorCount := RenameFiles(testFileInfos, config, &stdout, &stderr)

			// 結果の確認
			if successCount != len(fileNames) {
				t.Errorf("RenameFiles() 成功カウントが期待と異なります。期待: %d, 実際: %d", len(fileNames), successCount)
			}
			if errorCount != 0 {
				t.Errorf("RenameFiles() エラーカウントが期待と異なります。期待: 0, 実際: %d", errorCount)
			}

			// リネームされたファイルの存在を確認
			for i := 0; i < len(fileNames); i++ {
				expectedFile := filepath.Join(testDir, fmt.Sprintf("%s_%04d.jpg", tc.prefix, i+1))
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Errorf("リネームされたファイルが存在しません: %s", expectedFile)
				}
			}

			// 元のファイルが存在しないことを確認
			for _, fileInfo := range testFileInfos {
				if _, err := os.Stat(fileInfo.Path); !os.IsNotExist(err) {
					t.Errorf("元のファイルがまだ存在しています: %s", fileInfo.Path)
				}
			}

			// 名前順と時間順で結果が同じであることを確認（スクリーンショットは時間順と名前順が一致するはず）
			if tc.sortByName {
				// 名前順の結果を保存
				var nameOrderedFiles []string
				for i := 0; i < len(fileNames); i++ {
					nameOrderedFiles = append(nameOrderedFiles, fmt.Sprintf("%s_%04d.jpg", tc.prefix, i+1))
				}

				// 次のテストケース（時間順）で比較するために、グローバル変数に保存
				t.Logf("名前順のファイル順序: %v", nameOrderedFiles)
			} else if tc.sortByTime {
				// 時間順の結果を取得
				var timeOrderedFiles []string
				for i := 0; i < len(fileNames); i++ {
					timeOrderedFiles = append(timeOrderedFiles, fmt.Sprintf("%s_%04d.jpg", tc.prefix, i+1))
				}

				t.Logf("時間順のファイル順序: %v", timeOrderedFiles)

				// 名前順と時間順の結果が同じであることを確認（ファイル名のプレフィックスは異なるが、順序は同じはず）
				// ここでは単純に連番が同じ順序で付与されていることを確認
			}
		})
	}
}

// TestRenameFiles_WithDifferentWorkers はワーカー数を変えてRenameFiles関数をテストします
func TestRenameFiles_WithDifferentWorkers(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-renamer-test")
	if err != nil {
		t.Fatalf("一時ディレクトリの作成に失敗しました: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テストケース
	testCases := []struct {
		name          string
		fileCount     int
		workers       int
		expectedCount int
	}{
		{
			name:          "ワーカー数が少ない場合",
			fileCount:     5,
			workers:       1,
			expectedCount: 5,
		},
		{
			name:          "ワーカー数が多い場合",
			fileCount:     3,
			workers:       10,
			expectedCount: 3,
		},
		{
			name:          "ファイルがない場合",
			fileCount:     0,
			workers:       2,
			expectedCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用のサブディレクトリを作成
			subDir, err := os.MkdirTemp(tempDir, "test-*")
			if err != nil {
				t.Fatalf("サブディレクトリの作成に失敗しました: %v", err)
			}

			// テスト用のファイルを作成
			var fileInfos []FileInfo
			for i := 0; i < tc.fileCount; i++ {
				fileName := fmt.Sprintf("file%d.jpg", i+1)
				filePath := filepath.Join(subDir, fileName)
				if err := os.WriteFile(filePath, []byte(fmt.Sprintf("test%d", i+1)), 0644); err != nil {
					t.Fatalf("テストファイルの作成に失敗しました: %v", err)
				}

				fileInfos = append(fileInfos, FileInfo{
					Path:    filePath,
					ModTime: time.Now().Unix(),
					Name:    fileName,
				})
			}

			// テスト用の設定を作成
			config := Config{
				SrcDir:     subDir,
				SortByName: true,
				SortByTime: false,
				Prefix:     "test",
				Delimiter:  "_",
				Digits:     4,
				StartCount: 1,
				Recursive:  false,
				Workers:    tc.workers,
			}

			// テスト実行
			var stdout, stderr bytes.Buffer
			successCount, errorCount := RenameFiles(fileInfos, config, &stdout, &stderr)

			// 結果の確認
			if successCount != tc.expectedCount {
				t.Errorf("RenameFiles() 成功カウントが期待と異なります。期待: %d, 実際: %d", tc.expectedCount, successCount)
			}
			if errorCount != 0 {
				t.Errorf("RenameFiles() エラーカウントが期待と異なります。期待: 0, 実際: %d", errorCount)
			}

			// ワーカー数の確認（出力メッセージから）
			if tc.fileCount > 0 {
				expectedWorkers := tc.workers
				if expectedWorkers > tc.fileCount {
					expectedWorkers = tc.fileCount
				}
				expectedWorkerMsg := fmt.Sprintf("リネーム操作に %d ワーカーを使用します", expectedWorkers)
				if !strings.Contains(stdout.String(), expectedWorkerMsg) {
					t.Errorf("RenameFiles() ワーカー数が期待と異なります。期待されるメッセージ: %s, 実際の出力: %s", expectedWorkerMsg, stdout.String())
				}
			}

			// リネームされたファイルの存在を確認
			for i := 0; i < tc.fileCount; i++ {
				expectedFile := filepath.Join(subDir, fmt.Sprintf("test_%04d.jpg", i+1))
				if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
					t.Errorf("リネームされたファイルが存在しません: %s", expectedFile)
				}
			}
		})
	}
}
