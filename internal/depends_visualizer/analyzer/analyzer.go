package analyzer

import (
	"fmt"
	"os"
	"strings"

	"github.com/landmaster135/devbox/internal/depends_visualizer/config"
)

// AnalyzeFile はファイルから関数依存関係を解析します
func AnalyzeFile(filePath string, extension string) (map[string][]string, error) {
	// ファイルを読み込み
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗: %w", err)
	}

	// 行に分割
	lines := strings.Split(string(data), "\n")

	// 関数を抽出
	functions, err := ExtractFunctions(lines, extension)
	if err != nil {
		return nil, fmt.Errorf("関数の抽出に失敗: %w", err)
	}

	// 依存関係を解析
	dependencies, err := AnalyzeDependencies(lines, functions, extension)
	if err != nil {
		return nil, fmt.Errorf("依存関係の解析に失敗: %w", err)
	}

	return dependencies, nil
}

// AnalyzeDependencies はソースコードの行から依存関係を解析します
func AnalyzeDependencies(lines []string, functions []string, extension string) (map[string][]string, error) {
	// 言語設定を取得
	lang, ok := config.GetLanguageConfig(extension)
	if !ok {
		return nil, fmt.Errorf("サポートされていない拡張子: %s", extension)
	}

	// 結果を格納するマップ
	result := make(map[string][]string)

	// 現在解析中の関数
	var currentFunction string

	for _, line := range lines {
		// 空行はスキップ
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 先頭の空白を削除
		cleanLine := removeHeadSpaces(line, config.GetSpaces())

		// 関数の定義行を検出
		if strings.HasPrefix(cleanLine, lang.FunctionHeader) {
			tailIdx := strings.Index(cleanLine, lang.FunctionTail)
			if tailIdx > 0 {
				funcName := cleanLine[len(lang.FunctionHeader):tailIdx]
				currentFunction = funcName
				result[currentFunction] = []string{}
				continue
			}
		}

		// メインマーカーを検出した場合は関数外に出る
		if cleanLine == lang.MainMarker {
			currentFunction = ""
		}

		// 関数内でない場合はスキップ
		if currentFunction == "" {
			continue
		}

		// 行内の関数参照を検出
		refs := findFunctionReferences(line, functions)
		for _, ref := range refs {
			// 自分自身への参照はスキップ
			if ref == currentFunction {
				continue
			}

			// 重複を避けて追加
			if !contains(result[currentFunction], ref) {
				result[currentFunction] = append(result[currentFunction], ref)
			}
		}
	}

	return result, nil
}

// スライスに要素が含まれているか確認
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// 解析結果を表す構造体
type AnalysisResult struct {
	FilePath     string
	Dependencies map[string][]string
}
