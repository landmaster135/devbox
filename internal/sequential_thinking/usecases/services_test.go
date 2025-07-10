package usecases

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// #==============================================================#
// ##          Mock Implementations                              ##
// #==============================================================#

// MockJSONMarshaler は JSONMarshaler のモック実装です
type MockJSONMarshaler struct {
	MarshalIndentFunc func(v interface{}, prefix, indent string) ([]byte, error)
}

func (m *MockJSONMarshaler) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	if m.MarshalIndentFunc != nil {
		return m.MarshalIndentFunc(v, prefix, indent)
	}
	return json.MarshalIndent(v, prefix, indent)
}

// MockOutputWriter は OutputWriter のモック実装です
type MockOutputWriter struct {
	FprintlnFunc func(w interface{}, a ...interface{}) (n int, err error)
	Output       []string
}

func (m *MockOutputWriter) Fprintln(w interface{}, a ...interface{}) (n int, err error) {
	if m.FprintlnFunc != nil {
		return m.FprintlnFunc(w, a...)
	}
	// 出力内容を記録
	if len(a) > 0 {
		m.Output = append(m.Output, fmt.Sprint(a...))
	}
	return len(a), nil
}

// #==============================================================#
// ##          Test Classes                                      ##
// #==============================================================#

// TestSequentialThinkingService はSequentialThinkingServiceのテストクラスです
type TestSequentialThinkingService struct {
	service *SequentialThinkingService
}

// NewTestSequentialThinkingService はテスト用のSequentialThinkingServiceを作成します
func NewTestSequentialThinkingService() *TestSequentialThinkingService {
	return &TestSequentialThinkingService{
		service: NewSequentialThinkingService(),
	}
}

// NewTestSequentialThinkingServiceWithMocks はモック付きのテスト用SequentialThinkingServiceを作成します
func NewTestSequentialThinkingServiceWithMocks(jsonMarshaler JSONMarshaler, outputWriter OutputWriter) *TestSequentialThinkingService {
	return &TestSequentialThinkingService{
		service: NewSequentialThinkingServiceWithDependencies(jsonMarshaler, outputWriter),
	}
}

// #==============================================================#
// ##          Constructor Tests                                 ##
// #==============================================================#

// TestNewSequentialThinkingService_Normal はNewSequentialThinkingServiceの正常系をテストします
func TestNewSequentialThinkingService_Normal(t *testing.T) {
	// Arrange & Act
	service := NewSequentialThinkingService()

	// Assert
	if service == nil {
		t.Fatal("NewSequentialThinkingService() returned nil")
	}

	if service.ThoughtHistory == nil {
		t.Error("ThoughtHistory should be initialized")
	}

	if service.Branches == nil {
		t.Error("Branches should be initialized")
	}

	if len(service.ThoughtHistory) != 0 {
		t.Errorf("ThoughtHistory should be empty, got length: %d", len(service.ThoughtHistory))
	}

	if len(service.Branches) != 0 {
		t.Errorf("Branches should be empty, got length: %d", len(service.Branches))
	}
}

// TestNewSequentialThinkingServiceWithDependencies_Normal はNewSequentialThinkingServiceWithDependenciesの正常系をテストします
func TestNewSequentialThinkingServiceWithDependencies_Normal(t *testing.T) {
	// Arrange
	mockJSON := &MockJSONMarshaler{}
	mockOutput := &MockOutputWriter{}

	// Act
	service := NewSequentialThinkingServiceWithDependencies(mockJSON, mockOutput)

	// Assert
	if service == nil {
		t.Fatal("NewSequentialThinkingServiceWithDependencies() returned nil")
	}

	if service.jsonMarshaler != mockJSON {
		t.Error("jsonMarshaler should be set to provided mock")
	}

	if service.outputWriter != mockOutput {
		t.Error("outputWriter should be set to provided mock")
	}
}

// #==============================================================#
// ##          ValidateThoughtData Tests                         ##
// #==============================================================#

// TestSequentialThinkingService_ValidateThoughtData_Normal はValidateThoughtDataの正常系をテストします
func TestSequentialThinkingService_ValidateThoughtData_Normal(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought":           "This is a test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}

	// Act
	result, err := testService.service.ValidateThoughtData(args)

	// Assert
	if err != nil {
		t.Errorf("ValidateThoughtData() error = %v, want nil", err)
	}

	if result.Thought != "This is a test thought" {
		t.Errorf("Thought = %v, want %v", result.Thought, "This is a test thought")
	}

	if result.ThoughtNumber != 1 {
		t.Errorf("ThoughtNumber = %v, want %v", result.ThoughtNumber, 1)
	}

	if result.TotalThoughts != 3 {
		t.Errorf("TotalThoughts = %v, want %v", result.TotalThoughts, 3)
	}

	if result.NextThoughtNeeded != true {
		t.Errorf("NextThoughtNeeded = %v, want %v", result.NextThoughtNeeded, true)
	}
}

