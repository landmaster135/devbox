// internal/analyzer/file.go
package analyzer

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/landmaster135/devbox/internal/code_analyzer/models"
)

// FileAnalyzer はファイル分析機能を提供します
type FileAnalyzer struct{}

// NewFileAnalyzer は新しいFileAnalyzerインスタンスを作成します
func NewFileAnalyzer() *FileAnalyzer {
	return &FileAnalyzer{}
}

// AnalyzeFile は1つのファイルを分析します
func (a *FileAnalyzer) AnalyzeFile(path string) (models.FileMetrics, string, error) {
	metrics := models.FileMetrics{Path: path}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return metrics, "", err
	}

	contentStr := string(content)
	total, code, comment, blank := countLines(contentStr)

	metrics.TotalLines = total
	metrics.CodeLines = code
	metrics.CommentLines = comment
	metrics.BlankLines = blank

	if code > 0 {
		metrics.CommentRatio = float64(comment) / float64(code) * 100.0
	}

	// Go言語の場合のみ構文解析を実行
	if strings.HasSuffix(strings.ToLower(path), ".go") {
		a.analyzeGoFile(&metrics, path)
	}

	// コードクローン検出用のトークンハッシュを作成
	metrics.TokenHashes = make(map[string]string)

	// 各行のハッシュを計算
	lines := strings.Split(contentStr, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		lineTokens := tokenizeCode(line)
		if len(lineTokens) > 0 {
			h := md5.New()
			for _, token := range lineTokens {
				h.Write([]byte(token))
			}
			metrics.TokenHashes[fmt.Sprintf("%d", i+1)] = hex.EncodeToString(h.Sum(nil))
		}
	}

	return metrics, contentStr, nil
}

// analyzeGoFile はGo言語ファイルの詳細分析を行います
func (a *FileAnalyzer) analyzeGoFile(metrics *models.FileMetrics, path string) {
	// 構文解析
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return // 構文解析エラーがあっても行数情報は返す
	}

	// 複雑度計算
	cv := &complexityVisitor{}
	ast.Walk(cv, f)
	metrics.Complexity = cv.complexity

	// 関数分析
	fv := &funcVisitor{}
	ast.Walk(fv, f)
	metrics.FunctionCount = fv.count

	if fv.count > 0 {
		metrics.AvgFunctionSize = float64(fv.totalSize) / float64(fv.count)
		metrics.MaxFunctionSize = fv.maxSize
	}
}

// 行数計測関数
func countLines(content string) (total, code, comment, blank int) {
	lines := strings.Split(content, "\n")
	total = len(lines)

	inBlockComment := false

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			blank++
			continue
		}

		if strings.HasPrefix(trimmedLine, "//") {
			comment++
			continue
		}

		if strings.HasPrefix(trimmedLine, "/*") {
			inBlockComment = true
			comment++

			if strings.Contains(trimmedLine, "*/") {
				inBlockComment = false
			}
			continue
		}

		if inBlockComment {
			comment++
			if strings.Contains(trimmedLine, "*/") {
				inBlockComment = false
			}
			continue
		}

		code++
	}

	return
}

// コードトークン化関数（コードクローン検出用）
func tokenizeCode(content string) []string {
	// 単純なトークン化
	re := regexp.MustCompile(`[\s\t\n\r]+`)
	content = re.ReplaceAllString(content, " ")

	// コメント除去
	content = regexp.MustCompile(`//.*`).ReplaceAllString(content, "")
	content = regexp.MustCompile(`/\*[\s\S]*?\*/`).ReplaceAllString(content, "")

	// 識別子、演算子、括弧などでトークン化
	tokenRegex := regexp.MustCompile(`([a-zA-Z_]\w*|\d+\.\d+|\d+|"[^"]*"|'[^']*'|==|!=|<=|>=|&&|\|\||[+\-*/=<>!&|;:,.(){}[\]])`)
	return tokenRegex.FindAllString(content, -1)
}
