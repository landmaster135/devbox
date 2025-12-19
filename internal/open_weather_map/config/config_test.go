package config

import (
	"bytes"
	"io"
	"os"
	"strings"
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

// TestParseFlags_FlagInputs は標準パーサー経由の解析をテストする
func TestParseFlags_FlagInputs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{
		"cmd",
		"-api-key", "abc123",
		"-city", "Osaka,JP",
		"-max-days", "4",
	}

	config, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v, want nil", err)
	}

	if config.APIKey != "abc123" || config.City != "Osaka,JP" || config.MaxDays != 4 {
		t.Errorf("ParseFlags() returned unexpected config: %+v", config)
	}
}

// TestParseFlags_Help はヘルプフラグ指定時の挙動をテストする
func TestParseFlags_Help(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"cmd", "-help"}

	config, err := ParseFlags()
	if err != nil {
		t.Fatalf("ParseFlags() error = %v, want nil", err)
	}

	if !config.Help {
		t.Error("Help flag should set Config.Help to true")
	}
}

// TestParseFlags_MissingRequiredArgs は必須値不足時のエラーをテストする
func TestParseFlags_MissingRequiredArgs(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = []string{"cmd"}

	_, err := ParseFlags()
	if err == nil {
		t.Fatal("ParseFlags() error = nil, want error")
	}

	if !strings.Contains(err.Error(), "APIキーが指定されていません") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestPrintUsage はPrintUsageの出力をテストする
func TestPrintUsage(t *testing.T) {
	origStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	PrintUsage()
	w.Close()
	os.Stderr = origStderr

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read stderr: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "OpenWeather API CLIツール") {
		t.Errorf("PrintUsage() output missing header: %s", output)
	}
	if !strings.Contains(output, "-api-key") {
		t.Errorf("PrintUsage() output missing flag description: %s", output)
	}
}
