package infrastructures

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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
	return &AttachmentFile{
		Path:        cleanPath,
		Filename:    filepath.Base(cleanPath),
		Content:     data,
		ContentType: detectContentType(cleanPath, data),
	}, nil
}

func detectContentType(filePath string, content []byte) string {
	extType := strings.TrimSpace(mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath))))
	if extType != "" {
		return strings.Split(extType, ";")[0]
	}
	if len(content) > 0 {
		return http.DetectContentType(content)
	}
	return "application/octet-stream"
}
