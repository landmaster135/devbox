package usecases

import (
	"github.com/landmaster135/devbox/internal/arithmetic_calculator/config"
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

// CalculatorService は算術計算を行うサービスです。
// 実処理は operations/basic_calculation へ委譲します。
type CalculatorService struct{}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (c *CalculatorService) add(x float64, y float64) float64 {
	return basiccalculation.Add(x, y)
}

func (c *CalculatorService) subtract(x float64, y float64) float64 {
	return basiccalculation.Subtract(x, y)
}

func (c *CalculatorService) multiply(x float64, y float64) float64 {
	return basiccalculation.Multiply(x, y)
}

func (c *CalculatorService) divide(x float64, y float64) float64 {
	return basiccalculation.Divide(x, y)
}

func (c *CalculatorService) sum(numbers []float64) float64 {
	return basiccalculation.Sum(numbers)
}

func (c *CalculatorService) HandleToCalculate(op string, x, y float64) (float64, error) {
	return basiccalculation.HandleToCalculate(op, x, y)
}

func (c *CalculatorService) HandleToCalculateWithArray(op string, numbers []float64) (float64, error) {
	return basiccalculation.HandleToCalculateWithArray(op, numbers)
}

type FileOpener = evaluatelinecount.FileOpener
type DefaultFileOpener = evaluatelinecount.DefaultFileOpener
type BufioScanner = evaluatelinecount.BufioScanner
type JSONMarshaler = evaluatelinecount.JSONMarshaler
type DefaultJSONMarshaler = evaluatelinecount.DefaultJSONMarshaler

// FileEvaluatorService はファイル行数評価サービスです。
// 実処理は operations/evaluate_line_count へ委譲します。
type FileEvaluatorService struct {
	service *evaluatelinecount.Service
}

func NewFileEvaluatorService() *FileEvaluatorService {
	return &FileEvaluatorService{service: evaluatelinecount.NewService()}
}

func NewFileEvaluatorServiceWithDependencies(fileOpener FileOpener, bufioScanner BufioScanner, jsonMarshaler JSONMarshaler) *FileEvaluatorService {
	return &FileEvaluatorService{service: evaluatelinecount.NewServiceWithDependencies(fileOpener, bufioScanner, jsonMarshaler)}
}

func (e *FileEvaluatorService) CountLines(filePath string) (int, error) {
	return e.service.CountLines(filePath)
}

func (e *FileEvaluatorService) IsLineCountGreaterThan(filePath string, threshold int) (bool, int, error) {
	return e.service.IsLineCountGreaterThan(filePath, threshold)
}

func isGreaterDescription(isGreater bool) string {
	return evaluatelinecount.IsGreaterDescription(isGreater)
}

func (e *FileEvaluatorService) HandleToEvaluateLineCount(filePath string, threshold int) (string, error) {
	return e.service.HandleToEvaluateLineCount(filePath, threshold)
}

// ApiCostExtractorService はAPI料金抽出サービスです。
// 実処理は operations/parse_api_cost へ委譲します。
type ApiCostExtractorService struct {
	service *parseapicost.Service
}

func NewApiCostExtractorService() *ApiCostExtractorService {
	return &ApiCostExtractorService{service: parseapicost.NewService()}
}

func NewApiCostExtractorServiceWithFileReader(fileReader config.FileReader) *ApiCostExtractorService {
	return &ApiCostExtractorService{service: parseapicost.NewServiceWithFileReader(fileReader)}
}

func (s *ApiCostExtractorService) extractApiCostFromText(text string) (float64, error) {
	return s.service.ExtractAPICostFromText(text)
}

func (s *ApiCostExtractorService) HandleApiCostExtraction(filePath, textInput string) (float64, error) {
	return s.service.HandleApiCostExtraction(filePath, textInput)
}

// AdvancedMathService は高度な数学演算サービスです。
// 実処理は operations/power, square_root, factorial へ委譲します。
type AdvancedMathService struct{}

func NewAdvancedMathService() *AdvancedMathService {
	return &AdvancedMathService{}
}

func (a *AdvancedMathService) power(base, exponent float64) float64 {
	return power.Calculate(base, exponent)
}

func (a *AdvancedMathService) squareRoot(number float64) (float64, error) {
	return squareroot.Calculate(number)
}

func (a *AdvancedMathService) factorial(n int) (float64, error) {
	return factorial.Calculate(n)
}

func (a *AdvancedMathService) HandleToPower(base, exponent float64) (float64, error) {
	return power.HandleToPower(base, exponent)
}

func (a *AdvancedMathService) HandleToSquareRoot(number float64) (float64, error) {
	return squareroot.HandleToSquareRoot(number)
}

func (a *AdvancedMathService) HandleToFactorial(n int) (float64, error) {
	return factorial.HandleToFactorial(n)
}

// TrigonometryService は三角関数計算サービスです。
// 実処理は operations/trigonometry へ委譲します。
type TrigonometryService struct{}

func NewTrigonometryService() *TrigonometryService {
	return &TrigonometryService{}
}

func (t *TrigonometryService) trigonometry(function string, angle float64, unit string) (float64, error) {
	return trigonometry.Calculate(function, angle, unit)
}

func (t *TrigonometryService) HandleToTrigonometry(function string, angle float64, unit string) (float64, error) {
	return trigonometry.HandleToTrigonometry(function, angle, unit)
}

// MathConstantsService は数学定数提供サービスです。
// 実処理は operations/get_constants へ委譲します。
type MathConstantsService struct{}

func NewMathConstantsService() *MathConstantsService {
	return &MathConstantsService{}
}

func (m *MathConstantsService) getConstants() map[string]float64 {
	return getconstants.GetConstants()
}

func (m *MathConstantsService) HandleToGetConstants() (string, error) {
	return getconstants.HandleToGetConstants()
}

// ExpressionEvaluatorService は安全な数式評価サービスです。
// 実処理は operations/calculate_expression へ委譲します。
type ExpressionEvaluatorService struct {
	mathConstants *MathConstantsService
	service       *calculateexpression.Service
}

func NewExpressionEvaluatorService() *ExpressionEvaluatorService {
	mathConstants := NewMathConstantsService()
	return &ExpressionEvaluatorService{
		mathConstants: mathConstants,
		service:       calculateexpression.NewServiceWithConstantsGetter(mathConstants.getConstants),
	}
}

func (e *ExpressionEvaluatorService) safeEvaluate(expression string) (float64, error) {
	return e.service.SafeEvaluate(expression)
}

func (e *ExpressionEvaluatorService) evaluateBasicExpression(expression string) (float64, error) {
	return e.service.EvaluateBasicExpression(expression)
}

func (e *ExpressionEvaluatorService) evaluateArithmeticExpression(expression string) (float64, error) {
	return e.service.EvaluateArithmeticExpression(expression)
}

func (e *ExpressionEvaluatorService) checkOsPattern(expression string) error {
	return e.service.CheckOSPattern(expression)
}

func (e *ExpressionEvaluatorService) getAllIndices(text, pattern string) []int {
	return e.service.GetAllIndices(text, pattern)
}

func (e *ExpressionEvaluatorService) HandleToCalculateExpression(expression string) (float64, error) {
	return e.service.HandleToCalculateExpression(expression)
}
