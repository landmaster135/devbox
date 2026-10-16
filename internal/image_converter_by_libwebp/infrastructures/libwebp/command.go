package libwebp

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
)

// Converter は libwebp による変換処理を抽象化します。
type Converter interface {
	CheckAvailable() error
	ConvertToWebP(ctx context.Context, inputPath string, outputPath string, quality int, method int, lossless bool) error
}

// CommandConverter は cwebp コマンドを使う Converter 実装です。
type CommandConverter struct {
	commandName string
	lookPath    func(file string) (string, error)
	runCommand  func(ctx context.Context, name string, args ...string) (string, error)
}

// NewCommandConverter は cwebp 実行器を作成します。
func NewCommandConverter() *CommandConverter {
	return &CommandConverter{
		commandName: "cwebp",
		lookPath:    exec.LookPath,
		runCommand:  runCommand,
	}
}

// CheckAvailable は cwebp が PATH に存在することを確認します。
func (c *CommandConverter) CheckAvailable() error {
	if _, err := c.lookPath(c.commandName); err != nil {
		return fmt.Errorf("libwebp パッケージが見つかりません: %s をインストールしてください", c.commandName)
	}
	return nil
}

// ConvertToWebP は cwebp で単一ファイルを WebP へ変換します。
func (c *CommandConverter) ConvertToWebP(ctx context.Context, inputPath string, outputPath string, quality int, method int, lossless bool) error {
	args := c.buildArgs(inputPath, outputPath, quality, method, lossless)
	output, err := c.runCommand(ctx, c.commandName, args...)
	if err != nil {
		return fmt.Errorf("cwebp 実行に失敗しました: %w: %s", err, output)
	}
	return nil
}

func runCommand(ctx context.Context, name string, args ...string) (string, error) {
	output, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	return string(output), err
}

func (c *CommandConverter) buildArgs(inputPath string, outputPath string, quality int, method int, lossless bool) []string {
	args := []string{
		"-preset", "photo",
		"-metadata", "icc",
		"-sharp_yuv",
		"-progress",
		"-short",
	}
	if lossless {
		args = append(args, "-lossless")
	}
	args = append(args, "-m", strconv.Itoa(method))
	args = append(args, "-q", strconv.Itoa(quality), inputPath, "-o", outputPath)
	return args
}
