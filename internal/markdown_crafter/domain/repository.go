package domain

// Repository は markdown-crafter のファイル入出力を抽象化します。
type Repository interface {
	ReadFile(filePath string) (string, error)
	WriteFile(filePath string, content string) error
	CreateDir(dirPath string) error
}
