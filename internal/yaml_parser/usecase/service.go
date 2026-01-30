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

// FileReader は任意のソースからファイルを読み込むためのインターフェースです。
type FileReader interface {
	ReadFile(path string) ([]byte, error)
}

// osFileReader はローカルファイルシステムからデータを読み込みます。
type osFileReader struct{}

func (osFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// Service はYAMLの解析機能を提供します。
type Service struct {
	fileReader FileReader
}

// NewService は標準のFileReaderを使ってサービスを生成します。
func NewService() *Service {
	return NewServiceWithFileReader(osFileReader{})
}

// NewServiceWithFileReader はカスタムFileReaderを注入できます（主にテスト用）。
func NewServiceWithFileReader(reader FileReader) *Service {
	if reader == nil {
		reader = osFileReader{}
	}
	return &Service{fileReader: reader}
}

// ReadFromFile は指定されたファイルからYAMLを読み取り、オブジェクトとして返します。
func (s *Service) ReadFromFile(path string) (any, error) {
	if strings.TrimSpace(path) == "" {
		return nil, errors.New("YAMLファイルのパスが指定されていません")
	}

	data, err := s.fileReader.ReadFile(path)
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
