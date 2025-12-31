# Web Scraper

`github.com/go-rod/rod` を利用してブラウザを自動操作し、WebページのDOMツリーのうち`<main>`要素のみを取得するCLIツールです。

## 機能
- `get_dom_tree`: 指定したURLを開き、DOMツリーから`<main>`要素のみ（外側のタグを含む）を抽出して出力します。現状は`<faceplate-tracker>`要素を自動的に除去します。
- リダイレクトを想定し、ページを開いてからDOMを取得するまでの待機秒数を指定できます。

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
```

### オプション
| フラグ | 必須 | 説明 |
| --- | --- | --- |
| `-operation` | ✅ | 実行する操作。現在は `get_dom_tree` のみ対応。 |
| `-url` | ✅ (`get_dom_tree`) | DOMを取得したいページのURL。 |
| `-wait-seconds` | 任意 | ページロード完了後にDOM抽出を行うまでの待機時間(秒)。リダイレクトや動的描画待ちに使用。 |

## 実行例
```bash
# 直ちにDOMを取得
go run ./cmd/cli/web-scraper -operation=get_dom_tree -url=https://example.com

# 3秒待ってからDOMを取得
go run ./cmd/cli/web-scraper -operation=get_dom_tree -url=https://news.ycombinator.com -wait-seconds=3
```

## 出力
標準出力へ`<main>`要素のHTML文字列をそのまま書き出します。対象ページに`<main>`が存在しない場合や内部処理での解析に失敗した場合はエラーを表示して終了コード `1` を返します。
