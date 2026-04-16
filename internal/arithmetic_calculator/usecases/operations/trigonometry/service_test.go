package trigonometry

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name        string
		function    string
		angle       float64
		unit        string
		expected    float64
		expectError bool
	}{
		{name: "sin degrees", function: "sin", angle: 30, unit: "degrees", expected: 0.5},
		{name: "cos radians", function: "cos", angle: math.Pi, unit: "radians", expected: -1},
		{name: "tan degrees", function: "tan", angle: 45, unit: "degrees", expected: 1},
		{name: "invalid", function: "bad", angle: 0, unit: "degrees", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Calculate(tt.function, tt.angle, tt.unit)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.InDelta(t, tt.expected, result, 1e-10)
		})
	}
}

func TestHandleToTrigonometry(t *testing.T) {
	result, err := HandleToTrigonometry("sin", math.Pi/2, "radians")
	assert.NoError(t, err)
	assert.InDelta(t, 1.0, result, 1e-10)
}

func TestServiceExecute(t *testing.T) {
	service := NewService()
	result, err := service.Execute("sin", 90, "degrees")
	assert.NoError(t, err)
	assert.Equal(t, "sin(90.00 degrees) = 1.000000\n", result)
}
