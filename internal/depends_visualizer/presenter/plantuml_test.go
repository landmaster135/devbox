package presenter

import (
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/depends_visualizer/analyzer"
)

// TestRenderPlantUML_Normal はRenderPlantUML関数の正常系テストです
func TestRenderPlantUML_Normal(t *testing.T) {
	// テスト用の解析結果
	results := []analyzer.AnalysisResult{
		{
			FilePath: "/path/to/file1.go",
			Dependencies: map[string][]string{
				"main":  {"hello", "world"},
				"hello": {},
				"world": {"hello"},
			},
		},
		{
			FilePath: "/path/to/file2.go",
			Dependencies: map[string][]string{
				"init":   {"setup"},
				"setup":  {},
				"helper": {"setup"},
			},
		},
	}

	// テスト実行
	output, err := RenderPlantUML(results)
	if err != nil {
		t.Fatalf("RenderPlantUML関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	expectedContents := []string{
		"@startuml",
		"skinparam defaultTextAlignment center",
		"' File: file1.go",
		"' File: file2.go",
		"main --> hello",
		"main --> world",
		"world --> hello",
		"init --> setup",
		"helper --> setup",
		"@enduml",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("出力に %s が含まれていません: %s", expected, output)
		}
	}
}

// TestRenderPlantUML_EmptyResults はRenderPlantUML関数の空結果テストです
func TestRenderPlantUML_EmptyResults(t *testing.T) {
	// 空の解析結果
	results := []analyzer.AnalysisResult{}

	// テスト実行
	output, err := RenderPlantUML(results)
	if err != nil {
		t.Fatalf("RenderPlantUML関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	expectedContents := []string{
		"@startuml",
		"skinparam defaultTextAlignment center",
		"@enduml",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("出力に %s が含まれていません: %s", expected, output)
		}
	}

	// 依存関係がないことを確認
	if strings.Contains(output, "-->") {
		t.Errorf("空の結果に対する出力に依存関係が含まれています: %s", output)
	}
}

// TestRenderPlantUML_SpecialCharacters はRenderPlantUML関数の特殊文字テストです
func TestRenderPlantUML_SpecialCharacters(t *testing.T) {
	// 特殊文字を含む解析結果
	results := []analyzer.AnalysisResult{
		{
			FilePath: "/path/to/file.go",
			Dependencies: map[string][]string{
				"func_with-dash":     {"normal_func"},
				"func_with__dunder":  {},
				"func_with space":    {},
				"normal_func":        {"func_with space"},
				"func_with(parens)":  {},
				"func_with[bracket]": {},
			},
		},
	}

	// テスト実行
	output, err := RenderPlantUML(results)
	if err != nil {
		t.Fatalf("RenderPlantUML関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	if !strings.Contains(output, "@startuml") {
		t.Errorf("出力にPlantUMLヘッダーが含まれていません: %s", output)
	}

	// 特殊文字が適切に処理されていることを確認
	// PlantUML形式では特殊文字を含む識別子は引用符で囲まれる
	specialChars := []string{
		"func_with-dash --> normal_func",
		"normal_func --> func_with space",
	}

	for _, special := range specialChars {
		if !strings.Contains(output, special) {
			t.Errorf("出力に特殊文字を含む依存関係 %s が適切に処理されていません: %s", special, output)
		}
	}

	// 終了タグの検証
	if !strings.Contains(output, "@enduml") {
		t.Errorf("出力に終了タグ @enduml が含まれていません: %s", output)
	}
}

// TestRenderPlantUML_ComplexDependencies はRenderPlantUML関数の複雑な依存関係テストです
func TestRenderPlantUML_ComplexDependencies(t *testing.T) {
	// 複雑な依存関係を持つ解析結果
	results := []analyzer.AnalysisResult{
		{
			FilePath: "/path/to/file1.go",
			Dependencies: map[string][]string{
				"main":    {"init", "helper1", "helper2", "helper3"},
				"init":    {},
				"helper1": {"util1", "util2"},
				"helper2": {"util2", "util3"},
				"helper3": {"util1", "util3"},
				"util1":   {},
				"util2":   {"util1"},
				"util3":   {"util2"},
			},
		},
	}

	// テスト実行
	output, err := RenderPlantUML(results)
	if err != nil {
		t.Fatalf("RenderPlantUML関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	expectedDependencies := []string{
		"main --> init",
		"main --> helper1",
		"main --> helper2",
		"main --> helper3",
		"helper1 --> util1",
		"helper1 --> util2",
		"helper2 --> util2",
		"helper2 --> util3",
		"helper3 --> util1",
		"helper3 --> util3",
		"util2 --> util1",
		"util3 --> util2",
	}

	for _, dep := range expectedDependencies {
		if !strings.Contains(output, dep) {
			t.Errorf("出力に依存関係 %s が含まれていません: %s", dep, output)
		}
	}

	// 循環依存関係が適切に処理されていることを確認
	// util2 -> util1, util3 -> util2, helper1 -> util2, helper3 -> util1 の循環
	circularDeps := []string{
		"util2 --> util1",
		"util3 --> util2",
	}

	for _, dep := range circularDeps {
		if !strings.Contains(output, dep) {
			t.Errorf("出力に循環依存関係 %s が含まれていません: %s", dep, output)
		}
	}
}
