package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_compute/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_compute/usecases"
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
	case cfg.OperationListDiskTypes:
		result, err := service.ExecuteListDiskTypes(usecases.ListDiskTypesParams{
			Zones:          config.Zones,
			MinDiskSizeGiB: config.MinDiskSizeGiB,
			MaxDiskSizeGiB: config.MaxDiskSizeGiB,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	case cfg.OperationListMachineTypes:
		result, err := service.ExecuteListMachineTypes(usecases.ListMachineTypesParams{
			Zones:          config.Zones,
			MinDiskSizeGiB: config.MinDiskSizeGiB,
			MaxDiskSizeGiB: config.MaxDiskSizeGiB,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(result)
	default:
		command, err := service.BuildCommand(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
			os.Exit(1)
		}

		service.PrintHighlightedCommand(command)
	}
}
