package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	config "github.com/landmaster135/devbox/internal/withings/config"
	usecases "github.com/landmaster135/devbox/internal/withings/usecases"
)

const (
	cliUserAgent       = "devbox-withings-cli/0.1"
	authorizeBaseURL   = "https://account.withings.com/oauth2_user/authorize2"
	oauthTokenEndpoint = "https://wbsapi.withings.net/v2/oauth2"
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

	healthService := usecases.NewServiceWithHTTPClient(cfg.BaseURL, httpClient)
	healthService.SetUserAgent(cliUserAgent)

	authService := usecases.NewAuthServiceWithEndpoints(httpClient, authorizeBaseURL, oauthTokenEndpoint)
	authService.WithRequestTimeout(cfg.Timeout)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout+5*time.Second)
	defer cancel()

	switch cfg.Operation {
	case config.OperationAuthURL:
		handleAuthURL(cfg, authService)
	case config.OperationRequestToken:
		handleRequestToken(ctx, authService, cfg)
	case config.OperationRefreshToken:
		handleRefreshToken(ctx, authService, cfg)
	case config.OperationDailySummary:
		handleDailySummary(ctx, healthService, cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の operation が指定されました: %s\n", cfg.Operation)
		os.Exit(1)
	}
}

func handleAuthURL(cfg *config.Config, authService *usecases.AuthService) {
	url, err := authService.BuildAuthorizationURL(usecases.AuthorizationURLParams{
		ClientID:     cfg.ClientID,
		RedirectURI:  cfg.RedirectURI,
		Scope:        cfg.Scope,
		State:        cfg.State,
		Mode:         cfg.Mode,
		ResponseType: cfg.ResponseType,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(url)
}

func handleRequestToken(ctx context.Context, authService *usecases.AuthService, cfg *config.Config) {
	resp, err := authService.ExchangeAuthorizationCode(ctx, usecases.TokenRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Code:         cfg.AuthorizationCode,
		RedirectURI:  cfg.RedirectURI,
		State:        cfg.State,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	emitTokenResponse(resp)
}

func handleRefreshToken(ctx context.Context, authService *usecases.AuthService, cfg *config.Config) {
	resp, err := authService.RefreshAccessToken(ctx, usecases.RefreshRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RefreshToken: cfg.RefreshToken,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	emitTokenResponse(resp)
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

func emitTokenResponse(resp *usecases.TokenResponse) {
	output := struct {
		Status int                        `json:"status"`
		Error  any                        `json:"error"`
		Body   usecases.TokenResponseBody `json:"body"`
	}{
		Status: resp.Status,
		Error:  resp.Error,
		Body:   resp.Body,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(output); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: トークン結果のエンコードに失敗しました: %v\n", err)
		os.Exit(1)
	}

	if resp.Body.UserID != "" {
		fmt.Fprintf(os.Stderr, "ユーザーID: %s\n", resp.Body.UserID)
	}
	if resp.Body.AccessToken != "" {
		fmt.Fprintf(os.Stderr, "アクセストークン: %s\n", resp.Body.AccessToken)
	}
	if resp.Body.RefreshToken != "" {
		fmt.Fprintf(os.Stderr, "リフレッシュトークンを安全な場所に保管してください: %s\n", resp.Body.RefreshToken)
	}
}
