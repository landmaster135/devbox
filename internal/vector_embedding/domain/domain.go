package domain

// EmbedRequest はユースケース層が扱う埋め込み要求を表す。
type EmbedRequest struct {
	Operation string
	Model     string
	Inputs    []string
}

// EmbedResult はCLIに返却する結果を表す。
type EmbedResult struct {
	Provider   string      `json:"provider"`
	Model      string      `json:"model"`
	Embeddings [][]float64 `json:"embeddings"`
	InputCount int         `json:"input_count"`
	Dimensions int         `json:"dimensions"`
}
