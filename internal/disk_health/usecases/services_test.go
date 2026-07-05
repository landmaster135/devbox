package usecases

import (
	"strings"
	"testing"

	config "github.com/landmaster135/devbox/internal/disk_health/config"
	"github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"
)

func TestService_ExecuteByConfig_Normal(t *testing.T) {
	tests := []struct {
		name           string
		cfg            *config.Config
		expectedResult string
		expectedSmart  bool
		expectedDmesg  bool
	}{
		{
			name: "AssessSmart",
			cfg: &config.Config{
				Operation: config.OperationAssessSmart,
				SrcFile:   "smart.log",
				JSON:      true,
				Verbose:   true,
			},
			expectedResult: "smart result",
			expectedSmart:  true,
		},
		{
			name: "AssessDmesg",
			cfg: &config.Config{
				Operation: config.OperationAssessDmesg,
				SrcFile:   "dmesg.log",
			},
			expectedResult: "dmesg result",
			expectedDmesg:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smartOp := &mockAssessSmartOperation{result: "smart result"}
			dmesgOp := &mockAssessDmesgOperation{result: "dmesg result"}
			service := newServiceWithOperations(nil, smartOp, dmesgOp)

			result, err := service.ExecuteByConfig(tt.cfg)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if result != tt.expectedResult {
				t.Fatalf("expected %s, got %s", tt.expectedResult, result)
			}
			if smartOp.called != tt.expectedSmart {
				t.Fatalf("expected smart called %v, got %v", tt.expectedSmart, smartOp.called)
			}
			if dmesgOp.called != tt.expectedDmesg {
				t.Fatalf("expected dmesg called %v, got %v", tt.expectedDmesg, dmesgOp.called)
			}
		})
	}
}

func TestService_ExecuteByConfig_UnsupportedOperation(t *testing.T) {
	service := newServiceWithOperations(nil, &mockAssessSmartOperation{}, &mockAssessDmesgOperation{})

	_, err := service.ExecuteByConfig(&config.Config{Operation: "unknown"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "未対応のoperationです: unknown") {
		t.Fatalf("expected unsupported operation error, got %v", err)
	}
}

func TestService_NewService_AssessSmartDelegation_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				return []byte("SMART overall-health self-assessment test result: PASSED\n"), nil
			},
		},
	})

	result, err := service.AssessSmart("smart.log", false, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(result, "status: unknown") {
		t.Fatalf("expected unknown output, got %s", result)
	}
}

type mockAssessSmartOperation struct {
	result string
	called bool
}

func (m *mockAssessSmartOperation) Execute(srcFile string, outputJSON bool, verbose bool) (string, error) {
	m.called = true
	return m.result, nil
}

type mockAssessDmesgOperation struct {
	result string
	called bool
}

func (m *mockAssessDmesgOperation) Execute(srcFile string, outputJSON bool, verbose bool) (string, error) {
	m.called = true
	return m.result, nil
}
