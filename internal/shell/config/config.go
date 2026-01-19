package config

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/landmaster135/devbox/internal/shell/domain"
)

const (
	// OperationExecute はコマンド実行モード
	OperationExecute = "execute"
	// OperationListDenied は禁止コマンド一覧取得モード
	OperationListDenied = "list_denied"
)

// Config はshell CLIの入力設定を保持する
type Config struct {
	Operation          string
	Command            []string
	WorkDir            string
	BaseDir            string
	TimeoutMs          uint64
	Env                map[string]string
	SandboxPermissions domain.SandboxPermissions
	Justification      string
	Help               bool
}

// ParseFlags はos.Argsから設定を解析する
func ParseFlags() (*Config, error) {
	return ParseFlagsWithArgs(os.Args[1:])
}

// ParseFlagsWithArgs は任意の引数スライスを解析する（テスト用）
func ParseFlagsWithArgs(args []string) (*Config, error) {
	cfg := &Config{
		Operation: OperationExecute,
		BaseDir:   ".",
		Env:       map[string]string{},
	}

	flagSet := flag.NewFlagSet("shell", flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)

	var (
		commandParts commandSliceFlag
		sandboxValue string
	)

	flagSet.StringVar(&cfg.Operation, "operation", OperationExecute, "実行する操作 (execute, list_denied)")
	flagSet.StringVar(&cfg.WorkDir, "workdir", "", "作業ディレクトリ (ベースディレクトリからの相対パス可)")
	flagSet.StringVar(&cfg.WorkDir, "cwd", "", "作業ディレクトリのエイリアス")
	flagSet.StringVar(&cfg.BaseDir, "base-dir", cfg.BaseDir, "許可されたベースディレクトリ (デフォルト: カレントディレクトリ)")
	flagSet.StringVar(&cfg.BaseDir, "basedir", cfg.BaseDir, "ベースディレクトリのエイリアス")
	flagSet.Uint64Var(&cfg.TimeoutMs, "timeout-ms", 0, "タイムアウト (ミリ秒)。0はデフォルト値を使用")
	flagSet.Uint64Var(&cfg.TimeoutMs, "timeout", 0, "タイムアウト (ミリ秒) のエイリアス")
	flagSet.StringVar(&sandboxValue, "sandbox-permissions", domain.SandboxPermissionsUseDefault.String(), "サンドボックス権限 (use_default, require_escalated)")
	flagSet.StringVar(&sandboxValue, "sandbox", domain.SandboxPermissionsUseDefault.String(), "サンドボックス権限のエイリアス")
	flagSet.StringVar(&cfg.Justification, "justification", "", "require_escalated指定時の理由")
	flagSet.StringVar(&cfg.Justification, "reason", cfg.Justification, "justificationのエイリアス")

	flagSet.BoolVar(&cfg.Help, "help", false, "ヘルプを表示する")
	flagSet.BoolVar(&cfg.Help, "h", false, "ヘルプを表示する (短縮形)")

	flagSet.Var(&commandParts, "command", "コマンド要素 (複数指定でVec<String>を構成)")
	flagSet.Var(newEnvMapFlag(cfg.Env), "env", "KEY=VALUE形式の環境変数。複数指定可")

	if err := flagSet.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			cfg.Help = true
		} else {
			return nil, err
		}
	}

	if cfg.Help {
		return cfg, nil
	}

	sandbox, err := domain.ParseSandboxPermissions(sandboxValue)
	if err != nil {
		return nil, err
	}
	cfg.SandboxPermissions = sandbox

	cfg.Command = append(cfg.Command, commandParts...)
	cfg.Command = append(cfg.Command, flagSet.Args()...)

	cfg.Operation = strings.ToLower(strings.TrimSpace(cfg.Operation))
	if err := validateOperation(cfg.Operation); err != nil {
		return nil, err
	}

	if cfg.Operation == OperationExecute && len(cfg.Command) == 0 {
		return nil, fmt.Errorf("execute操作にはコマンドが必要です (--commandを複数指定するか \"--\" 以降に指定してください)")
	}

	if cfg.SandboxPermissions.RequiresJustification() && strings.TrimSpace(cfg.Justification) == "" {
		return nil, fmt.Errorf("require_escalatedを使用する場合は--justificationで理由を指定してください")
	}

	if cfg.BaseDir == "" {
		cfg.BaseDir = "."
	}

	if len(cfg.Env) == 0 {
		cfg.Env = nil
	}

	return cfg, nil
}

// PrintUsage はCLIの利用方法を標準エラーに出力する
func PrintUsage() {
	program := os.Args[0]
	fmt.Fprintf(os.Stderr, `shell CLI - Codex互換のシェル実行器

使用方法:
  %[1]s -operation=execute [フラグ] -- <command> [args...]
  %[1]s -operation=list_denied

主なフラグ:
  -operation              実行する操作 (execute, list_denied)
  -command                コマンド要素。複数回指定してVec<String>を構成
  -workdir,-cwd           作業ディレクトリ (ベースディレクトリからの相対パス可)
  -base-dir               許可されるベースディレクトリのルート
  -timeout-ms,-timeout    ミリ秒単位のタイムアウト (0は60秒のデフォルト)
  -env KEY=VALUE          環境変数を追加 (複数指定可)
  -sandbox-permissions    use_default または require_escalated
  -justification          require_escalated指定時の理由
  -help,-h                このヘルプを表示

例:
  %[1]s -operation=execute -workdir=project -- bash -lc "npm test"
  %[1]s -operation=execute -command bash -command -lc -command "ls -a"
  %[1]s -operation=list_denied
`, program)
}

func validateOperation(op string) error {
	switch op {
	case OperationExecute, OperationListDenied:
		return nil
	default:
		return fmt.Errorf("無効なoperationです: %s (有効値: %s)", op, strings.Join(allowedOperations(), ", "))
	}
}

func allowedOperations() []string {
	ops := []string{OperationExecute, OperationListDenied}
	sort.Strings(ops)
	return ops
}

// commandSliceFlag は反復可能な--commandフラグ用
type commandSliceFlag []string

func (c *commandSliceFlag) String() string {
	return strings.Join(*c, " ")
}

func (c *commandSliceFlag) Set(value string) error {
	*c = append(*c, value)
	return nil
}

// envMapFlag はKEY=VALUE形式をmapに格納するflag.Value
type envMapFlag struct {
	target map[string]string
}

func newEnvMapFlag(target map[string]string) *envMapFlag {
	return &envMapFlag{target: target}
}

func (e *envMapFlag) String() string {
	if len(e.target) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(e.target))
	for k, v := range e.target {
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, ",")
}

func (e *envMapFlag) Set(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("--envにはKEY=VALUE形式を指定してください")
	}
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("環境変数の形式が不正です: %s", value)
	}
	key := strings.TrimSpace(parts[0])
	if key == "" {
		return fmt.Errorf("環境変数のキーが空です: %s", value)
	}
	if e.target == nil {
		e.target = map[string]string{}
	}
	e.target[key] = parts[1]
	return nil
}
