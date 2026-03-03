package power

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	assert.Equal(t, 8.0, Calculate(2, 3))
	assert.Equal(t, -8.0, Calculate(-2, 3))
	assert.InDelta(t, 0.25, Calculate(2, -2), 1e-10)
}

func TestHandleToPower(t *testing.T) {
	result, err := HandleToPower(3, 4)
	assert.NoError(t, err)
	assert.Equal(t, 81.0, result)
}

func TestServiceExecute(t *testing.T) {
	service := NewService()
	result, err := service.Execute(2, 8)
	assert.NoError(t, err)
	assert.Equal(t, "2.00^8.00 = 256.00\n", result)
}