// TestSequentialThinkingService_ValidateThoughtData_WithOptionalFields はオプションフィールド付きのValidateThoughtDataをテストします
func TestSequentialThinkingService_ValidateThoughtData_WithOptionalFields(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought":           "Revision thought",
		"thoughtNumber":     float64(2),
		"totalThoughts":     float64(5),
		"nextThoughtNeeded": false,
		"isRevision":        true,
		"revisesThought":    float64(1),
		"branchFromThought": float64(1),
		"branchId":          "branch-1",
		"needsMoreThoughts": true,
	}

	// Act
	result, err := testService.service.ValidateThoughtData(args)

	// Assert
	if err != nil {
		t.Errorf("ValidateThoughtData() error = %v, want nil", err)
	}

	if result.IsRevision != true {
		t.Errorf("IsRevision = %v, want %v", result.IsRevision, true)
	}

	if result.RevisesThought != 1 {
		t.Errorf("RevisesThought = %v, want %v", result.RevisesThought, 1)
	}

	if result.BranchFromThought != 1 {
		t.Errorf("BranchFromThought = %v, want %v", result.BranchFromThought, 1)
	}

	if result.BranchID != "branch-1" {
		t.Errorf("BranchID = %v, want %v", result.BranchID, "branch-1")
	}

	if result.NeedsMoreThoughts != true {
		t.Errorf("NeedsMoreThoughts = %v, want %v", result.NeedsMoreThoughts, true)
	}
}

// TestSequentialThinkingService_ValidateThoughtData_InvalidThought は無効なthoughtのテストです
func TestSequentialThinkingService_ValidateThoughtData_InvalidThought(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Empty thought",
			args: map[string]interface{}{
				"thought":           "",
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Missing thought",
			args: map[string]interface{}{
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Non-string thought",
			args: map[string]interface{}{
				"thought":           123,
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := testService.service.ValidateThoughtData(tt.args)

			// Assert
			if err == nil {
				t.Error("ValidateThoughtData() should return error for invalid thought")
			}

			if !strings.Contains(err.Error(), "invalid thought") {
				t.Errorf("Error should contain 'invalid thought', got: %v", err.Error())
			}
		})
	}
}

// TestSequentialThinkingService_ValidateThoughtData_InvalidThoughtNumber は無効なthoughtNumberのテストです
func TestSequentialThinkingService_ValidateThoughtData_InvalidThoughtNumber(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Zero thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(0),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Negative thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(-1),
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Missing thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Non-number thoughtNumber",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     "invalid",
				"totalThoughts":     float64(3),
				"nextThoughtNeeded": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := testService.service.ValidateThoughtData(tt.args)

			// Assert
			if err == nil {
				t.Error("ValidateThoughtData() should return error for invalid thoughtNumber")
			}

			if !strings.Contains(err.Error(), "invalid thoughtNumber") {
				t.Errorf("Error should contain 'invalid thoughtNumber', got: %v", err.Error())
			}
		})
	}
}

// TestSequentialThinkingService_ValidateThoughtData_InvalidTotalThoughts は無効なtotalThoughtsのテストです
func TestSequentialThinkingService_ValidateThoughtData_InvalidTotalThoughts(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "Zero totalThoughts",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(1),
				"totalThoughts":     float64(0),
				"nextThoughtNeeded": true,
			},
		},
		{
			name: "Missing totalThoughts",
			args: map[string]interface{}{
				"thought":           "Test thought",
				"thoughtNumber":     float64(1),
				"nextThoughtNeeded": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			_, err := testService.service.ValidateThoughtData(tt.args)

			// Assert
			if err == nil {
				t.Error("ValidateThoughtData() should return error for invalid totalThoughts")
			}

			if !strings.Contains(err.Error(), "invalid totalThoughts") {
				t.Errorf("Error should contain 'invalid totalThoughts', got: %v", err.Error())
			}
		})
	}
}

