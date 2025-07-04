package repositories

// JSONRepository はJSONファイルを操作するためのインターフェースです
type JSONRepository interface {
	// ConvertFile は単一のJSONファイルを処理します
	ConvertFile(filePath, key string, dryRun bool) (bool, error)

	// ProcessJSONData はJSONデータを再帰的に処理します
	ProcessJSONData(data interface{}, targetKey string) (interface{}, bool)

	// FindJSONFiles はディレクトリ内のJSONファイルを検索します
	FindJSONFiles(dirPath string, recursive bool) ([]string, error)
}
