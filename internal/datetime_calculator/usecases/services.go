package usecases

import (
	"fmt"
)

// #==============================================================#
// ##          DatetimeCalculatorService                         ##
// #==============================================================#
// DatetimeCalculatorService は日時計算を行うサービスです
type DatetimeCalculatorService struct {
	calculator *DatetimeCalculator
}

// NewDatetimeCalculatorService は新しいDatetimeCalculatorServiceを作成します
func NewDatetimeCalculatorService() *DatetimeCalculatorService {
	return &DatetimeCalculatorService{
		calculator: &DatetimeCalculator{},
	}
}

// HandleDatetimeCalc はMCPリクエストを処理して日時計算結果を返すハンドラーです
func (s *DatetimeCalculatorService) HandleDatetimeCalc(op string, year1, month1, day1, hour1, minute1, second1 float64, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond float64) (string, error) {
	var result string
	switch op {
	case "add":
		result = s.calculator.addDatetimeFloat(year1, month1, day1, hour1, minute1, second1, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond)
	case "subtract":
		result = s.calculator.subtractDatetimeFloat(year1, month1, day1, hour1, minute1, second1, durationYear, durationMonth, durationDay, durationHour, durationMinute, durationSecond)
	default:
		return "", fmt.Errorf("unsupported operation: %s", op)
	}
	return result, nil
}

// HandleTimeUnitSum は時間単位での合計計算を処理するハンドラーです
func (s *DatetimeCalculatorService) HandleTimeUnitSum(figures []float64, inputUnit, outputUnit string) (float64, error) {
	if len(figures) == 0 {
		return 0, fmt.Errorf("数値の配列が空です")
	}

	// 入力単位で合計を計算
	sum := 0.0
	for _, figure := range figures {
		sum += figure
	}

	// 時間単位変換を実行
	result, err := s.calculator.convertTimeUnit(sum, inputUnit, outputUnit)
	if err != nil {
		return 0, fmt.Errorf("時間単位変換に失敗しました: %v", err)
	}

	return result, nil
}
