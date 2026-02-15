package infrastructures

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"
)

func TestSQLiteRepository_ListTables_Normal(t *testing.T) {
	t.Parallel()

	dbPath := createTestDB(t)
	repository := NewSQLiteRepository()

	tables, err := repository.ListTables(context.Background(), dbPath)
	require.NoError(t, err)
	require.Equal(t, []string{"posts", "users"}, tables)
}

func TestSQLiteRepository_ListTables_ErrorFileNotFound(t *testing.T) {
	t.Parallel()

	repository := NewSQLiteRepository()
	_, err := repository.ListTables(context.Background(), filepath.Join(t.TempDir(), "missing.db"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "SQLite ファイルが存在しません")
}

func TestSQLiteRepository_ListTables_ErrorEmptyPath(t *testing.T) {
	t.Parallel()

	repository := NewSQLiteRepository()
	_, err := repository.ListTables(context.Background(), "  ")
	require.Error(t, err)
	require.Contains(t, err.Error(), "dbPath が空です")
}

func TestSQLiteRepository_ListTables_ErrorInvalidDB(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	invalidPath := filepath.Join(tempDir, "invalid.db")
	err := os.WriteFile(invalidPath, []byte("not-a-sqlite"), 0o644)
	require.NoError(t, err)

	repository := NewSQLiteRepository()
	_, listErr := repository.ListTables(context.Background(), invalidPath)
	require.Error(t, listErr)
	require.Contains(t, listErr.Error(), "テーブル一覧の取得に失敗しました")
}

func createTestDB(t *testing.T) string {
	t.Helper()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "sample.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	_, err = db.Exec(`
CREATE TABLE users(id INTEGER PRIMARY KEY, name TEXT);
CREATE TABLE posts(id INTEGER PRIMARY KEY, title TEXT);
`)
	require.NoError(t, err)

	return dbPath
}
