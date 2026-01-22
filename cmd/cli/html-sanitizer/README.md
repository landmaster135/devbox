# HTML Sanitizer

`internal/html_sanitizer/usecases/sanitizer` に実装されている `SanitizeHTMLBody` をCLIから呼び出し、main/article要素を基点にHTMLをサニタイズします。

## ビルド

```bash
# プロジェクトルートから
go build -o bin/html-sanitizer ./cmd/cli/html-sanitizer
```

| フラグ | 必須 | 説明 |
|--------|------|------|
| `--input-file` | * | サニタイズ対象のHTMLファイルパス（`--html-file`でも可） |
| `--html-file` |  | `--input-file`のエイリアス。既存スクリプトから移行する際に利用できます |
| `--output-file` | * | サニタイズ結果を書き込むファイルパス |
| `--omits-full-body` |  | trueで、サニタイザーがエラーを返した場合に入力HTML全文を書き出しません（既定はfalse） |
| `--help` |  | ヘルプとフラグ一覧を表示して終了します |

## 使用例

```bash
# 成功時はサニタイズ済みHTMLを出力ファイルへ書き込み
go run ./cmd/cli/html-sanitizer --input-file ./sample.html --output-file ./sanitized.html

# エラー時に空文字を書き出したい場合
go run ./cmd/cli/html-sanitizer \
  --input-file ./sample_data/html/broken.html \
  --output-file ./tmp/broken_sanitized.html \
  --omits-full-body
```

標準エラーには処理中のエラー内容が出力され、失敗時は非0コードで終了します。出力ファイルには `SanitizeHTMLBody` が返却した文字列がそのまま保存されるため、必要に応じてさらに別処理へ渡せます。
