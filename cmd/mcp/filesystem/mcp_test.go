package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// getTextFromResult は結果からテキストを取得するヘルパー関数です
func getTextFromResult(result *mcp.CallToolResult) string {
	if result == nil || len(result.Content) == 0 {
		return ""
	}
	if textContent, ok := result.Content[0].(mcp.TextContent); ok {
		return textContent.Text
	}
	return ""
}

// #==============================================================#
// ##          Test Cases for Handler Functions                 ##
// #==============================================================#

// TestHandleReadFile_Normal は正常系のテストです
func TestHandleReadFile_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	ctx := context.Background()
	request := createCallToolRequest("read_file", map[string]interface{}{
		"path": testFile,
	})

	// Act
	result, err := handleReadFile(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, testContent) {
		t.Errorf("期待されたコンテンツが含まれていません。実際: %s", text)
	}
}

// TestHandleReadFile_MissingPath はパスパラメータが欠如している場合のテストです
func TestHandleReadFile_MissingPath(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := createCallToolRequest("read_file", map[string]interface{}{})

	// Act
	result, err := handleReadFile(ctx, request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、発生しませんでした")
	}
	if result != nil {
		t.Error("結果がnilでありません")
	}
}

// TestHandleReadFile_FileNotFound はファイルが存在しない場合のテストです
func TestHandleReadFile_FileNotFound(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := createCallToolRequest("read_file", map[string]interface{}{
		"path": "/nonexistent/file.txt",
	})

	// Act
	result, err := handleReadFile(ctx, request)

	// Assert
	if err == nil {
		t.Error("エラーが期待されましたが、発生しませんでした")
	}
	if result != nil {
		t.Error("結果がnilでありません")
	}
	if !strings.Contains(err.Error(), "ファイルの読み取りに失敗しました") {
		t.Errorf("期待されたエラーメッセージが含まれていません。実際: %s", err.Error())
	}
}

// TestHandleWriteFile_Normal は正常系のテストです
func TestHandleWriteFile_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	ctx := context.Background()
	request := createCallToolRequest("write_file", map[string]interface{}{
		"path":    testFile,
		"content": testContent,
	})

	// Act
	result, err := handleWriteFile(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "書き込みに成功しました") {
		t.Errorf("期待されたメッセージが含まれていません。実際: %s", text)
	}

	// ファイルが実際に作成されたことを確認
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("ファイルの読み取りに失敗しました: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("ファイル内容が期待値と異なります。期待値: %s, 実際: %s", testContent, string(content))
	}
}

// TestHandleWriteFile_MissingParameters はパラメータが欠如している場合のテストです
func TestHandleWriteFile_MissingParameters(t *testing.T) {
	testCases := []struct {
		name   string
		params map[string]interface{}
	}{
		{
			name:   "パスが欠如",
			params: map[string]interface{}{"content": "test"},
		},
		{
			name:   "コンテンツが欠如",
			params: map[string]interface{}{"path": "/tmp/test.txt"},
		},
		{
			name:   "両方が欠如",
			params: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("write_file", tc.params)

			// Act
			result, err := handleWriteFile(ctx, request)

			// Assert
			if err == nil {
				t.Error("エラーが期待されましたが、発生しませんでした")
			}
			if result != nil {
				t.Error("結果がnilでありません")
			}
		})
	}
}

// TestHandleCreateDirectory_Normal は正常系のテストです
func TestHandleCreateDirectory_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "newdir")

	ctx := context.Background()
	request := createCallToolRequest("create_directory", map[string]interface{}{
		"path": testDir,
	})

	// Act
	result, err := handleCreateDirectory(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "作成に成功しました") {
		t.Errorf("期待されたメッセージが含まれていません。実際: %s", text)
	}

	// ディレクトリが実際に作成されたことを確認
	info, err := os.Stat(testDir)
	if err != nil {
		t.Errorf("ディレクトリの確認に失敗しました: %v", err)
	}
	if !info.IsDir() {
		t.Error("ディレクトリが作成されていません")
	}
}

// TestHandleListDirectory_Normal は正常系のテストです
func TestHandleListDirectory_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テストファイルとディレクトリを作成
	testFile := filepath.Join(tempDir, "test.txt")
	testSubDir := filepath.Join(tempDir, "subdir")

	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	err = os.Mkdir(testSubDir, 0755)
	if err != nil {
		t.Fatalf("テストディレクトリの作成に失敗しました: %v", err)
	}

	ctx := context.Background()
	request := createCallToolRequest("list_directory", map[string]interface{}{
		"path": tempDir,
	})

	// Act
	result, err := handleListDirectory(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "[FILE] test.txt") {
		t.Error("ファイルが正しく表示されていません")
	}
	if !strings.Contains(text, "[DIR] subdir") {
		t.Error("ディレクトリが正しく表示されていません")
	}
}

