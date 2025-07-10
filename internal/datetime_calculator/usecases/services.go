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
