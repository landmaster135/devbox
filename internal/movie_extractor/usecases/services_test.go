package usecases

import (
	"bytes"
	"testing"

	dedupImagesOperation "github.com/landmaster135/devbox/internal/movie_extractor/usecases/operations/dedup_images"
	extractFramesOperation "github.com/landmaster135/devbox/internal/movie_extractor/usecases/operations/extract_frames"
)

type extractFramesServiceStub struct {
	input  extractFramesOperation.Input
	called bool
}

func (s *extractFramesServiceStub) Handle(input extractFramesOperation.Input) (string, error) {
	s.called = true
	s.input = input
	return "ok", nil
}

type dedupImagesServiceStub struct {
	input  dedupImagesOperation.Input
	called bool
}

func (s *dedupImagesServiceStub) Handle(input dedupImagesOperation.Input) (string, error) {
	s.called = true
	s.input = input
	return "ok", nil
}

func TestHandleExtractFrames_DelegatesToOperation(t *testing.T) {
	extractStub := &extractFramesServiceStub{}
	dedupStub := &dedupImagesServiceStub{}
	service := newServiceWithOperations(extractStub, dedupStub)

	input := ExtractFramesInput{
		SrcFile:       "input.mp4",
		FPS:           3,
		Quality:       2,
		StartPosition: "00:00:01.5",
		OutDir:        "frames",
	}
	if _, err := service.HandleExtractFrames(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !extractStub.called {
		t.Fatal("extract operation was not called")
	}
	if dedupStub.called {
		t.Fatal("dedup operation should not be called")
	}

	if extractStub.input.SrcFile != input.SrcFile {
		t.Fatalf("unexpected src-file: %s", extractStub.input.SrcFile)
	}
	if extractStub.input.FPS != input.FPS {
		t.Fatalf("unexpected fps: %d", extractStub.input.FPS)
	}
	if extractStub.input.Quality != input.Quality {
		t.Fatalf("unexpected quality: %d", extractStub.input.Quality)
	}
	if extractStub.input.StartPosition != input.StartPosition {
		t.Fatalf("unexpected start-position: %s", extractStub.input.StartPosition)
	}
	if extractStub.input.OutDir != input.OutDir {
		t.Fatalf("unexpected out-dir: %s", extractStub.input.OutDir)
	}
}

func TestHandleDedupImages_DelegatesToOperation(t *testing.T) {
	extractStub := &extractFramesServiceStub{}
	dedupStub := &dedupImagesServiceStub{}
	service := newServiceWithOperations(extractStub, dedupStub)

	logWriter := &bytes.Buffer{}
	input := DedupImagesInput{
		SrcDir:    "images",
		MatchRate: 98,
		Log:       true,
		LogWriter: logWriter,
		OutDir:    "unique-images",
	}
	if _, err := service.HandleDedupImages(input); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !dedupStub.called {
		t.Fatal("dedup operation was not called")
	}
	if extractStub.called {
		t.Fatal("extract operation should not be called")
	}

	if dedupStub.input.SrcDir != input.SrcDir {
		t.Fatalf("unexpected src-dir: %s", dedupStub.input.SrcDir)
	}
	if dedupStub.input.MatchRate != input.MatchRate {
		t.Fatalf("unexpected match-rate: %f", dedupStub.input.MatchRate)
	}
	if dedupStub.input.Log != input.Log {
		t.Fatalf("unexpected log flag: %v", dedupStub.input.Log)
	}
	if dedupStub.input.LogWriter != input.LogWriter {
		t.Fatal("log writer was not passed through")
	}
	if dedupStub.input.OutDir != input.OutDir {
		t.Fatalf("unexpected out-dir: %s", dedupStub.input.OutDir)
	}
}
