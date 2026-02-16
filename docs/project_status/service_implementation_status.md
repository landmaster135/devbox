# サービス実装状況

このドキュメントは、devboxプロジェクトにおける各サービスの実装状況を記録しています。

## 実装状況

以下の表は、各サービスが下記のツールとして実装されているかを示しています。
- CLIツール（`cmd/cli`）
- MCPツール（`cmd/mcp`）
- gRPCハンドラ（`grpc/handlers`）
- HTTP REST API ハンドラ（`http/handlers`）

### 実装状況一覧

| service                                     | cli | mcp | grpc/handlers | http/handlers |
| :-----------------------------------------: | :-: | :-: | :-: | :-: |
| anilist                                     | ✅  | ❌️  | ❌️  | ❌️ |
| arithmetic-calculator                       | ✅  | ✅  | ❌️  | ❌️ |
| arxiv                                       | ✅  | ❌️  | ❌️  | ❌️ |
| base64-extractor                            | ✅  | ❌️  | ❌️  | ❌️ |
| brave-search                                | ❌️  | ✅  | ❌️  | ❌️ |
| claude-code-usage                           | ✅  | ❌️  | ❌️  | ❌️ |
| code-analyzer                               | ✅  | ❌️  | ❌️  | ❌️ |
| color-code-converter                        | ✅  | ❌️  | ❌️  | ❌️ |
| context7                                    | ✅  | ✅  | ❌️  | ❌️ |
| cron-workflow                               | ✅  | ❌️  | ❌️  | ✅ |
| data-converter                              | ✅  | ❌️  | ❌️  | ❌️ |
| datetime-calculator                         | ✅  | ✅  | ❌️  | ❌️ |
| db-server-sync                              | ✅  | ❌️  | ❌️  | ❌️ |
| depends-visualizer                          | ✅  | ❌️  | ❌️  | ❌️ |
| diff-dreamer                                | ✅  | ❌️  | ❌️  | ❌️ |
| discord-webhook                             | ✅  | ❌️  | ❌️  | ❌️ |
| docker                                      | ✅  | ❌️  | ❌️  | ❌️ |
| duckduckgo-search                           | ❌️  | ✅  | ❌️  | ❌️ |
| env-loader                                  | ✅  | ❌️  | ❌️  | ❌️ |
| everart                                     | ❌️  | ✅  | ❌️  | ❌️ |
| exif-mirror                                 | ✅  | ❌️  | ❌️  | ❌️ |
| exif-modifier                               | ✅  | ❌️  | ❌️  | ❌️ |
| exif-viewer                                 | ✅  | ❌️  | ❌️  | ❌️ |
| figma                                       | ❌️  | ✅  | ❌️  | ❌️ |
| file-character-replacer                     | ✅  | ❌️  | ❌️  | ❌️ |
| file-line-deduper                           | ✅  | ❌️  | ❌️  | ❌️ |
| file-maneuver                               | ✅  | ❌️  | ❌️  | ❌️ |
| filesystem                                  | ✅  | ✅  | ❌️  | ❌️ |
| gcloud-genset-ai                            | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-billing                       | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-cloudsql                      | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-container                     | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-deployment                    | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-dns                           | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-iam                           | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-init                          | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-logging                       | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-monitoring                    | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-monitoring-dashboard          | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-scheduler                     | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-secret                        | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-spanner                       | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-storage                       | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-wrapper-workload-identity-federation | ✅  | ❌️  | ❌️  | ❌️ |
| gdrive                                      | ❌️  | ✅  | ❌️  | ❌️ |
| git-commit-history-retriever                | ✅  | ✅  | ❌️  | ❌️ |
| git-diff-recorder                           | ✅  | ✅  | ❌️  | ❌️ |
| git-info-retriever                          | ✅  | ❌️  | ❌️  | ❌️ |
| git-pre-commit-hooks                        | ✅  | ❌️  | ❌️  | ❌️ |
| github                                      | ✅  | ✅  | ❌️  | ❌️ |
| goo-scraper                                 | ✅  | ❌️  | ❌️  | ❌️ |
| grpc-request                                | ✅  | ❌️  | ❌️  | ❌️ |
| html-sanitizer                              | ✅  | ❌️  | ❌️  | ❌️ |
| http-request                                | ✅  | ✅  | ❌️  | ❌️ |
| image-converter                             | ✅  | ❌️  | ❌️  | ❌️ |
| image-filterer                              | ✅  | ❌️  | ❌️  | ❌️ |
| image-filterer-v2                           | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer                               | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer-for-content                   | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer-for-screenshot                | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer-with-exif                     | ✅  | ❌️  | ❌️  | ❌️ |
| image-rotator                               | ✅  | ❌️  | ❌️  | ❌️ |
| image-trim-describer                        | ✅  | ❌️  | ❌️  | ❌️ |
| image-trimmer                               | ✅  | ❌️  | ❌️  | ❌️ |
| interactive-input                           | ✅  | ❌️  | ❌️  | ❌️ |
| iso8601-converter                           | ✅  | ❌️  | ❌️  | ❌️ |
| json-file-merger                            | ✅  | ❌️  | ❌️  | ❌️ |
| json-formatter-for-agent-interaction        | ✅  | ❌️  | ❌️  | ❌️ |
| json-iso8601-converter                      | ✅  | ❌️  | ❌️  | ❌️ |
| json-modifier                               | ✅  | ❌️  | ❌️  | ❌️ |
| json-timestamp-modifier                     | ✅  | ❌️  | ❌️  | ❌️ |
| kana-converter                              | ✅  | ❌️  | ❌️  | ❌️ |
| machine-info                                | ✅  | ❌️  | ❌️  | ❌️ |
| markdown-crafter                            | ✅  | ❌️  | ❌️  | ❌️ |
| mcp-remote                                  | ✅  | ❌️  | ❌️  | ❌️ |
| memory                                      | ✅  | ❌️  | ❌️  | ❌️ |
| memos                                       | ✅  | ❌️  | ❌️  | ❌️ |
| movie-converter-for-gif                     | ✅  | ❌️  | ❌️  | ❌️ |
| movie-converter-for-webm                    | ✅  | ❌️  | ❌️  | ❌️ |
| notion-blog-content-extractor               | ✅  | ❌️  | ❌️  | ❌️ |
| notion-sync                                 | ✅  | ✅  | ❌️  | ❌️ |
| ocr-executor                                | ✅  | ❌️  | ❌️  | ❌️ |
| ocr-executor-with-ai                        | ✅  | ❌️  | ❌️  | ❌️ |
| ollama                                      | ✅  | ❌️  | ❌️  | ❌️ |
| open-weather-map                            | ✅  | ✅  | ❌️  | ❌️ |
| ops-for-golang                              | ✅  | ✅  | ❌️  | ❌️ |
| pdf-encrypter                               | ✅  | ❌️  | ❌️  | ❌️ |
| pdf-merger                                  | ✅  | ❌️  | ❌️  | ❌️ |
| persona-extraction                          | ❌️  | ✅  | ❌️  | ❌️ |
| plan                                        | ❌️  | ✅  | ❌️  | ❌️ |
| postgresql                                  | ✅  | ✅  | ❌️  | ❌️ |
| qdrant                                      | ✅  | ❌️  | ❌️  | ❌️ |
| script-generator-to-build                   | ✅  | ❌️  | ❌️  | ❌️ |
| sequentialthinking                          | ❌️  | ✅  | ❌️  | ❌️ |
| service-implementing-viewer                 | ✅  | ✅  | ❌️  | ❌️ |
| shell                                       | ✅  | ✅  | ❌️  | ❌️ |
| sqlite                                      | ✅  | ❌️  | ❌️  | ❌️ |
| steam                                       | ✅  | ❌️  | ❌️  | ❌️ |
| taskfile                                    | ✅  | ❌️  | ❌️  | ❌️ |
| timezone                                    | ❌️  | ✅  | ❌️  | ❌️ |
| unit-converter                              | ✅  | ❌️  | ❌️  | ❌️ |
| util                                        | ❌️  | ✅  | ❌️  | ❌️ |
| valkey                                      | ✅  | ❌️  | ❌️  | ❌️ |
| vector-embedding                            | ✅  | ❌️  | ❌️  | ❌️ |
| weather-notificator                         | ✅  | ✅  | ✅  | ✅ |
| web-scraper                                 | ✅  | ❌️  | ❌️  | ❌️ |
| withings                                    | ✅  | ❌️  | ❌️  | ❌️ |
| yaml-parser                                 | ✅  | ❌️  | ❌️  | ❌️ |
| youtube-downloader                          | ✅  | ❌️  | ❌️  | ❌️ |
| youtube-transcript                          | ❌️  | ✅  | ❌️  | ❌️ |
| zip-compressor                              | ✅  | ❌️  | ❌️  | ❌️ |

