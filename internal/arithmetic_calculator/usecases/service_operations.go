package usecases

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	arraySum "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/array_sum"
	basicCalculation "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/basic_calculation"
	calculateExpression "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/calculate_expression"
	evaluateLineCount "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/evaluate_line_count"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/factorial"
	getConstants "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/get_constants"
	parseAPICost "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/parse_api_cost"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/power"
	squareRoot "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/square_root"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/trigonometry"
)

// NewService はCLIから利用するoperation集約サービスを生成する
func NewService() *Service {
	return newServiceWithOperations(
		basicCalculation.NewService(),
		arraySum.NewService(),
		evaluateLineCount.NewService(),
		parseAPICost.NewService(),
		power.NewService(),
		squareRoot.NewService(),
		factorial.NewService(),
		trigonometry.NewService(),
		calculateExpression.NewService(),
		getConstants.NewService(),
	)
}

// ExecuteByConfig はoperationに応じて処理を振り分ける
func (s *Service) ExecuteByConfig(cfg *config.Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("configがnilです")
	}

	switch cfg.Operation {
	case config.OperationAdd, config.OperationSubtract, config.OperationMultiply, config.OperationDivide:
		return s.basicCalculationOperation.Execute(cfg.Operation, cfg.X, cfg.Y)
	case config.OperationSum:
		return s.arraySumOperation.Execute(cfg.Operation, cfg.Numbers)
	case config.OperationEvaluateLineCount:
		return s.lineCountOperation.Execute(cfg.FilePath, cfg.Threshold)
	case config.OperationParseAPICost:
		return s.apiCostOperation.Execute(cfg.FilePath, cfg.TextInput)
	case config.OperationPower:
		return s.powerOperation.Execute(cfg.Base, cfg.Exponent)
	case config.OperationSquareRoot:
		return s.squareRootOperation.Execute(cfg.Number)
	case config.OperationFactorial:
		return s.factorialOperation.Execute(cfg.N)
	case config.OperationTrigonometry:
		return s.trigonometryOperation.Execute(cfg.Function, cfg.Angle, cfg.Unit)
	case config.OperationCalculate:
		return s.calculateOperation.Execute(cfg.Expression)
	case config.OperationGetConstants:
		return s.constantsOperation.Execute()
	default:
		return "", fmt.Errorf("未サポートのoperationです: %s", cfg.Operation)
	}
}
