package writer

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
)

type testStringer struct {
	value string
}

func (s testStringer) String() string {
	return "stringer:" + s.value
}

func TestFormatCSVValue_VariousTypes(t *testing.T) {
	sampleTime := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "nil", input: nil, want: ""},
		{name: "string", input: "hello", want: "hello"},
		{name: "bytesEmpty", input: []byte{}, want: ""},
		{name: "bytes", input: []byte("hello"), want: base64.StdEncoding.EncodeToString([]byte("hello"))},
		{name: "time", input: sampleTime, want: sampleTime.Format(time.RFC3339Nano)},
		{name: "bool", input: true, want: "true"},
		{name: "int", input: int(-42), want: "-42"},
		{name: "int8", input: int8(-8), want: "-8"},
		{name: "int16", input: int16(-16), want: "-16"},
		{name: "int32", input: int32(-32), want: "-32"},
		{name: "int64", input: int64(-64), want: "-64"},
		{name: "uint", input: uint(42), want: "42"},
		{name: "uint8", input: uint8(8), want: "8"},
		{name: "uint16", input: uint16(16), want: "16"},
		{name: "uint32", input: uint32(32), want: "32"},
		{name: "uint64", input: uint64(64), want: "64"},
		{name: "float32", input: float32(1.5), want: "1.5"},
		{name: "float64", input: float64(2.25), want: "2.25"},
		{name: "stringer", input: testStringer{value: "csv"}, want: "stringer:csv"},
		{name: "default", input: map[string]int{"a": 1}, want: "map[a:1]"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatCSVValue(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFormatSQLValue_VariousTypes(t *testing.T) {
	sampleTime := time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{name: "nil", input: nil, want: "NULL"},
		{name: "string", input: "hello", want: "'hello'"},
		{name: "stringEscaped", input: "O'Reilly", want: "'O''Reilly'"},
		{name: "bytesEmpty", input: []byte{}, want: "decode('', 'hex')"},
		{name: "bytes", input: []byte{0xBE, 0xEF}, want: "decode('beef','hex')"},
		{name: "time", input: sampleTime, want: "'2024-01-02T03:04:05Z'"},
		{name: "boolTrue", input: true, want: "TRUE"},
		{name: "boolFalse", input: false, want: "FALSE"},
		{name: "int", input: int(-42), want: "-42"},
		{name: "int8", input: int8(-8), want: "-8"},
		{name: "int16", input: int16(-16), want: "-16"},
		{name: "int32", input: int32(-32), want: "-32"},
		{name: "int64", input: int64(-64), want: "-64"},
		{name: "uint", input: uint(42), want: "42"},
		{name: "uint8", input: uint8(8), want: "8"},
		{name: "uint16", input: uint16(16), want: "16"},
		{name: "uint32", input: uint32(32), want: "32"},
		{name: "uint64", input: uint64(64), want: "64"},
		{name: "float32", input: float32(1.5), want: "1.5"},
		{name: "float64", input: float64(2.25), want: "2.25"},
		{name: "stringer", input: testStringer{value: "sql"}, want: "'stringer:sql'"},
		{name: "default", input: map[string]int{"a": 1}, want: "'map[a:1]'"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatSQLValue(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSQLStreamWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "dump.sql")

	writer, err := NewSQLStreamWriter(&MockFileWriter{
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

	assert.NoError(t, writer.WriteBatch(rows))
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

	err = writer.WriteBatch(rows)
	assert.EqualError(t, err, "既にクローズされたライターに書き込めません")
}

func TestNewSQLStreamWriter_InvalidTable(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "invalid.sql")

	writer, err := NewSQLStreamWriter(&MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}, filePath, "invalid table", []string{"id"})

	assert.Error(t, err)
	assert.Nil(t, writer)
}

func TestSQLStreamWriter_WriteBatch(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "users.sql")
	columns := []string{"id", "name", "active", "created_at", "payload"}
	writer, err := NewSQLStreamWriter(&DefaultFileWriter{}, filePath, "users", columns)
	assert.NoError(t, err)

	createdAt := time.Date(2024, 1, 2, 3, 4, 5, 6000, time.UTC)
	rows := []map[string]any{
		{
			"id":         1,
			"name":       "John",
			"active":     true,
			"created_at": createdAt,
			"payload":    []byte{0xDE, 0xAD, 0xBE, 0xEF},
		},
		{
			"id":         2,
			"name":       "Jane",
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
	content := string(data)
	assert.Contains(t, content, "-- Table dump for users")
	assert.Contains(t, content, "INSERT INTO \"users\"")
	assert.Contains(t, content, "'John'")
	assert.Contains(t, content, "TRUE")
	assert.Contains(t, content, "FALSE")
	assert.Contains(t, content, createdAt.Format(time.RFC3339Nano))
	assert.Contains(t, content, "decode('deadbeef','hex')")
	assert.Contains(t, content, "NULL")
	assert.Equal(t, 2, writer.RowsWritten())
}

func TestSQLStreamWriter_CloseWithoutRows(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty.sql")
	columns := []string{"id"}
	writer, err := NewSQLStreamWriter(&DefaultFileWriter{}, filePath, "users", columns)
	assert.NoError(t, err)

	// Act
	assert.NoError(t, writer.Close())

	// Assert
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "-- No data to export\n", string(data))
	assert.Equal(t, 0, writer.RowsWritten())
}
