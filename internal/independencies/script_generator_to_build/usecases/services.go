package usecases

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/independencies/script_generator_to_build/config"
)

// 終了コード
const (
	ExitCodeOK = iota
	ExitCodeError
)

// App はアプリケーションの主要なロジックを表します
type App struct {
	Config *config.AppConfig
}

// NewApp は新しい App インスタンスを作成します
func NewApp(cfg *config.AppConfig) *App {
	return &App{
		Config: cfg,
	}
}

// Run はアプリケーションを実行します
func (a *App) Run(stdout, stderr io.Writer) int {
	// ログの出力先を設定
	log.SetOutput(stderr)

	// ヘルプオプションの確認
	if a.Config.ShowHelp {
		a.showHelp(stdout)
		return ExitCodeOK
	}

	var packageName string
	if a.Config.PackageName != "" {
		packageName = a.Config.PackageName
	} else {
		// パッケージ名が指定されていない場合、選択肢を表示
		var err error
		packageName, err = a.selectPackage(stdout)
		if err != nil {
			log.Printf("エラー: %v\n", err)
			return ExitCodeError
		}
	}

	// ビルドスクリプトを生成
	if err := a.generateBuildScript(packageName, stdout); err != nil {
		log.Printf("エラー: %v\n", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

// showHelp はヘルプメッセージを表示する
func (a *App) showHelp(w io.Writer) {
	fmt.Fprintln(w, "使用方法: script-generator-to-build [パッケージ名]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "このツールは、指定されたGoパッケージのビルドスクリプトを生成します。")
	fmt.Fprintln(w, "パッケージ名が指定されない場合は、利用可能なパッケージの一覧から選択できます。")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "例:")
	fmt.Fprintln(w, "  script-generator-to-build code-analyzer")
	fmt.Fprintln(w, "  script-generator-to-build")
	fmt.Fprintln(w, "")
}

// getAvailablePackages は利用可能なパッケージのリストを取得する
func (a *App) getAvailablePackages() ([]string, error) {
	// /home/nov/devbox/cmd/cli ディレクトリ内のサブディレクトリを検索
	entries, err := os.ReadDir("/home/nov/devbox/cmd/cli")
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み取りに失敗しました: %v", err)
	}

	var packages []string
	for _, entry := range entries {
		if entry.IsDir() {
			packages = append(packages, entry.Name())
		}
	}

	// パッケージ名をアルファベット順にソート
	sort.Strings(packages)
	return packages, nil
}

// selectPackage はユーザーにパッケージを選択させる
func (a *App) selectPackage(w io.Writer) (string, error) {
	packages, err := a.getAvailablePackages()
	if err != nil {
		return "", err
	}

	if len(packages) == 0 {
		return "", fmt.Errorf("利用可能なパッケージが見つかりません")
	}

	fmt.Fprintln(w, "利用可能なパッケージ:")
	for i, pkg := range packages {
		fmt.Fprintf(w, "  %d. %s\n", i+1, pkg)
	}

	// ユーザーに選択してもらう
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(w, "ビルドスクリプトを生成するパッケージの番号を入力してください (1-%d): ", len(packages))
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("入力の読み取りに失敗しました: %v", err)
		}

		input = strings.TrimSpace(input)
		selection, err := strconv.Atoi(input)
		if err != nil || selection < 1 || selection > len(packages) {
			fmt.Fprintf(w, "無効な選択です。1から%dまでの数字を入力してください。\n", len(packages))
			continue
		}

		return packages[selection-1], nil
	}
}

// generateBuildScript はビルドスクリプトを生成する
func (a *App) generateBuildScript(packageName string, w io.Writer) error {
	packagePath := fmt.Sprintf("cmd/cli/%s", packageName)
	outputName := strings.ToLower(strings.ReplaceAll(packageName, "-", "_"))
	scriptPath := fmt.Sprintf("/home/nov/devbox/scripts/build_%s.sh", outputName)

	// パッケージディレクトリが存在するか確認
	if _, err := os.Stat(fmt.Sprintf("/home/nov/devbox/%s", packagePath)); os.IsNotExist(err) {
		return fmt.Errorf("パッケージ '%s' が見つかりません", packageName)
	}

	// 使用例を取得するために、READMEファイルを確認
	usageExamples := ""
	readmePath := fmt.Sprintf("/home/nov/devbox/%s/README.md", packagePath)
	if _, err := os.Stat(readmePath); err == nil {
		// READMEファイルが存在する場合
		content, err := os.ReadFile(readmePath)
		if err == nil {
			// デバッグ用にREADMEの内容を出力
			fmt.Fprintf(w, "READMEファイルを読み込みました: %s\n", readmePath)

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
					fmt.Fprintf(w, "「使用例」セクションを見つけました: %s（行 %d）\n", line, i+1)
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
						fmt.Fprintf(w, "「使用方法」セクションを見つけました: %s（行 %d）\n", line, i+1)
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

			if startIndex >= 0 {
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
							fmt.Fprintf(w, "「使用例」サブセクションを見つけました: %s（行 %d）\n", line, i+1)
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
								fmt.Fprintf(w, "「使用例」セクション内のサブセクションを検出: %s（行 %d）\n", line, i+1)
								continue
							}
						}

						// それ以外の見出しなら終了
						fmt.Fprintf(w, "次のセクションに到達しました: %s（行 %d）\n", line, i+1)
						break
					}

					trimmedLine := strings.TrimSpace(line)

					// コードブロックの開始を検出
					if !inCodeBlock && strings.HasPrefix(trimmedLine, "```") {
						inCodeBlock = true
						fmt.Fprintf(w, "コードブロック開始を検出: %s（行 %d）\n", line, i+1)
						continue
					}

					// コードブロックの終了を検出
					if inCodeBlock && strings.HasPrefix(trimmedLine, "```") {
						inCodeBlock = false
						fmt.Fprintf(w, "コードブロック終了を検出: %s（行 %d）\n", line, i+1)

						// 最初のコードブロックが終了したら処理を終了
						if len(usageLines) > 0 {
							fmt.Fprintln(w, "最初のコードブロックを抽出しました。処理を終了します。")
							break
						}
						continue
					}

					// コードブロック内の行を追加
					if inCodeBlock {
						if trimmedLine != "" {
							usageLines = append(usageLines, fmt.Sprintf("echo \\\"  %s\\\"", trimmedLine))
							fmt.Fprintf(w, "使用例行を追加: %s\n", trimmedLine)
						}
					}
				}

				if len(usageLines) > 0 {
					usageExamples = strings.Join(usageLines, "\n")
					fmt.Fprintf(w, "使用例を抽出しました（%d行）\n", len(usageLines))
				} else {
					fmt.Fprintln(w, "使用例を抽出できませんでした")
				}
			} else {
				fmt.Fprintln(w, "使用例または使用方法セクションが見つかりませんでした")
			}
		}
	}

	// 使用例がない場合のデフォルト
	if usageExamples == "" {
		// パッケージ名のハイフンをアンダースコアに変換して出力名を生成
		exampleOutputName := strings.ToLower(strings.ReplaceAll(packageName, "-", "_"))
		usageExamples = fmt.Sprintf("echo \\\"  ./%s [options]\\\"", exampleOutputName)
	}

	// ビルドスクリプトの内容を生成
	scriptContent := fmt.Sprintf(`#!/bin/bash

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
`, packagePath, outputName, usageExamples)

	// ファイルに書き込み
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("ビルドスクリプトの書き込みに失敗しました: %v", err)
	}

	fmt.Fprintf(w, "ビルドスクリプトを生成しました: %s\n", scriptPath)
	return nil
}
