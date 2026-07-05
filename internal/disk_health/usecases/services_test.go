package usecases

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
	"github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"
)

const healthySmartLog = `Device Model:     ST5000LM000-2AN170
Serial Number:    WCJ925F8
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0
190 Airflow_Temperature_Cel 0x0022   070   051   040    Old_age   Always       -       30 (Min/Max 25/30)
194 Temperature_Celsius     0x0022   030   049   000    Old_age   Always       -       30 (0 12 0 0 0)
197 Current_Pending_Sector  0x0012   100   100   000    Old_age   Always       -       0
198 Offline_Uncorrectable   0x0010   100   100   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       0
`

const criticalSmartLog = `Device Model:     WDC WD50NDZM-11BCXS1
Serial Number:    WD-WX12D126617N
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  5 Reallocated_Sector_Ct   0x0033   200   200   140    Pre-fail  Always       -       0
194 Temperature_Celsius     0x0022   103   091   000    Old_age   Always       -       49
197 Current_Pending_Sector  0x0032   200   200   000    Old_age   Always       -       404
198 Offline_Uncorrectable   0x0030   100   253   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x0032   200   200   000    Old_age   Always       -       0
`

func TestService_ParseSmartReport_Normal(t *testing.T) {
	service := NewService(ServiceOptions{})

	report, err := service.ParseSmartReport(healthySmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.Model != "ST5000LM000-2AN170" {
		t.Fatalf("expected model, got %s", report.Model)
	}
	if report.SerialNumber != "WCJ925F8" {
		t.Fatalf("expected serial, got %s", report.SerialNumber)
	}
	if report.OverallHealth != "PASSED" {
		t.Fatalf("expected PASSED, got %s", report.OverallHealth)
	}
	if len(report.Attributes) != 6 {
		t.Fatalf("expected 6 attributes, got %d", len(report.Attributes))
	}
	if report.Attributes[1].RawValue != 30 {
		t.Fatalf("expected raw value 30, got %d", report.Attributes[1].RawValue)
	}
}

func TestService_ParseSmartReport_EmptyContent(t *testing.T) {
	service := NewService(ServiceOptions{})

	_, err := service.ParseSmartReport("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_AssessReport_Healthy_Normal(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(healthySmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusHealthy {
		t.Fatalf("expected healthy, got %s", assessment.Status)
	}
	if assessment.Score != 100 {
		t.Fatalf("expected score 100, got %d", assessment.Score)
	}
	if len(assessment.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(assessment.Findings))
	}
}

func TestService_AssessReport_CriticalPendingSector(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(criticalSmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
	if len(assessment.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(assessment.Findings))
	}
	if assessment.Findings[0].AttributeID != 197 {
		t.Fatalf("expected attribute 197, got %d", assessment.Findings[0].AttributeID)
	}
	compoundFinding := assessment.Findings[1]
	if compoundFinding.Severity != domain.SeverityCritical {
		t.Fatalf("expected critical, got %s", compoundFinding.Severity)
	}
	if compoundFinding.AttributeID != 5 {
		t.Fatalf("expected attribute 5, got %d", compoundFinding.AttributeID)
	}
	if compoundFinding.AttributeName != "Reallocated_Sector_Ct" {
		t.Fatalf("expected Reallocated_Sector_Ct, got %s", compoundFinding.AttributeName)
	}
	if compoundFinding.RawValue != 0 {
		t.Fatalf("expected raw value 0, got %d", compoundFinding.RawValue)
	}
	expectedMessage := "Current_Pending_Sector=404 による代替セクタへの移し替えに失敗、ドライブが自己修復できていません"
	if compoundFinding.Message != expectedMessage {
		t.Fatalf("expected %q, got %q", expectedMessage, compoundFinding.Message)
	}
}

func TestService_AssessReport_FailedOverallHealth(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "PASSED", "FAILED", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
}

func TestService_AssessReport_WarningTemperature(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "30 (0 12 0 0 0)", "55", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
}

func TestService_AssessReport_CriticalTemperature(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "30 (0 12 0 0 0)", "61", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
}

func TestService_AssessReport_CriticalOfflineUncorrectable(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "198 Offline_Uncorrectable   0x0010   100   100   000    Old_age   Offline      -       0", "198 Offline_Uncorrectable   0x0010   100   100   000    Old_age   Offline      -       3", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
}

func TestService_AssessReport_WarningReallocatedSector(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0", "5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       2", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
}

func TestService_AssessReport_WarningCRCError(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       0", "199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       7", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
}

func TestService_AssessReport_CriticalWhenFailed(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0", "5 Reallocated_Sector_Ct   0x0033   001   001   010    Pre-fail  Always       FAILING_NOW       0", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
}

func TestService_AssessReport_UnknownMissingAttributes(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport("SMART overall-health self-assessment test result: PASSED\n")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusUnknown {
		t.Fatalf("expected unknown, got %s", assessment.Status)
	}
}

func TestService_AssessSmart_TextOutput_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				if filePath != "smart.log" {
					t.Fatalf("expected smart.log, got %s", filePath)
				}
				return []byte(healthySmartLog), nil
			},
		},
	})

	result, err := service.AssessSmart("smart.log", false, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(result, "status: healthy") {
		t.Fatalf("expected healthy output, got %s", result)
	}
	if strings.Contains(result, "attributes:") {
		t.Fatalf("expected non-verbose output, got %s", result)
	}
}

func TestService_AssessSmart_JSONOutput_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				return []byte(criticalSmartLog), nil
			},
		},
	})

	result, err := service.AssessSmart("smart.log", true, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var assessment domain.Assessment
	if err := json.Unmarshal([]byte(result), &assessment); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
	if len(assessment.Findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(assessment.Findings))
	}
	if assessment.Findings[1].AttributeID != 5 {
		t.Fatalf("expected attribute 5, got %d", assessment.Findings[1].AttributeID)
	}
	if assessment.Findings[1].RawValue != 0 {
		t.Fatalf("expected raw value 0, got %d", assessment.Findings[1].RawValue)
	}
	if len(assessment.Attributes) != 0 {
		t.Fatalf("expected attributes omitted, got %d", len(assessment.Attributes))
	}
}

func TestService_AssessSmart_JSONVerboseOutput_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				return []byte(criticalSmartLog), nil
			},
		},
	})

	result, err := service.AssessSmart("smart.log", true, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var assessment domain.Assessment
	if err := json.Unmarshal([]byte(result), &assessment); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if len(assessment.Attributes) == 0 {
		t.Fatal("expected attributes in verbose json")
	}
}

func TestService_AssessSmart_TextVerboseOutput_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				return []byte(healthySmartLog), nil
			},
		},
	})

	result, err := service.AssessSmart("smart.log", false, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(result, "attributes:") {
		t.Fatalf("expected verbose attributes, got %s", result)
	}
}

func TestService_AssessSmart_ReadFileError(t *testing.T) {
	expectedErr := errors.New("read failed")
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				return nil, expectedErr
			},
		},
	})

	_, err := service.AssessSmart("smart.log", false, false)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
