package sequentialthinking

import (
	"context"
	"strings"
	"testing"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockSequentialThinkingService は SequentialThinkingService のモック実装です
type MockSequentialThinkingService struct {
	HandleSequentialThinkingFunc func(args map[string]interface{}) (string, error)
}

func (m *MockSequentialThinkingService) HandleSequentialThinking(args map[string]interface{}) (string, error) {
	if m.HandleSequentialThinkingFunc != nil {
		return m.HandleSequentialThinkingFunc(args)
	}
	return `{"thoughtNumber": 1, "totalThoughts": 3, "nextThoughtNeeded": true, "branches": [], "thoughtHistoryLength": 1}`, nil
}

// #==============================================================#
// ##          Test Classes                                      ##
// #==============================================================#

// TestSequentialThinkingMCP はSequentialThinkingMCPのテストクラスです
type TestSequentialThinkingMCP struct {
	server      *server.MCPServer
	mockService *MockSequentialThinkingService
}

// NewTestSequentialThinkingMCP はテスト用のSequentialThinkingMCPを作成します
func NewTestSequentialThinkingMCP() *TestSequentialThinkingMCP {
	return &TestSequentialThinkingMCP{
		mockService: &MockSequentialThinkingService{},
	}
}

// #==============================================================#
// ##          HandleSequentialThinking Tests                    ##
// #==============================================================#

// TestHandleSequentialThinking_Normal はhandleSequentialThinkingの正常系をテストします
func TestHandleSequentialThinking_Normal(t *testing.T) {
	// Arrange
	ctx := context.Background()
	args := map[string]interface{}{
		"thought":           "Test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: args,
		},
	}

	// Act
	result, err := handleSequentialThinking(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("handleSequentialThinking() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("handleSequentialThinking() returned nil result")
	}

	if result.Content == nil {
		t.Fatal("handleSequentialThinking() returned nil content")
	}

	// 結果がテキストコンテンツであることを確認
	if len(result.Content) == 0 {
		t.Error("handleSequentialThinking() returned empty content")
	}

	// JSONレスポンスの基本構造を確認
	content := result.Content[0]
	if textContent, ok := content.(mcp.TextContent); ok {
		if textContent.Text == "" {
			t.Error("handleSequentialThinking() returned empty text content")
		}
		// JSONの基本構造を確認
		if !strings.Contains(textContent.Text, "thoughtNumber") {
			t.Error("handleSequentialThinking() result should contain thoughtNumber")
		}
	} else {
		t.Errorf("handleSequentialThinking() should return TextContent, got %T", content)
	}
}

// TestHandleSequentialThinking_InvalidArgs は無効な引数のテストです
func TestHandleSequentialThinking_InvalidArgs(t *testing.T) {
	// Arrange
	ctx := context.Background()

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Empty thought",
			args: map[string]interface{}{
				"thought":           "",
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Missing thought",
			args: map[string]interface{}{
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Invalid thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(0),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Missing nextThoughtNeeded",
			args: map[string]interface{}{
				"thought":       "Test thought",
				"thoughtNumber": float64(1),
				"totalThoughts": float64(3),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			request := mcp.CallToolRequest{
				Request: mcp.Request{
					Method: "tools/call",
				},
				Params: mcp.CallToolParams{
					Name:      "sequentialthinking",
					Arguments: tt.args,
				},
			}

			// Act
			result, err := handleSequentialThinking(ctx, request)

			// Assert
			if err == nil {
				t.Error("handleSequentialThinking() should return error for invalid args")
			}

			if result != nil {
				t.Errorf("handleSequentialThinking() should return nil result on error, got: %v", result)
			}

			// エラーメッセージに適切な内容が含まれているか確認
			if !strings.Contains(err.Error(), "シーケンシャルシンキング処理に失敗しました") {
				t.Errorf("Error should contain expected message, got: %v", err.Error())
			}
		})
	}
}

// TestHandleSequentialThinking_WithOptionalFields はオプションフィールド付きのテストです
func TestHandleSequentialThinking_WithOptionalFields(t *testing.T) {
	// Arrange
	ctx := context.Background()
	args := map[string]interface{}{
		"thought":           "Revision thought",
		"thoughtNumber":     float64(2),
		"totalThoughts":     float64(5),
		"nextThoughtNeeded": false,
		"isRevision":        true,
		"revisesThought":    float64(1),
		"branchFromThought": float64(1),
		"branchId":          "branch-1",
		"needsMoreThoughts": true,
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: args,
		},
	}

	// Act
	result, err := handleSequentialThinking(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("handleSequentialThinking() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("handleSequentialThinking() returned nil result")
	}

	// 結果の内容を確認
	if len(result.Content) == 0 {
		t.Error("handleSequentialThinking() returned empty content")
	}
}

// #==============================================================#
// ##          AddPromptIntoServer Tests                         ##
// #==============================================================#

// TestAddPromptIntoServer_Normal はaddPromptIntoServerの正常系をテストします
func TestAddPromptIntoServer_Normal(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Server",
		"1.0.0",
		server.WithPromptCapabilities(true),
	)

	// Act
	result := addPromptIntoServer(s)

	// Assert
	if result == nil {
		t.Fatal("addPromptIntoServer() returned nil")
	}

	if result != s {
		t.Error("addPromptIntoServer() should return the same server instance")
	}

	// プロンプトが追加されているかは内部実装のため、エラーが発生しないことで確認
}

