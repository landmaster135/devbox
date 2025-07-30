# Memory CLI Tool

知識グラフを使用したメモリ管理を行うCLIツールです。エンティティ、リレーション、観察事項を管理し、永続的なメモリ機能を提供します。

## 機能

- **エンティティ管理**: 人、組織、イベントなどのエンティティを作成・削除
- **リレーション管理**: エンティティ間の関係性を定義・削除
- **観察事項管理**: エンティティに関する情報を追加・削除
- **知識グラフ読み取り**: 全体のグラフ構造を表示
- **ノード検索**: 名前、タイプ、観察事項による検索
- **特定ノード取得**: 指定したエンティティとその関係を取得
- **JSON形式**: 構造化されたデータの入出力
- **マルチストレージ対応**: ファイルまたはValkey DBでのデータ保存
- **ファイル永続化**: memory.jsonファイルでのデータ保存
- **Valkey DB対応**: Redis互換のValkey DBでのデータ保存
- **短縮オプション**: 全てのオプションに短縮形を提供

## インストール

```bash
# プロジェクトルートから
go build -o bin/memory ./cmd/cli/memory
```

## 使用方法

### 基本的な使用方法

```bash
./bin/memory -operation create-entities -entities '[{"name":"John","entityType":"person","observations":["speaks English","likes coffee"]}]'
```

出力例:
```json
作成されたエンティティ:
[
  {
    "name": "John",
    "entityType": "person",
    "observations": [
      "speaks English",
      "likes coffee"
    ]
  }
]
```

## オプション

### 基本オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-operation` | `-o` | メモリ操作 (create-entities, create-relations, add-observations, delete-entities, delete-observations, delete-relations, read-graph, search-nodes, open-nodes) | | `-o create-entities` |
| `-storage-type` | `-s` | ストレージタイプ (file, valkey) | | `-s valkey` |
| `-memory-file` | `-f` | メモリファイルパス | ./memory.json | `-f /path/to/memory.json` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

### Valkey接続オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-valkey-host` | `-vh` | Valkeyホスト | localhost | `-vh redis.example.com` |
| `-valkey-port` | `-vp` | Valkeyポート | 6379 | `-vp 6380` |
| `-valkey-password` | `-vpass` | Valkeyパスワード | | `-vpass mypassword` |
| `-valkey-database` | `-vdb` | データベース番号 | 0 | `-vdb 1` |
| `-valkey-key` | `-vk` | Valkeyキー | knowledge_graph:main | `-vk my_graph` |

### エンティティ関連オプション

| オプション | 短縮形 | 説明 | 例 |
|-----------|--------|------|-----|
| `-entities` | `-e` | JSON形式のエンティティ配列 | `-e '[{"name":"John","entityType":"person","observations":["speaks English"]}]'` |
| `-entity-names` | `-en` | 削除対象エンティティ名（カンマ区切り） | `-en "John,Company"` |

### リレーション関連オプション

| オプション | 短縮形 | 説明 | 例 |
|-----------|--------|------|-----|
| `-relations` | `-r` | JSON形式のリレーション配列 | `-r '[{"from":"John","to":"Company","relationType":"works_at"}]'` |

### 観察事項関連オプション

| オプション | 短縮形 | 説明 | 例 |
|-----------|--------|------|-----|
| `-observations` | `-obs` | JSON形式の観察事項配列 | `-obs '[{"entityName":"John","contents":["likes coffee"]}]'` |
| `-deletions` | `-del` | JSON形式の削除対象 | `-del '[{"entityName":"John","observations":["old info"]}]'` |

### 検索・取得関連オプション

| オプション | 短縮形 | 説明 | 例 |
|-----------|--------|------|-----|
| `-query` | `-q` | 検索クエリ | `-q "coffee"` |
| `-names` | `-n` | カンマ区切りの名前リスト | `-n "John,Company"` |

## 使用例

### ストレージタイプの選択

```bash
# ファイルベース（デフォルト）
./bin/memory -operation read-graph

# Valkey DBベース
./bin/memory -operation read-graph -storage-type valkey

# 短縮形でValkey DB使用
./bin/memory -o read-graph -s valkey
```

### Valkey DB使用例

