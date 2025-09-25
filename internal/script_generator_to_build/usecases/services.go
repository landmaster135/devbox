package usecases

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/landmaster135/devbox/internal/script_generator_to_build/config"
)

// 終了コード
const (
	ExitCodeOK = iota
	ExitCodeError
)

// Service はサービスの主要なロジックを表します
type Service struct {
	Config          *config.ServiceConfig
	FileSystem      FileSystem
	InputReader     InputReader
	READMEParser    READMEParser
	ScriptGenerator ScriptGenerator
	Logger          *log.Logger
}

// NewService は新しい Service インスタンスを作成します（後方互換性のため）
func NewService(cfg *config.ServiceConfig) *Service {
	return NewServiceWithDependencies(cfg, nil, nil, nil, nil)
}

// NewServiceWithDependencies は依存性注入を使用して新しい Service インスタンスを作成します
func NewServiceWithDependencies(
	cfg *config.ServiceConfig,
	fs FileSystem,
	reader InputReader,
	parser READMEParser,
	generator ScriptGenerator,
) *Service {
	// デフォルト値を設定
	cfg.SetDefaults()

	// デフォルト実装を注入
	if fs == nil {
		fs = &OSFileSystem{}
	}
	if reader == nil {
		reader = NewStdinReader()
	}
	if parser == nil {
		parser = &DefaultREADMEParser{}
	}
	if generator == nil {
		generator = &DefaultScriptGenerator{}
	}

	return &Service{
		Config:          cfg,
		FileSystem:      fs,
		InputReader:     reader,
		READMEParser:    parser,
		ScriptGenerator: generator,
		Logger:          log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Run はアプリケーションを実行します
func (s *Service) Run(stdout, stderr io.Writer) int {
	// ログの出力先を設定
	log.SetOutput(stderr)

	// ヘルプオプションの確認
	if s.Config.ShowHelp {
		s.showHelp(stdout)
		return ExitCodeOK
	}

	var packageName string
	if s.Config.PackageName != "" {
		packageName = s.Config.PackageName
	} else {
		// パッケージ名が指定されていない場合、選択肢を表示
		var err error
		packageName, err = s.selectPackage(stdout)
		if err != nil {
			log.Printf("エラー: %v\n", err)
			return ExitCodeError
		}
	}

	// ビルドスクリプトを生成
	if err := s.generateBuildScript(packageName, stdout); err != nil {
		log.Printf("エラー: %v\n", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

// showHelp はヘルプメッセージを表示する
func (s *Service) showHelp(w io.Writer) {
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
func (s *Service) getAvailablePackages() ([]string, error) {
	// CLIディレクトリ内のサブディレクトリを検索
	entries, err := s.FileSystem.ReadDir(s.Config.GetCLIPath())
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
func (s *Service) selectPackage(w io.Writer) (string, error) {
	packages, err := s.getAvailablePackages()
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
	for {
		fmt.Fprintf(w, "ビルドスクリプトを生成するパッケージの番号を入力してください (1-%d): ", len(packages))
		input, err := s.InputReader.ReadString('\n')
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

// validatePackage はパッケージの存在を確認します
func (s *Service) validatePackage(packageName string) error {
	packagePath := fmt.Sprintf("%s/%s", s.Config.CLIDir, packageName)
	fullPath := fmt.Sprintf("%s/%s", s.Config.BaseDir, packagePath)

	if _, err := s.FileSystem.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("パッケージ '%s' が見つかりません", packageName)
	}
	return nil
}

// parseREADMEFile はREADMEファイルから使用例を解析します
func (s *Service) parseREADMEFile(packageName string, w io.Writer) ([]string, error) {
	packagePath := fmt.Sprintf("%s/%s", s.Config.CLIDir, packageName)
	readmePath := fmt.Sprintf("%s/%s/README.md", s.Config.BaseDir, packagePath)

	if _, err := s.FileSystem.Stat(readmePath); os.IsNotExist(err) {
		// READMEファイルが存在しない場合は空のスライスを返す
		return []string{}, nil
	}

	content, err := s.FileSystem.ReadFile(readmePath)
	if err != nil {
		return []string{}, fmt.Errorf("READMEファイルの読み取りに失敗しました: %v", err)
	}

	fmt.Fprintf(w, "READMEファイルを読み込みました: %s\n", readmePath)

	usageExamples, err := s.READMEParser.ParseUsageExamples(content)
	if err != nil {
		return []string{}, fmt.Errorf("使用例の解析に失敗しました: %v", err)
	}

	if len(usageExamples) > 0 {
		fmt.Fprintf(w, "使用例を抽出しました（%d行）\n", len(usageExamples))
	} else {
		fmt.Fprintln(w, "使用例を抽出できませんでした")
	}

	return usageExamples, nil
}

// writeScriptFile はスクリプトファイルを書き込みます
func (s *Service) writeScriptFile(packageName, content string, w io.Writer) error {
	outputName := strings.ToLower(strings.ReplaceAll(packageName, "-", "_"))
	scriptPath := fmt.Sprintf("%s/build_%s.sh", s.Config.GetScriptsPath(), outputName)

	// スクリプトディレクトリを作成
	if err := s.FileSystem.MkdirAll(s.Config.GetScriptsPath(), 0755); err != nil {
		return fmt.Errorf("スクリプトディレクトリの作成に失敗しました: %v", err)
	}

	// ファイルに書き込み
	if err := s.FileSystem.WriteFile(scriptPath, []byte(content), 0755); err != nil {
		return fmt.Errorf("ビルドスクリプトの書き込みに失敗しました: %v", err)
	}

	fmt.Fprintf(w, "ビルドスクリプトを生成しました: %s\n", scriptPath)
	return nil
}

// generateBuildScript はビルドスクリプトを生成する
func (s *Service) generateBuildScript(packageName string, w io.Writer) error {
	// パッケージの存在確認
	if err := s.validatePackage(packageName); err != nil {
		return err
	}

	// READMEファイルから使用例を解析
	usageExamples, err := s.parseREADMEFile(packageName, w)
	if err != nil {
		return err
	}

	// スクリプト内容を生成
	packagePath := fmt.Sprintf("%s/%s", s.Config.CLIDir, packageName)
	scriptContent := s.ScriptGenerator.GenerateContent(packageName, packagePath, usageExamples)

	// スクリプトファイルを書き込み
	return s.writeScriptFile(packageName, scriptContent, w)
}
