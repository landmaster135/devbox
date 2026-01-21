package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func syncUserIntoCompose(envPath, composePath, userKey, serviceName string) error {
	if envPath == "" {
		return fmt.Errorf("env-yaml-path は必須です")
	}
	if composePath == "" {
		return fmt.Errorf("compose-path は必須です")
	}
	if userKey == "" {
		return fmt.Errorf("user-key は必須です")
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
	entry := findEnvEntry(entries, userKey)
	if entry == nil {
		return fmt.Errorf("%s が %s に存在しません", userKey, absEnvPath)
	}
	if entry.Node == nil || entry.Node.Kind != yaml.ScalarNode {
		return fmt.Errorf("%s はスカラー値である必要があります", userKey)
	}
	userValue := entry.Value

	composeContentBytes, err := os.ReadFile(absComposePath)
	if err != nil {
		return fmt.Errorf("docker-compose.yml の読み込みに失敗しました: %w", err)
	}
	composeContent := string(composeContentBytes)
	lines := splitComposeLines(composeContent)

	updatedLines, updated, err := updateServiceUser(lines, serviceName, userValue)
	if err != nil {
		return err
	}
	if !updated {
		return fmt.Errorf("%s の user フィールドを更新できませんでした", serviceName)
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

func updateServiceUser(lines []string, serviceName, userValue string) ([]string, bool, error) {
	for idx := 0; idx < len(lines); idx++ {
		line := lines[idx]
		trimmed := strings.TrimSpace(line)
		indent := leadingWhitespace(line)
		if len(indent) == 2 && strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, " ") {
			current := strings.TrimSuffix(trimmed, ":")
			if current != serviceName {
				continue
			}

			childIndent := indent + "  "
			insertPos := idx + 1
			for j := idx + 1; j < len(lines); j++ {
				currentLine := lines[j]
				trimmedCurrent := strings.TrimSpace(currentLine)
				currentIndent := leadingWhitespace(currentLine)
				if trimmedCurrent == "" {
					insertPos = j
					continue
				}
				if len(currentIndent) <= len(indent) {
					insertPos = j
					break
				}
				insertPos = j + 1
				if strings.HasPrefix(trimmedCurrent, "user:") {
					lines[j] = fmt.Sprintf("%suser: %q", currentIndent, userValue)
					return lines, true, nil
				}
			}

			if insertPos < idx+1 {
				insertPos = idx + 1
			}
			newLine := fmt.Sprintf("%suser: %q", childIndent, userValue)
			next := append([]string{}, lines[:insertPos]...)
			next = append(next, newLine)
			next = append(next, lines[insertPos:]...)
			return next, true, nil
		}
	}
	return nil, false, nil
}
