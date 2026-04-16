package evaluatelinecount

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFileOpener struct {
	errToReturn error
	file        *os.File
}

func (m *mockFileOpener) Open(name string) (*os.File, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	if m.file != nil {
		return m.file, nil
	}
	return os.Open(name)
}

type mockBufioScanner struct {
	boolToReturn bool
	errToReturn  error
}

func (m *mockBufioScanner) Scan() bool {
	return m.boolToReturn
}

func (m *mockBufioScanner) Err() error {
	return m.errToReturn
}

type mockJSONMarshaler struct {
	errToReturn error
}

func (m *mockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return []byte(`{"mocked":"json"}`), nil
}

func createTempFileWithLines(t *testing.T, lines []string) string {
	t.Helper()

	tmpFile, err := os.CreateTemp("", "test-line-count-*.txt")
	assert.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	for _, line := range lines {
		_, err := tmpFile.WriteString(line + "\n")
		assert.NoError(t, err)
	}

	err = tmpFile.Close()
	assert.NoError(t, err)

	return tmpFile.Name()
}

func TestCountLines(t *testing.T) {
	service := NewService()

	filePath := createTempFileWithLines(t, []string{"1行目", "2行目", "3行目"})
	count, err := service.CountLines(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)
}

func TestCountLinesFileOpenError(t *testing.T) {
	service := NewServiceWithDependencies(
		&mockFileOpener{errToReturn: errors.New("open error")},
		&mockBufioScanner{},
		&DefaultJSONMarshaler{},
	)

	_, err := service.CountLines("dummy")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ファイルを開けませんでした")
}

func TestCountLinesScannerError(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "mock-file-*.txt")
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})
	_ = tmpFile.Close()

	service := NewServiceWithDependencies(
		&mockFileOpener{file: tmpFile},
		&mockBufioScanner{boolToReturn: false, errToReturn: errors.New("scanner error")},
		&DefaultJSONMarshaler{},
	)

	_, err = service.CountLines("dummy")
	assert.Error(t, err)
	assert.Equal(t, "scanner error", err.Error())
}

func TestIsLineCountGreaterThan(t *testing.T) {
	service := NewService()
	filePath := createTempFileWithLines(t, []string{"1", "2", "3", "4"})

	isGreater, count, err := service.IsLineCountGreaterThan(filePath, 3)
	assert.NoError(t, err)
	assert.True(t, isGreater)
	assert.Equal(t, 4, count)
}

func TestHandleToEvaluateLineCount(t *testing.T) {
	service := NewService()
	filePath := createTempFileWithLines(t, []string{"1", "2", "3"})

	result, err := service.HandleToEvaluateLineCount(filePath, 2)
	assert.NoError(t, err)
	assert.Contains(t, result, `"is_greater": true`)
	assert.Contains(t, result, `"line_count": 3`)
	assert.Contains(t, result, `"threshold": 2`)
	assert.Contains(t, result, "より大きいです")
}

func TestHandleToEvaluateLineCountJSONError(t *testing.T) {
	filePath := createTempFileWithLines(t, []string{"1", "2"})
	service := NewServiceWithDependencies(
		&DefaultFileOpener{},
		&mockBufioScanner{},
		&mockJSONMarshaler{errToReturn: errors.New("json error")},
	)

	_, err := service.HandleToEvaluateLineCount(filePath, 1)
	assert.Error(t, err)
	assert.Equal(t, "json error", err.Error())
}

func TestExecute(t *testing.T) {
	service := NewService()
	filePath := createTempFileWithLines(t, []string{"1", "2"})

	result, err := service.Execute(filePath, 1)
	assert.NoError(t, err)
	assert.Contains(t, result, `"line_count": 2`)
}

func TestIsGreaterDescription(t *testing.T) {
	assert.Equal(t, "より大きいです。", IsGreaterDescription(true))
	assert.Equal(t, "以下です。", IsGreaterDescription(false))
}