// TestSequentialThinkingService_ValidateThoughtData_InvalidNextThoughtNeeded は無効なnextThoughtNeededのテストです
func TestSequentialThinkingService_ValidateThoughtData_InvalidNextThoughtNeeded(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought":       "Test thought",
		"thoughtNumber": float64(1),
		"totalThoughts": float64(3),
		// nextThoughtNeeded is missing
	}

	// Act
	_, err := testService.service.ValidateThoughtData(args)

	// Assert
	if err == nil {
		t.Error("ValidateThoughtData() should return error for missing nextThoughtNeeded")
	}

	if !strings.Contains(err.Error(), "invalid nextThoughtNeeded") {
		t.Errorf("Error should contain 'invalid nextThoughtNeeded', got: %v", err.Error())
	}
}

// #==============================================================#
// ##          FormatThought Tests                               ##
// #==============================================================#

// TestSequentialThinkingService_FormatThought_Normal は通常の思考フォーマットをテストします
func TestSequentialThinkingService_FormatThought_Normal(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	data := ThoughtData{
		Thought:           "This is a test thought",
		ThoughtNumber:     1,
		TotalThoughts:     3,
		NextThoughtNeeded: true,
	}

	// Act
	result := testService.service.FormatThought(data)

	// Assert
	if !strings.Contains(result, "💭 Thought 1/3") {
		t.Errorf("Result should contain '💭 Thought 1/3', got: %v", result)
	}

	if !strings.Contains(result, "This is a test thought") {
		t.Errorf("Result should contain thought text, got: %v", result)
	}

	if !strings.Contains(result, "┌") || !strings.Contains(result, "┐") {
		t.Errorf("Result should contain border characters, got: %v", result)
	}
}

// TestSequentialThinkingService_FormatThought_Revision はリビジョン思考のフォーマットをテストします
func TestSequentialThinkingService_FormatThought_Revision(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	data := ThoughtData{
		Thought:           "This is a revision",
		ThoughtNumber:     2,
		TotalThoughts:     3,
		NextThoughtNeeded: true,
		IsRevision:        true,
		RevisesThought:    1,
	}

	// Act
	result := testService.service.FormatThought(data)

	// Assert
	if !strings.Contains(result, "🔄 Revision 2/3") {
		t.Errorf("Result should contain '🔄 Revision 2/3', got: %v", result)
	}

	if !strings.Contains(result, "(revising thought 1)") {
		t.Errorf("Result should contain revision context, got: %v", result)
	}
}

// TestSequentialThinkingService_FormatThought_Branch はブランチ思考のフォーマットをテストします
func TestSequentialThinkingService_FormatThought_Branch(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	data := ThoughtData{
		Thought:           "This is a branch",
		ThoughtNumber:     2,
		TotalThoughts:     3,
		NextThoughtNeeded: true,
		BranchFromThought: 1,
		BranchID:          "branch-1",
	}

	// Act
	result := testService.service.FormatThought(data)

	// Assert
	if !strings.Contains(result, "🌿 Branch 2/3") {
		t.Errorf("Result should contain '🌿 Branch 2/3', got: %v", result)
	}

	if !strings.Contains(result, "(from thought 1, ID: branch-1)") {
		t.Errorf("Result should contain branch context, got: %v", result)
	}
}

// #==============================================================#
// ##          ProcessThought Tests                              ##
// #==============================================================#

// TestSequentialThinkingService_ProcessThought_Normal はProcessThoughtの正常系をテストします
func TestSequentialThinkingService_ProcessThought_Normal(t *testing.T) {
	// Arrange
	mockOutput := &MockOutputWriter{}
	testService := NewTestSequentialThinkingServiceWithMocks(&MockJSONMarshaler{}, mockOutput)
	args := map[string]interface{}{
		"thought":           "Test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}

	// Act
	result, err := testService.service.ProcessThought(args)

	// Assert
	if err != nil {
		t.Errorf("ProcessThought() error = %v, want nil", err)
	}

	if result == "" {
		t.Error("ProcessThought() should return non-empty result")
	}

	// 思考履歴に追加されているか確認
	if len(testService.service.ThoughtHistory) != 1 {
		t.Errorf("ThoughtHistory length = %v, want 1", len(testService.service.ThoughtHistory))
	}

	// 出力が呼ばれているか確認
	if len(mockOutput.Output) == 0 {
		t.Error("Output should be called")
	}

	// JSON結果の確認
	var processResult ProcessResult
	if err := json.Unmarshal([]byte(result), &processResult); err != nil {
		t.Errorf("Result should be valid JSON: %v", err)
	}

	if processResult.ThoughtNumber != 1 {
		t.Errorf("ProcessResult.ThoughtNumber = %v, want 1", processResult.ThoughtNumber)
	}
}


// TestSequentialThinkingService_ProcessThought_ThoughtNumberAdjustment は思考番号の自動調整をテストします
func TestSequentialThinkingService_ProcessThought_ThoughtNumberAdjustment(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought":           "Test thought",
		"thoughtNumber":     float64(5),
		"totalThoughts":     float64(3), // thoughtNumberより小さい
		"nextThoughtNeeded": true,
	}

	// Act
	result, err := testService.service.ProcessThought(args)

	// Assert
	if err != nil {
		t.Errorf("ProcessThought() error = %v, want nil", err)
	}

	var processResult ProcessResult
	if err := json.Unmarshal([]byte(result), &processResult); err != nil {
		t.Errorf("Result should be valid JSON: %v", err)
	}

	// totalThoughtsが自動調整されているか確認
	if processResult.TotalThoughts != 5 {
		t.Errorf("ProcessResult.TotalThoughts = %v, want 5", processResult.TotalThoughts)
	}
}