// TestHandleDirectoryTree_Normal は正常系のテストです
func TestHandleDirectoryTree_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テストファイルを作成
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	ctx := context.Background()
	request := createCallToolRequest("directory_tree", map[string]interface{}{
		"path": tempDir,
	})

	// Act
	result, err := handleDirectoryTree(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "test.txt") {
		t.Error("ファイルがJSONに含まれていません")
	}
	if !strings.Contains(text, `"type": "file"`) {
		t.Error("ファイルタイプが正しく設定されていません")
	}
}

// TestHandleMoveFile_Normal は正常系のテストです
func TestHandleMoveFile_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	sourceFile := filepath.Join(tempDir, "source.txt")
	destFile := filepath.Join(tempDir, "dest.txt")
	testContent := "test content"

	err := os.WriteFile(sourceFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	ctx := context.Background()
	request := createCallToolRequest("move_file", map[string]interface{}{
		"source":      sourceFile,
		"destination": destFile,
	})

	// Act
	result, err := handleMoveFile(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "移動に成功しました") {
		t.Errorf("期待されたメッセージが含まれていません。実際: %s", text)
	}

	// ファイルが移動されたことを確認
	if _, err := os.Stat(sourceFile); !os.IsNotExist(err) {
		t.Error("ソースファイルが削除されていません")
	}

	content, err := os.ReadFile(destFile)
	if err != nil {
		t.Errorf("移動先ファイルの読み取りに失敗しました: %v", err)
	}
	if string(content) != testContent {
		t.Errorf("ファイル内容が期待値と異なります。期待値: %s, 実際: %s", testContent, string(content))
	}
}

// TestHandleSearchFiles_Normal は正常系のテストです
func TestHandleSearchFiles_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テストファイルを作成
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.go")
	testFile3 := filepath.Join(tempDir, "other.txt")

	for _, file := range []string{testFile1, testFile2, testFile3} {
		err := os.WriteFile(file, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	ctx := context.Background()
	request := createCallToolRequest("search_files", map[string]interface{}{
		"path":    tempDir,
		"pattern": "test",
	})

	// Act
	result, err := handleSearchFiles(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "test1.txt") {
		t.Error("test1.txtが検索結果に含まれていません")
	}
	if !strings.Contains(text, "test2.go") {
		t.Error("test2.goが検索結果に含まれていません")
	}
	if strings.Contains(text, "other.txt") {
		t.Error("other.txtが検索結果に含まれています（含まれるべきではありません）")
	}
}

// TestHandleSearchFiles_WithExcludePattern は除外パターンのテストです
func TestHandleSearchFiles_WithExcludePattern(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// テストファイルを作成
	testFile1 := filepath.Join(tempDir, "test1.txt")
	testFile2 := filepath.Join(tempDir, "test2.go")

	for _, file := range []string{testFile1, testFile2} {
		err := os.WriteFile(file, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("テストファイルの作成に失敗しました: %v", err)
		}
	}

	ctx := context.Background()
	request := createCallToolRequest("search_files", map[string]interface{}{
		"path":            tempDir,
		"pattern":         "test",
		"exclude_pattern": "*.go",
	})

	// Act
	result, err := handleSearchFiles(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "test1.txt") {
		t.Error("test1.txtが検索結果に含まれていません")
	}
	if strings.Contains(text, "test2.go") {
		t.Error("test2.goが検索結果に含まれています（除外されるべきです）")
	}
}

// TestHandleGetFileInfo_Normal は正常系のテストです
func TestHandleGetFileInfo_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	testContent := "test content"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	ctx := context.Background()
	request := createCallToolRequest("get_file_info", map[string]interface{}{
		"path": testFile,
	})

	// Act
	result, err := handleGetFileInfo(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "サイズ:") {
		t.Error("ファイルサイズ情報が含まれていません")
	}
	if !strings.Contains(text, "ファイル: true") {
		t.Error("ファイルタイプ情報が正しくありません")
	}
}

