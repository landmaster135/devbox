package yamlfmt

import (
	"fmt"

	"github.com/landmaster135/devbox/internal/data_converter/usecases/formats/common"
	"gopkg.in/yaml.v3"
)

// Parse は YAML を key-value リストへ正規化します。
func Parse(content []byte) (*common.NormalizedData, error) {
	return common.ParseStructuredToKeyValueList(content, yaml.Unmarshal, "YAML")
}

// Serialize は key-value リストを YAML へ変換します。
func Serialize(records []map[string]string) ([]byte, error) {
	b, err := yaml.Marshal(records)
	if err != nil {
		return nil, fmt.Errorf("YAML生成に失敗しました: %w", err)
	}
	return b, nil
}
