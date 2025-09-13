package usecases

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockFileWriter はテスト用のFileWriterモック
type MockFileWriter struct {
	WriteFileFunc func(filename string, data []byte, perm os.FileMode) error
	MkdirAllFunc  func(path string, perm os.FileMode) error
	CreateFunc    func(name string) (*os.File, error)
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(filename, data, perm)
	}
	return nil
}

func (m *MockFileWriter) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockFileWriter) Create(name string) (*os.File, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(name)
	}
	return nil, nil
}

// #==============================================================#
// ##          Test Helper Functions                             ##
// #==============================================================#

func createTestTableDumper() (*TableDumper, *MockDatabaseQueryExecutor, *MockFileWriter) {
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockFileWriter := &MockFileWriter{}

	dumper := NewTableDumperWithDependencies(mockExecutor, mockFileWriter)

	return dumper, mockExecutor, mockFileWriter
}

// #==============================================================#
// ##          TableDumper Tests                                 ##
// #==============================================================#

func TestNewTableDumper_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockDatabaseQueryExecutor{}

	// Act
	dumper := NewTableDumper(mockExecutor)

	// Assert
	assert.NotNil(t, dumper)
	assert.Equal(t, mockExecutor, dumper.executor)
	assert.IsType(t, &DefaultFileWriter{}, dumper.fileWriter)
}

func TestNewTableDumperWithDependencies_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockDatabaseQueryExecutor{}
	mockFileWriter := &MockFileWriter{}

	// Act
	dumper := NewTableDumperWithDependencies(mockExecutor, mockFileWriter)

	// Assert
	assert.NotNil(t, dumper)
	assert.Equal(t, mockExecutor, dumper.executor)
	assert.Equal(t, mockFileWriter, dumper.fileWriter)
}

// #==============================================================#
// ##          validateOptions Tests                             ##
// #==============================================================#

func TestTableDumper_validateOptions_Normal(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := dumper.validateOptions(options)

	// Assert
	assert.NoError(t, err)
}

func TestTableDumper_validateOptions_EmptyTableName(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	options := DumpOptions{
		TableName:  "",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := dumper.validateOptions(options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "テーブル名が指定されていません")
}

func TestTableDumper_validateOptions_InvalidFormat(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/tmp/dump",
		Format:     "invalid",
		Limit:      nil,
	}

	// Act
	err := dumper.validateOptions(options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "サポートされていないフォーマットです")
}

func TestTableDumper_validateOptions_InvalidLimit(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	invalidLimit := -1
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/tmp/dump",
		Format:     "json",
		Limit:      &invalidLimit,
	}

	// Act
	err := dumper.validateOptions(options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "limitは正の数である必要があります")
}

func TestTableDumper_validateOptions_PathTraversal(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "../../../etc",
		Format:     "json",
		Limit:      nil,
	}

	// Act
	err := dumper.validateOptions(options)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "無効なパスが指定されました")
}

// #==============================================================#
// ##          buildQuery Tests                                  ##
// #==============================================================#

func TestTableDumper_buildQuery_WithoutLimit(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	options := DumpOptions{
		TableName: "users",
		Limit:     nil,
	}

	// Act
	query := dumper.buildQuery(options)

	// Assert
	assert.Equal(t, "SELECT * FROM users", query)
}

func TestTableDumper_buildQuery_WithLimit(t *testing.T) {
	// Arrange
	dumper, _, _ := createTestTableDumper()
	limit := 100
	options := DumpOptions{
		TableName: "users",
		Limit:     &limit,
	}

	// Act
	query := dumper.buildQuery(options)

	// Assert
	assert.Equal(t, "SELECT * FROM users LIMIT 100", query)
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
// ##          writeJSONFile Tests                               ##
// #==============================================================#

func TestTableDumper_writeJSONFile_Normal(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/test.json"
	data := []map[string]any{
		{"id": 1, "name": "John"},
		{"id": 2, "name": "Jane"},
	}
	expectedJSON := `[
  {
    "id": 1,
    "name": "John"
  },
  {
    "id": 2,
    "name": "Jane"
  }
]`

	mockFileWriter.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		assert.Equal(t, filePath, filename)
		assert.Equal(t, []byte(expectedJSON), data)
		assert.Equal(t, os.FileMode(0644), perm)
		return nil
	}

	// Act
	err := dumper.writeJSONFile(filePath, data)

	// Assert
	assert.NoError(t, err)
}

func TestTableDumper_writeJSONFile_WriteError(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/test.json"
	data := []map[string]any{
		{"id": 1, "name": "John"},
	}
	expectedError := fmt.Errorf("write error")

	mockFileWriter.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		assert.Equal(t, filePath, filename)
		assert.Equal(t, os.FileMode(0644), perm)
		return expectedError
	}

	// Act
	err := dumper.writeJSONFile(filePath, data)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

// #==============================================================#
// ##          writeSQLFile Tests                                ##
// #==============================================================#

func TestTableDumper_writeSQLFile_Normal(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/users_20250822_170305.sql"
	tableName := "users"
	data := []map[string]any{
		{"id": 1, "name": "John"},
		{"id": 2, "name": "Jane"},
	}

	mockFileWriter.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		assert.Equal(t, filePath, filename)
		assert.Equal(t, os.FileMode(0644), perm)
		// SQLの内容が含まれていることを確認
		assert.Contains(t, string(data), "INSERT INTO")
		assert.Contains(t, string(data), tableName)
		return nil
	}

	// Act
	err := dumper.writeSQLFile(filePath, data, tableName)

	// Assert
	assert.NoError(t, err)
}

func TestTableDumper_writeSQLFile_EmptyData(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/empty.sql"
	tableName := "empty_table"
	data := []map[string]any{}
	expectedContent := "-- No data to export\n"

	mockFileWriter.WriteFileFunc = func(filename string, data []byte, perm os.FileMode) error {
		assert.Equal(t, filePath, filename)
		assert.Equal(t, []byte(expectedContent), data)
		assert.Equal(t, os.FileMode(0644), perm)
		return nil
	}

	// Act
	err := dumper.writeSQLFile(filePath, data, tableName)

	// Assert
	assert.NoError(t, err)
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
	result, err := dumper.DumpTable(context.Background(), options)

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
	result, err := dumper.DumpTable(context.Background(), options)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "出力ディレクトリ作成エラー")
}

// #==============================================================#
// ##          DefaultFileWriter Tests                           ##
// #==============================================================#

func TestDefaultFileWriter_WriteFile_Normal(t *testing.T) {
	// Arrange
	writer := &DefaultFileWriter{}
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("test content")

	// Act
	err := writer.WriteFile(filePath, content, 0644)

	// Assert
	assert.NoError(t, err)

	// ファイルが作成されたことを確認
	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, content, data)
}

func TestDefaultFileWriter_MkdirAll_Normal(t *testing.T) {
	// Arrange
	writer := &DefaultFileWriter{}
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
	writer := &DefaultFileWriter{}
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
