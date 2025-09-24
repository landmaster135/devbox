# ISO-8601 コンバーター

UNIXタイムスタンプとISO-8601形式の相互変換を行うコマンドラインツールです。

## 機能

- UNIXタイムスタンプからISO-8601形式への変換
- ISO-8601形式からUNIXタイムスタンプへの変換
- 日付文字列からUNIXタイムスタンプへの変換（UTCまたはJST）

## 使用方法

### UNIXタイムスタンプからISO-8601形式への変換

```bash
go run ./cmd/cli/iso8601-converter --to-iso --input <unix_timestamp>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --to-iso --input 1619712000
# 出力例: 2021-04-30T01:00:00+09:00
```

### ISO-8601形式からUNIXタイムスタンプへの変換

```bash
go run ./cmd/cli/iso8601-converter --to-unix --input <iso8601_time>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --to-unix --input "2021-04-30T00:00:00Z"
# 出力例: 1619740800
```

### 日付からUNIXタイムスタンプへの変換（UTC）

```bash
go run ./cmd/cli/iso8601-converter --to-unix --input <date>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --to-unix --input "2021-04-30"
# 出力例: 1619740800
```

### 日付からUNIXタイムスタンプへの変換（JST）

```bash
go run ./cmd/cli/iso8601-converter --to-unix --is-jst --input <date>
```

例:
```bash
go run ./cmd/cli/iso8601-converter --to-unix --is-jst --input "2021-04-30"
# 出力例: 1619708400
```

## サポートしている日付形式

- ISO-8601形式（例：2021-04-30T00:00:00Z）
- ハイフン区切りの日付（例：2021-04-30）
- スラッシュ区切りの日付（例：2021/04/30）
- 区切りなしの日付（例：20210430）

## オプション

- `--to-iso`: UNIXタイムスタンプからISO-8601形式への変換
- `--to-unix`: ISO-8601形式またはシンプルな日付からUNIXタイムスタンプへの変換
- `--is-jst`: 日付をJSTタイムゾーンとして扱う（デフォルトはUTC）
- `--input`: 変換する値
- `--help`: ヘルプメッセージの表示
