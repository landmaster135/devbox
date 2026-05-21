package usecases

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	commandExecutor "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/command_executor"
)

func TestHandleExtractFrames_Success(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "input.mp4")
	if err := os.WriteFile(srcFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	outDir := filepath.Join(tempDir, "frames")
	mock := &commandExecutor.MockRepository{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("ok"), nil
		},
	}

	service := NewServiceWithExecutor(mock)
	result, err := service.HandleExtractFrames(ExtractFramesInput{
		SrcFile:       srcFile,
		FPS:           3,
		Quality:       2,
		StartPosition: "00:00:01.5",
		OutDir:        outDir,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.Calls) != 1 {
		t.Fatalf("expected 1 command call, got %d", len(mock.Calls))
	}
	call := mock.Calls[0]
	if call.Name != "ffmpeg" {
		t.Fatalf("unexpected command: %s", call.Name)
	}

	argsText := strings.Join(call.Args, " ")
	if !strings.Contains(argsText, "-ss 00:00:01.5") {
		t.Fatalf("start-position argument missing: %s", argsText)
	}
	if !strings.Contains(argsText, "-vf fps=3") {
		t.Fatalf("fps argument missing: %s", argsText)
	}
	if !strings.Contains(argsText, "-q:v 2") {
		t.Fatalf("quality argument missing: %s", argsText)
	}

	if _, err := os.Stat(outDir); err != nil {
		t.Fatalf("out-dir should be created: %v", err)
	}
	if !strings.Contains(result, "フレーム抽出が完了しました") {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestHandleExtractFrames_InputValidation(t *testing.T) {
	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})

	_, err := service.HandleExtractFrames(ExtractFramesInput{})
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "src-file は必須です") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleExtractFrames_SourceNotFound(t *testing.T) {
	service := NewServiceWithExecutor(&commandExecutor.MockRepository{})

	_, err := service.HandleExtractFrames(ExtractFramesInput{
		SrcFile: "not-found.mp4",
		FPS:     1,
		Quality: 2,
		OutDir:  t.TempDir(),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "入力動画ファイルが存在しません") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleExtractFrames_FFmpegNotFound(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "input.mp4")
	if err := os.WriteFile(srcFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	mock := &commandExecutor.MockRepository{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return nil, exec.ErrNotFound
		},
	}

	service := NewServiceWithExecutor(mock)
	_, err := service.HandleExtractFrames(ExtractFramesInput{
		SrcFile: srcFile,
		FPS:     1,
		Quality: 2,
		OutDir:  filepath.Join(tempDir, "frames"),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "ffmpeg コマンドが見つかりません") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleExtractFrames_FFmpegExitError(t *testing.T) {
	tempDir := t.TempDir()
	srcFile := filepath.Join(tempDir, "input.mp4")
	if err := os.WriteFile(srcFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("failed to create source file: %v", err)
	}

	mock := &commandExecutor.MockRepository{
		ExecuteFunc: func(name string, args ...string) ([]byte, error) {
			return []byte("ffmpeg failed"), makeExitError(t)
		},
	}

	service := NewServiceWithExecutor(mock)
	_, err := service.HandleExtractFrames(ExtractFramesInput{
		SrcFile: srcFile,
		FPS:     1,
		Quality: 2,
		OutDir:  filepath.Join(tempDir, "frames"),
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "ffmpeg の実行に失敗しました") {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(err.Error(), "ffmpeg failed") {
		t.Fatalf("expected ffmpeg output in error: %v", err)
	}
}

func TestBuildFFmpegError_OtherError(t *testing.T) {
	err := buildFFmpegError(errors.New("boom"), nil)
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func makeExitError(t *testing.T) error {
	t.Helper()
	cmd := exec.Command("sh", "-c", "exit 1")
	err := cmd.Run()
	if err == nil {
		t.Fatal("expected exit error")
	}
	return err
}
