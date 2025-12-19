package usecases

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/image_filterer_v2/config"
)

type processor interface {
	Process() (string, error)
}

// Service はフィルタ処理のエントリーポイントを表します。
type Service struct {
	processor processor
}

// NewService は環境に応じた実装を選択して初期化します。
func NewService(cfg config.Config) (*Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg.Normalise()

	proc, err := newProcessor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare processor: %w", err)
	}

	return &Service{processor: proc}, nil
}

// Process は選択された実装に処理を委譲します。
func (s *Service) Process() (string, error) {
	if s.processor == nil {
		return "", fmt.Errorf("processor is not initialized")
	}
	return s.processor.Process()
}
