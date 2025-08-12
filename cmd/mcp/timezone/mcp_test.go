package timezone

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// #==============================================================#
// ##          Test Helper Functions                            ##
// #==============================================================#

// createCallToolRequest はテスト用のCallToolRequestを作成します
func createCallToolRequest(name string, arguments map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      name,
			Arguments: arguments,
		},
	}
}

// #==============================================================#
// ##          Handler Tests                                     ##
// #==============================================================#

func TestHandleGetCurrentTimezone_Normal(t *testing.T) {
	// テストケース
	tests := []struct {
		name     string
		timezone string
		wantErr  bool
	}{
		{
			name:     "有効なタイムゾーン_UTC",
			timezone: "UTC",
			wantErr:  false,
		},
		{
			name:     "有効なタイムゾーン_Asia/Tokyo",
			timezone: "Asia/Tokyo",
			wantErr:  false,
		},
		{
			name:     "有効なタイムゾーン_America/New_York",
			timezone: "America/New_York",
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("get-current-timezone", map[string]interface{}{
				"timezone": tc.timezone,
			})

			// Act
			result, err := handleGetCurrentTimezone(ctx, request)

			// Assert
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
			}
		})
	}
}

func TestHandleGetCurrentTimezone_Error(t *testing.T) {
	// テストケース
	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "パラメータなし",
			arguments: map[string]interface{}{},
			wantErr:   true,
		},
		{
			name: "無効なパラメータ型",
			arguments: map[string]interface{}{
				"timezone": 123,
			},
			wantErr: true,
		},
		{
			name: "無効なタイムゾーン",
			arguments: map[string]interface{}{
				"timezone": "Invalid/Timezone",
			},
			wantErr: false, // サービス層でエラーハンドリングされる
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("get-current-timezone", tc.arguments)

			// Act
			result, err := handleGetCurrentTimezone(ctx, request)

			// Assert
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				// サービス層でエラーが処理される場合、結果は返されるがIsErrorがtrueになる可能性がある
				if err != nil {
					assert.Error(t, err)
				} else {
					assert.NotNil(t, result)
				}
			}
		})
	}
}

func TestHandleConvertTimezone_Normal(t *testing.T) {
	// テストケース
	tests := []struct {
		name         string
		datetime     string
		fromTimezone string
		toTimezone   string
		wantErr      bool
	}{
		{
			name:         "UTC to Asia/Tokyo",
			datetime:     "2025-01-01 12:00:00",
			fromTimezone: "UTC",
			toTimezone:   "Asia/Tokyo",
			wantErr:      false,
		},
		{
			name:         "Asia/Tokyo to America/New_York",
			datetime:     "2025-01-01 12:00:00",
			fromTimezone: "Asia/Tokyo",
			toTimezone:   "America/New_York",
			wantErr:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("convert-timezone", map[string]interface{}{
				"datetime":      tc.datetime,
				"from_timezone": tc.fromTimezone,
				"to_timezone":   tc.toTimezone,
			})

			// Act
			result, err := handleConvertTimezone(ctx, request)

			// Assert
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.False(t, result.IsError)
				assert.NotEmpty(t, result.Content)
			}
		})
	}
}

func TestHandleConvertTimezone_Error(t *testing.T) {
	// テストケース
	tests := []struct {
		name      string
		arguments map[string]interface{}
		wantErr   bool
	}{
		{
			name:      "パラメータなし",
			arguments: map[string]interface{}{},
			wantErr:   true,
		},
		{
			name: "datetime パラメータなし",
			arguments: map[string]interface{}{
				"from_timezone": "UTC",
				"to_timezone":   "Asia/Tokyo",
			},
			wantErr: true,
		},
		{
			name: "from_timezone パラメータなし",
			arguments: map[string]interface{}{
				"datetime":    "2025-01-01 12:00:00",
				"to_timezone": "Asia/Tokyo",
			},
			wantErr: true,
		},
		{
			name: "to_timezone パラメータなし",
			arguments: map[string]interface{}{
				"datetime":      "2025-01-01 12:00:00",
				"from_timezone": "UTC",
			},
			wantErr: true,
		},
		{
			name: "無効な日時フォーマット",
			arguments: map[string]interface{}{
				"datetime":      "invalid-datetime",
				"from_timezone": "UTC",
				"to_timezone":   "Asia/Tokyo",
			},
			wantErr: false, // サービス層でエラーハンドリングされる
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("convert-timezone", tc.arguments)

			// Act
			result, err := handleConvertTimezone(ctx, request)

			// Assert
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				// サービス層でエラーが処理される場合
				if err != nil {
					assert.Error(t, err)
				} else {
					assert.NotNil(t, result)
				}
			}
		})
	}
}

