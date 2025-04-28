package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/usecases/services"
)

// モード定数
const (
	modeAddTimestamp = "add"
	modeToUnix       = "to-unix"
	modeToISO8601    = "to-iso"
)

// exitCode はプログラムの終了コードを表します
type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

// run はJSON Timestamp Modifierツールの主要なロジックを実行します
func run(args []string, stdout, stderr io.Writer) exitCode {
	// フラグセットを作成
	fs := flag.NewFlagSet("json-timestamp-modifier", flag.ContinueOnError)
	fs.SetOutput(stderr)

	// コマンドライン引数の定義
	filePath := fs.String("file", "", "操作するJSONファイルのパス")
	key := fs.String("key", "timestamp", "操作するキー")
	mode := fs.String("mode", modeAddTimestamp, "操作モード: add (タイムスタンプ追加), to-unix (ISO-8601→UNIX), to-iso (UNIX→ISO-8601)")
	isJst := fs.Bool("is-jst", false, "日付のみの文字列をJSTとして扱う（to-unixモードのみ）")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// ファイルパスが指定されていない場合はエラー
	if *filePath == "" {
		fmt.Fprintln(stderr, "エラー: JSONファイルのパスを指定してください（-file オプション）")
		fs.Usage()
		return exitCodeError
	}

	// モードの検証
	if *mode != modeAddTimestamp && *mode != modeToUnix && *mode != modeToISO8601 {
		fmt.Fprintf(stderr, "エラー: 無効なモード '%s' が指定されました\n", *mode)
		fs.Usage()
		return exitCodeError
	}

	// 依存関係の注入
	fileRepo := repositories.NewFileRepository()
	jsonService := services.NewJSONService(fileRepo)
	timestampService := services.NewTimestampService(jsonService)

	var err error

	// モードに応じた処理を実行
	switch *mode {
	case modeAddTimestamp:
		// タイムスタンプを追加
		err = timestampService.AddTimestamp(*filePath, *key)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintf(stdout, "JSONファイル '%s' にキー '%s' と現在の日時のタイムスタンプを追加しました\n", *filePath, *key)

	case modeToUnix:
		// ISO-8601形式からUNIXタイムスタンプに変換
		err = jsonService.ConvertISO8601ToUnix(*filePath, *key, *isJst)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintf(stdout, "JSONファイル '%s' のキー '%s' の値をISO-8601形式からUNIXタイムスタンプに変換しました\n", *filePath, *key)

	case modeToISO8601:
		// UNIXタイムスタンプからISO-8601形式に変換
		err = jsonService.ConvertUnixToISO8601(*filePath, *key)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintf(stdout, "JSONファイル '%s' のキー '%s' の値をUNIXタイムスタンプからISO-8601形式に変換しました\n", *filePath, *key)
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
