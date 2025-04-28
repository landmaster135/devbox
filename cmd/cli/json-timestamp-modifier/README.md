# JSON Timestamp Modifier ツール使用方法

## 概要

JSON Timestamp Modifierは、JSONファイルのタイムスタンプ関連の操作を行うためのコマンドラインツールです。以下の機能を提供します：

1. JSONファイルに現在の日時のタイムスタンプを追加
2. JSONファイル内の指定したキーの値をISO-8601形式からUNIXタイムスタンプに変換
3. JSONファイル内の指定したキーの値をUNIXタイムスタンプからISO-8601形式に変換

これらの操作は、単一のJSONファイルだけでなく、ディレクトリ内の全てのJSONファイルに対しても実行できます。

## インストール

ビルドスクリプトを使用してツールをビルドします：

```bash
cd devbox
./scripts/build_json_timestamp_modifier.sh
```

ビルドが成功すると、以下の場所に実行ファイルが生成されます：

- Linux用: `pkg/bin/linux_amd64/json-timestamp_modifier`
- Windows用: `pkg/bin/win_amd64/json-timestamp_modifier.exe`

## 基本的な使い方

### 単一ファイルの操作

#### タイムスタンプの追加

JSONファイルに現在の日時のタイムスタンプを追加するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -file <JSONファイルパス> [-key <キー>] [-mode add]
```

例：

```bash
./json-timestamp-modifier -file data.json
```

このコマンドは、`data.json`ファイルに`timestamp`キー（デフォルト）と現在の日時のタイムスタンプを追加します。

```bash
./json-timestamp-modifier -file data.json -key created_at -mode add
```

このコマンドは、`data.json`ファイルに`created_at`キーと現在の日時のタイムスタンプを追加します。

- ファイルが存在しない場合は、新しいJSONファイルが作成されます。
- 指定されたキーが既に存在する場合は、値が上書きされます。
- タイムスタンプはUNIXタイムスタンプ（1970年1月1日からの経過秒数）として整数で追加されます（例：`1744010395`）。

#### ISO-8601形式からUNIXタイムスタンプへの変換

JSONファイル内の指定したキーの値をISO-8601形式からUNIXタイムスタンプに変換するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -file <JSONファイルパス> -key <キー> -mode to-unix [-is-jst]
```

例：

```bash
./json-timestamp-modifier -file data.json -key date -mode to-unix
```

このコマンドは、`data.json`ファイル内の`date`キーの値をISO-8601形式からUNIXタイムスタンプに変換します。

```bash
./json-timestamp-modifier -file data.json -key date -mode to-unix -is-jst
```

このコマンドは、`data.json`ファイル内の`date`キーの値を日本標準時（JST）として扱い、UNIXタイムスタンプに変換します。

- `-is-jst`オプションは、時刻情報が含まれていない日付文字列（例：`2025-04-28`）をJSTタイムゾーンとして扱う場合に使用します。
- 時刻情報が含まれている場合（例：`2025-04-28T09:00:00Z`）は、指定されたタイムゾーン情報に基づいて変換されます。

#### UNIXタイムスタンプからISO-8601形式への変換

JSONファイル内の指定したキーの値をUNIXタイムスタンプからISO-8601形式に変換するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -file <JSONファイルパス> -key <キー> -mode to-iso
```

例：

```bash
./json-timestamp-modifier -file data.json -key timestamp -mode to-iso
```

このコマンドは、`data.json`ファイル内の`timestamp`キーの値をUNIXタイムスタンプからISO-8601形式に変換します。

- 変換結果はISO-8601形式の文字列として保存されます（例：`2025-04-28T09:00:00+09:00`）。

### ディレクトリ内の全てのJSONファイルの操作

#### ディレクトリ内の全てのJSONファイルにタイムスタンプを追加

ディレクトリ内の全てのJSONファイルに現在の日時のタイムスタンプを追加するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -dir <ディレクトリパス> [-key <キー>] [-mode add] [-recursive]
```

例：

```bash
./json-timestamp-modifier -dir ./data
```

このコマンドは、`./data`ディレクトリ内の全てのJSONファイルに`timestamp`キー（デフォルト）と現在の日時のタイムスタンプを追加します。

```bash
./json-timestamp-modifier -dir ./data -key created_at -recursive
```

このコマンドは、`./data`ディレクトリとそのサブディレクトリ内の全てのJSONファイルに`created_at`キーと現在の日時のタイムスタンプを追加します。

- `-recursive`オプションを指定すると、サブディレクトリ内のJSONファイルも処理されます。
- 処理されたファイル数が出力されます。

#### ディレクトリ内の全てのJSONファイルのISO-8601形式をUNIXタイムスタンプに変換

