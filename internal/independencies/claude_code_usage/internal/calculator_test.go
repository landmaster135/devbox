package internal

import (
	"testing"
	"time"
)

func TestCalculateTotal(t *testing.T) {
	calc := NewCalculator()

	testData := []UsageData{
		{
			Timestamp:           time.Now(),
			InputTokens:         100,
			OutputTokens:        200,
			CacheCreationTokens: 50,
			CacheReadTokens:     25,
			Cost:                1.50,
		},
		{
			Timestamp:           time.Now(),
			InputTokens:         150,
			OutputTokens:        300,
			CacheCreationTokens: 75,
			CacheReadTokens:     50,
			Cost:                2.25,
		},
	}

	input, output, cacheCreate, cacheRead, totalCost := calc.CalculateTotalUsage(testData)

	expectedInput := 250
	expectedOutput := 500
	expectedCacheCreate := 125
	expectedCacheRead := 75
	expectedCost := 3.75

	if input != expectedInput {
		t.Errorf("Expected input tokens %d, got %d", expectedInput, input)
	}
	if output != expectedOutput {
		t.Errorf("Expected output tokens %d, got %d", expectedOutput, output)
	}
	if cacheCreate != expectedCacheCreate {
		t.Errorf("Expected cache creation tokens %d, got %d", expectedCacheCreate, cacheCreate)
	}
	if cacheRead != expectedCacheRead {
		t.Errorf("Expected cache read tokens %d, got %d", expectedCacheRead, cacheRead)
	}
	if totalCost != expectedCost {
		t.Errorf("Expected total cost %.2f, got %.2f", expectedCost, totalCost)
	}
}

func TestAggregateDailyUsage(t *testing.T) {
	calc := NewCalculator()

	date1, _ := time.Parse("2006-01-02", "2025-05-30")
	date2, _ := time.Parse("2006-01-02", "2025-05-31")

	testData := []UsageData{
		{
			Timestamp:    date1,
			InputTokens:  100,
			OutputTokens: 200,
			Cost:         1.50,
		},
		{
			Timestamp:    date1,
			InputTokens:  50,
			OutputTokens: 100,
			Cost:         0.75,
		},
		{
			Timestamp:    date2,
			InputTokens:  75,
			OutputTokens: 150,
			Cost:         1.125,
		},
	}

	dailyUsage := calc.AggregateDailyUsage(testData)

	if len(dailyUsage) != 2 {
		t.Errorf("Expected 2 daily entries, got %d", len(dailyUsage))
	}

	// Check first day aggregation
	day1 := dailyUsage[0]
	if day1.Date != "2025-05-30" {
		t.Errorf("Expected date 2025-05-30, got %s", day1.Date)
	}
	if day1.InputTokens != 150 {
		t.Errorf("Expected 150 input tokens for day1, got %d", day1.InputTokens)
	}
	if day1.OutputTokens != 300 {
		t.Errorf("Expected 300 output tokens for day1, got %d", day1.OutputTokens)
	}
	if day1.Cost != 2.25 {
		t.Errorf("Expected cost 2.25 for day1, got %.2f", day1.Cost)
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
