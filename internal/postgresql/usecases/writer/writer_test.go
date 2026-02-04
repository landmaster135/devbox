package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONStreamWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dump.json")

	writer, err := NewJSONStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath)
	require.NoError(t, err)

	rows := []map[string]any{
		{"id": 1, "name": "Alice"},
		{"id": 2, "logged_in_at": time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)},
	}

	assert.NoError(t, writer.WriteBatch(rows))
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

	err = writer.WriteBatch(rows)
	assert.EqualError(t, err, "既にクローズされたライターに書き込めません")
}

func TestJSONStreamWriter_CloseEmpty(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "empty.json")

	writer, err := NewJSONStreamWriter(&MockFileWriter{
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

	writer, err := NewCSVStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, []string{"id", "name", "active"})
	require.NoError(t, err)

	rows := []map[string]any{
		{"id": 1, "name": "Alice", "active": true},
		{"id": 2, "name": "Bob", "active": false},
	}

	assert.NoError(t, writer.WriteBatch(rows))
	assert.Equal(t, 2, writer.RowsWritten())

	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "id,name,active\n1,Alice,true\n2,Bob,false\n", string(data))

	assert.NoError(t, writer.Close())

	err = writer.WriteBatch(rows)
	assert.EqualError(t, err, "既にクローズされたライターに書き込めません")
}

func TestCSVStreamWriter_WithoutHeaders(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "noheaders.csv")

	writer, err := NewCSVStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, nil)
	require.NoError(t, err)

	assert.NoError(t, writer.WriteBatch(nil))
	assert.Equal(t, 0, writer.RowsWritten())

	assert.NoError(t, writer.Close())

	data, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "", string(data))
}

// #==============================================================#
// ##          Stream Writer Tests                               ##
// #==============================================================#

func TestJSONStreamWriter_WriteBatch(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.json")
	writer, err := NewJSONStreamWriter(&DefaultFileWriter{}, filePath)
	assert.NoError(t, err)

	batch1 := []map[string]any{{"id": 1, "name": "John"}}
	batch2 := []map[string]any{{"id": 2, "name": "Jane"}}

	// Act
	assert.NoError(t, writer.WriteBatch(batch1))
	assert.NoError(t, writer.WriteBatch(batch2))
	assert.NoError(t, writer.Close())

	// Assert
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	expected := `[
  {
    "id": 1,
    "name": "John"
  },
  {
    "id": 2,
    "name": "Jane"
  }
]`
	assert.Equal(t, expected, string(data))
	assert.Equal(t, 2, writer.RowsWritten())
}

func TestJSONStreamWriter_CloseWithoutRows(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty.json")
	writer, err := NewJSONStreamWriter(&DefaultFileWriter{}, filePath)
	assert.NoError(t, err)

	// Act
	assert.NoError(t, writer.Close())

	// Assert
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "[]", string(data))
	assert.Equal(t, 0, writer.RowsWritten())
}

func TestCSVStreamWriter_WriteBatch(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.csv")
	headers := []string{"id", "name", "score", "active", "created_at", "payload"}
	writer, err := NewCSVStreamWriter(&DefaultFileWriter{}, filePath, headers)
	assert.NoError(t, err)

	createdAt := time.Date(2024, 1, 2, 3, 4, 5, 6000, time.UTC)
	rows := []map[string]any{
		{
			"id":         1,
			"name":       "John",
			"score":      1.2345,
			"active":     true,
			"created_at": createdAt,
			"payload":    []byte("hi"),
		},
		{
			"id":         2,
			"name":       "Jane",
			"score":      2.5,
			"active":     false,
			"created_at": createdAt.Add(time.Minute),
			"payload":    nil,
		},
	}

	// Act
	assert.NoError(t, writer.WriteBatch(rows))
	assert.NoError(t, writer.Close())

	// Assert
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	expected := strings.Join([]string{
		"id,name,score,active,created_at,payload",
		fmt.Sprintf("1,John,1.2345,true,%s,aGk=", createdAt.Format(time.RFC3339Nano)),
		fmt.Sprintf("2,Jane,2.5,false,%s,", createdAt.Add(time.Minute).Format(time.RFC3339Nano)),
	}, "\n") + "\n"
	assert.Equal(t, expected, string(data))
	assert.Equal(t, 2, writer.RowsWritten())
}
