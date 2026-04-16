package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	OperationAdd               = "add"
	OperationSubtract          = "subtract"
	OperationMultiply          = "multiply"
	OperationDivide            = "divide"
	OperationSum               = "sum"
	OperationEvaluateLineCount = "evaluate_line_count"
	OperationParseAPICost      = "parse-api-cost"
	OperationPower             = "power"
	OperationSquareRoot        = "square_root"
	OperationFactorial         = "factorial"
	OperationTrigonometry      = "trigonometry"
	OperationCalculate         = "calculate"
	OperationGetConstants      = "get_constants"
)

var supportedOperations = []string{
	OperationAdd,
	OperationSubtract,
	OperationMultiply,
	OperationDivide,
	OperationSum,
	OperationEvaluateLineCount,
	OperationParseAPICost,
	OperationPower,
	OperationSquareRoot,
	OperationFactorial,
	OperationTrigonometry,
	OperationCalculate,
	OperationGetConstants,
}

// SupportedOperations はサポート済みoperation一覧を返す
func SupportedOperations() []string {
	operations := make([]string, len(supportedOperations))
	copy(operations, supportedOperations)
	return operations
}

func isSupportedOperation(operation string) bool {
	for _, supported := range supportedOperations {
		if operation == supported {
			return true
		}
	}
	return false
}

// Config は算術計算CLIの設定を保持する構造体
type Config struct {
	Operation  string    // 操作タイプ (add, subtract, multiply, divide, sum, evaluate_line_count, parse-api-cost, power, square_root, factorial, trigonometry, calculate, get_constants)
	X          float64   // 第一オペランド
	Y          float64   // 第二オペランド
	Numbers    []float64 // 複数の数値（sum操作用）
	FilePath   string    // ファイルパス（evaluate_line_count, parse-api-cost操作用）
	Threshold  int       // 閾値（evaluate_line_count操作用）
	TextInput  string    // テキスト入力（parse-api-cost操作用）
	Base       float64   // べき乗の底（power操作用）
	Exponent   float64   // べき乗の指数（power操作用）
	Number     float64   // 単一の数値（square_root操作用）
	N          int       // 階乗の数値（factorial操作用）
	Function   string    // 三角関数の種類（trigonometry操作用）
	Angle      float64   // 角度（trigonometry操作用）
	Unit       string    // 角度の単位（trigonometry操作用）
	Expression string    // 数式（calculate操作用）
	Help       bool      // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation string, x, y float64, numbers []float64, filePath string, threshold int, base, exponent, number, angle float64, n int, function, unit, expression, textInput string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	if !isSupportedOperation(operation) {
		return nil, fmt.Errorf("無効な操作タイプです: %s", operation)
	}

	// 操作タイプ別の検証
	switch operation {
	case OperationAdd, OperationSubtract, OperationMultiply, OperationDivide:
		// 二項演算の場合、x, yが必要
		// デフォルト値0.0でも計算可能なので、特別な検証は不要
	case OperationSum:
		if len(numbers) == 0 {
			return nil, fmt.Errorf("sum操作には数値の配列が必要です")
		}
	case OperationEvaluateLineCount:
		if filePath == "" {
			return nil, fmt.Errorf("evaluate_line_count操作にはファイルパスが必要です")
		}
		if threshold < 0 {
			return nil, fmt.Errorf("閾値は0以上である必要があります")
		}
	case OperationPower:
		// べき乗計算では base と exponent が必要
	case OperationSquareRoot:
		// 平方根計算では number が必要
	case OperationFactorial:
		// 階乗計算では n が必要
		if n < 0 {
			return nil, fmt.Errorf("階乗は負数では定義されていません")
		}
	case OperationTrigonometry:
		// 三角関数計算では function, angle, unit が必要
		if function == "" {
			return nil, fmt.Errorf("三角関数の種類が指定されていません")
		}
		if unit == "" {
			unit = "radians" // デフォルトはラジアン
		}
	case OperationCalculate:
		// 数式評価では expression が必要
		if expression == "" {
			return nil, fmt.Errorf("評価する数式が指定されていません")
		}
	}

	return &Config{
		Operation:  operation,
		X:          x,
		Y:          y,
		Numbers:    numbers,
		FilePath:   filePath,
		Threshold:  threshold,
		TextInput:  textInput,
		Base:       base,
		Exponent:   exponent,
		Number:     number,
		N:          n,
		Function:   function,
		Angle:      angle,
		Unit:       unit,
		Expression: expression,
	}, nil
}

