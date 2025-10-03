package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	config "github.com/landmaster135/devbox/internal/withings/config"
	usecases "github.com/landmaster135/devbox/internal/withings/usecases"
)

const (
	cliUserAgent        = "devbox-withings-cli/0.1"
	authorizeBaseURL    = "https://account.withings.com/oauth2_user/authorize2"
	oauthTokenEndpoint  = "https://wbsapi.withings.net/v2/oauth2"
	defaultOAuthTimeout = 30 * time.Second
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	httpClient := &http.Client{Timeout: cfg.Timeout}
	service := usecases.NewServiceWithHTTPClient(cfg.BaseURL, httpClient)
	service.SetUserAgent(cliUserAgent)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout+5*time.Second)
	defer cancel()

	switch cfg.Operation {
	case config.OperationAuthURL:
		handleAuthURL(cfg)
	case config.OperationRequestToken:
		handleRequestToken(ctx, httpClient, cfg)
	case config.OperationRefreshToken:
		handleRefreshToken(ctx, httpClient, cfg)
	case config.OperationDailySummary:
		handleDailySummary(ctx, service, cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の operation が指定されました: %s\n", cfg.Operation)
		os.Exit(1)
	}
}

func handleAuthURL(cfg *config.Config) {
	values := url.Values{}
	values.Set("response_type", firstNonEmpty(cfg.ResponseType, "code"))
	values.Set("client_id", cfg.ClientID)
	values.Set("redirect_uri", cfg.RedirectURI)
	values.Set("scope", cfg.Scope)
	if cfg.State != "" {
		values.Set("state", cfg.State)
	}
	if cfg.Mode != "" {
		values.Set("mode", cfg.Mode)
	}

	authURL := authorizeBaseURL + "?" + values.Encode()
	fmt.Println(authURL)
}

func handleRequestToken(ctx context.Context, httpClient *http.Client, cfg *config.Config) {
	form := url.Values{}
	form.Set("action", "requesttoken")
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("code", cfg.AuthorizationCode)
	form.Set("redirect_uri", cfg.RedirectURI)

	payload, err := executeOAuthRequest(ctx, httpClient, form)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	emitTokenResponse(payload)
}

func handleRefreshToken(ctx context.Context, httpClient *http.Client, cfg *config.Config) {
	form := url.Values{}
	form.Set("action", "refresh")
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", cfg.ClientID)
	form.Set("client_secret", cfg.ClientSecret)
	form.Set("refresh_token", cfg.RefreshToken)

	payload, err := executeOAuthRequest(ctx, httpClient, form)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	emitTokenResponse(payload)
}

func handleDailySummary(ctx context.Context, service *usecases.Service, cfg *config.Config) {
	req := usecases.DailySummaryRequest{
		AccessToken:     cfg.AccessToken,
		UserID:          cfg.UserID,
		StartDate:       cfg.StartDate,
		EndDate:         cfg.EndDate,
		MeasureTypes:    cfg.MeasureTypes,
		IncludeActivity: cfg.IncludeActivity,
	}

	resp, err := service.FetchDailySummary(ctx, req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 出力のエンコードに失敗しました: %v\n", err)
		os.Exit(1)
	}
}

func executeOAuthRequest(ctx context.Context, httpClient *http.Client, form url.Values) (*tokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthTokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("OAuth リクエストの構築に失敗しました: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	oauthClient := httpClient
	if oauthClient.Timeout < defaultOAuthTimeout {
		oauthClient = &http.Client{Timeout: defaultOAuthTimeout}
	}

	resp, err := oauthClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OAuth リクエストの送信に失敗しました: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("OAuth レスポンスの読み取りに失敗しました: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OAuth エンドポイントがエラーを返しました: status=%d body=%s", resp.StatusCode, truncateForLog(bodyBytes))
	}

	var payload tokenResponse
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		return nil, fmt.Errorf("OAuth レスポンスの解析に失敗しました: %w", err)
	}

	if payload.Status != 0 {
		return nil, fmt.Errorf("OAuth エンドポイントがエラーを返しました: status=%d error=%v", payload.Status, payload.Error)
	}

	return &payload, nil
}

func emitTokenResponse(payload *tokenResponse) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: トークン結果のエンコードに失敗しました: %v\n", err)
		os.Exit(1)
	}

	if payload.Body.UserID != "" {
		fmt.Fprintf(os.Stderr, "ユーザーID: %s\n", payload.Body.UserID)
	}

	if payload.Body.AccessToken != "" {
		fmt.Fprintf(os.Stderr, "アクセストークン: %s\n", payload.Body.AccessToken)
	}

	if payload.Body.RefreshToken != "" {
		fmt.Fprintf(os.Stderr, "リフレッシュトークンを安全な場所に保管してください: %s\n", payload.Body.RefreshToken)
	}
}

func truncateForLog(data []byte) string {
	const limit = 512
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) <= limit {
		return string(trimmed)
	}
	return string(trimmed[:limit]) + "..."
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

type tokenResponse struct {
	Status int         `json:"status"`
	Error  any `json:"error"`
	Body   struct {
		UserID       string `json:"userid"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	} `json:"body"`
}
