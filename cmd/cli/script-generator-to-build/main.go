package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Config はアプリケーションの設定を保持する構造体
type Config struct {
	PackageName string
	ShowHelp    bool
}

// App はアプリケーションのメイン構造体
type App struct {
	Config *Config
}

// NewApp は新しいアプリケーションインスタンスを作成する
func NewApp(config *Config) *App {
	return &App{
		Config: config,
	}
}

// Run はアプリケーションのメイン処理を実行する
func (a *App) Run() int {
	// ヘルプオプションの確認
	if a.Config.ShowHelp {
		a.showHelp()
		return 0
	}

	var packageName string
	if a.Config.PackageName != "" {
		packageName = a.Config.PackageName
	} else {
		// パッケージ名が指定されていない場合、選択肢を表示
		var err error
		packageName, err = a.selectPackage()
		if err != nil {
			log.Printf("エラー: %v\n", err)
			return 1
		}
	}

	// ビルドスクリプトを生成
	if err := a.generateBuildScript(packageName); err != nil {
		log.Printf("エラー: %v\n", err)
		return 1
	}

	return 0
}

// showHelp はヘルプメッセージを表示する
func (a *App) showHelp() {
	fmt.Println("使用方法: script-generator-to-build [パッケージ名]")
	fmt.Println("")
	fmt.Println("このツールは、指定されたGoパッケージのビルドスクリプトを生成します。")
	fmt.Println("パッケージ名が指定されない場合は、利用可能なパッケージの一覧から選択できます。")
	fmt.Println("")
	fmt.Println("例:")
	fmt.Println("  script-generator-to-build code-analyzer")
	fmt.Println("  script-generator-to-build")
	fmt.Println("")
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
func (a *App) selectPackage() (string, error) {
	packages, err := a.getAvailablePackages()
	if err != nil {
		return "", err
	}

	if len(packages) == 0 {
		return "", fmt.Errorf("利用可能なパッケージが見つかりません")
	}

	fmt.Println("利用可能なパッケージ:")
	for i, pkg := range packages {
		fmt.Printf("  %d. %s\n", i+1, pkg)
	}

	// ユーザーに選択してもらう
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("ビルドスクリプトを生成するパッケージの番号を入力してください (1-%d): ", len(packages))
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("入力の読み取りに失敗しました: %v", err)
		}

		input = strings.TrimSpace(input)
		selection, err := strconv.Atoi(input)
		if err != nil || selection < 1 || selection > len(packages) {
			fmt.Printf("無効な選択です。1から%dまでの数字を入力してください。\n", len(packages))
			continue
		}

		return packages[selection-1], nil
	}
}

// generateBuildScript はビルドスクリプトを生成する
func (a *App) generateBuildScript(packageName string) error {
	packagePath := fmt.Sprintf("cmd/cli/%s", packageName)
	outputName := strings.ToLower(packageName)
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
			lines := strings.Split(string(content), "\n")
			inUsageSection := false
			usageLines := []string{}

			for _, line := range lines {
				if strings.Contains(line, "Usage:") || strings.Contains(line, "Example:") || strings.Contains(line, "Examples:") {
					inUsageSection = true
					continue
				}

				if inUsageSection {
					if strings.HasPrefix(line, "#") || line == "" {
						// 次のセクションに到達したか、空行の場合は終了
						if len(usageLines) > 0 {
							break
						}
					} else {
						// 使用例の行を追加
						trimmedLine := strings.TrimSpace(line)
						if trimmedLine != "" && !strings.HasPrefix(trimmedLine, "-") {
							usageLines = append(usageLines, fmt.Sprintf("echo \"  %s\"", trimmedLine))
						}
					}
				}
			}

			if len(usageLines) > 0 {
				usageExamples = strings.Join(usageLines, "\n")
			}
		}
	}

	// 使用例がない場合のデフォルト
	if usageExamples == "" {
		usageExamples = fmt.Sprintf("echo \"  ./%s [options]\"", outputName)
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

	fmt.Printf("ビルドスクリプトを生成しました: %s\n", scriptPath)
	return nil
}

func main() {
	// コマンドライン引数の解析
	var config Config
	flag.BoolVar(&config.ShowHelp, "h", false, "ヘルプメッセージを表示")
	flag.BoolVar(&config.ShowHelp, "help", false, "ヘルプメッセージを表示")
	flag.Parse()

	// 残りの引数があれば、最初の引数をパッケージ名として扱う
	args := flag.Args()
	if len(args) > 0 {
		config.PackageName = args[0]
	}

	// アプリケーションを実行
	app := NewApp(&config)
	exitCode := app.Run()

	// 終了コードを設定
	os.Exit(exitCode)
}
