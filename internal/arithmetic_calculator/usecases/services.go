package usecases

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	config "github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	reversePolishNotation "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/reverse_polish_notation"
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

// #==============================================================#
// ##          ApiCostExtractorService                           ##
// #==============================================================#
// ApiCostExtractorService はAPI料金抽出を行うサービスです
type ApiCostExtractorService struct {
	fileReader config.FileReader
}

// NewApiCostExtractorService は新しいApiCostExtractorServiceを作成します
func NewApiCostExtractorService() *ApiCostExtractorService {
	return &ApiCostExtractorService{
		fileReader: &config.StandardFileReader{},
	}
}

// NewApiCostExtractorServiceWithFileReader はFileReaderを注入した新しいApiCostExtractorServiceを作成します
func NewApiCostExtractorServiceWithFileReader(fileReader config.FileReader) *ApiCostExtractorService {
	return &ApiCostExtractorService{
		fileReader: fileReader,
	}
}

// extractApiCostFromText は文字列から「API料金が[数値]円掛かった」パターンを抽出し合計を計算する
func (s *ApiCostExtractorService) extractApiCostFromText(text string) (float64, error) {
	// 「API料金が[数値]円掛かった」を抽出する正規表現
	pattern := `API料金が(\d+)円掛かった`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(text, -1)

	var costs []float64
	for _, match := range matches {
		if len(match) > 1 {
			if cost, err := strconv.Atoi(match[1]); err == nil {
				costs = append(costs, float64(cost))
			}
		}
	}

	// 合計を計算
	total := 0.0
	for _, cost := range costs {
		total += cost
	}

	return total, nil
}

// HandleApiCostExtraction はファイルまたはテキストからAPI料金を抽出し合計を計算する
func (s *ApiCostExtractorService) HandleApiCostExtraction(filePath, textInput string) (float64, error) {
	// 排他制御
	if filePath != "" && textInput != "" {
		return 0, fmt.Errorf("ファイルパスとテキスト入力は同時に指定できません")
	}
	if filePath == "" && textInput == "" {
		return 0, fmt.Errorf("ファイルパスまたはテキスト入力のいずれかを指定してください")
	}

	var content string
	if filePath != "" {
		// ファイル拡張子の検証
		if !strings.HasSuffix(filePath, ".md") && !strings.HasSuffix(filePath, ".txt") {
			return 0, fmt.Errorf("ファイルは.mdまたは.txt形式である必要があります")
		}

		// ファイル読み込み（依存性注入されたFileReaderを使用）
		data, err := s.fileReader.ReadFile(filePath)
		if err != nil {
			return 0, fmt.Errorf("ファイル読み込みエラー: %v", err)
		}
		content = string(data)
	} else {
		content = textInput
	}

	return s.extractApiCostFromText(content)
}

// #==============================================================#
// ##          AdvancedMathService                               ##
// #==============================================================#
// AdvancedMathService は高度な数学演算を行うサービスです
type AdvancedMathService struct{}

// NewAdvancedMathService は新しいAdvancedMathServiceを作成します
func NewAdvancedMathService() *AdvancedMathService {
	return &AdvancedMathService{}
}

// power はべき乗計算を行うメソッドです
func (a *AdvancedMathService) power(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// squareRoot は平方根計算を行うメソッドです
func (a *AdvancedMathService) squareRoot(number float64) (float64, error) {
	if number < 0 {
		return 0, fmt.Errorf("負数の平方根は計算できません")
	}
	return math.Sqrt(number), nil
}

// factorial は階乗計算を行うメソッドです
func (a *AdvancedMathService) factorial(n int) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("負数の階乗は定義されていません")
	}
	if n > 170 {
		return 0, fmt.Errorf("数値が大きすぎて階乗計算でオーバーフローします")
	}

	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result, nil
}

// HandleToPower はべき乗計算のハンドラーです
func (a *AdvancedMathService) HandleToPower(base, exponent float64) (float64, error) {
	result := a.power(base, exponent)
	return result, nil
}

// HandleToSquareRoot は平方根計算のハンドラーです
func (a *AdvancedMathService) HandleToSquareRoot(number float64) (float64, error) {
	return a.squareRoot(number)
}

// HandleToFactorial は階乗計算のハンドラーです
func (a *AdvancedMathService) HandleToFactorial(n int) (float64, error) {
	return a.factorial(n)
}

// #==============================================================#
// ##          TrigonometryService                               ##
// #==============================================================#
// TrigonometryService は三角関数計算を行うサービスです
type TrigonometryService struct{}

// NewTrigonometryService は新しいTrigonometryServiceを作成します
func NewTrigonometryService() *TrigonometryService {
	return &TrigonometryService{}
}

// trigonometry は三角関数計算を行うメソッドです
func (t *TrigonometryService) trigonometry(function string, angle float64, unit string) (float64, error) {
	// 度数をラジアンに変換
	angleRad := angle
	if strings.ToLower(unit) == "degrees" {
		angleRad = angle * math.Pi / 180
	}

	switch strings.ToLower(function) {
	case "sin":
		return math.Sin(angleRad), nil
	case "cos":
		return math.Cos(angleRad), nil
	case "tan":
		return math.Tan(angleRad), nil
	default:
		return 0, fmt.Errorf("未知の三角関数です: %s", function)
	}
}

