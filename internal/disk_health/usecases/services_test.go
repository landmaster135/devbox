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
  1 Raw_Read_Error_Rate     0x002f   189   184   051    Pre-fail  Always       -       1061
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
	if len(report.Attributes) != 11 {
		t.Fatalf("expected 11 attributes, got %d", len(report.Attributes))
	}
	if report.Attributes[4].RawValue != 30 {
		t.Fatalf("expected raw value 30, got %d", report.Attributes[4].RawValue)
	}
	if report.DiskInfo == nil {
		t.Fatal("expected disk info")
	}
	assertInt64Pointer(t, report.DiskInfo.RotationRateRPM, 5526)
	assertInt64Pointer(t, report.DiskInfo.PowerOnHours, 87)
	assertInt64Pointer(t, report.DiskInfo.PowerCycleCount, 159)
	assertInt64Pointer(t, report.DiskInfo.TemperatureCelsius, 30)
	assertInt64Pointer(t, report.DiskInfo.TotalLBAsWritten, 3836222920)
	assertInt64Pointer(t, report.DiskInfo.TotalBytesWritten, 1964146135040)
	assertInt64Pointer(t, report.DiskInfo.TotalLBAsRead, 1647232405)
	assertInt64Pointer(t, report.DiskInfo.TotalBytesRead, 843382991360)
}

