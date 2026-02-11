package domain

// VersionResponse は /api/version のレスポンスを表す。
type VersionResponse struct {
	Version string `json:"version"`
}

// TagsResponse は /api/tags のレスポンスを表す。
type TagsResponse struct {
	Models []TagModel `json:"models"`
}

// TagModel はインストール済みモデルのメタデータを示す。
type TagModel struct {
	Name     string `json:"name"`
	Modified string `json:"modified"`
	Size     int64  `json:"size"`
	Digest   string `json:"digest"`
}

// ProcessesResponse は /api/ps のレスポンスを表す。
type ProcessesResponse struct {
	Models []ProcessModel `json:"models"`
}

// ProcessModel は稼働中モデルのメタデータを示す。
type ProcessModel struct {
	Name      string `json:"name"`
	Model     string `json:"model"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	ExpiresAt string `json:"expires_at"`
}

// EmbedRequest は /api/embed のリクエスト。
type EmbedRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// EmbedResponse は /api/embed のレスポンス。
type EmbedResponse struct {
	Model      string      `json:"model"`
	Embeddings [][]float64 `json:"embeddings"`
}

// GenerateRequest は /api/generate のリクエスト。
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateChunk はストリーミングレスポンスの1行分を表す。
type GenerateChunk struct {
	Model      string `json:"model"`
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	DoneReason string `json:"done_reason"`
	Error      string `json:"error"`
}

// PullRequest は /api/pull のリクエスト。
type PullRequest struct {
	Model  string `json:"model"`
	Stream bool   `json:"stream"`
}

// PullChunk は /api/pull のストリーミングレスポンス 1 行分。
type PullChunk struct {
	Status    string `json:"status"`
	Digest    string `json:"digest"`
	Total     int64  `json:"total"`
	Completed int64  `json:"completed"`
	Error     string `json:"error"`
}
