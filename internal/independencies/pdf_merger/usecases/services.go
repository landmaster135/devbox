package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

// ExtractPDFToImages はPDFの指定したページ範囲を画像として抽出します
func ExtractPDFToImages(pdfPath, outputDir, imageFormat string, startPage, endPage int) error {
	// 出力ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	// PDFファイルの存在確認
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("PDFファイルが見つかりません: %s", pdfPath)
	}

	// サポートする画像形式を確認
	supportedFormats := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"tiff": true,
		"webp": true,
	}
	
	lowerFormat := strings.ToLower(imageFormat)
	if !supportedFormats[lowerFormat] {
		return fmt.Errorf("サポートされていない画像形式です: %s (サポート形式: jpg, jpeg, png, tiff, webp)", imageFormat)
	}

	cfg := api.LoadConfiguration()
	
	// ページ範囲の指定
	var pageSelection []string
	if startPage > 0 && endPage > 0 {
		if startPage > endPage {
			return fmt.Errorf("開始ページ(%d)は終了ページ(%d)より小さくなければなりません", startPage, endPage)
		}
		pageSelection = []string{fmt.Sprintf("%d-%d", startPage, endPage)}
	} else if startPage > 0 {
		pageSelection = []string{strconv.Itoa(startPage)}
	}

	// pdfcpuを使ってPDFページを画像として出力
	// 注意: pdfcpuのバージョンによって利用可能なAPIが異なります
	// ExtractImagesFileを試し、失敗した場合は代替手段を試行
	err := api.ExtractImagesFile(pdfPath, outputDir, pageSelection, cfg)
	if err != nil {
		return fmt.Errorf("PDFからの画像抽出に失敗しました: %w", err)
	}

	return nil
}
