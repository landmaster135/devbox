# ドキュメント作業ガイド

このドキュメントは、`cmd/` 配下のツールを実装または改修した後に必要なドキュメント更新作業を定義します。

## 対象

- `cmd/` 配下の新規ツールを追加したとき
- `cmd/` 配下の既存ツールの仕様（フラグ、出力、挙動）を変更したとき

## cmd配下ツール実装・改修後の手順

1. `cmd/cli/<tool>/README.md` を更新する（CLIツールの場合）
- 必須で追記する内容
  - ツール概要
  - フラグ一覧（必須/任意、デフォルト値）
  - 使用方法
  - 使用例
  - 出力例（成功時/エラー時）
- 参考:
  - `cmd/cli/arithmetic-calculator/README.md`
  - `cmd/cli/service-implementing-viewer/README.md`

2. 実装ガイドを必要に応じて更新する
- 新しい実装パターンや注意点を追加した場合:
  - `docs/tool_implementation/implementation_guide.md` を更新する

3. 最終確認を行う
- 変更したMarkdownの見出し構造・リンク切れを確認する
- 実装とREADMEの操作例が一致していることを確認する
- ドキュメントの配置/目次管理が必要な場合は `docs/docs_management/index.md` を参照する

## チェックリスト

- CLIツールの場合、`cmd/cli/<tool>/README.md` を更新した
- `docs/tool_implementation/implementation_guide.md` の更新要否を確認した
