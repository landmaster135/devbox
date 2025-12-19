package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_cloudsql/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_cloudsql/usecases"
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
	case cfg.OperationDeleteInstance:
		return service.BuildDeleteInstanceCommand(usecases.DeleteInstanceParams{InstanceName: conf.InstanceName})
	case cfg.OperationPatchDeletionProtection:
		return service.BuildPatchDeletionProtectionCommand(usecases.PatchDeletionProtectionParams{
			InstanceName: conf.InstanceName,
			Mode:         conf.DeletionProtectionMode,
		})
	case cfg.OperationPatchActivationPolicy:
		return service.BuildPatchActivationPolicyCommand(usecases.PatchActivationPolicyParams{
			InstanceName: conf.InstanceName,
			Policy:       conf.ActivationPolicy,
		})
	case cfg.OperationStartInstance:
		return service.BuildStartInstanceCommand(usecases.InstanceParams{InstanceName: conf.InstanceName})
	case cfg.OperationStopInstance:
		return service.BuildStopInstanceCommand(usecases.InstanceParams{InstanceName: conf.InstanceName})
	default:
		return "", fmt.Errorf("未対応のoperationです: %s", conf.Operation)
	}
}
