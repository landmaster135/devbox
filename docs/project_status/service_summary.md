## CLI tools
*anilist*: AniListからアニメ・マンガ情報を取得するためのコマンドラインツールです。
*arithmetic-calculator*: 算術計算を行うCLIツールです。基本的な四則演算から高度な数学関数まで、幅広い計算機能を提供します。
*arxiv*: ArXivから論文を検索・取得するためのコマンドラインツールです。
*base64-extractor*: 任意のパス（ファイルもしくはディレクトリ）にある画像ファイルのbase64形式のバイト列を抽出するCLIツールです。
*claude-code-usage*: Claude Code使用状況分析ツールのCLIインターフェース
*code-analyzer*: コードベースを分析し、複雑度、コメント率、コードクローンなどの重要なメトリクスを収集・可視化するツールです。
*color-code-converter*: 任意のカラーコード形式から別のカラーコード形式に変換するためのCLIツールです。
*context7*: Context7 CLIは、最新のライブラリドキュメントを取得するためのコマンドラインツールです。Context7 APIを使用して、ライブラリの検索とドキュメントの取得を行います。
*cron-workflow*: あらかじめ定義したバックグラウンドジョブを継続的に動かし続ける CLI ツールです。バイナリは単一のスケジューラを起動し、`workflow/core.go` に宣言されたワークフローを登録してから `SIGINT` もしくは `SIGTERM` を受け取るまで待機します。
*data-converter*: 入力ファイルを別形式ファイルへ変換する CLI ツールです。
*datetime-calculator*: 日時計算を行うCLIツールです。指定された基準日時に対して年月日時分秒の加算・減算を行い、結果を表示します。
*db-server-sync*: AniListデータをデータベースサーバー用のリクエストボディ形式に変換するCLIツールです。
*depends-visualizer*: プログラムファイル内の関数間の依存関係を解析して可視化するツールです。
*diff-dreamer*: テキスト比較用のGUIをHTMLページで提供します。
*discord-webhook*: Discord WebhookでメッセージやEmbed付き通知を送信するためのCLIツールです。
*docker*: env.yml に定義した情報を `docker-compose.yml` へ同期するユーティリティです。`npm run docker:env` 相当の環境変数反映に加え、指定サービスのポート番号やボリューム定義を env.yml ベースで更新できます。
*env-loader*: 
*exif-mirror*: 指定されたソースフォルダ内の画像ファイルから、ターゲットフォルダ内の同名ファイルにEXIFデータをコピーするCLIツールです。
*exif-modifier*: 任意のディレクトリ内にある画像・動画ファイルのExifプロパティを編集するCLIツールです。
*exif-viewer*: 画像ファイルのEXIF・メタデータ情報を参照・表示するためのCLIツールです。
*file-character-replacer*: ファイルまたはディレクトリ内のファイルに対して、指定した文字列を別の文字列に置換するCLIツールです。
*file-maneuver*: 複数のディレクトリから指定した拡張子やファイル名部分一致条件に合致するファイルを検索し、単一の宛先ディレクトリに移動またはコピーするCLIツールです。
*file-processor*: `file-processor` は、指定されたファイルの各行から指定範囲の文字列を抽出し、その文字列が重複する行を削除するコマンドラインツールです。
*filesystem*: ファイル読み書きや検索、移動などの基本操作を安全に行うCLIツールです。
*gcloud-genset-ai*: Document AI の運用で利用する `gcloud`/`curl` コマンドを整形して出力する CLI ツールです。手動入力時に煩雑になりがちなアクセストークン取得やエンドポイントの組み立てを安全に自動化します。
*gcloud-genset-billing*: Google Cloud Billing の運用で利用する `gcloud billing` コマンドを安全に組み立てる CLI ツールです。請求アカウントや予算、プロジェクトに関する情報取得コマンドを素早く生成できます。
*gcloud-genset-cloudsql*: Cloud SQL インスタンスの運用で利用する `gcloud sql` コマンドを安全に組み立てて提示する CLI ツールです。
*gcloud-genset-container*: Cloud Run コンテナ／Cloud Functions (Gen2)／Pub/Sub の日常運用コマンドを生成するgcloud コマンドジェネレーターです。
*gcloud-genset-deployment*: Google Cloud Deployment Manager の運用で利用する `gcloud deployment-manager` コマンドを安全に組み立てて実行する CLI ツールです。
*gcloud-genset-dns*: Google Cloud DNS の `gcloud dns` コマンドを組み立てる CLI ツールです。
*gcloud-genset-iam*: Google Cloud IAM / Workload Identity Federation の `gcloud` コマンドを素早く組み立てる CLI ツールです。サービスアカウントの操作から Workload Identity Pool の構築・破棄まで、ツール内で管理しているスクリプトを Go 実装として再現し、コピペ可能なコマンド列を出力します。
*gcloud-genset-init*: Google Cloud プロジェクトの初期設定で利用する `gcloud` コマンドを生成する CLI ツールです。
*gcloud-genset-logging*: Google Cloud Logging の `gcloud` コマンドを生成する CLI ツールです。
*gcloud-genset-monitoring*: Google Cloud Monitoring の `gcloud` コマンドを生成する CLI ツールです。
*gcloud-genset-monitoring-dashboard*: Google Cloud Run サービス用のモニタリングダッシュボードを自動作成するCLIツールです。
*gcloud-genset-scheduler*: Google Cloud Scheduler 向けの gcloud コマンドを生成する CLI ツールです。`operation` を指定すると、必要なオプションを組み立てたコマンドを表示します。Cloud Run / Cloud Functions / Cloud SQL 向けのジョブ生成から、ジョブの更新・停止・削除まで幅広い操作をカバーします。
*gcloud-genset-secret*: Google Cloud Secret Manager の `gcloud` コマンドを生成する CLI ツールです。シークレットの作成・値の登録・バージョン取得・ラベル/エイリアス更新に必要なフラグを整理し、実行可能なコマンドラインを出力します。
*gcloud-genset-spanner*: Cloud Spanner の日常運用で使う `gcloud spanner` コマンドを安全に組み立て、コピーしやすい形で提示する CLI ツールです。インスタンス・データベースの準備で毎回調べがちなフラグをテンプレ化し、入力値の検証と整形を自動化します。
*gcloud-genset-storage*: Google Cloud Storage (GCS) 操作用の `gsutil` / `gcloud` コマンドを生成する CLI ツールです。Cloud Storage へのアップロード / ダウンロード / バケット作成 / ACL 操作など、日常的な運用コマンドを安全に組み立てます。
*gcloud-wrapper-workload-identity-federation*: Google Cloud Workload Identity FederationとGitHub Actions認証の設定を自動化するCLIツールです。
*git-commit-history-retriever*: Gitリポジトリのコミット履歴を取得し、フィルタリングして表示するCLIツールです。
*git-diff-recorder*: Git差分を記録・読み取りするCLIツールです。リポジトリの差分情報を構造化されたフォーマットでファイルに出力し、後で読み取ることができます。
*git-info-retriever*: GitHubからリポジトリ情報を取得し、Bash関数を生成してリポジトリのクローンとZip圧縮を行うCLIツールです。
*git-pre-commit-hooks*: Git pre-commit hook用のシークレット等に対する検知ツールです。JSON設定ファイル内の機密情報と、全ファイル内の禁止されたホームパスを自動検知し、コミット前にブロックします。
*github*: GitHubのイシューを取得するCLIツールです。
*goo-scraper*: 
*grpc-request*: gRPCサーバーにリクエストを送信するためのコマンドラインツールです。
*html-sanitizer*: main/article要素を基点にHTMLをサニタイズするCLIツールです。
*http-request*: 任意のURLのAPIにリクエストを送信するコマンドラインツールです。
*image-converter*: 複数の画像ファイルを一括で別のフォーマットに変換するコマンドラインツールです。
*image-filterer*: 画像の指定領域にフィルター効果を適用するためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。
*image-filterer-v2*: CPUベースで画像フィルタ処理を行うCLIツールです。グレースケール化、ティント付与、ビネット効果といった基本的なエフェクトをワンショットで適用できます。
*image-renamer*: 画像ファイルを指定したプレフィックスとシリアル番号でリネームするツールです。
*image-renamer-for-content*: 任意のディレクトリにある画像ファイルを、コンテンツIDと連番を基にした命名規則でリネームするCLIツールです。
*image-renamer-for-screenshot*: このツールは、VLCスナップショットファイル、Windowsスクリーンショットファイル、およびPixel端末で取得したスクリーンショット/録画ファイルを統一された命名規則でリネームするためのコマンドラインユーティリティです。
*image-renamer-with-exif*: 画像ファイルのEXIF CreateDateプロパティまたはファイルの更新時刻を使用して、ファイル名を年月日時分秒の形式（YYYYMMDDHHMMSS.拡張子）にリネームするCLIツールです。
*image-rotator*: 画像を回転させるためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。
*image-trim-describer*: 画像のトリミング座標を視覚的に選択し、トリミング処理を行うためのツールです。
*image-trimmer*: 画像をトリミングするためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。
*interactive-input*: ユーザーに1つの質問を投げかけ、標準出力へ機械的に扱いやすいキー付きの値を返すCLIツールです。Windowsバッチの `set /p` や `choice` の代替として設計されており、案内文やリトライ通知は標準エラーに集約されます。
*iso8601-converter*: UNIXタイムスタンプとISO-8601形式の相互変換を行うコマンドラインツールです。
*json-file-merger*: JSONファイルを読み込み、APIリクエスト用のリクエストボディを作成するツールです。
*json-formatter-for-agent-interaction*: 
*json-iso8601-converter*: 
*json-modifier*: 
*json-timestamp-modifier*: 
*kana-converter*: カタカナを含む文字列を全角または半角に変換するコマンドラインツールです。また、濁音・半濁音の追加や除去、濁音と半濁音の変換ペア処理も行えます。
*machine-info*: Ubuntu系LinuxでPCのハードウェア／ネットワーク情報を収集し、JSONとして保存するCLIツールです。
*mcp-remote*: MCP Remote CLIツールです。
*memory*: 知識グラフを使用したメモリ管理を行うCLIツールです。エンティティ、リレーション、観察事項を管理し、永続的なメモリ機能を提供します。
*memos*: Memos API（`/api/v1`）を操作するCLIツールです。
*movie-converter-for-gif*: GIFとMP4を相互に変換するCLIツールです。
*movie-converter-for-webm*: WEBMとMP4の相互変換を行うCLIツールです。
*notion-blog-content-extractor*: Markdownファイルからブログコンテンツを抽出するCLIツールです。特定のマーカーで区切られたコンテンツ部分のみを抽出し、新しいファイルとして保存します。
*notion-sync*: Notionのページにブロックを追加するためのCLIツールです。notion-synchronizerサーバーのページパッチAPIエンドポイントにHTTPリクエストを送信します。
*ocr-executor*: TesseractOCRを使用して画像ファイルからテキストを抽出するCLIツールです。
*ocr-executor-with-ai*: AI（Gemini API / Vertex AI / Ollama）を使用して画像からテキストを抽出するCLIツールです。
*ollama*: Ollama サーバーの HTTP API を手元から叩くための CLI です。
*open-weather-map*: OpenWeather APIを使用して天気予報を取得するコマンドラインツールです。
*ops-for-golang*: Go開発でよく実行するコマンドを自動化するCLIツールです。テストカバレッジの取得、カバレッジ分析、CLIツールの実行を効率化し、出力フィルタリング機能も提供します。
*pdf-encrypter*: 
*pdf-merger*: このツールは、画像からPDFを作成したり、PDFファイルから画像を抽出したりするためのコマンドラインツールです。
*postgresql*: PostgreSQLデータベースのテーブルダンプ機能を提供するCLIツールです。
*qdrant*: Qdrant の gRPC API を叩いてコレクション操作やベクトル upsert / 検索を行う CLI です。
*script-generator-to-build*: このツールは、指定されたGoパッケージのビルドスクリプトを生成するためのコマンドラインツールです。
*service-implementing-viewer*: 複数のディレクトリ内にあるサービスの実装状況を表形式で表示するCLIツールです。
*shell*: Codexの`shell`ツール互換でコマンドを安全に実行するCLIです。
*sqlite*: SQLite ファイルに対して操作を行う CLI です。
*steam*: Steam Web APIを使用してユーザーのゲーム情報を取得し、JSONファイルに出力するCLIツールです。
*taskfile*: Taskfile を検証・補完・新規作成するための CLI ツールです。指定した Taskfile が、プロジェクト標準 (`internal/taskfile/usecases/taskfiles/root.yml`) に含まれるフィールドと乖離しているかどうかを確認できます。
*unit-converter*: 物理量を **長さ・質量（重さ）・温度・面積・体積** の 5 カテゴリで相互変換できる高速 CLI ツールです。SI 接頭語 **yotta 〜 yocto** を自動解釈するため、`µm`, `nm2`, `pL` など自由に記述できます。
*valkey*: Valkey データベースを操作するためのコマンドラインツールです。高性能なキー・バリューストアであるValkeyに対して、包括的なデータ操作機能を提供します。
*vector-embedding*: 任意のテキストを埋め込みベクトルに変換する CLI ツールです。
*weather-notificator*: 指定した都市の天気予報をDiscord Webhookを通じて通知するCLIツールです。
*web-scraper*: ブラウザを自動操作し、WebページのDOMツリーを取得するCLIツールです。
*withings*: Withings Public Health Data API の OAuth フローに沿って認可 URL の生成、アクセストークン／リフレッシュトークンの取得、そして日次のヘルスデータ取得を行う CLI ツールです。
*yaml-parser*: YAMLファイルや直接指定したYAML文字列を解析して、構造化データをJSONとして出力するCLIツールです。設定値の確認やテストデータの検証に利用できます。
*youtube-downloader*: YouTube動画やプレイリストをダウンロードするためのCLIツールです。
*zip-compressor*: ファイルやディレクトリをZip形式で圧縮・展開するCLIツールです。

