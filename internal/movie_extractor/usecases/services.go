package usecases

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	commandExecutor "github.com/landmaster135/devbox/internal/movie_extractor/infrastructures/command_executor"
)

const frameFilePattern = "frame_%06d.jpg"

// ExtractFramesInput は extract-frames 操作の入力です。
type ExtractFramesInput struct {
	SrcFile       string
	FPS           int
	Quality       int
	StartPosition string
	OutDir        string
}

// Service は movie extractor のユースケースです。
type Service struct {
	commandExecutor commandExecutor.Repository
}

// NewService は標準依存で Service を生成します。
func NewService() *Service {
	return &Service{
		commandExecutor: commandExecutor.NewOSCommandExecutor(),
	}
}

// NewServiceWithExecutor はテスト用の依存注入を行います。
func NewServiceWithExecutor(executor commandExecutor.Repository) *Service {
	if executor == nil {
		executor = commandExecutor.NewOSCommandExecutor()
	}
	return &Service{commandExecutor: executor}
}

// HandleExtractFrames は動画からフレーム画像を抽出します。
func (s *Service) HandleExtractFrames(input ExtractFramesInput) (string, error) {
	if err := validateInput(input); err != nil {
		return "", err
	}

	sourceInfo, err := os.Stat(input.SrcFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("入力動画ファイルが存在しません: %s", input.SrcFile)
		}
		return "", fmt.Errorf("入力動画ファイルの確認に失敗しました: %w", err)
	}
	if sourceInfo.IsDir() {
		return "", fmt.Errorf("src-file にディレクトリは指定できません: %s", input.SrcFile)
	}

	if err := os.MkdirAll(input.OutDir, 0755); err != nil {
		return "", fmt.Errorf("出力ディレクトリの作成に失敗しました: %w", err)
	}

	absOutDir, err := filepath.Abs(input.OutDir)
	if err != nil {
		return "", fmt.Errorf("出力ディレクトリの絶対パス変換に失敗しました: %w", err)
	}

	outputPattern := filepath.Join(absOutDir, frameFilePattern)
	args := buildFFmpegArgs(input, outputPattern)

	output, err := s.commandExecutor.Execute("ffmpeg", args...)
	if err != nil {
		return "", buildFFmpegError(err, output)
	}

	var result strings.Builder
	result.WriteString("フレーム抽出が完了しました。\n")
	result.WriteString(fmt.Sprintf("出力ディレクトリ: %s\n", absOutDir))
	result.WriteString(fmt.Sprintf("出力ファイルパターン: %s\n", frameFilePattern))
	if len(output) > 0 {
		result.WriteString("\nffmpeg出力:\n")
		result.Write(output)
		if output[len(output)-1] != '\n' {
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}

func validateInput(input ExtractFramesInput) error {
	if strings.TrimSpace(input.SrcFile) == "" {
		return fmt.Errorf("src-file は必須です")
	}
	if input.FPS <= 0 {
		return fmt.Errorf("fps は1以上の整数を指定してください")
	}
	if input.Quality < 1 || input.Quality > 31 {
		return fmt.Errorf("quality は1から31の範囲で指定してください")
	}
	if strings.TrimSpace(input.OutDir) == "" {
		return fmt.Errorf("out-dir は必須です")
	}
	return nil
}

func buildFFmpegArgs(input ExtractFramesInput, outputPattern string) []string {
	args := []string{"-hide_banner", "-loglevel", "error"}
	if input.StartPosition != "" {
		args = append(args, "-ss", input.StartPosition)
	}

	args = append(args,
		"-i", input.SrcFile,
		"-vf", fmt.Sprintf("fps=%d", input.FPS),
		"-q:v", strconv.Itoa(input.Quality),
		"-y",
		outputPattern,
	)
	return args
}

func buildFFmpegError(err error, output []byte) error {
	if errors.Is(err, exec.ErrNotFound) {
		return fmt.Errorf("ffmpeg コマンドが見つかりません。ffmpeg をインストールして PATH を確認してください")
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("ffmpeg の実行に失敗しました(exit=%d): %s", exitError.ExitCode(), strings.TrimSpace(string(output)))
	}

	return fmt.Errorf("ffmpeg の実行に失敗しました: %w", err)
}
