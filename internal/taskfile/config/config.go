package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	// OperationInspect は Taskfile を検証する操作
	OperationInspect = "inspect"
	// OperationFill は Taskfile の空欄を参照Taskfileの値で補完する操作
	OperationFill = "fill"
	// OperationNew は参照Taskfileをもとに新規 Taskfile を作成する操作
	OperationNew = "new"
	// TaskTypeRoot は root Taskfile の検証を指すタスクタイプ
	TaskTypeRoot = "root"
)

var (
	allowedOperations = []string{OperationInspect, OperationFill, OperationNew}
	allowedTaskTypes  = []string{TaskTypeRoot}
)

// Config は taskfile CLI の設定を保持する
type Config struct {
	Operation    string
	TaskType     string
	TaskfilePath string
	Help         bool
}

// NewConfig は入力値を検証して Config を生成する
func NewConfig(operation, taskType, taskfilePath string, help bool) (*Config, error) {
	cfg := &Config{
		Operation:    operation,
		TaskType:     taskType,
		TaskfilePath: taskfilePath,
		Help:         help,
	}

	if cfg.Help {
		return cfg, nil
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.Operation == "" {
		return fmt.Errorf("--operation は必須です")
	}

	if !isAllowed(c.Operation, allowedOperations) {
		return fmt.Errorf("--operation には %s のいずれかを指定してください", strings.Join(allowedOperations, ", "))
	}

	if c.TaskType == "" {
		return fmt.Errorf("--task-type は必須です")
	}

	if !isAllowed(c.TaskType, allowedTaskTypes) {
		return fmt.Errorf("--task-type には %s のいずれかを指定してください", strings.Join(allowedTaskTypes, ", "))
	}

	if c.TaskfilePath == "" {
		return fmt.Errorf("--taskfile-path は必須です")
	}

	return nil
}

func isAllowed(value string, candidates []string) bool {
	for _, candidate := range candidates {
		if candidate == value {
			return true
		}
	}
	return false
}

// ParseFlags は CLI フラグを解析して Config を返す
func ParseFlags() (*Config, error) {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flagSet.Usage = func() {
		PrintUsage()
	}

	var (
		operation    string
		taskType     string
		taskfilePath string
		help         bool
	)

	flagSet.StringVar(&operation, "operation", "", fmt.Sprintf("実行する操作 (%s)", strings.Join(allowedOperations, ", ")))
	flagSet.StringVar(&taskType, "task-type", "", fmt.Sprintf("Taskfileの種類 (%s)", strings.Join(allowedTaskTypes, ", ")))
	flagSet.StringVar(&taskfilePath, "taskfile-path", "", "検証/補完対象のTaskfileパス")
	flagSet.BoolVar(&help, "help", false, "ヘルプを表示")
	flagSet.BoolVar(&help, "h", false, "ヘルプを表示 (短縮形)")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	return NewConfig(operation, taskType, taskfilePath, help)
}

// PrintUsage は CLI の使用方法を表示する
func PrintUsage() {
	exeName := os.Args[0]

	fmt.Fprintf(os.Stderr, "Taskfile CLI ツール\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n  %s --operation inspect --task-type root --taskfile-path ./Taskfile.yml\n  %s --operation fill --task-type root --taskfile-path ./Taskfile.yml\n  %s --operation new --task-type root --taskfile-path ./Taskfile.yml\n\n", exeName, exeName, exeName)
	fmt.Fprintf(os.Stderr, "オプション:\n")
	fmt.Fprintf(os.Stderr, "  --operation string\n        実行する操作 (%s)\n", strings.Join(allowedOperations, ", "))
	fmt.Fprintf(os.Stderr, "  --task-type string\n        Taskfileの種類 (%s)\n", strings.Join(allowedTaskTypes, ", "))
	fmt.Fprintf(os.Stderr, "  --taskfile-path string\n        検証/補完/作成対象のTaskfileへのパス\n")
	fmt.Fprintf(os.Stderr, "  --help\n        このヘルプを表示\n")
}
