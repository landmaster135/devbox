package usecases

import (
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
	expectedResults := map[string]string{
		"1-1_image_renamer_with_exif/PXL_20250629_040139164.jpg": "20250629130139.jpg", // 最初のファイルは元の時刻維持
		"1-1_image_renamer_with_exif/PXL_20250629_040139378.jpg": "20250629130140.jpg", // 1秒後
		"1-1_image_renamer_with_exif/PXL_20250629_040139801.jpg": "20250629130141.jpg", // 2秒後
		"1-1_image_renamer_with_exif/PXL_20250629_040138106.jpg": "20250629130138.jpg", // 最初のファイルは元の時刻維持
		"1-1_image_renamer_with_exif/PXL_20250629_040138353.jpg": "20250629130142.jpg", // 4秒後（130139,130140,130141は使用済み）
		"1-1_image_renamer_with_exif/PXL_20250629_040138704.jpg": "20250629130143.jpg", // 5秒後
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
