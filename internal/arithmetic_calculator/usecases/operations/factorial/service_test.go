package factorial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	result, err := Calculate(5)
	assert.NoError(t, err)
	assert.Equal(t, 120.0, result)

	_, err = Calculate(-1)
	assert.Error(t, err)

	_, err = Calculate(171)
	assert.Error(t, err)
}

func TestHandleToFactorial(t *testing.T) {
	result, err := HandleToFactorial(6)
	assert.NoError(t, err)
	assert.Equal(t, 720.0, result)
}

func TestServiceExecute(t *testing.T) {
	service := NewService()
	result, err := service.Execute(5)
	assert.NoError(t, err)
	assert.Equal(t, "5! = 120\n", result)
}