// TestAddPromptIntoServer_PromptHandler はプロンプトハンドラーの動作をテストします
func TestAddPromptIntoServer_PromptHandler(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Server",
		"1.0.0",
		server.WithPromptCapabilities(true),
	)
	s = addPromptIntoServer(s)

	// プロンプトハンドラーの直接テストは困難なため、
	// サーバーが正常に初期化されることで間接的に確認
	if s == nil {
		t.Error("Server should not be nil after adding prompt")
	}
}

// #==============================================================#
// ##          CreateSequentialThinkingServer Tests             ##
// #==============================================================#

// TestCreateSequentialThinkingServer_Normal はcreateSequentialThinkingServerの正常系をテストします
func TestCreateSequentialThinkingServer_Normal(t *testing.T) {
	// Act
	server := createSequentialThinkingServer()

	// Assert
	if server == nil {
		t.Fatal("createSequentialThinkingServer() returned nil")
	}

	// サーバーが正常に作成されていることを確認
	// 内部実装の詳細は確認できないが、エラーが発生しないことで確認
}

// TestCreateSequentialThinkingServer_ServerConfiguration はサーバー設定をテストします
func TestCreateSequentialThinkingServer_ServerConfiguration(t *testing.T) {
	// Act
	server := createSequentialThinkingServer()

	// Assert
	if server == nil {
		t.Fatal("createSequentialThinkingServer() returned nil")
	}

	// サーバーが正常に設定されていることを間接的に確認
	// 実際のMCPサーバーの内部状態は直接アクセスできないため、
	// 作成時にエラーが発生しないことで確認
}

// #==============================================================#
// ##          BuildSequentialThinkingServer Tests              ##
// #==============================================================#

// TestBuildSequentialThinkingServer_ServerCreation はBuildSequentialThinkingServerのサーバー作成をテストします
func TestBuildSequentialThinkingServer_ServerCreation(t *testing.T) {
	// BuildSequentialThinkingServer は server.ServeStdio を呼び出すため、
	// 直接テストすることは困難です。
	// 代わりに、createSequentialThinkingServer が正常に動作することで
	// 間接的にテストします。

	// Act
	server := createSequentialThinkingServer()

	// Assert
	if server == nil {
		t.Fatal("Server creation failed, BuildSequentialThinkingServer would fail")
	}
}

// #==============================================================#
// ##          Integration Tests                                 ##
// #==============================================================#

// TestSequentialThinkingMCP_Integration_ToolCall は統合テストです
func TestSequentialThinkingMCP_Integration_ToolCall(t *testing.T) {
	// Arrange
	ctx := context.Background()
	args := map[string]interface{}{
		"thought":           "Integration test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(2),
		"nextThoughtNeeded": true,
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: args,
		},
	}

	// Act
	result, err := handleSequentialThinking(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("Integration test failed with error: %v", err)
	}

	if result == nil {
		t.Fatal("Integration test returned nil result")
	}

	// 結果の構造を確認
	if len(result.Content) == 0 {
		t.Error("Integration test returned empty content")
	}

	// テキストコンテンツの確認
	if textContent, ok := result.Content[0].(mcp.TextContent); ok {
		if textContent.Text == "" {
			t.Error("Integration test returned empty text")
		}

		// JSONレスポンスの基本要素を確認
		expectedFields := []string{
			"thoughtNumber",
			"totalThoughts",
			"nextThoughtNeeded",
			"branches",
			"thoughtHistoryLength",
		}

		for _, field := range expectedFields {
			if !strings.Contains(textContent.Text, field) {
				t.Errorf("Integration test result should contain field: %s", field)
			}
		}
	} else {
		t.Errorf("Integration test should return TextContent, got %T", result.Content[0])
	}
}

// TestSequentialThinkingMCP_Integration_MultipleRequests は複数リクエストの統合テストです
func TestSequentialThinkingMCP_Integration_MultipleRequests(t *testing.T) {
	// Arrange
	ctx := context.Background()

	requests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "First thought",
			args: map[string]interface{}{
				"thought":           "First integration thought",
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Second thought",
			args: map[string]interface{}{
				"thought":           "Second integration thought",
				"thoughtNumber":     float64(2),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Final thought",
			args: map[string]interface{}{
				"thought":           "Final integration thought",
				"thoughtNumber":     float64(3),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": false,
			},
		},
	}

	// Act & Assert
	for _, req := range requests {
		t.Run(req.name, func(t *testing.T) {
			request := mcp.CallToolRequest{
				Request: mcp.Request{
					Method: "tools/call",
				},
				Params: mcp.CallToolParams{
					Name:      "sequentialthinking",
					Arguments: req.args,
				},
			}

			result, err := handleSequentialThinking(ctx, request)

			if err != nil {
				t.Errorf("Multiple requests test failed for %s: %v", req.name, err)
			}

			if result == nil {
				t.Fatalf("Multiple requests test returned nil result for %s", req.name)
			}

			if len(result.Content) == 0 {
				t.Errorf("Multiple requests test returned empty content for %s", req.name)
			}
		})
	}
}

