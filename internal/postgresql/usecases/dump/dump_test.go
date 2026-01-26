package dump

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	model "github.com/landmaster135/devbox/internal/postgresql/domain/model"
	writer "github.com/landmaster135/devbox/internal/postgresql/usecases/writer"
)

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

func createTestTableDumper() (*TableDumper, *MockDatabaseQueryExecutor, *writer.MockFileWriter) {
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockFileWriter := &writer.MockFileWriter{}

	mockExecutor.QueryContextRowsFunc = func(ctx context.Context, query string, args ...any) (model.RowsInterface, error) {
		trimmed := strings.TrimSpace(query)
		if strings.Contains(trimmed, "FROM information_schema.tables") {
			return NewMockRows([]string{"table_name"}, [][]any{{"users"}}), nil
		}
		return nil, fmt.Errorf("unexpected query: %s", trimmed)
	}

	dumper := NewTableDumperWithDependencies(mockExecutor, mockFileWriter, "")

	return dumper, mockExecutor, mockFileWriter
}

// #==============================================================#
// ##          TableDumper Tests                                 ##
// #==============================================================#

func TestNewTableDumper_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockDatabaseQueryExecutor{}

	// Act
	dumper := NewTableDumper(mockExecutor, "")

	// Assert
	assert.NotNil(t, dumper)
	assert.Equal(t, mockExecutor, dumper.executor)
	assert.IsType(t, &writer.DefaultFileWriter{}, dumper.fileWriter)
}

func TestNewTableDumperWithDependencies_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockFileWriter := &writer.MockFileWriter{}

	// Act
	dumper := NewTableDumperWithDependencies(mockExecutor, mockFileWriter, "")

	// Assert
	assert.NotNil(t, dumper)
	assert.Equal(t, mockExecutor, dumper.executor)
	assert.Equal(t, mockFileWriter, dumper.fileWriter)
}

// #==============================================================#
// ##          generateFileName Tests                            ##
// #==============================================================#

func TestTableDumper_generateFileName_JSON(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()

	// Act
	fileName := dumper.generateFileName("users", "json")

	// Assert
	assert.Contains(t, fileName, "users_")
	assert.Contains(t, fileName, ".json")
}

func TestTableDumper_generateFileName_CSV(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()

	// Act
	fileName := dumper.generateFileName("products", "csv")

	// Assert
	assert.Contains(t, fileName, "products_")
	assert.Contains(t, fileName, ".csv")
}

func TestTableDumper_generateFileName_SQL(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()

	// Act
	fileName := dumper.generateFileName("orders", "sql")

	// Assert
	assert.Contains(t, fileName, "orders_")
	assert.Contains(t, fileName, ".sql")
}

func TestTableDumper_generateFileName_UsesTimezone(t *testing.T) {
	// Arrange
	dumper := NewTableDumperWithDependencies(&MockDatabaseQueryExecutor{}, &writer.MockFileWriter{}, "Asia/Tokyo")
	loc, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	before := time.Now().In(loc)
	fileName := dumper.generateFileName("users", "json")
	after := time.Now().In(loc)

	assert.True(t, strings.HasPrefix(fileName, "users_"))
	assert.True(t, strings.HasSuffix(fileName, ".json"))

	timestamp := strings.TrimSuffix(strings.TrimPrefix(fileName, "users_"), ".json")
	generatedTime, err := time.ParseInLocation("20060102_150405", timestamp, loc)
	require.NoError(t, err)

	if generatedTime.Before(before.Add(-2*time.Second)) || generatedTime.After(after.Add(2*time.Second)) {
		t.Fatalf("expected timestamp within timezone window (before=%v after=%v actual=%v)", before, after, generatedTime)
	}
}

// #==============================================================#
// ##          getFileExtension Tests                            ##
// #==============================================================#

func TestTableDumper_getFileExtension_Normal(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()

	tests := []struct {
		format   string
		expected string
	}{
		{"json", "json"},
		{"csv", "csv"},
		{"sql", "sql"},
		{"unknown", "txt"},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			// Act
			extension := dumper.getFileExtension(tt.format)

			// Assert
			assert.Equal(t, tt.expected, extension)
		})
	}
}

// #==============================================================#
// ##          ensureOutputDirectory Tests                       ##
// #==============================================================#

func TestTableDumper_ensureOutputDirectory_CurrentDirectory(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()

	// Act
	err := dumper.ensureOutputDirectory(".")

	// Assert
	assert.NoError(t, err)
	// MkdirAllFuncが呼ばれないことを確認（nilのまま）
	assert.Nil(t, mockFileWriter.MkdirAllFunc)
}

