# Diff Dreamer

Diff Dreamerは、difff《ﾃﾞｭﾌﾌ》と同様のUIを提供するテキスト比較ツールです。CLIからHTMLファイルを生成し、ブラウザで差分を可視化できます。

## 概要

- **オフライン動作**: インターネット接続不要
- **ファイル入力対応**: CLIからファイルを直接指定可能
- **プライバシー保護**: データが外部に送信されない
- **カスタマイズ可能**: ローカルファイルなので自由に改変可能
- **レスポンシブデザイン**: デスクトップ・モバイル対応

## 機能

### 差分表示機能
- **行単位の差分**: 追加・削除・変更・同一行を色分け表示
- **文字単位の差分**: 変更行内の詳細な差分をハイライト
- **サイドバイサイド表示**: 左右に並べて比較
- **行番号表示**: 各行に行番号を表示
- **統計情報**: 追加・削除・変更・同一行数の表示

### UI機能
- **リアルタイム比較**: テキスト入力後すぐに比較可能
- **キーボードショートカット**: 
  - `Ctrl+Enter`: 比較実行
  - `Ctrl+Shift+C`: すべてクリア
- **個別クリア**: 左右それぞれのテキストエリアをクリア可能
- **自動リサイズ**: テキストエリアの高さが内容に応じて調整

## インストール

### 前提条件
- Go 1.23.5以上

### ビルド
```bash
cd devbox
go build -o bin/diff-dreamer ./cmd/cli/diff-dreamer/
```

## 使用方法

### 基本的な使用方法
```bash
# 空のテキストエリアでツール起動
go run ./cmd/cli/diff-dreamer

# ファイル指定での比較
go run ./cmd/cli/diff-dreamer -left file1.txt -right file2.txt

# 片方だけファイル指定
go run ./cmd/cli/diff-dreamer -left file1.txt
go run ./cmd/cli/diff-dreamer -right file2.txt

# 出力ファイル名指定
go run ./cmd/cli/diff-dreamer -output my_diff.html
```

### コマンドライン引数

| 引数 | 説明 | デフォルト値 |
|------|------|-------------|
| `-left` | 左側に表示するテキストファイルのパス | なし |
| `-right` | 右側に表示するテキストファイルのパス | なし |
| `-output` | 出力HTMLファイル名 | `diff_dreamer.html` |

### 使用例

#### 例1: 基本的な使用
```bash
go run ./cmd/cli/diff-dreamer
```
- 空のテキストエリアでツールが起動
- ブラウザで手動でテキストを入力して比較

#### 例2: ファイル比較
```bash
go run ./cmd/cli/diff-dreamer -left version1.txt -right version2.txt
```
- 指定したファイルの内容が自動的に読み込まれる
- ブラウザで差分が表示される

#### 例3: カスタム出力ファイル
```bash
go run ./cmd/cli/diff-dreamer -left old.txt -right new.txt -output comparison.html
```
- `comparison.html`として保存される

## 技術仕様

### アーキテクチャ
- **言語**: Go 1.23.5
- **フロントエンド**: Vanilla JavaScript + CSS
- **差分アルゴリズム**: LCS（Longest Common Subsequence）ベース
- **ファイル埋め込み**: `go:embed`を使用してWebアセットを組み込み

### ディレクトリ構造
```
cmd/cli/diff-dreamer/
├── main.go                    # CLIエントリーポイント
├── web/                       # 静的ファイル（embed用）
│   ├── index.html            # メインUI
│   ├── style.css             # スタイルシート
│   └── script.js             # JavaScript（diff処理）
└── README.md                 # このファイル

internal/diff_dreamer/
└── usecases/
    ├── services.go           # HTML生成、ブラウザ起動機能
    └── services_test.go      # テストファイル
```

### 差分アルゴリズム
- **行レベル**: LCS（Longest Common Subsequence）アルゴリズム
- **文字レベル**: 変更行内の詳細差分計算
- **分類**: 追加・削除・変更・同一の4種類

## 開発

### テスト実行
```bash
cd devbox
go test -v ./internal/diff_dreamer/usecases/
```

### テストカバレッジ確認
```bash
cd devbox
go test -coverprofile=coverage.out ./internal/diff_dreamer/usecases/
go tool cover -html=coverage.out -o coverage.html
```

### コード品質チェック
```bash
cd devbox
go vet ./internal/diff_dreamer/...
go fmt ./internal/diff_dreamer/...
```

## 制限事項

- **ファイルサイズ**: 非常に大きなファイルの場合、ブラウザのメモリ制限により動作が重くなる可能性があります
- **文字エンコーディング**: UTF-8以外のエンコーディングは正しく表示されない場合があります
- **バイナリファイル**: テキストファイル専用です

## トラブルシューティング

### ブラウザが自動で開かない
- 手動で生成された`diff_dreamer.html`ファイルをブラウザで開いてください
- ファイルパスがコンソールに表示されます

### ファイルが読み込めない
- ファイルパスが正しいか確認してください
- ファイルの読み取り権限があるか確認してください
- UTF-8エンコーディングのテキストファイルか確認してください

### 差分が正しく表示されない
- ブラウザのJavaScriptが有効になっているか確認してください
- ブラウザのコンソールでエラーが発生していないか確認してください

## ライセンス

このプロジェクトは、devboxプロジェクトの一部として提供されています。

## 関連リンク

- [difff《ﾃﾞｭﾌﾌ》](https://difff.jp/) - 参考にしたオリジナルのWebツール
- [GitHub Issue #4](https://github.com/landmaster135/devbox/issues/4) - 実装要件

## 更新履歴

### v1.0.0 (2025-07-06)
- 初回リリース
- difff《ﾃﾞｭﾌﾌ》風のUI実装
- LCSベースの差分計算アルゴリズム
- ファイル入力対応
- レスポンシブデザイン
- キーボードショートカット対応
