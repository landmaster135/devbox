package usecases

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAuthServiceExportDailySummary(t *testing.T) {
	svc := NewAuthService(2 * time.Second)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "withings-summary.json")

	steps := 1234
	isTracker := true
	hrAvg := 65.5

	resp := &DailySummaryResponse{
		Summaries: []DailySummary{
			{
				Date:     "2025-10-01",
				Timezone: "Europe/Paris",
				Measures: map[string]float64{
					"bone_mass_kg":     2.6,
					"fat_free_mass_kg": 50.613,
				},
				Activity: &ActivitySummary{
					Steps:        &steps,
					HrAverageBPM: &hrAvg,
					IsTracker:    &isTracker,
				},
			},
		},
	}

	if err := svc.ExportDailySummary(resp, outputPath); err != nil {
		t.Fatalf("ExportDailySummary returned error: %v", err)
	}

	raw, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	var payload struct {
		Data struct {
			HealthMates []map[string]any `json:"health_mates"`
		} `json:"data"`
		Description string `json:"description"`
		Name        string `json:"name"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("failed to unmarshal output json: %v", err)
	}

	if payload.Description != "Health Mate data from Withings" {
		t.Fatalf("unexpected description: %s", payload.Description)
	}
	if payload.Name != "My Health Mate Data" {
		t.Fatalf("unexpected name: %s", payload.Name)
	}
	if len(payload.Data.HealthMates) != 1 {
		t.Fatalf("unexpected health mate entry count: %d", len(payload.Data.HealthMates))
	}

	entry := payload.Data.HealthMates[0]
	assertFloat(t, entry["measures_bone_mass_kg"], 2.6)
	assertFloat(t, entry["measures_fat_free_mass_kg"], 50.613)
	if stepsVal, ok := entry["activity_steps"].(float64); !ok || int(stepsVal) != steps {
		t.Fatalf("unexpected activity_steps value: %#v", entry["activity_steps"])
	}
	if avg, ok := entry["activity_hr_average_bpm"].(float64); !ok || avg != hrAvg {
		t.Fatalf("unexpected activity_hr_average_bpm: %#v", entry["activity_hr_average_bpm"])
	}
	if tracker, ok := entry["activity_is_tracker"].(bool); !ok || tracker != isTracker {
		t.Fatalf("unexpected activity_is_tracker: %#v", entry["activity_is_tracker"])
	}
	if measuredAt, ok := entry["measured_at"].(string); !ok || measuredAt != "2025-10-01" {
		t.Fatalf("unexpected measured_at: %#v", entry["measured_at"])
	}
	if _, ok := entry["timezone"]; ok {
		t.Fatalf("timezone should not be present: %#v", entry["timezone"])
	}
}

func assertFloat(t *testing.T, value any, expected float64) {
	t.Helper()
	v, ok := value.(float64)
	if !ok {
		t.Fatalf("value is not float64: %#v", value)
	}
	if v != expected {
		t.Fatalf("unexpected float value: got %v want %v", v, expected)
	}
}
