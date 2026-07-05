package assessdmesg

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/disk_health/domain"
	"github.com/landmaster135/devbox/internal/disk_health/infrastructures/filesystem"
)

const healthyDmesgLog = `(標準入力):379:[    0.198970] ACPI BIOS Error (bug): Failure creating named object
(標準入力):806:[    0.527869] Adaptive Deadline I/O Scheduler 3.2.0
(標準入力):1478:[ 6230.321645] sd 2:0:0:0: [sdc] Attached SCSI disk
`

const criticalDmesgLog = `(標準入力):1:[17568.682761] critical medium error, dev sdb, sector 8491832904 op 0x0:(READ)
(標準入力):2:[17571.754409] sd 0:0:0:0: [sdb] tag#0 FAILED Result: hostbyte=DID_OK driverbyte=DRIVER_OK
(標準入力):3:[17571.754422] sd 0:0:0:0: [sdb] tag#0 Sense Key : Medium Error [current]
(標準入力):4:[17571.754426] sd 0:0:0:0: [sdb] tag#0 Add. Sense: Unrecovered read error
`

const warningDmesgLog = `[18000.100000] blk_update_request: I/O error, dev sdb, sector 42 op 0x0:(READ)
[18001.100000] sd 0:0:0:0: [sdb] tag#1 failed command: READ FPDMA QUEUED
[18002.100000] ata1: link is slow to respond, please be patient
`

func TestService_ParseDmesgLog_Normal(t *testing.T) {
	service := NewService(nil)

	events, err := service.ParseDmesgLog(criticalDmesgLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(events) != 4 {
		t.Fatalf("expected 4 events, got %d", len(events))
	}
	if events[0].Severity != domain.SeverityCritical {
		t.Fatalf("expected critical, got %s", events[0].Severity)
	}
}

func TestService_ParseDmesgLog_HealthyNoiseOnly_Normal(t *testing.T) {
	service := NewService(nil)

	events, err := service.ParseDmesgLog(healthyDmesgLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no events, got %d", len(events))
	}
}

func TestService_ParseDmesgLog_EmptyContent(t *testing.T) {
	service := NewService(nil)

	_, err := service.ParseDmesgLog("")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestService_AssessDmesgEvents_Warning_Normal(t *testing.T) {
	service := NewService(nil)
	events, err := service.ParseDmesgLog(warningDmesgLog)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	assessment := service.AssessDmesgEvents(events)
	if assessment.Status != domain.StatusWarning {
		t.Fatalf("expected warning, got %s", assessment.Status)
	}
}

func TestService_FormatJSON_OmitsLineWhenNonVerbose_Normal(t *testing.T) {
	service := NewService(nil)
	assessment := domain.DmesgAssessment{
		Status: domain.StatusCritical,
		Findings: []domain.DmesgEvent{
			{Severity: domain.SeverityCritical, Line: "[1.0] raw line"},
		},
	}

	result, err := service.FormatJSON(assessment, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var decoded domain.DmesgAssessment
	if err := json.Unmarshal([]byte(result), &decoded); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if decoded.Findings[0].Line != "" {
		t.Fatalf("expected line omitted, got %s", decoded.Findings[0].Line)
	}
}

func TestService_Execute_TextOutput_Normal(t *testing.T) {
	service := NewService(&filesystem.MockRepository{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			if filePath != "dmesg.log" {
				t.Fatalf("expected dmesg.log, got %s", filePath)
			}
			return []byte(criticalDmesgLog), nil
		},
	})

	result, err := service.Execute("dmesg.log", false, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(result, "status: critical") {
		t.Fatalf("expected critical output, got %s", result)
	}
	if !strings.Contains(result, "raw:") {
		t.Fatalf("expected verbose raw output, got %s", result)
	}
}

func TestService_Execute_ReadFileError(t *testing.T) {
	expectedErr := errors.New("read failed")
	service := NewService(&filesystem.MockRepository{
		ReadFileFunc: func(filePath string) ([]byte, error) {
			return nil, expectedErr
		},
	})

	_, err := service.Execute("dmesg.log", false, false)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}
