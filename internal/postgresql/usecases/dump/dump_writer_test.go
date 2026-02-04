package dump

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	writer "github.com/landmaster135/devbox/internal/postgresql/usecases/writer"
)

func TestTableDumper_newStreamWriter_SupportedFormats(t *testing.T) {
	dir := t.TempDir()
	mockWriter := &writer.MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}

	dumper := NewTableDumperWithDependencies(&MockDatabaseQueryExecutor{}, mockWriter, "")
	columns := []string{"id", "name"}

	jsonWriter, err := dumper.newStreamWriter("json", filepath.Join(dir, "data.json"), "users", columns)
	require.NoError(t, err)
	require.IsType(t, &writer.JSONStreamWriter{}, jsonWriter)
	require.NoError(t, jsonWriter.Close())

	csvWriter, err := dumper.newStreamWriter("csv", filepath.Join(dir, "data.csv"), "users", columns)
	require.NoError(t, err)
	require.IsType(t, &writer.CSVStreamWriter{}, csvWriter)
	require.NoError(t, csvWriter.Close())

	sqlWriter, err := dumper.newStreamWriter("sql", filepath.Join(dir, "data.sql"), "public.users", columns)
	require.NoError(t, err)
	require.IsType(t, &writer.SQLStreamWriter{}, sqlWriter)
	require.NoError(t, sqlWriter.Close())
}

func TestTableDumper_newStreamWriter_InvalidFormat(t *testing.T) {
	dir := t.TempDir()
	mockWriter := &writer.MockFileWriter{
		CreateFunc: func(name string) (*os.File, error) {
			return os.Create(name)
		},
	}

	dumper := NewTableDumperWithDependencies(&MockDatabaseQueryExecutor{}, mockWriter, "")

	_, err := dumper.newStreamWriter("unsupported", filepath.Join(dir, "data.bin"), "users", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サポートされていないフォーマット")
}
