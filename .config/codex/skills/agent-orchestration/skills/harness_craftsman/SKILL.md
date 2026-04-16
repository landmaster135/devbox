---
name: harness-craftsman
description: エージェント用ハーネスを新規作成・改修するスキル。リポジトリの現状把握、AGENTS.md の作成、docs 配下の初期ドキュメント整備を段階的に実施し、既存ドキュメントがある場合は .agents/tmp に改修案を書き出す。
---

# Harness Craftsman

## 実行手順

1. 入力情報を確認する。
2. ユーザーから指定が無ければ、対象リポジトリのルートパスを確認する。
3. リポジトリ構成を把握する。最低限、次を確認する。
  - ルート直下の `AGENTS.md` の有無
  - `docs/` の有無と主要インデックスの有無
  - `.agents/tmp/` の有無
4. `AGENTS.md` が無い場合は新規作成する。存在する場合は改修案を作成する。
5. `docs/` が無い場合はディレクトリと初期ドキュメント群を新規作成する。存在する場合は不足・重複・リンク不整合の改修案を作成する。
6. `scripts/directory_validate.go` を実行して、ベースライン構成が満たされていることを検証する。
7. `scripts/empty_docs_inspection.go` を実行して、`docs/` 配下に空ファイルが無いことを検証する。
8. 既存ドキュメントがある場合、改修案を `.agents/tmp` に Markdown で書き出す。
9. 最終出力として、作成または改修案出力したファイル一覧と、次に実施すべき作業を簡潔に報告する。

## 基本方針

- 回答・ドキュメント記述は日本語で行う。
- 詳細実装ルールは `AGENTS.md` に肥大化させず、`docs/` 側へ分離する。
- 既存資産は消さずに活かし、差分ベースで改修案を提示する。
- 既存の利用導線を壊さないよう、インデックス文書を入口として整備する。

## 作成フロー

### 1. リポジトリ現状把握

次の観点で調査する。

- プロジェクト概要: 主要ディレクトリ、主要エントリポイント、言語・実行基盤
- エージェント向け入口: `AGENTS.md` の有無、内容、リンク先
- ドキュメント構成: `docs/` 配下のカテゴリ整理、indexの有無
- 一時作業領域: `.agents/tmp/` の有無

### 2. AGENTS.md 整備

`AGENTS.md` は「概要と案内」に限定し、詳細ルールは `docs/` に逃がす。

最低限含める。

- Summary
- Prerequisites
- Quick Decision Trees
  - タスク着手
  - 実装・改修
  - docs管理
  - プロジェクト状況確認

既存 `AGENTS.md` がある場合の対応。

- 既存見出し構造は維持し、リンク切れと重複記述を優先修正候補にする。
- ルール追加が `AGENTS.md` 側に偏っている場合は、`docs/` へ移管する改修案を作る。

具体例（新規作成時）。

```md
# AGENTS.md (docs マップ)

## Summary
- このファイルは入口。詳細ルールは `docs/` 側を正本とする。
- `docs/tool_implementation/index.md`: 実装・改修手順
- `docs/docs_management/index.md`: docs 管理ルール

## Prerequisites
- 回答は日本語で行う。
- 詳細ルールは AGENTS.md に書きすぎない。

## Quick Decision Trees
### 「タスクに着手したい」
- 指示書を確認したい → `.agents/tmp/instructions.md`
```

具体例（既存改修時）。
- NG: AGENTS.md に長いコーディング規約を追記する。
- OK: AGENTS.md には「規約は `docs/task_implementation/index.md` を参照」とだけ書く。

### 3. docs 初期整備

`docs/` が無い場合は、次のような入口中心構成を作る。

- `docs/changelog/README.md`
- `docs/tool_implementation/index.md`
- `docs/docs_management/index.md`
- `docs/exec_plans/index.md`
- `docs/project_status/index.md`
- `docs/project_overview/index.md`
- `docs/task_implementation/index.md`
- `docs/user_prompt/index.md`

必要に応じて配下の詳細ドキュメントを追加する。追加時は必ず index から辿れるようにする。

既存 `docs/` がある場合の対応。

- 新規大量作成は避ける。
- まず index の入口品質を揃える。
- 不足情報は「改修案」として `.agents/tmp` に出力する。

### 3.1 ディレクトリ構成のバリデーション

次のコマンドで、`docs` 配下がベースライン構成を満たすか検証する。

```bash
go run $HOME/.codex/skills/agent-orchestration/skills/harness_craftsman/scripts/directory_validate.go --docs-dir <対象ディレクトリ>
```

- 検証に失敗した場合は不足項目を埋めるか、既存構成を残す判断をしたうえで `.agents/tmp/harness_doc_improvement_plan.md` に改修案として記録する。

### 3.2 空ファイル検査

次のコマンドで、`docs/` 配下に空ファイルが存在しないか検証する。

```bash
go run $HOME/.codex/skills/agent-orchestration/skills/harness_craftsman/scripts/empty_docs_inspection.go --docs-dir <対象ディレクトリ>
```

- 空ファイルが検出された場合は対象ファイルを修正し、再実行で `OK` を確認する。
- 既存構成を優先して修正を保留する場合は、理由を `.agents/tmp/harness_doc_improvement_plan.md` に記録する。

### 4. 既存ドキュメントがある場合の改修案出力

出力先は `.agents/tmp` 固定とし、次のファイル名を使う。

- `.agents/tmp/harness_doc_improvement_plan.md`

改修案には次を含める。

1. 現状サマリー
2. 問題点
3. 改修方針
4. 具体的な修正候補ファイル
5. 実施順序

## 改修案テンプレート

既存ドキュメントがある場合、次の形式で `.agents/tmp/harness_doc_improvement_plan.md` を作成する。

```md
# Harness Documentation Improvement Plan

## 現状サマリー
- AGENTS.md: <要約>
- docs/: <要約>
- .agents/tmp/: <要約>

## 問題点
1. <問題1>
2. <問題2>

## 改修方針
1. <方針1>
2. <方針2>

## 修正候補ファイル
- <path>: <変更概要>
- <path>: <変更概要>

## 実施順序
1. <step1>
2. <step2>
3. <step3>
```

## 最終チェック

- `AGENTS.md` が入口として機能しているか。
- `AGENTS.md` に詳細ルールを過度に書いていないか。
- `docs/` の主要 index へ遷移できるか。
- `scripts/directory_validate.go` でベースライン構成が検証済みか。
- `scripts/empty_docs_inspection.go` で空ファイル検査が完了しているか。
- 既存ドキュメントがある場合、`.agents/tmp/harness_doc_improvement_plan.md` が生成されているか。
- 生成または提案したファイルの一覧を最終報告で明示したか。
