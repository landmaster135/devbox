package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/file_maneuver/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はfile-maneuverツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("file-maneuver", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドライン引数の定義
	srcDirsStr := flagSet.String("src-dirs", "", "source directories to scan (comma-separated)")
	extensionsStr := flagSet.String("extensions", "", "target file extensions (comma-separated)")
	destDir := flagSet.String("dest-dir", "", "destination directory")
	recursive := flagSet.Bool("recursive", false, "recursively scan sub-directories")
	workers := flagSet.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	dryRun := flagSet.Bool("dry-run", false, "show what would be moved without actually moving files")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// 必須パラメータのチェック
	if *srcDirsStr == "" {
		fmt.Fprintln(stderr, "エラー: --src-dirs パラメータが必要です")
		flagSet.Usage()
		return exitCodeError
	}

	if *extensionsStr == "" {
		fmt.Fprintln(stderr, "エラー: --extensions パラメータが必要です")
		flagSet.Usage()
		return exitCodeError
	}

	if *destDir == "" {
		fmt.Fprintln(stderr, "エラー: --dest-dir パラメータが必要です")
		flagSet.Usage()
		return exitCodeError
	}

	// カンマ区切り文字列を配列に変換
	srcDirs := parseCommaSeparatedString(*srcDirsStr)
	extensions := parseCommaSeparatedString(*extensionsStr)

	// 設定の作成（この時点で全バリデーション完了）
	config, err := usecases.NewConfig(
		srcDirs,
		extensions,
		*destDir,
		*recursive,
		*workers,
		*dryRun,
	)
	if err != nil {
		fmt.Fprintf(stderr, "設定エラー: %v\n", err)
		return exitCodeError
	}

	// サービスの作成（configは既にバリデーション済み）
	service := usecases.NewFileManeuverService(config)

	// 処理開始の通知
	fmt.Fprintf(stdout, "✨ ファイル移動処理を開始します\n")
	fmt.Fprintf(stdout, "  ソースディレクトリ: %s\n", strings.Join(srcDirs, ", "))
	fmt.Fprintf(stdout, "  対象拡張子: %s\n", strings.Join(extensions, ", "))
	fmt.Fprintf(stdout, "  宛先ディレクトリ: %s\n", *destDir)
	fmt.Fprintf(stdout, "  再帰的検索: %t\n", *recursive)
	fmt.Fprintf(stdout, "  ワーカー数: %d\n", *workers)
	if *dryRun {
		fmt.Fprintf(stdout, "  モード: ドライラン（実際の移動は行いません）\n")
	}
	fmt.Fprintln(stdout)

	// ファイル移動処理を一括実行
	successCount, errorCount, err := service.ExecuteFileManeuver(stdout, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "処理エラー: %v\n", err)
		return exitCodeError
	}

	// 処理結果の出力
	fmt.Fprintln(stdout)
	fmt.Fprintf(stdout, "✅ ファイル移動処理が完了しました\n")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", successCount)

	if errorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", errorCount)
		return exitCodeError
	}

	return exitCodeOK
}

// parseCommaSeparatedString はカンマ区切り文字列を配列に変換します
func parseCommaSeparatedString(str string) []string {
	if str == "" {
		return []string{}
	}

	parts := strings.Split(str, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
