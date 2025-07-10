package datetime_calc

import (
	"context"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
)

func TestHandleDatetimeCalc_Normal(t *testing.T) {
	tests := []struct {
		name     string
		request  map[string]interface{}
		expected string
	}{
		{
			name: "加算テスト - 1年追加",
			request: map[string]interface{}{
				"operation":          "add",
				"year_1":             2023.0,
				"month_1":            5.0,
				"day_1":              15.0,
				"hour_1":             10.0,
				"minute_1":           30.0,
				"second_1":           45.0,
				"duration_of_year":   1.0,
				"duration_of_month":  0.0,
				"duration_of_day":    0.0,
				"duration_of_hour":   0.0,
				"duration_of_minute": 0.0,
				"duration_of_second": 0.0,
			},
			expected: "2024-05-15 10:30:45",
		},
		{
			name: "減算テスト - 1年減算",
			request: map[string]interface{}{
				"operation":          "subtract",
				"year_1":             2023.0,
				"month_1":            5.0,
				"day_1":              15.0,
				"hour_1":             10.0,
				"minute_1":           30.0,
				"second_1":           45.0,
				"duration_of_year":   1.0,
				"duration_of_month":  0.0,
				"duration_of_day":    0.0,
				"duration_of_hour":   0.0,
				"duration_of_minute": 0.0,
				"duration_of_second": 0.0,
			},
			expected: "2022-05-15 10:30:45",
		},
		{
			name: "複合的な加算テスト",
			request: map[string]interface{}{
				"operation":          "add",
				"year_1":             2023.0,
				"month_1":            5.0,
				"day_1":              15.0,
				"hour_1":             10.0,
				"minute_1":           30.0,
				"second_1":           45.0,
				"duration_of_year":   1.0,
				"duration_of_month":  2.0,
				"duration_of_day":    3.0,
				"duration_of_hour":   4.0,
				"duration_of_minute": 5.0,
				"duration_of_second": 6.0,
			},
			expected: "2024-07-18 14:35:51",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// MCPリクエストを作成
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "datetime_calc",
					Arguments: tt.request,
				},
			}

			// ハンドラー関数を直接テスト
			ctx := context.Background()
			result, err := handleDatetimeCalc(ctx, request)

			if err != nil {
				t.Errorf("Handler returned error: %v", err)
				return
			}

			// 結果を検証
			if result == nil {
				t.Error("Result is nil")
				return
			}

			if len(result.Content) == 0 {
				t.Error("Result content is empty")
				return
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Error("Result content is not TextContent")
				return
			}

			if textContent.Text != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, textContent.Text)
			}
		})
	}
}

func TestHandleDatetimeCalc_Error(t *testing.T) {
	tests := []struct {
		name    string
		request map[string]interface{}
	}{
		{
			name: "無効な操作",
			request: map[string]interface{}{
				"operation":          "invalid",
				"year_1":             2023.0,
				"month_1":            5.0,
				"day_1":              15.0,
				"hour_1":             10.0,
				"minute_1":           30.0,
				"second_1":           45.0,
				"duration_of_year":   1.0,
				"duration_of_month":  0.0,
				"duration_of_day":    0.0,
				"duration_of_hour":   0.0,
				"duration_of_minute": 0.0,
				"duration_of_second": 0.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// MCPリクエストを作成
			request := mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "datetime_calc",
					Arguments: tt.request,
				},
			}

			// ハンドラー関数を直接テスト
			ctx := context.Background()
			_, err := handleDatetimeCalc(ctx, request)

			// エラーが発生することを期待
			if err == nil {
				t.Error("Expected error but got none")
			}
		})
	}
}
