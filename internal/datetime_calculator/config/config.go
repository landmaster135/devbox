package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config は日時計算CLIの設定を保持する構造体
type Config struct {
	Operation      string    // 操作タイプ (add, subtract, sum, parse-time)
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
	FilePath       string    // ファイルパス (parse-time操作用)
	TextInput      string    // テキスト入力 (parse-time操作用)
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

// NewConfigForParseTime はparse-time操作用の新しいConfigを作成する
func NewConfigForParseTime(operation, filePath, textInput, outputUnit string) (*Config, error) {
	if operation == "" {
		return nil, fmt.Errorf("操作タイプが指定されていません")
	}

	if operation != "parse-time" {
		return nil, fmt.Errorf("この関数はparse-time操作専用です: %s", operation)
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

	// outputUnitのデフォルト値設定
	if outputUnit == "" {
		outputUnit = "minute"
	}

	// 時間単位の検証
	validUnits := []string{"year", "month", "day", "hour", "minute", "second"}
	isValidOutput := false
	for _, unit := range validUnits {
		if outputUnit == unit {
			isValidOutput = true
			break
		}
	}
	if !isValidOutput {
		return nil, fmt.Errorf("無効な出力時間単位です: %s", outputUnit)
	}

	return &Config{
		Operation:  operation,
		FilePath:   filePath,
		TextInput:  textInput,
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
		filePath          = ""
		textInput         = ""
		help              = false
	)

	parser.StringVar(&operation, "operation", operation, "日時操作 (add, subtract, sum, parse-time)")
	parser.StringVar(&operation, "o", operation, "操作の短縮形")

	// 基準日時のパラメータ
	parser.StringVar(&year1Str, "year", year1Str, "基準日時の年")
	parser.StringVar(&year1Str, "y", year1Str, "年の短縮形")
	parser.StringVar(&month1Str, "month", month1Str, "基準日時の月")
	parser.StringVar(&month1Str, "m", month1Str, "月の短縮形")
	parser.StringVar(&day1Str, "day", day1Str, "基準日時の日")
	parser.StringVar(&day1Str, "d", day1Str, "日の短縮形")
	parser.StringVar(&hour1Str, "hour", hour1Str, "基準日時の時")
	parser.StringVar(&hour1Str, "hr", hour1Str, "時の短縮形")
	parser.StringVar(&minute1Str, "minute", minute1Str, "基準日時の分")
	parser.StringVar(&minute1Str, "min", minute1Str, "分の短縮形")
	parser.StringVar(&second1Str, "second", second1Str, "基準日時の秒")
	parser.StringVar(&second1Str, "s", second1Str, "秒の短縮形")

	// 期間のパラメータ
	parser.StringVar(&durationYearStr, "duration-year", durationYearStr, "加算/減算する年")
	parser.StringVar(&durationYearStr, "dy", durationYearStr, "期間年の短縮形")
	parser.StringVar(&durationMonthStr, "duration-month", durationMonthStr, "加算/減算する月")
	parser.StringVar(&durationMonthStr, "dm", durationMonthStr, "期間月の短縮形")
	parser.StringVar(&durationDayStr, "duration-day", durationDayStr, "加算/減算する日")
	parser.StringVar(&durationDayStr, "dd", durationDayStr, "期間日の短縮形")
	parser.StringVar(&durationHourStr, "duration-hour", durationHourStr, "加算/減算する時")
	parser.StringVar(&durationHourStr, "dh", durationHourStr, "期間時の短縮形")
	parser.StringVar(&durationMinuteStr, "duration-minute", durationMinuteStr, "加算/減算する分")
	parser.StringVar(&durationMinuteStr, "dmin", durationMinuteStr, "期間分の短縮形")
	parser.StringVar(&durationSecondStr, "duration-second", durationSecondStr, "加算/減算する秒")
	parser.StringVar(&durationSecondStr, "ds", durationSecondStr, "期間秒の短縮形")

	// 時間単位計算用のパラメータ (sum操作用)
	parser.StringVar(&figures, "figures", figures, "カンマ区切りの数値リスト (sum操作用)")
	parser.StringVar(&figures, "f", figures, "数値リストの短縮形")
	parser.StringVar(&inputUnit, "input-unit", inputUnit, "入力時間単位 (sum操作用)")
	parser.StringVar(&inputUnit, "iu", inputUnit, "入力時間単位の短縮形")
	parser.StringVar(&outputUnit, "output-unit", outputUnit, "出力時間単位 (sum操作用)")
	parser.StringVar(&outputUnit, "ou", outputUnit, "出力時間単位の短縮形")

	// 時間抽出用のパラメータ (parse-time操作用)
	parser.StringVar(&filePath, "file-path", filePath, "ファイルパス (parse-time操作用)")
	parser.StringVar(&filePath, "fp", filePath, "ファイルパスの短縮形")
	parser.StringVar(&textInput, "text-input", textInput, "テキスト入力 (parse-time操作用)")
	parser.StringVar(&textInput, "ti", textInput, "テキスト入力の短縮形")

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

	// parse-time操作の場合は専用の処理を行う
	if operation == "parse-time" {
		return NewConfigForParseTime(operation, filePath, textInput, outputUnit)
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

  時間抽出（ファイルから）:
    %s -operation parse-time -file-path /path/to/file.txt

  時間抽出（テキストから）:
    %s -operation parse-time -text-input "作業は合計30分掛かった。別の作業は合計45分掛かった。"

  時間抽出（単位変換）:
    %s -operation parse-time -text-input "合計120分掛かった。" -output-unit hour

  短縮形:
    %s -o add -y 2025 -m 1 -d 15 -hr 12 -min 30 -s 0 -dy 1 -dm 2 -dd 10 -dh 5 -dmin 30 -ds 45
    %s -o sum -iu second -ou hour -f 3600,1800,7200
    %s -o parse-time -fp /path/to/file.md
    %s -o parse-time -ti "合計120分掛かった。" -ou second

オプション:
  -operation, -o       日時操作 (add, subtract, sum, parse-time)
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
  -output-unit, -ou   出力時間単位 (sum, parse-time操作用: year, month, day, hour, minute, second, デフォルト: minute)
  -file-path, -fp     ファイルパス (parse-time操作用: .mdまたは.txt)
  -text-input, -ti    テキスト入力 (parse-time操作用)
  -help, -h           このヘルプを表示

`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}
