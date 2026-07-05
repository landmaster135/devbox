package assesssmart

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
Rotation Rate:    5526 rpm
Sector Sizes:     512 bytes logical, 4096 bytes physical
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  1 Raw_Read_Error_Rate     0x002f   189   184   051    Pre-fail  Always       -       1061
  5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0
  9 Power_On_Hours          0x0032   100   100   000    Old_age   Always       -       87 (239 223 0)
 12 Power_Cycle_Count       0x0032   100   100   020    Old_age   Always       -       159
190 Airflow_Temperature_Cel 0x0022   070   051   040    Old_age   Always       -       30 (Min/Max 25/30)
194 Temperature_Celsius     0x0022   030   049   000    Old_age   Always       -       30 (0 12 0 0 0)
197 Current_Pending_Sector  0x0012   100   100   000    Old_age   Always       -       0
198 Offline_Uncorrectable   0x0010   100   100   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       0
241 Total_LBAs_Written      0x0000   100   253   000    Old_age   Offline      -       3836222920
242 Total_LBAs_Read         0x0000   100   253   000    Old_age   Offline      -       1647232405
`

const criticalSmartLog = `Device Model:     WDC WD50NDZM-11BCXS1
Serial Number:    WD-WX12D126617N
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  1 Raw_Read_Error_Rate     0x002f   051   050   051    Pre-fail  Always       -       1061
  5 Reallocated_Sector_Ct   0x0033   200   200   140    Pre-fail  Always       -       0
194 Temperature_Celsius     0x0022   103   091   000    Old_age   Always       -       61
197 Current_Pending_Sector  0x0032   200   200   000    Old_age   Always       -       404
198 Offline_Uncorrectable   0x0030   100   253   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x0032   200   200   000    Old_age   Always       -       0
`

func TestService_ParseSmartReport_Normal(t *testing.T) {
	service := NewService(nil)

	report, err := service.ParseSmartReport(healthySmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.Model != "ST5000LM000-2AN170" {
		t.Fatalf("expected model, got %s", report.Model)
	}
	if len(report.Attributes) != 11 {
		t.Fatalf("expected 11 attributes, got %d", len(report.Attributes))
	}
	assertInt64Pointer(t, report.DiskInfo.RotationRateRPM, 5526)
	assertInt64Pointer(t, report.DiskInfo.TotalBytesWritten, 1964146135040)
}

func TestService_ParseSmartReport_EmptyContent(t *testing.T) {
	service := NewService(nil)

	_, err := service.ParseSmartReport("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_AssessReport_Critical_Normal(t *testing.T) {
	service := NewService(nil)
	report, err := service.ParseSmartReport(criticalSmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)
	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
	if len(assessment.Findings) == 0 {
		t.Fatal("expected findings")
	}
}

func TestService_AssessReport_WarningBranches_Normal(t *testing.T) {
	service := NewService(nil)
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       0", "199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       7", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)
	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
}

func TestService_AssessReport_UnknownMissingAttributes_Normal(t *testing.T) {
	service := NewService(nil)
	report, err := service.ParseSmartReport("SMART overall-health self-assessment test result: PASSED\n")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)
	if assessment.Status != domain.StatusUnknown {
		t.Fatalf("expected unknown, got %s", assessment.Status)
	}
}

func TestService_FormatText_WithVerboseAttributes_Normal(t *testing.T) {
	service := NewService(nil)
	report, err := service.ParseSmartReport(healthySmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assessment := service.AssessReport(report)

	result := service.FormatText(assessment, true)
	expectedFragments := []string{
		"disk_info:",
		"  rotation_rate_rpm: 5526",
		"attributes:",
		"id=1 name=Raw_Read_Error_Rate",
	}
	for _, fragment := range expectedFragments {
		if !strings.Contains(result, fragment) {
			t.Fatalf("expected %q in result, got %s", fragment, result)
		}
	}
}

func TestService_FormatJSON_OmitsAttributesWhenNonVerbose_Normal(t *testing.T) {
	service := NewService(nil)
	report, err := service.ParseSmartReport(healthySmartLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	assessment := service.AssessReport(report)

	result, err := service.FormatJSON(assessment, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded domain.Assessment
	if err := json.Unmarshal([]byte(result), &decoded); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if len(decoded.Attributes) != 0 {
		t.Fatalf("expected attributes omitted, got %d", len(decoded.Attributes))
	}
}

func TestService_Execute_TextOutput_Normal(t *testing.T) {
	service := NewService(&filesystem.MockRepository{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			if filePath != "smart.log" {
				t.Fatalf("expected smart.log, got %s", filePath)
			}
			return []byte(healthySmartLog), nil
		},
	})

	result, err := service.Execute("smart.log", false, false)
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

func TestService_Execute_JSONOutput_Normal(t *testing.T) {
	service := NewService(&filesystem.MockRepository{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			return []byte(criticalSmartLog), nil
		},
	})

	result, err := service.Execute("smart.log", true, false)
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
	if len(assessment.Attributes) != 0 {
		t.Fatalf("expected attributes omitted, got %d", len(assessment.Attributes))
	}
}

func TestService_Execute_ReadFileError(t *testing.T) {
	expectedErr := errors.New("read failed")
	service := NewService(&filesystem.MockRepository{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			return nil, expectedErr
		},
	})

	_, err := service.Execute("smart.log", false, false)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}

func assertInt64Pointer(t *testing.T, actual *int64, expected int64) {
	t.Helper()
	if actual == nil {
		t.Fatalf("expected %d, got nil", expected)
	}
	if *actual != expected {
		t.Fatalf("expected %d, got %d", expected, *actual)
	}
}
