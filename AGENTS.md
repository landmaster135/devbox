# AGENTS.md (docs マップ)

このファイルは、エージェント向けの「入口（目次）」です。  
実装・設計・運用の詳細は `docs/` を正本として参照してください。

## 参照優先順

1. `docs/tool_implementation/implementation_guide.md`
2. `docs/tool_implementation/documentation_guide.md`
3. `docs/docs_management/index.md`
4. `docs/project_status/service_implementation_status.md`
5. `docs/project_overview/project_overview.md`
6. `docs/user_prompt/prompt_sample.md`

## docs マップ

- `docs/tool_implementation/implementation_guide.md`
  - 用途: CLI/MCP実装時の注意点、実装パターン、失敗しやすいポイントの確認
  - 読むタイミング: 新規ツール実装、既存実装の修正、MCP実装時

- `docs/tool_implementation/documentation_guide.md`
  - 用途: `cmd/` 配下ツールの実装・改修後に必要なドキュメント更新手順の確認
  - 読むタイミング: README更新や実装ガイド更新時

- `docs/docs_management/index.md`
  - 用途: docs配置方針とAGENTSマップのリンク整合性確認
  - 読むタイミング: ドキュメント追加・移動・リネーム時

- `docs/project_status/service_implementation_status.md`
  - 用途: サービス実装状況（CLI/MCP/gRPC/HTTP）の確認と重複実装の回避
  - 読むタイミング: 新機能追加前、同名/類似ツール調査時

- `docs/project_overview/project_overview.md`
  - 用途: 全体アーキテクチャ、ディレクトリ構成、ビルド/テスト/品質方針の確認
  - 読むタイミング: 設計方針確認、ビルド/テスト手順確認、構成変更時

- `docs/user_prompt/prompt_sample.md`
  - 用途: 過去の依頼テンプレートやプロンプト例の参照
  - 読むタイミング: 依頼文作成やワークフロー定義時の参考

## 最小運用ルール（このファイルに残す項目）

- 回答は日本語で行うこと。

## 更新ルール

- 詳細ルールをこのファイルへ増やさないこと。詳細は `docs/` 側へ追記する。
- `docs/` に新規ドキュメントを追加した場合は、このマップにリンクと用途を追加する。
