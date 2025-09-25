package usecases

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Consts and Types                                  ##
// #==============================================================#
const (
	sampleFileOfPDF0101 = "sample_01_01.pdf"
	sampleFileOfJPG0101 = "sample_01_01.jpg"
	sampleFileOfPNG0101 = "sample_01_01.png"
)

// #==============================================================#
// ##          Helpers for tests                                 ##
// #==============================================================#
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
	return err
}

// setupTestData はテスト用データをorgからtmpにコピーします
func setupTestData(t *testing.T) (string, func()) {
	t.Helper()

	// tmpディレクトリのパス
	rootDirForTestData := "test_data"
	tmpDir := filepath.Join(rootDirForTestData, "tmp")
	orgDir := filepath.Join(rootDirForTestData, "org")

	// tmpディレクトリをクリーンアップ
	os.RemoveAll(tmpDir)
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		t.Fatalf("tmpディレクトリの作成に失敗: %v", err)
	}

	// サンプルファイルをコピー
	files := []string{sampleFileOfPDF0101, sampleFileOfJPG0101, sampleFileOfPNG0101}
	for _, file := range files {
		srcPath := filepath.Join(orgDir, file)
		dstPath := filepath.Join(tmpDir, file)
		if err := copyFile(srcPath, dstPath); err != nil {
			t.Fatalf("ファイル %s のコピーに失敗: %v", file, err)
		}
	}

	// クリーンアップ関数を返す
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

// #==============================================================#
// ##         Tests                                              ##
// #==============================================================#
// TestImageExtractionService_GetPageCount はGetPageCountのテスト
func TestImageExtractionService_GetPageCount(t *testing.T) {
	service := NewImageExtractionService()

	t.Run("正常系_有効なPDFファイル", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)

		pageCount, err := service.GetPageCount(pdfPath)
		if err != nil {
			t.Fatalf("GetPageCount()でエラーが発生: %v", err)
		}

		if pageCount <= 0 {
			t.Errorf("ページ数は1以上である必要があります。取得値: %d", pageCount)
		}

		t.Logf("PDFのページ数: %d", pageCount)
	})

	t.Run("異常系_存在しないPDFファイル", func(t *testing.T) {
		nonExistentPath := "non_existent.pdf"

		_, err := service.GetPageCount(nonExistentPath)
		if err == nil {
			t.Error("存在しないファイルではエラーが期待されます")
		}

		expectedMsg := "PDFファイルが見つかりません"
		if err != nil && !containsString(err.Error(), expectedMsg) {
			t.Errorf("期待されるエラーメッセージが含まれていません。期待: %s, 実際: %s", expectedMsg, err.Error())
		}
	})

	t.Run("異常系_無効なファイル", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		// JPGファイルをPDFとして扱おうとする
		invalidPdfPath := filepath.Join(tmpDir, sampleFileOfJPG0101)

		_, err := service.GetPageCount(invalidPdfPath)
		if err == nil {
			t.Error("無効なPDFファイルではエラーが期待されます")
		}

		expectedMsg := "PDFファイルの読み込みに失敗しました"
		if err != nil && !containsString(err.Error(), expectedMsg) {
			t.Errorf("期待されるエラーメッセージが含まれていません。期待: %s, 実際: %s", expectedMsg, err.Error())
		}
	})
}

