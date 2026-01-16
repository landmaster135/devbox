package config

import (
	"os"
	"path/filepath"
)

// ServiceConfig はアプリケーションの設定を表します
type ServiceConfig struct {
	PackageName string
	ShowHelp    bool
	BaseDir     string
	CLIDir      string
	ScriptsDir  string
	OutputDir   string
}

const (
	defaultCLIDir     = "cmd/cli"
	defaultScriptsDir = "scripts/build"
	defaultOutputDir  = "./pkg/bin/cli"
)

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
		c.CLIDir = defaultCLIDir
	}
	if c.ScriptsDir == "" {
		c.ScriptsDir = defaultScriptsDir
	}
	if c.OutputDir == "" {
		c.OutputDir = defaultOutputDir
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
