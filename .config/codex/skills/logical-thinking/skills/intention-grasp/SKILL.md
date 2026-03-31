---
name: intention-grasp
description: Clarify ambiguous user intent by inferring likely goals from user actions and asking concise follow-up questions before making irreversible assumptions. Use when instructions are vague, actions appear contradictory, directory/role swaps occur, or the user seems dissatisfied with interpretation.
---

# Intention Grasp

相手の意図が不明瞭なとき、行動の背景を確認し、合意した前提で作業を進める。

## Core Rules

- 先に解釈し切らない。曖昧さを検知したら確認質問を優先する。
- 相手が取った行動を事実として要約してから質問する。
- 質問は短く具体的にする。1回で1〜2問までに抑える。
- 変更前に「採用する前提」を明文化し、相手の承認を得る。
- 承認前に破壊的変更や大きな再配置を行わない。

## Ambiguity Signals

以下のいずれかがあれば、確認フェーズに入る。

- 指示と実際の操作結果が矛盾している。
- 同名リソースの役割が入れ替わっている可能性がある。
- ユーザーが不満・違和感を示している。
- 複数の解釈で結果が大きく変わる。

## Workflow

1. 観測事実を2〜4行で要約する。
2. 意図の候補を最大2つに絞って提示する。
3. 候補を分岐させる確認質問を1〜2問する。
4. 回答を受けたら、採用前提を1文で宣言する（必要なルールを複数列挙する）。
5. 宣言した前提に沿って実装・操作する。

## Question Patterns

- 「現在の目的は `A` と `B` のどちらですか？」
- 「`X` を本番、`Y` を検証用として固定して進めてよいですか？」
- 「今回の変更で優先すべきのは `命名の整合` と `動作の復旧` のどちらですか？」

## Confirmation Template

以下の形式を使って前提を固定する。

`前提: <対象1> を <役割1> とし、<対象2> を <役割2> とし、...、<対象n> を <役割n> として進めます。`

## Confirmation Rules

- 対象/役割ペアは `n` 個（n>=1）で定義できる。必要な数だけ列挙する。
- 重要度が高い順に並べる（環境定義 → 変更範囲 → 変更不可条件）。
- 制約も対象/役割として書く（例: `既存API挙動` を `変更不可条件`）。
- 前提文は必ず1文で完結させる。

## Confirmation Examples

- `前提: apps/blog を Twilight（本番側） とし、apps/sample を 旧ブログ（検証側） とし、ルート scripts を 共通運用対象 として進めます。`
- `前提: authモジュール を リファクタリング対象 とし、既存API挙動 を 変更不可条件 とし、テスト失敗0件 を 完了条件 として進めます。`

## Anti-Patterns

- 相手の行動意図を聞かずに「正しそうな形」に戻す。
- 長い尋問を行う。
- 確認なしに前提を途中変更する。
