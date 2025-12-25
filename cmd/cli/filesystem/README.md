# filesystem

ファイル読み書きや検索、移動などの基本操作を安全に行うCLIツールです。MCP版と同じ `internal/filesystem/usecases` を利用し、許可したディレクトリの範囲内で各操作を実行します。

## 主な機能

- **read_file**: テキストファイルを読み取り内容を標準出力に表示
- **write_file**: ファイルの作成または上書き
- **create_directory**: ネストしたディレクトリの作成
- **list_directory**: ディレクトリのフラット一覧を表示
- **directory_tree**: もしくはYAML形式のツリー構造を表示
- **move_file**: ファイル／ディレクトリの移動やリネーム
- **search_files**: 正規表現に一致する項目を再帰検索
- **get_file_info**: サイズや権限などのメタ情報を整形表示
- **list_allowed_directories**: 現在許可されているディレクトリを確認

## インストール

```bash
# プロジェクトルートでビルド
go build -o bin/filesystem ./cmd/cli/filesystem
```

## 基本的な使い方

```bash
# operationと必要なフラグを指定
go run ./cmd/cli/filesystem -operation=read_file -path=/path/to/file
```

指定した `-path` や `-source`/`-destination` は自動的に許可ディレクトリとして扱われます。明示的に追加で許可する設定は不要です。

## フラグ

| フラグ | 説明 | 例 |
| ------ | ---- | --- |
| `-operation` | 実行する操作を指定 | `-operation=read_file` |
| `-path` | 処理対象のパス。読取/書込/一覧などで利用 | `-path=./README.md` |
| `-source` | 移動元パス (`move_file`) | `-source=./old.txt` |
| `-destination` | 移動先パス (`move_file`) | `-destination=./archive/old.txt` |
| `-content` | 書き込む内容 (`write_file`) | `-content="hello"` |
| `-pattern` | 検索パターン (`search_files`, Go正規表現) | `-pattern="(?i)config"` |
| `-offset` | `read_file`/`list_directory`/`directory_tree`での開始位置（1始まり、既定値1） | `-offset=120` |
| `-limit` | `read_file`/`list_directory`/`directory_tree`で返す最大件数（既定値2000） | `-limit=200` |
| `-depth` | `directory_tree`で辿る最大階層（0で無制限、既定値0） | `-depth=2` |

## 利用可能なoperation

| operation | 説明 | 必須フラグ |
| --------- | ---- | ---------- |
| `read_file` | ファイル内容を表示 | `-path` |
| `write_file` | ファイルへ書き込み | `-path`, `-content` (空文字を書き込む場合は `-content=""`) |
| `create_directory` | ディレクトリを作成 | `-path` |
| `list_directory` | 指定フォルダ直下を一覧表示 | `-path` (`-offset`/`-limit`でページング可能) |
| `directory_tree` | YAML形式のツリーを表示 | `-path` |
| `move_file` | ファイル/ディレクトリを移動 | `-source`, `-destination` |
| `search_files` | パターンに一致する項目を検索 | `-path`, `-pattern` |
| `get_file_info` | ファイル/ディレクトリの詳細を表示 | `-path` |
| `list_allowed_directories` | 許可済みディレクトリを表示 | `-path` (省略時はカレントディレクトリを許可) |

## 使用例

```bash
# ファイルの読み取り
go run ./cmd/cli/filesystem -operation=read_file -path=./README.md

# 任意範囲のみ読み取り（offset/limit）
go run ./cmd/cli/filesystem -operation=read_file -path=./README.md -offset=34 -limit=20

# ディレクトリツリーをページング表示
go run ./cmd/cli/filesystem -operation=directory_tree -path=. -offset=1 -limit=10 -depth=2

# ディレクトリ一覧をページング表示
go run ./cmd/cli/filesystem -operation=list_directory -path=. -offset=26 -limit=10

# ファイルの作成
go run ./cmd/cli/filesystem -operation=write_file -path=./notes/todo.txt -content="- [ ] task"

# ディレクトリ検索
go run ./cmd/cli/filesystem -operation=search_files -path=./notes -pattern=todo

# 正規表現でのディレクトリ検索（例: cmd以下のGoファイル）
go run ./cmd/cli/filesystem -operation=search_files -path=. -pattern='^cmd/.+\\.go$'

# ファイルの移動
go run ./cmd/cli/filesystem -operation=move_file -source=./notes/todo.txt -destination=./notes/archive/todo.txt

# 許可ディレクトリの確認（path省略時は作業ディレクトリ）
go run ./cmd/cli/filesystem -operation=list_allowed_directories
```

## 出力仕様

### read_file

`read_file` は内部サービスと同一フォーマットで結果を整形します。

- 各行を `L{行番号}: {内容}` で表示（1始まりの行番号）
- 改行コードはLFへ正規化し、CRLF末尾の `\r` は除去
- 1行あたり500文字を上限にサロゲートを壊さない形でトリム
- 末尾の空行も同じ形式で表示されるため、ファイル終端の空行確認が容易
- `-offset`/`-limit` フラグで開始行と取得行数を制御（既定値は1行目/最大2000行）

```bash
$ go run ./cmd/cli/filesystem -operation=read_file -path=./examples/sample.txt
L1: package main
L2: import "fmt"
L3: func main() {
...
L9: }
```

長文の場合も先頭500文字だけを返すため、MCPツールと同じく過度な出力でタイムアウトするリスクを抑えられます。

### directory_tree

`directory_tree` はYAML形式でツリーを表示します。`-depth`で探索を途中で打ち切ったノードには、残りの直下要素数を示す`truncated_children`フィールドが自動で追加され、0のときはフィールド自体が出力されません。`-limit`を設定し結果が途中で切れる場合は、`list_directory`と同様に末尾へ「More than ...」メッセージを付けて次の`-offset`値を案内します。

```yaml
- name: tmp
  type: directory
  truncated_children: 5   # depth制限の先に5件存在する
```

## 補足

- すべての操作は `internal/filesystem/usecases.FileSystemService` を経由しており、許可ディレクトリ外のパスにはアクセスできません。
- 大量出力になりそうな場合は `-path` を最小限のディレクトリに指定して探索範囲を絞ってください。
- `pattern` は常にGoの正規表現として解釈されます。大文字小文字を区別しない検索が必要な場合は `(?i)` プレフィックスを付けるなど、正規表現の機能を利用してください。
- Goの正規表現エンジン（RE2）は後方参照や先読みなど一部構文をサポートしていません。複雑な条件は文字クラスやアンカーを組み合わせて表現してください。
- エラーは標準エラー出力に表示され、失敗時は終了コード1で終了します。
