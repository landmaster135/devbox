package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config は日時計算CLIの設定を保持する構造体
type Config struct {
	Operation      string    // 操作タイプ (add, subtract, sum)
	Year1          float64   // 基準日時の年
	Month1         float64   // 基準日時の月
	Day1           float64   // 基準日時の日
	Hour1          float64   // 基準日時の時
	Minute1        float64   // 基準日時の分
	Second1        float64   // 基準日時の秒
	DurationYear   float64   // 加算/減算する年
	DurationMonth  float64   // 加算/減算する月
	DurationDay    float64   // 加算/減算する日
	DurationHour   float64   // 加算/減算する時
	DurationMinute float64   // 加算/減算する分
	DurationSecond float64   // 加算/減算する秒
	Figures        []float64 // 時間単位計算用の数値配列 (sum操作用)
	InputUnit      string    // 入力時間単位 (sum操作用)
	OutputUnit     string    // 出力時間単位 (sum操作用)
	Help           bool      // ヘルプ表示フラグ
}

// NewConfig は新しいConfigを作成する
func NewConfig(operation string, year1, month1, day1, hour1, minute1, second1 float64, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond float64) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// 操作タイプの検証
	validOperations := []string{"add", "subtract", "sum"}
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

	return &Config{
		Operation:      operation,
		Year1:          year1,
		Month1:         month1,
		Day1:           day1,
		Hour1:          hour1,
		Minute1:        minute1,
		Second1:        second1,
		DurationYear:   durationYear,
		DurationMonth:  durationMonth,
		DurationDay:    durationDay,
		DurationHour:   durationHour,
		DurationMinute: durationMinute,
		DurationSecond: durationSecond,
	}, nil
}

// NewConfigForSum はsum操作用の新しいConfigを作成する
func NewConfigForSum(operation string, figures []float64, inputUnit, outputUnit string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	if operation != "sum" {
		return nil, fmt.Errorf("この関数はsum操作専用です: %s", operation)
	}

	if len(figures) == 0 {
		return nil, fmt.Errorf("sum操作には数値の配列が必要です")
	}

	if inputUnit == "" || outputUnit == "" {
		return nil, fmt.Errorf("sum操作には入力単位と出力単位が必要です")
	}

	// 時間単位の検証
	validUnits := []string{"year", "month", "day", "hour", "minute", "second"}
	isValidInput := false
	isValidOutput := false
	for _, unit := range validUnits {
		if inputUnit == unit {
			isValidInput = true
		}
		if outputUnit == unit {
			isValidOutput = true
		}
	}
	if !isValidInput {
		return nil, fmt.Errorf("無効な入力時間単位です: %s", inputUnit)
	}
	if !isValidOutput {
		return nil, fmt.Errorf("無効な出力時間単位です: %s", outputUnit)
	}

	return &Config{
		Operation:  operation,
		Figures:    figures,
		InputUnit:  inputUnit,
		OutputUnit: outputUnit,
	}, nil
}

