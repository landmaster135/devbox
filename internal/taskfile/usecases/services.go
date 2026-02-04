package usecases

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	// taskTypeRoot は root Taskfile を示す内部識別子
	taskTypeRoot = "root"
)

var rootReferenceCandidates = []string{
	filepath.Join("internal", "taskfile", "usecases", "taskfiles", "root.yml"),
	filepath.Join("taskfiles", "root.yml"),
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

// Fill writes reference Taskfile values into blank or missing target fields
func (s *Service) Fill(taskType, targetPath string) (bool, error) {
	referencePath, err := resolveReferencePath(taskType)
	if err != nil {
		return false, err
	}

	referenceDoc, err := readYAMLDocument(referencePath)
	if err != nil {
		return false, fmt.Errorf("参照Taskfileの読み込みに失敗しました: %w", err)
	}

	targetDoc, err := readYAMLDocument(targetPath)
	if err != nil {
		return false, fmt.Errorf("補完対象Taskfileの読み込みに失敗しました: %w", err)
	}

	referenceRoot := ensureDocumentRoot(referenceDoc)
	targetRoot := ensureDocumentRoot(targetDoc)
	if referenceRoot == nil || targetRoot == nil {
		return false, fmt.Errorf("Taskfileの解析に失敗しました")
	}

	changed := fillMappingNode(referenceRoot, targetRoot)
	if !changed {
		return false, nil
	}

	if err := writeYAMLDocument(targetPath, targetRoot); err != nil {
		return false, fmt.Errorf("Taskfileの書き込みに失敗しました: %w", err)
	}

	return true, nil
}

// Create writes a brand-new Taskfile from the reference template
func (s *Service) Create(taskType, targetPath string) error {
	referencePath, err := resolveReferencePath(taskType)
	if err != nil {
		return err
	}

	templateData, err := os.ReadFile(referencePath)
	if err != nil {
		return fmt.Errorf("参照Taskfileの読み込みに失敗しました: %w", err)
	}

	dir := filepath.Dir(targetPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("Taskfileのディレクトリ作成に失敗しました: %w", err)
		}
	}

	if info, err := os.Stat(targetPath); err == nil {
		if info.IsDir() {
			return fmt.Errorf("指定されたパスはディレクトリです: %s", targetPath)
		}
		return fmt.Errorf("指定されたパスには既にファイルが存在します: %s", targetPath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Taskfileの存在確認に失敗しました: %w", err)
	}

	if err := os.WriteFile(targetPath, templateData, 0o644); err != nil {
		return fmt.Errorf("Taskfileの書き込みに失敗しました: %w", err)
	}

	return nil
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

func diffFields(reference, target map[string]struct{}) []string {
	missing := make([]string, 0)
	for path := range reference {
		if _, ok := target[path]; !ok {
			missing = append(missing, path)
		}
	}
	return missing
}
