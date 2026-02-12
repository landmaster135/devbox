package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

type serviceFactory func(conf *cfg.Config) usecases.MemoService

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
	case cfg.OperationCreateMemo:
		content, contentErr := resolveContent(conf.Content, conf.ContentFile)
		if contentErr != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", contentErr)
			return 1
		}
		result, err = service.CreateMemo(
			ctx,
			conf.MemoID,
			content,
			conf.Visibility,
			conf.State,
			boolPointer(conf.Pinned, conf.PinnedSet),
			conf.DisplayTime,
		)
	case cfg.OperationGetMemo:
		result, err = service.GetMemo(ctx, conf.Memo)
	case cfg.OperationListMemos:
		result, err = service.ListMemos(
			ctx,
			conf.PageSize,
			conf.PageToken,
			conf.State,
			conf.OrderBy,
		)
	case cfg.OperationUpdateMemo:
		content, contentErr := resolveContent(conf.Content, conf.ContentFile)
		if contentErr != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", contentErr)
			return 1
		}
		result, err = service.UpdateMemo(
			ctx,
			conf.Memo,
			content,
			conf.Visibility,
			conf.State,
			boolPointer(conf.Pinned, conf.PinnedSet),
			splitByComma(conf.UpdateMask),
		)
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

func newServiceFromConfig(conf *cfg.Config) usecases.MemoService {
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

func boolPointer(value bool, set bool) *bool {
	if !set {
		return nil
	}
	v := value
	return &v
}

func splitByComma(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		field := strings.TrimSpace(part)
		if field == "" {
			continue
		}
		out = append(out, field)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func resolveContent(content, contentFile string) (string, error) {
	if content != "" && contentFile != "" {
		return "", fmt.Errorf("content と content-file は同時に指定できません")
	}
	if content != "" {
		return content, nil
	}
	if contentFile == "" {
		return "", fmt.Errorf("content または content-file の指定が必要です")
	}

	data, err := os.ReadFile(contentFile)
	if err != nil {
		return "", fmt.Errorf("content-file の読み込みに失敗しました: %w", err)
	}
	fileContent := string(data)
	if strings.TrimSpace(fileContent) == "" {
		return "", fmt.Errorf("content-file が空です: %s", contentFile)
	}
	return fileContent, nil
}
