package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	// DefaultBaseURL は Withings Public API のベース URL です。
	DefaultBaseURL     = "https://wbsapi.withings.net"
	defaultUserAgent   = "devbox-withings-cli/0.1"
	maxPaginationLoops = 100
)

// HTTPClient は http.Client 互換のインターフェースです。
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HealthService は Withings API から健康データを取得するユースケースを提供します。
type HealthService struct {
	httpClient HTTPClient
	baseURL    string
	userAgent  string
}

// NewHealthService は指定したタイムアウトで HTTP クライアントを作成して HealthService を生成します。
func NewHealthService(timeout time.Duration) *HealthService {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	client := &http.Client{Timeout: timeout}
	return NewHealthServiceWithHTTPClient(DefaultBaseURL, client)
}

// NewHealthServiceWithHTTPClient はベース URL と HTTP クライアントを差し替えて HealthService を生成します。
func NewHealthServiceWithHTTPClient(baseURL string, client HTTPClient) *HealthService {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	svc := &HealthService{
		httpClient: client,
		baseURL:    DefaultBaseURL,
		userAgent:  defaultUserAgent,
	}
	svc.SetBaseURL(baseURL)
	svc.SetUserAgent(defaultUserAgent)
	return svc
}

// SetUserAgent は API 呼び出し時の User-Agent を上書きします。
func (s *HealthService) SetUserAgent(agent string) {
	if trimmed := strings.TrimSpace(agent); trimmed != "" {
		s.userAgent = trimmed
	}
}

// SetBaseURL は API 呼び出し時のベース URL を上書きします。
func (s *HealthService) SetBaseURL(baseURL string) {
	trimmed := strings.TrimSpace(baseURL)
	if trimmed == "" {
		return
	}
	s.baseURL = strings.TrimRight(trimmed, "/")
}

// DailySummaryRequest は FetchDailySummary の入力を表します。
type DailySummaryRequest struct {
	AccessToken     string
	UserID          int64
	StartDate       time.Time
	EndDate         time.Time
	MeasureTypes    []int
	IncludeActivity bool
}

// DailySummaryResponse は日次データの結果です。
type DailySummaryResponse struct {
	Summaries []DailySummary `json:"summaries"`
	Timezone  string         `json:"timezone,omitempty"`
}

// DailySummary は 1 日分の測定値と活動サマリです。
type DailySummary struct {
	Date     string             `json:"date"`
	Timezone string             `json:"timezone,omitempty"`
	Measures map[string]float64 `json:"measures,omitempty"`
	Activity *ActivitySummary   `json:"activity,omitempty"`
}

// ActivitySummary は Withings の日次活動サマリをラップします。
type ActivitySummary struct {
	Steps             *int     `json:"steps,omitempty"`
	DistanceMeter     *float64 `json:"distance_meter,omitempty"`
	ElevationMeter    *float64 `json:"elevation_meter,omitempty"`
	CaloriesKcal      *float64 `json:"calories_kcal,omitempty"`
	TotalCaloriesKcal *float64 `json:"total_calories_kcal,omitempty"`
	SoftSeconds       *int     `json:"soft_seconds,omitempty"`
	ModerateSeconds   *int     `json:"moderate_seconds,omitempty"`
	IntenseSeconds    *int     `json:"intense_seconds,omitempty"`
	ActiveSeconds     *int     `json:"active_seconds,omitempty"`
	HrAverageBPM      *float64 `json:"hr_average_bpm,omitempty"`
	HrMinBPM          *float64 `json:"hr_min_bpm,omitempty"`
	HrMaxBPM          *float64 `json:"hr_max_bpm,omitempty"`
	DeviceBrand       *int     `json:"device_brand,omitempty"`
	DeviceModelID     *int     `json:"device_model_id,omitempty"`
	DeviceModelName   *string  `json:"device_model_name,omitempty"`
	IsTracker         *bool    `json:"is_tracker,omitempty"`
}

func (s *HealthService) ShouldRetryDailySummaryWithRefresh(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "status=401") ||
		strings.Contains(message, "status=403") ||
		strings.Contains(message, "unauthorized") ||
		strings.Contains(message, "invalid_token")
}

