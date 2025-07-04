package yaml_service

import (
	yaml_controller "github.com/landmaster135/devbox/internal/yaml_parser/adapter/controller"
	yaml_gateway "github.com/landmaster135/devbox/internal/yaml_parser/adapter/gateway"
	yaml_infrastructure "github.com/landmaster135/devbox/internal/yaml_parser/infrastructure"
	yaml_usecase "github.com/landmaster135/devbox/internal/yaml_parser/usecase"
)

type YamlService struct {
	YAMLController yaml_controller.YAMLController
}

func NewYamlService() *YamlService {
	// インフラストラクチャ層の初期化
	yamlParser := yaml_infrastructure.NewYAMLParser()

	// アダプター層の初期化
	yamlGateway := yaml_gateway.NewYAMLGateway(yamlParser)

	// ユースケース層の初期化
	yamlUseCase := yaml_usecase.NewYAMLUseCase(yamlGateway)

	// コントローラー層の初期化
	yamlController := yaml_controller.NewYAMLController(yamlUseCase)

	return &YamlService{
		YAMLController: *yamlController,
	}
}

func (s *YamlService) ParseYAML(yamlContent string) (interface{}, error) {
	parsedData, err := s.YAMLController.ParseYAML(yamlContent)
	if err != nil {
		return nil, err
	}
	return parsedData, nil
}

func (s *YamlService) ParseYAMLToStruct(yamlContent string, model interface{}) error {
	// モデルは既にポインタとして渡されるため、再度ポインタ化しない
	err := s.YAMLController.ParseYAMLToStruct(yamlContent, model)
	if err != nil {
		return err
	}
	return nil
}

func (s *YamlService) MarshalToYAML(content interface{}) (string, error) {
	yamlString, err := s.YAMLController.MarshalToYAML(content)
	if err != nil {
		return "", err
	}
	return yamlString, nil
}

func (s *YamlService) ReadYAMLFile(filePath string) (string, error) {
	yamlContent, err := s.YAMLController.ReadYAMLFile(filePath)
	if err != nil {
		return "", err
	}
	return yamlContent, nil
}

func (s *YamlService) ParseYAMLFile(filePath string) (interface{}, error) {
	parsedData, err := s.YAMLController.ParseYAMLFile(filePath)
	if err != nil {
		return nil, err
	}
	return parsedData, nil
}

func (s *YamlService) ParseYAMLFileToStruct(filePath string, model interface{}) error {
	err := s.YAMLController.ParseYAMLFileToStruct(filePath, model)
	if err != nil {
		return err
	}
	return nil
}
