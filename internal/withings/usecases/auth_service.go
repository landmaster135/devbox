package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultAuthorizeBaseURL   = "https://account.withings.com/oauth2_user/authorize2"
	defaultOAuthTokenEndpoint = "https://wbsapi.withings.net/v2/oauth2"
	defaultOAuthTimeout       = 30 * time.Second
)

// AuthService は Withings OAuth フローを操作する機能を提供します。
type AuthService struct {
	httpClient     HTTPClient
	authorizeURL   string
	tokenEndpoint  string
	requestTimeout time.Duration
}

// NewAuthService は指定されたタイムアウトで初期化された AuthService を返します。
func NewAuthService(timeout time.Duration) *AuthService {
	return NewAuthServiceWithClient(nil, timeout)
}

// NewAuthServiceWithClient は HTTP クライアントを差し替えて AuthService を返します。
func NewAuthServiceWithClient(client HTTPClient, timeout time.Duration) *AuthService {
	timeout = normalizeOAuthTimeout(timeout)
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	svc := &AuthService{
		httpClient:     client,
		authorizeURL:   defaultAuthorizeBaseURL,
		tokenEndpoint:  defaultOAuthTokenEndpoint,
		requestTimeout: timeout,
	}
	svc.SetEndpoints(defaultAuthorizeBaseURL, defaultOAuthTokenEndpoint)
	return svc
}

// NewAuthServiceWithEndpoints はベースURLを差し替えて AuthService を構築します。
func NewAuthServiceWithEndpoints(timeout time.Duration, authorizeURL, tokenEndpoint string) *AuthService {
	svc := NewAuthService(timeout)
	svc.SetEndpoints(authorizeURL, tokenEndpoint)
	return svc
}

// SetEndpoints は認可URLとトークンエンドポイントを上書きします。
func (s *AuthService) SetEndpoints(authorizeURL, tokenEndpoint string) {
	if trimmed := strings.TrimSpace(authorizeURL); trimmed != "" {
		s.authorizeURL = trimmed
	}
	if trimmed := strings.TrimSpace(tokenEndpoint); trimmed != "" {
		s.tokenEndpoint = trimmed
	}
}

// AuthorizationURLParams は認可URL生成に必要なパラメータを表します。
type AuthorizationURLParams struct {
	ClientID     string
	RedirectURI  string
	Scope        string
	State        string
	Mode         string
	ResponseType string
}

// BuildAuthorizationURL は認可URLを生成します。
func (s *AuthService) BuildAuthorizationURL(params AuthorizationURLParams) (string, error) {
	if strings.TrimSpace(params.ClientID) == "" {
		return "", errors.New("client_id が指定されていません")
	}
	if strings.TrimSpace(params.RedirectURI) == "" {
		return "", errors.New("redirect_uri が指定されていません")
	}

	values := url.Values{}
	responseType := strings.TrimSpace(params.ResponseType)
	if responseType == "" {
		responseType = "code"
	}
	values.Set("response_type", responseType)
	values.Set("client_id", strings.TrimSpace(params.ClientID))
	values.Set("redirect_uri", strings.TrimSpace(params.RedirectURI))

	scope := normalizeScope(params.Scope)
	if scope != "" {
		values.Set("scope", scope)
	}
	if params.State != "" {
		values.Set("state", params.State)
	}
	if params.Mode != "" {
		values.Set("mode", params.Mode)
	}

	return s.authorizeURL + "?" + values.Encode(), nil
}

// TokenRequest は認可コード交換に必要なパラメータです。
type TokenRequest struct {
	ClientID     string
	ClientSecret string
	Code         string
	RedirectURI  string
	State        string
}

// TokenResponse は Withings OAuth endpoint の結果を表します。
type TokenResponse struct {
	Status     int               `json:"status"`
	Error      any               `json:"error"`
	Body       TokenResponseBody `json:"body"`
	RawBody    []byte            `json:"-"`
	StatusCode int               `json:"-"`
}

