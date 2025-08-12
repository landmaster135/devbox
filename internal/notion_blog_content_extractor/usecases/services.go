package usecases

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileOperator はファイル操作のインターフェース
type FileOperator interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	WalkDir(root string, fn fs.WalkDirFunc) error
	Stat(name string) (os.FileInfo, error)
}

// DefaultFileOperator は標準のファイル操作実装
type DefaultFileOperator struct{}

func (d *DefaultFileOperator) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (d *DefaultFileOperator) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return os.WriteFile(filename, data, perm)
}

func (d *DefaultFileOperator) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (d *DefaultFileOperator) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}

func (d *DefaultFileOperator) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Service はブログコンテンツ抽出サービス
type Service struct {
	fileOperator FileOperator
}

// NewService は新しいServiceインスタンスを作成する
func NewService() *Service {
	return &Service{
		fileOperator: &DefaultFileOperator{},
	}
}

// NewServiceWithFileOperator はFileOperatorを注入してServiceインスタンスを作成する（テスト用）
func NewServiceWithFileOperator(fileOperator FileOperator) *Service {
	return &Service{
		fileOperator: fileOperator,
	}
}

// ExtractBlogContent は指定されたディレクトリからブログコンテンツを抽出する
func (s *Service) ExtractBlogContent(srcDir, destDir string) (string, error) {
	// ソースディレクトリの存在確認
	if err := s.validateDirectory(srcDir, "ソースディレクトリ"); err != nil {
		return "", err
	}

	// 出力ディレクトリの作成
	if err := s.createDestinationDirectory(destDir); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました: %v", err)
	}

	// コンテンツファイルの検索
	contentFiles, err := s.findContentFiles(srcDir)
	if err != nil {
		return "", fmt.Errorf("コンテンツファイルの検索に失敗しました: %v", err)
	}

	if len(contentFiles) == 0 {
		return "指定されたディレクトリにコンテンツマーカーを含むMarkdownファイルが見つかりませんでした。", nil
	}

	// 各ファイルからコンテンツを抽出
	extractedCount := 0
	var errors []string

	for _, filePath := range contentFiles {
		if err := s.extractContentFromFile(filePath, destDir); err != nil {
			errors = append(errors, fmt.Sprintf("ファイル %s: %v", filepath.Base(filePath), err))
		} else {
			extractedCount++
		}
	}

	// 結果の構築
	result := fmt.Sprintf("処理完了: %d件のファイルからコンテンツを抽出しました。", extractedCount)
	if len(errors) > 0 {
		result += fmt.Sprintf("\n\nエラーが発生したファイル:\n%s", strings.Join(errors, "\n"))
	}

	return result, nil
}

// validateDirectory はディレクトリの存在と種別を確認する
func (s *Service) validateDirectory(dirPath, description string) error {
	info, err := s.fileOperator.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%sが存在しません: %s", description, dirPath)
		}
		return fmt.Errorf("%sの確認に失敗しました: %v", description, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("指定されたパスはディレクトリではありません: %s", dirPath)
	}

	return nil
}

// createDestinationDirectory は出力ディレクトリを作成する
func (s *Service) createDestinationDirectory(destDir string) error {
	return s.fileOperator.MkdirAll(destDir, 0755)
}

// findContentFiles は指定されたディレクトリ内でコンテンツマーカーを含むファイルを検索する
func (s *Service) findContentFiles(srcDir string) ([]string, error) {
	var contentFiles []string
	var totalMdFiles int
	var readableFiles int
	var filesWithMarker int

	// コンテンツマーカーのパターン（改行を含む）
	// # Content\n\n## はじまり\n
	pattern := regexp.MustCompile(`# Content\s*\n\s*\n## はじまり\s*\n`)

	err := s.fileOperator.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// .mdファイルのみを対象とする
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		totalMdFiles++

		// ファイル内容を読み込み
		content, err := s.fileOperator.ReadFile(path)
		if err != nil {
			// 個別ファイルの読み込みエラーは、処理を中止
			return fmt.Errorf("ファイルの読み込みに失敗しました: %s (%v)", path, err)
		}

		readableFiles++

		// パターンマッチング
		if pattern.Match(content) {
			// パターンマッチングの場合は、処理を継続
			contentFiles = append(contentFiles, path)
			filesWithMarker++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("ディレクトリの走査に失敗しました: %v", err)
	}

	// デバッグ情報を出力
	fmt.Printf("デバッグ情報: 総Markdownファイル数=%d, 読み込み可能ファイル数=%d, マーカー付きファイル数=%d\n",
		totalMdFiles, readableFiles, filesWithMarker)

	return contentFiles, nil
}

// extractContentFromFile は指定されたファイルからコンテンツを抽出して保存する
func (s *Service) extractContentFromFile(filePath, destDir string) error {
	// ファイル内容を読み込み
	content, err := s.fileOperator.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ファイルの読み込みに失敗しました: %v", err)
	}

	// コンテンツマーカー以降の内容を抽出
	extractedContent, err := s.extractContentAfterMarker(string(content))
	if err != nil {
		return fmt.Errorf("コンテンツの抽出に失敗しました: %v", err)
	}

	// 出力ファイルパスの構築
	fileName := filepath.Base(filePath)
	outputPath := filepath.Join(destDir, fileName)

	// 抽出したコンテンツを保存
	if err := s.fileOperator.WriteFile(outputPath, []byte(extractedContent), 0644); err != nil {
		return fmt.Errorf("ファイルの保存に失敗しました: %v", err)
	}

	return nil
}

// extractContentAfterMarker はマーカー以降の内容を抽出する
func (s *Service) extractContentAfterMarker(content string) (string, error) {
	// マーカーパターンの定義
	// # Content\n\n## はじまり\n 以降のすべての内容を抽出（マーカーも含む）
	// (?s)フラグで.が改行にもマッチするようにする
	pattern := regexp.MustCompile(`(?s)(# Content\s*\n\s*\n## はじまり\s*\n.*)`)

	matches := pattern.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("指定されたマーカーが見つかりません")
	}

	return matches[1], nil
}