// TestHandleListAllowedDirectories_Normal は正常系のテストです
func TestHandleListAllowedDirectories_Normal(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	ctx := context.Background()
	request := createCallToolRequest("list_allowed_directories", map[string]interface{}{
		"path": tempDir,
	})

	// Act
	result, err := handleListAllowedDirectories(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "許可されたディレクトリ:") {
		t.Error("期待されたメッセージが含まれていません")
	}
}

// TestAddPromptIntoServer_Normal は正常系のテストです
func TestAddPromptIntoServer_Normal(t *testing.T) {
	// Arrange
	s := server.NewMCPServer(
		"Test Server",
		"1.0.0",
	)

	// Act
	result := addPromptIntoServer(s)

	// Assert
	if result == nil {
		t.Fatal("結果がnilです")
	}
	// サーバーが返されることを確認
	if result != s {
		t.Error("同じサーバーインスタンスが返されていません")
	}
}

// BuildFileSystemServer関数のテストはスキップ
// 実際のサーバーを起動するため、単体テストには適さない
func TestBuildFileSystemServer(t *testing.T) {
	t.Skip("このテストは実際にサーバーを起動するため、スキップします")
}

// #==============================================================#
// ##          Error Handling Tests                             ##
// #==============================================================#

// TestHandlers_ErrorHandling は各ハンドラーのエラーハンドリングをテストします
func TestHandlers_ErrorHandling(t *testing.T) {
	testCases := []struct {
		name    string
		handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		params  map[string]interface{}
	}{
		{
			name:    "handleReadFile - 無効なパス",
			handler: handleReadFile,
			params:  map[string]interface{}{"path": "/invalid/path/file.txt"},
		},
		{
			name:    "handleWriteFile - 無効なパス",
			handler: handleWriteFile,
			params:  map[string]interface{}{"path": "/invalid/path/file.txt", "content": "test"},
		},
		{
			name:    "handleCreateDirectory - 無効なパス",
			handler: handleCreateDirectory,
			params:  map[string]interface{}{"path": "/invalid/path/dir"},
		},
		{
			name:    "handleListDirectory - 無効なパス",
			handler: handleListDirectory,
			params:  map[string]interface{}{"path": "/invalid/path"},
		},
		{
			name:    "handleDirectoryTree - 無効なパス",
			handler: handleDirectoryTree,
			params:  map[string]interface{}{"path": "/invalid/path"},
		},
		{
			name:    "handleMoveFile - 無効なソース",
			handler: handleMoveFile,
			params:  map[string]interface{}{"source": "/invalid/source", "destination": "/tmp/dest"},
		},
		{
			name:    "handleGetFileInfo - 無効なパス",
			handler: handleGetFileInfo,
			params:  map[string]interface{}{"path": "/invalid/path/file.txt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("test_tool", tc.params)

			// Act
			result, err := tc.handler(ctx, request)

			// Assert
			if err == nil {
				t.Error("エラーが期待されましたが、発生しませんでした")
			}
			if result != nil {
				t.Error("結果がnilでありません")
			}
		})
	}
}

// TestHandleSearchFiles_ErrorHandling は検索ファイルのエラーハンドリングを個別にテストします
func TestHandleSearchFiles_ErrorHandling(t *testing.T) {
	// Arrange
	ctx := context.Background()
	request := createCallToolRequest("search_files", map[string]interface{}{
		"path": "/invalid/path",
		"pattern": "test",
	})

	// Act
	result, err := handleSearchFiles(ctx, request)

	// Assert
	// SearchFilesは無効なパスでもエラーを返さず、空の結果を返す可能性がある
	if err != nil {
		// エラーが返された場合
		if result != nil {
			t.Error("エラー時に結果がnilでありません")
		}
	} else {
		// エラーが返されなかった場合、結果が空であることを確認
		if result == nil {
			t.Error("結果がnilです")
		} else {
			text := getTextFromResult(result)
			if !strings.Contains(text, "一致するものが見つかりませんでした") && text != "" {
				t.Errorf("予期しない結果が返されました: %s", text)
			}
		}
	}
}

// #==============================================================#
// ##          Edge Case Tests                                  ##
// #==============================================================#