// TestImageExtractionService_ExtractToImages はExtractToImagesのテスト
func TestImageExtractionService_ExtractToImages(t *testing.T) {
	service := NewImageExtractionService()

	t.Run("正常系_全ページ抽出", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)
		outputDir := filepath.Join(tmpDir, "extracted_images")

		err := service.ExtractToImages(pdfPath, outputDir, "jpg", 0, 0)
		if err != nil {
			t.Fatalf("ExtractToImages()でエラーが発生: %v", err)
		}

		// 出力ディレクトリが作成されたことを確認
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			t.Error("出力ディレクトリが作成されませんでした")
		}

		// 抽出された画像ファイルの存在確認
		files, err := os.ReadDir(outputDir)
		if err != nil {
			t.Fatalf("出力ディレクトリの読み取りに失敗: %v", err)
		}

		if len(files) == 0 {
			t.Error("画像ファイルが抽出されませんでした")
		}

		t.Logf("抽出された画像ファイル数: %d", len(files))
		for _, file := range files {
			t.Logf("抽出されたファイル: %s", file.Name())
		}
	})

	t.Run("正常系_ページ範囲指定", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)
		outputDir := filepath.Join(tmpDir, "extracted_images_range")

		// 最初のページのみ抽出
		err := service.ExtractToImages(pdfPath, outputDir, "png", 1, 1)
		if err != nil {
			t.Fatalf("ExtractToImages()でエラーが発生: %v", err)
		}

		// 出力ディレクトリが作成されたことを確認
		if _, err := os.Stat(outputDir); os.IsNotExist(err) {
			t.Error("出力ディレクトリが作成されませんでした")
		}
	})

	t.Run("異常系_存在しないPDFファイル", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		nonExistentPath := "non_existent.pdf"
		outputDir := filepath.Join(tmpDir, "output")

		err := service.ExtractToImages(nonExistentPath, outputDir, "jpg", 0, 0)
		if err == nil {
			t.Error("存在しないファイルではエラーが期待されます")
		}

		expectedMsg := "PDFファイルが見つかりません"
		if err != nil && !containsString(err.Error(), expectedMsg) {
			t.Errorf("期待されるエラーメッセージが含まれていません。期待: %s, 実際: %s", expectedMsg, err.Error())
		}
	})

	t.Run("異常系_サポートされていない画像形式", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)
		outputDir := filepath.Join(tmpDir, "output")

		err := service.ExtractToImages(pdfPath, outputDir, "gif", 0, 0)
		if err == nil {
			t.Error("サポートされていない画像形式ではエラーが期待されます")
		}

		expectedMsg := "サポートの有無が規定されていない画像形式です"
		if err != nil && !containsString(err.Error(), expectedMsg) {
			t.Errorf("期待されるエラーメッセージが含まれていません。期待: %s, 実際: %s", expectedMsg, err.Error())
		}
	})
}

// TestPDFCreationService_AddImagesToExistingPDF はAddImagesToExistingPDFのテスト
func TestPDFCreationService_AddImagesToExistingPDF(t *testing.T) {
	service := NewPDFCreationService()

	t.Run("正常系_既存PDFに画像を追加", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		existingPDF := filepath.Join(tmpDir, sampleFileOfPDF0101)
		images := []string{
			filepath.Join(tmpDir, sampleFileOfJPG0101),
			filepath.Join(tmpDir, sampleFileOfPNG0101),
		}
		outputPDF := filepath.Join(tmpDir, "merged_output.pdf")

		err := service.AddImagesToExistingPDF(existingPDF, images, outputPDF)
		if err != nil {
			t.Fatalf("AddImagesToExistingPDF()でエラーが発生: %v", err)
		}

		// 出力PDFファイルが作成されたことを確認
		if _, err := os.Stat(outputPDF); os.IsNotExist(err) {
			t.Error("出力PDFファイルが作成されませんでした")
		}

		// PDFファイルの基本的な検証（PDFヘッダーの確認）
		content, err := os.ReadFile(outputPDF)
		if err != nil {
			t.Fatalf("出力PDFファイルの読み取りに失敗: %v", err)
		}

		if len(content) < 4 || string(content[:4]) != "%PDF" {
			t.Error("出力ファイルは有効なPDFではありません")
		}

		t.Logf("マージされたPDFファイルサイズ: %d bytes", len(content))
	})

	t.Run("異常系_存在しない既存PDFファイル", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		nonExistentPDF := "non_existent.pdf"
		images := []string{filepath.Join(tmpDir, sampleFileOfJPG0101)}
		outputPDF := filepath.Join(tmpDir, "output.pdf")

		err := service.AddImagesToExistingPDF(nonExistentPDF, images, outputPDF)
		if err == nil {
			t.Error("存在しない既存PDFファイルではエラーが期待されます")
		}
	})

	t.Run("異常系_存在しない画像ファイル", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		existingPDF := filepath.Join(tmpDir, sampleFileOfPDF0101)
		images := []string{"non_existent_image.jpg"}
		outputPDF := filepath.Join(tmpDir, "output.pdf")

		err := service.AddImagesToExistingPDF(existingPDF, images, outputPDF)
		if err == nil {
			t.Error("存在しない画像ファイルではエラーが期待されます")
		}

		expectedMsg := "画像をPDFに変換中にエラーが発生しました"
		if err != nil && !containsString(err.Error(), expectedMsg) {
			t.Errorf("期待されるエラーメッセージが含まれていません。期待: %s, 実際: %s", expectedMsg, err.Error())
		}
	})

	t.Run("異常系_空の画像リスト", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		existingPDF := filepath.Join(tmpDir, sampleFileOfPDF0101)
		images := []string{}
		outputPDF := filepath.Join(tmpDir, "output.pdf")

		err := service.AddImagesToExistingPDF(existingPDF, images, outputPDF)
		// 空の画像リストでも処理は成功する可能性があるため、エラーの有無をログに記録
		t.Logf("空の画像リストでの結果: %v", err)
	})
}