ディレクトリ内の全てのJSONファイルの指定したキーの値をISO-8601形式からUNIXタイムスタンプに変換するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -dir <ディレクトリパス> -key <キー> -mode to-unix [-is-jst] [-recursive]
```

例：

```bash
./json-timestamp-modifier -dir ./data -key date -mode to-unix
```

このコマンドは、`./data`ディレクトリ内の全てのJSONファイルの`date`キーの値をISO-8601形式からUNIXタイムスタンプに変換します。

```bash
./json-timestamp-modifier -dir ./data -key date -mode to-unix -is-jst -recursive
```

このコマンドは、`./data`ディレクトリとそのサブディレクトリ内の全てのJSONファイルの`date`キーの値を日本標準時（JST）として扱い、UNIXタイムスタンプに変換します。

#### ディレクトリ内の全てのJSONファイルのUNIXタイムスタンプをISO-8601形式に変換

ディレクトリ内の全てのJSONファイルの指定したキーの値をUNIXタイムスタンプからISO-8601形式に変換するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -dir <ディレクトリパス> -key <キー> -mode to-iso [-recursive]
```

例：

```bash
./json-timestamp-modifier -dir ./data -key timestamp -mode to-iso
```

このコマンドは、`./data`ディレクトリ内の全てのJSONファイルの`timestamp`キーの値をUNIXタイムスタンプからISO-8601形式に変換します。

```bash
./json-timestamp-modifier -dir ./data -key timestamp -mode to-iso -recursive
```

このコマンドは、`./data`ディレクトリとそのサブディレクトリ内の全てのJSONファイルの`timestamp`キーの値をUNIXタイムスタンプからISO-8601形式に変換します。

## オプション

| オプション | 説明 |
|------------|------|
| `-file` | 操作するJSONファイルのパス |
| `-dir` | 操作するJSONファイルが含まれるディレクトリのパス |
| `-recursive` | ディレクトリを再帰的に処理する（`-dir`オプションと共に使用）（デフォルト: `false`） |
| `-key` | 操作するキー（デフォルト: `timestamp`） |
| `-mode` | 操作モード: `add`（タイムスタンプ追加）, `to-unix`（ISO-8601→UNIX）, `to-iso`（UNIX→ISO-8601）（デフォルト: `add`） |
| `-is-jst` | 日付のみの文字列をJSTとして扱う（`to-unix`モードのみ有効）（デフォルト: `false`） |

**注意**: `-file`と`-dir`オプションは同時に指定できません。どちらか一方を指定してください。

## エラーメッセージ

以下のような場合にエラーメッセージが表示されます：

- ファイルパスが指定されていない場合
- キーが空の場合
- JSONファイルの読み込みや書き込みに失敗した場合

## 使用例

### 単一ファイルの操作

```bash
# 新しいJSONファイルの作成
./json-timestamp-modifier -file new_data.json

# 既存のJSONファイルにタイムスタンプを追加
./json-timestamp-modifier -file data.json

# カスタムキーでタイムスタンプを追加
./json-timestamp-modifier -file data.json -key updated_at -mode add

# ISO-8601形式からUNIXタイムスタンプに変換
./json-timestamp-modifier -file data.json -key created_at -mode to-unix

# 日付のみの文字列をJSTとして扱い、UNIXタイムスタンプに変換
./json-timestamp-modifier -file data.json -key date -mode to-unix -is-jst

# UNIXタイムスタンプからISO-8601形式に変換
./json-timestamp-modifier -file data.json -key timestamp -mode to-iso
```

### ディレクトリ内の全てのJSONファイルの操作

```bash
# ディレクトリ内の全てのJSONファイルにタイムスタンプを追加
./json-timestamp-modifier -dir ./data

# ディレクトリとサブディレクトリ内の全てのJSONファイルにタイムスタンプを追加
./json-timestamp-modifier -dir ./data -recursive

# ディレクトリ内の全てのJSONファイルのISO-8601形式をUNIXタイムスタンプに変換
./json-timestamp-modifier -dir ./data -key date -mode to-unix

# ディレクトリとサブディレクトリ内の全てのJSONファイルの日付をJSTとして扱い、UNIXタイムスタンプに変換
./json-timestamp-modifier -dir ./data -key date -mode to-unix -is-jst -recursive

# ディレクトリ内の全てのJSONファイルのUNIXタイムスタンプをISO-8601形式に変換
./json-timestamp-modifier -dir ./data -key timestamp -mode to-iso -recursive
```

## 注意事項

- JSONファイルは常にUTF-8エンコーディングで処理されます。
- タイムスタンプはシステムの現在の日時に基づいて生成されます。
- ファイルパスとディレクトリパスは相対パスまたは絶対パスで指定できます。
- ISO-8601形式からUNIXタイムスタンプへの変換では、時刻情報が含まれていない場合は、デフォルトでUTCとして扱われます。JSTとして扱う場合は `-is-jst` オプションを使用してください。
- UNIXタイムスタンプからISO-8601形式への変換では、結果はローカルタイムゾーンで表示されます。
- ディレクトリモードでは、指定したキーが存在しないJSONファイルはスキップされます。
- 処理中にエラーが発生した場合でも、可能な限り処理を継続し、最終的に処理されたファイル数が表示されます。
