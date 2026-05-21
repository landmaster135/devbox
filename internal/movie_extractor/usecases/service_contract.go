package usecases

import (
	commandExecutor "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/command_executor"
	dedupImagesOperation "github.com/landmaster135/devbox/internal/movie_extractor/usecases/operations/dedup_images"
	extractFramesOperation "github.com/landmaster135/devbox/internal/movie_extractor/usecases/operations/extract_frames"
)

type extractFramesService interface {
	Handle(input extractFramesOperation.Input) (string, error)
}

type dedupImagesService interface {
	Handle(input dedupImagesOperation.Input) (string, error)
}

func newDefaultOperations(executor commandExecutor.Repository) (extractFramesService, dedupImagesService) {
	return extractFramesOperation.NewService(executor), dedupImagesOperation.NewService()
}
