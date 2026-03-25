package dump_binary

import (
	"context"
	"strings"
	"time"
)

const (
	defaultPgDumpRetryMaxAttempts = 3
	defaultPgDumpRetryBaseDelay   = 2 * time.Second
	defaultPgDumpRetryMaxDelay    = 8 * time.Second
)

type SleepWithContextFunc func(ctx context.Context, d time.Duration) error

// RetryConfig は pg_dump 実行時のリトライ設定です。
type RetryConfig struct {
	MaxAttempts      int
	BaseDelay        time.Duration
	MaxDelay         time.Duration
	SleepWithContext SleepWithContextFunc
}

func isRetriablePgDumpError(err error, commandOutput string) bool {
	if err == nil {
		return false
	}

	body := strings.ToLower(strings.TrimSpace(err.Error() + " " + commandOutput))
	if body == "" {
		return false
	}

	retriablePatterns := []string{
		"control plane request failed",
		"connection to server",
		"could not connect to server",
		"server closed the connection unexpectedly",
		"connection reset by peer",
		"connection refused",
		"timeout",
		"timed out",
		"temporarily unavailable",
		"the database system is starting up",
	}

	for _, p := range retriablePatterns {
		if strings.Contains(body, p) {
			return true
		}
	}

	return false
}

func pgDumpRetryDelay(retryConfig RetryConfig, attempt int) time.Duration {
	if attempt <= 0 {
		return retryConfig.BaseDelay
	}

	delay := retryConfig.BaseDelay
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= retryConfig.MaxDelay {
			return retryConfig.MaxDelay
		}
	}

	if delay > retryConfig.MaxDelay {
		return retryConfig.MaxDelay
	}

	return delay
}

func normalizeRetryConfig(config RetryConfig) RetryConfig {
	if config.MaxAttempts <= 0 {
		config.MaxAttempts = defaultPgDumpRetryMaxAttempts
	}
	if config.BaseDelay <= 0 {
		config.BaseDelay = defaultPgDumpRetryBaseDelay
	}
	if config.MaxDelay <= 0 {
		config.MaxDelay = defaultPgDumpRetryMaxDelay
	}
	if config.MaxDelay < config.BaseDelay {
		config.MaxDelay = config.BaseDelay
	}
	if config.SleepWithContext == nil {
		config.SleepWithContext = sleepWithContext
	}

	return config
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
