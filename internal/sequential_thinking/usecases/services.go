package usecases

import (
	"encoding/json"
	"fmt"
	"os"
)

// ThoughtData は思考データを表す構造体
type ThoughtData struct {
	Thought           string `json:"thought"`
	ThoughtNumber     int    `json:"thoughtNumber"`
	TotalThoughts     int    `json:"totalThoughts"`
	NextThoughtNeeded bool   `json:"nextThoughtNeeded"`
	IsRevision        bool   `json:"isRevision,omitempty"`
	RevisesThought    int    `json:"revisesThought,omitempty"`
	BranchFromThought int    `json:"branchFromThought,omitempty"`
	BranchID          string `json:"branchId,omitempty"`
	NeedsMoreThoughts bool   `json:"needsMoreThoughts,omitempty"`
}

// ProcessResult は思考処理の結果を表す構造体
type ProcessResult struct {
	ThoughtNumber        int      `json:"thoughtNumber"`
	TotalThoughts        int      `json:"totalThoughts"`
	NextThoughtNeeded    bool     `json:"nextThoughtNeeded"`
	Branches             []string `json:"branches"`
	ThoughtHistoryLength int      `json:"thoughtHistoryLength"`
}

// JSONMarshaler インターフェースを定義
type JSONMarshaler interface {
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
}

// DefaultJSONMarshaler は標準のjson.MarshalIndentを使用する実装
type DefaultJSONMarshaler struct{}

