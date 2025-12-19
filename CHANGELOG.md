# Changelog

## v0.05 — 2025-12-19

PR: #18

### Features
- **service-implementing-viewer (CLI/MCP)**: `-operation`/`-write`フラグを追加し、出力と同時にドキュメントへ反映できる書き込みモードを実装。
- **自動ドキュメント生成**: `internal/service_implementing_viewer/usecases/document_updater`を新設し、docs配下の実装状況表をコマンドから更新可能に。
- **Issue/PRテンプレート**: `.github/ISSUE_TEMPLATE`と`PULL_REQUEST_TEMPLATE`を整備し、提出時のチェックリストと概要入力を標準化。

### Improvements
- `cmd/cli/service-implementing-viewer/README.md`の利用手順やサンプルを刷新し、writeモードの例を追記。
- `docs/project_overview.md`と`docs/service_implementation_status.md`に自動更新フローや統計テーブルの最新値を反映。
- READMEのパッケージ概要セクションを再整理し、読みやすさを向上。

### Refactor
- 旧`util`パッケージ（ロガー/ユーティリティ）の撤去に伴い、依存箇所を直接実装へ移行して構成を単純化。

### CI / Testing
- GitHub Actionsからツールの実装状況を自動更新するコミットを追加し、実装状況のドキュメントの陳腐化を防止。

## v0.04 — 2025-12-19
PR: #13

### Features
| service | cli | mcp | grpc | http | 説明 |
| :-- | :--: | :--: | :--: | :--: | :-- |
| anilist | ✅ | ❌️ | ❌️ | ❌️ | AniList GraphQL APIからユーザーのアニメ／マンガ進捗を取得しJSONや表形式で出力。 |
| color-code-converter | ✅ | ❌️ | ❌️ | ❌️ | HEX・RGB・HSL・HSVの相互変換に対応したカラーコード変換CLI。 |
| data-converter | ✅ | ❌️ | ❌️ | ❌️ | JSON/CSV/TSV/HTML/Markdownリストを自在に相互変換できるデータ変換CLI。 |
| db-server-sync | ✅ | ❌️ | ❌️ | ❌️ | AniListエクスポートと追加データを結合し、DBサーバー用リクエストボディを生成。 |
| git-pre-commit-hooks | ✅ | ❌️ | ❌️ | ❌️ | MCP設定を含むリポジトリ全体のシークレットと禁止パスを検知するプレコミットフック。 |
| grpc-request | ✅ | ❌️ | ❌️ | ❌️ | gRPCリフレクションを用いたサービス探索とJSONベースのリクエスト送信を行うCLI。 |
| postgresql | ✅ | ✅ | ❌️ | ❌️ | テーブル一覧・全テーブルダンプ・並列エクスポートを備えたPostgreSQL CLI（MCPも継続）。 |
| steam | ✅ | ❌️ | ❌️ | ❌️ | Steam Web APIから所有ゲーム・統計・実績を取得しJSONに出力するCLI。 |
| youtube-downloader | ✅ | ❌️ | ❌️ | ❌️ | YouTube動画やプレイリストを並列ダウンロードしFFmpegで高品質結合できるCLI。 |
| weather-notificator | ✅ | ✅ | ✅ | ✅ | OpenWeatherMap予報をCLI/MCP/HTTP/gRPC経由で配信し、Discord埋め込み通知にも対応。 |
| gcloud-genset-ai | ✅ | ❌️ | ❌️ | ❌️ | Document AIプロセッサのアンデプロイ用curl/gcloudコマンドを安全に組み立てるCLI。 |
| gcloud-genset-billing | ✅ | ❌️ | ❌️ | ❌️ | Billingの予算・プロジェクト情報取得コマンドを自動生成し、フィルタやリミットも補助。 |
| gcloud-genset-cloudsql | ✅ | ❌️ | ❌️ | ❌️ | Cloud SQLの削除/削除保護/起動ポリシー変更を手順どおりに並べたgcloudコマンドを出力。 |
| gcloud-genset-container | ✅ | ❌️ | ❌️ | ❌️ | Cloud Run/Functions/Pub/Sub操作とDiscord通知スニペットを一括生成し、環境変数更新にも対応。 |
| gcloud-genset-deployment | ✅ | ❌️ | ❌️ | ❌️ | Deployment Managerのデプロイ一覧取得コマンドをフィルターやフォーマット付きで実行。 |
| gcloud-genset-dns | ✅ | ❌️ | ❌️ | ❌️ | Cloud DNS managed-zonesリスト用のgcloudコマンドを柔軟なフラグ付きで生成。 |
| gcloud-genset-iam | ✅ | ❌️ | ❌️ | ❌️ | サービスアカウント操作とWorkload Identity Federationセットアップ/クリーンアップを支援するコマンド群。 |
| gcloud-genset-init | ✅ | ❌️ | ❌️ | ❌️ | プロジェクト初期設定で使う`gcloud auth login`や`config set project`を追加引数込みで整形。 |
| gcloud-genset-logging | ✅ | ❌️ | ❌️ | ❌️ | Logging readやシンク作成コマンドをseverity/フィルター/追加引数から組み立てるCLI。 |
| gcloud-genset-monitoring | ✅ | ❌️ | ❌️ | ❌️ | Monitoringダッシュボード/スヌーズ/アップタイム設定の`gcloud monitoring`コマンドをテンプレ化。 |
| gcloud-genset-monitoring-dashboard | ✅ | ❌️ | ❌️ | ❌️ | Cloud Runサービス向けに16ウィジェット構成の監視ダッシュボードをAPI経由で作成。 |
| gcloud-genset-scheduler | ✅ | ❌️ | ❌️ | ❌️ | Cloud SchedulerのHTTP/PubSub/CloudSQLジョブ生成と更新/停止/削除コマンドを自動化。 |
| gcloud-genset-secret | ✅ | ❌️ | ❌️ | ❌️ | Secret Managerの作成・値登録・エイリアス/ラベル更新コマンドとDiscord通知スクリプトを生成。 |
| gcloud-genset-storage | ✅ | ❌️ | ❌️ | ❌️ | GCSアップロード/ダウンロード/ACL/バケット作成などのgsutilコマンドを整形出力。 |
| image-filterer-v2 | ✅ | ❌️ | ❌️ | ❌️ | グレースケール/ティント/ビネットをCPUのみで適用できる第2世代画像フィルタCLI。 |
| image-renamer-for-content | ✅ | ❌️ | ❌️ | ❌️ | コンテンツIDプリセットと並列ワーカーでWeb素材を一括リネームするCLI。 |
| withings | ✅ | ❌️ | ❌️ | ❌️ | Withings Public Health Data APIでOAuth/日次サマリ取得/トークン自動更新を行うヘルスデータCLI。 |

