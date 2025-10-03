package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	config "github.com/landmaster135/devbox/internal/withings/config"
	usecases "github.com/landmaster135/devbox/internal/withings/usecases"
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

	healthService := usecases.NewHealthService(cfg.Timeout)
	authService := usecases.NewAuthService(cfg.Timeout)

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
		handleDailySummary(ctx, healthService, authService, cfg)
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

func handleDailySummary(ctx context.Context, healthService *usecases.HealthService, authService *usecases.AuthService, cfg *config.Config) {
	req := usecases.DailySummaryRequest{
		AccessToken:     cfg.AccessToken,
		UserID:          cfg.UserID,
		StartDate:       cfg.StartDate,
		EndDate:         cfg.EndDate,
		MeasureTypes:    cfg.MeasureTypes,
		IncludeActivity: cfg.IncludeActivity,
	}

	resp, err := healthService.FetchDailySummary(ctx, req)
	if err != nil {
		if healthService.ShouldRetryDailySummaryWithRefresh(err) {
			resp, err = retryDailySummaryWithRefresh(ctx, healthService, authService, cfg, req, err)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(resp); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 出力のエンコードに失敗しました: %v\n", err)
		os.Exit(1)
	}
}

func retryDailySummaryWithRefresh(
	ctx context.Context,
	healthService *usecases.HealthService,
	authService *usecases.AuthService,
	cfg *config.Config,
	req usecases.DailySummaryRequest,
	originalErr error,
) (*usecases.DailySummaryResponse, error) {
	if cfg.RefreshToken == "" {
		return nil, originalErr
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("日次サマリの取得に失敗しました: %v\nリフレッシュトークンで再試行するには client-id と client-secret を指定してください", originalErr)
	}

	fmt.Fprintf(os.Stderr, "警告: アクセストークンでの取得に失敗したため refresh-token で再試行します...\n")
	refreshResp, err := authService.RefreshAccessToken(ctx, usecases.RefreshRequest{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RefreshToken: cfg.RefreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("アクセストークンのリフレッシュに失敗しました: %w\n元のエラー: %v", err, originalErr)
	}
	if refreshResp == nil || refreshResp.Body.AccessToken == "" {
		return nil, fmt.Errorf("リフレッシュレスポンスにアクセストークンが含まれていません\n元のエラー: %v", originalErr)
	}

	cfg.AccessToken = refreshResp.Body.AccessToken
	req.AccessToken = refreshResp.Body.AccessToken

	fmt.Fprintf(os.Stderr, "情報: 新しいアクセストークンで日次サマリ取得を再試行します。\nアクセストークン: %s\n", cfg.AccessToken)
	if refreshResp.Body.RefreshToken != "" {
		cfg.RefreshToken = refreshResp.Body.RefreshToken
		fmt.Fprintf(os.Stderr, "新しいリフレッシュトークンを安全な場所に保管してください: %s\n", cfg.RefreshToken)
	}

	resp, err := healthService.FetchDailySummary(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("リフレッシュ後の日次サマリ取得にも失敗しました: %v\n元のエラー: %v", err, originalErr)
	}

	return resp, nil
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
