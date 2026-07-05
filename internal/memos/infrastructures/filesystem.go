package infrastructures

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var contentTypeByExtension = map[string]string{
	".md":       "text/markdown",
	".markdown": "text/markdown",
}

// AttachmentFile は添付対象ファイルの読み込み結果。
type AttachmentFile struct {
	Path        string
	Filename    string
	Content     []byte
	ContentType string
}

// FileSystem は Memos 用のファイル操作契約。
type FileSystem interface {
	ReadFile(filePath string) ([]byte, error)
	ReadAttachmentFile(filePath string) (*AttachmentFile, error)
	EnsureDirectory(dirPath string) (string, error)
	WriteFile(filePath string, content []byte) error
}

// OSFileSystem は実ファイルシステム実装。
type OSFileSystem struct{}

// NewOSFileSystem は OSFileSystem を生成する。
func NewOSFileSystem() *OSFileSystem {
	return &OSFileSystem{}
}

// ReadFile は指定ファイルを読み込む。
func (fs *OSFileSystem) ReadFile(filePath string) ([]byte, error) {
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return nil, fmt.Errorf("filePath が空です")
	}

	data, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました (%s): %w", cleanPath, err)
	}
	return data, nil
}

// ReadAttachmentFile は添付用にファイル内容と MIME type を返す。
func (fs *OSFileSystem) ReadAttachmentFile(filePath string) (*AttachmentFile, error) {
	data, err := fs.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	cleanPath := strings.TrimSpace(filePath)
	contentType, err := normalizeContentType(detectContentType(cleanPath, data))
	if err != nil {
		return nil, fmt.Errorf("MIME type の判定に失敗しました (%s): %w", cleanPath, err)
	}

	return &AttachmentFile{
		Path:        cleanPath,
		Filename:    filepath.Base(cleanPath),
		Content:     data,
		ContentType: contentType,
	}, nil
}

// EnsureDirectory は指定パスが存在するディレクトリであることを確認し、絶対パスを返す。
func (fs *OSFileSystem) EnsureDirectory(dirPath string) (string, error) {
	cleanPath := strings.TrimSpace(dirPath)
	if cleanPath == "" {
		return "", fmt.Errorf("dirPath が空です")
	}

	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", fmt.Errorf("ディレクトリパスの絶対パス変換に失敗しました (%s): %w", cleanPath, err)
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", fmt.Errorf("ディレクトリの確認に失敗しました (%s): %w", absPath, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("指定されたパスはディレクトリではありません: %s", absPath)
	}
	return absPath, nil
}

// WriteFile は指定ファイルへ内容を書き込む。
func (fs *OSFileSystem) WriteFile(filePath string, content []byte) error {
	cleanPath := strings.TrimSpace(filePath)
	if cleanPath == "" {
		return fmt.Errorf("filePath が空です")
	}
	if err := os.WriteFile(cleanPath, content, 0o600); err != nil {
		return fmt.Errorf("ファイルの書き込みに失敗しました (%s): %w", cleanPath, err)
	}
	return nil
}

func detectContentType(filePath string, content []byte) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if mappedType, ok := contentTypeByExtension[ext]; ok {
		return mappedType
	}

	extType := strings.TrimSpace(mime.TypeByExtension(ext))
	if extType != "" {
		return extType
	}
	if len(content) > 0 {
		return http.DetectContentType(content)
	}
	return "application/octet-stream"
}

func normalizeContentType(contentType string) (string, error) {
	cleanType := strings.TrimSpace(contentType)
	if cleanType == "" {
		return "", fmt.Errorf("MIME type が空です")
	}

	mediaType, _, err := mime.ParseMediaType(cleanType)
	if err != nil {
		return "", fmt.Errorf("MIME type の解析に失敗しました: %w", err)
	}
	normalizedType := strings.TrimSpace(mediaType)
	if normalizedType == "" {
		return "", fmt.Errorf("MIME type が空です")
	}
	return normalizedType, nil
}
