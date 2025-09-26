package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	config "github.com/landmaster135/devbox/internal/gcloud_genset_deployment/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_deployment/usecases"
)

func main() {
	cfg, err := config.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n\n", err)
		fmt.Fprint(os.Stderr, config.Usage())
		os.Exit(1)
	}

	req, err := buildRequestFromConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	svc := usecases.NewService()
	command, err := svc.BuildCommand(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	switch req.Operation {
	case usecases.OperationListDeployments:
		if cfg.ListDeployments.ShowCommand {
			fmt.Printf("[INFO] list-deployments: 実行コマンド: %s\n", command.String())
		}

		if err := runCommand(command); err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("[INFO] list-deployments: デプロイメント一覧の取得が完了しました")
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の operation です: %s\n", req.Operation)
		os.Exit(1)
	}
}

func buildRequestFromConfig(cfg *config.Config) (usecases.BuildRequest, error) {
	req := usecases.BuildRequest{Operation: usecases.Operation(cfg.Operation)}

	switch req.Operation {
	case usecases.OperationListDeployments:
		req.ListDeployments = &usecases.ListDeploymentsOptions{
			Project: cfg.ListDeployments.Project,
			Filter:  cfg.ListDeployments.Filter,
			Format:  cfg.ListDeployments.Format,
			Limit:   cfg.ListDeployments.Limit,
		}
	default:
		return usecases.BuildRequest{}, fmt.Errorf("未対応の operation です: %s", cfg.Operation)
	}

	return req, nil
}

func runCommand(cmd usecases.Command) error {
	if cmd.Binary == "" {
		return errors.New("実行可能ファイルが指定されていません")
	}

	if _, err := exec.LookPath(cmd.Binary); err != nil {
		return fmt.Errorf("%s コマンドが見つかりません。Google Cloud SDK のインストールを確認してください", cmd.Binary)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	execCmd := exec.CommandContext(ctx, cmd.Binary, cmd.Args...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	if err := execCmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.Canceled) {
			return errors.New("コマンドの実行が中断されました")
		}
		return fmt.Errorf("コマンドの実行に失敗しました: %w", err)
	}

	return nil
}
