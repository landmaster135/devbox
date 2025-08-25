package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataConverterService_ConvertToHTML_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name                 string
		input                [][]string
		isTheadContained     bool
		textReplacingIfBlank string
		expected             string
	}{
		{
			name: "NormalCase_WithHeader",
			input: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
				{"Charlie", "", "Tokyo"},
			},
			isTheadContained: true,
			textReplacingIfBlank: "💩",
			expected: "<table>\n<thead>\n<tr><th>Name</th><th>Age</th><th>City</th></tr>\n</thead>\n<tbody>\n<tr><td>Alice</td><td>25</td><td>New York</td></tr>\n<tr><td>Bob</td><td>30</td><td>London</td></tr>\n<tr><td>Charlie</td><td>💩</td><td>Tokyo</td></tr>\n</tbody>\n</table>",
		},
		{
			name: "NormalCase_WithoutHeader",
			input: [][]string{
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
			isTheadContained: false,
			textReplacingIfBlank: "💩",
			expected: "<table>\n<thead>\n<tr><th>Alice</th><th>25</th><th>New York</th></tr>\n</thead>\n<tbody>\n<tr><td>Bob</td><td>30</td><td>London</td></tr>\n</tbody>\n</table>",
		},
		{
			name:                 "EmptyCase",
			input:                [][]string{},
			isTheadContained:     true,
			textReplacingIfBlank: "💩",
			expected:             "<table></table>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToHTML(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertData_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name          string
		inputFormat   string
		outputFormat  string
		input         string
		inputFilePath string
		expected      string
	}{
		{
			name:          "JSONToHTML_Normal",
			inputFormat:   "json",
			outputFormat:  "html",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expected:      "<table>\n<thead>\n<tr><th>Name</th><th>Age</th><th>City</th></tr>\n</thead>\n<tbody>\n<tr><td>Alice</td><td>25</td><td>New York</td></tr>\n<tr><td>Bob</td><td>30</td><td>London</td></tr>\n</tbody>\n</table>",
		},
		{
			name:          "JSONToCSV_Normal",
			inputFormat:   "json",
			outputFormat:  "csv",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expected:      "Name,Age,City\nAlice,25,New York\nBob,30,London",
		},
		{
			name:          "JSONToTSV_Normal",
			inputFormat:   "json",
			outputFormat:  "tsv",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expected:      "Name\tAge\tCity\nAlice\t25\tNew York\nBob\t30\tLondon",
		},
		{
			name:          "JSONToArray_Normal",
			inputFormat:   "json",
			outputFormat:  "array",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expected:      "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
		},
		{
			name:          "JSONToJSON_Normal",
			inputFormat:   "json",
			outputFormat:  "json",
			input:         "[[\"Name\",\"Age\",\"Score\"],[\"Alice\",\"25\",\"85.5\"],[\"Bob\",\"30\",\"92\"],[\"Charlie\",\"\",\"78\"]]",
			inputFilePath: "",
			expected:      "[{\"Age\":25,\"Name\":\"Alice\",\"Score\":85.5},{\"Age\":30,\"Name\":\"Bob\",\"Score\":92},{\"Age\":null,\"Name\":\"Charlie\",\"Score\":78}]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ConvertData(tt.inputFormat, tt.outputFormat, tt.input, tt.inputFilePath)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertData_Error(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name          string
		inputFormat   string
		outputFormat  string
		input         string
		inputFilePath string
		expectedError string
	}{
		{
			name:          "UnsupportedInputFormat_Error",
			inputFormat:   "xml",
			outputFormat:  "html",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expectedError: "入力データの解析に失敗しました: 未対応の入力形式です: xml",
		},
		{
			name:          "UnsupportedOutputFormat_Error",
			inputFormat:   "json",
			outputFormat:  "xml",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expectedError: "未対応の出力形式です: xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ConvertData(tt.inputFormat, tt.outputFormat, tt.input, tt.inputFilePath)
			assert.Error(t, err)
			assert.Equal(t, tt.expectedError, err.Error())
			assert.Equal(t, "", result)
		})
	}
}

func TestDataConverterService_ConvertToArray_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "NormalCase",
			input: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
			expected: "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToArray(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertToJSON_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "NormalCase",
			input: [][]string{
				{"Name", "Age", "Score"},
				{"Alice", "25", "85.5"},
				{"Bob", "30", "92"},
				{"Charlie", "", "78"},
			},
			expected: "[{\"Age\":25,\"Name\":\"Alice\",\"Score\":85.5},{\"Age\":30,\"Name\":\"Bob\",\"Score\":92},{\"Age\":null,\"Name\":\"Charlie\",\"Score\":78}]",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToJSON(tt.input)
			assert.NoError(t, err)

			// JSONの比較は文字列比較ではなく、構造体の比較の方が適切
			// ここでは文字列比較で実装するが、実際には構造体の比較の方が望ましい
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertToCSV_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "NormalCase",
			input: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
			expected: "Name,Age,City\nAlice,25,New York\nBob,30,London",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToCSV(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertToTSV_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "NormalCase",
			input: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
			expected: "Name\tAge\tCity\nAlice\t25\tNew York\nBob\t30\tLondon",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToTSV(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
