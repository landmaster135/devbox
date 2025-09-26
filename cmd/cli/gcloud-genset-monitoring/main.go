package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_monitoring/config"
	"github.com/landmaster135/devbox/internal/gcloud_genset_monitoring/usecases"
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
	case cfg.OperationListDashboards:
		command, err := service.BuildListDashboardsCommand(usecases.ListDashboardsParams{
			Project:  config.Project,
			Filter:   config.Filter,
			Format:   config.Format,
			PageSize: config.PageSize,
			SortBy:   config.SortBy,
			Limit:    config.Limit,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)

	case cfg.OperationDescribeDashboard:
		command, err := service.BuildDescribeDashboardCommand(usecases.DescribeDashboardParams{
			DashboardID: config.DashboardID,
			Project:     config.Project,
			Format:      config.Format,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)

	case cfg.OperationListSnoozes:
		command, err := service.BuildListSnoozesCommand(usecases.ListSnoozesParams{
			Project:    config.Project,
			Filter:     config.Filter,
			Format:     config.Format,
			PageSize:   config.PageSize,
			SortBy:     config.SortBy,
			Limit:      config.Limit,
			IncludeURI: config.URI,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		service.PrintHighlightedCommand(command)

	case cfg.OperationListUptimeConfigs:
		command, err := service.BuildListUptimeConfigsCommand(usecases.ListUptimeConfigsParams{
			Project:    config.Project,
			Filter:     config.Filter,
			Format:     config.Format,
			PageSize:   config.PageSize,
			SortBy:     config.SortBy,
			Limit:      config.Limit,
			IncludeURI: config.URI,
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
