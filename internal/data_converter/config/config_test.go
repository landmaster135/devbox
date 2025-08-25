package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Validate_Normal(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "ValidConfig_WithInput",
			config: &Config{
				InputFormat:   "json",
				OutputFormat:  "html",
				Input:         "test data",
				InputFilePath: "",
				Help:          false,
			},
		},
		{
			name: "ValidConfig_WithFile",
			config: &Config{
				InputFormat:   "csv",
				OutputFormat:  "json",
				Input:         "",
				InputFilePath: "test.csv",
				Help:          false,
			},
		},
		{
			name: "ValidConfig_Help",
			config: &Config{
				InputFormat:   "",
				OutputFormat:  "",
				Input:         "",
				InputFilePath: "",
				Help:          true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			assert.NoError(t, err)
		})
	}
}

func TestConfig_Validate_Error(t *testing.T) {
	tests := []struct {
		name          string
		config        *Config
		expectedError string
	}{
		{
			name: "MissingInputFormat_Error",
			config: &Config{
				InputFormat:   "",
				OutputFormat:  "html",
				Input:         "data",
				InputFilePath: "",
				Help:          false,
			},
			expectedError: "input-formatは必須です",
		},
		{
			name: "MissingOutputFormat_Error",
			config: &Config{
				InputFormat:   "json",
				OutputFormat:  "",
				Input:         "data",
				InputFilePath: "",
				Help:          false,
			},
			expectedError: "output-formatは必須です",
		},
		{
			name: "BothInputAndFile_Error",
			config: &Config{
				InputFormat:   "json",
				OutputFormat:  "html",
				Input:         "data",
				InputFilePath: "file.json",
				Help:          false,
			},
			expectedError: "inputとinput-file-pathは同時に指定できません",
		},
		{
			name: "NoInputOrFile_Error",
			config: &Config{
				InputFormat:   "json",
				OutputFormat:  "html",
				Input:         "",
				InputFilePath: "",
				Help:          false,
			},
			expectedError: "inputまたはinput-file-pathのいずれかを指定してください",
		},
		{
			name: "InvalidInputFormat_Error",
			config: &Config{
				InputFormat:   "xml",
				OutputFormat:  "html",
				Input:         "data",
				InputFilePath: "",
				Help:          false,
			},
			expectedError: "未対応の入力形式です: xml",
		},
		{
			name: "InvalidOutputFormat_Error",
			config: &Config{
				InputFormat:   "json",
				OutputFormat:  "xml",
				Input:         "data",
				InputFilePath: "",
				Help:          false,
			},
			expectedError: "未対応の出力形式です: xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestIsValidInputFormat_Normal(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "ValidJSON",
			format:   "json",
			expected: true,
		},
		{
			name:     "ValidCSV",
			format:   "csv",
			expected: true,
		},
		{
			name:     "ValidTSV",
			format:   "tsv",
			expected: true,
		},
		{
			name:     "ValidHTML",
			format:   "html",
			expected: true,
		},
		{
			name:     "InvalidXML",
			format:   "xml",
			expected: false,
		},
		{
			name:     "EmptyString",
			format:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidInputFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidOutputFormat_Normal(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "ValidHTML",
			format:   "html",
			expected: true,
		},
		{
			name:     "ValidCSV",
			format:   "csv",
			expected: true,
		},
		{
			name:     "ValidTSV",
			format:   "tsv",
			expected: true,
		},
		{
			name:     "ValidJSON",
			format:   "json",
			expected: true,
		},
		{
			name:     "InvalidXML",
			format:   "xml",
			expected: false,
		},
		{
			name:     "EmptyString",
			format:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidOutputFormat(tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPrintUsage(t *testing.T) {
	// PrintUsageは標準出力に出力するため、パニックしないことを確認
	assert.NotPanics(t, func() {
		PrintUsage()
	})
}