// TestHandleSearchFiles_NoResults は検索結果が0件の場合のテストです
func TestHandleSearchFiles_NoResults(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// 検索パターンに一致しないファイルを作成
	testFile := filepath.Join(tempDir, "nomatch.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("テストファイルの作成に失敗しました: %v", err)
	}

	ctx := context.Background()
	request := createCallToolRequest("search_files", map[string]interface{}{
		"path":    tempDir,
		"pattern": "nonexistent",
	})

	// Act
	result, err := handleSearchFiles(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	if !strings.Contains(text, "一致するものが見つかりませんでした") {
		t.Error("期待されたメッセージが含まれていません")
	}
}

// TestHandleDirectoryTree_EmptyDirectory は空のディレクトリのテストです
func TestHandleDirectoryTree_EmptyDirectory(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	ctx := context.Background()
	request := createCallToolRequest("directory_tree", map[string]interface{}{
		"path": tempDir,
	})

	// Act
	result, err := handleDirectoryTree(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	text := getTextFromResult(result)
	// 空のディレクトリの場合、空の配列が返されることを確認
	// JSONフォーマットで空の配列、または空白のみの場合も許可
	if !strings.Contains(text, "[]") && !strings.Contains(text, "[ ]") && !strings.Contains(text, "null") && strings.TrimSpace(text) != "" {
		t.Errorf("空のディレクトリに対して期待される結果が返されていません。実際: '%s'", text)
	}
}

// TestHandleCreateDirectory_NestedPath はネストされたパスのテストです
func TestHandleCreateDirectory_NestedPath(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	nestedDir := filepath.Join(tempDir, "level1", "level2", "level3")

	ctx := context.Background()
	request := createCallToolRequest("create_directory", map[string]interface{}{
		"path": nestedDir,
	})

	// Act
	result, err := handleCreateDirectory(ctx, request)

	// Assert
	if err != nil {
		t.Errorf("エラーが発生しました: %v", err)
	}
	if result == nil {
		t.Fatal("結果がnilです")
	}

	// ネストされたディレクトリが作成されたことを確認
	info, err := os.Stat(nestedDir)
	if err != nil {
		t.Errorf("ネストされたディレクトリの確認に失敗しました: %v", err)
	}
	if !info.IsDir() {
		t.Error("ネストされたディレクトリが作成されていません")
	}
}

// #==============================================================#
// ##          Parameter Validation Tests                       ##
// #==============================================================#

// TestParameterValidation は各ハンドラーのパラメータ検証をテストします
func TestParameterValidation(t *testing.T) {
	testCases := []struct {
		name    string
		handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "handleReadFile - パラメータなし",
			handler: handleReadFile,
			params:  map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "handleWriteFile - contentパラメータなし",
			handler: handleWriteFile,
			params:  map[string]interface{}{"path": "/tmp/test.txt"},
			wantErr: true,
		},
		{
			name:    "handleMoveFile - sourceパラメータなし",
			handler: handleMoveFile,
			params:  map[string]interface{}{"destination": "/tmp/dest.txt"},
			wantErr: true,
		},
		{
			name:    "handleMoveFile - destinationパラメータなし",
			handler: handleMoveFile,
			params:  map[string]interface{}{"source": "/tmp/source.txt"},
			wantErr: true,
		},
		{
			name:    "handleSearchFiles - patternパラメータなし",
			handler: handleSearchFiles,
			params:  map[string]interface{}{"path": "/tmp"},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("test_tool", tc.params)

			// Act
			result, err := tc.handler(ctx, request)

			// Assert
			if tc.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				if result != nil {
					t.Error("エラー時に結果がnilでありません")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("結果がnilです")
				}
			}
		})
	}
}

// TestParameterValidation_WithValidPaths は有効なパスでのパラメータ検証をテストします
func TestParameterValidation_WithValidPaths(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	testCases := []struct {
		name    string
		handler func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error)
		params  map[string]interface{}
		wantErr bool
	}{
		{
			name:    "handleReadFile - 存在しないファイル",
			handler: handleReadFile,
			params:  map[string]interface{}{"path": filepath.Join(tempDir, "nonexistent.txt")},
			wantErr: true, // ファイルが存在しないためエラーになる
		},
		{
			name:    "handleWriteFile - 有効なパラメータ",
			handler: handleWriteFile,
			params:  map[string]interface{}{"path": filepath.Join(tempDir, "test.txt"), "content": "test"},
			wantErr: false, // 有効なパスなので成功する
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctx := context.Background()
			request := createCallToolRequest("test_tool", tc.params)

			// Act
			result, err := tc.handler(ctx, request)

			// Assert
			if tc.wantErr {
				if err == nil {
					t.Error("エラーが期待されましたが、発生しませんでした")
				}
				if result != nil {
					t.Error("エラー時に結果がnilでありません")
				}
			} else {
				if err != nil {
					t.Errorf("予期しないエラーが発生しました: %v", err)
				}
				if result == nil {
					t.Error("結果がnilです")
				}
			}
		})
	}
}
