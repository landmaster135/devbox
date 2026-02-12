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
| `-operation` | 実行する操作。`create-collection` / `describe-collection` / `delete-collection` / `list-collections` / `upsert-texts` / `query-points` / `overwrite-payload` | - |
| `-db-host` | Qdrant gRPC ホスト | `127.0.0.1` |
| `-db-port` | Qdrant gRPC ポート | `6334` |
| `-collection-name` | 操作対象のコレクション名（create/list を除き必須） | - |
| `-help` | ヘルプ表示 | `false` |

### create-collection 固有

| フラグ | 説明 | 既定値 |
|--------|------|--------|
| `-size` | コレクションに登録するベクトル次元 | `4096` |

### describe-collection / delete-collection 固有

| フラグ | 説明 | 既定値 |
|--------|------|--------|
| `-collection-name` | 対象コレクション名（必須） | - |

`describe-collection` は `qdrant.CollectionInfo` を JSON として出力します。`delete-collection` は削除結果のみを出力します。

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

### overwrite-payload 固有

| フラグ | 説明 | 既定値 |
|--------|------|--------|
| `-payload` | `key=value` 形式で指定する上書き後の payload。CLI 上では 1 キーのみ受け付けますが、ユースケース層では複数フィールド分を扱えるようになっています。 | - |
| `-filter-must` | 上書き対象を絞り込む `must` 条件。`key=value` 形式で複数回指定可能。 | - |
| `-filter-must-not` | 除外する条件 (複数指定可)。 | - |
| `-filter-should` | `should` 条件 (複数指定可)。 | - |
| `-filter-min-should` | `should` 条件のうち満たすべき最小件数 (`--filter-should` と併用)。 | `0` |

`overwrite-payload` ではフィルタを指定しない場合、コレクション内の全ポイントが対象になります。`filter-*` フラグで payload 条件に基づくターゲットを段階的に絞り込めます。

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

### コレクション情報の取得
```bash
go run ./cmd/cli/qdrant \
  -operation describe-collection \
  -collection-name documents \
  -db-host 127.0.0.1 -db-port 6334
```

### コレクションの削除
```bash
go run ./cmd/cli/qdrant \
  -operation delete-collection \
  -collection-name documents \
  -db-host 127.0.0.1 -db-port 6334
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

### payload の上書き
```bash
go run ./cmd/cli/qdrant \
  -operation overwrite-payload \
  -collection-name documents \
  -db-host 127.0.0.1 -db-port 6334 \
  -payload status=reviewed
```

```bash
go run ./cmd/cli/qdrant \
  -operation overwrite-payload \
  -collection-name documents \
  -db-host 127.0.0.1 -db-port 6334 \
  -payload status=reviewed \
  -filter-must topic=travel \
  -filter-must-not status=archived \
  -filter-should lang=ja \
  -filter-min-should 1
```

> ℹ️ フィルタを指定しない場合はコレクション全体が対象です。`filter-should` と `filter-min-should` を併用すると、「should 条件を N 件以上満たすポイント」だけに payload を適用できます。