```bash
# 基本的なValkey DB使用
./bin/memory -o create-entities -s valkey -e '[{"name":"John","entityType":"person","observations":["speaks English"]}]'

# カスタムホストとポート
./bin/memory -o read-graph -s valkey -valkey-host redis.example.com -valkey-port 6380

# パスワード認証付き
./bin/memory -o read-graph -s valkey -valkey-password mypassword

# 特定のデータベースを使用
./bin/memory -o read-graph -s valkey -valkey-database 1

# カスタムキーを使用
./bin/memory -o read-graph -s valkey -valkey-key my_custom_graph

# 全ての設定を指定
./bin/memory -o read-graph -s valkey \
  -valkey-host redis.example.com \
  -valkey-port 6380 \
  -valkey-password mypassword \
  -valkey-database 1 \
  -valkey-key production_graph

# 短縮形を使用
./bin/memory -o read-graph -s valkey -vh redis.example.com -vp 6380 -vpass mypass -vdb 1 -vk prod_graph
```

### エンティティ作成

```bash
# 基本的なエンティティ作成
./bin/memory -operation create-entities -entities '[{"name":"John","entityType":"person","observations":["speaks English","likes coffee"]}]'

# 複数のエンティティを同時作成
./bin/memory -operation create-entities -entities '[{"name":"John","entityType":"person","observations":["speaks English"]},{"name":"Company","entityType":"organization","observations":["tech company","located in Tokyo"]}]'

# 短縮形を使用
./bin/memory -o create-entities -e '[{"name":"Alice","entityType":"person","observations":["speaks Japanese","works remotely"]}]'
```

### リレーション作成

```bash
# 基本的なリレーション作成
./bin/memory -operation create-relations -relations '[{"from":"John","to":"Company","relationType":"works_at"}]'

# 複数のリレーションを同時作成
./bin/memory -operation create-relations -relations '[{"from":"John","to":"Company","relationType":"works_at"},{"from":"Alice","to":"Company","relationType":"works_at"}]'

# 短縮形を使用
./bin/memory -o create-relations -r '[{"from":"John","to":"Alice","relationType":"colleague_of"}]'
```

### 観察事項追加

```bash
# 基本的な観察事項追加
./bin/memory -operation add-observations -observations '[{"entityName":"John","contents":["works remotely","enjoys programming"]}]'

# 複数のエンティティに観察事項を追加
./bin/memory -operation add-observations -observations '[{"entityName":"John","contents":["team leader"]},{"entityName":"Company","contents":["growing startup"]}]'

# 短縮形を使用
./bin/memory -o add-observations -obs '[{"entityName":"Alice","contents":["project manager","speaks three languages"]}]'
```

### 知識グラフ読み取り

```bash
# 全体のグラフを表示
./bin/memory -operation read-graph

# 短縮形を使用
./bin/memory -o read-graph
```

### ノード検索

```bash
# キーワードで検索
./bin/memory -operation search-nodes -query "coffee"

# エンティティタイプで検索
./bin/memory -operation search-nodes -query "person"

# 観察事項で検索
./bin/memory -operation search-nodes -query "programming"

# 短縮形を使用
./bin/memory -o search-nodes -q "Tokyo"
```

### 特定ノード取得

```bash
# 単一のノードを取得
./bin/memory -operation open-nodes -names "John"

# 複数のノードを取得
./bin/memory -operation open-nodes -names "John,Company"

# 短縮形を使用
./bin/memory -o open-nodes -n "Alice,Company"
```

### エンティティ削除

```bash
# 単一のエンティティを削除
./bin/memory -operation delete-entities -entity-names "John"

# 複数のエンティティを削除
./bin/memory -operation delete-entities -entity-names "John,Alice"

# 短縮形を使用
./bin/memory -o delete-entities -en "OldCompany"
```

### 観察事項削除

```bash
# 特定の観察事項を削除
./bin/memory -operation delete-observations -deletions '[{"entityName":"John","observations":["old information"]}]'

# 複数のエンティティから観察事項を削除
./bin/memory -operation delete-observations -deletions '[{"entityName":"John","observations":["outdated info"]},{"entityName":"Company","observations":["old status"]}]'

# 短縮形を使用
./bin/memory -o delete-observations -del '[{"entityName":"Alice","observations":["temporary role"]}]'
```

### リレーション削除

