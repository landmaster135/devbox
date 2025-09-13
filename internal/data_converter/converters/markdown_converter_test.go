package converters

import (
	"strings"
	"testing"
)

func TestMarkdownConverter_ConvertToList_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 1列のテーブルデータ
	data := [][]string{
		{"項目"},
		{"項目1"},
		{"項目2"},
		{"項目3"},
	}

	result := converter.ConvertToList(data)
	expected := "- 項目1\n- 項目2\n- 項目3\n"

	if result != expected {
		t.Errorf("ConvertToList() = %q, want %q", result, expected)
	}
}

func TestMarkdownConverter_ConvertToList_MultiColumn_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 複数列のテーブルデータ
	data := [][]string{
		{"名前", "年齢"},
		{"田中", "25"},
		{"佐藤", "30"},
	}

	result := converter.ConvertToList(data)
	lines := strings.Split(strings.TrimSpace(result), "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if !strings.HasPrefix(lines[0], "- ") {
		t.Errorf("First line should start with '- ', got %q", lines[0])
	}

	if !strings.Contains(lines[0], "名前: 田中") {
		t.Errorf("First line should contain '名前: 田中', got %q", lines[0])
	}

	if !strings.Contains(lines[0], "年齢: 25") {
		t.Errorf("First line should contain '年齢: 25', got %q", lines[0])
	}
}

func TestMarkdownConverter_ConvertToList_EmptyData_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	data := [][]string{}
	result := converter.ConvertToList(data)

	if result != "" {
		t.Errorf("ConvertToList() with empty data = %q, want empty string", result)
	}
}

func TestMarkdownConverter_ConvertToOrderedList_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 1列のテーブルデータ
	data := [][]string{
		{"項目"},
		{"項目1"},
		{"項目2"},
		{"項目3"},
	}

	result := converter.ConvertToOrderedList(data)
	expected := "1. 項目1\n2. 項目2\n3. 項目3\n"

	if result != expected {
		t.Errorf("ConvertToOrderedList() = %q, want %q", result, expected)
	}
}

func TestMarkdownConverter_ConvertToOrderedList_MultiColumn_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 複数列のテーブルデータ
	data := [][]string{
		{"名前", "年齢"},
		{"田中", "25"},
		{"佐藤", "30"},
	}

	result := converter.ConvertToOrderedList(data)
	lines := strings.Split(strings.TrimSpace(result), "\n")

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if !strings.HasPrefix(lines[0], "1. ") {
		t.Errorf("First line should start with '1. ', got %q", lines[0])
	}

	if !strings.HasPrefix(lines[1], "2. ") {
		t.Errorf("Second line should start with '2. ', got %q", lines[1])
	}

	if !strings.Contains(lines[0], "名前: 田中") {
		t.Errorf("First line should contain '名前: 田中', got %q", lines[0])
	}

	if !strings.Contains(lines[0], "年齢: 25") {
		t.Errorf("First line should contain '年齢: 25', got %q", lines[0])
	}
}

func TestMarkdownConverter_ConvertToOrderedList_EmptyData_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	data := [][]string{}
	result := converter.ConvertToOrderedList(data)

	if result != "" {
		t.Errorf("ConvertToOrderedList() with empty data = %q, want empty string", result)
	}
}

func TestMarkdownConverter_ConvertToList_SingleRow_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 1行のみのデータ（ヘッダーなしと判断される）
	data := [][]string{
		{"項目1"},
	}

	result := converter.ConvertToList(data)
	expected := "- 項目1\n"

	if result != expected {
		t.Errorf("ConvertToList() with single row = %q, want %q", result, expected)
	}
}

func TestMarkdownConverter_ConvertToOrderedList_SingleRow_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 1行のみのデータ（ヘッダーなしと判断される）
	data := [][]string{
		{"項目1"},
	}

	result := converter.ConvertToOrderedList(data)
	expected := "1. 項目1\n"

	if result != expected {
		t.Errorf("ConvertToOrderedList() with single row = %q, want %q", result, expected)
	}
}
