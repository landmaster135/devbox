package calculateexpression

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeEvaluate(t *testing.T) {
	service := NewService()

	tests := []struct {
		name        string
		expression  string
		expected    float64
		expectError bool
		errorSubstr string
	}{
		{name: "arithmetic", expression: "2+3*4", expected: 14},
		{name: "constant", expression: "pi", expected: 3.141593},
		{name: "sin", expression: "sin(pi/2)", expected: 1},
		{name: "dangerous", expression: "import os", expectError: true, errorSubstr: "危険なパターン"},
		{name: "os dangerous", expression: "os.system('ls')", expectError: true, errorSubstr: "危険なパターン"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.SafeEvaluate(tt.expression)
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorSubstr)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, result, 1e-6)
		})
	}
}

func TestEvaluateMethods(t *testing.T) {
	service := NewService()

	result, err := service.EvaluateBasicExpression("sqrt(16)")
	assert.NoError(t, err)
	assert.Equal(t, 4.0, result)

	result, err = service.EvaluateArithmeticExpression("3**2")
	assert.NoError(t, err)
	assert.Equal(t, 9.0, result)
}

func TestCheckOSPatternAndIndices(t *testing.T) {
	service := NewService()

	assert.NoError(t, service.CheckOSPattern("cos(0)"))
	assert.Error(t, service.CheckOSPattern("os.system('ls')"))

	indices := service.GetAllIndices(strings.ToLower("cos(cos(0))"), "os")
	assert.NotEmpty(t, indices)
}

func TestExecute(t *testing.T) {
	service := NewService()

	result, err := service.Execute("2+3*4")
	assert.NoError(t, err)
	assert.Equal(t, "2+3*4 = 14.00\n", result)
}
