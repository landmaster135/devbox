# project_status ドキュメントガイド

`docs/project_status` は、プロジェクトのサービス構成・実装状況・利用可能な入口情報を把握するためのディレクトリです。

## ドキュメント一覧（概要と責務）

| ドキュメント | 主な責務 |
|---|---|
| `index.md` | `project_status` 配下の参照順を定義し、必要な文書へ素早く到達できるようにする |
| `service_implementation_status.md` | サービス単位の実装有無（CLI/MCP/gRPC/HTTP）を管理し、追加・移行時の判断を補助する |
| `service_summary.md` | 各サービスの用途と機能要点を管理し、候補選定の初期判断と重複機能の早期検知に役立つ |
| `entrypoint_overview.md` | インストール・ビルド・起動の入口を管理し、利用者が環境別の導入手順を迷わず実行できるようにする |
