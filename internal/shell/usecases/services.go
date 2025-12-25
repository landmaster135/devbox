package usecases

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/shell/domain"
)

const (
	defaultCommandTimeout = 60 * time.Second
)

// CommandExecutor はexec.CommandContextを抽象化する
type CommandExecutor interface {
	Execute(ctx context.Context, command []string, workDir string, env map[string]string) (*ExecutionOutput, error)
}

// ExecutionOutput は外部コマンドの生出力を保持する
type ExecutionOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
	TimedOut bool
}

// DefaultCommandExecutor は実際にOSコマンドを実行する実装
type DefaultCommandExecutor struct{}

// Execute はOSコマンドを実行し、標準出力・標準エラーをバッファリングする
func (e *DefaultCommandExecutor) Execute(ctx context.Context, command []string, workDir string, env map[string]string) (*ExecutionOutput, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("コマンドが指定されていません")
	}

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	cmd.Env = mergeEnvironments(env)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := &ExecutionOutput{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			output.ExitCode = -1
			output.TimedOut = true
			return output, nil
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			output.ExitCode = exitErr.ExitCode()
			return output, nil
		}

		return nil, fmt.Errorf("コマンド実行に失敗しました: %w", err)
	}

	return output, nil
}

// ShellService はshellツールのユースケースをまとめる
type ShellService struct {
	executor   CommandExecutor
	maxTimeout time.Duration
}

// NewShellService はデフォルト依存を組み立てる
func NewShellService() *ShellService {
	return &ShellService{
		executor:   &DefaultCommandExecutor{},
		maxTimeout: defaultCommandTimeout,
	}
}

// NewShellServiceWithExecutor はテスト用の依存注入を行う
func NewShellServiceWithExecutor(executor CommandExecutor) *ShellService {
	return &ShellService{
		executor:   executor,
		maxTimeout: defaultCommandTimeout,
	}
}

// ExecuteCommandInput はExecuteCommandへの入力値
type ExecuteCommandInput struct {
	Command            []string
	WorkDir            string
	BaseDir            string
	TimeoutMs          uint64
	Env                map[string]string
	SandboxPermissions domain.SandboxPermissions
	Justification      string
}

// CommandResult はCLIへ返す構造体
type CommandResult struct {
	Command             []string                  `json:"command"`
	BaseDir             string                    `json:"base_dir"`
	WorkDir             string                    `json:"workdir"`
	Stdout              string                    `json:"stdout"`
	Stderr              string                    `json:"stderr"`
	ExitCode            int                       `json:"exit_code"`
	Success             bool                      `json:"success"`
	TimedOut            bool                      `json:"timed_out"`
	DurationMs          int64                     `json:"duration_ms"`
	TimeoutMs           int64                     `json:"timeout_ms"`
	SandboxPermissions  domain.SandboxPermissions `json:"sandbox_permissions"`
	EscalationRequested bool                      `json:"escalation_requested"`
	Justification       string                    `json:"justification,omitempty"`
}

// ExecuteCommand は構成済みコマンドを実行する
func (s *ShellService) ExecuteCommand(input *ExecuteCommandInput) (*CommandResult, error) {
	if input == nil {
		return nil, fmt.Errorf("入力が指定されていません")
	}
	if len(input.Command) == 0 {
		return nil, fmt.Errorf("コマンドを1つ以上指定してください")
	}

	sandbox := input.SandboxPermissions
	if sandbox == "" {
		sandbox = domain.SandboxPermissionsUseDefault
	}

	if sandbox.RequiresJustification() && strings.TrimSpace(input.Justification) == "" {
		return nil, fmt.Errorf("require_escalatedを利用する場合はjustificationが必須です")
	}

	baseDir := input.BaseDir
	if baseDir == "" {
		baseDir = "."
	}

	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return nil, fmt.Errorf("ベースディレクトリの解決に失敗しました: %w", err)
	}

	if err := ensureDir(absBase); err != nil {
		return nil, fmt.Errorf("ベースディレクトリが不正です: %w", err)
	}

	workDir, err := resolveWorkDir(absBase, input.WorkDir)
	if err != nil {
		return nil, err
	}

	timeout := s.maxTimeout
	if input.TimeoutMs > 0 {
		requested := time.Duration(input.TimeoutMs) * time.Millisecond
		if requested < s.maxTimeout {
			timeout = requested
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	envCopy := cloneEnv(input.Env)
	output, err := s.executor.Execute(ctx, input.Command, workDir, envCopy)
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}

	success := output.ExitCode == 0 && !output.TimedOut

	result := &CommandResult{
		Command:             append([]string(nil), input.Command...),
		BaseDir:             absBase,
		WorkDir:             workDir,
		Stdout:              output.Stdout,
		Stderr:              output.Stderr,
		ExitCode:            output.ExitCode,
		Success:             success,
		TimedOut:            output.TimedOut,
		DurationMs:          duration.Milliseconds(),
		TimeoutMs:           int64(timeout / time.Millisecond),
		SandboxPermissions:  sandbox,
		EscalationRequested: sandbox.RequiresJustification(),
		Justification:       strings.TrimSpace(input.Justification),
	}

	if !result.Success && result.ExitCode == 0 {
		result.ExitCode = 1
	}

	if result.DurationMs < 0 {
		result.DurationMs = 0
	}

	return result, nil
}

// ListAllowedCommands は許可済みコマンド一覧を返す
func (s *ShellService) ListAllowedCommands() []string {
	commands := make([]string, len(defaultAllowedCommands))
	copy(commands, defaultAllowedCommands)
	sort.Strings(commands)
	return commands
}

// defaultAllowedCommands はMCP実装と共有される許可コマンド
var defaultAllowedCommands = []string{
	"npm",
	"yarn",
	"pnpm",
	"bun",
	"git",
	"ls",
	"dir",
	"find",
	"mkdir",
	"rmdir",
	"cp",
	"mv",
	"rm",
	"cat",
	"awk",
	"wc",
	"node",
	"python",
	"python3",
	"tsc",
	"eslint",
	"prettier",
	"make",
	"cargo",
	"go",
	"docker",
	"docker-compose",
	"echo",
	"touch",
	"grep",
	"bash",
	"sh",
	"powershell.exe",
	"pwsh",
}

func resolveWorkDir(baseDir, subDir string) (string, error) {
	if strings.TrimSpace(subDir) == "" {
		return baseDir, nil
	}

	candidate := subDir
	if !filepath.IsAbs(candidate) {
		candidate = filepath.Join(baseDir, candidate)
	}

	absCandidate, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("作業ディレクトリの解決に失敗しました: %w", err)
	}

	rel, err := filepath.Rel(baseDir, absCandidate)
	if err != nil {
		return "", fmt.Errorf("作業ディレクトリの相対パス取得に失敗しました: %w", err)
	}
	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("作業ディレクトリが許可されたベースディレクトリの外です: %s", absCandidate)
	}

	if err := ensureDir(absCandidate); err != nil {
		return "", fmt.Errorf("作業ディレクトリが不正です: %w", err)
	}

	return absCandidate, nil
}

func ensureDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("ディレクトリが存在しません: %s", path)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("ディレクトリではありません: %s", path)
	}
	return nil
}

func mergeEnvironments(custom map[string]string) []string {
	base := os.Environ()
	if len(custom) == 0 {
		return base
	}

	keys := make([]string, 0, len(custom))
	for k := range custom {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		base = append(base, fmt.Sprintf("%s=%s", key, custom[key]))
	}
	return base
}

func cloneEnv(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
