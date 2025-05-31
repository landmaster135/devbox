# Claude Code Usage Analyzer

Claude Code使用状況分析ツールのCLIインターフェース

## 概要

このツールは、Claude Codeの使用状況を分析し、トークン使用量とコストをレポートします。
ローカルに保存されたJSONLファイルを解析して、日別やセッション別の使用状況を表示できます。

## 使用方法

### 基本的な使用

```bash
# 日別レポートを表示（デフォルト）
go run /home/nov/devbox/cmd/cli/claude-code-usage/main.go daily

# セッション別レポートを表示
go run /home/nov/devbox/cmd/cli/claude-code-usage/main.go session

# 日付でフィルター
go run /home/nov/devbox/cmd/cli/claude-code-usage/main.go daily --since 20250525 --until 20250530

# JSON形式で出力
go run /home/nov/devbox/cmd/cli/claude-code-usage/main.go daily --json

# カスタムパスを指定
go run /home/nov/devbox/cmd/cli/claude-code-usage/main.go daily --path /custom/path/to/.claude
```

### ビルドして実行

```bash
# ビルド
cd /home/nov/devbox
go build -o claude-code-usage ./cmd/cli/claude-code-usage/main.go

# 実行
./claude-code-usage daily
./claude-code-usage session --json
```

### オプション

- `-s, --since <date>`: 開始日フィルター (YYYYMMDD形式)
- `-u, --until <date>`: 終了日フィルター (YYYYMMDD形式)  
- `-p, --path <path>`: Claudeデータディレクトリのカスタムパス (デフォルト: ~/.claude)
- `-j, --json`: JSON形式で結果を出力
- `-h, --help`: ヘルプメッセージを表示
- `-v, --version`: バージョン情報を表示

## 出力例

### 日別レポート

```
╭──────────────────────────────────────────╮
│                                          │
│  Claude Code Token Usage Report - Daily  │
│                                          │
╰──────────────────────────────────────────╯

┌──────────────┬────────┬─────────┬──────────────┬────────────┬──────────────┬────────────┐
│ Date         │ Input  │ Output  │ Cache Create │ Cache Read │ Total Tokens │ Cost (USD) │
├──────────────┼────────┼─────────┼──────────────┼────────────┼──────────────┼────────────┤
│ 2025-05-30   │    277 │  31,456 │          512 │      1,024 │       33,269 │     $17.58 │
│ 2025-05-29   │    959 │  39,662 │          256 │        768 │       41,645 │     $16.42 │
└──────────────┴────────┴─────────┴──────────────┴────────────┴──────────────┴────────────┘
```

### セッション別レポート

```
╭───────────────────────────────────────────────╮
│                                               │
│  Claude Code Token Usage Report - By Session  │
│                                               │
╰───────────────────────────────────────────────╯

┌─────────────┬────────────┬────────┬─────────┬──────────────┬────────────┬──────────────┬────────────┬───────────────┐
│ Project     │ Session    │ Input  │ Output  │ Cache Create │ Cache Read │ Total Tokens │ Cost (USD) │ Last Activity │
├─────────────┼────────────┼────────┼─────────┼──────────────┼────────────┼──────────────┼────────────┼───────────────┤
│ myproject   │ session-1  │  4,512 │ 350,846 │          512 │      1,024 │      356,894 │    $156.40 │ 2025-05-24    │
└─────────────┴────────────┴────────┴─────────┴──────────────┴────────────┴──────────────┴────────────┴───────────────┘
```

## 機能

- 📊 **日別レポート**: 日付別にトークン使用量とコストを集計表示
- 💬 **セッション別レポート**: 会話セッション別に使用量を表示
- 📅 **日付フィルタリング**: `--since`と`--until`で期間指定
- 📁 **カスタムパス**: Claudeデータディレクトリの場所を指定可能
- 🎨 **美麗な表示**: カラフルなテーブル形式での表示
- 📄 **JSON出力**: 構造化されたJSON形式での出力サポート
- 💰 **コスト追跡**: 各日/セッションのUSDコストを表示
- 🔄 **キャッシュトークン対応**: キャッシュ作成・読み込みトークンを別途表示

## 制限事項

- Claude CodeのローカルJSONLファイルのみを読み取ります
- Web検索、コード実行、画像解析などのAPI使用量は含まれません
- 言語モデルのトークン使用量のみを追跡します

## 必要要件

- Claude Codeの使用履歴ファイル (`~/.claude/projects/**/*.jsonl`)
- Go 1.19以降

## アーキテクチャ

実装は以下のディレクトリ構造に従っています：

```
/home/nov/devbox/internal/independencies/claude_code_usage/
├── cmd/app.go                 # CLIアプリケーションロジック
├── internal/
│   ├── types.go              # データ構造定義
│   ├── dataloader.go         # JSONLファイル解析・読み込み
│   ├── calculator.go         # データ集計・計算
│   ├── formatter.go          # 出力フォーマット（テーブル/JSON）
│   └── calculator_test.go    # テストコード
└── main.go                   # メインエントリーポイント
```
