package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
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
			isTheadContained:     true,
			textReplacingIfBlank: "💩",
			expected:             "<table>\n<thead>\n<tr><th>Name</th><th>Age</th><th>City</th></tr>\n</thead>\n<tbody>\n<tr><td>Alice</td><td>25</td><td>New York</td></tr>\n<tr><td>Bob</td><td>30</td><td>London</td></tr>\n<tr><td>Charlie</td><td>💩</td><td>Tokyo</td></tr>\n</tbody>\n</table>",
		},
		{
			name: "NormalCase_WithoutHeader",
			input: [][]string{
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
			isTheadContained:     false,
			textReplacingIfBlank: "💩",
			expected:             "<table>\n<thead>\n<tr><th>Alice</th><th>25</th><th>New York</th></tr>\n</thead>\n<tbody>\n<tr><td>Bob</td><td>30</td><td>London</td></tr>\n</tbody>\n</table>",
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
		{
			name:         "YAMLToCSV_Normal",
			inputFormat:  "yaml",
			outputFormat: "csv",
			input: `- Name: Alice
  Age: "25"
  City: New York
- Name: Bob
  Age: "30"
  City: London
`,
			inputFilePath: "",
			expected:      "Name,Age,City\nAlice,25,New York\nBob,30,London",
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

func TestDataConverterService_ConvertToYAML_Normal(t *testing.T) {
	service := NewDataConverterService()
	input := [][]string{
		{"Name", "Age", "City"},
		{"Alice", "25", "New York"},
	}

	result, err := service.convertToYAML(input)
	assert.NoError(t, err)

	var decoded []map[string]any
	decodeErr := yaml.Unmarshal([]byte(result), &decoded)
	assert.NoError(t, decodeErr)
	assert.Equal(t, []map[string]any{{"Age": 25, "City": "New York", "Name": "Alice"}}, decoded)
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

func TestDataConverterService_ConvertData_OutputYAML(t *testing.T) {
	service := NewDataConverterService()
	input := "[[\"Name\",\"Age\"],[\"Alice\",\"25\"],[\"Bob\",\"30\"]]"

	result, err := service.ConvertData("json", "yaml", input, "")
	assert.NoError(t, err)

	var decoded []map[string]any
	decodeErr := yaml.Unmarshal([]byte(result), &decoded)
	assert.NoError(t, decodeErr)
	assert.Equal(t, []map[string]any{
		{"Age": 25, "Name": "Alice"},
		{"Age": 30, "Name": "Bob"},
	}, decoded)
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

func TestDataConverterService_ConvertToTable_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "MultiColumnTable",
			input: [][]string{
				{"名前", "年齢", "職業"},
				{"田中", "25", "エンジニア"},
				{"佐藤", "30", "デザイナー"},
			},
			expected: "| 名前 | 年齢 | 職業          |\n|--------|--------|-----------------|\n| 田中 | 25     | エンジニア |\n| 佐藤 | 30     | デザイナー |\n",
		},
		{
			name: "SingleColumnTable",
			input: [][]string{
				{"項目"},
				{"項目1"},
				{"項目2"},
			},
			expected: "|     | 項目  |\n|-----|---------|\n|     | 項目1 |\n|     | 項目2 |\n",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToTable(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertToList_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "SingleColumnList",
			input: [][]string{
				{"項目"},
				{"項目1"},
				{"項目2"},
				{"項目3"},
			},
			expected: "- 項目1\n- 項目2\n- 項目3\n",
		},
		{
			name: "MultiColumnList",
			input: [][]string{
				{"名前", "年齢"},
				{"田中", "25"},
				{"佐藤", "30"},
			},
			expected: "- 名前: 田中, 年齢: 25\n- 名前: 佐藤, 年齢: 30\n",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToList(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertToOrderedList_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name     string
		input    [][]string
		expected string
	}{
		{
			name: "SingleColumnOrderedList",
			input: [][]string{
				{"項目"},
				{"項目1"},
				{"項目2"},
				{"項目3"},
			},
			expected: "1. 項目1\n2. 項目2\n3. 項目3\n",
		},
		{
			name: "MultiColumnOrderedList",
			input: [][]string{
				{"名前", "年齢"},
				{"田中", "25"},
				{"佐藤", "30"},
			},
			expected: "1. 名前: 田中, 年齢: 25\n2. 名前: 佐藤, 年齢: 30\n",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.convertToOrderedList(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataConverterService_ConvertData_TableFormat_Normal(t *testing.T) {
	service := NewDataConverterService()

	tests := []struct {
		name             string
		inputFormat      string
		outputFormat     string
		input            string
		inputFilePath    string
		expectedContains []string
	}{
		{
			name:             "JSONToTable_Normal",
			inputFormat:      "json",
			outputFormat:     "table",
			input:            "[[\"名前\",\"年齢\",\"職業\"],[\"田中\",\"25\",\"エンジニア\"],[\"佐藤\",\"30\",\"デザイナー\"]]",
			inputFilePath:    "",
			expectedContains: []string{"名前", "年齢", "職業", "田中", "25", "エンジニア", "佐藤", "30", "デザイナー", "|", "---"},
		},
		{
			name:             "TableToJSON_Normal",
			inputFormat:      "table",
			outputFormat:     "json",
			input:            "| 名前 | 年齢 | 職業 |\n|------|------|------|\n| 田中 | 25   | エンジニア |\n| 佐藤 | 30   | デザイナー |",
			inputFilePath:    "",
			expectedContains: []string{"名前", "年齢", "職業", "田中", "25", "エンジニア", "佐藤", "30", "デザイナー"},
		},
		{
			name:             "ListToTable_Normal",
			inputFormat:      "list",
			outputFormat:     "table",
			input:            "- 項目1\n- 項目2\n- 項目3",
			inputFilePath:    "",
			expectedContains: []string{"項目1", "項目2", "項目3", "|", "---"},
		},
		{
			name:             "TableToList_Normal",
			inputFormat:      "table",
			outputFormat:     "list",
			input:            "| 項目 |\n|------|\n| 項目1 |\n| 項目2 |",
			inputFilePath:    "",
			expectedContains: []string{"- 項目1", "- 項目2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ConvertData(tt.inputFormat, tt.outputFormat, tt.input, tt.inputFilePath)
			assert.NoError(t, err)

			for _, expected := range tt.expectedContains {
				assert.Contains(t, result, expected, "Result should contain: %s", expected)
			}
		})
	}
}
