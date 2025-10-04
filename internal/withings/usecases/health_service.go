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
	DefaultBaseURL      = "https://wbsapi.withings.net"
	EndpointOfMeasure   = "/measure"
	EndpointOfMeasureV2 = "/v2/measure"
	defaultUserAgent    = "devbox-withings-cli/0.1"
	maxPaginationLoops  = 100
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
	Date     string                `json:"date"`
	Timezone string                `json:"timezone,omitempty"`
	Measures *DailySummaryMeasures `json:"measures,omitempty"`
	Activity *ActivitySummary      `json:"activity,omitempty"`
}

// DailySummaryMeasures は測定値をラベルごとに保持します。
type DailySummaryMeasures struct {
	WeightKg                   *float64 `json:"weight_kg,omitempty"`
	HeightMeter                *float64 `json:"height_meter,omitempty"`
	FatFreeMassKg              *float64 `json:"fat_free_mass_kg,omitempty"`
	FatRatioPercent            *float64 `json:"fat_ratio_percent,omitempty"`
	FatMassKg                  *float64 `json:"fat_mass_kg,omitempty"`
	DiastolicBpMmhg            *float64 `json:"diastolic_bp_mmhg,omitempty"`
	SystolicBpMmhg             *float64 `json:"systolic_bp_mmhg,omitempty"`
	HeartPulseBpm              *float64 `json:"heart_pulse_bpm,omitempty"`
	TemperatureC               *float64 `json:"temperature_c,omitempty"`
	Spo2Percent                *float64 `json:"spo2_percent,omitempty"`
	BodyTemperatureC           *float64 `json:"body_temperature_c,omitempty"`
	SkinTemperatureC           *float64 `json:"skin_temperature_c,omitempty"`
	MuscleMassKg               *float64 `json:"muscle_mass_kg,omitempty"`
	HydrationKg                *float64 `json:"hydration_kg,omitempty"`
	BoneMassKg                 *float64 `json:"bone_mass_kg,omitempty"`
	PulseWaveVelocityMPerS     *float64 `json:"pulse_wave_velocity_m_per_s,omitempty"`
	Vo2MaxMlPerMinPerKg        *float64 `json:"vo2max_ml_per_min_per_kg,omitempty"`
	AtrialFibrillationResult   *float64 `json:"atrial_fibrillation_result,omitempty"`
	QrsDurationMs              *float64 `json:"qrs_duration_ms,omitempty"`
	PrDurationMs               *float64 `json:"pr_duration_ms,omitempty"`
	QtDurationMs               *float64 `json:"qt_duration_ms,omitempty"`
	QtCorrectedDurationMs      *float64 `json:"qt_corrected_duration_ms,omitempty"`
	AtrialFibrillationPpg      *float64 `json:"atrial_fibrillation_ppg,omitempty"`
	VascularAgeYears           *float64 `json:"vascular_age_years,omitempty"`
	NerveHealthConductanceFeet *float64 `json:"nerve_health_conductance_feet,omitempty"`
	ExtracellularWaterKg       *float64 `json:"extracellular_water_kg,omitempty"`
	IntracellularWaterKg       *float64 `json:"intracellular_water_kg,omitempty"`
	VisceralFatIndex           *float64 `json:"visceral_fat_index,omitempty"`
	SegmentFatFreeMassKg       *float64 `json:"segment_fat_free_mass_kg,omitempty"`
	SegmentFatMassKg           *float64 `json:"segment_fat_mass_kg,omitempty"`
	SegmentMuscleMassKg        *float64 `json:"segment_muscle_mass_kg,omitempty"`
	ElectrodermalActivityFeet  *float64 `json:"electrodermal_activity_feet,omitempty"`
	BasalMetabolicRate         *float64 `json:"basal_metabolic_rate,omitempty"`
	MetabolicAgeYears          *float64 `json:"metabolic_age_years,omitempty"`
	ElectrochemicalSkinConduct *float64 `json:"electrochemical_skin_conductance,omitempty"`
}

