package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	cfg "github.com/landmaster135/devbox/internal/ollama/config"
	"github.com/landmaster135/devbox/internal/ollama/domain"
	usecases "github.com/landmaster135/devbox/internal/ollama/usecases"
)

func main() {
	conf, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if conf.Help {
		cfg.PrintUsage()
		return
	}

	baseURL := fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)
	service := usecases.NewService(usecases.ServiceOptions{
		BaseURL: baseURL,
		Timeout: time.Duration(conf.TimeoutSeconds) * time.Second,
	})

	ctx := context.Background()

	switch conf.Operation {
	case cfg.OperationVersion:
		handleVersion(ctx, service)
	case cfg.OperationListModels:
		handleListModels(ctx, service, conf.RunningOnly)
	case cfg.OperationEmbed:
		handleEmbed(ctx, service, conf)
	case cfg.OperationGenerate:
		handleGenerate(ctx, service, conf)
	case cfg.OperationPull:
		handlePull(ctx, service, conf)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の operation です: %s\n", conf.Operation)
		os.Exit(1)
	}
}

func handleVersion(ctx context.Context, service *usecases.Service) {
	resp, err := service.GetVersion(ctx)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(resp); err != nil {
		exitWithError(err)
	}
}

func handleListModels(ctx context.Context, service *usecases.Service, runningOnly bool) {
	if runningOnly {
		resp, err := service.ListRunningModels(ctx)
		if err != nil {
			exitWithError(err)
		}
		if err := printJSON(resp); err != nil {
			exitWithError(err)
		}
		return
	}

	resp, err := service.ListInstalledModels(ctx)
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(resp); err != nil {
		exitWithError(err)
	}
}

func handleEmbed(ctx context.Context, service *usecases.Service, conf *cfg.Config) {
	resp, err := service.CreateEmbeddings(ctx, domain.EmbedRequest{
		Model: conf.Model,
		Input: conf.Inputs,
	})
	if err != nil {
		exitWithError(err)
	}
	if err := printJSON(resp); err != nil {
		exitWithError(err)
	}
}

func handleGenerate(ctx context.Context, service *usecases.Service, conf *cfg.Config) {
	output, err := service.Generate(ctx, domain.GenerateRequest{
		Model:  conf.Model,
		Prompt: conf.Prompt,
	})
	if err != nil {
		exitWithError(err)
	}
	printText(output)
}

func handlePull(ctx context.Context, service *usecases.Service, conf *cfg.Config) {
	if err := service.StreamPull(ctx, domain.PullRequest{Model: conf.Model}, os.Stdout); err != nil {
		exitWithError(err)
	}
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON への整形に失敗しました: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

func printText(text string) {
	if text == "" {
		fmt.Println()
		return
	}
	if strings.HasSuffix(text, "\n") {
		fmt.Print(text)
	} else {
		fmt.Println(text)
	}
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
	os.Exit(1)
}
