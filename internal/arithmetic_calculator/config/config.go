package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config は算術計算CLIの設定を保持する構造体
type Config struct {
	Operation string    // 操作タイプ (add, subtract, multiply, divide, sum, evaluate_line_count)
	X         float64   // 第一オペランド
	Y         float64   // 第二オペランド
	Numbers   []float64 // 複数の数値（sum操作用）
	FilePath  string    // ファイルパス（evaluate_line_count操作用）
	Threshold int       // 閾値（evaluate_line_count操作用）
	Help      bool      // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation string, x, y float64, numbers []float64, filePath string, threshold int) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"add", "subtract", "multiply", "divide", "sum", "evaluate_line_count"}
	isValid := false
	for _, op := range validOperations {
		if operation == op {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, fmt.Errorf("無効な操作タイプです: %s", operation)
	}

	// 操作タイプ別の検証
	switch operation {
	case "add", "subtract", "multiply", "divide":
		// 二項演算の場合、x, yが必要
		// デフォルト値0.0でも計算可能なので、特別な検証は不要
	case "sum":
		if len(numbers) == 0 {
			return nil, fmt.Errorf("sum操作には数値の配列が必要です")
		}
	case "evaluate_line_count":
		if filePath == "" {
			return nil, fmt.Errorf("evaluate_line_count操作にはファイルパスが必要です")
		}
		if threshold < 0 {
			return nil, fmt.Errorf("閾値は0以上である必要があります")
		}
	}

	return &Config{
		Operation: operation,
		X:         x,
		Y:         y,
		Numbers:   numbers,
		FilePath:  filePath,
		Threshold: threshold,
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
		help         = false
	)

	parser.StringVar(&operation, "operation", operation, "算術操作 (add, subtract, multiply, divide, sum, evaluate_line_count)")
	parser.StringVar(&operation, "o", operation, "算術操作の短縮形")

	// 基本計算用のパラメータ
	parser.StringVar(&xStr, "x", xStr, "第一オペランド")
	parser.StringVar(&yStr, "y", yStr, "第二オペランド")

	// 配列計算用のパラメータ
	parser.StringVar(&numbers, "numbers", numbers, "カンマ区切りの数値リスト (sum操作用)")
	parser.StringVar(&numbers, "n", numbers, "数値リストの短縮形")

	// ファイル評価用のパラメータ
	parser.StringVar(&filePath, "file", filePath, "評価するファイルのパス")
	parser.StringVar(&filePath, "f", filePath, "ファイルパスの短縮形")
	parser.StringVar(&thresholdStr, "threshold", thresholdStr, "行数の閾値")
	parser.StringVar(&thresholdStr, "t", thresholdStr, "閾値の短縮形")

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

	// 残りの引数から x, y, threshold を取得（位置引数として）
	args := parser.Args()
	if len(args) >= 1 && operation != "sum" && operation != "evaluate_line_count" {
		if val, err := strconv.ParseFloat(args[0], 64); err == nil {
			x = val
		}
	}
	if len(args) >= 2 && operation != "sum" && operation != "evaluate_line_count" {
		if val, err := strconv.ParseFloat(args[1], 64); err == nil {
			y = val
		}
	}
	if len(args) >= 1 && operation == "evaluate_line_count" {
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

	return NewConfig(operation, x, y, numbersSlice, filePath, threshold)
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

オプション:
  -operation, -o    算術操作 (add, subtract, multiply, divide, sum, evaluate_line_count)
  -x               第一オペランド (基本計算用)
  -y               第二オペランド (基本計算用)
  -numbers, -n     カンマ区切りの数値リスト (sum操作用)
  -file, -f        評価するファイルのパス (evaluate_line_count操作用)
  -threshold, -t   行数の閾値 (evaluate_line_count操作用)
  -help, -h        このヘルプを表示

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
