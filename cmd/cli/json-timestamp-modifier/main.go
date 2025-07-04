package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/independencies/json_timestamp_modifier/interfaces/repositories"
	"github.com/landmaster135/devbox/internal/independencies/json_timestamp_modifier/usecases/services"
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
	dirPath := fs.String("dir", "", "操作するJSONファイルが含まれるディレクトリのパス")
	recursive := fs.Bool("recursive", false, "ディレクトリを再帰的に処理する（-dirオプションと共に使用）")
	key := fs.String("key", "timestamp", "操作するキー")
	mode := fs.String("mode", modeAddTimestamp, "操作モード: add (タイムスタンプ追加), to-unix (ISO-8601→UNIX), to-iso (UNIX→ISO-8601)")
	isJst := fs.Bool("is-jst", false, "日付のみの文字列をJSTとして扱う（to-unixモードのみ）")

	// 引数の解析
	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	// ファイルパスとディレクトリパスの両方が指定されていない場合はエラー
	if *filePath == "" && *dirPath == "" {
		fmt.Fprintln(stderr, "エラー: JSONファイルのパス（-file）またはディレクトリのパス（-dir）を指定してください")
		fs.Usage()
		return exitCodeError
	}

	// ファイルパスとディレクトリパスの両方が指定されている場合はエラー
	if *filePath != "" && *dirPath != "" {
		fmt.Fprintln(stderr, "エラー: -fileと-dirオプションは同時に指定できません")
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

	// 単一ファイルモードとディレクトリモードで処理を分岐
	if *filePath != "" {
		// 単一ファイルモード
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
	} else {
		// ディレクトリモード
		var count int
		var err error

		// ディレクトリが存在するか確認
		fileInfo, err := os.Stat(*dirPath)
		if err != nil {
			fmt.Fprintf(stderr, "エラー: ディレクトリ '%s' にアクセスできません: %v\n", *dirPath, err)
			return exitCodeError
		}
		if !fileInfo.IsDir() {
			fmt.Fprintf(stderr, "エラー: '%s' はディレクトリではありません\n", *dirPath)
			return exitCodeError
		}

		// モードに応じた処理を実行
		switch *mode {
		case modeAddTimestamp:
			// ディレクトリ内の全てのJSONファイルにタイムスタンプを追加
			count, err = timestampService.AddTimestampToAllFiles(*dirPath, *key, *recursive)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}
			fmt.Fprintf(stdout, "ディレクトリ '%s' 内の %d 個のJSONファイルにキー '%s' と現在の日時のタイムスタンプを追加しました\n", *dirPath, count, *key)

		case modeToUnix:
			// ディレクトリ内の全てのJSONファイルのISO-8601形式をUNIXタイムスタンプに変換
			count, err = jsonService.ConvertISO8601ToUnixInAllFiles(*dirPath, *key, *isJst, *recursive)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}
			fmt.Fprintf(stdout, "ディレクトリ '%s' 内の %d 個のJSONファイルのキー '%s' の値をISO-8601形式からUNIXタイムスタンプに変換しました\n", *dirPath, count, *key)

		case modeToISO8601:
			// ディレクトリ内の全てのJSONファイルのUNIXタイムスタンプをISO-8601形式に変換
			count, err = jsonService.ConvertUnixToISO8601InAllFiles(*dirPath, *key, *recursive)
			if err != nil {
				fmt.Fprintf(stderr, "エラー: %v\n", err)
				return exitCodeError
			}
			fmt.Fprintf(stdout, "ディレクトリ '%s' 内の %d 個のJSONファイルのキー '%s' の値をUNIXタイムスタンプからISO-8601形式に変換しました\n", *dirPath, count, *key)
		}
	}

	return exitCodeOK
}

func main() {
	// run関数を呼び出し、結果に応じて終了コードを設定
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
