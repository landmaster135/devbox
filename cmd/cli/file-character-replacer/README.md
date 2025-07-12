# File Character Replacer

ファイルまたはディレクトリ内のファイルに対して、指定した文字列を別の文字列に置換するCLIツールです。

## 機能

- **単一ファイル置換**: 指定したファイルの文字列を置換
- **ディレクトリ内一括置換**: ディレクトリ内の全ファイルを再帰的に処理
- **文字エンコーディング対応**: UTF-8、Shift_JIS、EUC-JP、ISO-2022-JPに対応
- **バックアップ機能**: 元ファイルのバックアップを自動作成（タイムスタンプ付き）
- **ドライラン機能**: 実際の変更を行わず、変更予定を表示
- **テキストファイル自動判定**: 拡張子に基づいてテキストファイルのみを処理
- **詳細なログ出力**: 処理結果とエラー情報を表示

## インストール

```bash
# devboxディレクトリでビルド
cd devbox
go build -o bin/file-character-replacer ./cmd/cli/file-character-replacer
```

## 使用方法

### 基本構文

```bash
file-character-replacer [オプション]
```

### オプション

| オプション | 型 | 必須 | デフォルト | 説明 |
|-----------|---|------|-----------|------|
| `-target` | string | ✓ | - | 対象パス（ファイルまたはディレクトリ） |
| `-from` | string | ✓ | - | 置換元文字列 |
| `-to` | string | ✓ | - | 置換先文字列 |
| `-encoding` | string | - | utf-8 | 文字エンコーディング |
| `-recursive` | bool | - | false | ディレクトリの場合、再帰的に処理 |
| `-backup` | bool | - | false | バックアップファイルを作成 |
| `-backup-dir` | string | - | - | バックアップディレクトリ（ディレクトリ処理時は必須） |
| `-dry-run` | bool | - | false | 実際の変更を行わず、変更予定を表示のみ |

### 対応エンコーディング

- `utf-8`: UTF-8（デフォルト）
- `shift_jis`: Shift_JIS
- `euc-jp`: EUC-JP
- `iso-2022-jp`: ISO-2022-JP

### 対応ファイル拡張子

以下の拡張子のファイルがテキストファイルとして処理されます：

```
.txt, .md, .go, .py, .js, .ts, .html, .css, .xml, .json,
.yaml, .yml, .toml, .ini, .cfg, .conf, .log, .sql, .sh,
.bat, .ps1, .php, .rb, .java, .c, .cpp, .h, .hpp,
.cs, .vb, .pl, .r, .scala, .kt, .swift, .dart, .rs
```

## 使用例

```bash
# 単一ファイルの文字列置換
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./test.txt -from="old" -to="new"
# バックアップ付きで置換
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./test.txt -from="old" -to="new" -backup
# バックアップディレクトリを指定
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go  -target=./test.txt -from="old" -to="new" -backup -backup-dir=./backups
# ディレクトリ内のファイルを再帰的に処理
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./src -from="TODO" -to="DONE" -recursive
# バックアップ付きで再帰的に処理
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./src -from="TODO" -to="DONE" -recursive -backup
# エスケープが必要なケース
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=$HOME/devbox/pkg/dos/test_file.bat -from=".\\pkg\\bin\\" -to=".\\pkg\\bin\cli\\" -encoding=shift_jis
# Shift_JISファイルの処理
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./data.txt -from="古い" -to="新しい" -encoding=shift_jis
# EUC-JPファイルの処理
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./legacy.txt -from="旧" -to="新" -encoding=euc-jp
# ドライラン（事前確認）
go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./project -from="debug" -to="release" -recursive -dry-run
```

## 出力例

### 成功時の出力

```
=== ファイル文字列置換結果 ===
処理されたファイル数: 2
置換された箇所数: 3

=== 処理詳細 ===
ディレクトリ内のファイル数: 5
ファイル: src/main.go, 置換回数: 2
バックアップを作成しました: src/main.go
ファイル: src/config.go, 置換回数: 1
バックアップを作成しました: src/config.go
```

### ドライラン時の出力

```
=== ファイル文字列置換結果 ===
処理されたファイル数: 2
置換された箇所数: 3
モード: ドライラン（実際の変更は行われていません）

=== 処理詳細 ===
ディレクトリ内のファイル数: 5
ファイル: src/main.go, 置換回数: 2
[ドライラン] src/main.go の置換をスキップしました
ファイル: src/config.go, 置換回数: 1
[ドライラン] src/config.go の置換をスキップしました
```

### エラー時の出力

```
エラー: --target は必須パラメータです
使用方法: ./bin/file-character-replacer [オプション]
...
```

## バックアップファイル

`-backup` オプションを使用すると、元ファイルのバックアップが作成されます。

### バックアップの動作

**単一ファイル処理時:**
- `-backup-dir`が指定されていない場合：元ファイルと同じディレクトリにバックアップを作成
- `-backup-dir`が指定されている場合：指定されたディレクトリにバックアップを作成

**ディレクトリ処理時:**
- `-backup-dir`の指定が**必須**
- 元のディレクトリ構造を保持してバックアップディレクトリ内に再現

### バックアップファイル名の形式

```
{元ファイル名}.backup_{YYYYMMDD_HHMMSS}
```

### バックアップ構造の例

**単一ファイル処理時（バックアップディレクトリ指定なし）:**
```
./test.txt → ./test.txt.backup_20250712_162007
```

**単一ファイル処理時（バックアップディレクトリ指定あり）:**
```
./test.txt → ./backups/test.txt.backup_20250712_162007
```

**ディレクトリ処理時（ディレクトリ構造保持）:**
```
元の構造:
./src/main.go
./src/config/app.json

バックアップ構造（./backups指定時）:
./backups/src/main.go.backup_20250712_162007
./backups/src/config/app.json.backup_20250712_162007
```

## 注意事項

1. **バイナリファイルは処理されません**: 対応する拡張子のテキストファイルのみが処理対象です
2. **バックアップファイルの管理**: バックアップファイルは自動削除されないため、必要に応じて手動で削除してください
3. **大きなファイルの処理**: メモリ上でファイル全体を読み込むため、非常に大きなファイルの処理には注意が必要です
4. **文字エンコーディングの自動判定**: エンコーディングが指定されていない場合、ファイルから推測を試みますが、正確でない場合があります

## アーキテクチャ

このツールはClean Architectureに基づいて設計されています：

- **Domain層**: ドメインモデルとリポジトリインターフェース
- **Usecases層**: ビジネスロジック（依存関係注入を含む）
- **Interfaces層**: ファイル操作とエンコーディング変換の具象実装
- **Config層**: CLI設定管理
- **CLI層**: エントリーポイント

## ライセンス

MIT License