// TokenResponseBody はレスポンス body の内容です。
type TokenResponseBody struct {
	UserID       string `json:"userid"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// ExchangeAuthorizationCode は認可コードからアクセストークンを取得します。
func (s *AuthService) ExchangeAuthorizationCode(ctx context.Context, req TokenRequest) (*TokenResponse, error) {
	if err := validateTokenRequest(req); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("action", "requesttoken")
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", req.ClientID)
	form.Set("client_secret", req.ClientSecret)
	form.Set("code", req.Code)
	form.Set("redirect_uri", req.RedirectURI)
	if req.State != "" {
		form.Set("state", req.State)
	}

	return s.doTokenRequest(ctx, form)
}

// RefreshAccessToken はリフレッシュトークンでアクセストークンを再取得します。
type RefreshRequest struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, req RefreshRequest) (*TokenResponse, error) {
	if strings.TrimSpace(req.ClientID) == "" || strings.TrimSpace(req.ClientSecret) == "" {
		return nil, errors.New("client_id と client_secret を指定してください")
	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		return nil, errors.New("refresh_token が指定されていません")
	}

	form := url.Values{}
	form.Set("action", "requesttoken")
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", strings.TrimSpace(req.ClientID))
	form.Set("client_secret", strings.TrimSpace(req.ClientSecret))
	form.Set("refresh_token", strings.TrimSpace(req.RefreshToken))

	return s.doTokenRequest(ctx, form)
}

func validateTokenRequest(req TokenRequest) error {
	if strings.TrimSpace(req.ClientID) == "" || strings.TrimSpace(req.ClientSecret) == "" {
		return errors.New("client_id と client_secret を指定してください")
	}
	if strings.TrimSpace(req.Code) == "" {
		return errors.New("authorization code が指定されていません")
	}
	if strings.TrimSpace(req.RedirectURI) == "" {
		return errors.New("redirect_uri が指定されていません")
	}
	return nil
}

func (s *AuthService) doTokenRequest(ctx context.Context, form url.Values) (*TokenResponse, error) {
	timeout := s.requestTimeout
	if timeout <= 0 {
		timeout = defaultOAuthTimeout
	}
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, s.tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("OAuth リクエストの構築に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OAuth リクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("OAuth レスポンスの読み取りに失敗しました: %w", err)
	}

	tokenResp := &TokenResponse{RawBody: bodyBytes, StatusCode: resp.StatusCode}
	if resp.StatusCode != http.StatusOK {
		tokenResp.Status = resp.StatusCode
		tokenResp.Error = string(bytes.TrimSpace(bodyBytes))
		return tokenResp, fmt.Errorf("OAuth エンドポイントがエラーを返しました: status=%d body=%s", resp.StatusCode, truncateForOAuthLog(bodyBytes))
	}

	type rawResponse struct {
		Status int             `json:"status"`
		Error  any             `json:"error"`
		Body   json.RawMessage `json:"body"`
	}

	var raw rawResponse
	if err := json.Unmarshal(bodyBytes, &raw); err != nil {
		return nil, fmt.Errorf("OAuth レスポンスの解析に失敗しました: %w", err)
	}

	tokenResp.Status = raw.Status
	tokenResp.Error = raw.Error
	if raw.Status != 0 {
		return tokenResp, fmt.Errorf("OAuth エンドポイントがエラーを返しました: status=%d error=%v", raw.Status, raw.Error)
	}

	var body struct {
		UserID       json.RawMessage `json:"userid"`
		AccessToken  string          `json:"access_token"`
		RefreshToken string          `json:"refresh_token"`
		ExpiresIn    int             `json:"expires_in"`
		TokenType    string          `json:"token_type"`
		Scope        string          `json:"scope"`
	}
	if err := json.Unmarshal(raw.Body, &body); err != nil {
		return nil, fmt.Errorf("OAuth body の解析に失敗しました: %w", err)
	}

	userID, err := parseFlexibleString(body.UserID)
	if err != nil {
		return nil, fmt.Errorf("userid の解析に失敗しました: %w", err)
	}

	tokenResp.Body = TokenResponseBody{
		UserID:       userID,
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
		ExpiresIn:    body.ExpiresIn,
		TokenType:    body.TokenType,
		Scope:        body.Scope,
	}

	return tokenResp, nil
}

func parseFlexibleString(raw json.RawMessage) (string, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return "", nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return "", err
		}
		return s, nil
	}
	var num json.Number
	if err := json.Unmarshal(trimmed, &num); err != nil {
		return "", err
	}
	return num.String(), nil
}

func truncateForOAuthLog(data []byte) string {
	const limit = 512
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) <= limit {
		return string(trimmed)
	}
	return string(trimmed[:limit]) + "..."
}

func normalizeScope(scope string) string {
	trimmed := strings.TrimSpace(scope)
	if trimmed == "" {
		return ""
	}
	tokens := strings.Split(trimmed, ",")
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t != "" {
			result = append(result, t)
		}
	}
	return strings.Join(result, ",")
}

func normalizeOAuthTimeout(timeout time.Duration) time.Duration {
	if timeout <= 0 {
		return defaultOAuthTimeout
	}
	return timeout
}

// ExportDailySummary は日次サマリの結果を Withings Health Mate 形式の JSON として保存します。
func (s *AuthService) ExportDailySummary(resp *DailySummaryResponse, path string) error {
	if resp == nil {
		return fmt.Errorf("daily summary レスポンスが空です")
	}
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("出力ファイルパスが指定されていません")
	}

	entries := buildHealthMateEntries(resp)
	export := healthMateExport{
		Data:        healthMateExportData{HealthMates: entries},
		Description: "Health Mate data from Withings",
		Name:        "My Health Mate Data",
	}

	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("ディレクトリの作成に失敗しました: %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ファイルの作成に失敗しました: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(export); err != nil {
		return fmt.Errorf("JSON のエンコードに失敗しました: %w", err)
	}

	return nil
}

type healthMateExport struct {
	Data        healthMateExportData `json:"data"`
	Description string               `json:"description"`
	Name        string               `json:"name"`
}

type healthMateExportData struct {
	HealthMates []map[string]any `json:"health_mates"`
}

func buildHealthMateEntries(resp *DailySummaryResponse) []map[string]any {
	if resp == nil {
		return nil
	}

	entries := make([]map[string]any, 0, len(resp.Summaries))
	for _, summary := range resp.Summaries {
		entry := make(map[string]any)
		if summary.Date != "" {
			entry["date"] = summary.Date
		}
		for key, value := range summary.Measures {
			entry[key] = value
		}
		if summary.Activity != nil {
			for key, value := range activitySummaryToMap(summary.Activity) {
				entry[key] = value
			}
		}
		entries = append(entries, entry)
	}

	return entries
}

func activitySummaryToMap(activity *ActivitySummary) map[string]any {
	if activity == nil {
		return nil
	}

	result := make(map[string]any)
	if activity.Steps != nil {
		result["steps"] = *activity.Steps
	}
	if activity.DistanceMeter != nil {
		result["distance_meter"] = *activity.DistanceMeter
	}
	if activity.ElevationMeter != nil {
		result["elevation_meter"] = *activity.ElevationMeter
	}
	if activity.CaloriesKcal != nil {
		result["calories_kcal"] = *activity.CaloriesKcal
	}
	if activity.TotalCaloriesKcal != nil {
		result["total_calories_kcal"] = *activity.TotalCaloriesKcal
	}
	if activity.SoftSeconds != nil {
		result["soft_seconds"] = *activity.SoftSeconds
	}
	if activity.ModerateSeconds != nil {
		result["moderate_seconds"] = *activity.ModerateSeconds
	}
	if activity.IntenseSeconds != nil {
		result["intense_seconds"] = *activity.IntenseSeconds
	}
	if activity.ActiveSeconds != nil {
		result["active_seconds"] = *activity.ActiveSeconds
	}
	if activity.HrAverageBPM != nil {
		result["hr_average_bpm"] = *activity.HrAverageBPM
	}
	if activity.HrMinBPM != nil {
		result["hr_min_bpm"] = *activity.HrMinBPM
	}
	if activity.HrMaxBPM != nil {
		result["hr_max_bpm"] = *activity.HrMaxBPM
	}
	if activity.DeviceBrand != nil {
		result["device_brand"] = *activity.DeviceBrand
	}
	if activity.DeviceModelID != nil {
		result["device_model_id"] = *activity.DeviceModelID
	}
	if activity.DeviceModelName != nil {
		result["device_model_name"] = *activity.DeviceModelName
	}
	if activity.IsTracker != nil {
		result["is_tracker"] = *activity.IsTracker
	}

	return result
}
