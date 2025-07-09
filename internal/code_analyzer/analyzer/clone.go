// internal/analyzer/clone.go
package analyzer

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/landmaster135/devbox/internal/code_analyzer/models"
)

// CloneDetector はコードクローン検出機能を提供します
type CloneDetector struct {
	MinBlockSize  int
	MinSimilarity float64
}

// NewCloneDetector は新しいCloneDetectorインスタンスを作成します
func NewCloneDetector(minBlockSize int, minSimilarity float64) *CloneDetector {
	return &CloneDetector{
		MinBlockSize:  minBlockSize,
		MinSimilarity: minSimilarity,
	}
}

// トークン化して類似コードブロックを検出するための構造体
type codeBlock struct {
	file   string
	line   int
	code   string
	hash   string
	size   int
	tokens []string
}

// DetectClones はプロジェクト内のコードクローンを検出します
func (d *CloneDetector) DetectClones(files []models.FileMetrics, contents map[string]string) []models.CodeClone {
	var allBlocks []codeBlock
	var clones []models.CodeClone

	// 各ファイルからコードブロックを抽出
	for _, file := range files {
		if content, ok := contents[file.Path]; ok {
			blocks := d.extractCodeBlocks(file.Path, content)
			allBlocks = append(allBlocks, blocks...)
		}
	}

	// ハッシュによるブロックのグループ化
	blocksByHash := make(map[string][]codeBlock)
	for _, block := range allBlocks {
		if block.hash != "" {
			blocksByHash[block.hash] = append(blocksByHash[block.hash], block)
		}
	}

	// 完全一致するブロックをクローンとして検出
	for _, blocks := range blocksByHash {
		if len(blocks) < 2 {
			continue
		}

		// 同じハッシュを持つブロック同士を比較
		for i := 0; i < len(blocks); i++ {
			for j := i + 1; j < len(blocks); j++ {
				// 同じファイル内のクローンは無視する
				if blocks[i].file == blocks[j].file {
					continue
				}

				// 類似度計算
				similarity := d.calculateSimilarity(blocks[i].tokens, blocks[j].tokens)

				if similarity >= d.MinSimilarity {
					sourceCode := strings.Join(blocks[i].tokens, " ")

					clone := models.CodeClone{
						SourceFile: blocks[i].file,
						TargetFile: blocks[j].file,
						SourceLine: blocks[i].line,
						TargetLine: blocks[j].line,
						LineCount:  blocks[i].size / 5, // 概算：トークン数/行あたりの平均トークン数
						Similarity: similarity,
						Code:       sourceCode[:min(len(sourceCode), 100)] + "...", // プレビュー用に先頭部分だけ
					}
					clones = append(clones, clone)
				}
			}
		}
	}

	return clones
}

// コードブロック抽出
func (d *CloneDetector) extractCodeBlocks(path string, content string) []codeBlock {
	lines := strings.Split(content, "\n")
	tokens := tokenizeCode(content)

	// 行番号とトークンのマッピングを作成
	tokenPositions := make(map[int]int)
	lineCount := 0
	tokenCount := 0

	for i, line := range lines {
		lineTokens := tokenizeCode(line)
		tokenPositions[i] = tokenCount
		tokenCount += len(lineTokens)
		lineCount++
	}

	// コードブロックの抽出
	var blocks []codeBlock

	for i := 0; i <= len(tokens)-d.MinBlockSize; i++ {
		blockTokens := tokens[i : i+d.MinBlockSize]
		blockHash := d.hashCodeBlock(blockTokens)

		// 対応する行番号を見つける
		var lineNumber int
		for j := 0; j < lineCount; j++ {
			if tokenPositions[j] <= i && (j == lineCount-1 || tokenPositions[j+1] > i) {
				lineNumber = j + 1 // 1から始まる行番号
				break
			}
		}

		blocks = append(blocks, codeBlock{
			file:   path,
			line:   lineNumber,
			hash:   blockHash,
			size:   d.MinBlockSize,
			tokens: blockTokens,
		})
	}

	return blocks
}

// コード行からハッシュを生成
func (d *CloneDetector) hashCodeBlock(tokens []string) string {
	if len(tokens) < d.MinBlockSize {
		return ""
	}

	hash := md5.New()
	for i := 0; i < d.MinBlockSize; i++ {
		hash.Write([]byte(tokens[i]))
	}
	return hex.EncodeToString(hash.Sum(nil))
}

// 類似度計算
func (d *CloneDetector) calculateSimilarity(tokensA, tokensB []string) float64 {
	// トークンの出現頻度を計算
	freqA := make(map[string]int)
	freqB := make(map[string]int)

	for _, token := range tokensA {
		freqA[token]++
	}

	for _, token := range tokensB {
		freqB[token]++
	}

	// Jaccard類似度の計算
	intersection := 0
	for token, countA := range freqA {
		if countB, exists := freqB[token]; exists {
			intersection += min(countA, countB)
		}
	}

	union := 0
	for _, count := range freqA {
		union += count
	}
	for _, count := range freqB {
		union += count
	}
	union -= intersection

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