func (m *DefaultJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// OutputWriter インターフェースを定義
type OutputWriter interface {
	Fprintln(w interface{}, a ...interface{}) (n int, err error)
}

// DefaultOutputWriter は標準のfmt.Fprintlnを使用する実装
type DefaultOutputWriter struct{}

func (w *DefaultOutputWriter) Fprintln(writer interface{}, a ...interface{}) (n int, err error) {
	return fmt.Fprintln(writer.(interface{ Write([]byte) (int, error) }), a...)
}

// SequentialThinkingService はシーケンシャルシンキングを行うサービスです
type SequentialThinkingService struct {
	ThoughtHistory []ThoughtData
	Branches       map[string][]ThoughtData
	jsonMarshaler  JSONMarshaler
	outputWriter   OutputWriter
}

// NewSequentialThinkingService は新しいSequentialThinkingServiceを作成します
func NewSequentialThinkingService() *SequentialThinkingService {
	return &SequentialThinkingService{
		ThoughtHistory: []ThoughtData{},
		Branches:       make(map[string][]ThoughtData),
		jsonMarshaler:  &DefaultJSONMarshaler{},
		outputWriter:   &DefaultOutputWriter{},
	}
}

// NewSequentialThinkingServiceWithDependencies はテスト用に依存性を注入できるSequentialThinkingServiceを作成します
func NewSequentialThinkingServiceWithDependencies(jsonMarshaler JSONMarshaler, outputWriter OutputWriter) *SequentialThinkingService {
	return &SequentialThinkingService{
		ThoughtHistory: []ThoughtData{},
		Branches:       make(map[string][]ThoughtData),
		jsonMarshaler:  jsonMarshaler,
		outputWriter:   outputWriter,
	}
}

// ValidateThoughtData は思考データを検証します
func (s *SequentialThinkingService) ValidateThoughtData(args map[string]interface{}) (ThoughtData, error) {
	thought, ok := args["thought"].(string)
	if !ok || thought == "" {
		return ThoughtData{}, fmt.Errorf("invalid thought: must be a string")
	}

	thoughtNumber, ok := args["thoughtNumber"].(float64)
	if !ok || thoughtNumber < 1 {
		return ThoughtData{}, fmt.Errorf("invalid thoughtNumber: must be a number greater than 0")
	}

	totalThoughts, ok := args["totalThoughts"].(float64)
	if !ok || totalThoughts < 1 {
		return ThoughtData{}, fmt.Errorf("invalid totalThoughts: must be a number greater than 0")
	}

	nextThoughtNeeded, ok := args["nextThoughtNeeded"].(bool)
	if !ok {
		return ThoughtData{}, fmt.Errorf("invalid nextThoughtNeeded: must be a boolean")
	}

	data := ThoughtData{
		Thought:           thought,
		ThoughtNumber:     int(thoughtNumber),
		TotalThoughts:     int(totalThoughts),
		NextThoughtNeeded: nextThoughtNeeded,
	}

	// オプションフィールドの処理
	if isRevision, ok := args["isRevision"].(bool); ok {
		data.IsRevision = isRevision
	}

	if revisesThought, ok := args["revisesThought"].(float64); ok {
		data.RevisesThought = int(revisesThought)
	}

	if branchFromThought, ok := args["branchFromThought"].(float64); ok {
		data.BranchFromThought = int(branchFromThought)
	}

	if branchID, ok := args["branchId"].(string); ok {
		data.BranchID = branchID
	}

	if needsMoreThoughts, ok := args["needsMoreThoughts"].(bool); ok {
		data.NeedsMoreThoughts = needsMoreThoughts
	}

	return data, nil
}

// FormatThought は思考データをフォーマットします
func (s *SequentialThinkingService) FormatThought(data ThoughtData) string {
	var prefix, context string

	if data.IsRevision {
		prefix = "🔄 Revision"
		context = fmt.Sprintf(" (revising thought %d)", data.RevisesThought)
	} else if data.BranchFromThought > 0 {
		prefix = "🌿 Branch"
		context = fmt.Sprintf(" (from thought %d, ID: %s)", data.BranchFromThought, data.BranchID)
	} else {
		prefix = "💭 Thought"
		context = ""
	}

	header := fmt.Sprintf("%s %d/%d%s", prefix, data.ThoughtNumber, data.TotalThoughts, context)
	border := "─"
	borderLength := max(len(header), len(data.Thought)) + 4
	borderLine := "┌" + repeatString(border, borderLength) + "┐"

	return fmt.Sprintf(`
%s
│ %s │
├%s┤
│ %s │
└%s┘`,
		borderLine,
		header,
		repeatString(border, borderLength),
		data.Thought,
		repeatString(border, borderLength))
}

// ProcessThought は思考データを処理します
func (s *SequentialThinkingService) ProcessThought(args map[string]interface{}) (string, error) {
	data, err := s.ValidateThoughtData(args)
	if err != nil {
		return "", err
	}

	// 思考番号が総思考数を超える場合、総思考数を更新
	if data.ThoughtNumber > data.TotalThoughts {
		data.TotalThoughts = data.ThoughtNumber
	}

	// 思考履歴に追加
	s.ThoughtHistory = append(s.ThoughtHistory, data)

	// ブランチ処理
	if data.BranchFromThought > 0 && data.BranchID != "" {
		if _, ok := s.Branches[data.BranchID]; !ok {
			s.Branches[data.BranchID] = []ThoughtData{}
		}
		s.Branches[data.BranchID] = append(s.Branches[data.BranchID], data)
	}

	// フォーマットされた思考を表示
	formattedThought := s.FormatThought(data)
	s.outputWriter.Fprintln(os.Stderr, formattedThought)

	// 結果を作成
	result := ProcessResult{
		ThoughtNumber:        data.ThoughtNumber,
		TotalThoughts:        data.TotalThoughts,
		NextThoughtNeeded:    data.NextThoughtNeeded,
		Branches:             getKeys(s.Branches),
		ThoughtHistoryLength: len(s.ThoughtHistory),
	}

	jsonResult, err := s.jsonMarshaler.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonResult), nil
}

// HandleSequentialThinking はシーケンシャルシンキングのハンドラーです
func (s *SequentialThinkingService) HandleSequentialThinking(args map[string]interface{}) (string, error) {
	return s.ProcessThought(args)
}

// ヘルパー関数

// max は2つの整数の最大値を返します
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// repeatString は文字列を指定された回数繰り返します
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// getKeys はマップのキーをスライスとして返します
func getKeys(m map[string][]ThoughtData) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
