package workflow

import (
	infraEnv "github.com/landmaster135/devbox/internal/cron_workflow/infrastructure/env"
)

type EnvKey string

const (
	EnvKeyHeartOwner        EnvKey = "HEART_OWNER"
	EnvKeyDiscordWebhookURL EnvKey = "DISCORD_WEBHOOK_URL"
	EnvKeyOpenWeatherAPIKey EnvKey = "OPEN_WEATHER_API_KEY"
)

func getEnvVars(repo infraEnv.Repository, envKey EnvKey) (string, error) {
	k := string(envKey)
	return repo.GetEnv(k)
}
