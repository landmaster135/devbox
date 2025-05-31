package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewFormatter(t *testing.T) {
	formatter := NewFormatter()
	if formatter == nil {
		t.Error("NewFormatter() should return a valid Formatter instance")
	}
}

func TestFormatNumber(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		input    int
		expected string
	}{
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
	}

	for _, test := range tests {
		result := formatter.formatNumber(test.input)
		if result != test.expected {
			t.Errorf("formatNumber(%d) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestFormatCost(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		input    float64
		expected string
	}{
		{1.5, "$1.50"},
		{10.0, "$10.00"},
		{123.456, "$123.46"},
		{0.99, "$0.99"},
	}

	for _, test := range tests {
		result := formatter.formatCost(test.input)
		if result != test.expected {
			t.Errorf("formatCost(%.3f) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestOutputJSON(t *testing.T) {
	formatter := NewFormatter()

	// テスト用のデータ
	testData := []DailyUsage{
		{
			Date:                "2025-05-30",
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			TotalTokens:         375,
			Cost:                1.50,
		},
		{
			Date:                "2025-05-29",
			InputTokens:         150,
			OutputTokens:        300,
			CacheCreationTokens: 75,
			CacheReadTokens:     50,
			TotalTokens:         575,
			Cost:                2.25,
		},
	}

	// 標準出力をキャプチャするための準備
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// JSON出力をテスト
	err := formatter.outputJSON(testData)
	if err != nil {
		t.Fatalf("outputJSON() returned error: %v", err)
	}

	// 出力をキャプチャして復元
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// JSON形式として正しいかテスト
	var parsedData []DailyUsage
	err = json.Unmarshal([]byte(output), &parsedData)
	if err != nil {
		t.Fatalf("outputJSON() produced invalid JSON: %v", err)
	}

	// データが正しく出力されているかテスト
	if len(parsedData) != len(testData) {
		t.Errorf("outputJSON() length mismatch: got %d, expected %d", len(parsedData), len(testData))
	}

	for i, data := range parsedData {
		if data.Date != testData[i].Date {
			t.Errorf("outputJSON() date mismatch at index %d: got %s, expected %s", i, data.Date, testData[i].Date)
		}
		if data.InputTokens != testData[i].InputTokens {
			t.Errorf("outputJSON() input tokens mismatch at index %d: got %d, expected %d", i, data.InputTokens, testData[i].InputTokens)
		}
	}
}

func TestFormatDailyReport_JSON(t *testing.T) {
	formatter := NewFormatter()

	testData := []DailyUsage{
		{
			Date:                "2025-05-30",
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			TotalTokens:         375,
			Cost:                1.50,
		},
	}

	config := &Config{JSON: true}

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := formatter.FormatDailyReport(testData, config)
	if err != nil {
		t.Fatalf("FormatDailyReport() returned error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// JSON形式として正しいかテスト
	var parsedData []DailyUsage
	err = json.Unmarshal([]byte(output), &parsedData)
	if err != nil {
		t.Fatalf("FormatDailyReport() produced invalid JSON: %v", err)
	}

	if len(parsedData) != 1 {
		t.Errorf("FormatDailyReport() length mismatch: got %d, expected 1", len(parsedData))
	}
}

func TestFormatSessionReport_JSON(t *testing.T) {
	formatter := NewFormatter()

	testData := []SessionUsage{
		{
			Project:             "myproject",
			Session:             "session-1",
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			TotalTokens:         375,
			Cost:                1.50,
			LastActivity:        time.Date(2025, 5, 30, 12, 0, 0, 0, time.UTC),
		},
	}

	config := &Config{JSON: true}

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := formatter.FormatSessionReport(testData, config)
	if err != nil {
		t.Fatalf("FormatSessionReport() returned error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// JSON形式として正しいかテスト
	var parsedData []SessionUsage
	err = json.Unmarshal([]byte(output), &parsedData)
	if err != nil {
		t.Fatalf("FormatSessionReport() produced invalid JSON: %v", err)
	}

	if len(parsedData) != 1 {
		t.Errorf("FormatSessionReport() length mismatch: got %d, expected 1", len(parsedData))
	}

	if parsedData[0].Project != "myproject" {
		t.Errorf("FormatSessionReport() project mismatch: got %s, expected myproject", parsedData[0].Project)
	}
}

func TestFormatDailyReport_Table(t *testing.T) {
	formatter := NewFormatter()

	testData := []DailyUsage{
		{
			Date:                "2025-05-30",
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			TotalTokens:         375,
			Cost:                1.50,
		},
		{
			Date:                "Total",
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			TotalTokens:         375,
			Cost:                1.50,
		},
	}

	config := &Config{JSON: false}

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := formatter.FormatDailyReport(testData, config)
	if err != nil {
		t.Fatalf("FormatDailyReport() returned error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// テーブル出力の基本的な要素がある事を確認
	expectedStrings := []string{
		"Claude Code Token Usage Report - Daily",
		"Date",
		"Input",
		"Output",
		"Cache Create",
		"Cache Read",
		"Total Tokens",
		"Cost (USD)",
		"2025-05-30",
		"Total",
		"100",
		"200",
		"$1.50",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("FormatDailyReport() table output missing expected string: %s", expected)
		}
	}

	// 罫線文字がある事を確認
	boxDrawingChars := []string{"┌", "┐", "├", "┤", "└", "┘", "│", "─"}
	for _, char := range boxDrawingChars {
		if !strings.Contains(output, char) {
			t.Errorf("FormatDailyReport() table output missing box drawing character: %s", char)
		}
	}
}

func TestFormatSessionReport_Table(t *testing.T) {
	formatter := NewFormatter()

	testData := []SessionUsage{
		{
			Project:             "myproject",
			Session:             "session-1",
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			TotalTokens:         375,
			Cost:                1.50,
			LastActivity:        time.Date(2025, 5, 30, 12, 0, 0, 0, time.UTC),
		},
	}

	config := &Config{JSON: false}

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := formatter.FormatSessionReport(testData, config)
	if err != nil {
		t.Fatalf("FormatSessionReport() returned error: %v", err)
	}

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// テーブル出力の基本的な要素がある事を確認
	expectedStrings := []string{
		"Claude Code Token Usage Report - By Session",
		"Project",
		"Session",
		"Input",
		"Output",
		"Last Activity",
		"myproject",
		"session-1",
		"2025-05-30",
		"$1.50",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("FormatSessionReport() table output missing expected string: %s", expected)
		}
	}
}

func TestPrintHeader(t *testing.T) {
	formatter := NewFormatter()

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	title := "Test Header"
	formatter.printHeader(title)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// ヘッダーの要素が含まれている事を確認
	if !strings.Contains(output, title) {
		t.Errorf("printHeader() output missing title: %s", title)
	}

	// ボックス描画文字が含まれている事を確認
	headerChars := []string{"╭", "╮", "╰", "╯", "─", "│"}
	for _, char := range headerChars {
		if !strings.Contains(output, char) {
			t.Errorf("printHeader() output missing header character: %s", char)
		}
	}
}

func TestPrintTableBorder(t *testing.T) {
	formatter := NewFormatter()

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	widths := []int{10, 8, 12}
	formatter.printTableBorder(widths...)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// ボーダー文字が含まれている事を確認
	borderChars := []string{"┌", "┐", "┬", "─"}
	for _, char := range borderChars {
		if !strings.Contains(output, char) {
			t.Errorf("printTableBorder() output missing border character: %s", char)
		}
	}
}

func TestPrintTableSeparator(t *testing.T) {
	formatter := NewFormatter()

	// 標準出力をキャプチャ
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	widths := []int{10, 8, 12}
	formatter.printTableSeparator(widths...)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	// セパレーター文字が含まれている事を確認
	separatorChars := []string{"├", "┤", "┼", "─"}
	for _, char := range separatorChars {
		if !strings.Contains(output, char) {
			t.Errorf("printTableSeparator() output missing separator character: %s", char)
		}
	}
}

// ベンチマークテスト
func BenchmarkFormatNumber(b *testing.B) {
	formatter := NewFormatter()
	for i := 0; i < b.N; i++ {
		formatter.formatNumber(1234567)
	}
}

func BenchmarkFormatCost(b *testing.B) {
	formatter := NewFormatter()
	for i := 0; i < b.N; i++ {
		formatter.formatCost(123.456)
	}
}

// エッジケースのテスト
func TestFormatNumber_EdgeCases(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		input    int
		expected string
	}{
		{-1, "-1"},
		{-1234, "-1,234"},
		{-12345678, "-12,345,678"},
	}

	for _, test := range tests {
		result := formatter.formatNumber(test.input)
		if result != test.expected {
			t.Errorf("formatNumber(%d) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestFormatCost_EdgeCases(t *testing.T) {
	formatter := NewFormatter()

	tests := []struct {
		input    float64
		expected string
	}{
		{-1.5, "$-1.50"},
		{0.0, "$0.00"},
		{-0.001, "$-0.00"},
	}

	for _, test := range tests {
		result := formatter.formatCost(test.input)
		if result != test.expected {
			t.Errorf("formatCost(%.3f) = %s; expected %s", test.input, result, test.expected)
		}
	}
}
