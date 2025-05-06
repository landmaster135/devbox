package analyzer

import (
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/config"
)

// ExtractFunctions はテキスト行から関数を抽出します
func ExtractFunctions(lines []string, extension string) ([]string, error) {
	// 言語設定を取得
	lang, ok := config.GetLanguageConfig(extension)
	if !ok {
		return nil, nil
	}

	spaces := config.GetSpaces()
	headOfTarget := lang.FunctionHeader
	tailOfTarget := lang.FunctionTail

	var functions []string

	for _, line := range lines {
		// 先頭の空白を削除
		cleanLine := removeHeadSpaces(line, spaces)

		// 関数ヘッダーが先頭にあるか確認
		headIndex := strings.Index(cleanLine, headOfTarget)
		if headIndex != 0 {
			continue
		}

		// 関数テイルを検索
		tailIndex := strings.Index(cleanLine, tailOfTarget)
		if tailIndex == -1 {
			continue
		}

		// 関数名を抽出
		funcName := cleanLine[len(headOfTarget):tailIndex]
		functions = append(functions, funcName)
	}

	return functions, nil
}

// 先頭の空白を削除
func removeHeadSpaces(text string, spaces []string) string {
	if len(text) == 0 {
		return text
	}

	for _, space := range spaces {
		if strings.HasPrefix(text, space) {
			return removeHeadSpaces(text[len(space):], spaces)
		}
	}

	return text
}

// テキスト内の関数参照を検索
func findFunctionReferences(line string, functions []string) []string {
	var references []string

	for _, function := range functions {
		if strings.Contains(line, function) {
			references = append(references, function)
		}
	}

	return references
}
