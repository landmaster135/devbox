package internal

import "time"

// UsageData represents a usage entry from JSONL files
type UsageData struct {
	Timestamp           time.Time `json:"timestamp"`
	InputTokens         int       `json:"input_tokens"`
	OutputTokens        int       `json:"output_tokens"`
	CacheCreationTokens int       `json:"cache_creation_tokens"`
	CacheReadTokens     int       `json:"cache_read_tokens"`
	Cost                float64   `json:"cost"`
}

// DailyUsage represents aggregated usage for a single day
type DailyUsage struct {
	Date                string  `json:"date"`
	InputTokens         int     `json:"input_tokens"`
	OutputTokens        int     `json:"output_tokens"`
	CacheCreationTokens int     `json:"cache_creation_tokens"`
	CacheReadTokens     int     `json:"cache_read_tokens"`
	TotalTokens         int     `json:"total_tokens"`
	Cost                float64 `json:"cost"`
}

// SessionUsage represents aggregated usage for a conversation session
type SessionUsage struct {
	Project             string    `json:"project"`
	Session             string    `json:"session"`
	InputTokens         int       `json:"input_tokens"`
	OutputTokens        int       `json:"output_tokens"`
	CacheCreationTokens int       `json:"cache_creation_tokens"`
	CacheReadTokens     int       `json:"cache_read_tokens"`
	TotalTokens         int       `json:"total_tokens"`
	Cost                float64   `json:"cost"`
	LastActivity        time.Time `json:"last_activity"`
}

// Config holds application configuration
type Config struct {
	ClaudePath string
	Since      *time.Time
	Until      *time.Time
	JSON       bool
}
