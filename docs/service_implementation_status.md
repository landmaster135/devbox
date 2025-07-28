# サービス実装状況

このドキュメントは、devboxプロジェクトにおける各サービスの実装状況を記録しています。

## 実装状況一覧

以下の表は、各サービスがCLIツール（`cmd/cli`）とMCPツール（`cmd/mcp`）として実装されているかを示しています。

| service                              | cmd/cli | cmd/mcp |
| :----------------------------------: | :-: | :-: |
| arithmetic-calculator                | ✅  | ✅ |
| base64-extractor                     | ✅  | ❌️ |
| brave-search                         | ❌️  | ✅ |
| claude-code-usage                    | ✅  | ❌️ |
| code-analyzer                        | ✅  | ❌️ |
| context7                             | ✅  | ✅ |
| datetime-calculator                  | ✅  | ✅ |
| depends-visualizer                   | ✅  | ❌️ |
| diff-dreamer                         | ✅  | ❌️ |
| duckduckgo-search                    | ❌️  | ✅ |
| env-loader                           | ✅  | ❌️ |
| everart                              | ❌️  | ✅ |
| exif-mirror                          | ✅  | ❌️ |
| exif-modifier                        | ✅  | ❌️ |
| exif-viewer                          | ✅  | ❌️ |
| figma                                | ❌️  | ✅ |
| file-character-replacer              | ✅  | ❌️ |
| file-maneuver                        | ✅  | ❌️ |
| file-processor                       | ✅  | ❌️ |
| filesystem                           | ❌️  | ✅ |
| gdrive                               | ❌️  | ✅ |
| git-commit-history-retriever         | ✅  | ✅ |
| git-diff-recorder                    | ✅  | ✅ |
| github                               | ❌️  | ✅ |
| goo-scraper                          | ✅  | ❌️ |
| http-request                         | ✅  | ✅ |
| image-converter                      | ✅  | ❌️ |
| image-filterer                       | ✅  | ❌️ |
| image-renamer                        | ✅  | ❌️ |
| image-renamer-for-screenshot         | ✅  | ❌️ |
| image-renamer-with-exif              | ✅  | ❌️ |
| image-rotator                        | ✅  | ❌️ |
| image-trim-describer                 | ✅  | ❌️ |
| image-trimmer                        | ✅  | ❌️ |
| iso8601-converter                    | ✅  | ❌️ |
| json-file-merger                     | ✅  | ❌️ |
| json-formatter-for-agent-interaction | ✅  | ❌️ |
| json-iso8601-converter               | ✅  | ❌️ |
| json-modifier                        | ✅  | ❌️ |
| json-timestamp-modifier              | ✅  | ❌️ |
| kana-converter                       | ✅  | ❌️ |
| movie-converter-for-gif              | ✅  | ❌️ |
| movie-converter-for-webm             | ✅  | ❌️ |
| ocr-executor                         | ✅  | ❌️ |
| pdf-encrypter                        | ✅  | ❌️ |
| pdf-merger                           | ✅  | ❌️ |
| postgresql                           | ❌️  | ✅ |
| script-generator-to-build            | ✅  | ❌️ |
| sequentialthinking                   | ❌️  | ✅ |
| service-implementing-viewer          | ✅  | ✅ |
| shell                                | ❌️  | ✅ |
| timezone                             | ❌️  | ✅ |
| unit-converter                       | ✅  | ❌️ |
| util                                 | ❌️  | ✅ |
| yaml-parser                          | ✅  | ❌️ |
| youtube-transcript                   | ❌️  | ✅ |

## 統計情報

### 実装済みサービス数
- **CLIツール**: 38サービス
- **MCPツール**: 20サービス
- **両方実装済み**: 8サービス

### 実装パターン
- **CLIのみ**: 30サービス
- **MCPのみ**: 12サービス
- **両方実装**: 8サービス

## 注意事項

- ✅: 実装済み
- ❌️: 未実装
- この状況は`service-implementing-viewer`ツールを使用して自動生成されています
- 最終更新: 2025年7月29日

## ドキュメント更新手順

このドキュメントを最新の状態に更新するには、以下の手順を実行してください。`$HOME`は、実行環境によって異なることでしょう：

### 1. MCPツールを使用した自動更新

```bash
# devboxディレクトリに移動
cd $HOME/devbox

# service-implementing-viewerのMCPツールを使用して最新の実装状況を取得
# ClineなどのAIツールから以下のMCPツールを実行：
```

MCPツール実行例：
```json
{
  "server_name": "service_implementing_viewer",
  "tool_name": "get_service_implementing_status",
  "arguments": {
    "root_dir": "$HOME/devbox",
    "target_dirs": "cmd/cli,cmd/mcp"
  }
}
```

### 2. 手動での確認方法

MCPツールが利用できない場合は、以下のコマンドで手動確認できます：

```bash
# CLIツールの一覧を取得
ls -1 $HOME/devbox/cmd/cli/

# MCPツールの一覧を取得
ls -1 $HOME/devbox/cmd/mcp/
```

### 3. ドキュメントの更新

1. 上記の方法で取得した最新の実装状況を基に、このドキュメントの表を更新
2. 統計情報セクションの数値を再計算して更新
3. 最終更新日を現在の日付に変更

### 4. 更新頻度の推奨

- 新しいサービスを追加した際
- 既存サービスのMCP化を完了した際
- 月次での定期更新
