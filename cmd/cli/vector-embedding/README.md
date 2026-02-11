# Vector Embedding CLI

任意のテキストを埋め込みベクトルに変換する CLI ツールです。現在は Ollama API を経由した `--operation=ollama` のみをサポートしています。

## 主な特徴
- Ollama の `/api/embed` エンドポイントをラップし、複数テキストを一括でベクトル化
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

成功すると以下のような JSON が表示されます:

```json
{
  "provider": "ollama",
  "model": "all-minilm",
  "embeddings": [[0.12, -0.98, ...], [...]],
  "input_count": 2,
  "dimensions": 1024
}
```

### フラグ一覧

| フラグ | 必須 | 説明 | 例 |
|--------|------|------|----|
| `-operation` | * | 実行モード。現在は `ollama` のみ | `-operation ollama` |
| `-host` | * | Ollama API のホスト名 | `-host 127.0.0.1` |
| `-port` | * | Ollama API のポート番号 | `-port 11434` |
| `-model` | * | 埋め込みに使用するモデル名 | `-model all-minilm` |
| `-input` | * (複数可) | ベクトル化する文字列。複数指定で一括処理 | `-input "サンプル文"` |
| `-timeout` | 任意 | HTTP タイムアウト秒数 (デフォルト: 60s) | `-timeout 120` |
| `-help` | 任意 | ヘルプを表示 | `-help` |

## 使用例
Ollama で 2 件のテキストをベクトル化
```bash
go run ./cmd/cli/vector-embedding \
  -operation ollama \
  -model all-minilm \
  -input "Go は楽しい" \
  -input "ベクトル化とは？"
```

## 注意事項
- Ollama がローカルで起動し、指定ポートで `/api/embed` を受け付けている必要があります。