// HandleToTrigonometry は三角関数計算のハンドラーです
func (t *TrigonometryService) HandleToTrigonometry(function string, angle float64, unit string) (float64, error) {
	return t.trigonometry(function, angle, unit)
}

// #==============================================================#
// ##          MathConstantsService                              ##
// #==============================================================#
// MathConstantsService は数学定数を提供するサービスです
type MathConstantsService struct{}

// NewMathConstantsService は新しいMathConstantsServiceを作成します
func NewMathConstantsService() *MathConstantsService {
	return &MathConstantsService{}
}

// getConstants は利用可能な数学定数を返すメソッドです
func (m *MathConstantsService) getConstants() map[string]float64 {
	return map[string]float64{
		"pi":  math.Pi,
		"e":   math.E,
		"tau": 2 * math.Pi,
	}
}

// HandleToGetConstants は数学定数取得のハンドラーです
func (m *MathConstantsService) HandleToGetConstants() (string, error) {
	constants := m.getConstants()

	var result strings.Builder
	result.WriteString("利用可能な数学定数:\n")
	for name, value := range constants {
		result.WriteString(fmt.Sprintf("%s = %f\n", name, value))
	}

	return result.String(), nil
}

// #==============================================================#
// ##          ExpressionEvaluatorService                        ##
// #==============================================================#
// ExpressionEvaluatorService は安全な数式評価を行うサービスです
type ExpressionEvaluatorService struct {
	mathConstants *MathConstantsService
}

// NewExpressionEvaluatorService は新しいExpressionEvaluatorServiceを作成します
func NewExpressionEvaluatorService() *ExpressionEvaluatorService {
	return &ExpressionEvaluatorService{
		mathConstants: NewMathConstantsService(),
	}
}

// safeEvaluate は安全に数式を評価するメソッドです
func (e *ExpressionEvaluatorService) safeEvaluate(expression string) (float64, error) {
	// 空白を除去
	expression = strings.ReplaceAll(expression, " ", "")

	// 危険なパターンをチェック
	dangerousPatterns := []string{
		"__", "import", "exec", "eval", "open", "file", "input", "sys",
	}

	// 通常の危険パターンをチェック
	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(expression), pattern) {
			return 0, fmt.Errorf("危険なパターンが検出されました: %s", pattern)
		}
	}

	// "os"パターンの特別処理（cos関数内のosは許可）
	if err := e.checkOsPattern(expression); err != nil {
		return 0, err
	}

	// 数学定数を置換
	constants := e.mathConstants.getConstants()
	for name, value := range constants {
		expression = strings.ReplaceAll(expression, name, fmt.Sprintf("%f", value))
	}

	// RPNベースの評価器で式を安全に評価
	result, err := e.evaluateBasicExpression(expression)
	if err != nil {
		return 0, fmt.Errorf("数式の評価に失敗しました: %v", err)
	}

	return result, nil
}

// evaluateBasicExpression は基本的な数式を評価するメソッドです
func (e *ExpressionEvaluatorService) evaluateBasicExpression(expression string) (float64, error) {
	return e.evaluateUsingReversePolish(expression)
}

// evaluateArithmeticExpression は四則演算を評価するメソッドです
func (e *ExpressionEvaluatorService) evaluateArithmeticExpression(expression string) (float64, error) {
	return e.evaluateUsingReversePolish(expression)
}

func (e *ExpressionEvaluatorService) evaluateUsingReversePolish(expression string) (float64, error) {
	cleaned := strings.ReplaceAll(expression, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "**", "^")
	if cleaned == "" {
		return 0, fmt.Errorf("無効な数式です: %s", expression)
	}
	return reversePolishNotation.Evaluate(cleaned)
}

// checkOsPattern は"os"パターンの特別処理を行う（cos関数内のosは許可）
func (e *ExpressionEvaluatorService) checkOsPattern(expression string) error {
	lowerExpr := strings.ToLower(expression)

	// "os"の全ての出現位置を取得
	osIndices := e.getAllIndices(lowerExpr, "os")
	if len(osIndices) == 0 {
		return nil // "os"が含まれていない場合は問題なし
	}

	// "cos"の全ての出現位置を取得
	cosIndices := e.getAllIndices(lowerExpr, "cos")

	// 各"os"が"cos"の一部かチェック
	for _, osIndex := range osIndices {
		isPartOfCos := false
		for _, cosIndex := range cosIndices {
			if osIndex == cosIndex+1 {
				isPartOfCos = true
				break
			}
		}
		if !isPartOfCos {
			return fmt.Errorf("危険なパターンが検出されました: os")
		}
	}

	return nil
}

// getAllIndices は文字列内の指定されたパターンの全ての出現位置を返す
func (e *ExpressionEvaluatorService) getAllIndices(text, pattern string) []int {
	var indices []int
	start := 0

	for {
		index := strings.Index(text[start:], pattern)
		if index == -1 {
			break
		}
		actualIndex := start + index
		indices = append(indices, actualIndex)
		start = actualIndex + 1
	}

	return indices
}

// HandleToCalculateExpression は数式評価のハンドラーです
func (e *ExpressionEvaluatorService) HandleToCalculateExpression(expression string) (float64, error) {
	return e.safeEvaluate(expression)
}
