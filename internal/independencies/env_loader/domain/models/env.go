package models

// EnvConfig は環境変数の設定を表す構造体です
type EnvConfig struct {
	// 環境変数のマップ（キーは環境変数名、値は環境変数の値）
	Variables map[string]string
}

// NewEnvConfig は新しいEnvConfigインスタンスを作成します
func NewEnvConfig() *EnvConfig {
	return &EnvConfig{
		Variables: make(map[string]string),
	}
}

// AddVariable は環境変数を追加します
func (ec *EnvConfig) AddVariable(key, value string) {
	ec.Variables[key] = value
}

// GetVariable は環境変数の値を取得します
func (ec *EnvConfig) GetVariable(key string) (string, bool) {
	value, exists := ec.Variables[key]
	return value, exists
}

// GetAllVariables は全ての環境変数を取得します
func (ec *EnvConfig) GetAllVariables() map[string]string {
	return ec.Variables
}