// TestSequentialThinkingService_ProcessThought_BranchHandling はブランチ処理をテストします
func TestSequentialThinkingService_ProcessThought_BranchHandling(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought":           "Branch thought",
		"thoughtNumber":     float64(2),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
		"branchFromThought": float64(1),
		"branchId":          "test-branch",
	}

	// Act
	result, err := testService.service.ProcessThought(args)

	// Assert
	if err != nil {
		t.Errorf("ProcessThought() error = %v, want nil", err)
	}

	// ブランチが作成されているか確認
	if _, exists := testService.service.Branches["test-branch"]; !exists {
		t.Error("Branch 'test-branch' should be created")
	}

	if len(testService.service.Branches["test-branch"]) != 1 {
		t.Errorf("Branch 'test-branch' should have 1 thought, got %v", len(testService.service.Branches["test-branch"]))
	}

	var processResult ProcessResult
	if err := json.Unmarshal([]byte(result), &processResult); err != nil {
		t.Errorf("Result should be valid JSON: %v", err)
	}

	// ブランチがキーに含まれているか確認
	found := false
	for _, branch := range processResult.Branches {
		if branch == "test-branch" {
			found = true
			break
		}
	}
	if !found {
		t.Error("ProcessResult.Branches should contain 'test-branch'")
	}
}

// TestSequentialThinkingService_ProcessThought_ValidationError は検証エラーのテストです
func TestSequentialThinkingService_ProcessThought_ValidationError(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought": "", // 無効な思考
	}

	// Act
	result, err := testService.service.ProcessThought(args)

	// Assert
	if err == nil {
		t.Error("ProcessThought() should return error for invalid input")
	}

	if result != "" {
		t.Errorf("ProcessThought() should return empty result on error, got: %v", result)
	}
}

