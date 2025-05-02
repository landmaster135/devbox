package infrastructure

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/adapter/gateway"
)

// YAMLParserInterfaceを実装していることを確認
var _ gateway.YAMLParserInterface = (*YAMLParser)(nil)

// YAMLParser はYAMLデータのパースと変換を行います
type YAMLParser struct{}

// NewYAMLParser は新しいYAMLParserを作成します
func NewYAMLParser() *YAMLParser {
	return &YAMLParser{}
}

// ParseYAML はYAML文字列をパースします
func (p *YAMLParser) ParseYAML(yamlContent string) (interface{}, error) {
	var result interface{}
	err := yaml.Unmarshal([]byte(yamlContent), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// MarshalToYAML はデータをYAML文字列に変換します
func (p *YAMLParser) MarshalToYAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ParseYAMLToStruct はYAML文字列を指定された構造体にパースします
func (p *YAMLParser) ParseYAMLToStruct(yamlContent string, out interface{}) error {
	return yaml.Unmarshal([]byte(yamlContent), out)
}

// ReadYAMLFile はファイルからYAMLデータを読み込みます
func (p *YAMLParser) ReadYAMLFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// ParseYAMLFile はファイルからYAMLデータを読み込み、パースします
func (p *YAMLParser) ParseYAMLFile(filePath string) (interface{}, error) {
	content, err := p.ReadYAMLFile(filePath)
	if err != nil {
		return nil, err
	}

	return p.ParseYAML(content)
}

// ParseYAMLFileToStruct はファイルからYAMLデータを読み込み、指定された構造体にパースします
func (p *YAMLParser) ParseYAMLFileToStruct(filePath string, out interface{}) error {
	content, err := p.ReadYAMLFile(filePath)
	if err != nil {
		return err
	}

	return p.ParseYAMLToStruct(content, out)
}
