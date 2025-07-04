package repository

import "github.com/landmaster135/devbox/internal/yaml_parser/domain/entity"

// YAMLRepository はYAMLデータの操作に関するインターフェースです
type YAMLRepository interface {
	Parse(yamlContent string) (*entity.YAMLData, error)
	Marshal(data interface{}) (string, error)
	ParseToStruct(yamlContent string, out interface{}) error
	ReadFile(filePath string) (string, error)
	ParseFile(filePath string) (*entity.YAMLData, error)
	ParseFileToStruct(filePath string, out interface{}) error
}
