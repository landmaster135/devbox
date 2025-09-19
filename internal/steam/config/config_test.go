package config

import (
	"testing"
)

func TestNewConfig_Normal(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		steamAPIKey string
		steamID     string
		gameID      int
		wantErr     bool
	}{
		{
			name:        "valid games operation",
			operation:   "games",
			steamAPIKey: "test-api-key",
			steamID:     "76561198000000000",
			gameID:      0,
			wantErr:     false,
		},
		{
			name:        "valid game-stats operation",
			operation:   "game-stats",
			steamAPIKey: "test-api-key",
			steamID:     "76561198000000000",
			gameID:      123,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := NewConfig(tt.operation, tt.steamAPIKey, tt.steamID, tt.gameID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if config.Operation != tt.operation {
					t.Errorf("NewConfig() Operation = %v, want %v", config.Operation, tt.operation)
				}
				if config.SteamAPIKey != tt.steamAPIKey {
					t.Errorf("NewConfig() SteamAPIKey = %v, want %v", config.SteamAPIKey, tt.steamAPIKey)
				}
				if config.SteamID != tt.steamID {
					t.Errorf("NewConfig() SteamID = %v, want %v", config.SteamID, tt.steamID)
				}
				if config.GameID != tt.gameID {
					t.Errorf("NewConfig() GameID = %v, want %v", config.GameID, tt.gameID)
				}
			}
		})
	}
}

func TestValidateConfig_Normal(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		steamAPIKey string
		steamID     string
		gameID      int
		wantErr     bool
		errContains string
	}{
		{
			name:        "empty operation",
			operation:   "",
			steamAPIKey: "test-api-key",
			steamID:     "76561198000000000",
			gameID:      0,
			wantErr:     true,
			errContains: "operation is required",
		},
		{
			name:        "empty steam api key",
			operation:   "games",
			steamAPIKey: "",
			steamID:     "76561198000000000",
			gameID:      0,
			wantErr:     true,
			errContains: "steam-api-key is required",
		},
		{
			name:        "empty steam id",
			operation:   "games",
			steamAPIKey: "test-api-key",
			steamID:     "",
			gameID:      0,
			wantErr:     true,
			errContains: "steam-id is required",
		},
		{
			name:        "unsupported operation",
			operation:   "invalid-operation",
			steamAPIKey: "test-api-key",
			steamID:     "76561198000000000",
			gameID:      0,
			wantErr:     true,
			errContains: "unsupported operation",
		},
		{
			name:        "invalid steam id format",
			operation:   "games",
			steamAPIKey: "test-api-key",
			steamID:     "123",
			gameID:      0,
			wantErr:     true,
			errContains: "invalid Steam ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.operation, tt.steamAPIKey, tt.steamID, tt.gameID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errContains != "" && err.Error() != "" {
					// エラーメッセージに期待する文字列が含まれているかチェック
					found := false
					if len(err.Error()) >= len(tt.errContains) {
						for i := 0; i <= len(err.Error())-len(tt.errContains); i++ {
							if err.Error()[i:i+len(tt.errContains)] == tt.errContains {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("validateConfig() error = %v, want error containing %v", err, tt.errContains)
					}
				}
			}
		})
	}
}
