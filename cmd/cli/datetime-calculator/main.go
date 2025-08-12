package main

import (
	"fmt"
	"os"
	"time"

	config "github.com/landmaster135/devbox/internal/datetime_calculator/config"
	usecases "github.com/landmaster135/devbox/internal/datetime_calculator/usecases"
)

// getOperationSymbol は操作タイプに対応する記号を返す
func getOperationSymbol(operation string) string {
	switch operation {
	case "add":
		return "+"
	case "subtract":
		return "-"
	default:
		return operation
	}
}

// handleDatetimeCalculation は日時計算を処理する
func handleDatetimeCalculation(cfg *config.Config) {
	// DatetimeCalculatorServiceを初期化
	service := usecases.NewDatetimeCalculatorService()

	// 計算を実行
	result, err := service.HandleDatetimeCalc(
		cfg.Operation,
		cfg.Year1, cfg.Month1, cfg.Day1, cfg.Hour1, cfg.Minute1, cfg.Second1,
		cfg.DurationYear, cfg.DurationMonth, cfg.DurationDay, cfg.DurationHour, cfg.DurationMinute, cfg.DurationSecond,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を整形して出力
	baseDateTime := fmt.Sprintf("%.0f-%.0f-%.0f %.0f:%.0f:%.0f",
		cfg.Year1, cfg.Month1, cfg.Day1, cfg.Hour1, cfg.Minute1, cfg.Second1)

	duration := formatDuration(cfg.DurationYear, cfg.DurationMonth, cfg.DurationDay, cfg.DurationHour, cfg.DurationMinute, cfg.DurationSecond)

	operationSymbol := getOperationSymbol(cfg.Operation)

	fmt.Printf("%s %s %s = %s\n", baseDateTime, operationSymbol, duration, result)
}

// formatDuration は期間を読みやすい形式でフォーマットする
func formatDuration(year, month, day, hour, minute, second float64) string {
	var parts []string

	if year != 0 {
		parts = append(parts, fmt.Sprintf("%.0f年", year))
	}
	if month != 0 {
		parts = append(parts, fmt.Sprintf("%.0f月", month))
	}
	if day != 0 {
		parts = append(parts, fmt.Sprintf("%.0f日", day))
	}
	if hour != 0 {
		parts = append(parts, fmt.Sprintf("%.0f時間", hour))
	}
	if minute != 0 {
		parts = append(parts, fmt.Sprintf("%.0f分", minute))
	}
	if second != 0 {
		parts = append(parts, fmt.Sprintf("%.0f秒", second))
	}

	if len(parts) == 0 {
		return "0秒"
	}

	result := ""
	for i, part := range parts {
		if i > 0 {
			result += ""
		}
		result += part
	}

	return result
}

// handleTimeUnitSum は時間単位の合計計算を処理する
func handleTimeUnitSum(cfg *config.Config) {
	// DatetimeCalculatorServiceを初期化
	service := usecases.NewDatetimeCalculatorService()

	// 計算を実行
	result, err := service.HandleTimeUnitSum(cfg.Figures, cfg.InputUnit, cfg.OutputUnit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 結果を出力
	fmt.Printf("sum(%v %s) = %.6f %s\n", cfg.Figures, cfg.InputUnit, result, cfg.OutputUnit)
}

// handleTimeExtraction は時間抽出処理を処理する
func handleTimeExtraction(cfg *config.Config) {
	// DatetimeCalculatorServiceを初期化
	service := usecases.NewDatetimeCalculatorService()

	// 時間抽出を実行
	result, err := service.HandleTimeExtraction(cfg.FilePath, cfg.TextInput, cfg.OutputUnit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	// 単位名を取得
	unitName := service.GetCalculator().GetUnitName(cfg.OutputUnit)

	// 結果を出力
	fmt.Printf("抽出された時間の合計: %.6f%s\n", result, unitName)
}

// getWeekdayJapanese は曜日を日本語で返す
func getWeekdayJapanese(weekday time.Weekday) string {
	switch weekday {
	case time.Sunday:
		return "Sun"
	case time.Monday:
		return "Mon"
	case time.Tuesday:
		return "Tue"
	case time.Wednesday:
		return "Wed"
	case time.Thursday:
		return "Thu"
	case time.Friday:
		return "Fri"
	case time.Saturday:
		return "Sat"
	default:
		return ""
	}
}

// handleGenerateDailyHeading は日次見出し生成を処理する
func handleGenerateDailyHeading(cfg *config.Config) {
	// 現在日時を取得
	now := time.Now()

	// day-offsetを適用して開始日を計算
	startDate := now.AddDate(0, 0, cfg.DayOffset)
	endDate := startDate.AddDate(0, 0, 1)

	// 日付と曜日をフォーマット
	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")
	startWeekday := getWeekdayJapanese(startDate.Weekday())
	endWeekday := getWeekdayJapanese(endDate.Weekday())

	// 指定されたフォーマットで出力
	fmt.Printf("## %s(%s)から%s(%s)にかけて（）\n", startDateStr, startWeekday, endDateStr, endWeekday)
	fmt.Printf("- [ ]  %s(%s)から%s(%s)にかけて進める。・・・合計0分掛かった。\n", startDateStr, startWeekday, endDateStr, endWeekday)
}

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
	case "add", "subtract":
		handleDatetimeCalculation(cfg)
	case "sum":
		handleTimeUnitSum(cfg)
	case "parse-time":
		handleTimeExtraction(cfg)
	case "generate-daily-heading":
		handleGenerateDailyHeading(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}
