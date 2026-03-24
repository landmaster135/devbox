package usecases

import (
	"context"
	"time"

	infrastructures "github.com/landmaster135/devbox/internal/memos/infrastructures"
	memos "github.com/landmaster135/devbox/internal/memos/usecases"
	common "github.com/landmaster135/devbox/internal/memos_utility/usecases/common"
	createClip "github.com/landmaster135/devbox/internal/memos_utility/usecases/operations/create_clip"
	createClips "github.com/landmaster135/devbox/internal/memos_utility/usecases/operations/create_clips"
	createCommonMemos "github.com/landmaster135/devbox/internal/memos_utility/usecases/operations/create_common_memos"
)

// ServiceOptions は Service 生成時の入力。
type ServiceOptions struct {
	BaseURL                           string
	APIToken                          string
	Timeout                           time.Duration
	MemosService                      MemosService
	FileSystem                        infrastructures.FileSystem
	CreateClipsProgressReporter       func(progress CreateClipsProgress)
	CreateCommonMemosProgressReporter func(progress CreateClipsProgress)
}

// Service は memos-utility のユースケースを提供する。
type Service struct {
	memosService        MemosService
	fileSystem          infrastructures.FileSystem
	createClipOp        createClipOperation
	createClipsOp       createClipsOperation
	createCommonMemosOp createCommonMemosOperation
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

	createClipOp := createClip.NewService(createClip.ServiceOptions{
		MemosService: memosService,
		FileSystem:   fileSystem,
	})
	createClipsOp := createClips.NewService(createClips.ServiceOptions{
		CreateClipService: createClipOp,
		FileSystem:        fileSystem,
		ProgressReporter:  opts.CreateClipsProgressReporter,
	})
	createCommonMemosOp := createCommonMemos.NewService(createCommonMemos.ServiceOptions{
		MemosService:     memosService,
		FileSystem:       fileSystem,
		ProgressReporter: opts.CreateCommonMemosProgressReporter,
	})

	return &Service{
		memosService:        memosService,
		fileSystem:          fileSystem,
		createClipOp:        createClipOp,
		createClipsOp:       createClipsOp,
		createCommonMemosOp: createCommonMemosOp,
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

// CreateCommonMemos は operation を委譲して共通メモを一括作成する。
func (s *Service) CreateCommonMemos(ctx context.Context, input CreateCommonMemosInput) (*CreateCommonMemosOutput, error) {
	return s.createCommonMemosOp.Execute(ctx, common.CreateCommonMemosInput(input))
}