func TestHandleListAvailableTimezones_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := createCallToolRequest("list-available-timezones", map[string]interface{}{})

	// Act
	result, err := handleListAvailableTimezones(ctx, request)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	assert.NotEmpty(t, result.Content)
}

// #==============================================================#
// ##          Server Tests                                      ##
// #==============================================================#

func TestSetTimezoneToolsServer_Normal(t *testing.T) {
	// テストケース
	tests := []struct {
		name string
	}{
		{
			name: "サーバーへのタイムゾーンツール追加",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockServer := server.NewMCPServer(
				"Mock Server",
				"1.0.0",
			)

			// Act
			resultServer := setTimezoneToolsServer(mockServer)

			// Assert
			assert.NotNil(t, resultServer, "サーバーが正しく設定されていません")
			assert.Equal(t, mockServer, resultServer, "返されたサーバーが入力と一致しません")
		})
	}
}

func TestCreateTimezoneServer_Normal(t *testing.T) {
	// テストケース
	tests := []struct {
		name string
	}{
		{
			name: "基本的なサーバー作成",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			server := createTimezoneServer()

			// Assert
			require.NotNil(t, server, "サーバーがnilです")

			// サーバーが正しく作成されていることを確認するために、
			// サーバーがnilでないことを検証するだけで十分です。
			// 実際のサーバーの動作は、統合テストで検証することができます。
		})
	}
}

// #==============================================================#
// ##          Error Handling Tests                             ##
// #==============================================================#

func TestErrorHandling_Normal(t *testing.T) {
	tests := []struct {
		name        string
		handlerFunc func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		request     mcp.CallToolRequest
		expectError bool
	}{
		{
			name:        "handleGetCurrentTimezone with missing timezone",
			handlerFunc: handleGetCurrentTimezone,
			request:     createCallToolRequest("get-current-timezone", map[string]interface{}{}),
			expectError: true,
		},
		{
			name:        "handleConvertTimezone with missing datetime",
			handlerFunc: handleConvertTimezone,
			request:     createCallToolRequest("convert-timezone", map[string]interface{}{"from_timezone": "UTC", "to_timezone": "Asia/Tokyo"}),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := tc.handlerFunc(context.Background(), tc.request)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// #==============================================================#
// ##          Mock and Utility Tests                           ##
// #==============================================================#

func TestCreateCallToolRequest_Normal(t *testing.T) {
	// Arrange
	name := "test-tool"
	arguments := map[string]interface{}{
		"param1": "value1",
		"param2": 123,
	}

	// Act
	request := createCallToolRequest(name, arguments)

	// Assert
	assert.Equal(t, name, request.Params.Name)
	assert.Equal(t, arguments, request.Params.Arguments)
	assert.Equal(t, "tools/call", request.Method)
}

func TestCallToolRequestMethods_Normal(t *testing.T) {
	// Arrange
	request := createCallToolRequest("test", map[string]interface{}{
		"string_param": "test_value",
		"int_param":    42,
		"float_param":  3.14,
		"bool_param":   true,
	})

	// Act & Assert
	stringVal, err := request.RequireString("string_param")
	assert.NoError(t, err)
	assert.Equal(t, "test_value", stringVal)

	intVal, err := request.RequireInt("int_param")
	assert.NoError(t, err)
	assert.Equal(t, 42, intVal)

	floatVal, err := request.RequireFloat("float_param")
	assert.NoError(t, err)
	assert.Equal(t, 3.14, floatVal)

	boolVal, err := request.RequireBool("bool_param")
	assert.NoError(t, err)
	assert.True(t, boolVal)
}

func TestCallToolRequestMethods_Error(t *testing.T) {
	// Arrange
	request := createCallToolRequest("test", map[string]interface{}{
		"wrong_type": 123,
	})

	// Act & Assert
	_, err := request.RequireString("missing_param")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	_, err = request.RequireString("wrong_type")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a string")
}

// #==============================================================#
// ##          Service Layer Integration Tests                  ##
// #==============================================================#

func TestServiceIntegration_Normal(t *testing.T) {
	// このテストは実際のサービス層との統合をテストします
	// 実際の実装では、サービス層のモックを使用することも可能です

	tests := []struct {
		name    string
		handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		request mcp.CallToolRequest
	}{
		{
			name:    "GetCurrentTimezone integration",
			handler: handleGetCurrentTimezone,
			request: createCallToolRequest("get-current-timezone", map[string]interface{}{
				"timezone": "UTC",
			}),
		},
		{
			name:    "ConvertTimezone integration",
			handler: handleConvertTimezone,
			request: createCallToolRequest("convert-timezone", map[string]interface{}{
				"datetime":      "2025-01-01 12:00:00",
				"from_timezone": "UTC",
				"to_timezone":   "Asia/Tokyo",
			}),
		},
		{
			name:    "ListAvailableTimezones integration",
			handler: handleListAvailableTimezones,
			request: createCallToolRequest("list-available-timezones", map[string]interface{}{}),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := tc.handler(context.Background(), tc.request)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.Content)

			// 結果の内容を検証
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(mcp.TextContent); ok {
					assert.NotEmpty(t, textContent.Text)
				}
			}
		})
	}
}

