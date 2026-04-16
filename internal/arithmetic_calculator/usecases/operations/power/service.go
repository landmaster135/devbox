package power

import (
	"fmt"
	"math"
)

func Calculate(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

func HandleToPower(base, exponent float64) (float64, error) {
	return Calculate(base, exponent), nil
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(base, exponent float64) (string, error) {
	result, err := HandleToPower(base, exponent)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f^%.2f = %.2f\n", base, exponent, result), nil
}
