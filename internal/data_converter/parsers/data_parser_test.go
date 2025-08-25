package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataParser_ParseJSON_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:  "NormalCase",
			input: "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:     "EmptyCase",
			input:    "[]",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseJSON(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseCSV_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:  "NormalCase",
			input: "Name,Age,City\nAlice,25,New York\nBob,30,London",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:     "EmptyCase",
			input:    "",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseCSV(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseTSV_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:  "NormalCase",
			input: "Name\tAge\tCity\nAlice\t25\tNew York\nBob\t30\tLondon",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:     "EmptyCase",
			input:    "",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseTSV(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseInput_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name         string
		inputFormat  string
		input        string
		inputFilePath string
		expected     [][]string
	}{
		{
			name:         "JSONInput_Normal",
			inputFormat:  "json",
			input:        "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:         "CSVInput_Normal",
			inputFormat:  "csv",
			input:        "Name,Age,City\nAlice,25,New York\nBob,30,London",
			inputFilePath: "",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:         "TSVInput_Normal",
			inputFormat:  "tsv",
			input:        "Name\tAge\tCity\nAlice\t25\tNew York\nBob\t30\tLondon",
			inputFilePath: "",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseInput(tt.inputFormat, tt.input, tt.inputFilePath)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseInput_Error(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name         string
		inputFormat  string
		input        string
		inputFilePath string
		expectedError string
	}{
		{
			name:         "UnsupportedFormat_Error",
			inputFormat:  "xml",
			input:        "test",
			inputFilePath: "",
			expectedError: "未対応の入力形式です: xml",
		},
		{
			name:         "InvalidJSON_Error",
			inputFormat:  "json",
			input:        "invalid json",
			inputFilePath: "",
			expectedError: "JSON解析エラー:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.ParseInput(tt.inputFormat, tt.input, tt.inputFilePath)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, result)
		})
	}
}

func TestDataParser_ParseHTML_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name:  "SimpleTable",
			input: "<table><tr><th>Name</th><th>Age</th></tr><tr><td>Alice</td><td>25</td></tr></table>",
			expected: [][]string{
				{"Name", "Age"},
				{"Alice", "25"},
			},
		},
		{
			name:  "TableWithSpaces",
			input: "<table> <tr> <th>Name</th> <th>Age</th> </tr> <tr> <td>Alice</td> <td>25</td> </tr> </table>",
			expected: [][]string{
				{"Name", "Age"},
				{"Alice", "25"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseHTML(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseHTML_Error(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name         string
		input        string
		expectedError string
	}{
		{
			name:         "NoTable_Error",
			input:        "<div>No table here</div>",
			expectedError: "HTMLテーブルが見つかりません",
		},
		{
			name:         "NoClosingTable_Error",
			input:        "<table><tr><td>test</td></tr>",
			expectedError: "HTMLテーブルの終了タグが見つかりません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseHTML(tt.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, result)
		})
	}
}

func TestDataParser_ParseDelimitedLine_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name      string
		line      string
		delimiter string
		expected  []string
	}{
		{
			name:      "SimpleCSV",
			line:      "Alice,25,New York",
			delimiter: ",",
			expected:  []string{"Alice", "25", "New York"},
		},
		{
			name:      "QuotedCSV",
			line:      "\"Alice\",\"25\",\"New York\"",
			delimiter: ",",
			expected:  []string{"Alice", "25", "New York"},
		},
		{
			name:      "EscapedQuotes",
			line:      "\"She said \"\"Hello\"\"\",25",
			delimiter: ",",
			expected:  []string{"She said \"Hello\"", "25"},
		},
		{
			name:      "TSV",
			line:      "Alice\t25\tNew York",
			delimiter: "\t",
			expected:  []string{"Alice", "25", "New York"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseDelimitedLine(tt.line, tt.delimiter)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_StripHTMLTags_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SimpleTag",
			input:    "<b>Bold</b>",
			expected: "Bold",
		},
		{
			name:     "MultipleTag",
			input:    "<div><span>Text</span></div>",
			expected: "Text",
		},
		{
			name:     "WithEntities",
			input:    "&lt;test&gt; &amp; &quot;quote&quot;",
			expected: "<test> & \"quote\"",
		},
		{
			name:     "NoTags",
			input:    "Plain text",
			expected: "Plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.stripHTMLTags(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
