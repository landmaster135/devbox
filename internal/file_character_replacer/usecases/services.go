package usecases

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/landmaster135/devbox/internal/file_character_replacer/config"
	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
	"github.com/landmaster135/devbox/internal/file_character_replacer/interfaces/repositories"
)

// FileReplacerService はファイル文字列置換のビジネスロジックを実装します
type FileReplacerService struct {
	fileRepo          domain.FileRepository
	encodingConverter domain.EncodingConverter
	config            *domain.ReplacementConfig
}

// NewFileReplacerService は新しいFileReplacerServiceを作成します
func NewFileReplacerService() *FileReplacerService {
	// 依存関係を内部で注入
	encodingConverter := repositories.NewEncodingConverter()
	fileRepo := repositories.NewFileRepository(encodingConverter)

	return &FileReplacerService{
		fileRepo:          fileRepo,
		encodingConverter: encodingConverter,
		config:            nil, // 後でSetConfigで設定
	}
}

// NewFileReplacerServiceWithConfig はコマンドライン引数を解析してサービスを作成します
func NewFileReplacerServiceWithConfig() (*FileReplacerService, error) {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		return nil, err
	}

	// ドメイン設定に変換
	replacementConfig := cfg.ToReplacementConfig()

	// 依存関係を内部で注入
	encodingConverter := repositories.NewEncodingConverter()
	fileRepo := repositories.NewFileRepository(encodingConverter)

	return &FileReplacerService{
		fileRepo:          fileRepo,
		encodingConverter: encodingConverter,
		config:            replacementConfig,
	}, nil
}

// SetConfig は設定を設定します
func (s *FileReplacerService) SetConfig(config *domain.ReplacementConfig) {
	s.config = config
}

// ReplaceStrings は文字列置換を実行します
func (s *FileReplacerService) ReplaceStrings() (*domain.FileProcessResult, error) {
	result := &domain.FileProcessResult{
		ProcessedFiles: 0,
		ReplacedCount:  0,
		Errors:         []error{},
		Messages:       []string{},
	}

	// 設定の妥当性を検証
	if err := s.config.Validate(); err != nil {
		result.AddError(err)
		return result, err
	}

	// 対象パスが存在するかチェック
	if !s.fileRepo.FileExists(s.config.Target) {
		err := fmt.Errorf("対象パスが存在しません: %s", s.config.Target)
		result.AddError(err)
		return result, err
	}

	// ディレクトリかファイルかで処理を分岐
	if s.fileRepo.IsDirectory(s.config.Target) {
		return s.replaceInDirectory(result)
	} else {
		return s.replaceInFile(s.config.Target, result)
	}
}

// replaceInDirectory はディレクトリ内のファイルを処理します
func (s *FileReplacerService) replaceInDirectory(result *domain.FileProcessResult) (*domain.FileProcessResult, error) {
	files, err := s.fileRepo.ListFiles(s.config.Target, s.config.Recursive)
	if err != nil {
		result.AddError(err)
		return result, err
	}

	result.AddMessage(fmt.Sprintf("ディレクトリ内のファイル数: %d", len(files)))

	for _, filePath := range files {
		// テキストファイルのみを処理
		if !s.isTextFile(filePath) {
			continue
		}

		_, err := s.replaceInFile(filePath, result)
		if err != nil {
			log.Printf("ファイル処理中にエラーが発生しました: %s, エラー: %v", filePath, err)
			// エラーが発生してもその他のファイルの処理は続行
		}
	}

	return result, nil
}

// replaceInFile は単一ファイルの文字列置換を実行します
func (s *FileReplacerService) replaceInFile(filePath string, result *domain.FileProcessResult) (*domain.FileProcessResult, error) {
	// ファイル情報を取得
	fileInfo, err := s.fileRepo.GetFileInfo(filePath)
	if err != nil {
		result.AddError(err)
		return result, err
	}

	// 使用するエンコーディングを決定
	encoding := s.config.Encoding
	if encoding == "" {
		encoding = fileInfo.Encoding // ファイルから推測されたエンコーディングを使用
	}

	// ファイルを読み込み
	content, err := s.fileRepo.ReadFile(filePath, encoding)
	if err != nil {
		result.AddError(err)
		return result, err
	}

	// 文字列置換を実行
	originalContent := content
	replacedContent := strings.ReplaceAll(content, s.config.From, s.config.To)

	// 置換が発生したかチェック
	if originalContent == replacedContent {
		// 置換が発生しなかった場合
		return result, nil
	}

	// 置換回数をカウント
	replacedCount := strings.Count(originalContent, s.config.From)
	result.ReplacedCount += replacedCount

	result.AddMessage(fmt.Sprintf("ファイル: %s, 置換回数: %d", filePath, replacedCount))

	// ドライランの場合は実際の書き込みをスキップ
	if s.config.DryRun {
		result.AddMessage(fmt.Sprintf("[ドライラン] %s の置換をスキップしました", filePath))
		result.ProcessedFiles++
		return result, nil
	}

	// バックアップを作成
	if s.config.Backup {
		if err := s.fileRepo.CreateBackup(filePath); err != nil {
			result.AddError(fmt.Errorf("バックアップ作成に失敗しました: %s, エラー: %w", filePath, err))
			return result, err
		}
		result.AddMessage(fmt.Sprintf("バックアップを作成しました: %s", filePath))
	}

	// ファイルに書き込み
	if err := s.fileRepo.WriteFile(filePath, replacedContent, encoding); err != nil {
		result.AddError(err)
		return result, err
	}

	result.ProcessedFiles++
	return result, nil
}

// isTextFile はファイルがテキストファイルかどうかを判定します
func (s *FileReplacerService) isTextFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	textExtensions := []string{
		".txt", ".md", ".go", ".py", ".js", ".ts", ".html", ".css", ".xml", ".json",
		".yaml", ".yml", ".toml", ".ini", ".cfg", ".conf", ".log", ".sql", ".sh",
		".bat", ".ps1", ".php", ".rb", ".java", ".c", ".cpp", ".h", ".hpp",
		".cs", ".vb", ".pl", ".r", ".scala", ".kt", ".swift", ".dart", ".rs",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}

	return false
}

// GetSummary は処理結果のサマリーを取得します
func (s *FileReplacerService) GetSummary(result *domain.FileProcessResult) string {
	var summary strings.Builder

	summary.WriteString("=== ファイル文字列置換結果 ===\n")
	summary.WriteString(fmt.Sprintf("処理されたファイル数: %d\n", result.ProcessedFiles))
	summary.WriteString(fmt.Sprintf("置換された箇所数: %d\n", result.ReplacedCount))

	if s.config.DryRun {
		summary.WriteString("モード: ドライラン（実際の変更は行われていません）\n")
	}

	if len(result.Messages) > 0 {
		summary.WriteString("\n=== 処理詳細 ===\n")
		for _, msg := range result.Messages {
			summary.WriteString(fmt.Sprintf("%s\n", msg))
		}
	}

	if result.HasErrors() {
		summary.WriteString("\n=== エラー ===\n")
		for _, err := range result.Errors {
			summary.WriteString(fmt.Sprintf("エラー: %v\n", err))
		}
	}

	return summary.String()
}