```bash
# 特定のリレーションを削除
./bin/memory -operation delete-relations -relations '[{"from":"John","to":"OldCompany","relationType":"works_at"}]'

# 複数のリレーションを削除
./bin/memory -operation delete-relations -relations '[{"from":"John","to":"OldCompany","relationType":"works_at"},{"from":"Alice","to":"OldCompany","relationType":"works_at"}]'

# 短縮形を使用
./bin/memory -o delete-relations -r '[{"from":"John","to":"Alice","relationType":"old_colleague"}]'
```

## データ形式

### エンティティのJSON形式

```json
[
  {
    "name": "John",
    "entityType": "person",
    "observations": [
      "speaks English",
      "likes coffee",
      "works remotely"
    ]
  }
]
```

### リレーションのJSON形式

```json
[
  {
    "from": "John",
    "to": "Company",
    "relationType": "works_at"
  }
]
```

### 観察事項のJSON形式

```json
[
  {
    "entityName": "John",
    "contents": [
      "enjoys programming",
      "team leader"
    ]
  }
]
```

### 削除対象のJSON形式

```json
[
  {
    "entityName": "John",
    "observations": [
      "old information",
      "outdated status"
    ]
  }
]
```

## 出力フォーマット

### エンティティ作成の出力

```json
作成されたエンティティ:
[
  {
    "name": "John",
    "entityType": "person",
    "observations": [
      "speaks English",
      "likes coffee"
    ]
  }
]
```

### リレーション作成の出力

```json
作成されたリレーション:
[
  {
    "from": "John",
    "to": "Company",
    "relationType": "works_at"
  }
]
```

### 観察事項追加の出力

```json
追加された観察事項:
[
  {
    "entityName": "John",
    "addedObservations": [
      "works remotely",
      "enjoys programming"
    ]
  }
]
```

### 知識グラフ読み取りの出力

```json
知識グラフ:
{
  "entities": [
    {
      "name": "John",
      "entityType": "person",
      "observations": [
        "speaks English",
        "likes coffee",
        "works remotely",
        "enjoys programming"
      ]
    },
    {
      "name": "Company",
      "entityType": "organization",
      "observations": [
        "tech company",
        "located in Tokyo"
      ]
    }
  ],
  "relations": [
    {
      "from": "John",
      "to": "Company",
      "relationType": "works_at"
    }
  ]
}
```

### 検索結果の出力

```json
検索結果 (クエリ: coffee):
{
  "entities": [
    {
      "name": "John",
      "entityType": "person",
      "observations": [
        "speaks English",
        "likes coffee"
      ]
    }
  ],
  "relations": null
}
```

### 削除操作の出力

```
エンティティを削除しました: John,Company
観察事項を削除しました
リレーションを削除しました
```

## エラーハンドリング

### 無効な操作タイプ

```bash
./bin/memory -operation invalid
```

```
エラー: 無効な操作タイプです: invalid
```

### 必須パラメータの不足

```bash
./bin/memory
```

```
エラー: 操作タイプが指定されていません
メモリCLIツール

使用方法:
  エンティティ作成:
    ./bin/memory -operation create-entities -entities '[{"name":"John","entityType":"person","observations":["speaks English"]}]'
...
```

### 無効なJSON形式

```bash
./bin/memory -o create-entities -e 'invalid json'
```

```
エラー: エンティティのJSON解析エラー: invalid character 'i' looking for beginning of value
```

### 存在しないエンティティへの操作

```bash
./bin/memory -o add-observations -obs '[{"entityName":"NonExistent","contents":["test"]}]'
```

```
エラー: エンティティが見つかりません: NonExistent
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **依存性注入**: テスタビリティを考慮した設計

### ディレクトリ構造

```
internal/memory/
├── config/              # 設定管理
│   ├── config.go        # 設定構造体とパーサー
│   ├── flag_parser.go   # フラグ解析
│   └── interfaces.go    # インターフェース定義
└── usecases/            # ビジネスロジック
    ├── models.go        # データ構造定義
    ├── memory_manager.go # 知識グラフ管理
    └── services.go      # サービス層
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `encoding/json`, `flag`, `fmt`, `os`など
- **テストフレームワーク**: Go標準のtestingパッケージ

### データ永続化の仕組み

