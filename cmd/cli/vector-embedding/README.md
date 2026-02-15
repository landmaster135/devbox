# Vector Embedding CLI

任意のテキストを埋め込みベクトルに変換する CLI ツールです。Ollama (`--operation=ollama`) に加えて OpenAI Embeddings API (`--operation=openai`) を利用できます。

## 主な特徴
- Ollama の `/api/embed` エンドポイントや OpenAI の `POST /v1/embeddings` をラップして複数テキストを一括ベクトル化
- JSON 形式でベクトル・入力数・次元数などのメタ情報を出力
- CLI では設定値のみを受け取り、ネットワーク処理は内部ユースケース層で完結

## インストール / ビルド

```bash
# プロジェクトルートからビルド
go build -o bin/vector-embedding ./cmd/cli/vector-embedding

# もしくは直接実行
go run ./cmd/cli/vector-embedding -operation ollama -model mistral -input "こんにちは"
```

## 使用方法
```bash
go run ./cmd/cli/vector-embedding \
  -operation ollama \
  -host 127.0.0.1 \
  -port 11434 \
  -model all-minilm \
  -input "今日の天気は？" \
  -input "ゴーファーについて教えて"
```

成功すると以下のような JSON が表示されます (provider には実行した operation が入ります):

```json
{
  "provider": "openai",
  "model": "text-embedding-3-small",
  "embeddings": [[0.12, -0.98, ...], [...]],
  "input_count": 2,
  "dimensions": 1536
}
```

### フラグ一覧

| フラグ | 必須 | 対応 operation | 説明 | 例 |
|--------|------|----------------|------|----|
| `-operation` | * | 全て | 実行モード。`ollama` / `openai` を指定 | `-operation openai` |
| `-model` | * | 全て | 埋め込みに使用するモデル名 | `-model text-embedding-3-small` |
| `-input` | * (複数可) | 全て | ベクトル化する文字列。複数指定で一括処理 | `-input "サンプル文"` |
| `-host` | Ollama | `ollama` | Ollama API のホスト名 | `-host 127.0.0.1` |
| `-port` | Ollama | `ollama` | Ollama API のポート番号 | `-port 11434` |
| `-api-key` | OpenAI | `openai` | OpenAI API キー。フラグに直接渡すか、`$OPENAI_API_KEY` を参照 | `-api-key "$OPENAI_API_KEY"` |
| `-timeout` | 任意 | 全て | HTTP タイムアウト秒数 (デフォルト: 60s) | `-timeout 120` |
| `-help` | 任意 | 全て | ヘルプを表示 | `-help` |

## 使用例
Ollama で 2 件のテキストをベクトル化
```bash
go run ./cmd/cli/vector-embedding \
  -operation ollama \
  -model all-minilm \
  -input "Go は楽しい" \
  -input "ベクトル化とは？"
```

OpenAI でベクトル化
```bash
go run ./cmd/cli/vector-embedding \
  -operation openai \
  -api-key "$OPENAI_API_KEY" \
  -model text-embedding-3-small \
  -input "ベクトルデータベースとは？"
```

## 注意事項
- Ollama がローカルで起動し、指定ポートで `/api/embed` を受け付けている必要があります。
- OpenAI を利用する場合はアカウントの API キーと対応モデル (例: `text-embedding-3-small`) を用意してください。
