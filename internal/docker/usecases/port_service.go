package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// syncPortsIntoCompose は指定したport-keyの値をもとに、サービスのportsおよびtsdproxy.container_portを更新する
func syncPortsIntoCompose(envPath, composePath, portKey, serviceName string) error {
	if envPath == "" {
		return fmt.Errorf("env-yaml-path は必須です")
	}
	if composePath == "" {
		return fmt.Errorf("compose-path は必須です")
	}
	if portKey == "" {
		return fmt.Errorf("port-key は必須です")
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
	envLookup := entriesToMap(entries)
	portValue := strings.TrimSpace(envLookup[portKey])
	if portValue == "" {
		return fmt.Errorf("%s が %s に存在しません", portKey, absEnvPath)
	}

	composeContentBytes, err := os.ReadFile(absComposePath)
	if err != nil {
		return fmt.Errorf("docker-compose.yml の読み込みに失敗しました: %w", err)
	}
	composeContent := string(composeContentBytes)
	lines := splitComposeLines(composeContent)

	portUpdated, err := updateServicePorts(lines, serviceName, portValue)
	if err != nil {
		return err
	}
	if !portUpdated {
		return fmt.Errorf("%s の ports セクションを更新できませんでした", serviceName)
	}

	labelUpdated, err := updateServiceTsdproxyPort(lines, serviceName, portValue)
	if err != nil {
		return err
	}
	if !labelUpdated {
		return fmt.Errorf("%s の labels 内に tsdproxy.container_port が見つかりません", serviceName)
	}

	nextContent := strings.Join(lines, "\n")
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

func entriesToMap(entries []EnvEntry) map[string]string {
	lookup := make(map[string]string, len(entries))
	for _, entry := range entries {
		lookup[entry.Key] = entry.Value
	}
	return lookup
}

func updateServicePorts(lines []string, serviceName, portValue string) (bool, error) {
	currentService := ""
	inPorts := false
	portsIndentLen := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		indentLen := len(leadingWhitespace(line))

		if indentLen == 2 && strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, " ") {
			currentService = strings.TrimSuffix(trimmed, ":")
			inPorts = false
		}

		if currentService != serviceName {
			continue
		}

		if !inPorts {
			if trimmed == "ports:" {
				inPorts = true
				portsIndentLen = indentLen
			}
			continue
		}

		if trimmed == "" {
			continue
		}
		if indentLen <= portsIndentLen {
			inPorts = false
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "-") {
			indent := leadingWhitespace(line)
			lines[i] = fmt.Sprintf("%s- \"%s:%s\"", indent, portValue, portValue)
			return true, nil
		}
	}
	return false, nil
}

func updateServiceTsdproxyPort(lines []string, serviceName, portValue string) (bool, error) {
	currentService := ""
	inLabels := false
	labelsIndentLen := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		indentLen := len(leadingWhitespace(line))

		if indentLen == 2 && strings.HasSuffix(trimmed, ":") && !strings.Contains(trimmed, " ") {
			currentService = strings.TrimSuffix(trimmed, ":")
			if currentService != serviceName {
				inLabels = false
			}
		}

		if currentService != serviceName {
			continue
		}

		if !inLabels {
			if trimmed == "labels:" {
				inLabels = true
				labelsIndentLen = indentLen
			}
			continue
		}

		if trimmed == "" {
			continue
		}
		if indentLen <= labelsIndentLen {
			inLabels = false
			continue
		}
		if strings.HasPrefix(trimmed, "tsdproxy.container_port:") {
			indent := leadingWhitespace(line)
			lines[i] = fmt.Sprintf("%stsdproxy.container_port: %s", indent, portValue)
			return true, nil
		}
	}
	return false, nil
}
