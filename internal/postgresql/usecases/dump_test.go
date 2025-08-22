package usecases

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockFileWriter はテスト用のFileWriterモック
type MockFileWriter struct {
	mock.Mock
}

func (m *MockFileWriter) WriteFile(filename string, data []byte, perm os.FileMode) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

func (m *MockFileWriter) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

func (m *MockFileWriter) Create(name string) (*os.File, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*os.File), args.Error(1)
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
	mockFileWriter.AssertNotCalled(t, "MkdirAll")
}

func TestTableDumper_ensureOutputDirectory_EmptyPath(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()

	// Act
	err := dumper.ensureOutputDirectory("")

	// Assert
	assert.NoError(t, err)
	mockFileWriter.AssertNotCalled(t, "MkdirAll")
}

func TestTableDumper_ensureOutputDirectory_CreateDirectory(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	outputPath := "/tmp/dump"

	mockFileWriter.On("MkdirAll", outputPath, os.FileMode(0755)).Return(nil)

	// Act
	err := dumper.ensureOutputDirectory(outputPath)

	// Assert
	assert.NoError(t, err)
	mockFileWriter.AssertExpectations(t)
}

func TestTableDumper_ensureOutputDirectory_CreateDirectoryError(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	outputPath := "/tmp/dump"
	expectedError := fmt.Errorf("permission denied")

	mockFileWriter.On("MkdirAll", outputPath, os.FileMode(0755)).Return(expectedError)

	// Act
	err := dumper.ensureOutputDirectory(outputPath)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockFileWriter.AssertExpectations(t)
}

// #==============================================================#
// ##          writeJSONFile Tests                               ##
// #==============================================================#

func TestTableDumper_writeJSONFile_Normal(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/test.json"
	data := []map[string]interface{}{
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

	mockFileWriter.On("WriteFile", filePath, []byte(expectedJSON), os.FileMode(0644)).Return(nil)

	// Act
	err := dumper.writeJSONFile(filePath, data)

	// Assert
	assert.NoError(t, err)
	mockFileWriter.AssertExpectations(t)
}

func TestTableDumper_writeJSONFile_WriteError(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/test.json"
	data := []map[string]interface{}{
		{"id": 1, "name": "John"},
	}
	expectedError := fmt.Errorf("write error")

	mockFileWriter.On("WriteFile", filePath, mock.AnythingOfType("[]uint8"), os.FileMode(0644)).Return(expectedError)

	// Act
	err := dumper.writeJSONFile(filePath, data)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockFileWriter.AssertExpectations(t)
}

// #==============================================================#
// ##          writeSQLFile Tests                                ##
// #==============================================================#

func TestTableDumper_writeSQLFile_Normal(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/users_20250822_170305.sql"
	tableName := "users"
	data := []map[string]interface{}{
		{"id": 1, "name": "John"},
		{"id": 2, "name": "Jane"},
	}

	mockFileWriter.On("WriteFile", filePath, mock.AnythingOfType("[]uint8"), os.FileMode(0644)).Return(nil)

	// Act
	err := dumper.writeSQLFile(filePath, data, tableName)

	// Assert
	assert.NoError(t, err)
	mockFileWriter.AssertExpectations(t)
}

func TestTableDumper_writeSQLFile_EmptyData(t *testing.T) {
	// Arrange
	dumper, _, mockFileWriter := createTestTableDumper()
	filePath := "/tmp/empty.sql"
	tableName := "empty_table"
	data := []map[string]interface{}{}
	expectedContent := "-- No data to export\n"

	mockFileWriter.On("WriteFile", filePath, []byte(expectedContent), os.FileMode(0644)).Return(nil)

	// Act
	err := dumper.writeSQLFile(filePath, data, tableName)

	// Assert
	assert.NoError(t, err)
	mockFileWriter.AssertExpectations(t)
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
	dumper, mockExecutor, mockFileWriter := createTestTableDumper()
	options := DumpOptions{
		TableName:  "users",
		OutputPath: "/invalid/path",
		Format:     "json",
	}
	expectedError := fmt.Errorf("permission denied")

	// データ取得は成功するがディレクトリ作成で失敗するシナリオ
	// ただし、ディレクトリ作成はデータ取得の前に実行されるため、データ取得のモックは不要
	mockFileWriter.On("MkdirAll", "/invalid/path", os.FileMode(0755)).Return(expectedError)

	// Act
	result, err := dumper.DumpTable(context.Background(), options)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "出力ディレクトリ作成エラー")
	mockFileWriter.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
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
