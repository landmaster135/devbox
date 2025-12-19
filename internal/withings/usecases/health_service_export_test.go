package usecases

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExportDailySummaryWritesFile(t *testing.T) {
	svc := NewHealthService(0)
	weight := 68.5
	resp := &DailySummaryResponse{
		Timezone: "UTC",
		Summaries: []DailySummary{
			{
				Date:     "2024-10-01",
				Measures: &DailySummaryMeasures{WeightKg: &weight},
			},
		},
	}

	tmp := t.TempDir()
	path := filepath.Join(tmp, "exports", "summary.json")

	if err := svc.ExportDailySummary(resp, path); err != nil {
		t.Fatalf("unexpected export error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read exported file: %v", err)
	}

	var parsed struct {
		Data struct {
			HealthMates []FlattenedDailySummary `json:"health_mates"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}
	if len(parsed.Data.HealthMates) != 1 {
		t.Fatalf("expected one entry, got %d", len(parsed.Data.HealthMates))
	}
}

func TestExportDailySummaryValidatesInput(t *testing.T) {
	svc := NewHealthService(0)
	if err := svc.ExportDailySummary(nil, "out.json"); err == nil {
		t.Fatalf("expected error for nil response")
	}
	resp := &DailySummaryResponse{Summaries: []DailySummary{{}}}
	if err := svc.ExportDailySummary(resp, " "); err == nil {
		t.Fatalf("expected error for blank path")
	}
}
