# AGENTS.md (docs マップ)
Quick Decision Trees に則って、参照するべきドキュメントを全て確認すること。

## Summary

このファイルは、エージェント向けの入口です。  
実装・設計・運用の詳細ルールは `docs/` 側を正本とします。

- `docs/tool_implementation/index.md`
  - `cmd/` 配下ツールの実装・改修に関する統合入口
- `docs/docs_management/index.md`
  - docs配置方針とリンク整合性の管理
- `docs/exec_plans/index.md`
  - 実行計画ドキュメント（active/completed）の入口
- `docs/project_status/service_implementation_status.md`
  - サービス実装状況（CLI/MCP/gRPC/HTTP）の確認
- `docs/project_overview/project_overview.md`
  - 全体アーキテクチャとディレクトリ構成の確認
- `docs/user_prompt/prompt_sample.md`
  - 依頼テンプレートやプロンプト例の参照

## Prerequisites

- 回答は日本語で行うこと。
- 詳細ルールをこのファイルへ増やさないこと。詳細は `docs/` 側へ追記すること。

## Quick Decision Trees

### 「タスクに着手したい」

```text
タスクに着手したい
└─ タスクの指示書を確認したい → `.agents/tmp/instructions.md`
```

### 「cmd配下ツールを実装・改修したい」

```text
cmd配下ツールを実装・改修したい
├─ 実装・改修方針/ビルド/運用の全体像を見たい → docs/tool_implementation/index.md
├─ 実装・改修計画を立てたい → docs/tool_implementation/index.md (実装計画は`.agents/tmp/draft.md`に記載する)
├─ 実装・改修パターンや注意点を確認したい → docs/tool_implementation/index.md
└─ 実装・改修後のドキュメント更新手順を確認したい → docs/tool_implementation/documentation_guide.md
```

### 「ドキュメントを追加・移動・整理したい」

```text
ドキュメントを追加・移動・整理したい
└─ docs管理ルールとリンク整合性を確認したい → docs/docs_management/index.md
```

### 「Taskfileのタスクを実装・改修したい」

```text
Taskfileのタスクを実装・改修したい
└─ 実装フローを確認したい → docs/task_implementation/index.md
```

### 「プロジェクトの状況を確認したい」
```text
プロジェクトの状況を確認したい
├─ 実装済みサービス一覧（CLI/MCP/gRPC/HTTP）を確認したい → docs/project_status/service_implementation_status.md
├─ 実装済みサービスの概要を確認したい → docs/project_status/service_overview.md
├─ 全体のアーキテクチャ/構成を確認したい → docs/project_overview/project_overview.md
├─ 計画（active/completed）を確認したい → docs/exec_plans/index.md
└─ インストール可能なエントリポイント一覧を確認したい → docs/project_status/entrypoint_overview.md
```

### 「依頼文の例を見たい」

```text
依頼文の例を見たい
└─ プロンプトサンプルを参照したい → docs/user_prompt/prompt_sample.md
```