func TestService_ParseSmartReport_MissingDiskInfo_Normal(t *testing.T) {
	service := NewService(ServiceOptions{})

	report, err := service.ParseSmartReport(`Device Model:     ST5000LM000-2AN170
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0
`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.DiskInfo != nil {
		t.Fatalf("expected no disk info, got %#v", report.DiskInfo)
	}
}

func TestService_ParseSmartReport_DefaultSectorAndAirflowTemperature_Normal(t *testing.T) {
	service := NewService(ServiceOptions{})

	report, err := service.ParseSmartReport(`Device Model:     ST5000LM000-2AN170
Rotation Rate:    Solid State Device
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
190 Airflow_Temperature_Cel 0x0022   070   051   040    Old_age   Always       -       31 (Min/Max 25/31)
241 Total_LBAs_Written      0x0000   100   253   000    Old_age   Offline      -       10
242 Total_LBAs_Read         0x0000   100   253   000    Old_age   Offline      -       20
`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.DiskInfo == nil {
		t.Fatal("expected disk info")
	}
	if report.DiskInfo.RotationRateRPM != nil {
		t.Fatalf("expected no rotation rate, got %d", *report.DiskInfo.RotationRateRPM)
	}
	assertInt64Pointer(t, report.DiskInfo.TemperatureCelsius, 31)
	assertInt64Pointer(t, report.DiskInfo.TotalBytesWritten, 5120)
	assertInt64Pointer(t, report.DiskInfo.TotalBytesRead, 10240)
}

func TestService_ParseSmartReport_SectorSizeFallback_Normal(t *testing.T) {
	service := NewService(ServiceOptions{})

	report, err := service.ParseSmartReport(`Device Model:     ST5000LM000-2AN170
Sector Size:      4096 bytes
SMART overall-health self-assessment test result: PASSED

Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
241 Total_LBAs_Written      0x0000   100   253   000    Old_age   Offline      -       10
`)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if report.DiskInfo == nil {
		t.Fatal("expected disk info")
	}
	assertInt64Pointer(t, report.DiskInfo.TotalBytesWritten, 40960)
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

func TestService_AssessReport_WarningRawReadErrorRateValue(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "1 Raw_Read_Error_Rate     0x002f   189   184   051", "1 Raw_Read_Error_Rate     0x002f   061   184   051", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
	if len(assessment.Findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(assessment.Findings))
	}
	finding := assessment.Findings[0]
	if finding.Severity != domain.SeverityWarning {
		t.Fatalf("expected warning, got %s", finding.Severity)
	}
	if finding.AttributeID != 1 {
		t.Fatalf("expected attribute 1, got %d", finding.AttributeID)
	}
	expectedMessage := "Raw_Read_Error_Rate の現在正規化値がメーカー定義しきい値に近づいています: VALUE=61 THRESH=51"
	if finding.Message != expectedMessage {
		t.Fatalf("expected %q, got %q", expectedMessage, finding.Message)
	}
}

func TestService_AssessReport_CriticalRawReadErrorRateValue(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "1 Raw_Read_Error_Rate     0x002f   189   184   051", "1 Raw_Read_Error_Rate     0x002f   051   184   051", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
	finding := assessment.Findings[0]
	if finding.Severity != domain.SeverityCritical {
		t.Fatalf("expected critical, got %s", finding.Severity)
	}
	expectedMessage := "Raw_Read_Error_Rate の現在正規化値がメーカー定義しきい値以下です: VALUE=51 THRESH=51"
	if finding.Message != expectedMessage {
		t.Fatalf("expected %q, got %q", expectedMessage, finding.Message)
	}
}

func TestService_AssessReport_WarningRawReadErrorRateWorst(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "1 Raw_Read_Error_Rate     0x002f   189   184   051", "1 Raw_Read_Error_Rate     0x002f   189   060   051", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
	finding := assessment.Findings[0]
	if finding.Severity != domain.SeverityWarning {
		t.Fatalf("expected warning, got %s", finding.Severity)
	}
	expectedMessage := "Raw_Read_Error_Rate の過去最悪正規化値がメーカー定義しきい値に近づいています: WORST=60 THRESH=51"
	if finding.Message != expectedMessage {
		t.Fatalf("expected %q, got %q", expectedMessage, finding.Message)
	}
}

func TestService_AssessReport_CriticalRawReadErrorRateWorst(t *testing.T) {
	service := NewService(ServiceOptions{})
	report, err := service.ParseSmartReport(strings.Replace(healthySmartLog, "1 Raw_Read_Error_Rate     0x002f   189   184   051", "1 Raw_Read_Error_Rate     0x002f   189   050   051", 1))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessReport(report)

	if assessment.Status != domain.StatusCritical {
		t.Fatalf("expected critical, got %s", assessment.Status)
	}
	finding := assessment.Findings[0]
	if finding.Severity != domain.SeverityCritical {
		t.Fatalf("expected critical, got %s", finding.Severity)
	}
	expectedMessage := "Raw_Read_Error_Rate の過去最悪正規化値がメーカー定義しきい値以下です: WORST=50 THRESH=51"
	if finding.Message != expectedMessage {
		t.Fatalf("expected %q, got %q", expectedMessage, finding.Message)
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
	if !strings.Contains(result, "disk_info:") {
		t.Fatalf("expected disk info output, got %s", result)
	}
	if !strings.Contains(result, "  rotation_rate_rpm: 5526") {
		t.Fatalf("expected rotation rate output, got %s", result)
	}
	if !strings.Contains(result, "  total_bytes_written: 1964146135040") {
		t.Fatalf("expected total bytes written output, got %s", result)
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

func TestService_AssessSmart_JSONDiskInfoOutput_Normal(t *testing.T) {
	service := NewService(ServiceOptions{
		FileSystem: &filesystem.MockRepository{
			ReadFileFunc: func(filePath string) ([]byte, error) {
				return []byte(healthySmartLog), nil
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
	if assessment.DiskInfo == nil {
		t.Fatal("expected disk info")
	}
	assertInt64Pointer(t, assessment.DiskInfo.RotationRateRPM, 5526)
	assertInt64Pointer(t, assessment.DiskInfo.PowerOnHours, 87)
	assertInt64Pointer(t, assessment.DiskInfo.PowerCycleCount, 159)
	assertInt64Pointer(t, assessment.DiskInfo.TemperatureCelsius, 30)
	assertInt64Pointer(t, assessment.DiskInfo.TotalLBAsWritten, 3836222920)
	assertInt64Pointer(t, assessment.DiskInfo.TotalBytesWritten, 1964146135040)
	assertInt64Pointer(t, assessment.DiskInfo.TotalLBAsRead, 1647232405)
	assertInt64Pointer(t, assessment.DiskInfo.TotalBytesRead, 843382991360)
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

func assertInt64Pointer(t *testing.T, actual *int64, expected int64) {
	t.Helper()

	if actual == nil {
		t.Fatalf("expected %d, got nil", expected)
	}
	if *actual != expected {
		t.Fatalf("expected %d, got %d", expected, *actual)
	}
}
