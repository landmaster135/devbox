package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockFlagParser はテスト用のFlagParser実装
type MockFlagParser struct {
	args   []string
	values map[string]interface{}
}

// NewMockFlagParser は新しいMockFlagParserを作成する
func NewMockFlagParser(args []string) *MockFlagParser {
	return &MockFlagParser{
		args:   args,
		values: make(map[string]interface{}),
	}
}

// StringVar は文字列フラグを定義する（モック）
func (p *MockFlagParser) StringVar(ptr *string, name string, value string, usage string) {
	if val, exists := p.values[name]; exists {
		*ptr = val.(string)
	} else {
		*ptr = value
	}
}

// BoolVar はブールフラグを定義する（モック）
func (p *MockFlagParser) BoolVar(ptr *bool, name string, value bool, usage string) {
	if val, exists := p.values[name]; exists {
		*ptr = val.(bool)
	} else {
		*ptr = value
	}
}

// Parse はコマンドライン引数を解析する（モック）
func (p *MockFlagParser) Parse() error {
	return nil
}

// Args は解析後の残りの引数を返す（モック）
func (p *MockFlagParser) Args() []string {
	return p.args
}

// SetValue はモック用の値を設定する
func (p *MockFlagParser) SetValue(name string, value interface{}) {
	p.values[name] = value
}

// TestNewConfig は新しいConfigの作成をテストする
func TestNewConfig_Normal(t *testing.T) {
	serverURL := "https://example.com/mcp/sse"
	cfg, err := NewConfig(serverURL)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, serverURL, cfg.ServerURL)
	assert.Equal(t, 0, cfg.CallbackPort)
	assert.Equal(t, TransportHTTPFirst, cfg.TransportStrategy)
	assert.Equal(t, "localhost", cfg.Host)
	assert.False(t, cfg.Debug)
	assert.False(t, cfg.AllowHTTP)
	assert.NotNil(t, cfg.Headers)
}

// TestNewConfig_EmptyURL は空のURLでのConfig作成をテストする
func TestNewConfig_EmptyURL(t *testing.T) {
	cfg, err := NewConfig("")

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "サーバーURLが指定されていません")
}

// TestParseFlagsWithParser_Normal は正常なフラグ解析をテストする
func TestParseFlagsWithParser_Normal(t *testing.T) {
	parser := NewMockFlagParser([]string{"https://example.com/mcp/sse"})
	parser.SetValue("debug", true)
	parser.SetValue("transport", "sse-only")
	parser.SetValue("callback-port", "3334")

	cfg, err := ParseFlagsWithParser(parser)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "https://example.com/mcp/sse", cfg.ServerURL)
	assert.Equal(t, 3334, cfg.CallbackPort)
	assert.Equal(t, TransportSSEOnly, cfg.TransportStrategy)
	assert.True(t, cfg.Debug)
}

// TestParseFlagsWithParser_Help はヘルプフラグのテストする
func TestParseFlagsWithParser_Help(t *testing.T) {
	parser := NewMockFlagParser([]string{})
	parser.SetValue("help", true)

	cfg, err := ParseFlagsWithParser(parser)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.True(t, cfg.Help)
}

// TestParseFlagsWithParser_NoServerURL はサーバーURL未指定のテストする
func TestParseFlagsWithParser_NoServerURL(t *testing.T) {
	parser := NewMockFlagParser([]string{})

	cfg, err := ParseFlagsWithParser(parser)

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "サーバーURLが指定されていません")
}

// TestParseFlagsWithParser_InvalidPort は無効なポート番号のテストする
func TestParseFlagsWithParser_InvalidPort(t *testing.T) {
	parser := NewMockFlagParser([]string{"https://example.com/mcp/sse"})
	parser.SetValue("callback-port", "invalid")

	cfg, err := ParseFlagsWithParser(parser)

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "無効なポート番号です")
}

// TestParseFlagsWithParser_InvalidTransportStrategy は無効なトランスポート戦略のテストする
func TestParseFlagsWithParser_InvalidTransportStrategy(t *testing.T) {
	parser := NewMockFlagParser([]string{"https://example.com/mcp/sse"})
	parser.SetValue("transport", "invalid-strategy")

	cfg, err := ParseFlagsWithParser(parser)

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "無効なトランスポート戦略です")
}

// TestParseFlagsWithParser_Headers はヘッダー解析のテストする
func TestParseFlagsWithParser_Headers(t *testing.T) {
	parser := NewMockFlagParser([]string{"https://example.com/mcp/sse"})
	parser.SetValue("header", "Authorization:Bearer token123,Content-Type:application/json")

	cfg, err := ParseFlagsWithParser(parser)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "Bearer token123", cfg.Headers["Authorization"])
	assert.Equal(t, "application/json", cfg.Headers["Content-Type"])
}

// TestParseFlagsWithParser_InvalidHeaders は無効なヘッダー形式のテストする
func TestParseFlagsWithParser_InvalidHeaders(t *testing.T) {
	parser := NewMockFlagParser([]string{"https://example.com/mcp/sse"})
	parser.SetValue("header", "InvalidHeader")

	cfg, err := ParseFlagsWithParser(parser)

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "無効なヘッダー形式です")
}

// TestIsValidTransportStrategy はトランスポート戦略の検証をテストする
func TestIsValidTransportStrategy_Valid(t *testing.T) {
	validStrategies := []TransportStrategy{
		TransportSSEOnly,
		TransportHTTPOnly,
		TransportSSEFirst,
		TransportHTTPFirst,
	}

	for _, strategy := range validStrategies {
		assert.True(t, isValidTransportStrategy(strategy), "Strategy %s should be valid", strategy)
	}
}

// TestIsValidTransportStrategy_Invalid は無効なトランスポート戦略のテストする
func TestIsValidTransportStrategy_Invalid(t *testing.T) {
	invalidStrategy := TransportStrategy("invalid-strategy")
	assert.False(t, isValidTransportStrategy(invalidStrategy))
}
