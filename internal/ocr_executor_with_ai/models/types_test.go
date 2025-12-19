package models

import (
	"strings"
	"testing"
)

// TestOcrResult はOcrResult構造体のテストクラス
type TestOcrResult struct{}

func (t *TestOcrResult) TestOcrResult_Creation_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/test.jpg"
		testContent  = "テスト用のOCR結果"
		testError    = ""
	)

	result := OcrResult{
		FilePath: testFilePath,
		Content:  testContent,
		Error:    testError,
	}

	if result.FilePath != testFilePath {
		test.Errorf("FilePath = %v, want %v", result.FilePath, testFilePath)
	}
	if result.Content != testContent {
		test.Errorf("Content = %v, want %v", result.Content, testContent)
	}
	if result.Error != testError {
		test.Errorf("Error = %v, want %v", result.Error, testError)
	}
}

func (t *TestOcrResult) TestOcrResult_WithError_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/error.jpg"
		testContent  = ""
		testError    = "ファイルが見つかりません"
	)

	result := OcrResult{
		FilePath: testFilePath,
		Content:  testContent,
		Error:    testError,
	}

	if result.FilePath != testFilePath {
		test.Errorf("FilePath = %v, want %v", result.FilePath, testFilePath)
	}
	if result.Content != testContent {
		test.Errorf("Content = %v, want %v", result.Content, testContent)
	}
	if result.Error != testError {
		test.Errorf("Error = %v, want %v", result.Error, testError)
	}
}

func (t *TestOcrResult) TestOcrResult_WithoutError_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/success.jpg"
		testContent  = "成功したOCR結果"
		testError    = ""
	)

	result := OcrResult{
		FilePath: testFilePath,
		Content:  testContent,
		Error:    testError,
	}

	if result.FilePath != testFilePath {
		test.Errorf("FilePath = %v, want %v", result.FilePath, testFilePath)
	}
	if result.Content != testContent {
		test.Errorf("Content = %v, want %v", result.Content, testContent)
	}
	if result.Error != testError {
		test.Errorf("Error = %v, want %v", result.Error, testError)
	}
}

