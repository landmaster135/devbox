package usecases

import (
	"context"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/withings/config"
)

// CoreService は Withings CLI が用いる横断的なユースケースをまとめます。
type CoreService struct {
	healthService *HealthService
	authService   *AuthService
}

// NewCoreService はヘルス系・認可系サービスを受け取り CoreService を構築します。
func NewCoreService(healthService *HealthService, authService *AuthService) *CoreService {
	return &CoreService{
		healthService: healthService,
		authService:   authService,
	}
}

// GetHealthService は内部で保持している GetHealthService を返します。
func (s *CoreService) GetHealthService() *HealthService {
	return s.healthService
}

// GetAuthService は内部で保持している GetAuthService を返します。
func (s *CoreService) GetAuthService() *AuthService {
	return s.authService
}

// RetryDailySummaryWithRefresh はアクセストークンでの日次サマリ取得が失敗した際に
// リフレッシュトークンを用いて再試行します。
func (s *CoreService) RetryDailySummaryWithRefresh(
	ctx context.Context,
	cfg *config.Config,
	req DailySummaryRequest,
	originalErr error,
) (*DailySummaryResponse, error) {
	if cfg == nil {
		return nil, originalErr
	}
	if cfg.RefreshToken == "" {
		return nil, originalErr
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf(
			"日次サマリの取得に失敗しました: %v\nリフレッシュトークンで再試行するには client-id と client-secret を指定してください",
			originalErr,
		)
	}

	fmt.Fprintf(os.Stderr, "警告: アクセストークンでの取得に失敗したため refresh-token で再試行します...\n")
	refreshResp, err := s.authService.RefreshAccessToken(ctx, RefreshRequest{
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

	resp, err := s.healthService.FetchDailySummary(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("リフレッシュ後の日次サマリ取得にも失敗しました: %v\n元のエラー: %v", err, originalErr)
	}

	return resp, nil
}
