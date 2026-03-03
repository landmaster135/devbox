package usecases

// Service はoperationごとの実行サービスを集約する
// cmd層からは ExecuteByConfig のみを利用する。
type Service struct {
	basicCalculationOperation basicCalculationOperation
	arraySumOperation         arraySumOperation
	lineCountOperation        lineCountOperation
	apiCostOperation          apiCostOperation
	powerOperation            powerOperation
	squareRootOperation       squareRootOperation
	factorialOperation        factorialOperation
	trigonometryOperation     trigonometryOperation
	calculateOperation        calculateOperation
	constantsOperation        constantsOperation
}

type basicCalculationOperation interface {
	Execute(operation string, x, y float64) (string, error)
}

type arraySumOperation interface {
	Execute(operation string, numbers []float64) (string, error)
}

type lineCountOperation interface {
	Execute(filePath string, threshold int) (string, error)
}

type apiCostOperation interface {
	Execute(filePath, textInput string) (string, error)
}

type powerOperation interface {
	Execute(base, exponent float64) (string, error)
}

type squareRootOperation interface {
	Execute(number float64) (string, error)
}

type factorialOperation interface {
	Execute(n int) (string, error)
}

type trigonometryOperation interface {
	Execute(function string, angle float64, unit string) (string, error)
}

type calculateOperation interface {
	Execute(expression string) (string, error)
}

type constantsOperation interface {
	Execute() (string, error)
}

func newServiceWithOperations(
	basicCalculationOp basicCalculationOperation,
	arraySumOp arraySumOperation,
	lineCountOp lineCountOperation,
	apiCostOp apiCostOperation,
	powerOp powerOperation,
	squareRootOp squareRootOperation,
	factorialOp factorialOperation,
	trigonometryOp trigonometryOperation,
	calculateOp calculateOperation,
	constantsOp constantsOperation,
) *Service {
	return &Service{
		basicCalculationOperation: basicCalculationOp,
		arraySumOperation:         arraySumOp,
		lineCountOperation:        lineCountOp,
		apiCostOperation:          apiCostOp,
		powerOperation:            powerOp,
		squareRootOperation:       squareRootOp,
		factorialOperation:        factorialOp,
		trigonometryOperation:     trigonometryOp,
		calculateOperation:        calculateOp,
		constantsOperation:        constantsOp,
	}
}
