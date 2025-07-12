package domain

// FileRepository はファイル操作の抽象化インターフェースです
type FileRepository interface {
	// ReadFile はファイルを読み込み、指定されたエンコーディングでUTF-8文字列として返します
	ReadFile(path string, encoding EncodingType) (string, error)

	// WriteFile はUTF-8文字列を指定されたエンコーディングでファイルに書き込みます
	WriteFile(path string, content string, encoding EncodingType) error

	// ListFiles はディレクトリ内のファイル一覧を取得します
	ListFiles(dirPath string, recursive bool) ([]string, error)

	// CreateBackup はファイルのバックアップを作成します
	CreateBackup(filePath string, backupDir string) error

	// FileExists はファイルまたはディレクトリが存在するかを確認します
	FileExists(path string) bool

	// IsDirectory はパスがディレクトリかどうかを確認します
	IsDirectory(path string) bool

	// GetFileInfo はファイル情報を取得します
	GetFileInfo(path string) (*FileInfo, error)
}

// EncodingConverter は文字エンコーディング変換の抽象化インターフェースです
type EncodingConverter interface {
	// ConvertToUTF8 は指定されたエンコーディングのバイト列をUTF-8文字列に変換します
	ConvertToUTF8(content []byte, encoding EncodingType) (string, error)

	// ConvertFromUTF8 はUTF-8文字列を指定されたエンコーディングのバイト列に変換します
	ConvertFromUTF8(content string, encoding EncodingType) ([]byte, error)

	// DetectEncoding はバイト列から文字エンコーディングを推測します
	DetectEncoding(content []byte) (EncodingType, error)
}
