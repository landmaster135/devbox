package dump

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableDumper_newStreamWriter_SupportedFormats(t *testing.T) {
	dir := t.TempDir()
	mockWriter := &MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}

	dumper := NewTableDumperWithDependencies(&MockDatabaseQueryExecutor{}, mockWriter)
	columns := []string{"id", "name"}

	jsonWriter, err := dumper.newStreamWriter("json", filepath.Join(dir, "data.json"), "users", columns)
	require.NoError(t, err)
	require.IsType(t, &jsonStreamWriter{}, jsonWriter)
	require.NoError(t, jsonWriter.Close())

	csvWriter, err := dumper.newStreamWriter("csv", filepath.Join(dir, "data.csv"), "users", columns)
	require.NoError(t, err)
	require.IsType(t, &csvStreamWriter{}, csvWriter)
	require.NoError(t, csvWriter.Close())

	sqlWriter, err := dumper.newStreamWriter("sql", filepath.Join(dir, "data.sql"), "public.users", columns)
	require.NoError(t, err)
	require.IsType(t, &sqlStreamWriter{}, sqlWriter)
	require.NoError(t, sqlWriter.Close())
}

func TestTableDumper_newStreamWriter_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	mockWriter := &MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}

	dumper := NewTableDumperWithDependencies(&MockDatabaseQueryExecutor{}, mockWriter)

	_, err := dumper.newStreamWriter("unsupported", filepath.Join(dir, "data.bin"), "users", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サポートされていないフォーマット")
}

func TestJSONStreamWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dump.json")

	writer, err := newJSONStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath)
	require.NoError(t, err)

	rows := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "logged_in_at": time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)},
	}

	assert.NoError(t, writer.writeBatch(rows))
	assert.Equal(t, 2, writer.RowsWritten())

	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	content := string(data)
	assert.True(t, strings.HasPrefix(content, "[\n"))
	assert.Contains(t, content, "\"name\": \"Alice\"")
	assert.Contains(t, content, "\"logged_in_at\": \"2024-01-02T03:04:05Z\"")
	assert.True(t, strings.HasSuffix(strings.TrimSpace(content), "]"))

	assert.NoError(t, writer.Close())

	err = writer.writeBatch(rows)
	assert.EqualError(t, err, "既にクローズされたライターに書き込めません")
}

func TestJSONStreamWriter_CloseEmpty(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "empty.json")

	writer, err := newJSONStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath)
	require.NoError(t, err)

	assert.Equal(t, 0, writer.RowsWritten())
	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "[]", string(data))
}

func TestCSVStreamWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dump.csv")

	writer, err := newCSVStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, []string{"id", "name", "active"})
	require.NoError(t, err)

	rows := []map[string]any{
		{"id": 1, "name": "Alice", "active": true},
		{"id": 2, "name": "Bob", "active": false},
	}

	assert.NoError(t, writer.writeBatch(rows))
	assert.Equal(t, 2, writer.RowsWritten())

	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "id,name,active\n1,Alice,true\n2,Bob,false\n", string(data))

	assert.NoError(t, writer.Close())

	err = writer.writeBatch(rows)
	assert.EqualError(t, err, "既にクローズされたライターに書き込めません")
}

func TestCSVStreamWriter_WithoutHeaders(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "noheaders.csv")

	writer, err := newCSVStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, nil)
	require.NoError(t, err)

	assert.NoError(t, writer.writeBatch(nil))
	assert.Equal(t, 0, writer.RowsWritten())

	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "", string(data))
}

func TestSQLStreamWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dump.sql")

	writer, err := newSQLStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, "public.users", []string{"id", "name", "created_at", "payload", "active", "score", "amount", "note"})
	require.NoError(t, err)

	rows := []map[string]any{
		{
			"id":         int64(1),
			"name":       "Alice",
			"created_at": time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC),
			"payload":    []byte{0xBE, 0xEF},
			"active":     true,
			"score":      float64(9.5),
			"amount":     int(42),
			"note":       testStringer{value: "memo"},
		},
	}

	assert.NoError(t, writer.writeBatch(rows))
	assert.Equal(t, 1, writer.RowsWritten())

	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	content := string(data)

	query := "INSERT INTO \"" + model.DefaultTableSchema + "\".\"users\" (\"id\", \"name\", \"created_at\", \"payload\", \"active\", \"score\", \"amount\", \"note\") VALUES (1, "
	assert.Contains(t, content, query)
	assert.Contains(t, content, "decode('beef','hex')")
	assert.Contains(t, content, "TRUE")

	assert.NoError(t, writer.Close())

	err = writer.writeBatch(rows)
	assert.EqualError(t, err, "既にクローズされたライターに書き込めません")
}

func TestNewSQLStreamWriter_InvalidTable(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "invalid.sql")

	writer, err := newSQLStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, "invalid table", []string{"id"})

	assert.Error(t, err)
	assert.Nil(t, writer)
}
