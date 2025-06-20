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


// renameExtractedImagesWithFourDigits は抽出された画像ファイルの名前を4桁連番形式に変更します
func renameExtractedImagesWithFourDigits(outputDir, pdfBaseName string) error {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	// PDFファイル名から拡張子を除去
	pdfName := strings.TrimSuffix(pdfBaseName, filepath.Ext(pdfBaseName))

	// ファイルを分析してページ別に整理
	type fileInfo struct {
		originalPath string
		pageNum      int
		imageNum     int
		ext          string
	}

	var imageFiles []fileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()

		// pdfcpuが生成する可能性のあるファイル名パターンを解析
		// 例: document_001_Im0.jpg, document_page1_Im0.jpg, document_1_Im0.jpg
		if strings.Contains(filename, "_Im") {
			ext := filepath.Ext(filename)

			// _Im部分で分割
			parts := strings.Split(filename, "_Im")
			if len(parts) < 2 {
				continue
			}

			// 左側の部分からページ番号を抽出
			leftPart := parts[0]
			// 右側の部分から画像番号を抽出
			rightPart := parts[1]

			// 画像番号を抽出（拡張子を除去）
			imageNumStr := strings.TrimSuffix(rightPart, ext)
			imageNum, err := strconv.Atoi(imageNumStr)
			if err != nil {
				imageNum = 0 // デフォルト値
			}

			// ページ番号を抽出（様々なパターンに対応）
			var pageNum int

			// パターン1: document_001 のような形式
			leftPart = strings.TrimSuffix(leftPart, "_")

			// 最後のアンダースコア以降を数値として解析
			lastUnderscoreIndex := strings.LastIndex(leftPart, "_")
			if lastUnderscoreIndex >= 0 {
				pageNumStr := leftPart[lastUnderscoreIndex+1:]

				// "page"プレフィックスがある場合は除去
				pageNumStr = strings.TrimPrefix(pageNumStr, "page")

				if num, err := strconv.Atoi(pageNumStr); err == nil {
					pageNum = num
				} else {
					pageNum = 1 // デフォルト値
				}
			} else {
				pageNum = 1 // デフォルト値
			}

			imageFiles = append(imageFiles, fileInfo{
				originalPath: filepath.Join(outputDir, filename),
				pageNum:      pageNum,
				imageNum:     imageNum,
				ext:          ext,
			})
		}
	}

	// ページ番号でソート
	sort.Slice(imageFiles, func(i, j int) bool {
		if imageFiles[i].pageNum != imageFiles[j].pageNum {
			return imageFiles[i].pageNum < imageFiles[j].pageNum
		}
		return imageFiles[i].imageNum < imageFiles[j].imageNum
	})

	// 4桁連番でリネーム
	for i, fileInfo := range imageFiles {
		newName := fmt.Sprintf("%s_%04d%s", pdfName, i+1, fileInfo.ext)
		newPath := filepath.Join(outputDir, newName)

		if err := os.Rename(fileInfo.originalPath, newPath); err != nil {
			return fmt.Errorf("ファイル名の変更に失敗しました: %w", err)
		}
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

	// 抽出された画像ファイルの名前を4桁連番形式に整理
	err = renameExtractedImagesWithFourDigits(outputDir, filepath.Base(pdfPath))
	if err != nil {
		return fmt.Errorf("画像ファイル名の整理に失敗しました: %w", err)
	}

	return nil
}