// FetchDailySummary は指定期間の測定値と活動サマリをまとめて取得します。
func (s *HealthService) FetchDailySummary(ctx context.Context, req DailySummaryRequest) (*DailySummaryResponse, error) {
	if strings.TrimSpace(req.AccessToken) == "" {
		return nil, errors.New("access token が指定されていません")
	}
	if req.UserID <= 0 {
		return nil, errors.New("userID は正の整数で指定してください")
	}
	if req.StartDate.IsZero() {
		return nil, errors.New("開始日が指定されていません")
	}
	if req.EndDate.IsZero() {
		req.EndDate = req.StartDate
	}
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("終了日は開始日以降の日付を指定してください")
	}

	startUTC, endUTC := normalizeDateRange(req.StartDate, req.EndDate)

	measureResult, err := s.fetchMeasures(ctx, req.AccessToken, req.UserID, startUTC, endUTC, req.MeasureTypes)
	if err != nil {
		return nil, err
	}

	summaries := make(map[string]*DailySummary)

	for dateKey, measureMap := range measureResult.measurements {
		// map を複製して呼び出し元での意図せぬ変更を防ぐ
		measuresCopy := make(map[string]float64, len(measureMap))
		for k, v := range measureMap {
			measuresCopy[k] = v
		}
		summaries[dateKey] = &DailySummary{
			Date:     dateKey,
			Timezone: measureResult.timezone,
			Measures: measuresCopy,
		}
	}

	var timezone = measureResult.timezone

	if req.IncludeActivity {
		activityResult, err := s.fetchActivities(ctx, req.AccessToken, req.UserID, req.StartDate, req.EndDate)
		if err != nil {
			return nil, err
		}
		if timezone == "" {
			timezone = activityResult.defaultTimezone
		}
		for dateKey, activityData := range activityResult.activities {
			summary, exists := summaries[dateKey]
			if !exists {
				summary = &DailySummary{Date: dateKey}
				summaries[dateKey] = summary
			}
			if summary.Timezone == "" {
				summary.Timezone = activityData.timezone
			}
			summary.Activity = &activityData.summary
		}
	}

	// 集計結果をソートして返却
	dates := make([]string, 0, len(summaries))
	for date := range summaries {
		dates = append(dates, date)
	}
	sort.Strings(dates)

	ordered := make([]DailySummary, 0, len(dates))
	for _, dateKey := range dates {
		ordered = append(ordered, *summaries[dateKey])
	}

	return &DailySummaryResponse{
		Summaries: ordered,
		Timezone:  timezone,
	}, nil
}

func (s *HealthService) fetchMeasures(ctx context.Context, token string, userID int64, start, end time.Time, measureTypes []int) (*measureFetchResult, error) {
	offset := 0
	loops := 0
	result := &measureFetchResult{measurements: make(map[string]map[string]float64)}
	seen := make(map[string]map[int]measurementValue)

	for {
		if loops >= maxPaginationLoops {
			return nil, fmt.Errorf("measure API のページネーション制限を超えました (%d)", maxPaginationLoops)
		}
		loops++

		form := url.Values{}
		form.Set("action", "getmeas")
		form.Set("userid", strconv.FormatInt(userID, 10))
		form.Set("startdate", strconv.FormatInt(start.Unix(), 10))
		form.Set("enddate", strconv.FormatInt(end.Unix(), 10))
		form.Set("category", "1")
		if len(measureTypes) > 0 {
			form.Set("meastypes", joinInts(measureTypes))
		}
		if offset > 0 {
			form.Set("offset", strconv.Itoa(offset))
		}

		resp, err := s.postForm(ctx, "/measure", token, form)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("measure API のレスポンス読み取りに失敗しました: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("measure API がエラーを返しました: status=%d, body=%s", resp.StatusCode, truncateForLog(bodyBytes))
		}

		var payload measureAPIResponse
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			return nil, fmt.Errorf("measure API の JSON 解析に失敗しました: %w", err)
		}

		if payload.Status != 0 {
			return nil, fmt.Errorf("measure API でエラーが発生しました: status=%d, error=%v", payload.Status, payload.Error)
		}

		if payload.Body.Timezone != "" {
			result.timezone = payload.Body.Timezone
		}

		loc := time.UTC
		if payload.Body.Timezone != "" {
			if tz, err := time.LoadLocation(payload.Body.Timezone); err == nil {
				loc = tz
			}
		}

		for _, group := range payload.Body.Measuregrps {
			if group.Category != 1 {
				continue
			}
			recordedAt := time.Unix(group.Date, 0).In(loc)
			dateKey := recordedAt.Format("2006-01-02")
			if _, ok := seen[dateKey]; !ok {
				seen[dateKey] = make(map[int]measurementValue)
			}
			for _, m := range group.Measures {
				value := convertMeasureValue(m.Value, m.Unit)
				current := seen[dateKey][m.Type]
				if current.timestamp.IsZero() || current.timestamp.Before(recordedAt) {
					seen[dateKey][m.Type] = measurementValue{value: value, timestamp: recordedAt}
				}
			}
		}

		if !payload.Body.More.Bool() {
			break
		}
		if payload.Body.Offset == 0 {
			break
		}
		offset = payload.Body.Offset
	}

	for dateKey, typeMap := range seen {
		measurements := make(map[string]float64, len(typeMap))
		for measureType, stored := range typeMap {
			label := labelForMeasureType(measureType)
			measurements[label] = stored.value
		}
		result.measurements[dateKey] = measurements
	}

	return result, nil
}

