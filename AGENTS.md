# AGENTS.md (docs マップ)

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

### 「cmd配下ツールを実装・改修したい」

```text
cmd配下ツールを実装・改修したい
├─ 実装方針/ビルド/運用の全体像を見たい → docs/tool_implementation/index.md
├─ 実装パターンや注意点を確認したい → docs/tool_implementation/implementation_guide.md
└─ 実装後のドキュメント更新手順を確認したい → docs/tool_implementation/documentation_guide.md
```

### 「ドキュメントを追加・移動・整理したい」

```text
ドキュメントを追加・移動・整理したい
└─ docs管理ルールとリンク整合性を確認する → docs/docs_management/index.md
```

### 「実装済みサービスを調査したい」

```text
実装済みサービスを調査したい
└─ 実装状況一覧（CLI/MCP/gRPC/HTTP）を確認する → docs/project_status/service_implementation_status.md
```

### 「計画（active/completed）を確認したい」

```text
計画を確認したい
└─ 進行中および完了計画を見たい → docs/exec_plans/index.md
```

### 「全体設計を確認したい」

```text
全体設計を確認したい
└─ アーキテクチャ/構成を参照する → docs/project_overview/project_overview.md
```

### 「依頼文の例を見たい」

```text
依頼文の例を見たい
└─ プロンプトサンプルを参照する → docs/user_prompt/prompt_sample.md
```
