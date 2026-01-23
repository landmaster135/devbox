package workflow

import (
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
)

func getEnvVars(repo infraEnv.Repository, envKey EnvKey) (string, error) {
	k := string(envKey)
	return repo.GetEnv(k)
}
