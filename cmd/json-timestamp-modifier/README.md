# JSON Timestamp Modifier ツール使用方法

## 概要

JSON Timestamp Modifierは、JSONファイルに現在の日時のタイムスタンプを追加するためのコマンドラインツールです。

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

JSONファイルに現在の日時のタイムスタンプを追加するには、以下のコマンドを使用します：

```bash
./json-timestamp-modifier -file <JSONファイルパス> [-key <キー>]
```

例：

```bash
./json-timestamp-modifier -file data.json
```

このコマンドは、`data.json`ファイルに`timestamp`キー（デフォルト）と現在の日時のタイムスタンプを追加します。

```bash
./json-timestamp-modifier -file data.json -key created_at
```

このコマンドは、`data.json`ファイルに`created_at`キーと現在の日時のタイムスタンプを追加します。

- ファイルが存在しない場合は、新しいJSONファイルが作成されます。
- 指定されたキーが既に存在する場合は、値が上書きされます。
- タイムスタンプはUNIXタイムスタンプ（1970年1月1日からの経過秒数）として整数で追加されます（例：`1744010395`）。

## オプション

| オプション | 説明 |
|------------|------|
| `-file` | 操作するJSONファイルのパス（必須） |
| `-key` | タイムスタンプを追加するキー（デフォルト: `timestamp`） |

## エラーメッセージ

以下のような場合にエラーメッセージが表示されます：

- ファイルパスが指定されていない場合
- キーが空の場合
- JSONファイルの読み込みや書き込みに失敗した場合

## 使用例

### 新しいJSONファイルの作成

```bash
./json-timestamp-modifier -file new_data.json
```

### 既存のJSONファイルにタイムスタンプを追加

```bash
./json-timestamp-modifier -file data.json
```

### カスタムキーでタイムスタンプを追加

```bash
./json-timestamp-modifier -file data.json -key updated_at
```

## 注意事項

- JSONファイルは常にUTF-8エンコーディングで処理されます。
- タイムスタンプはシステムの現在の日時に基づいて生成されます。
- ファイルパスは相対パスまたは絶対パスで指定できます。
