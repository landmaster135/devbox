package factorial

import "fmt"

func Calculate(n int) (float64, error) {
	if n < 0 {
		return 0, fmt.Errorf("負数の階乗は定義されていません")
	}
	if n > 170 {
		return 0, fmt.Errorf("数値が大きすぎて階乗計算でオーバーフローします")
	}

	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	return result, nil
}

func HandleToFactorial(n int) (float64, error) {
	return Calculate(n)
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute(n int) (string, error) {
	result, err := HandleToFactorial(n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d! = %.0f\n", n, result), nil
}
