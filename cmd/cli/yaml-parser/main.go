package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	yamlconfig "github.com/landmaster135/devbox/internal/yaml_parser/config"
	"github.com/landmaster135/devbox/internal/yaml_parser/usecase"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func run(args []string, stdout, stderr io.Writer) exitCode {
	fs := flag.NewFlagSet("yaml-parser", flag.ContinueOnError)
	fs.SetOutput(stderr)

	operation := fs.String("operation", "", "実行する操作 (read|parse|edit-file)")
	filePath := fs.String("file-path", "", "YAMLファイルのパス (--operation=read|edit-file)")
	yamlContent := fs.String("yaml-content", "", "解析するYAML文字列 (--operation=parse)")
	keyValueList := fs.String("key-value-list", "", "key=value 形式をカンマ/改行区切りで列挙 (--operation=edit-file)")

	fs.Usage = func() {
		fmt.Fprintln(stderr, "YAML Parser CLI")
		fmt.Fprintln(stderr, "使用方法:")
		fmt.Fprintln(stderr, "  yaml-parser --operation read --file-path ./config.yaml")
		fmt.Fprintln(stderr, "  yaml-parser --operation parse --yaml-content \"key: value\"")
		fmt.Fprintln(stderr, "  yaml-parser --operation edit-file --file-path ./config.yaml --key-value-list \"server.port=8081\"")
		fmt.Fprintln(stderr)
		fmt.Fprintln(stderr, "フラグ一覧:")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return exitCodeError
	}

	cfg, err := yamlconfig.NewConfig(*operation, *filePath, *yamlContent, *keyValueList)
	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		fs.Usage()
		return exitCodeError
	}

	svc := usecase.NewService()

	var result any
	switch cfg.Operation {
	case yamlconfig.OperationRead:
		result, err = svc.ReadFromFile(cfg.FilePath)
	case yamlconfig.OperationParse:
		result, err = svc.ParseFromContent(cfg.YAMLContent)
	case yamlconfig.OperationEditFile:
		result, err = svc.EditFile(cfg.FilePath, cfg.KeyValueList)
	default:
		err = fmt.Errorf("未対応の操作です: %s", cfg.Operation)
	}

	if err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	if err := outputResult(stdout, result); err != nil {
		fmt.Fprintf(stderr, "エラー: %v\n", err)
		return exitCodeError
	}

	return exitCodeOK
}

func outputResult(w io.Writer, data any) error {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("結果の整形に失敗しました: %w", err)
	}
	if _, err := fmt.Fprintln(w, string(encoded)); err != nil {
		return fmt.Errorf("結果の出力に失敗しました: %w", err)
	}
	return nil
}

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
