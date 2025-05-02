package gateway

import (
	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/domain/entity"
)

// YAMLParserInterface はYAMLパーサーのインターフェースです
type YAMLParserInterface interface {
	ParseYAML(yamlContent string) (interface{}, error)
	MarshalToYAML(data interface{}) (string, error)
	ParseYAMLToStruct(yamlContent string, out interface{}) error
	ReadYAMLFile(filePath string) (string, error)
	ParseYAMLFile(filePath string) (interface{}, error)
	ParseYAMLFileToStruct(filePath string, out interface{}) error
}

// YAMLGateway はリポジトリの実装です
type YAMLGateway struct {
	parser YAMLParserInterface
}

// NewYAMLGateway は新しいYAMLGatewayを作成します
func NewYAMLGateway(parser YAMLParserInterface) *YAMLGateway {
	return &YAMLGateway{
		parser: parser,
	}
}

// Parse はYAML文字列をパースします
func (g *YAMLGateway) Parse(yamlContent string) (*entity.YAMLData, error) {
	parsedData, err := g.parser.ParseYAML(yamlContent)
	if err != nil {
		return nil, err
	}
	return entity.NewYAMLData(parsedData), nil
}

// Marshal はデータをYAML文字列に変換します
func (g *YAMLGateway) Marshal(data interface{}) (string, error) {
	return g.parser.MarshalToYAML(data)
}

// ParseToStruct はYAML文字列を指定された構造体にパースします
func (g *YAMLGateway) ParseToStruct(yamlContent string, out interface{}) error {
	return g.parser.ParseYAMLToStruct(yamlContent, out)
}

// ReadFile はファイルからYAMLデータを読み込みます
func (g *YAMLGateway) ReadFile(filePath string) (string, error) {
	return g.parser.ReadYAMLFile(filePath)
}

// ParseFile はファイルからYAMLデータを読み込み、パースします
func (g *YAMLGateway) ParseFile(filePath string) (*entity.YAMLData, error) {
	parsedData, err := g.parser.ParseYAMLFile(filePath)
	if err != nil {
		return nil, err
	}
	return entity.NewYAMLData(parsedData), nil
}

// ParseFileToStruct はファイルからYAMLデータを読み込み、指定された構造体にパースします
func (g *YAMLGateway) ParseFileToStruct(filePath string, out interface{}) error {
	return g.parser.ParseYAMLFileToStruct(filePath, out)
}
