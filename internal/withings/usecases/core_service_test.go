package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	configpkg "github.com/landmaster135/devbox/internal/withings/config"
)

type coreStubHTTPClient struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (c *coreStubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.handler == nil {
		return nil, fmt.Errorf("handler is not defined")
	}
	return c.handler(req)
}

func TestCoreServiceAccessors(t *testing.T) {
	health := NewHealthService(50 * time.Millisecond)
	auth := NewAuthService(50 * time.Millisecond)

	svc := NewCoreService(health, auth)

	if svc.GetHealthService() != health {
		t.Fatalf("unexpected health service pointer")
	}
	if svc.GetAuthService() != auth {
		t.Fatalf("unexpected auth service pointer")
	}
}

func TestCoreServiceDailySummaryRetriesOnUnauthorized(t *testing.T) {
	const (
		expiredToken = "expired-token"
		newAccess    = "new-access-token"
		newRefresh   = "new-refresh-token"
	)

	var authCalls int
	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		if err := req.ParseForm(); err != nil {
			t.Fatalf("failed to parse refresh form: %v", err)
		}
		if req.Form.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected grant type: %s", req.Form.Get("grant_type"))
		}
		if req.Form.Get("refresh_token") != "refresh-old" {
			t.Fatalf("unexpected refresh token: %s", req.Form.Get("refresh_token"))
		}
		authCalls++
		body := fmt.Sprintf(`{"status":0,"body":{"userid":"123","access_token":"%s","refresh_token":"%s","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`, newAccess, newRefresh)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	var measureCalls int
	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/measure" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		measureCalls++
		authHeader := req.Header.Get("Authorization")
		switch measureCalls {
		case 1:
			if authHeader != "Bearer "+expiredToken {
				t.Fatalf("expected expired token on first call, got %s", authHeader)
			}
			return &http.Response{
				StatusCode: http.StatusUnauthorized,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader("unauthorized")),
			}, nil
		case 2:
			if authHeader != "Bearer "+newAccess {
				t.Fatalf("expected refreshed token, got %s", authHeader)
			}
			payload := `{"status":0,"body":{"updatetime":1727740800,"timezone":"UTC","measuregrps":[{"grpid":1,"category":1,"date":1727740800,"measures":[{"value":70000,"type":1,"unit":-3}]}],"more":false,"offset":0}}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(strings.NewReader(payload)),
			}, nil
		default:
			return nil, fmt.Errorf("unexpected measure call count: %d", measureCalls)
		}
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:        "client-id",
		ClientSecret:    "client-secret",
		RefreshToken:    "refresh-old",
		AccessToken:     expiredToken,
		UserID:          42,
		StartDate:       time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:         time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		IncludeActivity: false,
	}

	resp, err := coreService.FetchDailySummaryWithRetry(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatal("response must not be nil")
	}
	if len(resp.Summaries) != 1 {
		t.Fatalf("expected one summary, got %d", len(resp.Summaries))
	}
	if cfg.AccessToken != newAccess {
		t.Fatalf("access token not updated: %s", cfg.AccessToken)
	}
	if cfg.RefreshToken != newRefresh {
		t.Fatalf("refresh token not updated: %s", cfg.RefreshToken)
	}
	if measureCalls != 2 {
		t.Fatalf("expected two measure calls, got %d", measureCalls)
	}
	if authCalls != 1 {
		t.Fatalf("expected one auth call, got %d", authCalls)
	}
}

func TestCoreServiceDailySummarySuccessWithoutRetry(t *testing.T) {
	var measureCalls int
	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		measureCalls++
		payload := `{"status":0,"body":{"timezone":"UTC","measuregrps":[{"grpid":1,"category":1,"date":1727740800,"measures":[{"value":80000,"type":1,"unit":-3}]}],"more":false,"offset":0}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
		}, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		t.Fatal("refresh should not be invoked on success")
		return nil, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:        "client",
		ClientSecret:    "secret",
		RefreshToken:    "refresh-token",
		AccessToken:     "valid-token",
		UserID:          99,
		StartDate:       time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC),
		EndDate:         time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC),
		IncludeActivity: false,
	}

	resp, err := coreService.FetchDailySummaryWithRetry(context.Background(), cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || len(resp.Summaries) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if cfg.AccessToken != "valid-token" {
		t.Fatalf("access token should remain unchanged, got %s", cfg.AccessToken)
	}
	if measureCalls != 1 {
		t.Fatalf("expected single measure call, got %d", measureCalls)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshSuccess(t *testing.T) {
	const (
		expiredToken = "expired-token"
		newAccess    = "new-access-token"
		newRefresh   = "new-refresh-token"
	)

	var authCalls int
	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		if err := req.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}
		if req.Form.Get("grant_type") != "refresh_token" {
			t.Fatalf("unexpected grant_type: %s", req.Form.Get("grant_type"))
		}
		if req.Form.Get("refresh_token") != "refresh-old" {
			t.Fatalf("unexpected refresh token: %s", req.Form.Get("refresh_token"))
		}
		authCalls++
		body := fmt.Sprintf(`{"status":0,"body":{"userid":"123","access_token":"%s","refresh_token":"%s","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`, newAccess, newRefresh)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	var measureCalls int
	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/measure" {
			t.Fatalf("unexpected path: %s", req.URL.Path)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer "+newAccess {
			t.Fatalf("authorization header not updated: %s", got)
		}
		if err := req.ParseForm(); err != nil {
			t.Fatalf("failed to parse measure form: %v", err)
		}
		measureCalls++
		payload := `{"status":0,"body":{"updatetime":1727740800,"timezone":"UTC","measuregrps":[{"grpid":1,"category":1,"date":1727740800,"measures":[{"value":70000,"type":1,"unit":-3}]}],"more":false,"offset":0}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
		}, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:     "client-id",
		ClientSecret: "client-secret",
		RefreshToken: "refresh-old",
		AccessToken:  expiredToken,
		UserID:       42,
		StartDate:    time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
	}

	req := DailySummaryRequest{
		AccessToken:     expiredToken,
		UserID:          cfg.UserID,
		StartDate:       cfg.StartDate,
		EndDate:         cfg.EndDate,
		IncludeActivity: false,
	}

	resp, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, req, errors.New("original error"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil {
		t.Fatalf("response is nil")
	}
	if len(resp.Summaries) != 1 {
		t.Fatalf("expected one summary, got %d", len(resp.Summaries))
	}
	if cfg.AccessToken != newAccess {
		t.Fatalf("access token not updated: %s", cfg.AccessToken)
	}
	if cfg.RefreshToken != newRefresh {
		t.Fatalf("refresh token not updated: %s", cfg.RefreshToken)
	}
	if measureCalls != 1 {
		t.Fatalf("unexpected measure call count: %d", measureCalls)
	}
	if authCalls != 1 {
		t.Fatalf("unexpected auth call count: %d", authCalls)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshConfigNil(t *testing.T) {
	coreService := NewCoreService(nil, nil)
	original := errors.New("original error")

	resp, err := coreService.RetryDailySummaryWithRefresh(context.Background(), nil, DailySummaryRequest{}, original)
	if resp != nil {
		t.Fatalf("expected nil response")
	}
	if err != original {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshRequiresRefreshToken(t *testing.T) {
	coreService := NewCoreService(nil, nil)
	cfg := &configpkg.Config{ClientID: "client", ClientSecret: "secret"}
	original := errors.New("unauthorized")

	_, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, DailySummaryRequest{}, original)
	if err != original {
		t.Fatalf("expected original error, got %v", err)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshRequiresClientCredentials(t *testing.T) {
	coreService := NewCoreService(nil, nil)
	cfg := &configpkg.Config{RefreshToken: "refresh-only"}

	_, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, DailySummaryRequest{}, errors.New("unauthorized"))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "client-id") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshRefreshFailure(t *testing.T) {
	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network outage")
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		t.Fatalf("health service should not be called when refresh fails")
		return nil, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "refresh-token",
		AccessToken:  "stale",
		UserID:       1,
		StartDate:    time.Date(2024, 10, 2, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 10, 2, 0, 0, 0, 0, time.UTC),
	}

	req := DailySummaryRequest{AccessToken: "stale", UserID: cfg.UserID, StartDate: cfg.StartDate, EndDate: cfg.EndDate}

	resp, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, req, errors.New("unauthorized"))
	if resp != nil {
		t.Fatalf("expected nil response")
	}
	if err == nil || !strings.Contains(err.Error(), "アクセストークンのリフレッシュに失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AccessToken != "stale" {
		t.Fatalf("access token should not change on refresh failure")
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshMissingAccessToken(t *testing.T) {
	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		body := `{"status":0,"body":{"userid":"1","access_token":"","refresh_token":"","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		t.Fatalf("health service should not be called when access token missing")
		return nil, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "refresh-token",
		AccessToken:  "expired",
		UserID:       7,
		StartDate:    time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 10, 3, 0, 0, 0, 0, time.UTC),
	}

	req := DailySummaryRequest{AccessToken: "expired", UserID: cfg.UserID, StartDate: cfg.StartDate, EndDate: cfg.EndDate}

	resp, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, req, errors.New("unauthorized"))
	if resp != nil {
		t.Fatalf("expected nil response")
	}
	if err == nil || !strings.Contains(err.Error(), "リフレッシュレスポンスにアクセストークンが含まれていません") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshFetchFailure(t *testing.T) {
	const newToken = "fresh-token"

	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		body := fmt.Sprintf(`{"status":0,"body":{"userid":"1","access_token":"%s","refresh_token":"r2","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`, newToken)
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("Authorization"); got != "Bearer "+newToken {
			t.Fatalf("expected refreshed token, got %s", got)
		}
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader("server error")),
		}, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "refresh-token",
		AccessToken:  "expired",
		UserID:       9,
		StartDate:    time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 10, 4, 0, 0, 0, 0, time.UTC),
	}

	req := DailySummaryRequest{AccessToken: "expired", UserID: cfg.UserID, StartDate: cfg.StartDate, EndDate: cfg.EndDate}

	resp, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, req, errors.New("unauthorized"))
	if resp != nil {
		t.Fatalf("expected nil response")
	}
	if err == nil || !strings.Contains(err.Error(), "リフレッシュ後の日次サマリ取得にも失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AccessToken != newToken {
		t.Fatalf("expected access token to update despite fetch failure")
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshCopiesRequestAccessToken(t *testing.T) {
	var cachedToken string
	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		cachedToken = req.Header.Get("Authorization")
		payload := `{"status":0,"body":{"updatetime":1727740800,"timezone":"UTC","measuregrps":[],"more":false,"offset":0}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
		}, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		body := `{"status":0,"body":{"userid":"1","access_token":"updated","refresh_token":"refresh","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	coreService := NewCoreService(healthService, authService)

	cfg := &configpkg.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "refresh",
		AccessToken:  "before",
		UserID:       101,
		StartDate:    time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC),
		EndDate:      time.Date(2024, 10, 5, 0, 0, 0, 0, time.UTC),
	}

	req := DailySummaryRequest{
		AccessToken: "before",
		UserID:      cfg.UserID,
		StartDate:   cfg.StartDate,
		EndDate:     cfg.EndDate,
	}

	if _, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, req, errors.New("unauthorized")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cachedToken != "Bearer updated" {
		t.Fatalf("expected request to use refreshed token, got %s", cachedToken)
	}
}

func TestCoreServiceRetryDailySummaryWithRefreshPropagatesRequestFields(t *testing.T) {
	var receivedForm url.Values
	healthClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		req.Body.Close()
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		receivedForm, err = url.ParseQuery(string(bodyBytes))
		if err != nil {
			t.Fatalf("failed to parse body: %v", err)
		}
		payload := `{"status":0,"body":{"updatetime":1727740800,"timezone":"UTC","measuregrps":[],"more":false,"offset":0}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(payload)),
		}, nil
	}}
	healthService := NewHealthServiceWithHTTPClient("https://api.example.com", healthClient)

	authClient := &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		body := `{"status":0,"body":{"userid":"1","access_token":"updated","refresh_token":"refresh","expires_in":3600,"token_type":"Bearer","scope":"user.metrics"}}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}}
	authService := NewAuthServiceWithClient(authClient, time.Second)

	coreService := NewCoreService(healthService, authService)

	start := time.Date(2024, 10, 6, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 10, 6, 0, 0, 0, 0, time.UTC)

	cfg := &configpkg.Config{
		ClientID:     "client",
		ClientSecret: "secret",
		RefreshToken: "refresh",
		AccessToken:  "before",
		UserID:       202,
		StartDate:    start,
		EndDate:      end,
	}

	req := DailySummaryRequest{
		AccessToken:  "before",
		UserID:       cfg.UserID,
		StartDate:    start,
		EndDate:      end,
		MeasureTypes: []int{1, 6, 9},
	}

	if _, err := coreService.RetryDailySummaryWithRefresh(context.Background(), cfg, req, errors.New("unauthorized")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if receivedForm == nil {
		t.Fatalf("expected form to be captured")
	}
	if receivedForm.Get("userid") != fmt.Sprintf("%d", cfg.UserID) {
		t.Fatalf("unexpected userid: %s", receivedForm.Get("userid"))
	}
	if receivedForm.Get("meastypes") != "1,6,9" {
		t.Fatalf("unexpected measure types: %s", receivedForm.Get("meastypes"))
	}
}