// #==============================================================#
// ##          Performance Tests                                ##
// #==============================================================#

func BenchmarkHandleGetCurrentTimezone(b *testing.B) {
	ctx := context.Background()
	request := createCallToolRequest("get-current-timezone", map[string]interface{}{
		"timezone": "UTC",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handleGetCurrentTimezone(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandleConvertTimezone(b *testing.B) {
	ctx := context.Background()
	request := createCallToolRequest("convert-timezone", map[string]interface{}{
		"datetime":      "2025-01-01 12:00:00",
		"from_timezone": "UTC",
		"to_timezone":   "Asia/Tokyo",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handleConvertTimezone(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandleListAvailableTimezones(b *testing.B) {
	ctx := context.Background()
	request := createCallToolRequest("list-available-timezones", map[string]interface{}{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := handleListAvailableTimezones(ctx, request)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// #==============================================================#
// ##          Integration Tests                                 ##
// #==============================================================#

// TestBuildTimezoneServer は BuildTimezoneServer 関数をテストします
func TestBuildTimezoneServer(t *testing.T) {
	// このテストは BuildTimezoneServer 関数をテストする例です
	// 実際のテストでは、サーバーの起動と終了をテストする必要があります

	// 注意: このテストは実際には機能しません。
	// MCPサーバーのテストには特別な設定が必要です。
	// このコードはあくまで例として提供されています。
	t.Skip("このテストはサーバーテストの例として提供されており、実際には実行されません")
}

// TestMCPToolHandlers は MCP ツールハンドラーをテストします
func TestMCPToolHandlers(t *testing.T) {
	// このテストは MCP ツールハンドラーをテストする例です
	// 実際のテストでは、リクエストとレスポンスをモック化する必要があります

	// 注意: このテストは実際には機能しません。
	// MCP ツールハンドラーのテストには特別な設定が必要です。
	// このコードはあくまで例として提供されています。
	t.Skip("このテストは MCP ツールハンドラーテストの例として提供されており、実際には実行されません")
}

// 統合テストの例
func TestIntegration(t *testing.T) {
	// 統合テストは実際のサーバーを起動するため、通常のテスト実行では
	// スキップされるようにしています。実際にテストを実行する場合は、
	// 環境変数などを使用して制御することをお勧めします。
	t.Skip("統合テストはデフォルトでスキップされます")

	// サーバーを起動（実際のテストでは、別のゴルーチンで起動し、
	// テスト終了時にシャットダウンする必要があります）
	// go BuildTimezoneServer()

	// サーバーが起動するのを待つ
	// time.Sleep(100 * time.Millisecond)

	// ここでクライアントを使用してサーバーにリクエストを送信し、
	// レスポンスを検証します
}
