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

func TestMarkdownConverter_ConvertToTable_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 複数列のテーブルデータ
	data := [][]string{
		{"名前", "年齢", "職業"},
		{"田中", "25", "エンジニア"},
		{"佐藤", "30", "デザイナー"},
	}

	result := converter.ConvertToTable(data)
	lines := strings.Split(strings.TrimSpace(result), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (header + separator + 2 data rows), got %d", len(lines))
	}

	// ヘッダー行をチェック
	if !strings.Contains(lines[0], "名前") || !strings.Contains(lines[0], "年齢") || !strings.Contains(lines[0], "職業") {
		t.Errorf("Header line should contain all column names, got %q", lines[0])
	}

	// セパレーター行をチェック
	if !strings.Contains(lines[1], "---") {
		t.Errorf("Separator line should contain '---', got %q", lines[1])
	}

	// データ行をチェック
	if !strings.Contains(lines[2], "田中") || !strings.Contains(lines[2], "25") || !strings.Contains(lines[2], "エンジニア") {
		t.Errorf("First data line should contain '田中', '25', 'エンジニア', got %q", lines[2])
	}

	if !strings.Contains(lines[3], "佐藤") || !strings.Contains(lines[3], "30") || !strings.Contains(lines[3], "デザイナー") {
		t.Errorf("Second data line should contain '佐藤', '30', 'デザイナー', got %q", lines[3])
	}
}

func TestMarkdownConverter_ConvertToTable_SingleColumn_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 1列のテーブルデータ（箇条書きリストから変換されたデータ）
	data := [][]string{
		{"項目"},
		{"項目1"},
		{"項目2"},
	}

	result := converter.ConvertToTable(data)
	lines := strings.Split(strings.TrimSpace(result), "\n")

	if len(lines) != 4 {
		t.Errorf("Expected 4 lines (header + separator + 2 data rows), got %d", len(lines))
	}

	// 1列目が空で2列目に項目が入っていることを確認
	if !strings.Contains(lines[0], "|     | 項目  |") {
		t.Errorf("Header should have empty first column and '項目' in second column, got %q", lines[0])
	}

	if !strings.Contains(lines[2], "|     | 項目1 |") {
		t.Errorf("First data line should have empty first column and '項目1' in second column, got %q", lines[2])
	}

	if !strings.Contains(lines[3], "|     | 項目2 |") {
		t.Errorf("Second data line should have empty first column and '項目2' in second column, got %q", lines[3])
	}
}

func TestMarkdownConverter_ConvertToTable_EmptyData_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	data := [][]string{}
	result := converter.ConvertToTable(data)

	if result != "" {
		t.Errorf("ConvertToTable() with empty data = %q, want empty string", result)
	}
}

func TestMarkdownConverter_AddEmptyColumnIfSingleColumn_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 1列のデータ
	data := [][]string{
		{"項目"},
		{"項目1"},
		{"項目2"},
	}

	result := converter.addEmptyColumnIfSingleColumn(data)

	if len(result) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(result))
	}

	for i, row := range result {
		if len(row) != 2 {
			t.Errorf("Row %d should have 2 columns, got %d", i, len(row))
		}

		if row[0] != "" {
			t.Errorf("Row %d first column should be empty, got %q", i, row[0])
		}

		if i == 0 && row[1] != "項目" {
			t.Errorf("Header row second column should be '項目', got %q", row[1])
		}
	}
}

func TestMarkdownConverter_AddEmptyColumnIfSingleColumn_MultiColumn_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 複数列のデータ（変更されないはず）
	data := [][]string{
		{"名前", "年齢"},
		{"田中", "25"},
	}

	result := converter.addEmptyColumnIfSingleColumn(data)

	// 元のデータと同じであることを確認
	if len(result) != len(data) {
		t.Errorf("Expected %d rows, got %d", len(data), len(result))
	}

	for i, row := range result {
		if len(row) != len(data[i]) {
			t.Errorf("Row %d should have %d columns, got %d", i, len(data[i]), len(row))
		}

		for j, cell := range row {
			if cell != data[i][j] {
				t.Errorf("Row %d column %d should be %q, got %q", i, j, data[i][j], cell)
			}
		}
	}
}

func TestMarkdownConverter_AddEmptyColumnIfSingleColumn_EmptyData_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	data := [][]string{}
	result := converter.addEmptyColumnIfSingleColumn(data)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d rows", len(result))
	}
}

func TestMarkdownConverter_CalculateColumnWidths_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	data := [][]string{
		{"名前", "年齢", "職業"},
		{"田中太郎", "25", "エンジニア"},
		{"佐藤", "30", "デザイナー"},
	}

	widths := converter.calculateColumnWidths(data)

	if len(widths) != 3 {
		t.Errorf("Expected 3 column widths, got %d", len(widths))
	}

	// "田中太郎"が最長なので4文字
	if widths[0] < 4 {
		t.Errorf("First column width should be at least 4, got %d", widths[0])
	}

	// "年齢"と"30"で"年齢"が長いので2文字（ただし最小幅3）
	if widths[1] < 3 {
		t.Errorf("Second column width should be at least 3, got %d", widths[1])
	}

	// "エンジニア"が最長なので5文字
	if widths[2] < 5 {
		t.Errorf("Third column width should be at least 5, got %d", widths[2])
	}
}

func TestMarkdownConverter_PadCell_Normal(t *testing.T) {
	converter := NewMarkdownConverter()

	// 短いセルをパディング
	result := converter.padCell("test", 10)
	if len(result) != 10 {
		t.Errorf("Padded cell should be 10 characters, got %d", len(result))
	}

	if !strings.HasPrefix(result, "test") {
		t.Errorf("Padded cell should start with 'test', got %q", result)
	}

	// 長いセルはそのまま
	longText := "this is a very long text"
	result = converter.padCell(longText, 10)
	if result != longText {
		t.Errorf("Long cell should remain unchanged, got %q", result)
	}

	// 同じ長さのセル
	result = converter.padCell("exact", 5)
	if result != "exact" {
		t.Errorf("Exact length cell should remain unchanged, got %q", result)
	}
}
