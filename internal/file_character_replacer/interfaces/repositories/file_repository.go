package repositories

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/file_character_replacer/domain"
)

// FileRepositoryImpl はFileRepositoryインターフェースの具象実装です
type FileRepositoryImpl struct {
	encodingConverter domain.EncodingConverter
}

// NewFileRepository は新しいFileRepositoryImplを作成します
func NewFileRepository(encodingConverter domain.EncodingConverter) domain.FileRepository {
	return &FileRepositoryImpl{
		encodingConverter: encodingConverter,
	}
}

// ReadFile はファイルを読み込み、指定されたエンコーディングでUTF-8文字列として返します
func (r *FileRepositoryImpl) ReadFile(path string, encoding domain.EncodingType) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("ファイルを開けませんでした: %s, エラー: %w", path, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("ファイルを読み込めませんでした: %s, エラー: %w", path, err)
	}

	// UTF-8の場合はそのまま返す
	if encoding == domain.EncodingUTF8 {
		return string(content), nil
	}

	// その他のエンコーディングの場合は変換
	utf8Content, err := r.encodingConverter.ConvertToUTF8(content, encoding)
	if err != nil {
		return "", fmt.Errorf("文字エンコーディング変換に失敗しました: %s, エラー: %w", path, err)
	}

	return utf8Content, nil
}

// WriteFile はUTF-8文字列を指定されたエンコーディングでファイルに書き込みます
func (r *FileRepositoryImpl) WriteFile(path string, content string, encoding domain.EncodingType) error {
	var writeContent []byte
	var err error

	// UTF-8の場合はそのまま書き込み
	if encoding == domain.EncodingUTF8 {
		writeContent = []byte(content)
	} else {
		// その他のエンコーディングの場合は変換
		writeContent, err = r.encodingConverter.ConvertFromUTF8(content, encoding)
		if err != nil {
			return fmt.Errorf("文字エンコーディング変換に失敗しました: %s, エラー: %w", path, err)
		}
	}

	err = os.WriteFile(path, writeContent, 0644)
	if err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました: %s, エラー: %w", path, err)
	}

	return nil
}

// ListFiles はディレクトリ内のファイル一覧を取得します
func (r *FileRepositoryImpl) ListFiles(dirPath string, recursive bool) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("ディレクトリの再帰的走査に失敗しました: %s, エラー: %w", dirPath, err)
		}
	} else {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, fmt.Errorf("ディレクトリの読み込みに失敗しました: %s, エラー: %w", dirPath, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				files = append(files, filepath.Join(dirPath, entry.Name()))
			}
		}
	}

	return files, nil
}

// CreateBackup はファイルのバックアップを作成します
func (r *FileRepositoryImpl) CreateBackup(filePath string, backupDir string) error {
	timestamp := time.Now().Format("20060102_150405")

	var backupPath string
	if backupDir == "" {
		// バックアップディレクトリが指定されていない場合は従来通り
		backupPath = fmt.Sprintf("%s.backup_%s", filePath, timestamp)
	} else {
		// バックアップディレクトリが指定されている場合
		// 相対パスを計算してディレクトリ構造を保持
		fileName := filepath.Base(filePath)
		dirPath := filepath.Dir(filePath)

		// バックアップディレクトリ内に元のディレクトリ構造を再現
		backupSubDir := filepath.Join(backupDir, dirPath)
		if err := os.MkdirAll(backupSubDir, 0755); err != nil {
			return fmt.Errorf("バックアップディレクトリの作成に失敗しました: %s, エラー: %w", backupSubDir, err)
		}

		backupPath = filepath.Join(backupSubDir, fmt.Sprintf("%s.backup_%s", fileName, timestamp))
	}

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("バックアップ元ファイルを開けませんでした: %s, エラー: %w", filePath, err)
	}
	defer sourceFile.Close()

	backupFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("バックアップファイルを作成できませんでした: %s, エラー: %w", backupPath, err)
	}
	defer backupFile.Close()

	_, err = io.Copy(backupFile, sourceFile)
	if err != nil {
		return fmt.Errorf("バックアップファイルのコピーに失敗しました: %s, エラー: %w", backupPath, err)
	}

	return nil
}

// FileExists はファイルまたはディレクトリが存在するかを確認します
func (r *FileRepositoryImpl) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// IsDirectory はパスがディレクトリかどうかを確認します
func (r *FileRepositoryImpl) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// GetFileInfo はファイル情報を取得します
func (r *FileRepositoryImpl) GetFileInfo(path string) (*domain.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("ファイル情報の取得に失敗しました: %s, エラー: %w", path, err)
	}

	fileInfo := &domain.FileInfo{
		Path:     path,
		IsDir:    info.IsDir(),
		Size:     info.Size(),
		Encoding: domain.EncodingUTF8, // デフォルトはUTF-8
	}

	// ファイルの場合、エンコーディングを推測
	if !info.IsDir() {
		content, err := os.ReadFile(path)
		if err == nil && len(content) > 0 {
			if detectedEncoding, err := r.encodingConverter.DetectEncoding(content); err == nil {
				fileInfo.Encoding = detectedEncoding
			}
		}
	}

	return fileInfo, nil
}

// IsTextFile はファイルがテキストファイルかどうかを判定します
func (r *FileRepositoryImpl) IsTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
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
