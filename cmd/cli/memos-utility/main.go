package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	cfg "github.com/landmaster135/devbox/internal/memos_utility/config"
	usecases "github.com/landmaster135/devbox/internal/memos_utility/usecases"
)

type serviceFactory func(conf *cfg.Config) usecases.MemosUtilityService

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, newServiceFromConfig))
}

func run(args []string, stdout, stderr io.Writer, factory serviceFactory) int {
	conf, err := cfg.ParseFlagsFromArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		return 1
	}
	if conf.Help {
		cfg.PrintUsage()
		return 0
	}

	service := factory(conf)
	ctx := context.Background()

	var result any
	switch conf.Operation {
	case cfg.OperationCreateWebClip, cfg.OperationCreateMovieClip:
		result, err = service.CreateClip(ctx, usecases.CreateClipInput{
			Operation:   conf.Operation,
			ContentFile: conf.ContentFile,
			Attachments: splitByComma(conf.Attachments),
		})
	default:
		fmt.Fprintf(stderr, "エラー: 未対応の operation です: %s\n", conf.Operation)
		return 1
	}
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return 1
	}

	if err := printJSON(stdout, result); err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return 1
	}
	return 0
}

func newServiceFromConfig(conf *cfg.Config) usecases.MemosUtilityService {
	return usecases.NewService(usecases.ServiceOptions{
		BaseURL:  conf.BaseURL,
		APIToken: conf.APIToken,
		Timeout:  time.Duration(conf.TimeoutSeconds) * time.Second,
	})
}

func printJSON(writer io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON への整形に失敗しました: %w", err)
	}
	fmt.Fprintln(writer, string(data))
	return nil
}

func splitByComma(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item == "" {
			continue
		}
		out = append(out, item)
	}

	if len(out) == 0 {
		return nil
	}
	return out
}