### Improvements
- Discord通知パイプラインにOpenWeatherMapスタイルの埋め込み・完了通知・JST表記を追加し配信品質を向上。
- Git差分／履歴系ツールへ高度なパスバリデーションとセキュリティ例外体系を導入し、アーカイブ生成も安全化。
- AniList・Steam・PostgreSQL各サービスを依存性注入構成へ再設計し、高負荷操作や並列処理に対するテストを強化。
- Data ConverterとDB Server SyncでMarkdownテーブルや配列整形を拡充し、AniList連携データの整合性を向上。
- Notion SyncやOps for Golangなど既存CLIにWebクリップパッチやCLIオプション拡張を追加し、開発体験を改善。
- サービス実装統計機能で作成されるレポートに統計情報に関する内容を追加。
- .config: Codex/Cagent向け設定テンプレートを追加。
- Go のバージョンアップデート: 1.23.5 -> 1.25.5
- gcloud-gensetファミリーをCloud Run/Secret Manager/Storage/CloudSQL/Billing/Logging/DNS/Deployment/Scheduler/Containerまで拡張し、Discord通知テンプレートとWorkload Identity Federation雛形もCLIで生成できるよう統合。
- WithingsヘルスサービスをAuth/Core/Health層へ再分割し、測定タイプ`all`、JSONファイルエクスポート、アクセストークン自動更新、キャッシュ最適化、包括的テストで日次サマリ取得を堅牢化。
- image-renamer-for-contentにmackerel/web_clip/date/wineプリセットやソート・開始番号・ワーカー数指定を加え、DOSバッチ／ビルドスクリプトとテストで画像命名ワークフローを高速化。
- setup-git-pre-commit-hooksスクリプトとgit commit message generatorをモジュール化し、cagentブログエージェント設定・Anthropicモデル・プロンプト再編・Shell環境ポリシーでワークフローを共通化。
- OCR ExecutorをAIクライアントから分離し、Ollama OCR対応・http.Client再利用・ストリーミングデコードでマルチモデル処理の遅延を削減。
- PostgreSQLダンプ機能の追加。

