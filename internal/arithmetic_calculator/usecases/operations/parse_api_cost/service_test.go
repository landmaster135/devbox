package parseapicost

import (
	"errors"
	"os"
	"testing"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	"github.com/stretchr/testify/assert"
)

type mockFileReader struct {
	dataToReturn []byte
	errToReturn  error
}

func (m *mockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.errToReturn != nil {
		return nil, m.errToReturn
	}
	return m.dataToReturn, nil
}

func TestExtractAPICostFromText(t *testing.T) {
	service := NewService()

	tests := []struct {
		name     string
		text     string
		expected float64
	}{
		{name: "single", text: "API料金が100円掛かった", expected: 100},
		{name: "multiple", text: "API料金が50円掛かった API料金が75円掛かった", expected: 125},
		{name: "none", text: "料金情報なし", expected: 0},
		{name: "zero", text: "API料金が0円掛かった", expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ExtractAPICostFromText(tt.text)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleApiCostExtractionTextInput(t *testing.T) {
	service := NewService()

	result, err := service.HandleApiCostExtraction("", "API料金が100円掛かった。API料金が200円掛かった。")
	assert.NoError(t, err)
	assert.Equal(t, 300.0, result)
}

func TestHandleApiCostExtractionValidationErrors(t *testing.T) {
	service := NewService()

	_, err := service.HandleApiCostExtraction("a.md", "API料金が100円掛かった")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "同時に指定できません")

	_, err = service.HandleApiCostExtraction("", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "いずれかを指定")

	_, err = service.HandleApiCostExtraction("a.json", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ".mdまたは.txt")
}

func TestHandleApiCostExtractionWithFile(t *testing.T) {
	service := NewService()

	tmpFile, err := os.CreateTemp("", "api-cost-*.md")
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(tmpFile.Name())
	})

	_, err = tmpFile.WriteString("API料金が250円掛かった。API料金が150円掛かった。")
	assert.NoError(t, err)
	assert.NoError(t, tmpFile.Close())

	result, err := service.HandleApiCostExtraction(tmpFile.Name(), "")
	assert.NoError(t, err)
	assert.Equal(t, 400.0, result)
}

func TestHandleApiCostExtractionFileReadError(t *testing.T) {
	service := NewServiceWithFileReader(&mockFileReader{errToReturn: errors.New("read error")})

	_, err := service.HandleApiCostExtraction("a.md", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ファイル読み込みエラー")
}

func TestExecute(t *testing.T) {
	service := NewServiceWithFileReader(&mockFileReader{dataToReturn: []byte("API料金が300円掛かった")})

	result, err := service.Execute("a.md", "")
	assert.NoError(t, err)
	assert.Equal(t, "抽出されたAPI料金の合計: 300円\n", result)
}

func TestConstructors(t *testing.T) {
	assert.NotNil(t, NewService())
	assert.NotNil(t, NewServiceWithFileReader(&config.StandardFileReader{}))
}
