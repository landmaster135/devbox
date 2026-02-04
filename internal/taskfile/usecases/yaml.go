package usecases

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func extractFieldPaths(data []byte) (map[string]struct{}, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	if len(node.Content) == 0 {
		return map[string]struct{}{}, nil
	}

	fields := make(map[string]struct{})
	collectFieldPaths(node.Content[0], "", fields)
	return fields, nil
}

func collectFieldPaths(node *yaml.Node, prefix string, fields map[string]struct{}) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i+1 < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			key := keyNode.Value
			var current string
			if prefix == "" {
				current = key
			} else {
				current = prefix + "." + key
			}
			fields[current] = struct{}{}
			collectFieldPaths(valueNode, current, fields)
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			collectFieldPaths(child, prefix, fields)
		}
	}
}

func readYAMLDocument(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var node yaml.Node
	if len(bytes.TrimSpace(data)) == 0 {
		node.Kind = yaml.DocumentNode
		node.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
		return &node, nil
	}

	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}

	if len(node.Content) == 0 {
		node.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}

	return &node, nil
}

func ensureDocumentRoot(doc *yaml.Node) *yaml.Node {
	if doc == nil {
		return nil
	}
	if doc.Kind != yaml.DocumentNode {
		return doc
	}
	if len(doc.Content) == 0 {
		doc.Content = []*yaml.Node{{Kind: yaml.MappingNode}}
	}
	if doc.Content[0] == nil {
		doc.Content[0] = &yaml.Node{Kind: yaml.MappingNode}
	}
	return doc.Content[0]
}

func writeYAMLDocument(path string, node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("書き込み対象のノードがありません")
	}
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(node); err != nil {
		return err
	}
	if err := encoder.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func fillMappingNode(reference, target *yaml.Node) bool {
	if reference == nil || target == nil {
		return false
	}
	if reference.Kind != yaml.MappingNode {
		return false
	}
	if target.Kind != yaml.MappingNode {
		if isBlankNode(target) {
			*target = *deepCopyNode(reference)
			return true
		}
		return false
	}

	changed := false
	for i := 0; i+1 < len(reference.Content); i += 2 {
		keyNode := reference.Content[i]
		valueNode := reference.Content[i+1]
		idx := findMappingKey(target, keyNode.Value)
		if idx == -1 {
			target.Content = append(target.Content, deepCopyNode(keyNode), deepCopyNode(valueNode))
			changed = true
			continue
		}

		targetValue := target.Content[idx+1]
		switch valueNode.Kind {
		case yaml.MappingNode:
			if targetValue.Kind != yaml.MappingNode {
				if isBlankNode(targetValue) {
					target.Content[idx+1] = deepCopyNode(valueNode)
					changed = true
				}
				continue
			}
			if fillMappingNode(valueNode, targetValue) {
				changed = true
			}
		case yaml.SequenceNode:
			if isBlankSequence(targetValue) {
				target.Content[idx+1] = deepCopyNode(valueNode)
				changed = true
			}
		case yaml.ScalarNode:
			if isBlankScalar(targetValue) {
				targetValue.Value = valueNode.Value
				targetValue.Tag = valueNode.Tag
				targetValue.Style = valueNode.Style
				changed = true
			}
		default:
			if isBlankNode(targetValue) {
				target.Content[idx+1] = deepCopyNode(valueNode)
				changed = true
			}
		}
	}

	return changed
}

func findMappingKey(node *yaml.Node, key string) int {
	if node == nil || node.Kind != yaml.MappingNode {
		return -1
	}
	for i := 0; i+1 < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return i
		}
	}
	return -1
}

func isBlankNode(node *yaml.Node) bool {
	if node == nil {
		return true
	}
	switch node.Kind {
	case yaml.ScalarNode:
		return isBlankScalar(node)
	case yaml.MappingNode:
		return len(node.Content) == 0
	case yaml.SequenceNode:
		return isBlankSequence(node)
	default:
		return false
	}
}

func isBlankScalar(node *yaml.Node) bool {
	if node == nil {
		return true
	}
	return strings.TrimSpace(node.Value) == ""
}

func isBlankSequence(node *yaml.Node) bool {
	if node == nil {
		return true
	}
	return len(node.Content) == 0
}

func deepCopyNode(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}
	copied := *node
	copied.Content = make([]*yaml.Node, len(node.Content))
	for i, child := range node.Content {
		copied.Content[i] = deepCopyNode(child)
	}
	return &copied
}