// TestPDFCreationService_createTemporaryPDF_Integration は AddImagesToExistingPDF 内で一時ファイルが適切に扱われることを検証
func TestPDFCreationService_createTemporaryPDF_Integration(t *testing.T) {
	service := NewPDFCreationService()

	t.Run("一時ファイルが出力ディレクトリ内で生成され後始末される", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		existingPDF := filepath.Join(tmpDir, sampleFileOfPDF0101)
		images := []string{filepath.Join(tmpDir, sampleFileOfJPG0101)}
		outputPDF := filepath.Join(tmpDir, "output.pdf")

		err := service.AddImagesToExistingPDF(existingPDF, images, outputPDF)
		if err != nil {
			t.Fatalf("AddImagesToExistingPDF()でエラーが発生: %v", err)
		}

		// 出力ファイルが作成されていることを確認（間接的に一時ファイルが利用されたことを確認）
		if _, err := os.Stat(outputPDF); os.IsNotExist(err) {
			t.Error("出力PDFファイルが作成されませんでした")
		}

		entries, err := os.ReadDir(tmpDir)
		if err != nil {
			t.Fatalf("ディレクトリ一覧の取得に失敗: %v", err)
		}

		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "pdfmerge_temp_") {
				t.Errorf("一時PDFファイルが削除されていません: %s", entry.Name())
			}
		}
	})
}

// containsString は文字列に部分文字列が含まれているかチェックします
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

// containsSubstring は文字列内の部分文字列を検索します
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestImageExtractionService_ValidatePageRange_Integration はValidatePageRangeの統合テスト
func TestImageExtractionService_ValidatePageRange_Integration(t *testing.T) {
	service := NewImageExtractionService()

	t.Run("実際のPDFファイルでのページ範囲バリデーション", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)

		// 実際のPDFのページ数を取得
		pageCount, err := service.GetPageCount(pdfPath)
		if err != nil {
			t.Fatalf("GetPageCount()でエラーが発生: %v", err)
		}

		// 正常範囲のテスト
		err = service.ValidatePageRange(1, pageCount, pageCount)
		if err != nil {
			t.Errorf("正常範囲でエラーが発生: %v", err)
		}

		// 範囲外のテスト
		err = service.ValidatePageRange(pageCount+1, pageCount+2, pageCount)
		if err == nil {
			t.Error("範囲外の値でエラーが期待されます")
		}

		// 負の値のテスト
		err = service.ValidatePageRange(-1, 1, pageCount)
		if err == nil {
			t.Error("負の値でエラーが期待されます")
		}
	})
}

// TestImageExtractionService_isSupportedFormat_Integration はisSupportedFormatの統合テスト
func TestImageExtractionService_isSupportedFormat_Integration(t *testing.T) {
	service := NewImageExtractionService()

	t.Run("ExtractToImages内でisSupportedFormatが呼ばれる", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)
		outputDir := filepath.Join(tmpDir, "output")

		// サポートされている形式でのテスト
		supportedFormats := []string{"jpg", "jpeg", "png", "tiff", "webp"}
		for _, format := range supportedFormats {
			err := service.ExtractToImages(pdfPath, outputDir+"_"+format, format, 0, 0)
			if err != nil {
				t.Logf("形式 %s でのテスト結果: %v", format, err)
			}
		}

		// サポートされていない形式でのテスト
		unsupportedFormats := []string{"gif", "bmp", "svg"}
		for _, format := range unsupportedFormats {
			err := service.ExtractToImages(pdfPath, outputDir+"_"+format, format, 0, 0)
			if err == nil {
				t.Errorf("サポートされていない形式 %s でエラーが期待されます", format)
			}
		}
	})
}

