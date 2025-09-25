package config

import (
	"os"
	"path/filepath"
)

// ServiceConfig はアプリケーションの設定を表します
type ServiceConfig struct {
	PackageName string
	ShowHelp    bool

	// 新しいフィールド（デフォルト値で後方互換性を保つ）
	BaseDir    string // デフォルト: ワーキングディレクトリ
	CLIDir     string // デフォルト: "cmd/cli"
	ScriptsDir string // デフォルト: "scripts"
	OutputDir  string // デフォルト: "./pkg/bin/cli"
}

// SetDefaults はデフォルト値を設定します
func (c *ServiceConfig) SetDefaults() {
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
func (c *ServiceConfig) GetCLIPath() string {
	return filepath.Join(c.BaseDir, c.CLIDir)
}

// GetScriptsPath はスクリプトディレクトリの完全パスを返します
func (c *ServiceConfig) GetScriptsPath() string {
	return filepath.Join(c.BaseDir, c.ScriptsDir)
}
