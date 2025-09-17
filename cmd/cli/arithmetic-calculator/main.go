package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	usecases "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases"
)

func main() {
	// コマンドライン引数を解析
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ヘルプが要求された場合
	if cfg.Help {
		config.PrintUsage()
		return
	}

	// 操作タイプに応じて処理を実行
	switch cfg.Operation {
	case "add", "subtract", "multiply", "divide":
		handleBasicCalculation(cfg)
	case "sum":
		handleArrayCalculation(cfg)
	case "evaluate_line_count":
		handleFileEvaluation(cfg)
	case "parse-api-cost":
		handleApiCostExtraction(cfg)
	case "power":
		handlePowerCalculation(cfg)
	case "square_root":
		handleSquareRootCalculation(cfg)
	case "factorial":
		handleFactorialCalculation(cfg)
	case "trigonometry":
		handleTrigonometryCalculation(cfg)
	case "calculate":
		handleExpressionCalculation(cfg)
	case "get_constants":
		handleGetConstants(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

// handleBasicCalculation は基本的な算術計算を処理する
func handleBasicCalculation(cfg *config.Config) {
	// CalculatorServiceを初期化
	service := usecases.NewCalculatorService()

	// 計算を実行
	result, err := service.HandleToCalculate(cfg.Operation, cfg.X, cfg.Y)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("%.2f %s %.2f = %.2f\n", cfg.X, getOperationSymbol(cfg.Operation), cfg.Y, result)
}

// handleArrayCalculation は配列を使った計算を処理する
func handleArrayCalculation(cfg *config.Config) {
	// CalculatorServiceを初期化
	service := usecases.NewCalculatorService()

	// 計算を実行
	result, err := service.HandleToCalculateWithArray(cfg.Operation, cfg.Numbers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("sum(%v) = %.2f\n", cfg.Numbers, result)
}

// handleFileEvaluation はファイルの行数評価を処理する
func handleFileEvaluation(cfg *config.Config) {
	// FileEvaluatorServiceを初期化
	service := usecases.NewFileEvaluatorService()

	// 評価を実行
	jsonResult, err := service.HandleToEvaluateLineCount(cfg.FilePath, cfg.Threshold)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(jsonResult)
}

// handleApiCostExtraction はAPI料金抽出を処理する
func handleApiCostExtraction(cfg *config.Config) {
	// ApiCostExtractorServiceを初期化
	service := usecases.NewApiCostExtractorService()

	// API料金抽出を実行
	result, err := service.HandleApiCostExtraction(cfg.FilePath, cfg.TextInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("抽出されたAPI料金の合計: %.0f円\n", result)
}

// handlePowerCalculation はべき乗計算を処理する
func handlePowerCalculation(cfg *config.Config) {
	// AdvancedMathServiceを初期化
	service := usecases.NewAdvancedMathService()

	// べき乗計算を実行
	result, err := service.HandleToPower(cfg.Base, cfg.Exponent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("%.2f^%.2f = %.2f\n", cfg.Base, cfg.Exponent, result)
}

// handleSquareRootCalculation は平方根計算を処理する
func handleSquareRootCalculation(cfg *config.Config) {
	// AdvancedMathServiceを初期化
	service := usecases.NewAdvancedMathService()

	// 平方根計算を実行
	result, err := service.HandleToSquareRoot(cfg.Number)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("√%.2f = %.2f\n", cfg.Number, result)
}

// handleFactorialCalculation は階乗計算を処理する
func handleFactorialCalculation(cfg *config.Config) {
	// AdvancedMathServiceを初期化
	service := usecases.NewAdvancedMathService()

	// 階乗計算を実行
	result, err := service.HandleToFactorial(cfg.N)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("%d! = %.0f\n", cfg.N, result)
}

// handleTrigonometryCalculation は三角関数計算を処理する
func handleTrigonometryCalculation(cfg *config.Config) {
	// TrigonometryServiceを初期化
	service := usecases.NewTrigonometryService()

	// 三角関数計算を実行
	result, err := service.HandleToTrigonometry(cfg.Function, cfg.Angle, cfg.Unit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("%s(%.2f %s) = %.6f\n", cfg.Function, cfg.Angle, cfg.Unit, result)
}

// handleExpressionCalculation は数式評価を処理する
func handleExpressionCalculation(cfg *config.Config) {
	// ExpressionEvaluatorServiceを初期化
	service := usecases.NewExpressionEvaluatorService()

	// 数式評価を実行
	result, err := service.HandleToCalculateExpression(cfg.Expression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 非常に小さな値を0として扱う（浮動小数点演算の誤差対策）
	if result > -1e-6 && result < 1e-6 {
		result = 0.0
	}

	// 結果を出力
	fmt.Printf("%s = %.2f\n", cfg.Expression, result)
}

// handleGetConstants は数学定数取得を処理する
func handleGetConstants(cfg *config.Config) {
	// MathConstantsServiceを初期化
	service := usecases.NewMathConstantsService()

	// 数学定数取得を実行
	result, err := service.HandleToGetConstants()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Print(result)
}

// getOperationSymbol は操作タイプに対応する記号を返す
func getOperationSymbol(operation string) string {
	switch operation {
	case "add":
		return "+"
	case "subtract":
		return "-"
	case "multiply":
		return "*"
	case "divide":
		return "/"
	default:
		return operation
	}
}
