package basiccalculation

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
)

func Add(x, y float64) float64 {
	return x + y
}

func Subtract(x, y float64) float64 {
	return x - y
}

func Multiply(x, y float64) float64 {
	return x * y
}

func Divide(x, y float64) float64 {
	return x / y
}

func Sum(numbers []float64) float64 {
	result := 0.0
	for _, n := range numbers {
		result += n
	}
	return result
}

func HandleToCalculate(operation string, x, y float64) (float64, error) {
	switch operation {
	case config.OperationAdd:
		return Add(x, y), nil
	case config.OperationSubtract:
		return Subtract(x, y), nil
	case config.OperationMultiply:
		return Multiply(x, y), nil
	case config.OperationDivide:
		if y == 0 {
			return 0, fmt.Errorf("division by zero is not allowed")
		}
		return Divide(x, y), nil
	default:
		return 0, nil
	}
}

func HandleToCalculateWithArray(operation string, numbers []float64) (float64, error) {
	switch operation {
	case config.OperationSum:
		return Sum(numbers), nil
	default:
		return 0, nil
	}
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(operation string, x, y float64) (string, error) {
	result, err := HandleToCalculate(operation, x, y)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f %s %.2f = %.2f\n", x, operationSymbol(operation), y, result), nil
}

func operationSymbol(operation string) string {
	switch operation {
	case config.OperationAdd:
		return "+"
	case config.OperationSubtract:
		return "-"
	case config.OperationMultiply:
		return "*"
	case config.OperationDivide:
		return "/"
	default:
		return operation
	}
}
