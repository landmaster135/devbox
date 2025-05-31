package config

import "time"

// AppConfig holds application configuration
type AppConfig struct {
	Command     string     // "daily" or "session"
	ClaudePath  string     // Path to Claude data directory
	Since       *time.Time // Start date filter
	Until       *time.Time // End date filter
	OutputJSON  bool       // Output in JSON format
	ShowHelp    bool       // Show help message
	ShowVersion bool       // Show version information
}

// SetDefaults sets default values for the configuration
func (c *AppConfig) SetDefaults() {
	if c.Command == "" {
		c.Command = "daily"
	}
}
