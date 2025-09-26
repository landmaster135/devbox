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
| data-converter                              | ✅  | ❌️  | ❌️  | ❌️ |
| datetime-calculator                         | ✅  | ✅  | ❌️  | ❌️ |
| db-server-sync                              | ✅  | ❌️  | ❌️  | ❌️ |
| depends-visualizer                          | ✅  | ❌️  | ❌️  | ❌️ |
| diff-dreamer                                | ✅  | ❌️  | ❌️  | ❌️ |
| discord-webhook                             | ✅  | ❌️  | ❌️  | ❌️ |
| duckduckgo-search                           | ❌️  | ✅  | ❌️  | ❌️ |
| env-loader                                  | ✅  | ❌️  | ❌️  | ❌️ |
| everart                                     | ❌️  | ✅  | ❌️  | ❌️ |
| exif-mirror                                 | ✅  | ❌️  | ❌️  | ❌️ |
| exif-modifier                               | ✅  | ❌️  | ❌️  | ❌️ |
| exif-viewer                                 | ✅  | ❌️  | ❌️  | ❌️ |
| figma                                       | ❌️  | ✅  | ❌️  | ❌️ |
| file-character-replacer                     | ✅  | ❌️  | ❌️  | ❌️ |
| file-maneuver                               | ✅  | ❌️  | ❌️  | ❌️ |
| file-processor                              | ✅  | ❌️  | ❌️  | ❌️ |
| filesystem                                  | ❌️  | ✅  | ❌️  | ❌️ |
| gcloud-genset-init                          | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-logging                       | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-monitoring                    | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-monitoring-dashboard          | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-genset-secret                        | ✅  | ❌️  | ❌️  | ❌️ |
| gcloud-wrapper-workload-identity-federation | ✅  | ❌️  | ❌️  | ❌️ |
| gdrive                                      | ❌️  | ✅  | ❌️  | ❌️ |
| git-commit-history-retriever                | ✅  | ✅  | ❌️  | ❌️ |
| git-diff-recorder                           | ✅  | ✅  | ❌️  | ❌️ |
| git-info-retriever                          | ✅  | ❌️  | ❌️  | ❌️ |
| git-pre-commit-hooks                        | ✅  | ❌️  | ❌️  | ❌️ |
| github                                      | ✅  | ✅  | ❌️  | ❌️ |
| goo-scraper                                 | ✅  | ❌️  | ❌️  | ❌️ |
| grpc-request                                | ✅  | ❌️  | ❌️  | ❌️ |
| http-request                                | ✅  | ✅  | ❌️  | ❌️ |
| image-converter                             | ✅  | ❌️  | ❌️  | ❌️ |
| image-filterer                              | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer                               | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer-for-screenshot                | ✅  | ❌️  | ❌️  | ❌️ |
| image-renamer-with-exif                     | ✅  | ❌️  | ❌️  | ❌️ |
| image-rotator                               | ✅  | ❌️  | ❌️  | ❌️ |
| image-trim-describer                        | ✅  | ❌️  | ❌️  | ❌️ |
| image-trimmer                               | ✅  | ❌️  | ❌️  | ❌️ |
| iso8601-converter                           | ✅  | ❌️  | ❌️  | ❌️ |
| json-file-merger                            | ✅  | ❌️  | ❌️  | ❌️ |
| json-formatter-for-agent-interaction        | ✅  | ❌️  | ❌️  | ❌️ |
| json-iso8601-converter                      | ✅  | ❌️  | ❌️  | ❌️ |
| json-modifier                               | ✅  | ❌️  | ❌️  | ❌️ |
| json-timestamp-modifier                     | ✅  | ❌️  | ❌️  | ❌️ |
| kana-converter                              | ✅  | ❌️  | ❌️  | ❌️ |
| mcp-remote                                  | ✅  | ❌️  | ❌️  | ❌️ |
| memory                                      | ✅  | ❌️  | ❌️  | ❌️ |
| movie-converter-for-gif                     | ✅  | ❌️  | ❌️  | ❌️ |
| movie-converter-for-webm                    | ✅  | ❌️  | ❌️  | ❌️ |
| notion-blog-content-extractor               | ✅  | ❌️  | ❌️  | ❌️ |
| notion-sync                                 | ✅  | ✅  | ❌️  | ❌️ |
| ocr-executor                                | ✅  | ❌️  | ❌️  | ❌️ |
| ocr-executor-with-ai                        | ✅  | ❌️  | ❌️  | ❌️ |
| open-weather-map                            | ✅  | ✅  | ❌️  | ❌️ |
| ops-for-golang                              | ✅  | ✅  | ❌️  | ❌️ |
| pdf-encrypter                               | ✅  | ❌️  | ❌️  | ❌️ |
| pdf-merger                                  | ✅  | ❌️  | ❌️  | ❌️ |
| postgresql                                  | ✅  | ✅  | ❌️  | ❌️ |
| qdrant                                      | ✅  | ❌️  | ❌️  | ❌️ |
| script-generator-to-build                   | ✅  | ❌️  | ❌️  | ❌️ |
| sequentialthinking                          | ❌️  | ✅  | ❌️  | ❌️ |
| service-implementing-viewer                 | ✅  | ✅  | ❌️  | ❌️ |
| shell                                       | ❌️  | ✅  | ❌️  | ❌️ |
| steam                                       | ✅  | ❌️  | ❌️  | ❌️ |
| timezone                                    | ❌️  | ✅  | ❌️  | ❌️ |
| unit-converter                              | ✅  | ❌️  | ❌️  | ❌️ |
| util                                        | ❌️  | ✅  | ❌️  | ❌️ |
| valkey                                      | ✅  | ❌️  | ❌️  | ❌️ |
| weather-notificator                         | ✅  | ✅  | ✅  | ✅ |
| yaml-parser                                 | ✅  | ❌️  | ❌️  | ❌️ |
| youtube-downloader                          | ✅  | ❌️  | ❌️  | ❌️ |
| youtube-transcript                          | ❌️  | ✅  | ❌️  | ❌️ |
| zip-compressor                              | ✅  | ❌️  | ❌️  | ❌️ |

## 統計情報

- **総サービス数**: 84
- **CLIツール実装数**: 73
- **MCPツール実装数**: 24
- **gRPCハンドラ実装数**: 1
- **HTTPハンドラ実装数**: 1
- **CLIのみ実装**: 60
- **MCPのみ実装**: 11
- **gRPCハンドラのみ実装**: 0
- **HTTPハンドラのみ実装**: 0
- **CLI+MCP両方実装**: 13
- **全て実装済み**: 1

## 注意事項

- ✅: 実装済み
- ❌️: 未実装
- この状況は`service-implementing-viewer`ツールを使用して自動生成されています

## ドキュメント更新手順

このドキュメントを最新の状態に更新するには、以下の手順を実行してください。`$HOME`は、実行環境によって異なることでしょう：

### 1. CLIツールを使用した更新

```bash
cd $HOME/devbox
./pkg/bash/generate_service_implementing_table.sh
```

### 2. MCPツールを使用した更新
MCPツール実行例：
```json
{
  "server_name": "service_implementing_viewer",
  "tool_name": "get_service_implementing_status",
  "arguments": {
    "root_dir": "/home/user/devbox",
    "target_dirs": "cmd/cli,cmd/mcp"
  }
}
```

### 4. 更新頻度の推奨

- 新しいサービスを追加した際
- 既存サービスのMCP化を完了した際
- 月次での定期更新
