package usecases

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// EnvEntry は単一の環境変数キーと値を表す
type EnvEntry struct {
	Key   string
	Value string
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

func parseEnvEntries(content string) ([]EnvEntry, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var entries []EnvEntry
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		key, rawValue, ok := splitKeyValueLine(line)
		if !ok {
			continue
		}
		withoutComment := stripInlineComment(rawValue)
		normalized := stripQuotes(strings.TrimSpace(withoutComment))
		entries = append(entries, EnvEntry{Key: key, Value: normalized})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envファイルの解析に失敗しました: %w", err)
	}

	return entries, nil
}

func splitKeyValueLine(line string) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx == -1 {
		return "", "", false
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" || !isValidKey(key) {
		return "", "", false
	}

	value := ""
	if idx+1 < len(line) {
		value = line[idx+1:]
	}

	return key, value, true
}

func isValidKey(key string) bool {
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func stripInlineComment(value string) string {
	inSingle := false
	inDouble := false
	var builder strings.Builder

	for i := 0; i < len(value); i++ {
		ch := value[i]
		switch ch {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
			builder.WriteByte(ch)
		case '"':
			if !inSingle {
				escaped := i > 0 && value[i-1] == '\\'
				if !escaped {
					inDouble = !inDouble
				}
			}
			builder.WriteByte(ch)
		case '#':
			if !inSingle && !inDouble {
				return strings.TrimRightFunc(builder.String(), unicode.IsSpace)
			}
			builder.WriteByte(ch)
		default:
			builder.WriteByte(ch)
		}
	}

	return strings.TrimRightFunc(builder.String(), unicode.IsSpace)
}

func stripQuotes(value string) string {
	trimmed := strings.TrimSpace(value)
	if len(trimmed) < 2 {
		return trimmed
	}

	first := trimmed[0]
	last := trimmed[len(trimmed)-1]
	if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
		return trimmed[1 : len(trimmed)-1]
	}
	return trimmed
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
