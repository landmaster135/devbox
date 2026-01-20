package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// taskTypeRoot は root Taskfile を示す内部識別子
	taskTypeRoot = "root"
)

var rootReferenceCandidates = []string{
	filepath.Join("internal", "taskfile", "usecases", "taskfiles", "root.yml"),
}

// InspectionResult represents diff outcome between reference and target Taskfile
type InspectionResult struct {
	MissingFields []string
}

// HasMissingFields returns true when there are missing paths
func (r *InspectionResult) HasMissingFields() bool {
	return len(r.MissingFields) > 0
}

// Service controls Taskfile validation logic
type Service struct{}

// NewService builds Service
func NewService() *Service {
	return &Service{}
}

// Inspect compares the reference Taskfile of given type with the target file
func (s *Service) Inspect(taskType, targetPath string) (*InspectionResult, error) {
	referencePath, err := resolveReferencePath(taskType)
	if err != nil {
		return nil, err
	}

	referenceFields, err := readFieldSet(referencePath)
	if err != nil {
		return nil, fmt.Errorf("参照Taskfileの読み込みに失敗しました: %w", err)
	}

	targetFields, err := readFieldSet(targetPath)
	if err != nil {
		return nil, fmt.Errorf("検証対象Taskfileの読み込みに失敗しました: %w", err)
	}

	missing := diffFields(referenceFields, targetFields)
	sort.Strings(missing)

	return &InspectionResult{MissingFields: missing}, nil
}

func resolveReferencePath(taskType string) (string, error) {
	switch taskType {
	case taskTypeRoot:
		return findExistingPath(rootReferenceCandidates)
	default:
		return "", fmt.Errorf("未サポートのtask-typeです: %s", taskType)
	}
}

func findExistingPath(candidates []string) (string, error) {
	tried := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		tried = append(tried, candidate)
		path := candidate
		if !filepath.IsAbs(path) {
			absPath, err := filepath.Abs(candidate)
			if err == nil {
				path = absPath
			}
		}
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path, nil
		}
	}
	return "", fmt.Errorf("参照Taskfileが見つかりません: %s", strings.Join(tried, ", "))
}

func readFieldSet(path string) (map[string]struct{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return extractFieldPaths(data)
}

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

func diffFields(reference, target map[string]struct{}) []string {
	missing := make([]string, 0)
	for path := range reference {
		if _, ok := target[path]; !ok {
			missing = append(missing, path)
		}
	}
	return missing
}
