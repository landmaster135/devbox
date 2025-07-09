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

	// スクリーンショットファイルのリネーム処理を統合メソッドで実行
	_, errorCount, err := usecases.ProcessScreenshotRename(config, stdout, stderr)
	if err != nil {
		if errorCount > 0 {
			return exitCodeError
		}
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