func TestTableDumper_ensureOutputDirectory_EmptyPath(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()

	// Act
	err := dumper.ensureOutputDirectory("")

	// Assert
	assert.NoError(t, err)
	// MkdirAllFuncが呼ばれないことを確認（nilのまま）
	assert.Nil(t, mockFileWriter.MkdirAllFunc)
}

func TestTableDumper_ensureOutputDirectory_CreateDirectory(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	outputPath := "/tmp/dump"

	mockFileWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, outputPath, path)
		assert.Equal(t, os.FileMode(0755), perm)
		return nil
	}

	// Act
	err := dumper.ensureOutputDirectory(outputPath)

	// Assert
	assert.NoError(t, err)
}

func TestTableDumper_ensureOutputDirectory_CreateDirectoryError(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	outputPath := "/tmp/dump"
	expectedError := fmt.Errorf("permission denied")

	mockFileWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, outputPath, path)
		assert.Equal(t, os.FileMode(0755), perm)
		return expectedError
	}

	// Act
	err := dumper.ensureOutputDirectory(outputPath)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// #==============================================================#
// ##          OutputResultIntoFile Tests                        ##
// #==============================================================#

func TestTableDumper_OutputResultIntoFile_Success(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	baseDir := t.TempDir()
	outputPath := filepath.Join(baseDir, "results")
	jsonBytes := []byte(`{"total_tables":1}`)

	mockFileWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, outputPath, path)
		return os.MkdirAll(path, perm)
	}

	var createdPath string
	mockFileWriter.CreateFunc = func(name string) (*os.File, error) {
		createdPath = name
		return os.Create(name)
	}

	// Act
	err := dumper.OutputResultIntoFile(jsonBytes, outputPath, "json")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, createdPath)
	assert.Equal(t, outputPath, filepath.Dir(createdPath))
	written, readErr := os.ReadFile(createdPath)
	assert.NoError(t, readErr)
	assert.Equal(t, jsonBytes, written)
}

func TestTableDumper_OutputResultIntoFile_CreateError(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	outputPath := "/tmp/dump"
	expectedErr := fmt.Errorf("permission denied")

	mockFileWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, outputPath, path)
		return nil
	}

	mockFileWriter.CreateFunc = func(name string) (*os.File, error) {
		return nil, expectedErr
	}

	// Act
	err := dumper.OutputResultIntoFile([]byte(`{}`), outputPath, "json")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "結果ファイル作成エラー")
}

// #==============================================================#
// ##          DumpTable Integration Tests                       ##
// #==============================================================#

func TestTableDumper_DumpTable_ValidationError(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	options := DumpOptions{
		TableName: "", // 無効なテーブル名
		Format:    "json",
	}

	// Act
	result, err := dumper.DumpTable(context.Background(), &options)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "オプション検証エラー")
}

func TestTableDumper_DumpTable_DirectoryCreationError(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/invalid/path",
		Format:     "json",
	}
	expectedError := fmt.Errorf("permission denied")

	// データ取得は成功するがディレクトリ作成で失敗するシナリオ
	// ただし、ディレクトリ作成はデータ取得の前に実行されるため、データ取得のモックは不要
	mockFileWriter.MkdirAllFunc = func(path string, perm os.FileMode) error {
		assert.Equal(t, "/invalid/path", path)
		assert.Equal(t, os.FileMode(0755), perm)
		return expectedError
	}

	// Act
	result, err := dumper.DumpTable(context.Background(), &options)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "出力ディレクトリ作成エラー")
}

// #==============================================================#
// ##          DefaultFileWriter Tests                           ##
// #==============================================================#

func TestDefaultFileWriter_MkdirAll_Normal(t *testing.T) {
	// Arrange
	writer := &writer.DefaultFileWriter{}
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "test", "nested", "dir")

	// Act
	err := writer.MkdirAll(dirPath, 0755)

	// Assert
	assert.NoError(t, err)

	// ディレクトリが作成されたことを確認
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestDefaultFileWriter_Create_Normal(t *testing.T) {
	// Arrange
	writer := &writer.DefaultFileWriter{}
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")

	// Act
	file, err := writer.Create(filePath)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, file)

	// ファイルを閉じる
	file.Close()

	// ファイルが作成されたことを確認
	_, err = os.Stat(filePath)
	assert.NoError(t, err)
}
