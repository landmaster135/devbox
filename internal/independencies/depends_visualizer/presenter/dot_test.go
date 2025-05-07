package presenter

import (
	"strings"
	"testing"

	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/analyzer"
)

// TestRenderDOT_Normal はRenderDOT関数の正常系テストです
func TestRenderDOT_Normal(t *testing.T) {
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
	output, err := RenderDOT(results)
	if err != nil {
		t.Fatalf("RenderDOT関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	expectedContents := []string{
		"digraph G {",
		"rankdir=BT;",
		"node [shape=box",
		"// File: file1.go",
		"// File: file2.go",
		"subgraph cluster_0",
		"subgraph cluster_1",
		"\"main\"",
		"\"hello\"",
		"\"world\"",
		"\"init\"",
		"\"setup\"",
		"\"helper\"",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("出力に %s が含まれていません: %s", expected, output)
		}
	}

	// 依存関係の検証
	dependencyPairs := []string{
		"\"main\" -> \"hello\"",
		"\"main\" -> \"world\"",
		"\"world\" -> \"hello\"",
		"\"init\" -> \"setup\"",
		"\"helper\" -> \"setup\"",
	}

	for _, pair := range dependencyPairs {
		if !strings.Contains(output, pair) {
			t.Errorf("出力に依存関係 %s が含まれていません: %s", pair, output)
		}
	}

	// 終了タグの検証
	if !strings.Contains(output, "}") {
		t.Errorf("出力に終了タグ } が含まれていません: %s", output)
	}
}

// TestRenderDOT_EmptyResults はRenderDOT関数の空結果テストです
func TestRenderDOT_EmptyResults(t *testing.T) {
	// 空の解析結果
	results := []analyzer.AnalysisResult{}

	// テスト実行
	output, err := RenderDOT(results)
	if err != nil {
		t.Fatalf("RenderDOT関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	expectedContents := []string{
		"digraph G {",
		"rankdir=BT;",
		"node [shape=box",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(output, expected) {
			t.Errorf("出力に %s が含まれていません: %s", expected, output)
		}
	}

	// サブグラフがないことを確認
	if strings.Contains(output, "subgraph cluster_") {
		t.Errorf("空の結果に対する出力にサブグラフが含まれています: %s", output)
	}

	// 終了タグの検証
	if !strings.Contains(output, "}") {
		t.Errorf("出力に終了タグ } が含まれていません: %s", output)
	}
}

// TestRenderDOT_SpecialCharacters はRenderDOT関数の特殊文字テストです
func TestRenderDOT_SpecialCharacters(t *testing.T) {
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
	output, err := RenderDOT(results)
	if err != nil {
		t.Fatalf("RenderDOT関数の実行に失敗しました: %v", err)
	}

	// 出力の検証
	if !strings.Contains(output, "digraph G {") {
		t.Errorf("出力にDOTヘッダーが含まれていません: %s", output)
	}

	// 特殊文字が適切に処理されていることを確認
	// DOT形式では特殊文字を含む識別子は引用符で囲まれる
	specialChars := []string{
		"\"func_with-dash\"",
		"\"func_with__dunder\"",
		"\"func_with space\"",
		"\"func_with(parens)\"",
		"\"func_with[bracket]\"",
	}

	for _, special := range specialChars {
		if !strings.Contains(output, special) {
			t.Errorf("出力に特殊文字を含む識別子 %s が適切に処理されていません: %s", special, output)
		}
	}
}