func (m *DailySummaryMeasures) set(label string, value float64) bool {
	if m == nil {
		return false
	}
	v := value
	switch label {
	case "weight_kg":
		m.WeightKg = &v
	case "height_meter":
		m.HeightMeter = &v
	case "fat_free_mass_kg":
		m.FatFreeMassKg = &v
	case "fat_ratio_percent":
		m.FatRatioPercent = &v
	case "fat_mass_kg":
		m.FatMassKg = &v
	case "diastolic_bp_mmhg":
		m.DiastolicBpMmhg = &v
	case "systolic_bp_mmhg":
		m.SystolicBpMmhg = &v
	case "heart_pulse_bpm":
		m.HeartPulseBpm = &v
	case "temperature_c":
		m.TemperatureC = &v
	case "spo2_percent":
		m.Spo2Percent = &v
	case "body_temperature_c":
		m.BodyTemperatureC = &v
	case "skin_temperature_c":
		m.SkinTemperatureC = &v
	case "muscle_mass_kg":
		m.MuscleMassKg = &v
	case "hydration_kg":
		m.HydrationKg = &v
	case "bone_mass_kg":
		m.BoneMassKg = &v
	case "pulse_wave_velocity_m_per_s":
		m.PulseWaveVelocityMPerS = &v
	case "vo2max_ml_per_min_per_kg":
		m.Vo2MaxMlPerMinPerKg = &v
	case "atrial_fibrillation_result":
		m.AtrialFibrillationResult = &v
	case "qrs_duration_ms":
		m.QrsDurationMs = &v
	case "pr_duration_ms":
		m.PrDurationMs = &v
	case "qt_duration_ms":
		m.QtDurationMs = &v
	case "qt_corrected_duration_ms":
		m.QtCorrectedDurationMs = &v
	case "atrial_fibrillation_ppg":
		m.AtrialFibrillationPpg = &v
	case "vascular_age_years":
		m.VascularAgeYears = &v
	case "nerve_health_conductance_feet":
		m.NerveHealthConductanceFeet = &v
	case "extracellular_water_kg":
		m.ExtracellularWaterKg = &v
	case "intracellular_water_kg":
		m.IntracellularWaterKg = &v
	case "visceral_fat_index":
		m.VisceralFatIndex = &v
	case "segment_fat_free_mass_kg":
		m.SegmentFatFreeMassKg = &v
	case "segment_fat_mass_kg":
		m.SegmentFatMassKg = &v
	case "segment_muscle_mass_kg":
		m.SegmentMuscleMassKg = &v
	case "electrodermal_activity_feet":
		m.ElectrodermalActivityFeet = &v
	case "basal_metabolic_rate":
		m.BasalMetabolicRate = &v
	case "metabolic_age_years":
		m.MetabolicAgeYears = &v
	case "electrochemical_skin_conductance":
		m.ElectrochemicalSkinConduct = &v
	default:
		return false
	}
	return true
}

func (m *DailySummaryMeasures) roundToTwoDecimalPlaces() {
	if m == nil {
		return
	}
	round := func(target **float64) {
		if target == nil || *target == nil {
			return
		}
		value := **target
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return
		}
		**target = math.Round(value*100) / 100
	}

	round(&m.WeightKg)
	round(&m.HeightMeter)
	round(&m.FatFreeMassKg)
	round(&m.FatRatioPercent)
	round(&m.FatMassKg)
	round(&m.DiastolicBpMmhg)
	round(&m.SystolicBpMmhg)
	round(&m.HeartPulseBpm)
	round(&m.TemperatureC)
	round(&m.Spo2Percent)
	round(&m.BodyTemperatureC)
	round(&m.SkinTemperatureC)
	round(&m.MuscleMassKg)
	round(&m.HydrationKg)
	round(&m.BoneMassKg)
	round(&m.PulseWaveVelocityMPerS)
	round(&m.Vo2MaxMlPerMinPerKg)
	round(&m.AtrialFibrillationResult)
	round(&m.QrsDurationMs)
	round(&m.PrDurationMs)
	round(&m.QtDurationMs)
	round(&m.QtCorrectedDurationMs)
	round(&m.AtrialFibrillationPpg)
	round(&m.VascularAgeYears)
	round(&m.NerveHealthConductanceFeet)
	round(&m.ExtracellularWaterKg)
	round(&m.IntracellularWaterKg)
	round(&m.VisceralFatIndex)
	round(&m.SegmentFatFreeMassKg)
	round(&m.SegmentFatMassKg)
	round(&m.SegmentMuscleMassKg)
	round(&m.ElectrodermalActivityFeet)
	round(&m.BasalMetabolicRate)
	round(&m.MetabolicAgeYears)
	round(&m.ElectrochemicalSkinConduct)
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