## MCP tools
*arithmetic_calculator*: Performs arithmetic operations for two or multiple numbers, and evaluates whether a file line count exceeds a threshold.
*brave_search*: Provides Brave Web Search and Brave Local Search tools.
*context7*: Searches Context7 libraries and fetches documentation by Context7-compatible library ID.
*datetime_calculator*: Provides datetime calculation, time value extraction, and unit-converted summation tools.
*duckduckgo_search*: Provides DuckDuckGo web search and instant answer tools.
*everart*: Generates images via EverArt models and returns links to view generated images.
*figma*: Provides tools to inspect Figma node data and download node images.
*filesystem*: Provides filesystem read, write, list, search, move, and metadata operations within allowed directories.
*gdrive*: Provides Google Drive file access and content retrieval tools.
*git_commit_history_retriever*: Retrieves Git commit history, with optional detailed commit metadata.
*git_diff_recorder*: Provides a tool to retrieve Git diff output from a target repository.
*github*: Provides GitHub repository, issue, pull request, and search operations.
*http_request*: Provides a tool to make HTTP requests to external APIs.
*notion_sync*: Provides tools to patch Markdown content and web clips into Notion pages.
*open_weather_map*: Provides weather forecast retrieval for a specified city using OpenWeather API.
*ops_for_golang*: Provides Go development helpers for test coverage, coverage reports, and `go run` execution.
*persona_extraction*: Extracts and summarizes structured persona information for characters.
*plan*: Provides a structured planning helper tool to update multi-step plans.
*postgresql*: Provides PostgreSQL tools for read-only queries, schema inspection, table listing, and table dump operations.
*sequentialthinking*: Provides a structured sequential-thinking tool for dynamic, revisable problem solving.
*service_implementing_viewer*: Provides a tool to inspect service implementation status across target directories.
*shell*: Provides a shell execution tool compatible with Codex CLI shell-style command execution.
*timezone*: Provides timezone utilities to get current time, convert timezones, and list available timezones.
*util*: 
*weather_notificator*: Sends weather forecast notifications to Discord for a specified city.
*youtube_transcript*: Extracts transcript text from a YouTube video URL or ID.

## GRPC/HANDLERS tools
*weather_notificator*: 天気予報をDiscordに通知するgRPCハンドラーを提供します。

## HTTP/HANDLERS tools
*cron_workflow*: Provides GUI for CRON workflow.
*weather_notificator*: 天気予報をDiscordに通知するHTTPハンドラを提供します。
