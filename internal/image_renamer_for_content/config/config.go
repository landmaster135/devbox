package config

import (
	"runtime"
	"strings"
)

// Config はimage-renamer-for-contentツールの設定を保持します。
type Config struct {
	SrcDir     string
	SortByName bool
	SortByTime bool
	ContentID  string
	Suffix     string
	Delimiter  string
	Digits     int
	Start      int
	Recursive  bool
	Workers    int
	Operation  string
}

// Normalize は設定値にデフォルト値を適用し、前後の空白を除去します。
func (c *Config) Normalize() {
	c.ContentID = strings.TrimSpace(c.ContentID)
	c.Suffix = strings.TrimSpace(c.Suffix)
	c.Delimiter = strings.TrimSpace(c.Delimiter)
	c.Operation = strings.TrimSpace(c.Operation)

	if c.Suffix == "" {
		c.Suffix = "01"
	}

	if c.Digits <= 0 {
		c.Digits = 4
	}

	if c.Start <= 0 {
		c.Start = 1
	}

	if c.Workers <= 0 {
		c.Workers = DefaultWorkers()
	}
}

// DefaultWorkers はデフォルトのワーカー数を計算します。
func DefaultWorkers() int {
	workers := runtime.NumCPU() - 1
	if workers < 1 {
		return 1
	}
	return workers
}
