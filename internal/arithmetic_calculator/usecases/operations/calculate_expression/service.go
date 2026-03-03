package calculateexpression

import (
	"fmt"
	"strings"

	getConstants "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/get_constants"
	reversePolishNotation "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/reverse_polish_notation"
)

type constantsGetter func() map[string]float64

type Service struct {
	getConstants constantsGetter
}

func NewService() *Service {
	return &Service{getConstants: getConstants.GetConstants}
}

func NewServiceWithConstantsGetter(getter constantsGetter) *Service {
	if getter == nil {
		getter = getConstants.GetConstants
	}
	return &Service{getConstants: getter}
}

func (s *Service) SafeEvaluate(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")

	dangerousPatterns := []string{
		"__", "import", "exec", "eval", "open", "file", "input", "sys",
	}
	for _, pattern := range dangerousPatterns {
		if strings.Contains(strings.ToLower(expression), pattern) {
			return 0, fmt.Errorf("危険なパターンが検出されました: %s", pattern)
		}
	}

	if err := s.CheckOSPattern(expression); err != nil {
		return 0, err
	}

	constants := s.getConstants()
	for name, value := range constants {
		expression = strings.ReplaceAll(expression, name, fmt.Sprintf("%f", value))
	}

	result, err := s.EvaluateBasicExpression(expression)
	if err != nil {
		return 0, fmt.Errorf("数式の評価に失敗しました: %v", err)
	}

	return result, nil
}

func (s *Service) EvaluateBasicExpression(expression string) (float64, error) {
	return s.evaluateUsingReversePolish(expression)
}

func (s *Service) EvaluateArithmeticExpression(expression string) (float64, error) {
	return s.evaluateUsingReversePolish(expression)
}

func (s *Service) evaluateUsingReversePolish(expression string) (float64, error) {
	cleaned := strings.ReplaceAll(expression, " ", "")
	cleaned = strings.ReplaceAll(cleaned, "**", "^")
	if cleaned == "" {
		return 0, fmt.Errorf("無効な数式です: %s", expression)
	}
	return reversePolishNotation.Evaluate(cleaned)
}

func (s *Service) CheckOSPattern(expression string) error {
	lowerExpr := strings.ToLower(expression)

	osIndices := s.GetAllIndices(lowerExpr, "os")
	if len(osIndices) == 0 {
		return nil
	}

	cosIndices := s.GetAllIndices(lowerExpr, "cos")

	for _, osIndex := range osIndices {
		isPartOfCos := false
		for _, cosIndex := range cosIndices {
			if osIndex == cosIndex+1 {
				isPartOfCos = true
				break
			}
		}
		if !isPartOfCos {
			return fmt.Errorf("危険なパターンが検出されました: os")
		}
	}

	return nil
}

func (s *Service) GetAllIndices(text, pattern string) []int {
	var indices []int
	start := 0

	for {
		index := strings.Index(text[start:], pattern)
		if index == -1 {
			break
		}
		actualIndex := start + index
		indices = append(indices, actualIndex)
		start = actualIndex + 1
	}

	return indices
}

func (s *Service) HandleToCalculateExpression(expression string) (float64, error) {
	return s.SafeEvaluate(expression)
}

func (s *Service) Execute(expression string) (string, error) {
	result, err := s.HandleToCalculateExpression(expression)
	if err != nil {
		return "", err
	}

	if result > -1e-6 && result < 1e-6 {
		result = 0.0
	}

	return fmt.Sprintf("%s = %.2f\n", expression, result), nil
}
