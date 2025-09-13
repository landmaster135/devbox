package domain

// 色定数
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

// GetSuspiciousPatterns は疑わしいキー名のパターンを返す
func GetSuspiciousPatterns() []string {
	return []string{
		`(?i)api[_-]?key`,
		`(?i)secret[_-]?key`,
		`(?i)access[_-]?token`,
		`(?i)private[_-]?key`,
		`(?i)client[_-]?secret`,
		`(?i)auth[_-]?token`,
		`(?i)bearer[_-]?token`,
		`(?i)integration[_-]?token`,
		`(?i)webhook[_-]?url`,
		`(?i)credentials[_-]?path`,
		`(?i)database[_-]?url`,
		`(?i)connection[_-]?string`,
		`(?i)password`,
		`(?i)secret`,
		`(?i)token`,
		`(?i)key`,
	}
}

// GetAllowedPlaceholders は許可されたプレースホルダー値を返す
func GetAllowedPlaceholders() []string {
	return []string{
		"YOUR_API_KEY",
		"YOUR_BRAVE_API_KEY",
		"YOUR_GITHUB_PAT",
		"YOUR_URL",
		"YOUR_KEY",
		"YOUR_VALUE",
		"YOUR_PATH",
		"YOUR_PERSONAL_ACCESS_TOKEN",
		"REPLACE_WITH_YOUR_KEY",
		"REPLACE_WITH_YOUR_TOKEN",
		"CHANGE_ME",
		"PLACEHOLDER",
		"EXAMPLE_KEY",
		"DUMMY_VALUE",
		"",
	}
}

// GetRealSecretPatterns は実際のAPIキーパターンを返す
func GetRealSecretPatterns() []string {
	return []string{
		`sk-[a-zA-Z0-9]{48}`, // OpenAI API key
		`xoxp-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-zA-Z0-9]{32}`,           // Slack token
		`ghp_[a-zA-Z0-9]{36}`,                                          // GitHub personal access token
		`AIza[0-9A-Za-z_-]{35}`,                                        // Google API key
		`AKIA[0-9A-Z]{16}`,                                             // AWS access key
		`ya29\.[a-zA-Z0-9_-]+`,                                         // Google OAuth access token
		`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`, // UUID format
	}
}

// GetConfigFilePatterns は設定ファイルのパターンを返す
func GetConfigFilePatterns() []string {
	return []string{
		"*.json",
		"*.config.js",
		"*.config.ts",
		"mcp_settings.json",
		"mcp_settings*.json",
		"claude_desktop_config.json",
		"cline_mcp_settings.json",
	}
}

// GetTestPatterns はテスト用の値のパターンを返す
func GetTestPatterns() []string {
	return []string{
		`(?i)(test|demo|example|dummy|fake|mock|sample)`,
	}
}

// GetProtocolPrefixes はプロトコル識別子のリストを返す
func GetProtocolPrefixes() []string {
	return []string{
		"http://",
		"https://",
		"postgresql://",
		"postgres://",
		"mysql://",
		"mongodb://",
		"sqlite://",
		"redis://",
		"ftp://",
		"sftp://",
		"file://",
		"ldap://",
		"ldaps://",
		"smtp://",
		"smtps://",
		"pop3://",
		"imap://",
		"ssh://",
		"tcp://",
		"udp://",
		"ws://",
		"wss://",
		"amqp://",
		"amqps://",
		"kafka://",
		"elasticsearch://",
		"memcached://",
		"cassandra://",
		"couchdb://",
		"neo4j://",
		"influxdb://",
	}
}

// GetHomePathPattern はホームパス検知用のパターンを返す
func GetHomePathPattern() string {
	return `/home` + `/`
}

// GetAllowedHomePathPatterns は許可されるホームパスのパターンリストを返す
func GetAllowedHomePathPatterns() []string {
	return []string{
		`/home/user`,
		`/home/[username]`,
		`/home/alice`,
	}
}

// GetBinaryFileExtensions はバイナリファイルの拡張子リストを返す
func GetBinaryFileExtensions() []string {
	return []string{
		".exe", ".dll", ".so", ".dylib", ".a", ".o", ".obj",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp",
		".mp3", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".mkv",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".zip", ".tar", ".gz", ".bz2", ".7z", ".rar",
		".bin", ".dat", ".db", ".sqlite", ".sqlite3",
	}
}
