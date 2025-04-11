# JSON ISO8601 Converter

## 概要

JSON ISO8601 Converterは、JSONファイル内のISO8601形式（RFC3339）の日時文字列をUNIXタイムスタンプに変換するためのコマンドラインツールです。ディレクトリ内の全てのJSONファイルを再帰的に処理することができます。

## インストール

ビルドスクリプトを使用してツールをビルドします：

```bash
cd devbox
./scripts/build_json_iso8601_converter.sh
```

ビルドが成功すると、以下の場所に実行ファイルが生成されます：

- Linux用: `bin/linux_amd64/json-iso8601-converter`
- Windows用: `bin/win_amd64/json-iso8601-converter.exe`

## 基本的な使い方

JSONファイル内のISO8601形式の日時文字列をUNIXタイムスタンプに変換するには、以下のコマンドを使用します：

```bash
./json-iso8601-converter -dir <ディレクトリパス> -key <キー名> [-recursive] [-dry-run]
```

例：

```bash
./json-iso8601-converter -dir ./data -key created_at
```

このコマンドは、`./data`ディレクトリ内のJSONファイルの`created_at`キーの値をISO8601形式からUNIXタイムスタンプに変換します。

```bash
./json-iso8601-converter -dir ./data -key timestamp -recursive
```

このコマンドは、`./data`ディレクトリとそのサブディレクトリ内のJSONファイルの`timestamp`キーの値をISO8601形式からUNIXタイムスタンプに変換します。

## オプション

| オプション | 説明 |
|------------|------|
| `-dir` | JSONファイルを検索するディレクトリパス（デフォルト: カレントディレクトリ） |
| `-key` | 変換対象のキー名（必須） |
| `-recursive` | サブディレクトリも検索するかどうか（デフォルト: false） |
| `-dry-run` | 変換をシミュレーションするだけで実際には変更しない（デフォルト: false） |

## 変換例

### 変換前のJSONファイル

```json
{
  "id": 123,
  "title": "サンプルデータ",
  "created_at": "2025-04-10T15:30:45Z"
}
```

### 変換後のJSONファイル

```json
{
  "id": 123,
  "title": "サンプルデータ",
  "created_at": 1744010395
}
```

## 注意事項

- JSONファイルは常にUTF-8エンコーディングで処理されます。
- 変換対象のキーの値がISO8601形式（RFC3339）の日時文字列でない場合は、変換されません。
- ドライランモード（`-dry-run`）を使用すると、実際のファイルを変更せずに変換対象のファイル数を確認できます。
- 再帰モード（`-recursive`）を使用すると、指定したディレクトリとそのすべてのサブディレクトリ内のJSONファイルが処理されます。
- 指定したキーの値が**全て**変換されます。
