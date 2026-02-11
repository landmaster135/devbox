package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_spanner/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_spanner/usecases"
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

	service := usecases.NewService()
	command, err := buildCommand(service, conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	service.PrintHighlightedCommand(command)
}

func buildCommand(service *usecases.Service, conf *cfg.Config) (string, error) {
	switch conf.Operation {
	case cfg.OperationInstanceList:
		return service.BuildInstanceListCommand()
	case cfg.OperationInstanceCreate:
		return service.BuildInstanceCreateCommand(usecases.InstanceCreateParams{
			InstanceID:     conf.InstanceID,
			InstanceConfig: conf.InstanceConfig,
			Description:    conf.Description,
			Nodes:          conf.Nodes,
		})
	case cfg.OperationDatabaseCreate:
		return service.BuildDatabaseCreateCommand(usecases.DatabaseCreateParams{
			InstanceID:  conf.InstanceID,
			DatabaseID:  conf.DatabaseID,
			DDLFilePath: conf.DDLFilePath,
		})
	case cfg.OperationDatabaseList:
		return service.BuildDatabaseListCommand(usecases.DatabaseListParams{InstanceID: conf.InstanceID})
	case cfg.OperationDatabaseDescribe:
		return service.BuildDatabaseDescribeCommand(usecases.DatabaseDescribeParams{
			InstanceID: conf.InstanceID,
			DatabaseID: conf.DatabaseID,
		})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}
