package getconstants

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConstants(t *testing.T) {
	constants := GetConstants()
	assert.Equal(t, math.Pi, constants["pi"])
	assert.Equal(t, math.E, constants["e"])
	assert.Equal(t, 2*math.Pi, constants["tau"])
}

func TestHandleToGetConstants(t *testing.T) {
	result, err := HandleToGetConstants()
	assert.NoError(t, err)
	assert.Contains(t, result, "利用可能な数学定数:")
	assert.Contains(t, result, "pi =")
	assert.Contains(t, result, "e =")
	assert.Contains(t, result, "tau =")
	assert.True(t, strings.HasSuffix(result, "\n"))
}

func TestServiceExecute(t *testing.T) {
	service := NewService()
	result, err := service.Execute()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
