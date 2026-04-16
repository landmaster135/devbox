package usecases

import (
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
)

type mockBasicCalculationOperation struct {
	called bool
}

func (m *mockBasicCalculationOperation) Execute(operation string, x, y float64) (string, error) {
	m.called = true
	return "basic\n", nil
}

type mockArraySumOperation struct {
	called bool
}

func (m *mockArraySumOperation) Execute(operation string, numbers []float64) (string, error) {
	m.called = true
	return "sum\n", nil
}

type mockLineCountOperation struct {
	called bool
}

func (m *mockLineCountOperation) Execute(filePath string, threshold int) (string, error) {
	m.called = true
	return "linecount", nil
}

type mockAPICostOperation struct {
	called bool
}

func (m *mockAPICostOperation) Execute(filePath, textInput string) (string, error) {
	m.called = true
	return "api-cost\n", nil
}

type mockPowerOperation struct {
	called bool
}

func (m *mockPowerOperation) Execute(base, exponent float64) (string, error) {
	m.called = true
	return "power\n", nil
}

type mockSquareRootOperation struct {
	called bool
}

func (m *mockSquareRootOperation) Execute(number float64) (string, error) {
	m.called = true
	return "sqrt\n", nil
}

type mockFactorialOperation struct {
	called bool
}

func (m *mockFactorialOperation) Execute(n int) (string, error) {
	m.called = true
	return "factorial\n", nil
}

type mockTrigonometryOperation struct {
	called bool
}

func (m *mockTrigonometryOperation) Execute(function string, angle float64, unit string) (string, error) {
	m.called = true
	return "trigonometry\n", nil
}

type mockCalculateOperation struct {
	called bool
}

func (m *mockCalculateOperation) Execute(expression string) (string, error) {
	m.called = true
	return "calculate\n", nil
}

type mockConstantsOperation struct {
	called bool
}

func (m *mockConstantsOperation) Execute() (string, error) {
	m.called = true
	return "constants\n", nil
}

func TestService_ExecuteByConfig_Normal(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		expected string
		assert   func(
			t *testing.T,
			basic *mockBasicCalculationOperation,
			sum *mockArraySumOperation,
			lineCount *mockLineCountOperation,
			apiCost *mockAPICostOperation,
			power *mockPowerOperation,
			squareRoot *mockSquareRootOperation,
			factorial *mockFactorialOperation,
			trigonometry *mockTrigonometryOperation,
			calculate *mockCalculateOperation,
			constants *mockConstantsOperation,
		)
	}{
		{
			name:     "basic calculation",
			cfg:      &config.Config{Operation: config.OperationAdd, X: 1, Y: 2},
			expected: "basic\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if !basic.called || sum.called || lineCount.called || apiCost.called || power.called || squareRoot.called || factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "sum",
			cfg:      &config.Config{Operation: config.OperationSum, Numbers: []float64{1, 2}},
			expected: "sum\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || !sum.called || lineCount.called || apiCost.called || power.called || squareRoot.called || factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "line count",
			cfg:      &config.Config{Operation: config.OperationEvaluateLineCount, FilePath: "x", Threshold: 1},
			expected: "linecount",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || !lineCount.called || apiCost.called || power.called || squareRoot.called || factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "api cost",
			cfg:      &config.Config{Operation: config.OperationParseAPICost, TextInput: "API料金が100円掛かった"},
			expected: "api-cost\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || !apiCost.called || power.called || squareRoot.called || factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "power",
			cfg:      &config.Config{Operation: config.OperationPower, Base: 2, Exponent: 8},
			expected: "power\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || apiCost.called || !power.called || squareRoot.called || factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "square root",
			cfg:      &config.Config{Operation: config.OperationSquareRoot, Number: 16},
			expected: "sqrt\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || apiCost.called || power.called || !squareRoot.called || factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "factorial",
			cfg:      &config.Config{Operation: config.OperationFactorial, N: 5},
			expected: "factorial\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || apiCost.called || power.called || squareRoot.called || !factorial.called || trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "trigonometry",
			cfg:      &config.Config{Operation: config.OperationTrigonometry, Function: "sin", Angle: 90, Unit: "degrees"},
			expected: "trigonometry\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || apiCost.called || power.called || squareRoot.called || factorial.called || !trigonometry.called || calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "calculate",
			cfg:      &config.Config{Operation: config.OperationCalculate, Expression: "2+3"},
			expected: "calculate\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || apiCost.called || power.called || squareRoot.called || factorial.called || trigonometry.called || !calculate.called || constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
		{
			name:     "get constants",
			cfg:      &config.Config{Operation: config.OperationGetConstants},
			expected: "constants\n",
			assert: func(t *testing.T, basic *mockBasicCalculationOperation, sum *mockArraySumOperation, lineCount *mockLineCountOperation, apiCost *mockAPICostOperation, power *mockPowerOperation, squareRoot *mockSquareRootOperation, factorial *mockFactorialOperation, trigonometry *mockTrigonometryOperation, calculate *mockCalculateOperation, constants *mockConstantsOperation) {
				t.Helper()
				if basic.called || sum.called || lineCount.called || apiCost.called || power.called || squareRoot.called || factorial.called || trigonometry.called || calculate.called || !constants.called {
					t.Fatalf("unexpected dispatch state")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_Normal", func(t *testing.T) {
			basic := &mockBasicCalculationOperation{}
			sum := &mockArraySumOperation{}
			lineCount := &mockLineCountOperation{}
			apiCost := &mockAPICostOperation{}
			power := &mockPowerOperation{}
			squareRoot := &mockSquareRootOperation{}
			factorial := &mockFactorialOperation{}
			trigonometry := &mockTrigonometryOperation{}
			calculate := &mockCalculateOperation{}
			constants := &mockConstantsOperation{}

			service := newServiceWithOperations(
				basic,
				sum,
				lineCount,
				apiCost,
				power,
				squareRoot,
				factorial,
				trigonometry,
				calculate,
				constants,
			)

			result, err := service.ExecuteByConfig(tt.cfg)
			if err != nil {
				t.Fatalf("ExecuteByConfig returned error: %v", err)
			}
			if result != tt.expected {
				t.Fatalf("unexpected result: got=%q want=%q", result, tt.expected)
			}

			tt.assert(t, basic, sum, lineCount, apiCost, power, squareRoot, factorial, trigonometry, calculate, constants)
		})
	}
}

func TestService_ExecuteByConfig_NilConfig(t *testing.T) {
	service := newServiceWithOperations(
		&mockBasicCalculationOperation{},
		&mockArraySumOperation{},
		&mockLineCountOperation{},
		&mockAPICostOperation{},
		&mockPowerOperation{},
		&mockSquareRootOperation{},
		&mockFactorialOperation{},
		&mockTrigonometryOperation{},
		&mockCalculateOperation{},
		&mockConstantsOperation{},
	)

	_, err := service.ExecuteByConfig(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
	if !strings.Contains(err.Error(), "configがnilです") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestService_ExecuteByConfig_UnsupportedOperation(t *testing.T) {
	service := NewService()

	_, err := service.ExecuteByConfig(&config.Config{Operation: "unknown"})
	if err == nil {
		t.Fatal("expected unsupported operation error")
	}
	if !strings.Contains(err.Error(), "未サポートのoperationです") {
		t.Fatalf("unexpected error: %v", err)
	}
}