func TestRoundDailySummaryResponseRounding(t *testing.T) {
	weight := 70.129
	spo2 := math.Inf(1)
	distance := 6543.219
	calories := 250.123
	totalCalories := 300.987
	nanValue := math.NaN()

	resp := &DailySummaryResponse{
		Summaries: []DailySummary{
			{
				Measures: &DailySummaryMeasures{
					WeightKg:    &weight,
					Spo2Percent: &spo2,
				},
				Activity: &ActivitySummary{
					DistanceMeter:     &distance,
					CaloriesKcal:      &calories,
					TotalCaloriesKcal: &totalCalories,
					HrAverageBPM:      &nanValue,
				},
			},
			{},
		},
	}

	roundDailySummaryResponse(resp)

	if resp.Summaries[0].Measures == nil || resp.Summaries[0].Measures.WeightKg == nil {
		t.Fatalf("weight should be rounded")
	}
	if got := *resp.Summaries[0].Measures.WeightKg; got != 70.13 {
		t.Fatalf("unexpected rounded weight: %v", got)
	}
	if resp.Summaries[0].Measures.Spo2Percent == nil || !math.IsInf(*resp.Summaries[0].Measures.Spo2Percent, 0) {
		t.Fatalf("infinite value should remain unchanged")
	}

	activity := resp.Summaries[0].Activity
	if activity == nil || activity.DistanceMeter == nil {
		t.Fatalf("activity should not be nil")
	}
	if got := *activity.DistanceMeter; got != 6543.22 {
		t.Fatalf("unexpected rounded distance: %v", got)
	}
	if got := *activity.CaloriesKcal; got != 250.12 {
		t.Fatalf("unexpected rounded calories: %v", got)
	}
	if got := *activity.TotalCaloriesKcal; got != 300.99 {
		t.Fatalf("unexpected rounded total calories: %v", got)
	}
	if activity.HrAverageBPM == nil || !math.IsNaN(*activity.HrAverageBPM) {
		t.Fatalf("NaN value should remain NaN")
	}

	roundDailySummaryResponse(nil)
}

