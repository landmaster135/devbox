package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	usecases "github.com/landmaster135/devbox/internal/image_renamer_for_screenshot/usecases"
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
	flagSet := flag.NewFlagSet("image-renamer-for-screenshot", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	// コマンドラインフラグの定義
	srcDir := flagSet.String("src", ".", "スキャンするソースディレクトリ")
	vlcPattern := flagSet.Bool("vlc", false, "VLCスナップショットファイル (vlcsnap-*.png) をリネーム")
	winPattern := flagSet.Bool("win", false, "Windowsスクリーンショットファイル (スクリーンショット *.png) をリネーム")
	androidPattern := flagSet.Bool("android", false, "Androidスクリーンレコードファイル (screen-*.mp4) をリネーム")
	toDateTime := flagSet.Bool("to-datetime", false, "ファイル名をYYYYMMDDHHMMSS形式にリネーム")
	recursive := flagSet.Bool("r", false, "サブディレクトリを再帰的にスキャン")
	workers := flagSet.Int("workers", runtime.NumCPU(), "並行ワーカー数")

	// 引数の解析
	if err := flagSet.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return usecases.Config{}, err
	}

	return usecases.Config{
		SrcDir:         *srcDir,
		Recursive:      *recursive,
		Workers:        *workers,
		VlcPattern:     *vlcPattern,
		WinPattern:     *winPattern,
		AndroidPattern: *androidPattern,
		ToDateTime:     *toDateTime,
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

	// スクリーンショットファイルの検索
	var files []string
	if config.ToDateTime {
		files, err = usecases.FindScreenshotFilesForDateTime(config.SrcDir, config.Recursive, stdout, stderr)
	} else {
		files, err = usecases.FindScreenshotFiles(config.SrcDir, config.Recursive, config.VlcPattern, config.WinPattern, config.AndroidPattern, stdout, stderr)
	}
	if err != nil {
		return exitCodeError
	}

	if len(files) == 0 {
		fmt.Fprintln(stdout, "スクリーンショットファイルが見つかりませんでした。")
		return exitCodeOK
	}

	fmt.Fprintf(stdout, "スクリーンショットファイルが %d 件見つかりました。\n", len(files))

	if config.ToDateTime {
		fmt.Fprintln(stdout, "YYYYMMDDHHMMSS形式でのリネームを実行します。")
	} else if config.VlcPattern {
		fmt.Fprintln(stdout, "VLCスナップショットパターンを使用します。")
	} else if config.WinPattern {
		fmt.Fprintln(stdout, "Windowsスクリーンショットパターンを使用します。")
	} else if config.AndroidPattern {
		fmt.Fprintln(stdout, "Androidスクリーンレコードパターンを使用します。")
	}

	// ファイル情報の取得
	fileInfos, err := usecases.GetFileInfos(files, stderr)
	if err != nil {
		// エラーがあっても続行するため、ここではエラーコードを返さない
	}

	// リネーム処理の実行
	successCount, errorCount := usecases.RenameScreenshotFiles(fileInfos, config, stdout, stderr)

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
