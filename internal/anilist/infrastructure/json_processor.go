package infrastructure

import (
	"encoding/json"
)

// JSONProcessor はJSON処理のインターフェース
type JSONProcessor interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

// JSONProcessorImpl は実際のJSON処理を行う実装
type JSONProcessorImpl struct{}

// NewJSONProcessor は新しいJSONProcessorImplを作成する
func NewJSONProcessor() *JSONProcessorImpl {
	return &JSONProcessorImpl{}
}

// MarshalIndent はJSONにマーシャルする
func (jp *JSONProcessorImpl) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
