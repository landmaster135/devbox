package usecase

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileAccessor はYAMLファイルの読み書きを抽象化します。
type FileAccessor interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

// osFileAccessor はローカルファイルシステムを利用します。
type osFileAccessor struct{}

func (osFileAccessor) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (osFileAccessor) WriteFile(path string, data []byte) error {
	perm := os.FileMode(0o644)
	if info, err := os.Stat(path); err == nil {
		perm = info.Mode()
	}
	return os.WriteFile(path, data, perm)
}

// KeyValuePair は編集対象のキーと値を表します。
type KeyValuePair struct {
	Key   string
	Value interface{}
}

// Service はYAMLの解析・編集機能を提供します。
type Service struct {
	fileAccessor FileAccessor
}

// NewService は標準のFileAccessorを使ってサービスを生成します。
func NewService() *Service {
	return NewServiceWithFileAccessor(osFileAccessor{})
}

// NewServiceWithFileAccessor はカスタムFileAccessorを注入できます（主にテスト用）。
func NewServiceWithFileAccessor(accessor FileAccessor) *Service {
	if accessor == nil {
		accessor = osFileAccessor{}
	}
	return &Service{fileAccessor: accessor}
}

// ReadFromFile は指定されたファイルからYAMLを読み取り、オブジェクトとして返します。
func (s *Service) ReadFromFile(path string) (any, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("YAMLファイルのパスが指定されていません")
	}

	data, err := s.fileAccessor.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("YAMLファイルの読み込みに失敗しました: %w", err)
	}

	return parseYAML(data)
}

// ParseFromContent は与えられた文字列（YAML）を解析してオブジェクトに変換します。
func (s *Service) ParseFromContent(content string) (any, error) {
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("YAMLコンテンツが空です")
	}

	return parseYAML([]byte(content))
}

// EditFile は指定されたキーと値でYAMLファイルを更新し、更新後の構造を返します。
func (s *Service) EditFile(path string, keyValueList string) (any, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("YAMLファイルのパスが指定されていません")
	}
	if strings.TrimSpace(keyValueList) == "" {
		return nil, errors.New("更新するキーと値を指定してください")
	}

	keyValues, err := parseKeyValueList(keyValueList)
	if err != nil {
		return nil, err
	}

	data, err := s.fileAccessor.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("YAMLファイルの読み込みに失敗しました: %w", err)
	}

	updated, err := applyEditsWithOriginalOrder(data, keyValues)
	if err != nil {
		return nil, err
	}

	if err := s.fileAccessor.WriteFile(path, updated); err != nil {
		return nil, fmt.Errorf("YAMLファイルの書き込みに失敗しました: %w", err)
	}

	result, err := parseYAML(updated)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseYAML(data []byte) (any, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return nil, errors.New("YAMLの内容が空です")
	}

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	var documents []any
	for {
		var doc any
		if err := decoder.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("YAMLの解析に失敗しました: %w", err)
		}
		if doc == nil {
			continue
		}
		documents = append(documents, normalize(doc))
	}

	if len(documents) == 0 {
		return nil, errors.New("YAMLの構造を特定できませんでした")
	}

	if len(documents) == 1 {
		return documents[0], nil
	}

	return documents, nil
}

func normalize(value any) any {
	switch v := value.(type) {
	case map[string]any:
		for key, val := range v {
			v[key] = normalize(val)
		}
		return v
	case map[any]any:
		normalized := make(map[string]any, len(v))
		for key, val := range v {
			normalized[fmt.Sprint(key)] = normalize(val)
		}
		return normalized
	case []any:
		for i, val := range v {
			v[i] = normalize(val)
		}
		return v
	case *yaml.Node:
		var decoded any
		if err := v.Decode(&decoded); err != nil {
			return v.Value
		}
		return normalize(decoded)
	default:
		return v
	}
}
