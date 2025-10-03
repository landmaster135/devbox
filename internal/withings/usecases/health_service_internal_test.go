package usecases

import (
	"errors"
	"testing"
)

func TestTruncateForLog(t *testing.T) {
	short := []byte("hello")
	if got := truncateForLog(short); got != "hello" {
		t.Fatalf("unexpected short result: %s", got)
	}

	long := make([]byte, 600)
	for i := range long {
		long[i] = 'a'
	}
	got := truncateForLog(long)
	if len(got) != 515 || got[len(got)-3:] != "..." {
		t.Fatalf("unexpected long result: %q", got[len(got)-10:])
	}
}

func TestSetBaseURL(t *testing.T) {
	svc := NewHealthService(0)
	original := svc.baseURL
	svc.SetBaseURL(" https://custom.example/api/ ")
	if svc.baseURL != "https://custom.example/api" {
		t.Fatalf("baseURL not trimmed: %s", svc.baseURL)
	}
	svc.SetBaseURL("")
	if svc.baseURL != "https://custom.example/api" {
		t.Fatalf("baseURL should remain unchanged when empty")
	}

	// ensure default was set initially
	if original == "" {
		t.Fatalf("default baseURL should not be empty")
	}
}

func TestHealthServiceShouldRetryDailySummaryWithRefresh(t *testing.T) {
	svc := NewHealthService(0)

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil", err: nil, want: false},
		{name: "401", err: errors.New("status=401 unauthorized"), want: true},
		{name: "403", err: errors.New("Status=403"), want: true},
		{name: "unauthorized", err: errors.New("UNAUTHORIZED"), want: true},
		{name: "invalid token", err: errors.New("INVALID_TOKEN"), want: true},
		{name: "other", err: errors.New("timeout"), want: false},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			got := svc.ShouldRetryDailySummaryWithRefresh(tc.err)
			if got != tc.want {
				t.Fatalf("ShouldRetryDailySummaryWithRefresh(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

func TestLabelForMeasureType(t *testing.T) {
	if got := labelForMeasureType(1); got != "weight_kg" {
		t.Fatalf("unexpected known label: %s", got)
	}
	if got := labelForMeasureType(999); got != "measure_999" {
		t.Fatalf("unexpected fallback label: %s", got)
	}
}

func TestMergeActivitySummaries(t *testing.T) {
	base := ActivitySummary{
		Steps:           intPtr(1000),
		CaloriesKcal:    floatPtr(200),
		ModerateSeconds: intPtr(600),
	}
	incoming := ActivitySummary{
		Steps:             intPtr(2500),
		DistanceMeter:     floatPtr(1500),
		ElevationMeter:    floatPtr(12.5),
		CaloriesKcal:      floatPtr(350),
		TotalCaloriesKcal: floatPtr(500),
		SoftSeconds:       intPtr(400),
		IntenseSeconds:    intPtr(50),
		ActiveSeconds:     intPtr(1200),
		HrAverageBPM:      floatPtr(62.5),
		HrMinBPM:          floatPtr(45),
		HrMaxBPM:          floatPtr(140),
		DeviceBrand:       intPtr(7),
		DeviceModelID:     intPtr(1234),
		DeviceModelName:   stringPtr("Withings Move"),
		IsTracker:         boolPtr(true),
	}

	merged := mergeActivitySummaries(base, incoming)

	if merged.Steps == nil || *merged.Steps != 2500 {
		t.Fatalf("incoming steps should override base")
	}
	if merged.DistanceMeter == nil || *merged.DistanceMeter != 1500 {
		t.Fatalf("incoming distance should be set")
	}
	if merged.ElevationMeter == nil || *merged.ElevationMeter != 12.5 {
		t.Fatalf("incoming elevation should be set")
	}
	if merged.CaloriesKcal == nil || *merged.CaloriesKcal != 350 {
		t.Fatalf("calories should update from incoming")
	}
	if merged.TotalCaloriesKcal == nil || *merged.TotalCaloriesKcal != 500 {
		t.Fatalf("total calories should update from incoming")
	}
	if merged.SoftSeconds == nil || *merged.SoftSeconds != 400 {
		t.Fatalf("soft seconds should update")
	}
	if merged.ModerateSeconds == nil || *merged.ModerateSeconds != 600 {
		t.Fatalf("moderate seconds should remain from base when incoming nil")
	}
	if merged.IntenseSeconds == nil || *merged.IntenseSeconds != 50 {
		t.Fatalf("intense seconds should update")
	}
	if merged.ActiveSeconds == nil || *merged.ActiveSeconds != 1200 {
		t.Fatalf("active seconds should update")
	}
	if merged.HrAverageBPM == nil || *merged.HrAverageBPM != 62.5 {
		t.Fatalf("hr average should update")
	}
	if merged.HrMinBPM == nil || *merged.HrMinBPM != 45 {
		t.Fatalf("hr min should update")
	}
	if merged.HrMaxBPM == nil || *merged.HrMaxBPM != 140 {
		t.Fatalf("hr max should update")
	}
	if merged.DeviceBrand == nil || *merged.DeviceBrand != 7 {
		t.Fatalf("device brand should update")
	}
	if merged.DeviceModelID == nil || *merged.DeviceModelID != 1234 {
		t.Fatalf("device model id should update")
	}
	if merged.DeviceModelName == nil || *merged.DeviceModelName != "Withings Move" {
		t.Fatalf("device model name should update")
	}
	if merged.IsTracker == nil || !*merged.IsTracker {
		t.Fatalf("is_tracker should update")
	}
}

func TestFlexibleBoolUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		data  string
		want  bool
		valid bool
	}{
		{name: "true", data: "true", want: true, valid: true},
		{name: "false", data: "false", want: false, valid: true},
		{name: "one", data: "1", want: true, valid: true},
		{name: "zero", data: "0", want: false, valid: true},
		{name: "whitespace number", data: " 5 ", want: true, valid: true},
		{name: "empty", data: "", want: false, valid: true},
		{name: "invalid", data: "\"yes\"", valid: false},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			var fb flexibleBool
			err := fb.UnmarshalJSON([]byte(tc.data))
			if !tc.valid {
				if err == nil {
					t.Fatalf("expected error for %s", tc.data)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fb.Bool() != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, fb.Bool())
			}
		})
	}
}
