package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ensureEditableMap(value any) (map[string]any, error) {
	switch v := value.(type) {
	case map[string]any:
		return v, nil
	case []any:
		return nil, errors.New("edit-file は単一ドキュメントのYAMLのみ対応しています")
	default:
		return nil, errors.New("YAMLをマップオブジェクトとして解釈できません")
	}
}

func applyKeyValuePairs(root map[string]any, pairs []KeyValuePair) error {
	for _, pair := range pairs {
		if strings.TrimSpace(pair.Key) == "" {
			return errors.New("キーが空の値が含まれています")
		}
		if err := setNestedValue(root, pair.Key, normalize(pair.Value)); err != nil {
			return err
		}
	}
	return nil
}

func setNestedValue(root map[string]any, dottedKey string, value any) error {
	trimmed := strings.TrimSpace(dottedKey)
	if trimmed == "" {
		return errors.New("更新対象のキーが空です")
	}
	partsRaw := strings.Split(dottedKey, ".")
	parts := make([]string, 0, len(partsRaw))
	for _, raw := range partsRaw {
		segment := strings.TrimSpace(raw)
		if segment == "" {
			return fmt.Errorf("キー '%s' に空のセグメントが含まれています", dottedKey)
		}
		parts = append(parts, segment)
	}
	_, err := setNestedValueRecursive(root, parts, dottedKey, value)
	return err
}

func setNestedValueRecursive(node interface{}, parts []string, originalKey string, value any) (interface{}, error) {
	if len(parts) == 0 {
		return node, fmt.Errorf("キー '%s' の解析に失敗しました", originalKey)
	}
	part := parts[0]
	isLast := len(parts) == 1
	idx, isIndex := parseIndex(part)

	switch current := node.(type) {
	case map[string]any:
		if isIndex {
			return node, fmt.Errorf("キー '%s' のセグメント '%s' は配列インデックスですが、ここはオブジェクトです", originalKey, part)
		}
		if isLast {
			current[part] = value
			return current, nil
		}
		child, exists := current[part]
		if !exists || child == nil {
			child = defaultContainer(parts[1])
		}
		updatedChild, err := setNestedValueRecursive(child, parts[1:], originalKey, value)
		if err != nil {
			return node, err
		}
		current[part] = updatedChild
		return current, nil
	case []any:
		if !isIndex {
			return node, fmt.Errorf("キー '%s' のセグメント '%s' はキーですが、ここは配列です", originalKey, part)
		}
		if idx < 0 {
			return node, fmt.Errorf("配列インデックスは0以上で指定してください: %s", part)
		}
		arr := current
		if idx >= len(arr) {
			newArr := make([]any, idx+1)
			copy(newArr, arr)
			arr = newArr
		}
		if isLast {
			arr[idx] = value
			return arr, nil
		}
		child := arr[idx]
		if child == nil {
			child = defaultContainer(parts[1])
		}
		updatedChild, err := setNestedValueRecursive(child, parts[1:], originalKey, value)
		if err != nil {
			return node, err
		}
		arr[idx] = updatedChild
		return arr, nil
	default:
		return node, fmt.Errorf("キー '%s' の途中にオブジェクト以外の値があります", originalKey)
	}
}

func defaultContainer(nextSegment string) interface{} {
	if nextSegment == "" {
		return map[string]any{}
	}
	if _, isIndex := parseIndex(strings.TrimSpace(nextSegment)); isIndex {
		return []any{}
	}
	return map[string]any{}
}

func parseIndex(segment string) (int, bool) {
	if segment == "" {
		return 0, false
	}
	for _, r := range segment {
		if r < '0' || r > '9' {
			return 0, false
		}
	}
	idx, err := strconv.Atoi(segment)
	if err != nil {
		return 0, false
	}
	return idx, true
}

func parseKeyValueList(raw string) ([]KeyValuePair, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("--key-value-list に有効な値がありません")
	}

	splitter := func(r rune) bool {
		return r == ',' || r == '\n'
	}
	segments := strings.FieldsFunc(raw, splitter)
	var pairs []KeyValuePair
	for _, segment := range segments {
		entry := strings.TrimSpace(segment)
		if entry == "" {
			continue
		}
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("key=value 形式で指定してください: %s", entry)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("キーが空の組み合わせがあります: %s", entry)
		}
		valueStr := strings.TrimSpace(parts[1])
		value, err := parseValue(valueStr)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, KeyValuePair{Key: key, Value: value})
	}

	if len(pairs) == 0 {
		return nil, fmt.Errorf("--key-value-list に有効なエントリがありません")
	}

	return pairs, nil
}

func parseValue(raw string) (any, error) {
	if raw == "" {
		return "", nil
	}
	var parsed any
	if err := yaml.Unmarshal([]byte(raw), &parsed); err != nil || parsed == nil {
		return raw, nil
	}
	return parsed, nil
}

func (s *Service) writeYAML(path string, value any) error {
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(value); err != nil {
		encoder.Close()
		return fmt.Errorf("YAMLの再構築に失敗しました: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("YAMLエンコーダーの終了に失敗しました: %w", err)
	}
	if err := s.fileAccessor.WriteFile(path, buf.Bytes()); err != nil {
		return fmt.Errorf("YAMLファイルの書き込みに失敗しました: %w", err)
	}
	return nil
}
