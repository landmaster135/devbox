package infrastructure

import (
	"encoding/json"
	"testing"
)

// TestNewJSONProcessor_Normal はNewJSONProcessorメソッドの正常系テスト
func TestNewJSONProcessor_Normal(t *testing.T) {
	// Act
	result := NewJSONProcessor()

	// Assert
	if result == nil {
		t.Error("結果がnilです")
	}
}

// TestJSONProcessorImpl_MarshalIndent はJSONProcessorImplのMarshalIndentメソッドテスト
func TestJSONProcessorImpl_MarshalIndent(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		prefix   string
		indent   string
		expected string
	}{
		{
			name:     "SimpleObject_Normal",
			input:    map[string]string{"key": "value"},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"key\": \"value\"\n}",
		},
		{
			name:     "EmptyObject_Normal",
			input:    map[string]string{},
			prefix:   "",
			indent:   "  ",
			expected: "{}",
		},
		{
			name:     "Array_Normal",
			input:    []string{"item1", "item2", "item3"},
			prefix:   "",
			indent:   "  ",
			expected: "[\n  \"item1\",\n  \"item2\",\n  \"item3\"\n]",
		},
		{
			name:     "NestedObject_Normal",
			input:    map[string]interface{}{"outer": map[string]string{"inner": "value"}},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"outer\": {\n    \"inner\": \"value\"\n  }\n}",
		},
		{
			name:     "WithPrefix_Normal",
			input:    map[string]string{"key": "value"},
			prefix:   ">>",
			indent:   "  ",
			expected: "{\n>>  \"key\": \"value\"\n>>}",
		},
		{
			name:     "WithTabIndent_Normal",
			input:    map[string]string{"key": "value"},
			prefix:   "",
			indent:   "\t",
			expected: "{\n\t\"key\": \"value\"\n}",
		},
		{
			name:     "NumberValue_Normal",
			input:    map[string]interface{}{"number": 42, "float": 3.14},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"float\": 3.14,\n  \"number\": 42\n}",
		},
		{
			name:     "BooleanValue_Normal",
			input:    map[string]interface{}{"true": true, "false": false},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"false\": false,\n  \"true\": true\n}",
		},
		{
			name:     "NullValue_Normal",
			input:    map[string]interface{}{"null": nil},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"null\": null\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			processor := NewJSONProcessor()

			// Act
			result, err := processor.MarshalIndent(tt.input, tt.prefix, tt.indent)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("期待される結果: %s, 実際: %s", tt.expected, string(result))
			}
		})
	}
}

// TestJSONProcessorImpl_MarshalIndent_ComplexStruct はJSONProcessorImplの複雑な構造体テスト
func TestJSONProcessorImpl_MarshalIndent_ComplexStruct(t *testing.T) {
	// Arrange
	type TestStruct struct {
		ID       int      `json:"id"`
		Name     string   `json:"name"`
		Tags     []string `json:"tags"`
		Metadata struct {
			Created string `json:"created"`
			Updated string `json:"updated"`
		} `json:"metadata"`
	}

	testData := TestStruct{
		ID:   123,
		Name: "test item",
		Tags: []string{"tag1", "tag2"},
	}
	testData.Metadata.Created = "2023-01-01"
	testData.Metadata.Updated = "2023-01-02"

	processor := NewJSONProcessor()

	// Act
	result, err := processor.MarshalIndent(testData, "", "  ")

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
		return
	}

	// JSONが有効であることを確認
	var parsed TestStruct
	if err := json.Unmarshal(result, &parsed); err != nil {
		t.Errorf("生成されたJSONが無効です: %v", err)
		return
	}

	// 元のデータと一致することを確認
	if parsed.ID != testData.ID {
		t.Errorf("期待されるID: %d, 実際: %d", testData.ID, parsed.ID)
	}
	if parsed.Name != testData.Name {
		t.Errorf("期待される名前: %s, 実際: %s", testData.Name, parsed.Name)
	}
	if len(parsed.Tags) != len(testData.Tags) {
		t.Errorf("期待されるタグ数: %d, 実際: %d", len(testData.Tags), len(parsed.Tags))
	}
}

// TestJSONProcessorImpl_MarshalIndent_InvalidInput はJSONProcessorImplの無効入力テスト
func TestJSONProcessorImpl_MarshalIndent_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "CircularReference_Error",
			input: createCircularReference(),
		},
		{
			name:  "InvalidChannel_Error",
			input: make(chan int),
		},
		{
			name:  "Function_Error",
			input: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			processor := NewJSONProcessor()

			// Act
			result, err := processor.MarshalIndent(tt.input, "", "  ")

			// Assert
			if err == nil {
				t.Error("エラーが期待されましたが、nilが返されました")
			}
			if result != nil {
				t.Error("エラー時は結果がnilである必要があります")
			}
		})
	}
}

// createCircularReference は循環参照を持つ構造体を作成する
func createCircularReference() map[string]interface{} {
	m := make(map[string]interface{})
	m["self"] = m
	return m
}

// TestJSONProcessorImpl_MarshalIndent_EmptyInput はJSONProcessorImplの空入力テスト
func TestJSONProcessorImpl_MarshalIndent_EmptyInput(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "NilInput_Normal",
			input:    nil,
			expected: "null",
		},
		{
			name:     "EmptyString_Normal",
			input:    "",
			expected: "\"\"",
		},
		{
			name:     "ZeroInt_Normal",
			input:    0,
			expected: "0",
		},
		{
			name:     "EmptySlice_Normal",
			input:    []string{},
			expected: "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			processor := NewJSONProcessor()

			// Act
			result, err := processor.MarshalIndent(tt.input, "", "  ")

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("期待される結果: %s, 実際: %s", tt.expected, string(result))
			}
		})
	}
}

// TestJSONProcessorImpl_MarshalIndent_SpecialCharacters はJSONProcessorImplの特殊文字テスト
func TestJSONProcessorImpl_MarshalIndent_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		prefix   string
		indent   string
		expected string
	}{
		{
			name:     "QuotesInString_Normal",
			input:    map[string]string{"key": "value with \"quotes\""},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"key\": \"value with \\\"quotes\\\"\"\n}",
		},
		{
			name:     "NewlineInString_Normal",
			input:    map[string]string{"key": "line1\nline2"},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"key\": \"line1\\nline2\"\n}",
		},
		{
			name:     "UnicodeCharacters_Normal",
			input:    map[string]string{"key": "こんにちは世界"},
			prefix:   "",
			indent:   "  ",
			expected: "{\n  \"key\": \"こんにちは世界\"\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			processor := NewJSONProcessor()

			// Act
			result, err := processor.MarshalIndent(tt.input, tt.prefix, tt.indent)

			// Assert
			if err != nil {
				t.Errorf("エラーが発生しました: %v", err)
				return
			}

			if string(result) != tt.expected {
				t.Errorf("期待される結果: %s, 実際: %s", tt.expected, string(result))
			}
		})
	}
}