// TestSequentialThinkingService_ProcessThought_JSONMarshalError はJSON変換エラーのテストです
func TestSequentialThinkingService_ProcessThought_JSONMarshalError(t *testing.T) {
	// Arrange
	mockJSON := &MockJSONMarshaler{
		MarshalIndentFunc: func(v interface{}, prefix, indent string) ([]byte, error) {
			return nil, fmt.Errorf("JSON marshal error")
		},
	}
	testService := NewTestSequentialThinkingServiceWithMocks(mockJSON, &MockOutputWriter{})
	args := map[string]interface{}{
		"thought":           "Test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}

	// Act
	result, err := testService.service.ProcessThought(args)

	// Assert
	if err == nil {
		t.Error("ProcessThought() should return error when JSON marshal fails")
	}

	if result != "" {
		t.Errorf("ProcessThought() should return empty result on error, got: %v", result)
	}
}

// #==============================================================#
// ##          HandleSequentialThinking Tests                    ##
// #==============================================================#

// TestSequentialThinkingService_HandleSequentialThinking_Normal はHandleSequentialThinkingの正常系をテストします
func TestSequentialThinkingService_HandleSequentialThinking_Normal(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()
	args := map[string]interface{}{
		"thought":           "Test thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}

	// Act
	result, err := testService.service.HandleSequentialThinking(args)

	// Assert
	if err != nil {
		t.Errorf("HandleSequentialThinking() error = %v, want nil", err)
	}

	if result == "" {
		t.Error("HandleSequentialThinking() should return non-empty result")
	}
}

// #==============================================================#
// ##          Helper Function Tests                             ##
// #==============================================================#

// TestMax_Normal はmax関数の正常系をテストします
func TestMax_Normal(t *testing.T) {
	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{
			name:     "a > b",
			a:        5,
			b:        3,
			expected: 5,
		},
		{
			name:     "a < b",
			a:        2,
			b:        7,
			expected: 7,
		},
		{
			name:     "a == b",
			a:        4,
			b:        4,
			expected: 4,
		},
		{
			name:     "negative numbers",
			a:        -2,
			b:        -5,
			expected: -2,
		},
		{
			name:     "zero and positive",
			a:        0,
			b:        3,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := max(tt.a, tt.b)

			// Assert
			if result != tt.expected {
				t.Errorf("max(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestRepeatString_Normal はrepeatString関数の正常系をテストします
func TestRepeatString_Normal(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		count    int
		expected string
	}{
		{
			name:     "normal case",
			s:        "a",
			count:    3,
			expected: "aaa",
		},
		{
			name:     "zero count",
			s:        "test",
			count:    0,
			expected: "",
		},
		{
			name:     "empty string",
			s:        "",
			count:    5,
			expected: "",
		},
		{
			name:     "single character",
			s:        "x",
			count:    1,
			expected: "x",
		},
		{
			name:     "multi-character string",
			s:        "ab",
			count:    3,
			expected: "ababab",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := repeatString(tt.s, tt.count)

			// Assert
			if result != tt.expected {
				t.Errorf("repeatString(%v, %v) = %v, want %v", tt.s, tt.count, result, tt.expected)
			}
		})
	}
}

// TestGetKeys_Normal はgetKeys関数の正常系をテストします
func TestGetKeys_Normal(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string][]ThoughtData
		expected []string
	}{
		{
			name:     "empty map",
			input:    map[string][]ThoughtData{},
			expected: []string{},
		},
		{
			name: "single key",
			input: map[string][]ThoughtData{
				"key1": {},
			},
			expected: []string{"key1"},
		},
		{
			name: "multiple keys",
			input: map[string][]ThoughtData{
				"key1": {},
				"key2": {},
				"key3": {},
			},
			expected: []string{"key1", "key2", "key3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := getKeys(tt.input)

			// Assert
			if len(result) != len(tt.expected) {
				t.Errorf("getKeys() length = %v, want %v", len(result), len(tt.expected))
			}

			// キーが全て含まれているか確認（順序は保証されない）
			for _, expectedKey := range tt.expected {
				found := false
				for _, actualKey := range result {
					if actualKey == expectedKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("getKeys() should contain key %v, got %v", expectedKey, result)
				}
			}
		})
	}
}

// #==============================================================#
// ##          Integration Tests                                 ##
// #==============================================================#

// TestSequentialThinkingService_Integration_MultipleThoughts は複数思考の統合テストです
func TestSequentialThinkingService_Integration_MultipleThoughts(t *testing.T) {
	// Arrange
	testService := NewTestSequentialThinkingService()

	// Act & Assert - 最初の思考
	args1 := map[string]interface{}{
		"thought":           "First thought",
		"thoughtNumber":     float64(1),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}
	result1, err1 := testService.service.ProcessThought(args1)
	if err1 != nil {
		t.Errorf("First ProcessThought() error = %v", err1)
	}

	// Act & Assert - 二番目の思考
	args2 := map[string]interface{}{
		"thought":           "Second thought",
		"thoughtNumber":     float64(2),
		"totalThoughts":     float64(3),
		"nextThoughtNeeded": true,
	}
	result2, err2 := testService.service.ProcessThought(args2)
	if err2 != nil {
		t.Errorf("Second ProcessThought() error = %v", err2)
	}

	// Assert - 思考履歴の確認
	if len(testService.service.ThoughtHistory) != 2 {
		t.Errorf("ThoughtHistory length = %v, want 2", len(testService.service.ThoughtHistory))
	}

	// Assert - 結果の確認
	if result1 == "" || result2 == "" {
		t.Error("ProcessThought() should return non-empty results")
	}

	// Assert - 最初の思考の内容確認
	if testService.service.ThoughtHistory[0].Thought != "First thought" {
		t.Errorf("First thought = %v, want 'First thought'", testService.service.ThoughtHistory[0].Thought)
	}

	// Assert - 二番目の思考の内容確認
	if testService.service.ThoughtHistory[1].Thought != "Second thought" {
		t.Errorf("Second thought = %v, want 'Second thought'", testService.service.ThoughtHistory[1].Thought)
	}
}
