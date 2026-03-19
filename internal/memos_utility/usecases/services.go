package usecases

import (
	"context"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
	createclip "github.com/landmaster135/devbox/internal/memos_utility/usecases/operations/create_clip"
	createclips "github.com/landmaster135/devbox/internal/memos_utility/usecases/operations/create_clips"
)

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	BaseURL      string
	APIToken     string
	Timeout      time.Duration
	MemosService MemosService
	FileSystem   infrastructures.FileSystem
}

// Service は memos-utility のユースケースを提供する。
type Service struct {
	memosService  MemosService
	fileSystem    infrastructures.FileSystem
	createClipOp  createClipOperation
	createClipsOp createClipsOperation
}

// NewService は Service を生成する。
func NewService(opts ServiceOptions) *Service {
	memosService := opts.MemosService
	if memosService == nil {
		memosService = memos.NewService(memos.ServiceOptions{
			BaseURL:  opts.BaseURL,
			APIToken: opts.APIToken,
			Timeout:  opts.Timeout,
		})
	}

	fileSystem := opts.FileSystem
	if fileSystem == nil {
		fileSystem = infrastructures.NewOSFileSystem()
	}

	createClipOp := createclip.NewService(createclip.ServiceOptions{
		MemosService: memosService,
		FileSystem:   fileSystem,
	})
	createClipsOp := createclips.NewService(createclips.ServiceOptions{
		CreateClipService: createClipOp,
		FileSystem:        fileSystem,
	})

	return &Service{
		memosService:  memosService,
		fileSystem:    fileSystem,
		createClipOp:  createClipOp,
		createClipsOp: createClipsOp,
	}
}

// CreateClip は operation を委譲して単一クリップを作成する。
func (s *Service) CreateClip(ctx context.Context, input CreateClipInput) (*CreateClipOutput, error) {
	return s.createClipOp.Execute(ctx, common.CreateClipInput(input))
}

// CreateClips は operation を委譲して一括クリップ作成する。
func (s *Service) CreateClips(ctx context.Context, input CreateClipsInput) (*CreateClipsOutput, error) {
	return s.createClipsOp.Execute(ctx, common.CreateClipsInput(input))
}
