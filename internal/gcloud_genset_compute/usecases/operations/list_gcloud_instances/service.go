package listgcloudinstances

import (
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases/common"
)

// Params はインスタンス一覧コマンド生成に必要な値。
type Params struct {
	Filter string
	Format string
}

// Service は list-gcloud-instances operation のコマンド生成を担当する。
type Service struct{}

// NewService は Service のインスタンスを返す。
func NewService() *Service {
	return &Service{}
}

// Build はインスタンス一覧コマンドを生成する。
func (s *Service) Build(params Params) (string, error) {
	if common.IsBlank(params.Format) {
		return "", fmt.Errorf("format は必須です")
	}

	parts := []string{"gcloud", "compute", "instances", "list"}
	if !common.IsBlank(params.Filter) {
		parts = append(parts, "--filter="+common.ShellQuote(params.Filter))
	}
	parts = append(parts, "--format="+common.ShellQuote(params.Format))

	return strings.Join(parts, " "), nil
}
