package squareroot

import (
	"fmt"
	"math"
)

func Calculate(number float64) (float64, error) {
	if number < 0 {
		return 0, fmt.Errorf("負数の平方根は計算できません")
	}
	return math.Sqrt(number), nil
}

func HandleToSquareRoot(number float64) (float64, error) {
	return Calculate(number)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(number float64) (string, error) {
	result, err := HandleToSquareRoot(number)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("√%.2f = %.2f\n", number, result), nil
}
