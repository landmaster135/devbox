package usecases

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"
	"time"
)

type stubHTTPClient struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (c *stubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.handler == nil {
		return nil, fmt.Errorf("handler is not defined")
	}
	return c.handler(req)
}

func TestFetchDailySummarySuccess(t *testing.T) {
	var measureCalls int32
	var activityCalls int32

	stubClient := &stubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

		form, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		switch form.Get("action") {
		case "getmeas":
			atomic.AddInt32(&measureCalls, 1)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":0,"body":{"timezone":"Europe/Paris","measuregrps":[{"grpid":1,"category":1,"date":1727740800,"measures":[{"value":70000,"type":1,"unit":-3},{"value":600,"type":11,"unit":-1}]}],"more":0,"offset":0}}`)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		case "getactivity":
			atomic.AddInt32(&activityCalls, 1)
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{"status":0,"body":{"activities":[{"date":"2024-10-01","timezone":"Europe/Paris","steps":8000,"calories":250.5,"totalcalories":300.0,"distance":6500.0,"elevation":10.5,"soft":1200,"moderate":600,"intense":300,"active":2100,"hr_average":58.5,"hr_min":45.0,"hr_max":120.0}],"more":false,"offset":0}}`)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		default:
			return nil, fmt.Errorf("unexpected action: %s", form.Get("action"))
		}
	}}

	service := NewHealthServiceWithHTTPClient("http://example.com", stubClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	start := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)

	resp, err := service.FetchDailySummary(ctx, DailySummaryRequest{
		AccessToken:     "dummy-token",
		UserID:          42,
		StartDate:       start,
		EndDate:         start,
		MeasureTypes:    []int{1, 11},
		IncludeActivity: true,
	})
	if err != nil {
		t.Fatalf("FetchDailySummary returned error: %v", err)
	}

	if measureCalls == 0 {
		t.Fatalf("measure endpoint was not invoked")
	}
	if activityCalls == 0 {
		t.Fatalf("activity endpoint was not invoked")
	}

	if resp.Timezone != "Europe/Paris" {
		t.Fatalf("unexpected timezone: %s", resp.Timezone)
	}

	if len(resp.Summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(resp.Summaries))
	}

	summary := resp.Summaries[0]
	if summary.Date != "2024-10-01" {
		t.Fatalf("unexpected summary date: %s", summary.Date)
	}

	weight, ok := summary.Measures["weight_kg"]
	if !ok {
		t.Fatalf("weight_kg not found in measures: %+v", summary.Measures)
	}
	if math.Abs(weight-70.0) > 1e-6 {
		t.Fatalf("unexpected weight value: %f", weight)
	}

	pulse, ok := summary.Measures["heart_pulse_bpm"]
	if !ok {
		t.Fatalf("heart_pulse_bpm not found in measures")
	}
	if math.Abs(pulse-60.0) > 1e-6 {
		t.Fatalf("unexpected heart pulse: %f", pulse)
	}

	if summary.Activity == nil {
		t.Fatalf("activity summary should not be nil")
	}
	if summary.Activity.Steps == nil || *summary.Activity.Steps != 8000 {
		t.Fatalf("unexpected steps value: %+v", summary.Activity.Steps)
	}
	if summary.Activity.CaloriesKcal == nil || math.Abs(*summary.Activity.CaloriesKcal-250.5) > 1e-6 {
		t.Fatalf("unexpected calories: %+v", summary.Activity.CaloriesKcal)
	}
	if summary.Activity.TotalCaloriesKcal == nil || math.Abs(*summary.Activity.TotalCaloriesKcal-300.0) > 1e-6 {
		t.Fatalf("unexpected total calories: %+v", summary.Activity.TotalCaloriesKcal)
	}
	if summary.Timezone != "Europe/Paris" {
		t.Fatalf("summary timezone mismatch: %s", summary.Timezone)
	}
}

func TestFetchDailySummaryMeasureAPIError(t *testing.T) {
	stubClient := &stubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"status":247,"error":"invalid_token","body":{}}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	}}

	service := NewHealthServiceWithHTTPClient("http://example.com", stubClient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	start := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)

	_, err := service.FetchDailySummary(ctx, DailySummaryRequest{
		AccessToken: "dummy-token",
		UserID:      42,
		StartDate:   start,
		EndDate:     start,
	})
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func TestFetchDailySummaryMissingToken(t *testing.T) {
	service := NewHealthService(0)
	_, err := service.FetchDailySummary(context.Background(), DailySummaryRequest{UserID: 1, StartDate: time.Now(), EndDate: time.Now()})
	if err == nil {
		t.Fatalf("expected error when token is missing")
	}
}
