package repositories

import (
	"github.com/landmaster135/devbox/internal/independencies/file_processor/domain/models"
)

// FileRepository はファイル操作を行うインターフェースです
type FileRepository interface {
	// ReadFile はファイルを読み込み、FileContentを返します
	ReadFile(path string) (*models.FileContent, error)

	// WriteFile はFileContentをファイルに書き込みます
	WriteFile(path string, content *models.FileContent) error

	// FileExists はファイルが存在するかどうかを確認します
	FileExists(path string) bool

	// FindFilesByExt はディレクトリ内の指定された拡張子のファイルのパスリストを返します
	FindFilesByExt(dirPath, ext string) ([]string, error)

	// HasFilesWithExt はディレクトリ内に指定された拡張子のファイルが存在するかどうかを確認します
	HasFilesWithExt(dirPath, ext string) (bool, error)

	// ReadJSONFile はJSONファイルを読み込み、インターフェース型として返します
	ReadJSONFile(path string) (interface{}, error)

	// GetDirectoryPath はファイルパスからディレクトリパスを取得します
	GetDirectoryPath(path string) string

	// CreateDirectory はディレクトリを作成します
	CreateDirectory(dirPath string) error

	// ReadDir はディレクトリ内のファイルとサブディレクトリのエントリを返します
	ReadDir(dirPath string) ([]*models.DirEntry, error)
}
