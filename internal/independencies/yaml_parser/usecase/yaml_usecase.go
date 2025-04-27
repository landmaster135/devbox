package usecase

import (
	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/domain/entity"
	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/domain/repository"
)

// YAMLUseCase はYAMLデータの処理に関するユースケースを定義します
type YAMLUseCase struct {
	yamlRepo repository.YAMLRepository
}

// NewYAMLUseCase は新しいYAMLUseCaseを作成します
func NewYAMLUseCase(yamlRepo repository.YAMLRepository) *YAMLUseCase {
	return &YAMLUseCase{
		yamlRepo: yamlRepo,
	}
}

// ParseYAML はYAML文字列をパースします
func (uc *YAMLUseCase) ParseYAML(yamlContent string) (*entity.YAMLData, error) {
	return uc.yamlRepo.Parse(yamlContent)
}

// MarshalToYAML はデータをYAML文字列に変換します
func (uc *YAMLUseCase) MarshalToYAML(data interface{}) (string, error) {
	return uc.yamlRepo.Marshal(data)
}

// ParseYAMLToStruct はYAML文字列を指定された構造体にパースします
func (uc *YAMLUseCase) ParseYAMLToStruct(yamlContent string, out interface{}) error {
	return uc.yamlRepo.ParseToStruct(yamlContent, out)
}

// ReadYAMLFile はファイルからYAMLデータを読み込みます
func (uc *YAMLUseCase) ReadYAMLFile(filePath string) (string, error) {
	return uc.yamlRepo.ReadFile(filePath)
}

// ParseYAMLFile はファイルからYAMLデータを読み込み、パースします
func (uc *YAMLUseCase) ParseYAMLFile(filePath string) (*entity.YAMLData, error) {
	return uc.yamlRepo.ParseFile(filePath)
}

// ParseYAMLFileToStruct はファイルからYAMLデータを読み込み、指定された構造体にパースします
func (uc *YAMLUseCase) ParseYAMLFileToStruct(filePath string, out interface{}) error {
	return uc.yamlRepo.ParseFileToStruct(filePath, out)
}
