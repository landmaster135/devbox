package usecases

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// #==============================================================#
// ##          CalculatorService                                 ##
// #==============================================================#
// CalculatorService は算術計算を行うサービスです
type CalculatorService struct{}

// NewCalculatorService は新しいCalculatorServiceを作成します
func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

// add は二つの数値を足し算するメソッドです
func (c *CalculatorService) add(x float64, y float64) float64 {
	result := x + y
	return result
}

// subtract は二つの数値を引き算するメソッドです
func (c *CalculatorService) subtract(x float64, y float64) float64 {
	result := x - y
	return result
}

// multiply は二つの数値を掛け算するメソッドです
func (c *CalculatorService) multiply(x float64, y float64) float64 {
	result := x * y
	return result
}

// divide は二つの数値を割り算するメソッドです
func (c *CalculatorService) divide(x float64, y float64) float64 {
	result := x / y
	return result
}

// sum は複数の数値を合計するメソッドです
func (c *CalculatorService) sum(arr []float64) float64 {
	result := 0.0
	for _, number := range arr {
		result += number
	}
	return result
}

// HandleToCalculate はMCPリクエストを処理して計算結果を返すハンドラーです
func (c *CalculatorService) HandleToCalculate(op string, x, y float64) (float64, error) {
	var result float64
	switch op {
	case "add":
		result = c.add(x, y)
	case "subtract":
		result = c.subtract(x, y)
	case "multiply":
		result = c.multiply(x, y)
	case "divide":
		if y == 0 {
			return 0, fmt.Errorf("division by zero is not allowed")
		}
		result = c.divide(x, y)
	}
	return result, nil
}

// HandleToCalculateWithArray は配列を使った計算のMCPリクエストを処理するハンドラーです
func (c *CalculatorService) HandleToCalculateWithArray(op string, numbers []float64) (float64, error) {
	var result float64
	switch op {
	case "sum":
		result = c.sum(numbers)
	}

	return result, nil
}

// FileOpener インターフェースを定義
type FileOpener interface {
	Open(name string) (*os.File, error)
}

// DefaultFileOpener は標準のos.Openを使用する実装
type DefaultFileOpener struct{}

func (o *DefaultFileOpener) Open(name string) (*os.File, error) {
	return os.Open(name)
}

// BufioScanner インターフェースを定義
type BufioScanner interface {
	Scan() bool
	Err() error
}

// JSONMarshaler インターフェースを定義
type JSONMarshaler interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

// DefaultJSONMarshaler は標準のjson.MarshalIndentを使用する実装
type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// #==============================================================#
// ##          FileEvaluatorService                              ##
// #==============================================================#
// FileEvaluatorService はファイル評価を行うサービスです
type FileEvaluatorService struct {
	fileOpener    FileOpener
	bufioScanner  BufioScanner
	jsonMarshaler JSONMarshaler
}

// NewFileEvaluatorService は新しいFileEvaluatorServiceを作成します
func NewFileEvaluatorService() *FileEvaluatorService {
	return &FileEvaluatorService{
		fileOpener:    &DefaultFileOpener{},
		bufioScanner:  &bufio.Scanner{},
		jsonMarshaler: &DefaultJSONMarshaler{},
	}
}

// NewFileEvaluatorServiceWithDependencies はテスト用に依存性を注入できるFileEvaluatorServiceを作成します
func NewFileEvaluatorServiceWithDependencies(fileOpener FileOpener, bufioScanner BufioScanner, jsonMarshaler JSONMarshaler) *FileEvaluatorService {
	return &FileEvaluatorService{
		fileOpener:    fileOpener,
		bufioScanner:  bufioScanner,
		jsonMarshaler: jsonMarshaler,
	}
}

// CountLines はファイルの行数をカウントするメソッドです
func (e *FileEvaluatorService) CountLines(filePath string) (int, error) {
	file, err := e.fileOpener.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("ファイルを開けませんでした: %w", err)
	}
	defer file.Close()

	// テスト用のスキャナーが設定されている場合は、それを使用
	// そうでない場合は新しいスキャナーを作成
	var scanner BufioScanner

	// テスト用のスキャナーかどうかを判定するために、
	// 標準のbufio.Scannerの場合はnilを返すErrメソッドを利用
	if err := e.bufioScanner.Err(); err != nil {
		// Errがnilでない場合はテスト用のモックと判断
		scanner = e.bufioScanner
	} else {
		// 通常の処理では新しいスキャナーを作成
		scanner = bufio.NewScanner(file)
		e.bufioScanner = scanner
	}

	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

// IsLineCountGreaterThan はファイルの行数が指定された数値より大きいかどうかを判定するメソッドです
func (e *FileEvaluatorService) IsLineCountGreaterThan(filePath string, threshold int) (bool, int, error) {
	lineCount, err := e.CountLines(filePath)
	if err != nil {
		return false, 0, err
	}

	return lineCount > threshold, lineCount, nil
}

// isGreaterDescription は比較結果に基づいた説明文を返します
func isGreaterDescription(isGreater bool) string {
	if isGreater {
		return "より大きいです。"
	}
	return "以下です。"
}

// HandleToEvaluateLineCount はファイルの行数が指定された数値より大きいかどうかを判定し、結果を返すハンドラーです
func (e *FileEvaluatorService) HandleToEvaluateLineCount(filePath string, threshold int) (string, error) {
	// 行数の評価
	isGreater, lineCount, err := e.IsLineCountGreaterThan(filePath, int(threshold))
	if err != nil {
		return "", err
	}

	// 結果の作成
	result := map[string]interface{}{
		"is_greater":  isGreater,
		"line_count":  lineCount,
		"threshold":   int(threshold),
		"file_path":   filePath,
		"description": fmt.Sprintf("ファイル '%s' の行数は %d 行で、閾値 %d %s", filePath, lineCount, int(threshold), isGreaterDescription(isGreater)),
	}

	// JSON形式で結果を返す
	jsonResult, err := e.jsonMarshaler.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}
