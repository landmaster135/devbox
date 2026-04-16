package arraysum

import (
	"testing"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	"github.com/stretchr/testify/assert"
)

func TestServiceExecuteSum(t *testing.T) {
	service := NewService()

	result, err := service.Execute(config.OperationSum, []float64{1, 2, 3, 4, 5})
	assert.NoError(t, err)
	assert.Equal(t, "sum([1 2 3 4 5]) = 15.00\n", result)
}

func TestServiceExecuteUnsupportedOperation(t *testing.T) {
	service := NewService()

	_, err := service.Execute("invalid", []float64{1, 2, 3})
	assert.Error(t, err)
}
