package repositories

// ISO8601Repository はISO8601形式の日時文字列をUNIXタイムスタンプに変換するためのインターフェースです
type ISO8601Repository interface {
	// ParseISO8601 はISO8601形式の日時文字列をUNIXタイムスタンプに変換します
	ParseISO8601(dateStr string) (int64, error)
}
