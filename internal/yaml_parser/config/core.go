package config

import (
	"fmt"
	"strings"
)

const (
	OperationRead  = "read"
	OperationParse = "parse"
)

type Config struct {
	Operation   string
	FilePath    string
	YAMLContent string
}

func NewConfig(operation, filePath, yamlContent string) (*Config, error) {
	op := strings.ToLower(strings.TrimSpace(operation))
	if op == "" {
		return nil, fmt.Errorf("--operation を指定してください (read|parse)")
	}

	cfg := &Config{
		Operation:   op,
		FilePath:    strings.TrimSpace(filePath),
		YAMLContent: yamlContent,
	}

	switch op {
	case OperationRead:
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("--file-path を指定してください (--operation=read)")
		}
	case OperationParse:
		if strings.TrimSpace(yamlContent) == "" {
			return nil, fmt.Errorf("--yaml-content を指定してください (--operation=parse)")
		}
	default:
		return nil, fmt.Errorf("--operation には read または parse を指定してください (指定値: %s)", operation)
	}

	return cfg, nil
}
