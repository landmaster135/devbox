package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          Config Tests                                      ##
// #==============================================================#

func TestParseFlags_Normal(t *testing.T) {
	const (
		testService      = "github"
		testToken        = "test-token"
		testSaveFilePath = "/tmp/test-output.json"
	)

	tests := []struct {
		name                 string
		args                 []string
		expectError          bool
		expectedErrorMsg     string
		expectedService      string
		expectedToken        string
		expectedSaveFilePath string
		expectedHelp         bool
	}{
		{
			name:                 "ValidArgsWithFilePath_Normal",
			args:                 []string{"-operation", "retrieve", "-service", testService, "-token", testToken, "-save-file-path", testSaveFilePath},
			expectError:          false,
			expectedService:      testService,
			expectedToken:        testToken,
			expectedSaveFilePath: testSaveFilePath,
			expectedHelp:         false,
		},
		{
			name:                 "ValidArgsWithoutFilePath_Normal",
			args:                 []string{"-operation", "retrieve", "-service", testService, "-token", testToken},
			expectError:          false,
			expectedService:      testService,
			expectedToken:        testToken,
			expectedSaveFilePath: "",
			expectedHelp:         false,
		},
		{
			name:                 "HelpFlag_Normal",
			args:                 []string{"-help"},
			expectError:          false,
			expectedService:      "",
			expectedToken:        "",
			expectedSaveFilePath: "",
			expectedHelp:         true,
		},
		{
			name:             "MissingOperation_Error",
			args:             []string{"-service", testService, "-token", testToken},
			expectError:      true,
			expectedErrorMsg: "operationパラメータは必須です",
		},
		{
			name:             "MissingService_Error",
			args:             []string{"-operation", "retrieve", "-token", testToken},
			expectError:      true,
			expectedErrorMsg: "serviceパラメータは必須です",
		},
		{
			name:             "MissingToken_Error",
			args:             []string{"-operation", "retrieve", "-service", testService},
			expectError:      true,
			expectedErrorMsg: "retrieveオペレーションの場合、tokenパラメータは必須です",
		},
		{
			name:             "UnsupportedService_Error",
			args:             []string{"-operation", "retrieve", "-service", "gitlab", "-token", testToken},
			expectError:      true,
			expectedErrorMsg: "サポートされていないサービスタイプです: gitlab",
		},
		{
			name:             "EmptyService_Error",
			args:             []string{"-operation", "retrieve", "-service", "", "-token", testToken},
			expectError:      true,
			expectedErrorMsg: "serviceパラメータは必須です",
		},
		{
			name:             "EmptyToken_Error",
			args:             []string{"-operation", "retrieve", "-service", testService, "-token", ""},
			expectError:      true,
			expectedErrorMsg: "retrieveオペレーションの場合、tokenパラメータは必須です",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// フラグをリセット
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// テスト用の引数を設定
			oldArgs := os.Args
			os.Args = append([]string{"test"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			cfg, err := ParseFlags()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, tt.expectedService, cfg.Service)
				assert.Equal(t, tt.expectedToken, cfg.Token)
				assert.Equal(t, tt.expectedSaveFilePath, cfg.SaveFilePath)
				assert.Equal(t, tt.expectedHelp, cfg.Help)
			}
		})
	}
}

func TestPrintUsage_Normal(t *testing.T) {
	// PrintUsage関数は標準出力に出力するため、実際の出力内容をテストするのは困難
	// ここでは関数が正常に実行されることを確認
	assert.NotPanics(t, func() {
		PrintUsage()
	})
}

// #==============================================================#
// ##          Config Struct Tests                               ##
// #==============================================================#