// #==============================================================#
// ##       ShouldRetryDailySummaryWithRefresh Process           ##
// #==============================================================#
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

// #==============================================================#
// ##       FetchDailySummary Process                            ##
// #==============================================================#
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

func truncateForLog(data []byte) string {
	const limit = 512
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) <= limit {
		return string(trimmed)
	}
	return string(trimmed[:limit]) + "..."
}

func joinInts(values []int) string {
	parts := make([]string, len(values))
	for i, v := range values {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}

func (s *HealthService) postFormJSON(ctx context.Context, path string, token string, form url.Values, apiName string) ([]byte, error) {
	resp, err := s.postForm(ctx, path, token, form)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("%s のレスポンス読み取りに失敗しました: %w", apiName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s がエラーを返しました: status=%d, body=%s", apiName, resp.StatusCode, truncateForLog(bodyBytes))
	}

	return bodyBytes, nil
}

func (s *HealthService) executePaginatedForm(
	ctx context.Context,
	endpointPath string,
	token string,
	form url.Values,
	apiName string,
	handler func([]byte) (bool, int, error),
) error {
	offset := 0
	loops := 0

	for {
		if loops >= maxPaginationLoops {
			return fmt.Errorf("%s のページネーション制限を超えました (%d)", apiName, maxPaginationLoops)
		}
		loops++

		if offset > 0 {
			form.Set("offset", strconv.Itoa(offset))
		} else {
			form.Del("offset")
		}

		bodyBytes, err := s.postFormJSON(ctx, endpointPath, token, form, apiName)
		if err != nil {
			return err
		}

		more, nextOffset, err := handler(bodyBytes)
		if err != nil {
			return err
		}
		if !more || nextOffset == 0 {
			break
		}
		offset = nextOffset
	}

	return nil
}

func (s *HealthService) convertMeasureValue(value int64, unit int) float64 {
	return float64(value) * math.Pow10(unit)
}

func (s *HealthService) labelForMeasureType(measureType int) string {
	if label, ok := measureTypeLabels[measureType]; ok {
		return label
	}
	return fmt.Sprintf("measure_%d", measureType)
}

func (s *HealthService) buildMeasureForm(userID int64, start, end time.Time, measureTypes []int) url.Values {
	form := url.Values{
		"action":    {"getmeas"},
		"userid":    {strconv.FormatInt(userID, 10)},
		"startdate": {strconv.FormatInt(start.Unix(), 10)},
		"enddate":   {strconv.FormatInt(end.Unix(), 10)},
		"category":  {"1"},
	}
	if len(measureTypes) > 0 {
		form.Set("meastypes", joinInts(measureTypes))
	}
	return form
}

func (s *HealthService) parseMeasureResponse(body []byte) (*measureAPIResponse, error) {
	var payload measureAPIResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("measure API の JSON 解析に失敗しました: %w", err)
	}
	if payload.Status != 0 {
		return nil, fmt.Errorf("measure API でエラーが発生しました: status=%d, error=%v", payload.Status, payload.Error)
	}
	return &payload, nil
}

func (s *HealthService) resolveTimezone(tzName string, current *time.Location, cachedName string) (*time.Location, string) {
	if tzName == "" || tzName == cachedName {
		return current, cachedName
	}
	if tz, err := time.LoadLocation(tzName); err == nil {
		return tz, tzName
	}
	return current, tzName
}

func (s *HealthService) ensureDailySummary(measurements map[string]*DailySummaryMeasures, dateKey string) *DailySummaryMeasures {
	summary := measurements[dateKey]
	if summary == nil {
		summary = &DailySummaryMeasures{}
		measurements[dateKey] = summary
	}
	return summary
}

func (s *HealthService) ensureLatestByDate(latestByDate map[string]map[int]time.Time, dateKey string) map[int]time.Time {
	latest := latestByDate[dateKey]
	if latest == nil {
		latest = make(map[int]time.Time)
		latestByDate[dateKey] = latest
	}
	return latest
}

func (s *HealthService) shouldSkipMeasurement(lastRecorded time.Time, recordedAt time.Time) bool {
	return !lastRecorded.IsZero() && !lastRecorded.Before(recordedAt)
}

