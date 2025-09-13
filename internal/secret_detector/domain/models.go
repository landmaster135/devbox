package domain

// Config はMCP設定ファイルの構造体
type Config struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// ServerConfig はサーバー設定の構造体
type ServerConfig struct {
	Command     string            `json:"command"`
	Args        []string          `json:"args"`
	Env         map[string]string `json:"env,omitempty"`
	Disabled    bool              `json:"disabled"`
	AutoApprove []string          `json:"autoApprove"`
}

// SecretResult はシークレット検知結果の構造体
type SecretResult struct {
	File           string
	Server         string
	Key            string
	Value          string
	IsPlaceholder  bool
	MatchedPattern string
}

// HomePathResult はホームパス検知結果の構造体
type HomePathResult struct {
	File         string
	LineNumber   int
	Content      string
	IsAllowed    bool
	MatchedPath  string
}

// ScanSummary はスキャン結果のサマリー構造体
type ScanSummary struct {
	TotalFiles         int
	TotalEnvVars       int
	SecretCount        int
	PlaceholderCount   int
	HomePathCount      int
	HasRealSecrets     bool
	HasForbiddenPaths  bool
}
