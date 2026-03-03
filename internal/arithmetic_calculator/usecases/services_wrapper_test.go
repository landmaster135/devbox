package usecases

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	"github.com/stretchr/testify/assert"
)

type wrapperMockFileReader struct {
	data []byte
	err  error
}

func (m *wrapperMockFileReader) ReadFile(filename string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data, nil
}

type wrapperMockFileOpener struct {
	err error
}

func (m *wrapperMockFileOpener) Open(name string) (*os.File, error) {
	if m.err != nil {
		return nil, m.err
	}
	return os.Open(name)
}

type wrapperMockScanner struct {
	scanResult bool
	err        error
}

func (m *wrapperMockScanner) Scan() bool { return m.scanResult }
func (m *wrapperMockScanner) Err() error { return m.err }

type wrapperMockJSONMarshaler struct {
	err error
}

func (m *wrapperMockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []byte("{}"), nil
}

func createTempTextFile(t *testing.T, content string) string {
	t.Helper()

	f, err := os.CreateTemp("", "arithmetic-wrapper-*.txt")
	assert.NoError(t, err)
	t.Cleanup(func() {
		_ = os.Remove(f.Name())
	})

	_, err = f.WriteString(content)
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
	return f.Name()
}

func TestCalculatorService_Wrappers(t *testing.T) {
	s := NewCalculatorService()
	assert.NotNil(t, s)

	assert.Equal(t, 8.0, s.add(5, 3))
	assert.Equal(t, 6.0, s.subtract(10, 4))
	assert.Equal(t, 42.0, s.multiply(6, 7))
	assert.Equal(t, 4.0, s.divide(20, 5))
	assert.Equal(t, 15.0, s.sum([]float64{1, 2, 3, 4, 5}))

	result, err := s.HandleToCalculate(config.OperationAdd, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, result)

	_, err = s.HandleToCalculate(config.OperationDivide, 1, 0)
	assert.Error(t, err)

	arrResult, err := s.HandleToCalculateWithArray(config.OperationSum, []float64{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 6.0, arrResult)
}

func TestFileEvaluatorService_Wrappers(t *testing.T) {
	filePath := createTempTextFile(t, "a\nb\nc\n")

	s := NewFileEvaluatorService()
	assert.NotNil(t, s)

	count, err := s.CountLines(filePath)
	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	isGreater, lineCount, err := s.IsLineCountGreaterThan(filePath, 2)
	assert.NoError(t, err)
	assert.True(t, isGreater)
	assert.Equal(t, 3, lineCount)

	jsonResult, err := s.HandleToEvaluateLineCount(filePath, 2)
	assert.NoError(t, err)
	assert.Contains(t, jsonResult, `"line_count": 3`)
	assert.Equal(t, "より大きいです。", isGreaterDescription(true))
	assert.Equal(t, "以下です。", isGreaterDescription(false))

	withDeps := NewFileEvaluatorServiceWithDependencies(
		&wrapperMockFileOpener{err: errors.New("open error")},
		&wrapperMockScanner{},
		&wrapperMockJSONMarshaler{},
	)
	assert.NotNil(t, withDeps)
	_, err = withDeps.CountLines("/tmp/not-found")
	assert.Error(t, err)
}

func TestApiCostExtractorService_Wrappers(t *testing.T) {
	s := NewApiCostExtractorService()
	assert.NotNil(t, s)

	sum, err := s.extractApiCostFromText("API料金が100円掛かった API料金が200円掛かった")
	assert.NoError(t, err)
	assert.Equal(t, 300.0, sum)

	handleResult, err := s.HandleApiCostExtraction("", "API料金が50円掛かった")
	assert.NoError(t, err)
	assert.Equal(t, 50.0, handleResult)

	withReader := NewApiCostExtractorServiceWithFileReader(&wrapperMockFileReader{data: []byte("API料金が300円掛かった")})
	assert.NotNil(t, withReader)
	result, err := withReader.HandleApiCostExtraction("a.md", "")
	assert.NoError(t, err)
	assert.Equal(t, 300.0, result)
}

func TestAdvancedMathService_Wrappers(t *testing.T) {
	s := NewAdvancedMathService()
	assert.NotNil(t, s)

	assert.Equal(t, 8.0, s.power(2, 3))
	sqrt, err := s.squareRoot(16)
	assert.NoError(t, err)
	assert.Equal(t, 4.0, sqrt)

	fact, err := s.factorial(5)
	assert.NoError(t, err)
	assert.Equal(t, 120.0, fact)

	_, err = s.squareRoot(-1)
	assert.Error(t, err)
	_, err = s.factorial(-1)
	assert.Error(t, err)

	powerResult, err := s.HandleToPower(3, 2)
	assert.NoError(t, err)
	assert.Equal(t, 9.0, powerResult)

	_, err = s.HandleToSquareRoot(-9)
	assert.Error(t, err)
	_, err = s.HandleToFactorial(171)
	assert.Error(t, err)
}

func TestTrigonometryService_Wrappers(t *testing.T) {
	s := NewTrigonometryService()
	assert.NotNil(t, s)

	result, err := s.trigonometry("sin", 90, "degrees")
	assert.NoError(t, err)
	assert.InDelta(t, 1.0, result, 1e-10)

	_, err = s.trigonometry("bad", 0, "degrees")
	assert.Error(t, err)

	handleResult, err := s.HandleToTrigonometry("cos", 0, "degrees")
	assert.NoError(t, err)
	assert.InDelta(t, 1.0, handleResult, 1e-10)
}

func TestMathConstantsService_Wrappers(t *testing.T) {
	s := NewMathConstantsService()
	assert.NotNil(t, s)

	constants := s.getConstants()
	assert.Contains(t, constants, "pi")
	assert.Contains(t, constants, "e")
	assert.Contains(t, constants, "tau")

	output, err := s.HandleToGetConstants()
	assert.NoError(t, err)
	assert.Contains(t, output, "利用可能な数学定数")
}

func TestExpressionEvaluatorService_Wrappers(t *testing.T) {
	s := NewExpressionEvaluatorService()
	assert.NotNil(t, s)

	value, err := s.safeEvaluate("2+3")
	assert.NoError(t, err)
	assert.Equal(t, 5.0, value)

	value, err = s.evaluateBasicExpression("sqrt(16)")
	assert.NoError(t, err)
	assert.Equal(t, 4.0, value)

	value, err = s.evaluateArithmeticExpression("3**2")
	assert.NoError(t, err)
	assert.Equal(t, 9.0, value)

	assert.NoError(t, s.checkOsPattern("cos(0)"))
	assert.Error(t, s.checkOsPattern("os.system('ls')"))

	indices := s.getAllIndices(strings.ToLower("cos(cos(0))"), "os")
	assert.NotEmpty(t, indices)

	handleValue, err := s.HandleToCalculateExpression("2+3*4")
	assert.NoError(t, err)
	assert.Equal(t, 14.0, handleValue)
}