- **ファイル形式**: JSON形式でのデータ保存
- **デフォルトパス**: `./memory.json`
- **カスタムパス**: `-memory-file`オプションで指定可能
- **自動作成**: ディレクトリが存在しない場合は自動作成
- **重複チェック**: エンティティ名、リレーション、観察事項の重複を防止

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/memory ./cmd/cli/memory

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/memory-linux ./cmd/cli/memory
GOOS=windows GOARCH=amd64 go build -o bin/memory.exe ./cmd/cli/memory
GOOS=darwin GOARCH=amd64 go build -o bin/memory-mac ./cmd/cli/memory
```

### テスト

```bash
# 単体テスト
go test ./internal/memory/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/memory/...
go tool cover -html=coverage.out -o coverage.html
```

### 関連ツール

このCLIツールと同じサービスを使用するMCPツールも利用可能です：
- Model Context Protocol経由で同じメモリ管理機能にアクセス可能

## 実用的な使用例

### 個人の知識管理

```bash
# 人物情報の管理
./bin/memory -o create-entities -e '[{"name":"田中太郎","entityType":"person","observations":["プロジェクトマネージャー","英語堪能","東京在住"]}]'

# 会社情報の追加
./bin/memory -o create-entities -e '[{"name":"ABC株式会社","entityType":"organization","observations":["IT企業","従業員500名","渋谷本社"]}]'

# 関係性の定義
./bin/memory -o create-relations -r '[{"from":"田中太郎","to":"ABC株式会社","relationType":"勤務"}]'
```

### プロジェクト管理

```bash
# プロジェクト情報の管理
./bin/memory -o create-entities -e '[{"name":"WebサイトリニューアルPJ","entityType":"project","observations":["期間6ヶ月","予算500万円","チーム5名"]}]'

# プロジェクトメンバーとの関係
./bin/memory -o create-relations -r '[{"from":"田中太郎","to":"WebサイトリニューアルPJ","relationType":"担当"}]'

# 進捗情報の追加
./bin/memory -o add-observations -obs '[{"entityName":"WebサイトリニューアルPJ","contents":["要件定義完了","デザイン進行中"]}]'
```

### 学習記録

```bash
# 技術スキルの管理
./bin/memory -o create-entities -e '[{"name":"Go言語","entityType":"skill","observations":["プログラミング言語","Google開発","並行処理得意"]}]'

# 学習関係の定義
./bin/memory -o create-relations -r '[{"from":"田中太郎","to":"Go言語","relationType":"学習中"}]'

# 学習進捗の記録
./bin/memory -o add-observations -obs '[{"entityName":"Go言語","contents":["基本文法習得","Webアプリ作成経験"]}]'
```

### 情報検索と分析

```bash
# 特定のスキルを持つ人を検索
./bin/memory -o search-nodes -q "英語"

# プロジェクト関連情報を検索
./bin/memory -o search-nodes -q "プロジェクト"

# 特定の人物とその関係を確認
./bin/memory -o open-nodes -n "田中太郎"

# 全体の知識グラフを確認
./bin/memory -o read-graph
```

### データメンテナンス

```bash
# 古い情報の削除
./bin/memory -o delete-observations -del '[{"entityName":"田中太郎","observations":["古い職歴"]}]'

# 終了したプロジェクトの削除
./bin/memory -o delete-entities -en "終了したプロジェクト"

# 変更された関係の更新
./bin/memory -o delete-relations -r '[{"from":"田中太郎","to":"旧会社","relationType":"勤務"}]'
./bin/memory -o create-relations -r '[{"from":"田中太郎","to":"新会社","relationType":"勤務"}]'
```

## ベストプラクティス

### エンティティ命名

- **一意性**: 同じタイプ内で重複しない名前を使用
- **識別性**: 明確で識別しやすい名前を選択
- **一貫性**: 命名規則を統一

### リレーション設計

- **能動態**: 関係性は能動態で記述（例：`works_at`, `manages`, `belongs_to`）
- **明確性**: 関係の方向性が明確になるよう設計
- **標準化**: 同じ種類の関係には統一した名前を使用

### 観察事項管理

- **原子性**: 1つの観察事項には1つの事実のみを記録
- **更新性**: 古い情報は削除し、新しい情報を追加
- **具体性**: 具体的で検索しやすい内容を記録

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