// #==============================================================#
// ##          Error Handling Tests                             ##
// #==============================================================#

// TestHandleSequentialThinking_ServiceError はサービスエラーのテストです
func TestHandleSequentialThinking_ServiceError(t *testing.T) {
	// このテストは実際のサービスがエラーを返すケースをテストしますが、
	// 現在の実装では直接モックを注入できないため、
	// 無効な引数でサービスエラーを発生させます

	// Arrange
	ctx := context.Background()
	args := map[string]interface{}{
		"thought": nil, // nilを渡してサービスエラーを発生させる
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: args,
		},
	}

	// Act
	result, err := handleSequentialThinking(ctx, request)

	// Assert
	if err == nil {
		t.Error("handleSequentialThinking() should return error when service fails")
	}

	if result != nil {
		t.Errorf("handleSequentialThinking() should return nil result on service error, got: %v", result)
	}

	// エラーメッセージの確認
	if !strings.Contains(err.Error(), "シーケンシャルシンキング処理に失敗しました") {
		t.Errorf("Error message should contain expected text, got: %v", err.Error())
	}
}

// #==============================================================#
// ##          Edge Case Tests                                   ##
// #==============================================================#

// TestHandleSequentialThinking_EmptyArgs は空の引数のテストです
func TestHandleSequentialThinking_EmptyArgs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	args := map[string]interface{}{}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: args,
		},
	}

	// Act
	result, err := handleSequentialThinking(ctx, request)

	// Assert
	if err == nil {
		t.Error("handleSequentialThinking() should return error for empty args")
	}

	if result != nil {
		t.Errorf("handleSequentialThinking() should return nil result on error, got: %v", result)
	}
}

// TestHandleSequentialThinking_NilArgs はnil引数のテストです
func TestHandleSequentialThinking_NilArgs(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: nil,
		},
	}

	// Act
	result, err := handleSequentialThinking(ctx, request)

	// Assert
	if err == nil {
		t.Error("handleSequentialThinking() should return error for nil args")
	}

	if result != nil {
		t.Errorf("handleSequentialThinking() should return nil result on error, got: %v", result)
	}
}

// #==============================================================#
// ##          Boundary Value Tests                              ##
// #==============================================================#

// TestHandleSequentialThinking_BoundaryValues は境界値のテストです
func TestHandleSequentialThinking_BoundaryValues(t *testing.T) {
	// Arrange
	ctx := context.Background()

	tests := []struct {
		name      string
		args      map[string]interface{}
		wantError bool
	}{
		{
			name: "Minimum valid thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(1),
				"nextThoughtNeeded": false,
			},
			wantError: false,
		},
		{
			name: "Large thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(1000),
				"totalThoughts":     float64(1000),
				"nextThoughtNeeded": false,
			},
			wantError: false,
		},
		{
			name: "Very long thought text",
			args: map[string]interface{}{
				"thought":           strings.Repeat("This is a very long thought. ", 100),
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(1),
				"nextThoughtNeeded": false,
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			request := mcp.CallToolRequest{
				Request: mcp.Request{
					Method: "tools/call",
				},
				Params: mcp.CallToolParams{
					Name:      "sequentialthinking",
					Arguments: tt.args,
				},
			}

			// Act
			result, err := handleSequentialThinking(ctx, request)

			// Assert
			if (err != nil) != tt.wantError {
				t.Errorf("handleSequentialThinking() error = %v, wantError %v", err, tt.wantError)
			}

			if !tt.wantError && result == nil {
				t.Error("handleSequentialThinking() should return result for valid boundary values")
			}

			if tt.wantError && result != nil {
				t.Error("handleSequentialThinking() should return nil result on error")
			}
		})
	}
}

// #==============================================================#
// ##          Performance Tests                                 ##
// #==============================================================#

// TestHandleSequentialThinking_Performance は性能テストです
func TestHandleSequentialThinking_Performance(t *testing.T) {
	// Arrange
	ctx := context.Background()
	args := map[string]interface{}{
		"thought":           "Performance test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(1),
		"nextThoughtNeeded": false,
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
		},
		Params: mcp.CallToolParams{
			Name:      "sequentialthinking",
			Arguments: args,
		},
	}

	// Act - 複数回実行して性能を確認
	const iterations = 100
	for i := 0; i < iterations; i++ {
		result, err := handleSequentialThinking(ctx, request)

		// Assert
		if err != nil {
			t.Errorf("Performance test failed at iteration %d: %v", i, err)
			break
		}

		if result == nil {
			t.Errorf("Performance test returned nil result at iteration %d", i)
			break
		}
	}
}
