package memos

import (
	"context"
	"fmt"
	"strings"
	"time"

	memosusecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

const defaultTimeout = 30 * time.Second

// Client は notion_to_memos_markdown から利用する最小の Memos API 契約。
type Client interface {
	CreateMemo(ctx context.Context, content string) (string, error)
	PatchFiles(ctx context.Context, memo string, filePaths []string) error
}

type apiClient struct {
	service memosusecases.MemoService
}

func NewClient(baseURL, apiToken string) Client {
	service := memosusecases.NewService(memosusecases.ServiceOptions{
		BaseURL:  strings.TrimSpace(baseURL),
		APIToken: strings.TrimSpace(apiToken),
		Timeout:  defaultTimeout,
	})
	return NewClientWithService(service)
}

func NewClientWithService(service memosusecases.MemoService) Client {
	return &apiClient{
		service: service,
	}
}

func (c *apiClient) CreateMemo(ctx context.Context, content string) (string, error) {
	result, err := c.service.CreateMemo(ctx, "", content, "", "", "", nil, "")
	if err != nil {
		return "", fmt.Errorf("メモ作成に失敗しました: %w", err)
	}
	if result == nil {
		return "", fmt.Errorf("メモ作成結果が空です")
	}
	memoName := strings.TrimSpace(result.Name)
	if memoName == "" {
		return "", fmt.Errorf("メモ作成結果に memo 名がありません")
	}
	return memoName, nil
}

func (c *apiClient) PatchFiles(ctx context.Context, memo string, filePaths []string) error {
	_, err := c.service.PatchFiles(ctx, memo, filePaths, false)
	if err != nil {
		return fmt.Errorf("メモ添付に失敗しました: %w", err)
	}
	return nil
}
