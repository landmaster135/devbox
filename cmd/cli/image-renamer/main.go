package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	usecases "github.com/landmaster135/devbox/internal/image_renamer/usecases"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// parseFlags はコマンドライン引数を解析し、設定を返します
func parseFlags(args []string, stderr io.Writer) (usecases.Config, error) {
	// フラグセットを作成
	flagSet := flag.NewFlagSet("image-renamer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドラインフラグの定義
	srcDir := flagSet.String("src", ".", "スキャンするソースディレクトリ")
	sortByName := flagSet.Bool("name", false, "画像ファイルをファイル名順に並べ替え")
	sortByTime := flagSet.Bool("time", false, "画像ファイルを更新日時順に並べ替え")
	prefix := flagSet.String("prefix", "", "記事番号のプレフィックス (必須)")
	delimiter := flagSet.String("delimiter", "_", "プレフィックスとシリアル番号の間の区切り文字")
	digits := flagSet.Int("digits", 4, "シリアル番号の桁数")
	startCount := flagSet.Int("start", 1, "リネーム操作の開始番号")
	recursive := flagSet.Bool("r", false, "サブディレクトリを再帰的にスキャン")
	workers := flagSet.Int("workers", runtime.NumCPU(), "並行ワーカー数")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return usecases.Config{}, err
	}

	return usecases.Config{
		SrcDir:     *srcDir,
		SortByName: *sortByName,
		SortByTime: *sortByTime,
		Prefix:     *prefix,
		Delimiter:  *delimiter,
		Digits:     *digits,
		StartCount: *startCount,
		Recursive:  *recursive,
		Workers:    *workers,
	}, nil
}

// run は画像ファイルリネームツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// コマンドライン引数の解析と設定の取得
	config, err := parseFlags(args, stderr)
	if err != nil {
		return exitCodeError
	}

	// 引数の検証
	if err := usecases.ValidateConfig(config, stderr); err != nil {
		return exitCodeError
	}

	// 画像ファイルの検索
	files, err := usecases.FindImageFiles(config.SrcDir, config.Recursive, stdout, stderr)
	if err != nil {
		return exitCodeError
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "画像ファイルが見つかりませんでした。")
		return exitCodeOK
	}

	fmt.Fprintf(stdout, "画像ファイルが %d 件見つかりました。\n", len(files))
	fmt.Fprintf(stdout, "プレフィックス: %s\n", config.Prefix)
	fmt.Fprintf(stdout, "区切り文字: %s\n", config.Delimiter)
	fmt.Fprintf(stdout, "開始番号: %d\n", config.StartCount)

	// ファイル情報の取得と並べ替え
	fileInfos, err := usecases.GetFileInfos(files, stderr)
	if err != nil {
		// エラーがあっても続行するため、ここではエラーコードを返さない
	}

	// ファイルの並べ替え
	usecases.SortFiles(fileInfos, config.SortByTime, stdout)

	// リネーム処理の実行
	successCount, errorCount := usecases.RenameFiles(fileInfos, config, stdout, stderr)

	// 処理結果の出力
	fmt.Fprintf(stdout, "✔ ファイルリネームが完了しました\n")
	fmt.Fprintf(stdout, "  成功: %d ファイル\n", successCount)

	if errorCount > 0 {
		fmt.Fprintf(stdout, "  失敗: %d ファイル\n", errorCount)
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
