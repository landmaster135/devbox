package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_dns/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_dns/usecases"
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
	case cfg.OperationManagedZonesList:
		command, err := service.BuildManagedZonesListCommand(usecases.ManagedZonesListParams{
			Project:        config.Project,
			Format:         config.Format,
			Filter:         config.Filter,
			Limit:          config.Limit,
			PageSize:       config.PageSize,
			SortBy:         config.SortBy,
			Verbosity:      config.Verbosity,
			URI:            config.URI,
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
