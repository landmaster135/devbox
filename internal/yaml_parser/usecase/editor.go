package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func applyEditsWithOriginalOrder(data []byte, pairs []KeyValuePair) ([]byte, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var doc yaml.Node
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("YAMLの解析に失敗しました: %w", err)
	}
	if len(doc.Content) == 0 {
		return nil, errors.New("YAMLの内容が空です")
	}

	// Ensure single document
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, errors.New("edit-file は単一ドキュメントのYAMLのみ対応しています")
	}

	root := doc.Content[0]
	if root == nil || root.Kind != yaml.MappingNode {
		return nil, errors.New("YAMLをマップオブジェクトとして解釈できません")
	}

	for _, pair := range pairs {
		if err := setNodeValue(root, pair.Key, pair.Value); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&doc); err != nil {
		encoder.Close()
		return nil, fmt.Errorf("YAMLの再構築に失敗しました: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("YAMLエンコーダーの終了に失敗しました: %w", err)
	}

	return buf.Bytes(), nil
}

func setNodeValue(root *yaml.Node, dottedKey string, value interface{}) error {
	segments, err := splitKeySegments(dottedKey)
	if err != nil {
		return err
	}
	valueNode, err := buildValueNode(value)
	if err != nil {
		return err
	}
	return setNodeValueRecursive(root, segments, valueNode, dottedKey)
}

func splitKeySegments(dottedKey string) ([]string, error) {
	trimmed := strings.TrimSpace(dottedKey)
	if trimmed == "" {
		return nil, errors.New("更新対象のキーが空です")
	}
	raw := strings.Split(trimmed, ".")
	segments := make([]string, 0, len(raw))
	for _, part := range raw {
		segment := strings.TrimSpace(part)
		if segment == "" {
			return nil, fmt.Errorf("キー '%s' に空のセグメントが含まれています", dottedKey)
		}
		segments = append(segments, segment)
	}
	return segments, nil
}

func setNodeValueRecursive(node *yaml.Node, parts []string, valueNode *yaml.Node, originalKey string) error {
	if len(parts) == 0 {
		return fmt.Errorf("キー '%s' の解析に失敗しました", originalKey)
	}

	currentPart := parts[0]
	isLast := len(parts) == 1
	index, isIndex := parseIndex(currentPart)

	switch node.Kind {
	case yaml.MappingNode:
		if isIndex {
			return fmt.Errorf("キー '%s' のセグメント '%s' は配列インデックスですが、ここはオブジェクトです", originalKey, currentPart)
		}
		keyIdx := findKeyIndex(node, currentPart)
		if isLast {
			if keyIdx >= 0 {
				node.Content[keyIdx+1] = cloneNode(valueNode)
			} else {
				node.Content = append(node.Content, newScalarKeyNode(currentPart), cloneNode(valueNode))
			}
			return nil
		}

		var child *yaml.Node
		if keyIdx >= 0 {
			child = node.Content[keyIdx+1]
		} else {
			child = defaultContainerNode(parts[1])
			node.Content = append(node.Content, newScalarKeyNode(currentPart), child)
		}

		if err := ensureNodeKind(child, parts[1], originalKey); err != nil {
			return err
		}
		return setNodeValueRecursive(child, parts[1:], valueNode, originalKey)

	case yaml.SequenceNode:
		if !isIndex {
			return fmt.Errorf("キー '%s' のセグメント '%s' はキーですが、ここは配列です", originalKey, currentPart)
		}
		if index < 0 {
			return fmt.Errorf("配列インデックスは0以上で指定してください: %s", currentPart)
		}

		// Extend slice if needed
		for len(node.Content) <= index {
			node.Content = append(node.Content, newNullNode())
		}

		if isLast {
			node.Content[index] = cloneNode(valueNode)
			return nil
		}

		child := node.Content[index]
		if child == nil || child.Kind == 0 || child.Kind == yaml.ScalarNode && child.Tag == "!!null" {
			child = defaultContainerNode(parts[1])
			node.Content[index] = child
		}
		if err := ensureNodeKind(child, parts[1], originalKey); err != nil {
			return err
		}
		return setNodeValueRecursive(child, parts[1:], valueNode, originalKey)

	default:
		return fmt.Errorf("キー '%s' の途中にオブジェクト以外の値があります", originalKey)
	}
}

func ensureNodeKind(node *yaml.Node, nextSegment string, originalKey string) error {
	if node == nil {
		return fmt.Errorf("キー '%s' の途中でノードが無効です", originalKey)
	}

	if isIndexSegment(nextSegment) {
		if node.Kind == yaml.SequenceNode {
			return nil
		}
		if node.Kind == 0 || (node.Kind == yaml.ScalarNode && node.Tag == "!!null") {
			*node = *defaultContainerNode(nextSegment)
			return nil
		}
		return fmt.Errorf("キー '%s' の途中に配列ではない値があります", originalKey)
	}

	if node.Kind == yaml.MappingNode {
		return nil
	}
	if node.Kind == 0 || (node.Kind == yaml.ScalarNode && node.Tag == "!!null") {
		*node = *defaultContainerNode(nextSegment)
		return nil
	}
	return fmt.Errorf("キー '%s' の途中にオブジェクト以外の値があります", originalKey)
}

func findKeyIndex(node *yaml.Node, key string) int {
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return i
		}
	}
	return -1
}

func defaultContainerNode(nextSegment string) *yaml.Node {
	if isIndexSegment(nextSegment) {
		return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: []*yaml.Node{}}
	}
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: []*yaml.Node{}}
}

func newScalarKeyNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func newNullNode() *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null", Value: "null"}
}

func isIndexSegment(segment string) bool {
	_, ok := parseIndex(segment)
	return ok
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

func buildValueNode(value interface{}) (*yaml.Node, error) {
	raw, err := yaml.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("値のエンコードに失敗しました: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(raw))
	var doc yaml.Node
	if err := decoder.Decode(&doc); err != nil {
		return nil, fmt.Errorf("値の解析に失敗しました: %w", err)
	}
	if len(doc.Content) == 0 {
		return newNullNode(), nil
	}
	return cloneNode(doc.Content[0]), nil
}

func cloneNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	copy := *node
	if len(node.Content) > 0 {
		copy.Content = make([]*yaml.Node, len(node.Content))
		for i, child := range node.Content {
			copy.Content[i] = cloneNode(child)
		}
	}
	return &copy
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