### 統計情報

- **総サービス数**: 111
- **CLIツール実装数**: 100
- **MCPツール実装数**: 26
- **gRPCハンドラ実装数**: 1
- **HTTPハンドラ実装数**: 2
- **CLIのみ実装**: 84
- **MCPのみ実装**: 11
- **gRPCハンドラのみ実装**: 0
- **HTTPハンドラのみ実装**: 0
- **CLI+MCP両方実装**: 15
- **全て実装済み**: 1

## 注意事項

- ✅: 実装済み
- ❌️: 未実装
- この状況は`service-implementing-viewer`ツールを使用して自動生成されています

## ドキュメント更新手順

このドキュメントを最新の状態に更新するには、以下の手順を実行してください。`$HOME`は、実行環境によって異なることでしょう：

### 1. CLIツールを使用した更新
ドキュメントの更新
```bash
cd $HOME/devbox
go run cmd/cli/service-implementing-viewer/main.go \
  -operation=write \
  -root-dir=$HOME/devbox/cmd \
  -target-dirs=cli,mcp,grpc/handlers,http/handlers \
  -write-file=./docs/service_implementation_status.md
```

プレビューのみの場合
```bash
cd $HOME/devbox
go run cmd/cli/service-implementing-viewer/main.go \
  -operation=output \
  -root-dir=$HOME/devbox/cmd \
  -target-dirs=cli,mcp,grpc/handlers,http/handlers
```

シェルファイルからも確認可能
```bash
./scripts/ops/generate_service_implementing_table.sh
```

### 2. MCPツールを使用した更新
MCPツール実行例：
```json
{
  "server_name": "service_implementing_viewer",
  "tool_name": "get_service_implementing_status",
  "arguments": {
    "operation": "output",
    "root_dir": "/home/user/devbox",
    "target_dirs": "cmd/cli,cmd/mcp"
  }
}
```

### 3. 更新頻度の推奨

- 新しいサービスを追加した際
- 既存サービスのMCP化を完了した際
- 月次での定期更新

## 本ドキュメントのロケーションが変わった際の対応

本ドキュメントのロケーションが変わった際は、以下の対応を行ってください：
- [.github/workflows/test_integration.yml](../../.github/workflows/test_integration.yml)の`DOCS_FILE_FOR_SERVICE_IMPLEMENTATION_STATUS`変数を更新する
