package controller

import (
	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/domain/entity"
	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/usecase"
)

// YAMLUseCaseInterface はYAMLUseCaseのインターフェースを定義します
type YAMLUseCaseInterface interface {
	ParseYAML(yamlContent string) (*entity.YAMLData, error)
	MarshalToYAML(data interface{}) (string, error)
	ParseYAMLToStruct(yamlContent string, out interface{}) error
	ReadYAMLFile(filePath string) (string, error)
	ParseYAMLFile(filePath string) (*entity.YAMLData, error)
	ParseYAMLFileToStruct(filePath string, out interface{}) error
}

// YAMLController はYAMLデータの入出力を制御します
type YAMLController struct {
	YAMLUseCase YAMLUseCaseInterface
}

// NewYAMLController は新しいYAMLControllerを作成します
func NewYAMLController(yamlUseCase *usecase.YAMLUseCase) *YAMLController {
	return &YAMLController{
		YAMLUseCase: yamlUseCase,
	}
}

// ParseYAML はYAML文字列をパースして結果を返します
func (c *YAMLController) ParseYAML(yamlContent string) (interface{}, error) {
	yamlData, err := c.YAMLUseCase.ParseYAML(yamlContent)
	if err != nil {
		return nil, err
	}
	return yamlData.GetData(), nil
}

// MarshalToYAML はデータをYAML文字列に変換します
func (c *YAMLController) MarshalToYAML(data interface{}) (string, error) {
	return c.YAMLUseCase.MarshalToYAML(data)
}

// ParseYAMLToStruct はYAML文字列を指定された構造体にパースします
func (c *YAMLController) ParseYAMLToStruct(yamlContent string, out interface{}) error {
	return c.YAMLUseCase.ParseYAMLToStruct(yamlContent, out)
}

// ReadYAMLFile はファイルからYAMLデータを読み込みます
func (c *YAMLController) ReadYAMLFile(filePath string) (string, error) {
	return c.YAMLUseCase.ReadYAMLFile(filePath)
}

// ParseYAMLFile はファイルからYAMLデータを読み込み、パースします
func (c *YAMLController) ParseYAMLFile(filePath string) (interface{}, error) {
	yamlData, err := c.YAMLUseCase.ParseYAMLFile(filePath)
	if err != nil {
		return nil, err
	}
	return yamlData.GetData(), nil
}

// ParseYAMLFileToStruct はファイルからYAMLデータを読み込み、指定された構造体にパースします
func (c *YAMLController) ParseYAMLFileToStruct(filePath string, out interface{}) error {
	return c.YAMLUseCase.ParseYAMLFileToStruct(filePath, out)
}
