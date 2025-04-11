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

// ISO8601Repository はISO8601形式の日時文字列をUNIXタイムスタンプに変換するためのインターフェースです
type ISO8601Repository interface {
	// ParseISO8601 はISO8601形式の日時文字列をUNIXタイムスタンプに変換します
	ParseISO8601(dateStr string) (int64, error)
}