// TestOcrExecutionResult はOcrExecutionResult構造体のテストクラス
type TestOcrExecutionResult struct{}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_Creation_Normal(test *testing.T) {
	result := OcrExecutionResult{
		Results: []OcrResult{},
		Total:   0,
		Success: 0,
		Failed:  0,
	}

	if len(result.Results) != 0 {
		test.Errorf("Results length = %v, want %v", len(result.Results), 0)
	}
	if result.Total != 0 {
		test.Errorf("Total = %v, want %v", result.Total, 0)
	}
	if result.Success != 0 {
		test.Errorf("Success = %v, want %v", result.Success, 0)
	}
	if result.Failed != 0 {
		test.Errorf("Failed = %v, want %v", result.Failed, 0)
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_AddResult_Success_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/success.jpg"
		testContent  = "成功したOCR結果"
	)

	executionResult := &OcrExecutionResult{}
	successResult := OcrResult{
		FilePath: testFilePath,
		Content:  testContent,
		Error:    "",
	}

	executionResult.AddResult(successResult)

	if len(executionResult.Results) != 1 {
		test.Errorf("Results length = %v, want %v", len(executionResult.Results), 1)
	}
	if executionResult.Total != 1 {
		test.Errorf("Total = %v, want %v", executionResult.Total, 1)
	}
	if executionResult.Success != 1 {
		test.Errorf("Success = %v, want %v", executionResult.Success, 1)
	}
	if executionResult.Failed != 0 {
		test.Errorf("Failed = %v, want %v", executionResult.Failed, 0)
	}
	if executionResult.Results[0].FilePath != testFilePath {
		test.Errorf("Results[0].FilePath = %v, want %v", executionResult.Results[0].FilePath, testFilePath)
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_AddResult_Failed_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/failed.jpg"
		testError    = "処理に失敗しました"
	)

	executionResult := &OcrExecutionResult{}
	failedResult := OcrResult{
		FilePath: testFilePath,
		Content:  "",
		Error:    testError,
	}

	executionResult.AddResult(failedResult)

	if len(executionResult.Results) != 1 {
		test.Errorf("Results length = %v, want %v", len(executionResult.Results), 1)
	}
	if executionResult.Total != 1 {
		test.Errorf("Total = %v, want %v", executionResult.Total, 1)
	}
	if executionResult.Success != 0 {
		test.Errorf("Success = %v, want %v", executionResult.Success, 0)
	}
	if executionResult.Failed != 1 {
		test.Errorf("Failed = %v, want %v", executionResult.Failed, 1)
	}
	if executionResult.Results[0].Error != testError {
		test.Errorf("Results[0].Error = %v, want %v", executionResult.Results[0].Error, testError)
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_AddResult_Multiple_Normal(test *testing.T) {
	executionResult := &OcrExecutionResult{}

	// 成功結果を追加
	successResult := OcrResult{
		FilePath: "/path/to/success.jpg",
		Content:  "成功したOCR結果",
		Error:    "",
	}
	executionResult.AddResult(successResult)

	// 失敗結果を追加
	failedResult := OcrResult{
		FilePath: "/path/to/failed.jpg",
		Content:  "",
		Error:    "処理に失敗しました",
	}
	executionResult.AddResult(failedResult)

	// もう一つ成功結果を追加
	anotherSuccessResult := OcrResult{
		FilePath: "/path/to/another_success.jpg",
		Content:  "別の成功したOCR結果",
		Error:    "",
	}
	executionResult.AddResult(anotherSuccessResult)

	if len(executionResult.Results) != 3 {
		test.Errorf("Results length = %v, want %v", len(executionResult.Results), 3)
	}
	if executionResult.Total != 3 {
		test.Errorf("Total = %v, want %v", executionResult.Total, 3)
	}
	if executionResult.Success != 2 {
		test.Errorf("Success = %v, want %v", executionResult.Success, 2)
	}
	if executionResult.Failed != 1 {
		test.Errorf("Failed = %v, want %v", executionResult.Failed, 1)
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_FormatAsText_Empty_Normal(test *testing.T) {
	const expectedOutput = "処理対象の画像ファイルが見つかりませんでした。"

	executionResult := &OcrExecutionResult{}
	output := executionResult.FormatAsText()

	if output != expectedOutput {
		test.Errorf("FormatAsText() = %v, want %v", output, expectedOutput)
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_FormatAsText_SuccessOnly_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/success.jpg"
		testContent  = "成功したOCR結果"
	)

	executionResult := &OcrExecutionResult{}
	successResult := OcrResult{
		FilePath: testFilePath,
		Content:  testContent,
		Error:    "",
	}
	executionResult.AddResult(successResult)

	output := executionResult.FormatAsText()

	// 期待される出力の検証
	if !strings.Contains(output, "=== AI OCR実行結果 ===") {
		test.Error("出力にヘッダーが含まれていません")
	}
	if !strings.Contains(output, "処理総数: 1件 (成功: 1件, 失敗: 0件)") {
		test.Error("出力に正しい統計情報が含まれていません")
	}
	if !strings.Contains(output, testFilePath) {
		test.Error("出力にファイルパスが含まれていません")
	}
	if !strings.Contains(output, testContent) {
		test.Error("出力にOCR結果が含まれていません")
	}
	if strings.Contains(output, "エラー:") {
		test.Error("出力にエラー情報が含まれています")
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_FormatAsText_FailedOnly_Normal(test *testing.T) {
	const (
		testFilePath = "/path/to/failed.jpg"
		testError    = "処理に失敗しました"
	)

	executionResult := &OcrExecutionResult{}
	failedResult := OcrResult{
		FilePath: testFilePath,
		Content:  "",
		Error:    testError,
	}
	executionResult.AddResult(failedResult)

	output := executionResult.FormatAsText()

	// 期待される出力の検証
	if !strings.Contains(output, "=== AI OCR実行結果 ===") {
		test.Error("出力にヘッダーが含まれていません")
	}
	if !strings.Contains(output, "処理総数: 1件 (成功: 0件, 失敗: 1件)") {
		test.Error("出力に正しい統計情報が含まれていません")
	}
	if !strings.Contains(output, testFilePath) {
		test.Error("出力にファイルパスが含まれていません")
	}
	if !strings.Contains(output, "エラー: "+testError) {
		test.Error("出力にエラー情報が含まれていません")
	}
}

func (t *TestOcrExecutionResult) TestOcrExecutionResult_FormatAsText_Mixed_Normal(test *testing.T) {
	executionResult := &OcrExecutionResult{}

	// 成功結果を追加
	successResult := OcrResult{
		FilePath: "/path/to/success.jpg",
		Content:  "成功したOCR結果",
		Error:    "",
	}
	executionResult.AddResult(successResult)

	// 失敗結果を追加
	failedResult := OcrResult{
		FilePath: "/path/to/failed.jpg",
		Content:  "",
		Error:    "処理に失敗しました",
	}
	executionResult.AddResult(failedResult)

	output := executionResult.FormatAsText()

	// 期待される出力の検証
	if !strings.Contains(output, "=== AI OCR実行結果 ===") {
		test.Error("出力にヘッダーが含まれていません")
	}
	if !strings.Contains(output, "処理総数: 2件 (成功: 1件, 失敗: 1件)") {
		test.Error("出力に正しい統計情報が含まれていません")
	}
	if !strings.Contains(output, "/path/to/success.jpg") {
		test.Error("出力に成功ファイルのパスが含まれていません")
	}
	if !strings.Contains(output, "/path/to/failed.jpg") {
		test.Error("出力に失敗ファイルのパスが含まれていません")
	}
	if !strings.Contains(output, "成功したOCR結果") {
		test.Error("出力に成功したOCR結果が含まれていません")
	}
	if !strings.Contains(output, "エラー: 処理に失敗しました") {
		test.Error("出力にエラー情報が含まれていません")
	}
}

// テスト実行用の関数
func TestOcrResult_Creation_Normal(t *testing.T) {
	testInstance := &TestOcrResult{}
	testInstance.TestOcrResult_Creation_Normal(t)
}

func TestOcrResult_WithError_Normal(t *testing.T) {
	testInstance := &TestOcrResult{}
	testInstance.TestOcrResult_WithError_Normal(t)
}

func TestOcrResult_WithoutError_Normal(t *testing.T) {
	testInstance := &TestOcrResult{}
	testInstance.TestOcrResult_WithoutError_Normal(t)
}

func TestOcrExecutionResult_Creation_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_Creation_Normal(t)
}

func TestOcrExecutionResult_AddResult_Success_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_AddResult_Success_Normal(t)
}

func TestOcrExecutionResult_AddResult_Failed_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_AddResult_Failed_Normal(t)
}

func TestOcrExecutionResult_AddResult_Multiple_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_AddResult_Multiple_Normal(t)
}

func TestOcrExecutionResult_FormatAsText_Empty_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_FormatAsText_Empty_Normal(t)
}

func TestOcrExecutionResult_FormatAsText_SuccessOnly_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_FormatAsText_SuccessOnly_Normal(t)
}

func TestOcrExecutionResult_FormatAsText_FailedOnly_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_FormatAsText_FailedOnly_Normal(t)
}

func TestOcrExecutionResult_FormatAsText_Mixed_Normal(t *testing.T) {
	testInstance := &TestOcrExecutionResult{}
	testInstance.TestOcrExecutionResult_FormatAsText_Mixed_Normal(t)
}