func TestFetchDailySummaryWithRetryRequiresHealthService(t *testing.T) {
	core := &CoreService{}
	_, err := core.FetchDailySummaryWithRetry(context.Background(), &configpkg.Config{})
	if err == nil || !strings.Contains(err.Error(), "health service") {
		t.Fatalf("expected health service error, got %v", err)
	}
}

func TestFetchDailySummaryWithRetryRequiresConfig(t *testing.T) {
	health := NewHealthService(0)
	core := NewCoreService(health, nil)
	_, err := core.FetchDailySummaryWithRetry(context.Background(), nil)
	if err == nil || !strings.Contains(err.Error(), "config") {
		t.Fatalf("expected config error, got %v", err)
	}
}

func TestFetchDailySummaryWithRetryPropagatesNonRetryableError(t *testing.T) {
	health := NewHealthServiceWithHTTPClient("https://api.example.com", &coreStubHTTPClient{handler: func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network down")
	}})
	core := NewCoreService(health, nil)

	cfg := &configpkg.Config{
		AccessToken:     "token",
		UserID:          123,
		StartDate:       time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC),
		EndDate:         time.Date(2024, 10, 7, 0, 0, 0, 0, time.UTC),
		IncludeActivity: false,
	}

	_, err := core.FetchDailySummaryWithRetry(context.Background(), cfg)
	if err == nil || !strings.Contains(err.Error(), "network down") {
		t.Fatalf("expected propagated error, got %v", err)
	}
}
