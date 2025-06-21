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

// ImageExtractionService は PDF から画像を抽出するためのサービス
type ImageExtractionService struct{}

// NewImageExtractionService は新しい ImageExtractionService のインスタンスを作成します
func NewImageExtractionService() *ImageExtractionService {
	return &ImageExtractionService{}
}

// PDFCreationService は PDF を作成・結合するためのサービス
type PDFCreationService struct{}

// NewPDFCreationService は新しい PDFCreationService のインスタンスを作成します
func NewPDFCreationService() *PDFCreationService {
	return &PDFCreationService{}
}

// GetSourceImages は指定されたディレクトリから画像ファイルを取得します
func (s *PDFCreationService) GetSourceImages(dir string, out string) ([]string, string, error) {
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

// MergeImagesIntoPDF は画像ファイルを1つのPDFに結合します
func (s *PDFCreationService) MergeImagesIntoPDF(images []string, output string) error {
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

// getNameOfTemporaryPDF は一時的なPDFファイル名を生成します
func (s *PDFCreationService) getNameOfTemporaryPDF() string {
	timestamp := time.Now().Format("20060102150405")
	tempPDF := fmt.Sprintf("added_%s.pdf", timestamp)
	return tempPDF
}

// AddImagesToExistingPDF は既存のPDFファイルに画像ページを追加します
func (s *PDFCreationService) AddImagesToExistingPDF(existingPDF string, images []string, output string) error {
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
	tempPDF := s.getNameOfTemporaryPDF()
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

// fileInfo は抽出された画像ファイルの情報を格納する構造体
type fileInfo struct {
	originalPath string
	pageNum      int
	imageNum     int
	ext          string
}

// extractPageNumber は左側部分からページ番号を抽出します
func (s *ImageExtractionService) extractPageNumber(leftPart string) int {
	// パターン1: document_001 のような形式
	leftPart = strings.TrimSuffix(leftPart, "_")

	// 最後のアンダースコア以降を数値として解析
	lastUnderscoreIndex := strings.LastIndex(leftPart, "_")
	if lastUnderscoreIndex >= 0 {
		pageNumStr := leftPart[lastUnderscoreIndex+1:]

		// "page"プレフィックスがある場合は除去
		pageNumStr = strings.TrimPrefix(pageNumStr, "page")

		if num, err := strconv.Atoi(pageNumStr); err == nil {
			return num
		}
	}
	return 1 // デフォルト値
}

// parseImageFileName はファイル名を解析してファイル情報を抽出します
func (s *ImageExtractionService) parseImageFileName(filename, outputDir string) (fileInfo, error) {
	// pdfcpuが生成する可能性のあるファイル名パターンを解析
	// 例: document_001_Im0.jpg, document_page1_Im0.jpg, document_1_Im0.jpg
	specifiedChr := "_Im"
	if !strings.Contains(filename, specifiedChr) {
		return fileInfo{}, fmt.Errorf("not supported without '%s' characters in the file name: %s", specifiedChr, filename)
	}

	ext := filepath.Ext(filename)

	// _Im部分で分割
	parts := strings.Split(filename, specifiedChr)
	if len(parts) < 2 {
		return fileInfo{}, fmt.Errorf("invalid filename format: %s", filename)
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

	// ページ番号を抽出
	pageNum := s.extractPageNumber(leftPart)

	return fileInfo{
		originalPath: filepath.Join(outputDir, filename),
		pageNum:      pageNum,
		imageNum:     imageNum,
		ext:          ext,
	}, nil
}

func (s *ImageExtractionService) sortImageFiles(imageFiles []fileInfo) []fileInfo {
	sort.Slice(imageFiles, func(i, j int) bool {
		if imageFiles[i].pageNum != imageFiles[j].pageNum {
			return imageFiles[i].pageNum < imageFiles[j].pageNum
		}
		return imageFiles[i].imageNum < imageFiles[j].imageNum
	})

	return imageFiles
}

// calculateActualPageNumber は実際のページ番号を計算します
func (s *ImageExtractionService) calculateActualPageNumber(fileInfo fileInfo, index int, startPageOffset int) int {
	if startPageOffset > 0 {
		// 開始ページが指定されている場合、実際のページ番号を使用
		return startPageOffset + index
	} else {
		// 全ページ抽出の場合、pdfcpuが抽出したページ番号または連番を使用
		if fileInfo.pageNum > 0 {
			return fileInfo.pageNum
		} else {
			return index + 1
		}
	}
}

// renameExtractedImagesWithFourDigits は抽出された画像ファイルの名前を4桁連番形式に変更します
func (s *ImageExtractionService) renameExtractedImagesWithFourDigits(outputDir, pdfBaseName string, startPageOffset int) error {
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}

	// PDFファイル名から拡張子を除去
	pdfName := strings.TrimSuffix(pdfBaseName, filepath.Ext(pdfBaseName))

	// ファイルを分析してページ別に整理
	var imageFiles []fileInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := file.Name()
		info, err := s.parseImageFileName(filename, outputDir)
		if err != nil {
			// ファイル名が対象パターンではない場合はスキップ
			continue
		}

		imageFiles = append(imageFiles, info)
	}

	// ページ番号でソート
	imageFiles = s.sortImageFiles(imageFiles)

	// 実際のページ番号に基づいてリネーム
	// startPageOffsetが0の場合は全ページ抽出なので、連番を使用
	// startPageOffsetが1以上の場合は、実際のページ番号を使用
	for i, fileInfo := range imageFiles {
		actualPageNum := s.calculateActualPageNumber(fileInfo, i, startPageOffset)

		newName := fmt.Sprintf("%s_%04d%s", pdfName, actualPageNum, fileInfo.ext)
		newPath := filepath.Join(outputDir, newName)

		if err := os.Rename(fileInfo.originalPath, newPath); err != nil {
			return fmt.Errorf("ファイル名の変更に失敗しました: %w", err)
		}
	}

	return nil
}

func (s *ImageExtractionService) isSupportedFormat(format string) (bool, string, error) {
	// サポートする画像形式を確認
	supportedFormats := map[string]bool{
		"jpg":  true,
		"jpeg": true,
		"png":  true,
		"tiff": true,
		"webp": true,
	}

	lowerFormat := strings.ToLower(format)
	value, exists := supportedFormats[lowerFormat]
	msg := "(サポート形式: jpg, jpeg, png, tiff, webp)"
	if !exists {
		return false, msg, fmt.Errorf("サポートの有無が規定されていない画像形式です: %s %s", format, msg)
	}
	return value, msg, nil
}

// PageRangeInfo はページ範囲の情報を格納する構造体
type PageRangeInfo struct {
	PageSelection []string // pdfcpuに渡すページ選択配列
	Message       string   // ユーザー向け表示メッセージ
}

// GetRangeOfPages はページ範囲の指定処理とメッセージ生成を行います
func (s *ImageExtractionService) GetRangeOfPages(startPage, endPage, totalPages int) PageRangeInfo {
	// ページ選択配列の生成（既存のspecifyRangeOfPages関数を使用）
	pageSelection := s.specifyRangeOfPages(startPage, endPage, totalPages)

	// メッセージの生成
	message := s.generatePageRangeMessage(startPage, endPage, totalPages)

	return PageRangeInfo{
		PageSelection: pageSelection,
		Message:       message,
	}
}

// specifyRangeOfPages はページ範囲の指定処理を行います
func (s *ImageExtractionService) specifyRangeOfPages(startPage, endPage, totalPages int) []string {
	// 0を適切な値に変換
	actualStart := startPage
	actualEnd := endPage

	if actualStart == 0 {
		actualStart = 1 // 最初のページ
	}
	if actualEnd == 0 {
		actualEnd = totalPages // 最後のページ
	}

	// 全ページ抽出の場合（両方とも0）は空配列を返す
	if startPage == 0 && endPage == 0 {
		return []string{}
	}

	return []string{fmt.Sprintf("%d-%d", actualStart, actualEnd)}
}

// generatePageRangeMessage はページ範囲表示用のメッセージを生成します
func (s *ImageExtractionService) generatePageRangeMessage(startPage, endPage, totalPages int) string {
	if startPage == 0 && endPage == 0 {
		return fmt.Sprintf("全ページ (ページ 1 から %d まで)", totalPages)
	}
	if startPage == 0 {
		return fmt.Sprintf("ページ 1 から %d まで", endPage)
	}
	if endPage == 0 {
		return fmt.Sprintf("ページ %d から %d まで (最終ページ)", startPage, totalPages)
	}
	return fmt.Sprintf("ページ %d から %d まで", startPage, endPage)
}

// ExtractToImages はPDFの指定したページ範囲を画像として抽出します
func (s *ImageExtractionService) ExtractToImages(pdfPath, outputDir, imageFormat string, startPage, endPage int) error {
	// 出力ディレクトリが存在しない場合は作成
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	// PDFファイルの存在確認
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return fmt.Errorf("PDFファイルが見つかりません: %s", pdfPath)
	}

	// サポートする画像形式を確認
	isValidFormat, messageAboutFormat, err := s.isSupportedFormat(imageFormat)
	if err != nil {
		return err
	}
	if !isValidFormat {
		return fmt.Errorf("サポートされていない画像形式です: %s %s", imageFormat, messageAboutFormat)
	}

	// PDFの総ページ数を取得
	totalPages, err := s.GetPageCount(pdfPath)
	if err != nil {
		return fmt.Errorf("PDFのページ数取得に失敗しました: %w", err)
	}

	// ページ範囲の指定
	rangeInfo := s.GetRangeOfPages(startPage, endPage, totalPages)

	cfg := api.LoadConfiguration()

	// pdfcpuを使ってPDFページを画像として出力
	// 注意: pdfcpuのバージョンによって利用可能なAPIが異なります
	// ExtractImagesFileを試し、失敗した場合は代替手段を試行
	err = api.ExtractImagesFile(pdfPath, outputDir, rangeInfo.PageSelection, cfg)
	if err != nil {
		return fmt.Errorf("PDFからの画像抽出に失敗しました: %w", err)
	}

	// 抽出された画像ファイルの名前を4桁連番形式に整理
	err = s.renameExtractedImagesWithFourDigits(outputDir, filepath.Base(pdfPath), startPage)
	if err != nil {
		return fmt.Errorf("画像ファイル名の整理に失敗しました: %w", err)
	}

	return nil
}

// GetPageCount はPDFファイルのページ数を取得します
func (s *ImageExtractionService) GetPageCount(pdfPath string) (int, error) {
	// PDFファイルの存在確認
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return 0, fmt.Errorf("PDFファイルが見つかりません: %s", pdfPath)
	}

	// pdfcpuを使ってPDFの情報を取得
	// ReadContextFileは内部で検証も行うため、これだけで十分
	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return 0, fmt.Errorf("PDFファイルの読み込みに失敗しました: %w", err)
	}

	return ctx.PageCount, nil
}