func TestConfig_Struct_Normal(t *testing.T) {
	const (
		testService      = "github"
		testToken        = "test-token"
		testSaveFilePath = "/tmp/test-output.json"
	)

	tests := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			name: "AllFieldsSet_Normal",
			config: Config{
				Service:      testService,
				Token:        testToken,
				SaveFilePath: testSaveFilePath,
				Help:         false,
			},
			expected: Config{
				Service:      testService,
				Token:        testToken,
				SaveFilePath: testSaveFilePath,
				Help:         false,
			},
		},
		{
			name: "MinimalFields_Normal",
			config: Config{
				Service: testService,
				Token:   testToken,
			},
			expected: Config{
				Service:      testService,
				Token:        testToken,
				SaveFilePath: "",
				Help:         false,
			},
		},
		{
			name: "HelpOnly_Normal",
			config: Config{
				Help: true,
			},
			expected: Config{
				Service:      "",
				Token:        "",
				SaveFilePath: "",
				Help:         true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.Service, tt.config.Service)
			assert.Equal(t, tt.expected.Token, tt.config.Token)
			assert.Equal(t, tt.expected.SaveFilePath, tt.config.SaveFilePath)
			assert.Equal(t, tt.expected.Help, tt.config.Help)
		})
	}
}

// #==============================================================#
// ##          Edge Case Tests                                   ##
// #==============================================================#

func TestParseFlags_EdgeCases_Normal(t *testing.T) {
	const (
		testService = "github"
		testToken   = "test-token"
	)

	tests := []struct {
		name             string
		args             []string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name:        "LongTokenValue_Normal",
			args:        []string{"-operation", "retrieve", "-service", testService, "-token", "ghp_1234567890abcdefghijklmnopqrstuvwxyz1234567890"},
			expectError: false,
		},
		{
			name:        "LongFilePathValue_Normal",
			args:        []string{"-operation", "retrieve", "-service", testService, "-token", testToken, "-save-file-path", "/very/long/path/to/some/directory/structure/output.json"},
			expectError: false,
		},
		{
			name:             "ServiceCaseSensitive_Normal",
			args:             []string{"-operation", "retrieve", "-service", "GitHub", "-token", testToken},
			expectError:      true,
			expectedErrorMsg: "サポートされていないサービスタイプです: GitHub",
		},
		{
			name:        "ServiceLowerCase_Normal",
			args:        []string{"-operation", "retrieve", "-service", "github", "-token", testToken},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// フラグをリセット
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// テスト用の引数を設定
			oldArgs := os.Args
			os.Args = append([]string{"test"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			cfg, err := ParseFlags()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
			}
		})
	}
}

// #==============================================================#
// ##          Integration Tests                                 ##
// #==============================================================#

func TestParseFlags_Integration_Normal(t *testing.T) {
	const (
		testService      = "github"
		testToken        = "ghp_test1234567890abcdef"
		testSaveFilePath = "/tmp/integration-test-output.json"
	)

	// 実際のCLI使用パターンをテスト
	tests := []struct {
		name           string
		args           []string
		expectError    bool
		validateConfig func(*testing.T, *Config)
	}{
		{
			name:        "TypicalUsageWithFile_Normal",
			args:        []string{"-operation", "retrieve", "-service", testService, "-token", testToken, "-save-file-path", testSaveFilePath},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				assert.Equal(t, testService, cfg.Service)
				assert.Equal(t, testToken, cfg.Token)
				assert.Equal(t, testSaveFilePath, cfg.SaveFilePath)
				assert.False(t, cfg.Help)
			},
		},
		{
			name:        "TypicalUsageWithoutFile_Normal",
			args:        []string{"-operation", "retrieve", "-service", testService, "-token", testToken},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				assert.Equal(t, testService, cfg.Service)
				assert.Equal(t, testToken, cfg.Token)
				assert.Empty(t, cfg.SaveFilePath)
				assert.False(t, cfg.Help)
			},
		},
		{
			name:        "HelpRequest_Normal",
			args:        []string{"-help"},
			expectError: false,
			validateConfig: func(t *testing.T, cfg *Config) {
				assert.True(t, cfg.Help)
				// ヘルプの場合、他のフィールドは検証しない
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// フラグをリセット
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			// テスト用の引数を設定
			oldArgs := os.Args
			os.Args = append([]string{"test"}, tt.args...)
			defer func() { os.Args = oldArgs }()

			cfg, err := ParseFlags()

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, cfg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				if tt.validateConfig != nil {
					tt.validateConfig(t, cfg)
				}
			}
		})
	}
}
