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

func TestDataParser_ParseYAML_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name: "NormalCase",
			input: `- Name: Alice
  Age: "25"
  City: New York
- Name: Bob
  Age: "30"
  City: London
`,
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
			result, err := parser.parseYAML(tt.input)
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
		name          string
		inputFormat   string
		input         string
		inputFilePath string
		expected      [][]string
	}{
		{
			name:          "JSONInput_Normal",
			inputFormat:   "json",
			input:         "[[\"Name\",\"Age\",\"City\"],[\"Alice\",\"25\",\"New York\"],[\"Bob\",\"30\",\"London\"]]",
			inputFilePath: "",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:          "CSVInput_Normal",
			inputFormat:   "csv",
			input:         "Name,Age,City\nAlice,25,New York\nBob,30,London",
			inputFilePath: "",
			expected: [][]string{
				{"Name", "Age", "City"},
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
		},
		{
			name:        "YAMLInput_Normal",
			inputFormat: "yaml",
			input: `- Name: Alice
  Age: "25"
- Name: Bob
  Age: "30"
`,
			inputFilePath: "",
			expected: [][]string{
				{"Name", "Age"},
				{"Alice", "25"},
				{"Bob", "30"},
			},
		},
		{
			name:          "TSVInput_Normal",
			inputFormat:   "tsv",
			input:         "Name\tAge\tCity\nAlice\t25\tNew York\nBob\t30\tLondon",
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
		name          string
		inputFormat   string
		input         string
		inputFilePath string
		expectedError string
	}{
		{
			name:          "InvalidYAML_Error",
			inputFormat:   "yaml",
			input:         ":: invalid",
			inputFilePath: "",
			expectedError: "YAML解析エラー",
		},
		{
			name:          "UnsupportedFormat_Error",
			inputFormat:   "xml",
			input:         "test",
			inputFilePath: "",
			expectedError: "未対応の入力形式です: xml",
		},
		{
			name:          "InvalidJSON_Error",
			inputFormat:   "json",
			input:         "invalid json",
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
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "NoTable_Error",
			input:         "<div>No table here</div>",
			expectedError: "HTMLテーブルが見つかりません",
		},
		{
			name:          "NoClosingTable_Error",
			input:         "<table><tr><td>test</td></tr>",
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

func TestDataParser_ParseMarkdownTable_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name: "SimpleTable",
			input: `| 名前 | 年齢 | 職業 |
|------|------|------|
| 田中 | 25   | エンジニア |
| 佐藤 | 30   | デザイナー |`,
			expected: [][]string{
				{"名前", "年齢", "職業"},
				{"田中", "25", "エンジニア"},
				{"佐藤", "30", "デザイナー"},
			},
		},
		{
			name: "TableWithSpaces",
			input: `|  名前  |  年齢  |  職業  |
|--------|--------|--------|
|  田中  |   25   | エンジニア |
|  佐藤  |   30   | デザイナー |`,
			expected: [][]string{
				{"名前", "年齢", "職業"},
				{"田中", "25", "エンジニア"},
				{"佐藤", "30", "デザイナー"},
			},
		},
		{
			name: "SingleColumnTable",
			input: `| 項目 |
|------|
| 項目1 |
| 項目2 |`,
			expected: [][]string{
				{"項目"},
				{"項目1"},
				{"項目2"},
			},
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseMarkdownTable(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseMarkdownTable_Error(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name: "MissingPipes",
			input: `名前 | 年齢 | 職業
|------|------|------|
田中 | 25   | エンジニア`,
			expectedError: "Markdownテーブルの形式が正しくありません",
		},
		{
			name: "MissingSeparator",
			input: `| 名前 | 年齢 | 職業 |
| 田中 | 25   | エンジニア |
| 佐藤 | 30   | デザイナー |`,
			expectedError: "Markdownテーブルのセパレーター行が見つかりません",
		},
		{
			name: "InvalidFormat",
			input: `| 名前 | 年齢 | 職業
|------|------|------|
| 田中 | 25   | エンジニア |`,
			expectedError: "Markdownテーブルの形式が正しくありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseMarkdownTable(tt.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, result)
		})
	}
}

func TestDataParser_ParseMarkdownTableRow_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "SimpleRow",
			input:    "| 田中 | 25 | エンジニア |",
			expected: []string{"田中", "25", "エンジニア"},
		},
		{
			name:     "RowWithSpaces",
			input:    "|  田中  |  25  |  エンジニア  |",
			expected: []string{"田中", "25", "エンジニア"},
		},
		{
			name:     "SingleCell",
			input:    "| 項目1 |",
			expected: []string{"項目1"},
		},
		{
			name:     "EmptyCell",
			input:    "| 田中 |  | エンジニア |",
			expected: []string{"田中", "", "エンジニア"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parser.parseMarkdownTableRow(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseMarkdownList_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name: "SimpleList",
			input: `- 項目1
- 項目2
- 項目3`,
			expected: [][]string{
				{"項目"},
				{"項目1"},
				{"項目2"},
				{"項目3"},
			},
		},
		{
			name: "ListWithAsterisk",
			input: `* 項目1
* 項目2
* 項目3`,
			expected: [][]string{
				{"項目"},
				{"項目1"},
				{"項目2"},
				{"項目3"},
			},
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseMarkdownList(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseMarkdownList_Error(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name: "InvalidFormat",
			input: `項目1
項目2
項目3`,
			expectedError: "箇条書きリストの形式が正しくありません",
		},
		{
			name: "MixedFormat",
			input: `- 項目1
項目2
- 項目3`,
			expectedError: "箇条書きリストの形式が正しくありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseMarkdownList(tt.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, result)
		})
	}
}

func TestDataParser_ParseMarkdownOrderedList_Normal(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name     string
		input    string
		expected [][]string
	}{
		{
			name: "SimpleOrderedList",
			input: `1. 項目1
2. 項目2
3. 項目3`,
			expected: [][]string{
				{"番号", "項目"},
				{"1", "項目1"},
				{"2", "項目2"},
				{"3", "項目3"},
			},
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: [][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseMarkdownOrderedList(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataParser_ParseMarkdownOrderedList_Error(t *testing.T) {
	parser := NewDataParser()

	tests := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name: "InvalidFormat",
			input: `項目1
項目2
項目3`,
			expectedError: "順序付きリストの形式が正しくありません",
		},
		{
			name: "InvalidNumber",
			input: `a. 項目1
b. 項目2`,
			expectedError: "番号が正しくありません",
		},
		{
			name: "MissingDot",
			input: `1 項目1
2 項目2`,
			expectedError: "順序付きリストの形式が正しくありません",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.parseMarkdownOrderedList(tt.input)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			assert.Nil(t, result)
		})
	}
}