// ValidatePageRange はページ範囲のバリデーションを行います
func (s *ImageExtractionService) ValidatePageRange(startPage, endPage, totalPages int) error {
	// 負の値のチェック
	if startPage < 0 {
		return fmt.Errorf("開始ページは0以上の値を指定してください（指定値: %d）", startPage)
	}
	if endPage < 0 {
		return fmt.Errorf("終了ページは0以上の値を指定してください（指定値: %d）", endPage)
	}

	// 1以上の値が指定された場合のページ範囲チェック（1からtotalPagesの範囲内）
	if startPage > 0 {
		if startPage > totalPages {
			return fmt.Errorf("開始ページがPDFの総ページ数を超えています（指定値: %d, 総ページ数: %d）", startPage, totalPages)
		}
	}

	if endPage > 0 {
		if endPage > totalPages {
			return fmt.Errorf("終了ページがPDFの総ページ数を超えています（指定値: %d, 総ページ数: %d）", endPage, totalPages)
		}
	}

	// 開始ページと終了ページの関係性チェック（両方が1以上の場合）
	if startPage > 0 && endPage > 0 {
		if startPage > endPage {
			return fmt.Errorf("開始ページが終了ページより後になっています（開始: %d, 終了: %d）", startPage, endPage)
		}
	}

	return nil
}