func (s *HealthService) fetchActivities(ctx context.Context, token string, userID int64, start, end time.Time) (*activityFetchResult, error) {
	offset := 0
	loops := 0
	result := &activityFetchResult{
		activities: make(map[string]activityEntryWithTZ),
	}

	for {
		if loops >= maxPaginationLoops {
			return nil, fmt.Errorf("activity API のページネーション制限を超えました (%d)", maxPaginationLoops)
		}
		loops++

		form := url.Values{}
		form.Set("action", "getactivity")
		form.Set("userid", strconv.FormatInt(userID, 10))
		form.Set("startdateymd", start.Format("2006-01-02"))
		form.Set("enddateymd", end.Format("2006-01-02"))
		if offset > 0 {
			form.Set("offset", strconv.Itoa(offset))
		}

		resp, err := s.postForm(ctx, "/v2/measure", token, form)
		if err != nil {
			return nil, err
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("activity API のレスポンス読み取りに失敗しました: %w", err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("activity API がエラーを返しました: status=%d, body=%s", resp.StatusCode, truncateForLog(bodyBytes))
		}

		var payload activityAPIResponse
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			return nil, fmt.Errorf("activity API の JSON 解析に失敗しました: %w", err)
		}

		if payload.Status != 0 {
			return nil, fmt.Errorf("activity API でエラーが発生しました: status=%d, error=%v", payload.Status, payload.Error)
		}

		if result.defaultTimezone == "" && len(payload.Body.Activities) > 0 {
			result.defaultTimezone = payload.Body.Activities[0].Timezone
		}

		for _, activity := range payload.Body.Activities {
			summary := buildActivitySummary(activity)
			existing, ok := result.activities[activity.Date]
			if ok {
				summary = mergeActivitySummaries(existing.summary, summary)
			}
			tz := activity.Timezone
			if tz == "" && ok {
				tz = existing.timezone
			}
			result.activities[activity.Date] = activityEntryWithTZ{summary: summary, timezone: tz}
		}

		if !payload.Body.More.Bool() {
			break
		}
		if payload.Body.Offset == 0 {
			break
		}
		offset = payload.Body.Offset
	}

	return result, nil
}

func (s *HealthService) postForm(ctx context.Context, path string, token string, form url.Values) (*http.Response, error) {
	endpoint := s.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("リクエストの構築に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	if s.userAgent != "" {
		req.Header.Set("User-Agent", s.userAgent)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API リクエストに失敗しました: %w", err)
	}
	return resp, nil
}

func normalizeDateRange(start, end time.Time) (time.Time, time.Time) {
	startUTC := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endUTC := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	return startUTC, endUTC
}

func convertMeasureValue(value int64, unit int) float64 {
	return float64(value) * math.Pow10(unit)
}

func joinInts(values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

func labelForMeasureType(measureType int) string {
	if label, ok := measureTypeLabels[measureType]; ok {
		return label
	}
	return fmt.Sprintf("measure_%d", measureType)
}

func truncateForLog(data []byte) string {
	const limit = 512
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) <= limit {
		return string(trimmed)
	}
	return string(trimmed[:limit]) + "..."
}

