package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_logging/config"
	"github.com/landmaster135/devbox/internal/gcloud_genset_logging/usecases"
)

func main() {
	config, err := cfg.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if config.Help {
		cfg.PrintUsage()
		return
	}

	service := usecases.NewService()

	switch config.Operation {
	case cfg.OperationLoggingRead:
		command, err := service.BuildLoggingReadCommand(usecases.LoggingReadParams{
			Severity:       config.Severity,
			Limit:          config.Limit,
			Query:          config.Query,
			ResourceType:   config.ResourceType,
			Filter:         config.Filter,
			AdditionalArgs: config.AdditionalArgs,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)
	case cfg.OperationCreateSink:
		command, err := service.BuildCreateSinkCommand(usecases.CreateSinkParams{
			SinkName:       config.SinkName,
			Destination:    config.Destination,
			LogFilter:      config.LogFilter,
			AdditionalArgs: config.AdditionalArgs,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作です: %s\n", config.Operation)
		cfg.PrintUsage()
		os.Exit(1)
	}
}