func (s *HealthService) storeMeasurement(summary *DailySummaryMeasures, latest map[int]time.Time, item measureItem, recordedAt time.Time) {
	value := s.convertMeasureValue(item.Value, item.Unit)
	label := s.labelForMeasureType(item.Type)
	if !summary.set(label, value) {
		return
	}
	latest[item.Type] = recordedAt
}

func (s *HealthService) collectMeasureGroups(groups []measureGroup, loc *time.Location, measurements map[string]*DailySummaryMeasures, latestByDate map[string]map[int]time.Time) {
	for _, group := range groups {
		if group.Category != 1 {
			continue
		}
		recordedAt := time.Unix(group.Date, 0).In(loc)
		dateKey := recordedAt.Format("2006-01-02")
		summary := s.ensureDailySummary(measurements, dateKey)
		latestForDate := s.ensureLatestByDate(latestByDate, dateKey)

		for _, item := range group.Measures {
			if s.shouldSkipMeasurement(latestForDate[item.Type], recordedAt) {
				continue
			}
			s.storeMeasurement(summary, latestForDate, item, recordedAt)
		}
	}
}

func (s *HealthService) fetchMeasures(ctx context.Context, token string, userID int64, start, end time.Time, measureTypes []int) (*measureFetchResult, error) {
	result := &measureFetchResult{measurements: make(map[string]*DailySummaryMeasures)}
	latestByDate := make(map[string]map[int]time.Time)
	form := s.buildMeasureForm(userID, start, end, measureTypes)

	location := time.UTC
	cachedLocationName := ""

	handler := func(body []byte) (bool, int, error) {
		payload, err := s.parseMeasureResponse(body)
		if err != nil {
			return false, 0, err
		}

		if tzName := payload.Body.Timezone; tzName != "" {
			result.timezone = tzName
			location, cachedLocationName = s.resolveTimezone(tzName, location, cachedLocationName)
		}

		s.collectMeasureGroups(payload.Body.Measuregrps, location, result.measurements, latestByDate)

		return payload.Body.More.Bool(), payload.Body.Offset, nil
	}

	if err := s.executePaginatedForm(ctx, EndpointOfMeasure, token, form, "measure API", handler); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *HealthService) fetchActivities(ctx context.Context, token string, userID int64, start, end time.Time) (*activityFetchResult, error) {
	result := &activityFetchResult{
		activities: make(map[string]activityEntryWithTZ),
	}
	form := url.Values{}
	form.Set("action", "getactivity")
	form.Set("userid", strconv.FormatInt(userID, 10))
	form.Set("startdateymd", start.Format("2006-01-02"))
	form.Set("enddateymd", end.Format("2006-01-02"))

	err := s.executePaginatedForm(ctx, EndpointOfMeasureV2, token, form, "activity API", func(bodyBytes []byte) (bool, int, error) {
		var payload activityAPIResponse
		if err := json.Unmarshal(bodyBytes, &payload); err != nil {
			return false, 0, fmt.Errorf("activity API の JSON 解析に失敗しました: %w", err)
		}

		if payload.Status != 0 {
			return false, 0, fmt.Errorf("activity API でエラーが発生しました: status=%d, error=%v", payload.Status, payload.Error)
		}

		if result.defaultTimezone == "" && len(payload.Body.Activities) > 0 {
			result.defaultTimezone = payload.Body.Activities[0].Timezone
		}

		for _, activity := range payload.Body.Activities {
			summary := buildActivitySummary(activity)
			existing, ok := result.activities[activity.Date]
			tz := activity.Timezone
			if tz == "" && ok {
				tz = existing.timezone
			}
			result.activities[activity.Date] = activityEntryWithTZ{summary: summary, timezone: tz}
		}

		return payload.Body.More.Bool(), payload.Body.Offset, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func normalizeDateRange(start, end time.Time) (time.Time, time.Time) {
	startUTC := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endUTC := time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, time.UTC)
	return startUTC, endUTC
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

	for dateKey, measures := range measureResult.measurements {
		summaries[dateKey] = &DailySummary{
			Date:     dateKey,
			Timezone: measureResult.timezone,
			Measures: measures,
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

type measureFetchResult struct {
	timezone     string
	measurements map[string]*DailySummaryMeasures
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
	Status int `json:"status"`
	Error  any `json:"error"`
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
	Status int `json:"status"`
	Error  any `json:"error"`
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
