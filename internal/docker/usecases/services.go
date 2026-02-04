package usecases

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// EnvEntry は単一の環境変数キーと値を表す
type EnvEntry struct {
	Key   string
	Value string
	Node  *yaml.Node
}

// EnvSyncService はenv.ymlの内容をdocker-compose.ymlへ同期する
type EnvSyncService struct{}

// NewEnvSyncService はEnvSyncServiceを作成する
func NewEnvSyncService() *EnvSyncService {
	return &EnvSyncService{}
}

// SyncEnvIntoCompose はenv.ymlから環境変数を読み、docker-compose.ymlのenvironmentセクションを書き換える
func (s *EnvSyncService) SyncEnvIntoCompose(envPath, composePath string) (int, error) {
	if envPath == "" {
		return 0, fmt.Errorf("env-yaml-path は必須です")
	}
	if composePath == "" {
		return 0, fmt.Errorf("compose-path は必須です")
	}

	absEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		return 0, fmt.Errorf("envファイルパスの解決に失敗しました: %w", err)
	}
	absComposePath, err := filepath.Abs(composePath)
	if err != nil {
		return 0, fmt.Errorf("composeファイルパスの解決に失敗しました: %w", err)
	}

	envContent, err := os.ReadFile(absEnvPath)
	if err != nil {
		return 0, fmt.Errorf("envファイルの読み込みに失敗しました: %w", err)
	}

	entries, err := parseEnvEntries(string(envContent))
	if err != nil {
		return 0, err
	}
	if len(entries) == 0 {
		return 0, fmt.Errorf("%s には有効な環境変数が含まれていません", absEnvPath)
	}

	composeContent, err := os.ReadFile(absComposePath)
	if err != nil {
		return 0, fmt.Errorf("docker-compose.yml の読み込みに失敗しました: %w", err)
	}

	nextContent, err := injectEnvironmentBlock(string(composeContent), entries)
	if err != nil {
		return 0, err
	}

	if nextContent != string(composeContent) {
		if err := os.WriteFile(absComposePath, []byte(nextContent), 0o644); err != nil {
			return 0, fmt.Errorf("docker-compose.yml への書き込みに失敗しました: %w", err)
		}
	}

	return len(entries), nil
}

// SyncPortsIntoCompose は指定したport-keyの値をもとに、サービスのportsおよびtsdproxy.container_portを更新する
func (s *EnvSyncService) SyncPortsIntoCompose(envPath, composePath, portKey, serviceName string) error {
	return syncPortsIntoCompose(envPath, composePath, portKey, serviceName)
}

// SyncVolumesIntoCompose は指定したvolume-keyの値をもとに、サービスのvolumesセクションを書き換える
func (s *EnvSyncService) SyncVolumesIntoCompose(envPath, composePath, volumeKey, serviceName string) error {
	return syncVolumesIntoCompose(envPath, composePath, volumeKey, serviceName)
}

// SyncUserIntoCompose は指定したuser-keyの値をもとに、サービスのuserフィールドを書き換える
func (s *EnvSyncService) SyncUserIntoCompose(envPath, composePath, userKey, serviceName string) error {
	return syncUserIntoCompose(envPath, composePath, userKey, serviceName)
}

func parseEnvEntries(content string) ([]EnvEntry, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		return nil, fmt.Errorf("envファイルの解析に失敗しました: %w", err)
	}

	if len(root.Content) == 0 {
		return nil, nil
	}

	mapping := root.Content[0]
	if mapping == nil {
		return nil, nil
	}
	if mapping.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("envファイルのルートはマッピングである必要があります")
	}

	entries := make([]EnvEntry, 0, len(mapping.Content)/2)
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode == nil || keyNode.Kind != yaml.ScalarNode {
			continue
		}
		key := strings.TrimSpace(keyNode.Value)
		if key == "" || !isValidKey(key) {
			continue
		}

		valueNode := mapping.Content[i+1]
		resolvedNode := resolveAlias(valueNode)
		value, err := renderYAMLValue(resolvedNode)
		if err != nil {
			return nil, err
		}
		entries = append(entries, EnvEntry{Key: key, Value: value, Node: resolvedNode})
	}

	return entries, nil
}

func findEnvEntry(entries []EnvEntry, key string) *EnvEntry {
	for i := range entries {
		if entries[i].Key == key {
			return &entries[i]
		}
	}
	return nil
}

