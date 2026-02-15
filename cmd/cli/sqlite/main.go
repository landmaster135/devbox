package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/sqlite/config"
	"github.com/landmaster135/devbox/internal/sqlite/usecases"
)

type sqliteService interface {
	HandleListTables(ctx context.Context, dbPath, format string) (string, error)
}

type serviceFactory func() sqliteService

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, defaultServiceFactory))
}

func run(args []string, stdout, stderr io.Writer, factory serviceFactory) int {
	cfg, err := config.ParseFlagsFromArgs(args)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		config.PrintUsageTo(stderr)
		return 1
	}

	if cfg.Help {
		config.PrintUsageTo(stderr)
		return 0
	}

	service := factory()

	switch cfg.Operation {
	case config.OperationListTables:
		result, handleErr := service.HandleListTables(context.Background(), cfg.DBPath, cfg.Format)
		if handleErr != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", handleErr)
			return 1
		}
		if result != "" {
			if _, writeErr := io.WriteString(stdout, result); writeErr != nil {
				fmt.Fprintf(stderr, "エラー: 出力に失敗しました: %v\n", writeErr)
				return 1
			}
			if !strings.HasSuffix(result, "\n") {
				if _, writeErr := io.WriteString(stdout, "\n"); writeErr != nil {
					fmt.Fprintf(stderr, "エラー: 出力に失敗しました: %v\n", writeErr)
					return 1
				}
			}
		}
		return 0
	default:
		fmt.Fprintf(stderr, "エラー: 未対応の operation です: %s\n", cfg.Operation)
		config.PrintUsageTo(stderr)
		return 1
	}
}

func defaultServiceFactory() sqliteService {
	return usecases.NewSQLiteService(nil)
}
