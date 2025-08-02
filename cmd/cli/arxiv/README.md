# ArXiv CLI Tool

ArXivから論文を検索・取得するためのコマンドラインツールです。

## 概要

このツールは、ArXiv APIを使用して学術論文の検索と取得を行います。キーワード検索、ID指定取得、高度な検索クエリに対応しており、結果はJSON形式で出力されます。

## 機能

- **論文検索**: キーワードによる論文検索
- **ID指定取得**: ArXiv IDによる直接取得
- **高度な検索**: 著者、タイトル、カテゴリ等の詳細検索
- **ページング**: 開始位置と最大結果数の指定
- **ソート機能**: 関連度、更新日時、投稿日時によるソート
- **JSON出力**: 構造化されたデータ形式での結果出力

## インストール

```bash
# リポジトリをクローン
git clone <repository-url>
cd devbox

# 依存関係をインストール
go mod tidy

# ビルド
go build -o arxiv-cli ./cmd/cli/arxiv
```

## 使用方法

### 基本的な使用方法

```bash
# ヘルプを表示
go run ./cmd/cli/arxiv -help

# 論文検索
go run ./cmd/cli/arxiv -operation search -query "quantum computing" -max_results 5

# ID指定取得
go run ./cmd/cli/arxiv -operation get_by_id -ids "2208.00733,2301.00001"
```

### オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-operation` | `-o` | 操作タイプ (search, get_by_id) | - |
| `-query` | `-q` | 検索クエリ (search操作用) | - |
| `-ids` | `-i` | カンマ区切りのarXiv IDリスト (get_by_id操作用) | - |
| `-start` | `-s` | 開始位置（0ベース） | 0 |
| `-max_results` | `-m` | 最大結果数 | 10 |
| `-sort_by` | - | ソート基準 (relevance, lastUpdatedDate, submittedDate) | relevance |
| `-sort_order` | - | ソート順 (ascending, descending) | descending |
| `-help` | `-h` | ヘルプを表示 | - |

## 使用例

### 1. 基本的な検索

```bash
# "quantum computing"で検索
go run ./cmd/cli/arxiv -operation search -query "quantum computing"
```

### 2. 結果数を指定した検索

```bash
# 最大20件の結果を取得
go run ./cmd/cli/arxiv -operation search -query "machine learning" -max_results 20
```

### 3. ページング

```bash
# 10件目から5件を取得
go run ./cmd/cli/arxiv -operation search -query "deep learning" -start 10 -max_results 5
```

### 4. ソート指定

```bash
# 最新更新日順で検索
go run ./cmd/cli/arxiv -operation search -query "neural networks" -sort_by lastUpdatedDate -sort_order descending
```

### 5. ID指定取得

```bash
# 単一IDで取得
go run ./cmd/cli/arxiv -operation get_by_id -ids "2208.00733"

# 複数IDで取得
go run ./cmd/cli/arxiv -operation get_by_id -ids "2208.00733,2301.00001,1234.5678"
```

### 6. 高度な検索クエリ

```bash
# タイトルで検索
go run ./cmd/cli/arxiv -operation search -query "ti:\"quantum computing\""

# 著者で検索
go run ./cmd/cli/arxiv -operation search -query "au:einstein"

# 要約で検索
go run ./cmd/cli/arxiv -operation search -query "abs:\"machine learning\""

# カテゴリで検索
go run ./cmd/cli/arxiv -operation search -query "cat:cs.AI"

# 複合検索
go run ./cmd/cli/arxiv -operation search -query "ti:quantum AND au:feynman"

# 除外検索
go run ./cmd/cli/arxiv -operation search -query "ti:physics ANDNOT cat:hep-th"
```

## 検索クエリの構文

ArXiv APIは以下の検索フィールドをサポートしています：

| フィールド | 説明 | 例 |
|-----------|------|-----|
| `all` | 全フィールド | `all:electron` |
| `ti` | タイトル | `ti:"quantum computing"` |
| `au` | 著者 | `au:einstein` |
| `abs` | 要約 | `abs:"machine learning"` |
| `co` | コメント | `co:"preliminary results"` |
| `jr` | ジャーナル参照 | `jr:"Phys Rev"` |
| `cat` | カテゴリ | `cat:cs.AI` |
| `rn` | レポート番号 | `rn:hep-th/9901001` |
| `id` | ID | `id:0706.0001` |

### 論理演算子

- `AND`: 両方の条件を満たす
- `OR`: いずれかの条件を満たす
- `ANDNOT`: 最初の条件を満たし、2番目の条件を満たさない

### 引用符

スペースを含む検索語句は引用符で囲みます：
```
ti:"quantum computing"
```

## 出力形式

結果はJSON形式で出力されます：

```json
{
  "operation": "search",
  "papers": [
    {
      "id": "2208.00733v1",
      "title": "The Rise of Quantum Internet Computing",
      "authors": ["Seng W. Loke"],
      "abstract": "This article highlights quantum Internet computing...",
      "published": "2022-08-01T10:36:13Z",
      "updated": "2022-08-01T10:36:13Z",
      "categories": ["cs.ET", "cs.DC"],
      "primary_category": "cs.ET",
      "comment": "Optional comment",
      "journal_ref": "Optional journal reference",
      "pdf_url": "http://arxiv.org/pdf/2208.00733v1",
      "abstract_url": "http://arxiv.org/abs/2208.00733v1"
    }
  ],
  "query_info": {
    "search_query": "quantum computing",
    "id_list": null,
    "start": 0,
    "max_results": 10,
    "sort_by": "relevance",
    "sort_order": "descending"
  },
  "total_count": 1
}
```

## エラーハンドリング

ツールは以下のエラー状況を適切に処理します：

- **ネットワークエラー**: ArXiv APIへの接続失敗
- **HTTPエラー**: サーバーエラー（5xx）やクライアントエラー（4xx）
- **XML解析エラー**: ArXiv APIからの不正なレスポンス
- **APIエラー**: ArXiv APIからのエラーメッセージ
- **設定エラー**: 無効なパラメータや必須パラメータの不足

## 制限事項

- ArXiv APIの制限により、最大結果数は30,000件です
- APIレート制限に注意してください（通常は1秒間に3リクエスト）
- 大量のデータを取得する場合は、適切な間隔を空けることを推奨します

## 開発

### テストの実行

```bash
# 全テストを実行
go test -v ./internal/arxiv/...

# カバレッジ付きでテスト実行
go test -coverprofile=coverage.out ./internal/arxiv/...

# カバレッジレポートを表示
go tool cover -html=coverage.out -o coverage.html
```

### プロジェクト構造

```
cmd/cli/arxiv/
├── main.go                    # CLIエントリーポイント
└── README.md                  # このファイル

internal/arxiv/
├── config/                    # 設定管理
│   ├── config.go             # 設定構造体と検証
│   ├── interfaces.go         # インターフェース定義
│   ├── flag_parser.go        # フラグパーサー実装
│   └── config_test.go        # 設定テスト
└── usecases/                  # ビジネスロジック
    ├── services.go           # ArXivサービス実装
    └── services_test.go      # サービステスト
```

## ライセンス

このプロジェクトは適切なライセンスの下で公開されています。

## 貢献

バグ報告や機能要求は、GitHubのIssueを通じてお知らせください。プルリクエストも歓迎します。

## 参考資料

- [ArXiv API Documentation](https://arxiv.org/help/api)
- [ArXiv Category Taxonomy](https://arxiv.org/category_taxonomy)
- [ArXiv Identifier Scheme](https://arxiv.org/help/arxiv_identifier)