### Bug Fixes
- git-info-retrieverのアーカイブ生成スクリプトをサニタイズし、任意のパス注入を防止。
- Gitコマンド用パスバリデータとgit-commit-history-executorの安全なワーキングディレクトリ検証でURLエンコードやシンボリックパス、ディレクトリ外アクセスを遮断し、差分／履歴取得時の安全性を向上。
- PostgreSQLテーブルダンプ系ツールでschema-qualifiedクエリ、識別子クオート、列順固定、並列実行ガード、CSV/SQLエンコーダの正規化を行い、エクスポート結果の破損を防止。
- script-generator-to-buildの使用例にバックスラッシュとクォートのエスケープを施し、生成スクリプトをコピー&ペーストしても安全に動作。

### Documentation
- docs/service_implementation_status.md を更新し、gRPC/HTTP列と統計を追加して最新のサービス実装状況を反映。
- docs/project_overview.md や docs/implementation_guide.md を刷新し、新規CLI・HTTP/gRPCサーバーの構成と導線を記載。
- cmd/grpc/docs/grpc_setting_guide.md や Weather Notificator ハンドラ README を追加し、サーバー構築と利用手順を文書化。

## v0.03 — 2025-08-13
PR: #7

### Features
| service | cli | mcp | 説明 |
| :-- | :--: | :--: | :-- |
| arithmetic-calculator | ✅ | ✅ | 数式評価とAPIコスト見積りを提供する計算ツール |
| arxiv | ✅ | ❌️ | arXiv論文を検索してメタデータを抽出するリサーチ支援CLI |
| base64-extractor | ✅ | ❌️ | Base64文字列を抽出・再利用するユーティリティ |
| brave-search | ❌️ | ✅ | Brave Search APIを呼び出して検索結果を返すMCPサーバー |
| context7 | ✅ | ✅ | ライブラリドキュメントを即時取得するContext7連携ツール |
| datetime-calculator | ✅ | ✅ | 日時演算・単位変換・自然言語解析を行う計算サービス |
| discord-webhook | ✅ | ❌️ | Discord Webhookへ任意メッセージを送信する通知CLI |
| duckduckgo-search | ❌️ | ✅ | DuckDuckGo検索をプロトコル越しに提供するMCPサーバー |
| everart | ❌️ | ✅ | EverArtモデルを使った画像生成MCPサーバー |
| file-character-replacer | ✅ | ❌️ | テキスト内の文字種や表記ゆれを一括置換 |
| figma | ❌️ | ✅ | Figmaドキュメントやアセットを取得するMCPサーバー |
| filesystem | ❌️ | ✅ | CLI/MCPツールのためのファイル操作を抽象化 |
| gcloud-monitoring | ✅ | ❌️ | Cloud Run監視ダッシュボードを自動生成 |
| gcloud-wrapper-workload-identity-federation | ✅ | ❌️ | Workload Identity Federation設定を生成 |
| gdrive | ❌️ | ✅ | Google Driveファイルの取得を行うMCPサーバー |
| git-commit-history-retriever | ✅ | ✅ | Git履歴を抽出しレポート化する履歴解析ツール |
| git-diff-recorder | ✅ | ✅ | 差分を構造化記録・再生できる差分管理ツール |
| git-info-retriever | ✅ | ❌️ | リポジトリ情報を収集し書き出すメタデータ収集CLI |
| github | ✅ | ✅ | GitHub APIでIssue/PR/コミットを扱う統合ツール |
| http-request | ✅ | ✅ | 任意のHTTPリクエストを実行する汎用クライアント |
| mcp-remote | ✅ | ❌️ | 複数MCPサーバーを中継するリモートプロキシCLI |
| memory | ✅ | ❌️ | Valkey＋ファイルを併用したナレッジ管理サービス |
| notion-blog-content-extractor | ✅ | ❌️ | NotionブログMarkdownから本文とメタ情報を抽出 |
| notion-sync | ✅ | ✅ | MarkdownをNotionページへ反映する同期ツール |
| ocr-executor | ✅ | ❌️ | 画像からテキストを抽出するローカルOCRツール (#10) |
| ocr-executor-with-ai | ✅ | ❌️ | Gemini/Vertex AIを使ってOCRと表生成を行うCLI |
| open-weather-map | ✅ | ✅ | OpenWeatherMapから天気予報を取得・パース |
| ops-for-golang | ✅ | ✅ | go run/test/buildの実行を自動化する開発支援ツール |
| postgresql | ❌️ | ✅ | PostgreSQLへの接続・操作を提供するMCPサーバー |
| qdrant | ✅ | ❌️ | Qdrantベクターストアへ接続してデータ管理 |
| sequentialthinking | ❌️ | ✅ | 段階的思考プロセスをMCP経由で提供 |
| service-implementing-viewer | ✅ | ✅ | サービス実装状況を集計する可視化ツール |
| shell | ❌️ | ✅ | シェルコマンド実行を抽象化するMCPサーバー |
| timezone | ❌️ | ✅ | タイムゾーン変換・取得を提供するMCPサーバー |
| util | ❌️ | ✅ | MCPサーバー共通のユーティリティ群 |
| valkey | ✅ | ❌️ | Valkeyの起動管理とデータ操作を行うCLI |
| weather-notificator | ✅ | ✅ | 天気予報を取得して通知するサービス |
| youtube-transcript | ❌️ | ✅ | YouTube動画の字幕を取得するMCPサーバー |
| zip-compressor | ✅ | ❌️ | 設定駆動でZIPアーカイブを生成するCLI |

### Improvements
- git-diff-recorder: 読み取りモード・生成モード・未追跡ファイル対応を追加し、出力を構造化フォーマットへ統一。
- ops-for-golang: 実行ディレクトリ指定と終了コード別エラーハンドリングを実装し、テスト自動化を強化。
- gcloud-monitoring/gcloud-wrapper-workload-identity-federation: ダッシュボード構成と認証スクリプト生成のバリデーションを整理。
- .config: Claude Desktop、Claude Code、Cline、Discord Bot向け設定テンプレートを追加。
- .github/workflows/test_integration.yml: ffmpegインストールをOS別に分岐しMCP追加ツールのCIを整備。
- scripts/build_mcp_tools.sh: MCPサーバー群の一括ビルドを追加。
- image-filterer: Blur処理とオプションを再構成し再利用性を向上。
- pdf-merger: 結合処理のロジックと引数レイアウトを整理。

### Bug Fixes
- git-diff-recorder: nilポインタ参照を解消し安定性を向上。

### Documentation
- arithmetic-calculator, arxiv, base64-extractor, context7, datetime-calculator, discord-webhook, file-character-replacer, gcloud-monitoring, gcloud-wrapper-workload-identity-federation, git-commit-history-retriever, git-info-retriever, github, http-request, mcp-remote, memory, notion-blog-content-extractor, notion-sync, ocr-executor, ocr-executor-with-ai, open-weather-map, ops-for-golang, qdrant, service-implementing-viewer, valkey, weather-notificator, zip-compressor 各READMEを新規追加または全面更新。
- README.md: MCP統合と新CLI群を反映した構成説明とアーキテクチャ図を追記。

## v0.02 — 2025-07-10
PR: #2

### Features
| service | cli | 説明 |
| :-- | :--: | :-- |
| claude-code-usage | ✅ | Claude APIの利用状況を集計してレポート化するCLI |
| code-analyzer | ✅ | プロジェクト内ソースコードの品質指標を分析するCLI |
| depends-visualizer | ✅ | 関数・モジュール依存関係をMermaid/PlantUML/Dotで可視化 |
| diff-dreamer | ✅ | オフライン差分閲覧UIを提供するCLI (#4) |
| exif-mirror | ✅ | EXIFデータを別ファイルへコピーするツール |
| exif-modifier | ✅ | EXIFプロパティを並列編集できるCLI |
| exif-viewer | ✅ | EXIF/メタデータをテーブル表示するCLI |
| file-maneuver | ✅ | 複数ディレクトリから条件一致ファイルを集約・移動するCLI |
| git-diff-recorder | ✅ | Git差分を構造化フォーマットで記録・読み出すCLI |
| goo-scraper | ✅ | goo系サイトから情報を取得するスクレイピングCLI |
| image-filterer | ✅ | 画像にぼかし等のフィルタ処理を適用するCLI |
| image-renamer | ✅ | 連番命名とプレフィックス付与を行う画像リネーマー |
| image-renamer-for-screenshot | ✅ | デバイス毎のスクリーンショット命名を統一するCLI |
| image-renamer-with-exif | ✅ | EXIF CreateDateを利用して日時命名するCLI |
| image-rotator | ✅ | 画像を指定角度で回転するツール |
| image-trim-describer | ✅ | 画像トリミング領域をブラウザUIで確認するCLI |
| image-trimmer | ✅ | 用途別プリセットで画像をトリミングするCLI |
| json-formatter-for-agent-interaction | ✅ | エージェント連携向けJSON整形を行うCLI |
| kana-converter | ✅ | 全角/半角カナやUnicode変換を行うCLI |
| movie-converter-for-gif | ✅ | GIFと動画間の相互変換を行うCLI |
| movie-converter-for-webm | ✅ | MP4/WebMの変換と品質調整を行うCLI |
| script-generator-to-build | ✅ | 指定CLIのビルド用シェルスクリプトを生成するCLI |
| unit-converter | ✅ | 一般的な単位換算を実行するCLI |

- cmd/powershell/Z9-11_machine_info_retriever.ps1: Windows向けマシン情報収集スクリプトを追加。

### Improvements
- image-converter: 引数体系とヘルプ表示を整理しWebP対応を強化。
- pdf-merger: 既存ファイル統合フローを刷新しエラーハンドリングを改善。
- env-loader / file-processor / json-* 系: I/OバリデーションとCRLFケースのテストを追加。
- internal/independenciesディレクトリの削除とアーキテクチャリファクタリング (#5)
- scripts/build.sh: 全CLIをビルド対象に含めるよう更新。
- scripts/create_project_files.sh: CLIテンプレートを自動生成するスクリプトを追加。
- .github/workflows/pull_request_stat.yml: プルリク統計収集ワークフローを追加。

### Bug Fixes
- 特筆すべき修正は報告されていません。

### Documentation
- 各新規CLI（claude-code-usage から unit-converter まで）のREADMEを追加。
- 既存ツール（image-converter, pdf-merger 等）の利用手順を更新。
- README.md と .clinerules: 新ツール群の方針とカバレッジ目標を追記。

## v0.01 — 2025-05-03
PR: #1

### Features
| service | cli | 説明 |
| :-- | :--: | :-- |
| http-request (旧api-client) | ✅ | REST/GraphQLエンドポイント向け汎用HTTPクライアント |
| env-loader | ✅ | .envファイルを読み込んで環境変数を整備するCLI |
| file-processor | ✅ | ファイル入出力や簡易加工を行うユーティリティ |
| image-converter | ✅ | 画像形式を相互変換するCLI |
| iso8601-converter | ✅ | ISO-8601日時フォーマットを変換するCLI |
| json-file-merger | ✅ | 複数JSONファイルを統合するCLI |
| json-iso8601-converter | ✅ | JSON内日時フィールドをISO-8601へ変換するCLI |
| json-modifier | ✅ | JSONの取得・更新・削除操作を行うCLI |
| json-timestamp-modifier | ✅ | JSON内タイムスタンプのUNIX変換を行うCLI |
| pdf-encrypter | ✅ | PDFへパスワードを設定するCLI |
| pdf-merger | ✅ | 複数PDFを結合するCLI |
| yaml-parser | ✅ | YAMLを解析し各種フォーマットへ出力するCLI |

### Improvements
- .gitignore: 生成物やIDE設定ファイルを除外対象に追加。
- pkg/dos/: Windows向けバッチスクリプト群の追加を開始。

### Bug Fixes
- 特筆すべき修正は報告されていません。

### Documentation
- 初期CLI群 (http-request, env-loader, file-processor, image-converter, iso8601-converter, json-file-merger, json-iso8601-converter, json-modifier, json-timestamp-modifier, pdf-encrypter, pdf-merger, yaml-parser) のREADMEを整備。
- pkg/dos/README.md: バッチスクリプトの使い方を追記。
- .clinerules: Clean Architecture方針とカバレッジ目標を定義。
- README.md: プロジェクト概要と初期ツールセットの導入手順を整理。
