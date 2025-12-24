# filesystem

ファイル読み書きや検索、移動などの基本操作を安全に行うCLIツールです。MCP版と同じ `internal/filesystem/usecases` を利用し、許可したディレクトリの範囲内で各操作を実行します。

## 主な機能

- **read_file**: テキストファイルを読み取り内容を標準出力に表示
- **write_file**: ファイルの作成または上書き
- **create_directory**: ネストしたディレクトリの作成
- **list_directory / directory_tree**: ディレクトリのフラット一覧またはツリー構造を表示
- **move_file**: ファイル／ディレクトリの移動やリネーム
- **search_files**: 名前に特定のパターンを含む項目の再帰検索（除外パターン対応）
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

必要に応じて `-allowed-dirs=/workspace,/tmp` のようにカンマ区切りで許可ディレクトリを追加できます。`move_file` のように複数パスを扱う操作では、指定した `source`/`destination` も自動的に許可対象へ追加されます。

## フラグ

| フラグ | 説明 | 例 |
| ------ | ---- | --- |
| `-operation` | 実行する操作を指定 | `-operation=read_file` |
| `-path` | 処理対象のパス。読取/書込/一覧などで利用 | `-path=./README.md` |
| `-source` | 移動元パス (`move_file`) | `-source=./old.txt` |
| `-destination` | 移動先パス (`move_file`) | `-destination=./archive/old.txt` |
| `-content` | 書き込む内容 (`write_file`) | `-content="hello"` |
| `-pattern` | 検索パターン (`search_files`) | `-pattern=config` |
| `-exclude-pattern` | 除外パターン (`search_files`) | `-exclude-pattern="*.git"` |
| `-allowed-dirs` | カンマ区切りの追加許可ディレクトリ | `-allowed-dirs=$HOME/projects` |

## 利用可能なoperation

| operation | 説明 | 必須フラグ |
| --------- | ---- | ---------- |
| `read_file` | ファイル内容を表示 | `-path` |
| `write_file` | ファイルへ書き込み | `-path`, `-content` (空文字を書き込む場合は `-content=""`) |
| `create_directory` | ディレクトリを作成 | `-path` |
| `list_directory` | 指定フォルダ直下を一覧表示 | `-path` |
| `directory_tree` | JSON形式のツリーを表示 | `-path` |
| `move_file` | ファイル/ディレクトリを移動 | `-source`, `-destination` |
| `search_files` | パターンに一致する項目を検索 | `-path`, `-pattern` |
| `get_file_info` | ファイル/ディレクトリの詳細を表示 | `-path` |
| `list_allowed_directories` | 許可済みディレクトリを表示 | `-path` または `-allowed-dirs` |

## 使用例

```bash
# ファイルの読み取り
go run ./cmd/cli/filesystem -operation=read_file -path=./README.md

# ファイルの作成
go run ./cmd/cli/filesystem -operation=write_file -path=./notes/todo.txt -content="- [ ] task"

# ディレクトリ検索（tmp配下を許可）
go run ./cmd/cli/filesystem -operation=search_files -path=/tmp -pattern=log -allowed-dirs=/tmp

# ファイルの移動
go run ./cmd/cli/filesystem -operation=move_file -source=./notes/todo.txt -destination=./notes/archive/todo.txt

# 許可ディレクトリの確認
go run ./cmd/cli/filesystem -operation=list_allowed_directories -allowed-dirs=/workspace/projects
```

## 補足

- すべての操作は `internal/filesystem/usecases.FileSystemService` を経由しており、許可ディレクトリ外のパスにはアクセスできません。
- 大量出力になる `directory_tree` や `search_files` では `allowed-dirs` を絞って利用すると意図しない探索を防げます。
- エラーは標準エラー出力に表示され、失敗時は終了コード1で終了します。
