package usecases

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// OSFileSystem は実際のファイルシステム操作を行う実装です
type OSFileSystem struct{}

// ReadDir はディレクトリの内容を読み取ります
func (fs *OSFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	return os.ReadDir(dirname)
}

// Stat はファイル情報を取得します
func (fs *OSFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// WriteFile はファイルに書き込みます
func (fs *OSFileSystem) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// ReadFile はファイルを読み取ります
func (fs *OSFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// MkdirAll はディレクトリを作成します
func (fs *OSFileSystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// StdinReader は標準入力から読み取る実装です
type StdinReader struct {
	reader *bufio.Reader
}

// NewStdinReader は新しい StdinReader を作成します
func NewStdinReader() *StdinReader {
	return &StdinReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadString は標準入力から文字列を読み取ります
func (r *StdinReader) ReadString(delim byte) (string, error) {
	return r.reader.ReadString(delim)
}

// DefaultREADMEParser はREADMEファイルの解析を行う実装です
type DefaultREADMEParser struct{}

// ParseUsageExamples はREADMEファイルから使用例を抽出します
func (p *DefaultREADMEParser) ParseUsageExamples(content []byte) ([]string, error) {
	lines := strings.Split(string(content), "\n")

	// まず「## 使用例」セクションを探す（最優先）
	usageExampleSectionFound := false
	usageExampleSectionIndex := -1

	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") &&
			strings.Contains(strings.ToLower(line), "使用例") &&
			!strings.Contains(strings.ToLower(line), "使用方法") {
			usageExampleSectionFound = true
			usageExampleSectionIndex = i
			break
		}
	}

	// 「使用例」セクションが見つからなかった場合は「使用方法」セクションを探す
	usageMethodSectionIndex := -1
	if !usageExampleSectionFound {
		for i, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "#") &&
				strings.Contains(strings.ToLower(line), "使用方法") {
				usageMethodSectionIndex = i
				break
			}
		}
	}

	// 使用例を抽出する開始位置を決定
	startIndex := -1
	if usageExampleSectionFound {
		startIndex = usageExampleSectionIndex
	} else if usageMethodSectionIndex >= 0 {
		startIndex = usageMethodSectionIndex
	}

	if startIndex < 0 {
		return []string{}, nil
	}

	inCodeBlock := false
	usageLines := []string{}

	// 開始位置から処理を開始
	for i := startIndex + 1; i < len(lines); i++ {
		line := lines[i]

		// 見出しに到達した場合の処理
		if strings.HasPrefix(strings.TrimSpace(line), "#") {
			// 現在「使用方法」セクションにいて、「使用例」セクションに到達した場合は続行
			if usageMethodSectionIndex >= 0 && !usageExampleSectionFound &&
				strings.Contains(strings.ToLower(line), "使用例") {
				usageExampleSectionFound = true
				continue
			}

			// コードブロック内にいる場合は、コードブロックが終了するまで続行
			if inCodeBlock {
				continue
			}

			// 「使用例」セクション内のサブセクションの場合は続行
			if usageExampleSectionFound {
				// 見出しレベルを確認（「##」より深い「###」などの場合は続行）
				currentLevel := 0
				for _, c := range strings.TrimSpace(line) {
					if c == '#' {
						currentLevel++
					} else {
						break
					}
				}

				// 「使用例」セクションの見出しレベルを確認
				exampleSectionLevel := 0
				if usageExampleSectionIndex >= 0 {
					exampleLine := lines[usageExampleSectionIndex]
					for _, c := range strings.TrimSpace(exampleLine) {
						if c == '#' {
							exampleSectionLevel++
						} else {
							break
						}
					}
				}

				// サブセクション（より深い見出し）の場合は続行
				if currentLevel > exampleSectionLevel {
					continue
				}
			}

			// それ以外の見出しなら終了
			break
		}

		trimmedLine := strings.TrimSpace(line)

		// コードブロックの開始を検出
		if !inCodeBlock && strings.HasPrefix(trimmedLine, "```") {
			inCodeBlock = true
			continue
		}

		// コードブロックの終了を検出
		if inCodeBlock && strings.HasPrefix(trimmedLine, "```") {
			inCodeBlock = false

			// 最初のコードブロックが終了したら処理を終了
			if len(usageLines) > 0 {
				break
			}
			continue
		}

		// コードブロック内の行を追加
		if inCodeBlock {
			if trimmedLine != "" {
				usageLines = append(usageLines, fmt.Sprintf("echo \\\"  %s\\\"", trimmedLine))
			}
		}
	}

	return usageLines, nil
}

// DefaultScriptGenerator はビルドスクリプトの生成を行う実装です
type DefaultScriptGenerator struct{}

// GenerateContent はビルドスクリプトの内容を生成します
func (g *DefaultScriptGenerator) GenerateContent(packageName, packagePath string, usageExamples []string) string {
	outputName := strings.ToLower(strings.ReplaceAll(packageName, "-", "_"))

	// 使用例がない場合のデフォルト
	usageExamplesStr := ""
	if len(usageExamples) > 0 {
		usageExamplesStr = strings.Join(usageExamples, "\n")
	} else {
		usageExamplesStr = fmt.Sprintf("echo \\\"  ./%s [options]\\\"", outputName)
	}

	return fmt.Sprintf(`#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="%s"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="%s"
WIN_OUTPUT_NAME="${OUTPUT_NAME}.exe"

echo "Building ${OUTPUT_NAME}..."

# Linux/AMD64向けビルド
echo "Building for Linux/AMD64..."
mkdir -p "${LINUX_AMD64_DIR}"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${LINUX_AMD64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

# Windows/AMD64向けビルド
echo "Building for Windows/AMD64..."
mkdir -p "${WIN_AMD64_DIR}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}" "${PACKAGE}"

# macOS/ARM64向けビルド
echo "Building for macOS/ARM64..."
mkdir -p "${MAC_ARM64_DIR}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o "${MAC_ARM64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

echo "Build completed successfully!"
echo "Usage as example:"
%s
`, packagePath, packageName, usageExamplesStr)
}