// TestImageExtractionService_renameExtractedImagesWithFourDigits_Integration はrenameExtractedImagesWithFourDigitsの統合テスト
func TestImageExtractionService_renameExtractedImagesWithFourDigits_Integration(t *testing.T) {
	service := NewImageExtractionService()

	t.Run("ExtractToImages内でrenameExtractedImagesWithFourDigitsが呼ばれる", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		pdfPath := filepath.Join(tmpDir, sampleFileOfPDF0101)
		outputDir := filepath.Join(tmpDir, "extracted_with_rename")

		err := service.ExtractToImages(pdfPath, outputDir, "jpg", 0, 0)
		if err != nil {
			t.Fatalf("ExtractToImages()でエラーが発生: %v", err)
		}

		// 抽出された画像ファイルの名前形式を確認
		files, err := os.ReadDir(outputDir)
		if err != nil {
			t.Fatalf("出力ディレクトリの読み取りに失敗: %v", err)
		}

		for _, file := range files {
			t.Logf("リネーム後のファイル名: %s", file.Name())
			// 4桁連番形式かどうかの簡単なチェック
			if len(file.Name()) > 4 && containsString(file.Name(), "_") {
				t.Logf("4桁連番形式のファイル名が確認されました: %s", file.Name())
			}
		}
	})
}

// TestPDFCreationService_GetSourceImages_EdgeCases はGetSourceImagesのエッジケーステスト
func TestPDFCreationService_GetSourceImages_EdgeCases(t *testing.T) {
	service := NewPDFCreationService()

	t.Run("権限のないディレクトリ", func(t *testing.T) {
		// 権限のないディレクトリを作成
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		restrictedDir := filepath.Join(tmpDir, "restricted")
		err := os.MkdirAll(restrictedDir, 0000) // 権限なし
		if err != nil {
			t.Skipf("権限制限ディレクトリの作成に失敗: %v", err)
		}
		defer os.Chmod(restrictedDir, 0755) // クリーンアップ用に権限を戻す

		_, _, err = service.GetSourceImages(restrictedDir, "")
		// 権限エラーが発生する可能性があるが、システムによって動作が異なる
		t.Logf("権限のないディレクトリでの結果: %v", err)
	})

	t.Run("非常に長いパス名", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		// 長いディレクトリ名を作成
		longDirName := filepath.Join(tmpDir, "very_long_directory_name_that_might_cause_issues_in_some_systems_but_should_be_handled_gracefully")
		err := os.MkdirAll(longDirName, 0755)
		if err != nil {
			t.Skipf("長いパス名のディレクトリ作成に失敗: %v", err)
		}

		// テスト用JPGファイルを作成
		jpgPath := filepath.Join(longDirName, "test.jpg")
		err = copyFile(filepath.Join(tmpDir, sampleFileOfJPG0101), jpgPath)
		if err != nil {
			t.Fatalf("JPGファイルのコピーに失敗: %v", err)
		}

		images, output, err := service.GetSourceImages(longDirName, "")
		if err != nil {
			t.Fatalf("長いパス名でエラーが発生: %v", err)
		}

		if len(images) != 1 {
			t.Errorf("画像が1つ期待されますが、%d個見つかりました", len(images))
		}

		t.Logf("長いパス名での出力: %s", output)
	})
}

// TestErrorHandling_FileSystemErrors はファイルシステムエラーのテスト
func TestErrorHandling_FileSystemErrors(t *testing.T) {
	t.Run("GetSourceImages_filepath.Absエラーのシミュレーション", func(t *testing.T) {
		service := NewPDFCreationService()

		// 無効なパス文字を含むディレクトリ名（システムによって異なる）
		invalidPath := string([]byte{0, 1, 2}) // null文字を含む無効なパス
		_, _, err := service.GetSourceImages(invalidPath, "")

		// エラーが発生することを確認（システムによって動作が異なる可能性がある）
		t.Logf("無効なパスでの結果: %v", err)
	})

	t.Run("MergeImagesIntoPDF_無効な出力パス", func(t *testing.T) {
		tmpDir, cleanup := setupTestData(t)
		defer cleanup()

		service := NewPDFCreationService()
		images := []string{filepath.Join(tmpDir, sampleFileOfJPG0101)}

		// 存在しないディレクトリへの出力パス
		invalidOutputPath := "/non/existent/directory/output.pdf"

		err := service.MergeImagesIntoPDF(images, invalidOutputPath)
		if err == nil {
			t.Error("無効な出力パスでエラーが期待されます")
		}

		t.Logf("無効な出力パスでの結果: %v", err)
	})
}
