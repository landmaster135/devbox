package config

import (
	"fmt"
	"strings"
)

const (
	OperationRead     = "read"
	OperationParse    = "parse"
	OperationEditFile = "edit-file"
)

type Config struct {
	Operation    string
	FilePath     string
	YAMLContent  string
	KeyValueList string
}

func NewConfig(operation, filePath, yamlContent, keyValueList string) (*Config, error) {
	op := strings.ToLower(strings.TrimSpace(operation))
	if op == "" {
		return nil, fmt.Errorf("--operation を指定してください (read|parse|edit-file)")
	}

	cfg := &Config{
		Operation:    op,
		FilePath:     strings.TrimSpace(filePath),
		YAMLContent:  yamlContent,
		KeyValueList: keyValueList,
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
	case OperationEditFile:
		if cfg.FilePath == "" {
			return nil, fmt.Errorf("--file-path を指定してください (--operation=edit-file)")
		}
		if strings.TrimSpace(cfg.KeyValueList) == "" {
			return nil, fmt.Errorf("--key-value-list を指定してください (--operation=edit-file)")
		}
	default:
		return nil, fmt.Errorf("--operation には read, parse, edit-file のいずれかを指定してください (指定値: %s)", operation)
	}

	return cfg, nil
}
