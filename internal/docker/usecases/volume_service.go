package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// syncVolumesIntoCompose は指定したvolume-keyの値をもとに、サービスのvolumesセクションを書き換える
func syncVolumesIntoCompose(envPath, composePath, volumeKey, serviceName string) error {
	if envPath == "" {
		return fmt.Errorf("env-yaml-path は必須です")
	}
	if composePath == "" {
		return fmt.Errorf("compose-path は必須です")
	}
	if volumeKey == "" {
		return fmt.Errorf("volume-key は必須です")
	}
	if serviceName == "" {
		return fmt.Errorf("service は必須です")
	}

	absEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		return fmt.Errorf("envファイルパスの解決に失敗しました: %w", err)
	}
	absComposePath, err := filepath.Abs(composePath)
	if err != nil {
		return fmt.Errorf("composeファイルパスの解決に失敗しました: %w", err)
	}

	envContent, err := os.ReadFile(absEnvPath)
	if err != nil {
		return fmt.Errorf("envファイルの読み込みに失敗しました: %w", err)
	}

	entries, err := parseEnvEntries(string(envContent))
	if err != nil {
		return err
	}

	entry := findEnvEntry(entries, volumeKey)
	if entry == nil {
		return fmt.Errorf("%s が %s に存在しません", volumeKey, absEnvPath)
	}
	node := resolveAlias(entry.Node)
	if node == nil {
		return fmt.Errorf("%s の値を解析できません", volumeKey)
	}
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("%s は配列である必要があります", volumeKey)
	}

	composeContentBytes, err := os.ReadFile(absComposePath)
	if err != nil {
		return fmt.Errorf("docker-compose.yml の読み込みに失敗しました: %w", err)
	}
	composeContent := string(composeContentBytes)
	lines := splitComposeLines(composeContent)

	updatedLines, updated, err := updateServiceVolumes(lines, serviceName, node, volumeKey)
	if err != nil {
		return err
	}
	if !updated {
		return fmt.Errorf("%s の volumes セクションを更新できませんでした", serviceName)
	}

	nextContent := strings.Join(updatedLines, "\n")
	if !strings.HasSuffix(nextContent, "\n") {
		nextContent += "\n"
	}
	if nextContent != composeContent {
		if err := os.WriteFile(absComposePath, []byte(nextContent), 0o644); err != nil {
			return fmt.Errorf("docker-compose.yml への書き込みに失敗しました: %w", err)
		}
	}

	return nil
}

func updateServiceVolumes(lines []string, serviceName string, node *yaml.Node, volumeKey string) ([]string, bool, error) {
	currentService := ""
	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		indentLen := len(leadingWhitespace(line))

		if indentLen == 2 && strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, " ") {
			currentService = strings.TrimSuffix(trimmed, ":")
		}
		if currentService != serviceName {
			continue
		}
		if trimmed != "volumes:" {
			continue
		}

		baseIndent := leadingWhitespace(line)
		start := idx + 1
		end := start
		for end < len(lines) {
			trimmedLine := strings.TrimSpace(lines[end])
			if trimmedLine == "" {
				end++
				continue
			}
			if len(leadingWhitespace(lines[end])) <= len(baseIndent) {
				break
			}
			end++
		}

		volumeLines, err := buildVolumeLinesFromNode(node, baseIndent+"  ", volumeKey)
		if err != nil {
			return nil, false, err
		}

		updated := append([]string{}, lines[:start]...)
		updated = append(updated, volumeLines...)
		updated = append(updated, lines[end:]...)
		return updated, true, nil
	}

	return nil, false, nil
}

func buildVolumeLinesFromNode(node *yaml.Node, indent, volumeKey string) ([]string, error) {
	node = resolveAlias(node)
	if node.Kind != yaml.SequenceNode {
		return nil, fmt.Errorf("%s は配列である必要があります", volumeKey)
	}
	if len(node.Content) == 0 {
		return nil, fmt.Errorf("%s の値が空です", volumeKey)
	}

	var lines []string
	for _, item := range node.Content {
		item = resolveAlias(item)
		switch item.Kind {
		case yaml.MappingNode:
			mappingLines, err := formatVolumeMapping(indent, item, volumeKey)
			if err != nil {
				return nil, err
			}
			lines = append(lines, mappingLines...)
		case yaml.ScalarNode:
			value, err := renderYAMLValue(item)
			if err != nil {
				return nil, err
			}
			lines = append(lines, fmt.Sprintf("%s- %s", indent, formatEnvValue(value)))
		default:
			return nil, fmt.Errorf("%s の各要素はマッピングまたはスカラーである必要があります", volumeKey)
		}
	}

	return lines, nil
}

func formatVolumeMapping(indent string, node *yaml.Node, volumeKey string) ([]string, error) {
	if len(node.Content) == 0 {
		return []string{fmt.Sprintf("%s- {}", indent)}, nil
	}
	if len(node.Content)%2 != 0 {
		return nil, fmt.Errorf("%s のボリューム定義が不正です", volumeKey)
	}

	lines := make([]string, 0, len(node.Content)/2)
	for i := 0; i+1 < len(node.Content); i += 2 {
		keyNode := resolveAlias(node.Content[i])
		valueNode := resolveAlias(node.Content[i+1])
		if keyNode.Kind != yaml.ScalarNode {
			return nil, fmt.Errorf("%s のボリューム定義のキーはスカラーである必要があります", volumeKey)
		}
		value, err := renderYAMLValue(valueNode)
		if err != nil {
			return nil, err
		}
		formatted := formatEnvValue(value)
		if i == 0 {
			lines = append(lines, fmt.Sprintf("%s- %s: %s", indent, keyNode.Value, formatted))
		} else {
			lines = append(lines, fmt.Sprintf("%s  %s: %s", indent, keyNode.Value, formatted))
		}
	}

	return lines, nil
}
