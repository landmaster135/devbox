package entity

// YAMLData はパースされたYAMLデータを表す構造体です
type YAMLData struct {
	Data interface{}
}

// NewYAMLData は新しいYAMLDataエンティティを作成します
func NewYAMLData(data interface{}) *YAMLData {
	return &YAMLData{
		Data: data,
	}
}

// GetData はパースされたデータを返します
func (y *YAMLData) GetData() interface{} {
	return y.Data
}
