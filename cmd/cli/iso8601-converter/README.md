# ISO-8601 コンバーター

UNIXタイムスタンプとISO-8601形式の相互変換を行うコマンドラインツールです。

## 機能

- UNIXタイムスタンプからISO-8601形式への変換
- ISO-8601形式からUNIXタイムスタンプへの変換
- 日付文字列からUNIXタイムスタンプへの変換（UTCまたはJST）

## 使用方法

### 現在日時を表示

```bash
go run ./cmd/cli/iso8601-converter --operation now
```

ISO-8601（UTC/JST）とUNIXタイムスタンプをまとめて表示します。

`--format`で出力内容を絞り込めます。

```bash
go run ./cmd/cli/iso8601-converter --operation now --format iso
# ISO-8601形式のみ表示

go run ./cmd/cli/iso8601-converter --operation now --format unix
# UNIXタイムスタンプのみ表示
```

### UNIXタイムスタンプからISO-8601形式への変換

```bash
go run ./cmd/cli/iso8601-converter --operation to-iso --input <unix_timestamp>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --operation to-iso --input 1619712000
# 出力例: 2021-04-30T01:00:00+09:00
```

### ISO-8601形式からUNIXタイムスタンプへの変換

```bash
go run ./cmd/cli/iso8601-converter --operation to-unix --input <iso8601_time>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --operation to-unix --input "2021-04-30T00:00:00Z"
# 出力例: 1619740800
```

### 日付からUNIXタイムスタンプへの変換（UTC）

```bash
go run ./cmd/cli/iso8601-converter --operation to-unix --input <date>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --operation to-unix --input "2021-04-30"
# 出力例: 1619740800
```

### 日付からUNIXタイムスタンプへの変換（JST）

```bash
go run ./cmd/cli/iso8601-converter --operation to-unix --is-jst --input <date>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --operation to-unix --is-jst --input "2021-04-30"
# 出力例: 1619708400
```

## サポートしている日付形式

- ISO-8601形式（例：2021-04-30T00:00:00Z）
- ハイフン区切りの日付（例：2021-04-30）
- スラッシュ区切りの日付（例：2021/04/30）
- 区切りなしの日付（例：20210430）

## オプション

- `--operation`: 実行する操作。`to-iso`、`to-unix`、`now` のいずれかを指定
- `--format`: `--operation now`の時に有効。`all` (デフォルト)、`iso`、`unix` を指定
- `--is-jst`: 日付をJSTタイムゾーンとして扱う（デフォルトはUTC）
- `--input`: 変換する値
- `--help`: ヘルプメッセージの表示
