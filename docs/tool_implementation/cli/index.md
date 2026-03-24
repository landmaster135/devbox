# CLIツール実装ガイド（入口）

このドキュメントは、CLIツール実装の入口です。

## CLI固有の実装手順

1. `cmd/cli/{tool_name}/main.go` を作成し、エントリーポイントを配置する
2. フラグ解析（`config.ParseFlags()`）と操作分岐（`switch cfg.Operation`）を実装する
3. CLI から `usecases` を呼び出す構成にする
4. 正常系の結果は標準出力（`fmt.Print`）に出力する
5. 異常系は標準エラー出力（`fmt.Fprintf(os.Stderr, ...)`）と `os.Exit(1)` で終了する
6. `-help` などの利用方法表示は `infrastructures/flag_parser.PrintUsage()` として実装し、`config` 層から分離する

## フラグ実装の責務分離

1. `config` 層は `Config` の生成と検証を担当する
2. `infrastructures/flag_parser` は `FlagParser` 実装と usage 出力を担当する
3. `MockFlagParser` は `infrastructures/flag_parser` に colocate し、`config` テストから再利用する

## 実装パターン

実装パターンとアンチパターンの詳細は `docs/tool_implementation/cli/guide.md` を参照してください。
