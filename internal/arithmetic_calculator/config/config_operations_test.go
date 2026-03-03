package config

import "testing"

func TestSupportedOperations_Normal(t *testing.T) {
	operations := SupportedOperations()
	if len(operations) == 0 {
		t.Fatal("supported operations should not be empty")
	}

	seen := map[string]struct{}{}
	for _, operation := range operations {
		if operation == "" {
			t.Fatal("operation should not be empty")
		}
		if _, ok := seen[operation]; ok {
			t.Fatalf("duplicated operation found: %s", operation)
		}
		seen[operation] = struct{}{}
	}
}

func TestNewConfig_AllSupportedOperations_Normal(t *testing.T) {
	for _, operation := range SupportedOperations() {
		t.Run(operation, func(t *testing.T) {
			cfg, err := buildConfigForOperation(operation)
			if err != nil {
				t.Fatalf("NewConfig returned error: %v", err)
			}
			if cfg.Operation != operation {
				t.Fatalf("unexpected operation: got=%s want=%s", cfg.Operation, operation)
			}
		})
	}
}

func buildConfigForOperation(operation string) (*Config, error) {
	switch operation {
	case OperationSum:
		return NewConfig(operation, 0, 0, []float64{1, 2}, "", 0, 0, 0, 0, 0, 0, "", "", "", "")
	case OperationEvaluateLineCount:
		return NewConfig(operation, 0, 0, nil, "/tmp/test.txt", 1, 0, 0, 0, 0, 0, "", "", "", "")
	case OperationParseAPICost:
		return NewConfigForParseApiCost(operation, "", "API料金が100円掛かった")
	case OperationPower:
		return NewConfig(operation, 0, 0, nil, "", 0, 2, 8, 0, 0, 0, "", "", "", "")
	case OperationSquareRoot:
		return NewConfig(operation, 0, 0, nil, "", 0, 0, 0, 16, 0, 0, "", "", "", "")
	case OperationFactorial:
		return NewConfig(operation, 0, 0, nil, "", 0, 0, 0, 0, 0, 5, "", "", "", "")
	case OperationTrigonometry:
		return NewConfig(operation, 0, 0, nil, "", 0, 0, 0, 0, 90, 0, "sin", "degrees", "", "")
	case OperationCalculate:
		return NewConfig(operation, 0, 0, nil, "", 0, 0, 0, 0, 0, 0, "", "", "2+3", "")
	case OperationGetConstants:
		return NewConfig(operation, 0, 0, nil, "", 0, 0, 0, 0, 0, 0, "", "", "", "")
	default:
		return NewConfig(operation, 1, 2, nil, "", 0, 0, 0, 0, 0, 0, "", "", "", "")
	}
}
