package main

import (
	"fmt"
	"os"

	cfg "github.com/landmaster135/devbox/internal/gcloud_genset_secret/config"
	usecases "github.com/landmaster135/devbox/internal/gcloud_genset_secret/usecases"
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

	var command string

	switch config.Operation {
	case cfg.OperationCreateSecret:
		command, err = service.BuildCreateSecretCommand(usecases.CreateSecretParams{
			SecretName:        config.SecretName,
			ReplicationPolicy: config.ReplicationPolicy,
			Locations:         config.Locations,
		})
	case cfg.OperationAddSecretVersion:
		command, err = service.BuildAddSecretVersionCommand(usecases.AddSecretVersionParams{
			SecretName:  config.SecretName,
			SecretValue: config.SecretValue,
		})
	case cfg.OperationCreateAndAddSecretVersion:
		command, err = service.BuildCreateAndAddSecretVersionCommand(usecases.CreateAndAddSecretVersionParams{
			SecretName:        config.SecretName,
			SecretValue:       config.SecretValue,
			ReplicationPolicy: config.ReplicationPolicy,
			Locations:         config.Locations,
		})
	case cfg.OperationAccessSecretVersion:
		command, err = service.BuildAccessSecretVersionCommand(usecases.AccessSecretVersionParams{
			SecretName: config.SecretName,
			Version:    config.Version,
		})
	case cfg.OperationUpdateSecretLabels:
		command, err = service.BuildUpdateSecretLabelsCommand(usecases.UpdateSecretLabelsParams{
			SecretName: config.SecretName,
			Labels:     config.Labels,
		})
	case cfg.OperationUpdateSecretVersionAliases:
		command, err = service.BuildUpdateSecretVersionAliasesCommand(usecases.UpdateSecretVersionAliasesParams{
			SecretName:  config.SecretName,
			AliasOption: config.AliasOption,
		})
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作です: %s\n", config.Operation)
		cfg.PrintUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	service.PrintHighlightedCommand(command)
	notificationScript, ok := service.BuildNotificationWrappedCommand(
		usecases.DiscordNotificationParams{
			Operation:  config.Operation,
			SecretName: config.SecretName,
		},
		command,
	)
	if ok {
		service.PrintNotificationScript(notificationScript)
	}
}
