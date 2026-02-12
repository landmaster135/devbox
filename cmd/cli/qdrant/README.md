# Qdrant CLI

Qdrant の gRPC API を叩いてコレクション操作やベクトル upsert / 検索を行う CLI です。`github.com/qdrant/go-client/qdrant` を利用し、埋め込み生成には既存の `vector-embedding` ツールと同じユースケース層を再利用しています。

## ビルド

```bash
# プロジェクトルートから
./scripts/build_qdrant.sh    # まだ無い場合は go build ./cmd/cli/qdrant
```

## 共通オプション

| フラグ | 説明 | 既定値 |
|--------|------|--------|
| `-operation` | 実行する操作。`create-collection` / `list-collections` / `upsert-texts` / `query-points` | - |
| `-db-host` | Qdrant gRPC ホスト | `127.0.0.1` |
| `-db-port` | Qdrant gRPC ポート | `6334` |
| `-collection-name` | 操作対象のコレクション名（create/list を除き必須） | - |
| `-help` | ヘルプ表示 | `false` |

### create-collection 固有

| フラグ | 説明 | 既定値 |
|--------|------|--------|
| `-size` | コレクションに登録するベクトル次元 | `4096` |

### upsert-texts / query-points 固有

| フラグ | 説明 | 既定値 |
|--------|------|--------|
| `-embedding-host` | 埋め込み API のホスト名 | `127.0.0.1` |
| `-embedding-port` | 埋め込み API のポート番号 | `11434` |
| `-embedding-model` | 埋め込みモデル名（`vector-embedding` CLI と同じ指定方法） | - |
| `-input` | 埋め込み対象テキスト。CLI 上では 1 件のみ受け付けますが、ユースケース層には input/payload ペアのリストとして渡され、将来的に複数指定へ拡張できる構成です。 | - |
| `-payload` | `key=value` 形式の payload 条件。CLI 上では 1 件のみ受け付け、内部でリスト化して upsert/payload フィルタに使います。 | - |
| `-limit` | `query-points` の取得件数 | `5` |

> ℹ️ `upsert-texts` では 1 件の `input` と `payload` を **ペアのリスト**（長さ 1）に変換してから `points.Upsert` に渡します。`query-points` でも同様に payload 条件をリスト化して `Filter.Must` に展開します。複数条件を取り扱いたい場合は CLI の拡張のみで済むように実装されています。

## 使用例

### コレクションの作成
```bash
go run ./cmd/cli/qdrant \
  -operation create-collection \
  -db-host 127.0.0.1 \
  -db-port 6334 \
  -collection-name documents \
  -size 4096
```

### コレクション一覧
```bash
go run ./cmd/cli/qdrant -operation list-collections -db-host 127.0.0.1 -db-port 6334
```

### テキストの upsert
```bash
go run ./cmd/cli/qdrant \
  -operation upsert-texts \
  -collection-name documents \
  -db-host 127.0.0.1 -db-port 6334 \
  -embedding-host 127.0.0.1 -embedding-port 11434 \
  -embedding-model nomic-embed-text \
  -input "東京タワーの歴史" \
  -payload topic=travel
```

### 類似ポイントの検索
```bash
go run ./cmd/cli/qdrant \
  -operation query-points \
  -collection-name documents \
  -db-host 127.0.0.1 -db-port 6334 \
  -embedding-host 127.0.0.1 -embedding-port 11434 \
  -embedding-model nomic-embed-text \
  -input "東京の観光名所" \
  -payload topic=travel \
  -limit 5
```

`query-points` では検索ベクトル生成にも `vector-embedding` ユースケースを利用し、Qdrant 側には `WithPayload=enable` をセットして payload を常に返すようにしています。
