package config

import (
	"testing"
)

// TestNewConfig_Normal は正常なConfigの作成をテストする
func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name     string
		apiKey   string
		city     string
		maxDays  int
		expected *Config
	}{
		{
			name:    "正常なConfig作成",
			apiKey:  "test-api-key",
			city:    "Tokyo,JP",
			maxDays: 3,
			expected: &Config{
				APIKey:  "test-api-key",
				City:    "Tokyo,JP",
				MaxDays: 3,
			},
		},
		{
			name:    "最小日数での作成",
			apiKey:  "test-key",
			city:    "London,UK",
			maxDays: 1,
			expected: &Config{
				APIKey:  "test-key",
				City:    "London,UK",
				MaxDays: 1,
			},
		},
		{
			name:    "最大日数での作成",
			apiKey:  "test-key",
			city:    "New York,US",
			maxDays: 5,
			expected: &Config{
				APIKey:  "test-key",
				City:    "New York,US",
				MaxDays: 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.apiKey, tt.city, tt.maxDays)
			if err != nil {
				t.Errorf("NewConfig() error = %v, want nil", err)
				return
			}

			if config.APIKey != tt.expected.APIKey {
				t.Errorf("APIKey = %v, want %v", config.APIKey, tt.expected.APIKey)
			}
			if config.City != tt.expected.City {
				t.Errorf("City = %v, want %v", config.City, tt.expected.City)
			}
			if config.MaxDays != tt.expected.MaxDays {
				t.Errorf("MaxDays = %v, want %v", config.MaxDays, tt.expected.MaxDays)
			}
		})
	}
}

// TestNewConfig_Error は異常なConfigの作成をテストする
func TestNewConfig_Error(t *testing.T) {
	tests := []struct {
		name    string
		apiKey  string
		city    string
		maxDays int
		wantErr bool
	}{
		{
			name:    "APIキーが空",
			apiKey:  "",
			city:    "Tokyo,JP",
			maxDays: 3,
			wantErr: true,
		},
		{
			name:    "都市名が空",
			apiKey:  "test-key",
			city:    "",
			maxDays: 3,
			wantErr: true,
		},
		{
			name:    "最大日数が0",
			apiKey:  "test-key",
			city:    "Tokyo,JP",
			maxDays: 0,
			wantErr: true,
		},
		{
			name:    "最大日数が6",
			apiKey:  "test-key",
			city:    "Tokyo,JP",
			maxDays: 6,
			wantErr: true,
		},
		{
			name:    "最大日数が負の値",
			apiKey:  "test-key",
			city:    "Tokyo,JP",
			maxDays: -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewConfig(tt.apiKey, tt.city, tt.maxDays)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
