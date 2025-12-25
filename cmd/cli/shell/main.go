package main

import (
	"encoding/json"
	"fmt"
	"os"

	shellconfig "github.com/landmaster135/devbox/internal/shell/config"
	shellusecases "github.com/landmaster135/devbox/internal/shell/usecases"
)

func main() {
	cfg, err := shellconfig.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		shellconfig.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		shellconfig.PrintUsage()
		return
	}

	service := shellusecases.NewShellService()

	switch cfg.Operation {
	case shellconfig.OperationExecute:
		exitCode := runExecute(service, cfg)
		os.Exit(exitCode)
	case shellconfig.OperationListDenied:
		if err := runListDenied(service); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応のoperationです: %s\n", cfg.Operation)
		shellconfig.PrintUsage()
		os.Exit(1)
	}
}

func runExecute(service *shellusecases.ShellService, cfg *shellconfig.Config) int {
	input := &shellusecases.ExecuteCommandInput{
		Command:            cfg.Command,
		WorkDir:            cfg.WorkDir,
		BaseDir:            cfg.BaseDir,
		TimeoutMs:          cfg.TimeoutMs,
		Env:                cfg.Env,
		SandboxPermissions: cfg.SandboxPermissions,
		Justification:      cfg.Justification,
	}

	result, err := service.ExecuteCommand(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		return 1
	}

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: 結果のシリアライズに失敗しました: %v\n", err)
		return 1
	}

	fmt.Println(string(payload))

	if result.Success {
		return 0
	}

	if result.TimedOut {
		fmt.Fprintf(os.Stderr, "コマンドがタイムアウトしました (上限 %dms)\n", result.TimeoutMs)
	} else {
		fmt.Fprintf(os.Stderr, "コマンドが終了コード%dで終了しました\n", result.ExitCode)
	}

	exitCode := result.ExitCode
	if exitCode <= 0 {
		exitCode = 1
	}
	return exitCode
}

func runListDenied(service *shellusecases.ShellService) error {
	commands := service.ListDeniedCommands()
	payload := map[string][]string{
		"commands": commands,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