func isValidKey(key string) bool {
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func renderYAMLValue(node *yaml.Node) (string, error) {
	if node == nil {
		return "", nil
	}
	resolved := resolveAlias(node)
	switch resolved.Kind {
	case yaml.ScalarNode:
		if resolved.Tag == "!!null" && resolved.Value == "" {
			return "", nil
		}
		return resolved.Value, nil
	case yaml.MappingNode, yaml.SequenceNode:
		var builder strings.Builder
		if err := encodeNodeAsJSON(&builder, resolved); err != nil {
			return "", err
		}
		return builder.String(), nil
	case yaml.DocumentNode:
		if len(resolved.Content) == 0 {
			return "", nil
		}
		return renderYAMLValue(resolved.Content[0])
	default:
		return "", fmt.Errorf("サポートされていないYAMLノードです: kind=%d", resolved.Kind)
	}
}

func resolveAlias(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.AliasNode && node.Alias != nil {
		return node.Alias
	}
	return node
}

func encodeNodeAsJSON(builder *strings.Builder, node *yaml.Node) error {
	node = resolveAlias(node)
	switch node.Kind {
	case yaml.SequenceNode:
		builder.WriteByte('[')
		for idx, child := range node.Content {
			if idx > 0 {
				builder.WriteByte(',')
			}
			if err := encodeNodeAsJSON(builder, child); err != nil {
				return err
			}
		}
		builder.WriteByte(']')
		return nil
	case yaml.MappingNode:
		builder.WriteByte('{')
		pairIdx := 0
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := resolveAlias(node.Content[i])
			valueNode := node.Content[i+1]
			if keyNode.Kind != yaml.ScalarNode {
				continue
			}
			if pairIdx > 0 {
				builder.WriteByte(',')
			}
			pairIdx++
			builder.WriteString(strconv.Quote(keyNode.Value))
			builder.WriteByte(':')
			if err := encodeNodeAsJSON(builder, valueNode); err != nil {
				return err
			}
		}
		builder.WriteByte('}')
		return nil
	case yaml.ScalarNode:
		builder.WriteString(jsonScalarLiteral(node))
		return nil
	case yaml.DocumentNode:
		if len(node.Content) == 0 {
			builder.WriteString("null")
			return nil
		}
		return encodeNodeAsJSON(builder, node.Content[0])
	default:
		return fmt.Errorf("サポートされていないYAMLノードです: kind=%d", node.Kind)
	}
}

func jsonScalarLiteral(node *yaml.Node) string {
	if node == nil {
		return "null"
	}
	switch node.Tag {
	case "!!bool":
		lower := strings.ToLower(node.Value)
		if lower == "true" || lower == "false" {
			return lower
		}
		return "false"
	case "!!int", "!!float":
		return node.Value
	case "!!null":
		return "null"
	default:
		return strconv.Quote(node.Value)
	}
}

func injectEnvironmentBlock(composeContent string, entries []EnvEntry) (string, error) {
	lines := splitComposeLines(composeContent)

	envLineIndex, endIndex, baseIndent, err := findEnvironmentBlock(lines)
	if err != nil {
		return "", err
	}

	itemIndent := baseIndent + "  "
	newEnvLines := buildEnvironmentLines(entries, itemIndent)

	startIndex := envLineIndex + 1
	nextLines := append([]string{}, lines[:startIndex]...)
	nextLines = append(nextLines, newEnvLines...)
	nextLines = append(nextLines, lines[endIndex:]...)

	result := strings.Join(nextLines, "\n")
	if !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	return result, nil
}

func splitComposeLines(content string) []string {
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	return lines
}

func findEnvironmentBlock(lines []string) (int, int, string, error) {
	for idx, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "environment:") {
			continue
		}
		if strings.Contains(trimmed, "*") && !strings.Contains(trimmed, "&") {
			// アンカー参照のみの場合はスキップ
			continue
		}
		baseIndent := leadingWhitespace(line)
		start := idx + 1
		end := start
		for end < len(lines) {
			current := lines[end]
			trimmedCurrent := strings.TrimSpace(current)
			if trimmedCurrent == "" {
				break
			}
			if len(leadingWhitespace(current)) <= len(baseIndent) {
				break
			}
			end++
		}
		return idx, end, baseIndent, nil
	}
	return -1, -1, "", errors.New("docker-compose.yml に environment セクションが見つかりませんでした")
}

func buildEnvironmentLines(entries []EnvEntry, indent string) []string {
	lines := make([]string, 0, len(entries))
	for _, entry := range entries {
		lines = append(lines, fmt.Sprintf("%s- %s=%s", indent, entry.Key, formatEnvValue(entry.Value)))
	}
	return lines
}

func formatEnvValue(value string) string {
	if value == "" {
		return strconv.Quote(value)
	}
	for _, r := range value {
		if unicode.IsSpace(r) {
			return strconv.Quote(value)
		}
	}
	return value
}

func leadingWhitespace(value string) string {
	idx := 0
	for idx < len(value) {
		if value[idx] != ' ' && value[idx] != '\t' {
			break
		}
		idx++
	}
	return value[:idx]
}
