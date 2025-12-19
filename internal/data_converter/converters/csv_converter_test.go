package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVConverter_ConvertToCSV_Normal(t *testing.T) {
	converter := NewCSVConverter()

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
			name: "WithQuotes",
			input: [][]string{
				{"Name", "Description"},
				{"Alice", "She said \"Hello\""},
				{"Bob", "He likes, commas"},
			},
			expected: "Name,Description\nAlice,\"She said \"\"Hello\"\"\"\nBob,\"He likes, commas\"",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertToCSV(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSVConverter_ConvertToTSV_Normal(t *testing.T) {
	converter := NewCSVConverter()

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
			name: "WithTabs",
			input: [][]string{
				{"Name", "Description"},
				{"Alice", "She likes\ttabs"},
				{"Bob", "Normal text"},
			},
			expected: "Name\tDescription\nAlice\t\"She likes\ttabs\"\nBob\tNormal text",
		},
		{
			name:     "EmptyCase",
			input:    [][]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertToTSV(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCSVConverter_EscapeCell_Normal(t *testing.T) {
	converter := NewCSVConverter()

	tests := []struct {
		name      string
		cell      string
		separator string
		expected  string
	}{
		{
			name:      "NormalText",
			cell:      "Alice",
			separator: ",",
			expected:  "Alice",
		},
		{
			name:      "WithComma",
			cell:      "Alice, Bob",
			separator: ",",
			expected:  "\"Alice, Bob\"",
		},
		{
			name:      "WithQuotes",
			cell:      "She said \"Hello\"",
			separator: ",",
			expected:  "\"She said \"\"Hello\"\"\"",
		},
		{
			name:      "WithNewline",
			cell:      "Line1\nLine2",
			separator: ",",
			expected:  "\"Line1\nLine2\"",
		},
		{
			name:      "WithTab",
			cell:      "Text\tTab",
			separator: "\t",
			expected:  "\"Text\tTab\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.escapeCell(tt.cell, tt.separator)
			assert.Equal(t, tt.expected, result)
		})
	}
}