type measurementValue struct {
	value     float64
	timestamp time.Time
}

type measureFetchResult struct {
	timezone     string
	measurements map[string]map[string]float64
}

type activityFetchResult struct {
	activities      map[string]activityEntryWithTZ
	defaultTimezone string
}

type activityEntryWithTZ struct {
	summary  ActivitySummary
	timezone string
}

type measureAPIResponse struct {
	Status int         `json:"status"`
	Error  interface{} `json:"error"`
	Body   struct {
		Updatetime  int64          `json:"updatetime"`
		Timezone    string         `json:"timezone"`
		Measuregrps []measureGroup `json:"measuregrps"`
		More        flexibleBool   `json:"more"`
		Offset      int            `json:"offset"`
	} `json:"body"`
}

type measureGroup struct {
	GroupID  int64         `json:"grpid"`
	Category int           `json:"category"`
	Date     int64         `json:"date"`
	Measures []measureItem `json:"measures"`
}

type measureItem struct {
	Value int64 `json:"value"`
	Type  int   `json:"type"`
	Unit  int   `json:"unit"`
}

type activityAPIResponse struct {
	Status int         `json:"status"`
	Error  interface{} `json:"error"`
	Body   struct {
		Activities []activityItem `json:"activities"`
		More       flexibleBool   `json:"more"`
		Offset     int            `json:"offset"`
	} `json:"body"`
}

type activityItem struct {
	Date          string   `json:"date"`
	Timezone      string   `json:"timezone"`
	Steps         *int     `json:"steps"`
	Calories      *float64 `json:"calories"`
	TotalCalories *float64 `json:"totalcalories"`
	Distance      *float64 `json:"distance"`
	Elevation     *float64 `json:"elevation"`
	Soft          *int     `json:"soft"`
	Moderate      *int     `json:"moderate"`
	Intense       *int     `json:"intense"`
	Active        *int     `json:"active"`
	HrAverage     *float64 `json:"hr_average"`
	HrMin         *float64 `json:"hr_min"`
	HrMax         *float64 `json:"hr_max"`
	Brand         *int     `json:"brand"`
	ModelID       *int     `json:"modelid"`
	Model         *string  `json:"model"`
	IsTracker     *bool    `json:"is_tracker"`
}

type flexibleBool bool

func (fb *flexibleBool) UnmarshalJSON(data []byte) error {
	stripped := bytes.TrimSpace(data)
	if len(stripped) == 0 {
		*fb = flexibleBool(false)
		return nil
	}
	switch stripped[0] {
	case 't', 'f':
		var b bool
		if err := json.Unmarshal(stripped, &b); err != nil {
			return err
		}
		*fb = flexibleBool(b)
		return nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-':
		var i int
		if err := json.Unmarshal(stripped, &i); err != nil {
			return err
		}
		*fb = flexibleBool(i != 0)
		return nil
	default:
		return fmt.Errorf("bool 表現を解釈できません: %s", string(stripped))
	}
}

func (fb flexibleBool) Bool() bool {
	return bool(fb)
}

func buildActivitySummary(item activityItem) ActivitySummary {
	summary := ActivitySummary{}
	if item.Steps != nil {
		summary.Steps = intPtr(*item.Steps)
	}
	if item.Distance != nil {
		summary.DistanceMeter = floatPtr(*item.Distance)
	}
	if item.Elevation != nil {
		summary.ElevationMeter = floatPtr(*item.Elevation)
	}
	if item.Calories != nil {
		summary.CaloriesKcal = floatPtr(*item.Calories)
	}
	if item.TotalCalories != nil {
		summary.TotalCaloriesKcal = floatPtr(*item.TotalCalories)
	}
	if item.Soft != nil {
		summary.SoftSeconds = intPtr(*item.Soft)
	}
	if item.Moderate != nil {
		summary.ModerateSeconds = intPtr(*item.Moderate)
	}
	if item.Intense != nil {
		summary.IntenseSeconds = intPtr(*item.Intense)
	}
	if item.Active != nil {
		summary.ActiveSeconds = intPtr(*item.Active)
	}
	if item.HrAverage != nil {
		summary.HrAverageBPM = floatPtr(*item.HrAverage)
	}
	if item.HrMin != nil {
		summary.HrMinBPM = floatPtr(*item.HrMin)
	}
	if item.HrMax != nil {
		summary.HrMaxBPM = floatPtr(*item.HrMax)
	}
	if item.Brand != nil {
		summary.DeviceBrand = intPtr(*item.Brand)
	}
	if item.ModelID != nil {
		summary.DeviceModelID = intPtr(*item.ModelID)
	}
	if item.Model != nil {
		summary.DeviceModelName = stringPtr(*item.Model)
	}
	if item.IsTracker != nil {
		summary.IsTracker = boolPtr(*item.IsTracker)
	}
	return summary
}

