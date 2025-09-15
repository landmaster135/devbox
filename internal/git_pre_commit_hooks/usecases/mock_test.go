package usecases

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockCommandExecutor_Execute_Normal(t *testing.T) {
	// Arrange
	mockExecutor := &MockCommandExecutor{}
	expectedOutput := []byte("test output")
	expectedError := errors.New("test error")
	commandName := "git"
	args := []string{"status", "--porcelain"}

	// Act & Assert - 正常系のテスト
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		assert.Equal(t, commandName, name)
		assert.Equal(t, []string{"status", "--porcelain"}, args)
		return expectedOutput, nil
	}
	output, err := mockExecutor.Execute(commandName, args...)

	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	// Act & Assert - エラー系のテスト
	mockExecutor.ExecuteFunc = func(name string, args ...string) ([]byte, error) {
		return nil, expectedError
	}
	output, err = mockExecutor.Execute(commandName, args...)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, output)
}

func TestMockOutputWriter_Printf_Normal(t *testing.T) {
	// Arrange
	mockWriter := &MockOutputWriter{
		Output: make([]string, 0),
	}
	format := "Hello %s, you have %d messages"
	args := []interface{}{"Alice", 5}
	expectedOutput := "Hello Alice, you have 5 messages"

	// Act & Assert - PrintfFunc設定時のテスト
	mockWriter.PrintfFunc = func(format string, args ...interface{}) {
		assert.Equal(t, "Hello %s, you have %d messages", format)
		assert.Equal(t, []interface{}{"Alice", 5}, args)
	}
	mockWriter.Printf(format, args...)

	// Output配列への記録確認
	assert.Len(t, mockWriter.Output, 1)
	assert.Equal(t, expectedOutput, mockWriter.Output[0])
}

func TestMockOutputWriter_Println_Normal(t *testing.T) {
	// Arrange
	mockWriter := &MockOutputWriter{
		Output: make([]string, 0),
	}
	args := []interface{}{"Debug message", "with", "multiple", "parts"}
	expectedOutput := "Debug message with multiple parts\n"

	// Act & Assert - PrintlnFunc設定時のテスト
	mockWriter.PrintlnFunc = func(args ...interface{}) {
		assert.Equal(t, []interface{}{"Debug message", "with", "multiple", "parts"}, args)
	}
	mockWriter.Println(args...)

	// Output配列への記録確認
	assert.Len(t, mockWriter.Output, 1)
	assert.Equal(t, expectedOutput, mockWriter.Output[0])
}

func TestMockOutputWriter_Printf_WithoutFunc_Normal(t *testing.T) {
	// Arrange
	mockWriter := &MockOutputWriter{
		Output: make([]string, 0),
	}
	format := "Simple %s"
	args := []interface{}{"test"}
	expectedOutput := "Simple test"

	// Act - PrintfFuncが設定されていない場合
	mockWriter.Printf(format, args...)

	// Assert - Output配列への記録のみ確認
	assert.Len(t, mockWriter.Output, 1)
	assert.Equal(t, expectedOutput, mockWriter.Output[0])
}

func TestMockOutputWriter_Println_WithoutFunc_Normal(t *testing.T) {
	// Arrange
	mockWriter := &MockOutputWriter{
		Output: make([]string, 0),
	}
	args := []interface{}{"Simple", "test"}
	expectedOutput := "Simple test\n"

	// Act - PrintlnFuncが設定されていない場合
	mockWriter.Println(args...)

	// Assert - Output配列への記録のみ確認
	assert.Len(t, mockWriter.Output, 1)
	assert.Equal(t, expectedOutput, mockWriter.Output[0])
}
