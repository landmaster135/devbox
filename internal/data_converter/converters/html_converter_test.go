package converters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTMLConverter_ConvertToHTML_Normal(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name                 string
		values               [][]string
		isTheadContained     bool
		textReplacingIfBlank string
		expected             string
	}{
		{
			name: "WithHeader_Normal",
			values: [][]string{
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
			name: "WithoutHeader_Normal",
			values: [][]string{
				{"Alice", "25", "New York"},
				{"Bob", "30", "London"},
			},
			isTheadContained:     false,
			textReplacingIfBlank: "💩",
			expected:             "<table>\n<tr><td>Alice</td><td>25</td><td>New York</td></tr>\n<tr><td>Bob</td><td>30</td><td>London</td></tr>\n</table>",
		},
		{
			name:                 "EmptyCase",
			values:               [][]string{},
			isTheadContained:     true,
			textReplacingIfBlank: "💩",
			expected:             "<table></table>",
		},
		{
			name: "SingleRow_WithHeader",
			values: [][]string{
				{"Name", "Age", "City"},
			},
			isTheadContained:     true,
			textReplacingIfBlank: "💩",
			expected:             "<table>\n<thead>\n<tr><th>Name</th><th>Age</th><th>City</th></tr>\n</thead>\n</table>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertToHTML(tt.values, tt.isTheadContained, tt.textReplacingIfBlank)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHTMLConverter_GetTrByRow_Normal(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name                 string
		row                  []string
		isTh                 bool
		textReplacingIfBlank string
		expected             string
	}{
		{
			name:                 "HeaderRow",
			row:                  []string{"Name", "Age", "City"},
			isTh:                 true,
			textReplacingIfBlank: "💩",
			expected:             "<tr><th>Name</th><th>Age</th><th>City</th></tr>",
		},
		{
			name:                 "DataRow",
			row:                  []string{"Alice", "25", "New York"},
			isTh:                 false,
			textReplacingIfBlank: "💩",
			expected:             "<tr><td>Alice</td><td>25</td><td>New York</td></tr>",
		},
		{
			name:                 "DataRowWithEmpty",
			row:                  []string{"Alice", "", "New York"},
			isTh:                 false,
			textReplacingIfBlank: "💩",
			expected:             "<tr><td>Alice</td><td>💩</td><td>New York</td></tr>",
		},
		{
			name:                 "HeaderRowWithEmpty",
			row:                  []string{"Name", "", "City"},
			isTh:                 true,
			textReplacingIfBlank: "💩",
			expected:             "<tr><th>Name</th><th></th><th>City</th></tr>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.getTrByRow(tt.row, tt.isTh, tt.textReplacingIfBlank)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHTMLConverter_CloseByTag_Normal(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name                 string
		tag                  string
		innerText            string
		textReplacingIfBlank string
		expected             string
	}{
		{
			name:                 "NormalTag",
			tag:                  "td",
			innerText:            "Alice",
			textReplacingIfBlank: "💩",
			expected:             "<td>Alice</td>",
		},
		{
			name:                 "EmptyInnerText",
			tag:                  "td",
			innerText:            "",
			textReplacingIfBlank: "💩",
			expected:             "<td>💩</td>",
		},
		{
			name:                 "HeaderTag",
			tag:                  "th",
			innerText:            "Name",
			textReplacingIfBlank: "💩",
			expected:             "<th>Name</th>",
		},
		{
			name:                 "EmptyReplacementText",
			tag:                  "th",
			innerText:            "",
			textReplacingIfBlank: "",
			expected:             "<th></th>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.closeByTag(tt.tag, tt.innerText, tt.textReplacingIfBlank)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHTMLConverter_CloseByTag_Panic(t *testing.T) {
	converter := NewHTMLConverter()

	tests := []struct {
		name      string
		tag       string
		innerText string
		textReplacingIfBlank string
	}{
		{
			name:      "TagStartsWithBracket",
			tag:       "<td",
			innerText: "test",
			textReplacingIfBlank: "",
		},
		{
			name:      "TagEndsWithBracket",
			tag:       "td>",
			innerText: "test",
			textReplacingIfBlank: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() {
				converter.closeByTag(tt.tag, tt.innerText, tt.textReplacingIfBlank)
			})
		})
	}
}