func mergeActivitySummaries(base, incoming ActivitySummary) ActivitySummary {
	if incoming.Steps != nil {
		base.Steps = incoming.Steps
	}
	if incoming.DistanceMeter != nil {
		base.DistanceMeter = incoming.DistanceMeter
	}
	if incoming.ElevationMeter != nil {
		base.ElevationMeter = incoming.ElevationMeter
	}
	if incoming.CaloriesKcal != nil {
		base.CaloriesKcal = incoming.CaloriesKcal
	}
	if incoming.TotalCaloriesKcal != nil {
		base.TotalCaloriesKcal = incoming.TotalCaloriesKcal
	}
	if incoming.SoftSeconds != nil {
		base.SoftSeconds = incoming.SoftSeconds
	}
	if incoming.ModerateSeconds != nil {
		base.ModerateSeconds = incoming.ModerateSeconds
	}
	if incoming.IntenseSeconds != nil {
		base.IntenseSeconds = incoming.IntenseSeconds
	}
	if incoming.ActiveSeconds != nil {
		base.ActiveSeconds = incoming.ActiveSeconds
	}
	if incoming.HrAverageBPM != nil {
		base.HrAverageBPM = incoming.HrAverageBPM
	}
	if incoming.HrMinBPM != nil {
		base.HrMinBPM = incoming.HrMinBPM
	}
	if incoming.HrMaxBPM != nil {
		base.HrMaxBPM = incoming.HrMaxBPM
	}
	if incoming.DeviceBrand != nil {
		base.DeviceBrand = incoming.DeviceBrand
	}
	if incoming.DeviceModelID != nil {
		base.DeviceModelID = incoming.DeviceModelID
	}
	if incoming.DeviceModelName != nil {
		base.DeviceModelName = incoming.DeviceModelName
	}
	if incoming.IsTracker != nil {
		base.IsTracker = incoming.IsTracker
	}
	return base
}

func intPtr(v int) *int {
	value := v
	return &value
}

func floatPtr(v float64) *float64 {
	value := v
	return &value
}

func boolPtr(v bool) *bool {
	value := v
	return &value
}

func stringPtr(v string) *string {
	value := v
	return &value
}

var measureTypeLabels = map[int]string{
	1:   "weight_kg",
	4:   "height_meter",
	5:   "fat_free_mass_kg",
	6:   "fat_ratio_percent",
	8:   "fat_mass_kg",
	9:   "diastolic_bp_mmhg",
	10:  "systolic_bp_mmhg",
	11:  "heart_pulse_bpm",
	12:  "temperature_c",
	54:  "spo2_percent",
	71:  "body_temperature_c",
	73:  "skin_temperature_c",
	76:  "muscle_mass_kg",
	77:  "hydration_kg",
	88:  "bone_mass_kg",
	91:  "pulse_wave_velocity_m_per_s",
	123: "vo2max_ml_per_min_per_kg",
	130: "atrial_fibrillation_result",
	135: "qrs_duration_ms",
	136: "pr_duration_ms",
	137: "qt_duration_ms",
	138: "qt_corrected_duration_ms",
	139: "atrial_fibrillation_ppg",
	155: "vascular_age_years",
	167: "nerve_health_conductance_feet",
	168: "extracellular_water_kg",
	169: "intracellular_water_kg",
	170: "visceral_fat_index",
	173: "segment_fat_free_mass_kg",
	174: "segment_fat_mass_kg",
	175: "segment_muscle_mass_kg",
	196: "electrodermal_activity_feet",
	226: "basal_metabolic_rate",
	227: "metabolic_age_years",
	229: "electrochemical_skin_conductance",
}
