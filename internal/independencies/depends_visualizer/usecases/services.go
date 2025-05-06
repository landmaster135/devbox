package usecases

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/analyzer"
	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/config"
	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/finder"
	"github.com/landmaster135/devbox/internal/independencies/depends_visualizer/presenter"
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
func (app *App) Run(stdout, stderr io.Writer) int {
	// ログの出力先を設定
	log.SetOutput(stderr)

	// 拡張子の処理
	if app.Config.Extension != "" && app.Config.Extension[0] != '.' {
		app.Config.Extension = "." + app.Config.Extension
	}

	// 設定ファイルの読み込み（もし指定されていれば）
	if app.Config.ConfigPath != "" {
		if err := config.LoadFromFile(app.Config.ConfigPath); err != nil {
			log.Printf("設定ファイルの読み込みに失敗しました: %v", err)
			return ExitCodeError
		}
	}

	// ファイル一覧を取得
	files, err := app.collectFiles()
	if err != nil {
		log.Printf("ファイルの収集に失敗しました: %v", err)
		return ExitCodeError
	}

	if len(files) == 0 {
		log.Printf("処理対象のファイルが見つかりません")
		return ExitCodeError
	}

	log.Printf("処理対象ファイル: %d ファイル", len(files))

	// ファイルを解析
	results, err := app.analyzeFiles(files)
	if err != nil {
		log.Printf("解析に失敗しました: %v", err)
		return ExitCodeError
	}

	// 結果を出力
	if err := app.renderResults(results, stdout); err != nil {
		log.Printf("出力の生成に失敗しました: %v", err)
		return ExitCodeError
	}

	return ExitCodeOK
}

// ファイル一覧を収集
func (app *App) collectFiles() ([]string, error) {
	var files []string

	if app.Config.SourceFile != "" {
		// 単一ファイルモード
		if _, err := os.Stat(app.Config.SourceFile); err != nil {
			return nil, fmt.Errorf("指定されたファイルが見つかりません: %v", err)
		}
		files = append(files, app.Config.SourceFile)

		// 拡張子が指定されていない場合は自動検出
		if app.Config.Extension == "" {
			app.Config.Extension = filepath.Ext(app.Config.SourceFile)
		}
	} else {
		// ディレクトリ処理モード
		extensions := []string{app.Config.Extension}
		if app.Config.Extension == "" {
			// 拡張子が指定されていない場合はデフォルトのすべてをサポート
			extensions = getSupportedExtensions()
		}

		var err error
		files, err = finder.FindFiles(app.Config.Directory, app.Config.Recursive, extensions)
		if err != nil {
			return nil, err
		}

		// 拡張子が指定されていない場合は最初のファイルから取得
		if app.Config.Extension == "" && len(files) > 0 {
			app.Config.Extension = filepath.Ext(files[0])
		}
	}

	// 拡張子が対応しているか確認
	if !isSupportedExtension(app.Config.Extension) {
		return nil, fmt.Errorf("サポートされていない拡張子です: %s", app.Config.Extension)
	}

	return files, nil
}

// ファイルを解析
func (app *App) analyzeFiles(files []string) ([]analyzer.AnalysisResult, error) {
	results := make([]analyzer.AnalysisResult, 0, len(files))
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// ワーカー数を設定
	workerCount := runtime.NumCPU()
	jobs := make(chan string, len(files))

	// エラーカウント
	errCount := 0

	// ワーカーを起動
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range jobs {
				// 単一ファイルを解析
				deps, err := analyzer.AnalyzeFile(file, app.Config.Extension)
				if err != nil {
					log.Printf("警告: ファイル %s の解析に失敗: %v", file, err)
					mutex.Lock()
					errCount++
					mutex.Unlock()
					continue
				}

				// 結果を保存
				mutex.Lock()
				results = append(results, analyzer.AnalysisResult{
					FilePath:     file,
					Dependencies: deps,
				})
				mutex.Unlock()

				log.Printf("情報: ファイル %s の解析が完了 (%d 関数)", file, len(deps))
			}
		}()
	}

	// ジョブをキューに追加
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	if len(results) == 0 {
		return nil, fmt.Errorf("すべてのファイルの解析に失敗しました")
	}

	return results, nil
}

// 結果を出力
func (app *App) renderResults(results []analyzer.AnalysisResult, w io.Writer) error {
	var outputStr string
	var err error

	// 出力形式に応じたプレゼンターを選択
	switch app.Config.Format {
	case "mermaid":
		outputStr, err = presenter.RenderMermaid(results)
	case "mermaid-flowchart":
		outputStr, err = presenter.RenderMermaidFlowchart(results)
	case "plantuml":
		outputStr, err = presenter.RenderPlantUML(results)
	case "dot":
		outputStr, err = presenter.RenderDOT(results)
	default:
		return fmt.Errorf("未サポートの出力形式です: %s", app.Config.Format)
	}

	if err != nil {
		return err
	}

	// 出力先を決定
	if app.Config.OutputPath != "" {
		// ファイルに出力
		err = os.WriteFile(app.Config.OutputPath, []byte(outputStr), 0644)
		if err != nil {
			return fmt.Errorf("ファイルへの書き込みに失敗しました: %v", err)
		}
		log.Printf("情報: 結果を %s に保存しました", app.Config.OutputPath)
	} else {
		// 標準出力に出力
		fmt.Fprintln(w, outputStr)
	}

	return nil
}

// サポートされている拡張子のリストを取得
func getSupportedExtensions() []string {
	return []string{".go", ".py", ".js"}
}

// 拡張子がサポートされているか確認
func isSupportedExtension(ext string) bool {
	for _, e := range getSupportedExtensions() {
		if e == ext {
			return true
		}
	}
	return false
}
