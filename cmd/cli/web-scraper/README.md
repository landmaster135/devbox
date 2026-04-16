# Web Scraper

ブラウザを自動操作し、WebページのDOMツリーを取得するCLIツールです。

## 機能
- `get_dom_tree`: 指定したURLを開き、DOMツリーから`<main>`要素のみ（外側のタグを含む）を抽出して出力します。`<faceplate-tracker>`要素を自動で除去し、全要素から`class`/`style`属性を削除、さらに連続する空行を1行に圧縮します。
- `get_meta_props`: 指定したURLのDOMツリーを取得してから最終的なアクセスURL（リダイレクト後を含む）とページタイトルを抽出し、JSONで返します。
- リダイレクトや動的描画に備え、ページを開いてからDOMを取得するまでの待機秒数を指定できます。

## ビルド
```bash
./scripts/build_web_scraper.sh
```
または個別に:
```bash
go build -o bin/web-scraper ./cmd/cli/web-scraper
```

## 使い方
```bash
go run ./cmd/cli/web-scraper \
  -operation=get_dom_tree \
  -url="https://example.com" \
  -wait-seconds=5

go run ./cmd/cli/web-scraper \
  -operation=get_meta_props \
  -url="https://example.com" \
  -wait-seconds=5
```

### オプション
| フラグ | 必須 | 説明 |
| --- | --- | --- |
| `-operation` | ✅ | 実行する操作。`get_dom_tree` または `get_meta_props` を指定。 |
| `-url` | ✅ | 対象ページのURL。操作に関わらず必須。 |
| `-wait-seconds` | 任意 | ページロード完了後にDOM抽出を行うまでの待機時間(秒)。リダイレクトや動的描画待ちに使用。 |
| `-output-file` | 任意 (`get_dom_tree`のみ) | `get_dom_tree`で取得したDOMを新規ファイルとして保存するパス。既存ファイルがある場合はエラー。`get_meta_props`では無視されます。 |

## 実行例
```bash
# 直ちにDOMを取得
go run ./cmd/cli/web-scraper -operation=get_dom_tree -url=https://example.com

# 3秒待ってからDOMを取得
go run ./cmd/cli/web-scraper -operation=get_dom_tree -url=https://news.ycombinator.com -wait-seconds=3

# タイトルと最終URLのみ取得
go run ./cmd/cli/web-scraper -operation=get_meta_props -url=https://example.com
```

## 出力
`get_dom_tree`は標準出力へ`<main>`要素のHTML文字列をそのまま書き出します（`<faceplate-tracker>`除去・`class`/`style`属性除去済み、連続空行は1行に圧縮）。`-output-file`を指定した場合はファイルへ保存し、成功メッセージのみ標準出力に表示します。

`get_meta_props`は以下のようなJSONを1行で標準出力します。

```json
{
  "url": "https://example.com/page",
  "title": "Example Domain"
}
```

対象ページに`<main>`が存在しない場合や内部処理での解析・ファイル出力に失敗した場合はエラーを表示して終了コード `1` を返します。
