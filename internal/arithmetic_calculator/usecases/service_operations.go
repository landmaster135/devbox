package usecases

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
	arraysum "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/array_sum"
	basiccalculation "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/basic_calculation"
	calculateexpression "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/calculate_expression"
	evaluatelinecount "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/evaluate_line_count"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/factorial"
	getconstants "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/get_constants"
	parseapicost "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/parse_api_cost"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/power"
	squareroot "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/square_root"
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases/operations/trigonometry"
)

// NewService はCLIから利用するoperation集約サービスを生成する
func NewService() *Service {
	return newServiceWithOperations(
		basiccalculation.NewService(),
		arraysum.NewService(),
		evaluatelinecount.NewService(),
		parseapicost.NewService(),
		power.NewService(),
		squareroot.NewService(),
		factorial.NewService(),
		trigonometry.NewService(),
		calculateexpression.NewService(),
		getconstants.NewService(),
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
