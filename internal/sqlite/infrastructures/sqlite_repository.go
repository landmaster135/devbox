package infrastructures

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

const sqliteDriverName = "sqlite"

// SQLiteRepository は SQLite ファイルを読み取る repository です。
type SQLiteRepository struct {
	driverName string
}

// NewSQLiteRepository は SQLiteRepository を生成します。
func NewSQLiteRepository() *SQLiteRepository {
	return &SQLiteRepository{driverName: sqliteDriverName}
}

// ListTables は SQLite ファイルからテーブル一覧を返します。
func (r *SQLiteRepository) ListTables(ctx context.Context, dbPath string) ([]string, error) {
	trimmed := strings.TrimSpace(dbPath)
	if trimmed == "" {
		return nil, fmt.Errorf("dbPath が空です")
	}

	if _, err := os.Stat(trimmed); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("SQLite ファイルが存在しません: %s", trimmed)
		}
		return nil, fmt.Errorf("SQLite ファイルの確認に失敗しました: %w", err)
	}

	db, err := sql.Open(r.driverName, trimmed)
	if err != nil {
		return nil, fmt.Errorf("SQLite への接続に失敗しました: %w", err)
	}
	defer db.Close()

	query := `
SELECT name
FROM sqlite_master
WHERE type = 'table' AND name NOT LIKE 'sqlite_%'
ORDER BY name
`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("テーブル一覧の取得に失敗しました: %w", err)
	}
	defer rows.Close()

	tables := make([]string, 0)
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("テーブル名の読み取りに失敗しました: %w", err)
		}
		tables = append(tables, tableName)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("テーブル一覧の走査に失敗しました: %w", err)
	}

	return tables, nil
}
