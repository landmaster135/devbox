package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	api "github.com/pdfcpu/pdfcpu/pkg/api"
	types "github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

func GetSourceImages(dir string, out string) ([]string, string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, "", err
	}

	// 出力名のデフォルト: <フォルダー名>.pdf
	if out == "" {
		base := filepath.Base(absDir)
		out = filepath.Join(absDir, base+".pdf")
	}

	// ---- 画像ファイルを収集 ----
	var images []string
	err = filepath.WalkDir(absDir, func(p string, d os.DirEntry, _ error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(p))
		if ext == ".jpg" {
			images = append(images, p)
		}
		return nil
	})
	if err != nil {
		return nil, "", err
	}

	if len(images) == 0 {
		fmt.Println("画像が見つかりませんでした。終了します。")
		return nil, "", nil
	}
	sort.Strings(images) // PowerShell の Sort-Object 相当

	return images, out, nil
}

func MergeImagesIntoPDF(images []string, output string) error {
	cfg := api.LoadConfiguration()
	// Unit is used in commands for layout
	cfg.Unit = types.POINTS
	// Compress non-stream object to stream object
	cfg.WriteObjectStream = true
	// Remove unused fonts and images from resource dictionary
	cfg.OptimizeResourceDicts = true
	// Share duplicated streams in all pages
	cfg.OptimizeDuplicateContentStreams = true
	if err := api.ImportImagesFile(images, output, nil, cfg); err != nil {
		return err
	}
	return nil
}

func getNameOfTemporaryPDF() string {
	timestamp := time.Now().Format("20060102150405")
	tempPDF := fmt.Sprintf("added_%s.pdf", timestamp)
	return tempPDF
}

// AddImagesToExistingPDF は既存のPDFファイルに画像ページを追加します
func AddImagesToExistingPDF(existingPDF string, images []string, output string) error {
	cfg := api.LoadConfiguration()
	// Unit is used in commands for layout
	cfg.Unit = types.POINTS
	// Compress non-stream object to stream object
	cfg.WriteObjectStream = true
	// Remove unused fonts and images from resource dictionary
	cfg.OptimizeResourceDicts = true
	// Share duplicated streams in all pages
	cfg.OptimizeDuplicateContentStreams = true

	// 一時的に画像からPDFを作成
	tempPDF := getNameOfTemporaryPDF()
	defer os.Remove(tempPDF) // 関数終了時に一時ファイルを削除

	// 画像をPDFに変換
	if err := api.ImportImagesFile(images, tempPDF, nil, cfg); err != nil {
		return fmt.Errorf("画像をPDFに変換中にエラーが発生しました: %w", err)
	}

	// 既存PDFと新規PDFをマージ
	inFiles := []string{existingPDF, tempPDF}
	if err := api.MergeCreateFile(inFiles, output, false, cfg); err != nil {
		return fmt.Errorf("PDFのマージ中にエラーが発生しました: %w", err)
	}

	return nil
}
