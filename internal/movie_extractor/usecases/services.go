package usecases

import (
	"io"

	commandExecutor "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/command_executor"
	dedupImagesOperation "github.com/landmaster135/devbox/internal/movie_extractor/usecases/operations/dedup_images"
	extractFramesOperation "github.com/landmaster135/devbox/internal/movie_extractor/usecases/operations/extract_frames"
)

// ExtractFramesInput は extract-frames 操作の入力です。
type ExtractFramesInput struct {
	SrcFile       string
	FPS           int
	Quality       int
	StartPosition string
	OutDir        string
}

// DedupImagesInput は dedup-images 操作の入力です。
type DedupImagesInput struct {
	SrcDir    string
	MatchRate float64
	Log       bool
	LogWriter io.Writer
	OutDir    string
}

// Service は movie extractor のユースケースです。
type Service struct {
	extractFrames extractFramesService
	dedupImages   dedupImagesService
}

// NewService は標準依存で Service を生成します。
func NewService() *Service {
	return NewServiceWithExecutor(commandExecutor.NewOSCommandExecutor())
}

// NewServiceWithExecutor はテスト用の依存注入を行います。
func NewServiceWithExecutor(executor commandExecutor.Repository) *Service {
	if executor == nil {
		executor = commandExecutor.NewOSCommandExecutor()
	}
	extractFrames, dedupImages := newDefaultOperations(executor)
	return &Service{
		extractFrames: extractFrames,
		dedupImages:   dedupImages,
	}
}

func newServiceWithOperations(extractFrames extractFramesService, dedupImages dedupImagesService) *Service {
	return &Service{
		extractFrames: extractFrames,
		dedupImages:   dedupImages,
	}
}

// HandleExtractFrames は動画からフレーム画像を抽出します。
func (s *Service) HandleExtractFrames(input ExtractFramesInput) (string, error) {
	return s.extractFrames.Handle(extractFramesOperation.Input{
		SrcFile:       input.SrcFile,
		FPS:           input.FPS,
		Quality:       input.Quality,
		StartPosition: input.StartPosition,
		OutDir:        input.OutDir,
	})
}

// HandleDedupImages は画像ディレクトリから重複画像を除去します。
func (s *Service) HandleDedupImages(input DedupImagesInput) (string, error) {
	return s.dedupImages.Handle(dedupImagesOperation.Input{
		SrcDir:    input.SrcDir,
		MatchRate: input.MatchRate,
		Log:       input.Log,
		LogWriter: input.LogWriter,
		OutDir:    input.OutDir,
	})
}
