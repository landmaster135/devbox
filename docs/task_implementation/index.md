# タスク実装ガイド

このディレクトリは、Taskfile 実装手順を用途別に管理します。

## 目次

- 開発用 Taskfile 実装フロー: `docs/task_implementation/dev/index.md`
- 配布用 Taskfile 実装フロー: `docs/task_implementation/pkg/index.md`

## 共通チェック

- 同名タスクが既に存在しない
- `desc` が実行内容を説明している
- 相対パスが Taskfile の位置から見て正しい
- 追加タスクが既存命名規則と整合している
- 新規ドキュメント追加時は `docs/docs_management/index.md` の運用ルールに従う
