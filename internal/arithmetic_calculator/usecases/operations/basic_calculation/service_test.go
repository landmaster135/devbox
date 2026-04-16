package basiccalculation

import (
	"math"
	"testing"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		x        float64
		y        float64
		expected float64
	}{
		{name: "positive", x: 5, y: 3, expected: 8},
		{name: "negative", x: -2, y: -3, expected: -5},
		{name: "mixed", x: 5, y: -3, expected: 2},
		{name: "decimal", x: 2.5, y: 3.5, expected: 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Add(tt.x, tt.y))
		})
	}
}

func TestSubtract(t *testing.T) {
	assert.Equal(t, 6.0, Subtract(10, 4))
	assert.Equal(t, 1.0, Subtract(-2, -3))
	assert.InDelta(t, 3.3, Subtract(5.5, 2.2), 1e-10)
}

func TestMultiply(t *testing.T) {
	assert.Equal(t, 42.0, Multiply(6, 7))
	assert.Equal(t, -15.0, Multiply(5, -3))
	assert.Equal(t, 0.0, Multiply(5, 0))
}

func TestDivide(t *testing.T) {
	assert.Equal(t, 4.0, Divide(20, 5))
	assert.Equal(t, -5.0, Divide(15, -3))
	assert.True(t, math.IsInf(Divide(5, 0), 1))
}

func TestSum(t *testing.T) {
	assert.Equal(t, 15.0, Sum([]float64{1, 2, 3, 4, 5}))
	assert.Equal(t, 0.0, Sum([]float64{-5, -3, 0, 3, 5}))
	assert.Equal(t, 0.0, Sum([]float64{}))
}

func TestHandleToCalculate(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		x           float64
		y           float64
		expected    float64
		expectError bool
	}{
		{name: "add", operation: config.OperationAdd, x: 5, y: 3, expected: 8},
		{name: "subtract", operation: config.OperationSubtract, x: 10, y: 4, expected: 6},
		{name: "multiply", operation: config.OperationMultiply, x: 6, y: 7, expected: 42},
		{name: "divide", operation: config.OperationDivide, x: 20, y: 5, expected: 4},
		{name: "divide by zero", operation: config.OperationDivide, x: 5, y: 0, expectError: true},
		{name: "unknown", operation: "invalid", x: 5, y: 3, expected: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HandleToCalculate(tt.operation, tt.x, tt.y)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandleToCalculateWithArray(t *testing.T) {
	result, err := HandleToCalculateWithArray(config.OperationSum, []float64{1.5, 2.5, 3.5})
	assert.NoError(t, err)
	assert.Equal(t, 7.5, result)

	unknown, err := HandleToCalculateWithArray("invalid", []float64{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, 0.0, unknown)
}

func TestServiceExecute(t *testing.T) {
	service := NewService()

	result, err := service.Execute(config.OperationAdd, 10, 5)
	assert.NoError(t, err)
	assert.Equal(t, "10.00 + 5.00 = 15.00\n", result)

	_, err = service.Execute(config.OperationDivide, 10, 0)
	assert.Error(t, err)
}
