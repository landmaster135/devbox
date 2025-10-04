package usecases

import (
	"encoding/json"
	"testing"
)

func TestFlattenDailySummaryResponse(t *testing.T) {
	weight := 70.1234
	hydration := 55.9876
	steps := 8500
	distance := 6543.219
	avg := 59.444
	min := 42.666
	max := 123.999
	date := "2024-10-01"

	resp := &DailySummaryResponse{
		Timezone: "Europe/Paris",
		Summaries: []DailySummary{
			{
				Date: date,
				Measures: &DailySummaryMeasures{
					WeightKg:     &weight,
					HydrationKg:  &hydration,
				},
				Activity: &ActivitySummary{
					Steps:             &steps,
					DistanceMeter:     &distance,
					HrAverageBPM:      &avg,
					HrMinBPM:          &min,
					HrMaxBPM:          &max,
				},
			},
		},
	}

	flattened := FlattenDailySummaryResponse(resp)
	if flattened.Timezone != "Europe/Paris" {
		t.Fatalf("unexpected timezone: %s", flattened.Timezone)
	}
	if len(flattened.Summaries) != 1 {
		t.Fatalf("expected one summary, got %d", len(flattened.Summaries))
	}
	item := flattened.Summaries[0]
	if item.MeasuredAt != date {
		t.Fatalf("unexpected measured_at: %s", item.MeasuredAt)
	}
	if item.FlattenedMeasures == nil || item.MeasuresWeightKg == nil {
		t.Fatalf("flattened measures missing weight: %+v", item.FlattenedMeasures)
	}
	if *item.MeasuresWeightKg != 70.1234 {
		t.Fatalf("unexpected weight: %f", *item.MeasuresWeightKg)
	}
	if item.MeasuresHydrationKg == nil || *item.MeasuresHydrationKg != 55.9876 {
		t.Fatalf("unexpected hydration: %+v", item.MeasuresHydrationKg)
	}
	if item.FlattenedActivity == nil || item.ActivitySteps == nil || *item.ActivitySteps != 8500 {
		t.Fatalf("flattened activity missing steps: %+v", item.FlattenedActivity)
	}
	if item.ActivityDistanceMeter == nil || *item.ActivityDistanceMeter != 6543.219 {
		t.Fatalf("unexpected distance: %+v", item.ActivityDistanceMeter)
	}
	if item.ActivityHrAverageBPM == nil || *item.ActivityHrAverageBPM != 59.444 {
		t.Fatalf("unexpected hr average: %+v", item.ActivityHrAverageBPM)
	}
}

func TestFlattenHelpersHandleNil(t *testing.T) {
	resp := FlattenDailySummaryResponse(nil)
	if resp.Summaries != nil {
		t.Fatalf("expected nil summaries, got %v", resp.Summaries)
	}
	if resp.Timezone != "" {
		t.Fatalf("expected empty timezone, got %s", resp.Timezone)
	}

	if toFlattenedMeasures(nil) != nil {
		t.Fatalf("expected nil flattened measures")
	}
	if toFlattenedActivity(nil) != nil {
		t.Fatalf("expected nil flattened activity")
	}
}

func TestFlattenDailySummaryResponseJSONMarshalling(t *testing.T) {
	resp := FlattenDailySummaryResponse(&DailySummaryResponse{
		Summaries: []DailySummary{{}},
	})

	if _, err := json.Marshal(resp); err != nil {
		t.Fatalf("flattened response should marshal, got error: %v", err)
	}
}
