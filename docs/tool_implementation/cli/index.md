# CLIツール実装ガイド（入口）

このドキュメントは、CLIツール実装の入口です。

## CLI固有の実装手順

1. `cmd/cli/{tool_name}/main.go` を作成し、エントリーポイントを配置する
2. フラグ解析（`config.ParseFlags()`）と操作分岐（`switch cfg.Operation`）を実装する
3. 正常系の結果は標準出力（`fmt.Print`）に出力する
4. 異常系は標準エラー出力（`fmt.Fprintf(os.Stderr, ...)`）と `os.Exit(1)` で終了する
5. `-help` などの利用方法表示（`config.PrintUsage()`）を実装する

## 実装パターン

実装パターンとアンチパターンの詳細は `docs/tool_implementation/cli/guide.md` を参照してください。