// ParseFlagsWithParser は指定されたFlagParserを使用してコマンドライン引数を解析する
func ParseFlagsWithParser(parser FlagParser) (*Config, error) {
	var (
		operation         = ""
		year1Str          = "2025"
		month1Str         = "1"
		day1Str           = "1"
		hour1Str          = "0"
		minute1Str        = "0"
		second1Str        = "0"
		durationYearStr   = "0"
		durationMonthStr  = "0"
		durationDayStr    = "0"
		durationHourStr   = "0"
		durationMinuteStr = "0"
		durationSecondStr = "0"
		figures           = ""
		inputUnit         = ""
		outputUnit        = ""
		help              = false
	)

	parser.StringVar(&operation, "operation", "", "日時操作 (add, subtract, sum)")
	parser.StringVar(&operation, "o", "", "操作の短縮形")

	// 基準日時のパラメータ
	parser.StringVar(&year1Str, "year", "2025", "基準日時の年")
	parser.StringVar(&year1Str, "y", "2025", "年の短縮形")
	parser.StringVar(&month1Str, "month", "1", "基準日時の月")
	parser.StringVar(&month1Str, "m", "1", "月の短縮形")
	parser.StringVar(&day1Str, "day", "1", "基準日時の日")
	parser.StringVar(&day1Str, "d", "1", "日の短縮形")
	parser.StringVar(&hour1Str, "hour", "0", "基準日時の時")
	parser.StringVar(&hour1Str, "hr", "0", "時の短縮形")
	parser.StringVar(&minute1Str, "minute", "0", "基準日時の分")
	parser.StringVar(&minute1Str, "min", "0", "分の短縮形")
	parser.StringVar(&second1Str, "second", "0", "基準日時の秒")
	parser.StringVar(&second1Str, "s", "0", "秒の短縮形")

	// 期間のパラメータ
	parser.StringVar(&durationYearStr, "duration-year", "0", "加算/減算する年")
	parser.StringVar(&durationYearStr, "dy", "0", "期間年の短縮形")
	parser.StringVar(&durationMonthStr, "duration-month", "0", "加算/減算する月")
	parser.StringVar(&durationMonthStr, "dm", "0", "期間月の短縮形")
	parser.StringVar(&durationDayStr, "duration-day", "0", "加算/減算する日")
	parser.StringVar(&durationDayStr, "dd", "0", "期間日の短縮形")
	parser.StringVar(&durationHourStr, "duration-hour", "0", "加算/減算する時")
	parser.StringVar(&durationHourStr, "dh", "0", "期間時の短縮形")
	parser.StringVar(&durationMinuteStr, "duration-minute", "0", "加算/減算する分")
	parser.StringVar(&durationMinuteStr, "dmin", "0", "期間分の短縮形")
	parser.StringVar(&durationSecondStr, "duration-second", "0", "加算/減算する秒")
	parser.StringVar(&durationSecondStr, "ds", "0", "期間秒の短縮形")

	// 時間単位計算用のパラメータ (sum操作用)
	parser.StringVar(&figures, "figures", "", "カンマ区切りの数値リスト (sum操作用)")
	parser.StringVar(&figures, "f", "", "数値リストの短縮形")
	parser.StringVar(&inputUnit, "input-unit", "", "入力時間単位 (sum操作用)")
	parser.StringVar(&inputUnit, "iu", "", "入力時間単位の短縮形")
	parser.StringVar(&outputUnit, "output-unit", "", "出力時間単位 (sum操作用)")
	parser.StringVar(&outputUnit, "ou", "", "出力時間単位の短縮形")

	parser.BoolVar(&help, "help", false, "ヘルプを表示")
	parser.BoolVar(&help, "h", false, "ヘルプの短縮形")

	if err := parser.Parse(); err != nil {
		return nil, fmt.Errorf("フラグの解析に失敗しました: %v", err)
	}

	// ヘルプが要求された場合
	if help {
		return &Config{Help: true}, nil
	}

	// 文字列から数値に変換
	year1, err := strconv.ParseFloat(year1Str, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な年の値です: %s", year1Str)
	}

	month1, err := strconv.ParseFloat(month1Str, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な月の値です: %s", month1Str)
	}

	day1, err := strconv.ParseFloat(day1Str, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な日の値です: %s", day1Str)
	}

	hour1, err := strconv.ParseFloat(hour1Str, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な時の値です: %s", hour1Str)
	}

	minute1, err := strconv.ParseFloat(minute1Str, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な分の値です: %s", minute1Str)
	}

	second1, err := strconv.ParseFloat(second1Str, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な秒の値です: %s", second1Str)
	}

	durationYear, err := strconv.ParseFloat(durationYearStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な期間年の値です: %s", durationYearStr)
	}

	durationMonth, err := strconv.ParseFloat(durationMonthStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な期間月の値です: %s", durationMonthStr)
	}

	durationDay, err := strconv.ParseFloat(durationDayStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な期間日の値です: %s", durationDayStr)
	}

	durationHour, err := strconv.ParseFloat(durationHourStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な期間時の値です: %s", durationHourStr)
	}

	durationMinute, err := strconv.ParseFloat(durationMinuteStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な期間分の値です: %s", durationMinuteStr)
	}

	durationSecond, err := strconv.ParseFloat(durationSecondStr, 64)
	if err != nil {
		return nil, fmt.Errorf("無効な期間秒の値です: %s", durationSecondStr)
	}

	// 操作タイプが指定されていない場合のエラーチェック
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	// sum操作の場合は専用の処理を行う
	if operation == "sum" {
		// figures文字列を[]float64に変換
		var figuresSlice []float64
		if figures != "" {
			parts := strings.Split(figures, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if num, err := strconv.ParseFloat(part, 64); err == nil {
					figuresSlice = append(figuresSlice, num)
				} else {
					return nil, fmt.Errorf("無効な数値です: %s", part)
				}
			}
		}

		return NewConfigForSum(operation, figuresSlice, inputUnit, outputUnit)
	}

	return NewConfig(operation, year1, month1, day1, hour1, minute1, second1, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond)
}

// ParseFlags はコマンドライン引数を解析してConfigを作成する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithParser(NewStandardFlagParser())
}

// PrintUsage は使用方法を表示する
func PrintUsage() {
	fmt.Fprintf(os.Stderr, `日時計算CLIツール

使用方法:
  日時加算:
    %s -operation add -year 2025 -month 1 -day 15 -hour 12 -minute 30 -second 0 -duration-year 1 -duration-month 2 -duration-day 10 -duration-hour 5 -duration-minute 30 -duration-second 45

  日時減算:
    %s -operation subtract -year 2025 -month 3 -day 25 -hour 18 -minute 0 -second 45 -duration-year 0 -duration-month 1 -duration-day 5 -duration-hour 2 -duration-minute 15 -duration-second 30

  時間単位計算（合計）:
    %s -operation sum -input-unit second -output-unit hour -figures 3600,1800,7200

  時間単位変換:
    %s -operation sum -input-unit hour -output-unit minute -figures 2.5

  短縮形:
    %s -o add -y 2025 -m 1 -d 15 -hr 12 -min 30 -s 0 -dy 1 -dm 2 -dd 10 -dh 5 -dmin 30 -ds 45
    %s -o sum -iu second -ou hour -f 3600,1800,7200

オプション:
  -operation, -o       日時操作 (add, subtract, sum)
  -year, -y           基準日時の年 (デフォルト: 2025)
  -month, -m          基準日時の月 (デフォルト: 1)
  -day, -d            基準日時の日 (デフォルト: 1)
  -hour, -hr          基準日時の時 (デフォルト: 0)
  -minute, -min       基準日時の分 (デフォルト: 0)
  -second, -s         基準日時の秒 (デフォルト: 0)
  -duration-year, -dy  加算/減算する年 (デフォルト: 0)
  -duration-month, -dm 加算/減算する月 (デフォルト: 0)
  -duration-day, -dd   加算/減算する日 (デフォルト: 0)
  -duration-hour, -dh  加算/減算する時 (デフォルト: 0)
  -duration-minute, -dmin 加算/減算する分 (デフォルト: 0)
  -duration-second, -ds 加算/減算する秒 (デフォルト: 0)
  -figures, -f        カンマ区切りの数値リスト (sum操作用)
  -input-unit, -iu    入力時間単位 (sum操作用: year, month, day, hour, minute, second)
  -output-unit, -ou   出力時間単位 (sum操作用: year, month, day, hour, minute, second)
  -help, -h           このヘルプを表示

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
