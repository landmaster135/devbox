package arraysum

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	basiccalculation "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/basic_calculation"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(operation string, numbers []float64) (string, error) {
	result, err := basiccalculation.HandleToCalculateWithArray(operation, numbers)
	if err != nil {
		return "", err
	}
	if operation != config.OperationSum {
		return "", fmt.Errorf("未対応の配列演算です: %s", operation)
	}

	return fmt.Sprintf("sum(%v) = %.2f\n", numbers, result), nil
}
