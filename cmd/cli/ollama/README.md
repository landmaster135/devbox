# Ollama CLI

Ollama ローカルサーバーの HTTP API を手元から叩くための CLI です。`version`/`list-models`/`embed`/`generate`/`pull` といった代表的なエンドポイント（[公式ドキュメント](https://docs.ollama.com/api)）をカバーし、簡単に API の結果を確認できます。

## 主な機能

- **version**: `/api/version` の結果を JSON で表示
- **list-models**: `/api/tags`（インストール済み）または `/api/ps`（稼働中のみ）を取得
- **embed**: `/api/embed` に複数テキストを送り、埋め込みベクトルを得る
- **generate**: `/api/generate` をストリーミングで呼び出し、生成テキストをまとめて表示
- **pull**: `/api/pull` の進捗を逐次受け取り、%表示付きでサマリ出力

## ビルド

```bash
# プロジェクトルートで実行
go build -o bin/ollama ./cmd/cli/ollama
```

## 使用方法

```bash
go run ./cmd/cli/ollama \
  --operation=list-models \
  --host=127.0.0.1 \
  --port=11434
```

### 共通フラグ

| フラグ | 説明 | 既定値 |
| --- | --- | --- |
| `--operation` | `version` / `list-models` / `embed` / `generate` / `pull` | **必須** |
| `--host` | Ollama サーバーのホスト名 | `127.0.0.1` |
| `--port` | Ollama サーバーのポート番号 | `11434` |
| `--timeout` | HTTP タイムアウト秒数。`embed`/`generate`/`pull` は未指定なら 300 秒、それ以外は 30 秒 | `30` |
| `--help` | ヘルプ表示 | `false` |

### 操作別フラグ

| operation | 追加フラグ | 説明 |
| --- | --- | --- |
| `list-models` | `--running-only` | `true` を指定すると `/api/ps` の結果のみ表示 |
| `embed` | `--model`, `--input` (複数可) | `--input` は最低 1 件必要 |
| `generate` | `--model`, `--prompt` | レスポンスは標準出力へまとめて出力 |
| `pull` | `--model` | 進捗ログを整形して表示 |

## 使用例
バージョン確認
```bash
go run ./cmd/cli/ollama --operation=version | jq
```

モデル一覧
```bash
# インストール済みモデル
go run ./cmd/cli/ollama --operation=list-models

# 稼働中のみ
go run ./cmd/cli/ollama --operation=list-models --running-only
```

埋め込み生成
```bash
go run ./cmd/cli/ollama \
  --operation=embed \
  --model=nomic-embed-text \
  --input="This is a pen." \
  --input="こんにちは世界"
```

テキスト生成
```bash
go run ./cmd/cli/ollama \
  --operation=generate \
  --model=llama3 \
  --prompt="日本語で自己紹介して"
```

モデル取得
```bash
go run ./cmd/cli/ollama \
  --operation=pull \
  --model=llama3:instruct
```

## 出力形式

- `version` / `list-models` / `embed`: `jq` で扱いやすい整形 JSON
- `generate`: ストリームを結合したテキスト（改行はモデル出力に従う）
- `pull`: `downloading 50.0% (50/100)` のような人間向けサマリを行単位で表示

## 参考資料

- [Ollama API Reference / Version](https://docs.ollama.com/api-reference/get-version)
- [Ollama API Reference / Tags](https://docs.ollama.com/api/tags)
- [Ollama API Reference / ps](https://docs.ollama.com/api/ps)
- [Ollama API Reference / Embed](https://docs.ollama.com/api/embed)
- [Ollama API Reference / Generate](https://docs.ollama.com/api/generate)
- [Ollama API Reference / Pull](https://docs.ollama.com/api/pull)
