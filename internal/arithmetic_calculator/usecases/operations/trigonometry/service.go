package trigonometry

import (
	"fmt"
	"math"
	"strings"
)

func Calculate(function string, angle float64, unit string) (float64, error) {
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

func HandleToTrigonometry(function string, angle float64, unit string) (float64, error) {
	return Calculate(function, angle, unit)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(function string, angle float64, unit string) (string, error) {
	result, err := HandleToTrigonometry(function, angle, unit)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s(%.2f %s) = %.6f\n", function, angle, unit, result), nil
}