// NewConfigForParseApiCost はparse-api-cost操作用の新しいConfigを作成する
func NewConfigForParseApiCost(operation, filePath, textInput string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	if operation != OperationParseAPICost {
		return nil, fmt.Errorf("この関数はparse-api-cost操作専用です: %s", operation)
	}

	// 排他制御
	if filePath != "" && textInput != "" {
		return nil, fmt.Errorf("ファイルパスとテキスト入力は同時に指定できません")
	}
	if filePath == "" && textInput == "" {
		return nil, fmt.Errorf("ファイルパスまたはテキスト入力のいずれかを指定してください")
	}

	// ファイル拡張子の検証
	if filePath != "" {
		if !strings.HasSuffix(filePath, ".md") && !strings.HasSuffix(filePath, ".txt") {
			return nil, fmt.Errorf("ファイルは.mdまたは.txt形式である必要があります")
		}
	}

	return &Config{
		Operation: operation,
		FilePath:  filePath,
		TextInput: textInput,
	}, nil
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation    = ""
		xStr         = "0"
		yStr         = "0"
		numbers      = ""
		filePath     = ""
		thresholdStr = "0"
		textInput    = ""
		baseStr      = "0"
		exponentStr  = "0"
		numberStr    = "0"
		nStr         = "0"
		function     = ""
		angleStr     = "0"
		unit         = ""
		expression   = ""
		help         = false
	)

	parser.StringVar(&operation, "operation", operation, fmt.Sprintf("算術操作 (%s)", strings.Join(SupportedOperations(), ", ")))
	parser.StringVar(&operation, "o", operation, "算術操作の短縮形")

	// 基本計算用のパラメータ
	parser.StringVar(&xStr, "x", xStr, "第一オペランド")
	parser.StringVar(&yStr, "y", yStr, "第二オペランド")

	// 配列計算用のパラメータ
	parser.StringVar(&numbers, "numbers", numbers, "カンマ区切りの数値リスト (sum操作用)")
	parser.StringVar(&numbers, "nums", numbers, "数値リストの短縮形")

	// ファイル評価用のパラメータ
	parser.StringVar(&filePath, "file", filePath, "評価するファイルのパス")
	parser.StringVar(&filePath, "f", filePath, "ファイルパスの短縮形")
	parser.StringVar(&thresholdStr, "threshold", thresholdStr, "行数の閾値")
	parser.StringVar(&thresholdStr, "t", thresholdStr, "閾値の短縮形")

	// API料金抽出用のパラメータ (parse-api-cost操作用)
	parser.StringVar(&textInput, "text-input", textInput, "テキスト入力 (parse-api-cost操作用)")
	parser.StringVar(&textInput, "ti", textInput, "テキスト入力の短縮形")

	// 高度な数学演算用のパラメータ
	parser.StringVar(&baseStr, "base", baseStr, "べき乗の底 (power操作用)")
	parser.StringVar(&baseStr, "b", baseStr, "べき乗の底の短縮形")
	parser.StringVar(&exponentStr, "exponent", exponentStr, "べき乗の指数 (power操作用)")
	parser.StringVar(&exponentStr, "exp", exponentStr, "べき乗の指数の短縮形")
	parser.StringVar(&numberStr, "number", numberStr, "単一の数値 (square_root操作用)")
	parser.StringVar(&nStr, "n", nStr, "階乗の数値 (factorial操作用)")

	// 三角関数用のパラメータ
	parser.StringVar(&function, "function", function, "三角関数の種類 (sin, cos, tan)")
	parser.StringVar(&function, "func", function, "三角関数の種類の短縮形")
	parser.StringVar(&angleStr, "angle", angleStr, "角度")
	parser.StringVar(&angleStr, "a", angleStr, "角度の短縮形")
	parser.StringVar(&unit, "unit", unit, "角度の単位 (radians, degrees)")
	parser.StringVar(&unit, "u", unit, "角度の単位の短縮形")

	// 数式評価用のパラメータ
	parser.StringVar(&expression, "expression", expression, "評価する数式 (calculate操作用)")
	parser.StringVar(&expression, "expr", expression, "数式の短縮形")

	parser.BoolVar(&help, "help", help, "ヘルプを表示")
	parser.BoolVar(&help, "h", help, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 文字列から数値に変換
	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なx値です: %s", xStr)
	}

	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なy値です: %s", yStr)
	}

	threshold, err := strconv.Atoi(thresholdStr)
	if err != nil {
		return nil, fmt.Errorf("無効な閾値です: %s", thresholdStr)
	}

	// 新しいパラメータの変換
	base, err := strconv.ParseFloat(baseStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なbase値です: %s", baseStr)
	}

	exponent, err := strconv.ParseFloat(exponentStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なexponent値です: %s", exponentStr)
	}

	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なnumber値です: %s", numberStr)
	}

	n, err := strconv.Atoi(nStr)
	if err != nil {
		return nil, fmt.Errorf("無効なn値です: %s", nStr)
	}

	angle, err := strconv.ParseFloat(angleStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効なangle値です: %s", angleStr)
	}

	// 残りの引数から x, y, threshold を取得（位置引数として）
	args := parser.Args()
	if len(args) >= 1 && operation != OperationSum && operation != OperationEvaluateLineCount {
		if val, err := strconv.ParseFloat(args[0], 64); err == nil {
			x = val
		}
	}
	if len(args) >= 2 && operation != OperationSum && operation != OperationEvaluateLineCount {
		if val, err := strconv.ParseFloat(args[1], 64); err == nil {
			y = val
		}
	}
	if len(args) >= 1 && operation == OperationEvaluateLineCount {
		if val, err := strconv.Atoi(args[0]); err == nil {
			threshold = val
		}
	}

	// numbers文字列を[]float64に変換
	var numbersSlice []float64
	if numbers != "" {
		parts := strings.Split(numbers, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if num, err := strconv.ParseFloat(part, 64); err == nil {
				numbersSlice = append(numbersSlice, num)
			} else {
				return nil, fmt.Errorf("無効な数値です: %s", part)
			}
		}
	}

	// parse-api-cost操作の場合は専用の処理を行う
	if operation == OperationParseAPICost {
		return NewConfigForParseApiCost(operation, filePath, textInput)
	}

	return NewConfig(operation, x, y, numbersSlice, filePath, threshold, base, exponent, number, angle, n, function, unit, expression, textInput)
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `算術計算CLIツール

使用方法:
  基本計算:
    %s -operation add -x 10 -y 5
    %s -o add 10 5

  配列計算:
    %s -operation sum -numbers 1,2,3,4,5
    %s -o sum -n 1,2,3,4,5

  ファイル行数評価:
    %s -operation evaluate_line_count -file /path/to/file -threshold 100
    %s -o evaluate_line_count -f /path/to/file 100

  API料金抽出（ファイルから）:
    %s -operation parse-api-cost -file /path/to/file.txt
    %s -o parse-api-cost -f /path/to/file.md

  API料金抽出（テキストから）:
    %s -operation parse-api-cost -text-input "API料金が100円掛かった。別のAPI料金が200円掛かった。"
    %s -o parse-api-cost -ti "API料金が150円掛かった。"

オプション:
  -operation, -o    算術操作 (%s)
  -x               第一オペランド (基本計算用)
  -y               第二オペランド (基本計算用)
  -numbers, -nums  カンマ区切りの数値リスト (sum操作用)
  -file, -f        評価するファイルのパス (evaluate_line_count, parse-api-cost操作用)
  -threshold, -t   行数の閾値 (evaluate_line_count操作用)
  -text-input, -ti テキスト入力 (parse-api-cost操作用)
  -help, -h        このヘルプを表示

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], strings.Join(SupportedOperations(), ", "))
}
