package config

import (
	"os"
	"path/filepath"
)

// AppConfig はアプリケーションの設定を表します
type AppConfig struct {
	PackageName string
	ShowHelp    bool

	// 新しいフィールド（デフォルト値で後方互換性を保つ）
	BaseDir    string // デフォルト: ワーキングディレクトリ
	CLIDir     string // デフォルト: "cmd/cli"
	ScriptsDir string // デフォルト: "scripts"
	OutputDir  string // デフォルト: "./pkg/bin/cli"
}

// SetDefaults はデフォルト値を設定します
func (c *AppConfig) SetDefaults() {
	if c.BaseDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			// エラーハンドリング: デフォルトで現在のディレクトリを使用
			c.BaseDir = "."
		} else {
			c.BaseDir = wd
		}
	}
	if c.CLIDir == "" {
		c.CLIDir = "cmd/cli"
	}
	if c.ScriptsDir == "" {
		c.ScriptsDir = "scripts"
	}
	if c.OutputDir == "" {
		c.OutputDir = "./pkg/bin/cli"
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
