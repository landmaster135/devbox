package squareroot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	result, err := Calculate(16)
	assert.NoError(t, err)
	assert.Equal(t, 4.0, result)

	_, err = Calculate(-1)
	assert.Error(t, err)
}

func TestHandleToSquareRoot(t *testing.T) {
	result, err := HandleToSquareRoot(25)
	assert.NoError(t, err)
	assert.Equal(t, 5.0, result)
}

func TestServiceExecute(t *testing.T) {
	service := NewService()
	result, err := service.Execute(64)
	assert.NoError(t, err)
	assert.Equal(t, "√64.00 = 8.00\n", result)
}
