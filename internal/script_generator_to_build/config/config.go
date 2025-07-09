package config

import "path/filepath"

// AppConfig はアプリケーションの設定を表します
type AppConfig struct {
	PackageName string
	ShowHelp    bool

	// 新しいフィールド（デフォルト値で後方互換性を保つ）
	BaseDir    string // デフォルト: "/home/nov/devbox"
	CLIDir     string // デフォルト: "cmd/cli"
	ScriptsDir string // デフォルト: "scripts"
	OutputDir  string // デフォルト: "./pkg/bin"
}

// SetDefaults はデフォルト値を設定します
func (c *AppConfig) SetDefaults() {
	if c.BaseDir == "" {
		c.BaseDir = "/home/nov/devbox"
	}
	if c.CLIDir == "" {
		c.CLIDir = "cmd/cli"
	}
	if c.ScriptsDir == "" {
		c.ScriptsDir = "scripts"
	}
	if c.OutputDir == "" {
		c.OutputDir = "./pkg/bin"
	}
}

// GetCLIPath はCLIディレクトリの完全パスを返します
func (c *AppConfig) GetCLIPath() string {
	return filepath.Join(c.BaseDir, c.CLIDir)
}

// GetScriptsPath はスクリプトディレクトリの完全パスを返します
func (c *AppConfig) GetScriptsPath() string {
	return filepath.Join(c.BaseDir, c.ScriptsDir)
}
