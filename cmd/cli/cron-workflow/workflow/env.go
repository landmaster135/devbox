package workflow

import (
	"errors"

	infraEnv "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/env"
)

type EnvKey string

const (
	EnvKeyHeartOwner                        EnvKey = "HEART_OWNER"
	EnvKeyDiscordWebhookURLForWeather       EnvKey = "DISCORD_WEBHOOK_URL_FOR_WEATHER"
	EnvKeyOpenWeatherAPIKey                 EnvKey = "OPEN_WEATHER_API_KEY"
	EnvKeyDiscordWebhookURLForDailyTemplate EnvKey = "DISCORD_WEBHOOK_URL_FOR_DAILY_TEMPLATE"
	EnvKeyDBURL01Staging                    EnvKey = "DATABASE_URL_01_STAGING"
	EnvKeyDBDirectory01Staging              EnvKey = "DATABASE_DUMP_DIR_01_STAGING"
	EnvKeyDBURL01Product                    EnvKey = "DATABASE_URL_01_PRODUCT"
	EnvKeyDBDirectory01Product              EnvKey = "DATABASE_DUMP_DIR_01_PRODUCT"
	EnvKeyDBURL01MemosStaging               EnvKey = "DATABASE_URL_01_MEMOS_STAGING"
	EnvKeyDBDirectory01MemosStaging         EnvKey = "DATABASE_DUMP_DIR_01_MEMOS_STAGING"
	EnvKeyPCInfoOutputDirectory             EnvKey = "PC_INFO_OUTPUT_DIR"
	EnvKeyPCInfoHostnameOfNAS01             EnvKey = "PC_INFO_MEMORY_HOSTNAME_OF_NAS_01"
	EnvKeyPCInfoMemoryNamesOfNAS01          EnvKey = "PC_INFO_MEMORY_NAMES_OF_NAS_01"
	EnvKeyPCInfoMemoryManufacturersOfNAS01  EnvKey = "PC_INFO_MEMORY_MANUFACTURERS_OF_NAS_01"
)

func getEnvVars(repo infraEnv.Repository, envKey EnvKey) (string, error) {
	k := string(envKey)
	return repo.GetEnv(k)
}

func getOptionalEnvVars(repo infraEnv.Repository, envKey EnvKey) (string, error) {
	value, err := getEnvVars(repo, envKey)
	if err == nil {
		return value, nil
	}

	var missing infraEnv.MissingEnvError
	if errors.As(err, &missing) {
		return "", nil
	}

	return "", err
}
