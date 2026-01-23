# Datetime Calculator

日時計算を行うCLIツールです。指定された基準日時に対して年月日時分秒の加算・減算を行い、結果を表示します。

## 機能

- **日時加算**: 基準日時に指定された期間を加算
- **日時減算**: 基準日時から指定された期間を減算
- **時間単位計算**: 複数の時間単位値を合計し、異なる時間単位に変換
- **時間単位変換**: 単一の時間単位値を別の時間単位に変換
- **時間抽出**: テキストやファイルから「合計[数値]分掛かった。」パターンを抽出して合計時間を計算
- **日次見出し生成**: 実行日からのオフセットを指定して日次タスク用のMarkdown見出しを生成
- **柔軟な期間指定**: 年、月、日、時、分、秒を個別に指定可能
- **短縮オプション**: 全てのオプションに短縮形を提供
- **直感的な出力**: 計算式と結果を分かりやすく表示

## インストール

```bash
# プロジェクトルートから
go build -o bin/datetime-calculator ./cmd/cli/datetime-calculator
```

## 使用方法

### 基本オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-operation` | `-o` | 日時操作 (add, subtract, sum, parse-time, generate-daily-heading) | | `-o add` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

### 日時加算・減算用オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-year` | `-y` | 基準日時の年 | 2025 | `-y 2025` |
| `-month` | `-m` | 基準日時の月 | 1 | `-m 1` |
| `-day` | `-d` | 基準日時の日 | 1 | `-d 15` |
| `-hour` | `-hr` | 基準日時の時 | 0 | `-hr 12` |
| `-minute` | `-min` | 基準日時の分 | 0 | `-min 30` |
| `-second` | `-s` | 基準日時の秒 | 0 | `-s 0` |
| `-duration-year` | `-dy` | 加算/減算する年 | 0 | `-dy 1` |
| `-duration-month` | `-dm` | 加算/減算する月 | 0 | `-dm 2` |
| `-duration-day` | `-dd` | 加算/減算する日 | 0 | `-dd 10` |
| `-duration-hour` | `-dh` | 加算/減算する時 | 0 | `-dh 5` |
| `-duration-minute` | `-dmin` | 加算/減算する分 | 0 | `-dmin 30` |
| `-duration-second` | `-ds` | 加算/減算する秒 | 0 | `-ds 45` |

### 時間単位計算・変換用オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-figures` | `-f` | カンマ区切りの数値リスト | | `-f 3600,1800,7200` |
| `-input-unit` | `-iu` | 入力時間単位 (year, month, day, hour, minute, second) | | `-iu second` |
| `-output-unit` | `-ou` | 出力時間単位 (year, month, day, hour, minute, second) | | `-ou hour` |

### 時間抽出用オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-file-path` | `-fp` | ファイルパス (.mdまたは.txt形式) | | `-fp /path/to/file.txt` |
| `-text-input` | `-ti` | テキスト入力 | | `-ti "合計30分掛かった。"` |
| `-output-unit` | `-ou` | 出力時間単位 (year, month, day, hour, minute, second) | minute | `-ou hour` |

### 日次見出し生成用オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-day-offset` | `-do` | 日付オフセット（実行日からの日数） | 0 | `-do -1` |

## 使用例

### 日時加算

```bash
# 基本的な加算
go run ./cmd/cli/datetime-calculator -operation add -year 2025 -month 1 -day 15 -hour 12 -minute 30 -second 0 -duration-year 1 -duration-month 2 -duration-day 10 -duration-hour 5 -duration-minute 30 -duration-second 45

# 短縮形を使用
go run ./cmd/cli/datetime-calculator -o add -y 2025 -m 1 -d 15 -hr 12 -min 30 -s 0 -dy 1 -dm 2 -dd 10 -dh 5 -dmin 30 -ds 45

# 一部の期間のみ指定
go run ./cmd/cli/datetime-calculator -o add -y 2025 -m 6 -d 15 -dy 2 -dm 3
```

### 日時減算

```bash
# 基本的な減算
go run ./cmd/cli/datetime-calculator -operation subtract -year 2025 -month 3 -day 25 -hour 18 -minute 0 -second 45 -duration-year 0 -duration-month 1 -duration-day 5 -duration-hour 2 -duration-minute 15 -duration-second 30

# 短縮形を使用
go run ./cmd/cli/datetime-calculator -o subtract -y 2025 -m 3 -d 25 -hr 18 -min 0 -s 45 -dy 0 -dm 1 -dd 5 -dh 2 -dmin 15 -ds 30

# 日数のみ減算
go run ./cmd/cli/datetime-calculator -o subtract -y 2025 -m 12 -d 31 -dd 10
```

### 時間単位計算（合計）

```bash
# 複数の秒数を合計して時間単位で表示
go run ./cmd/cli/datetime-calculator -operation sum -input-unit second -output-unit hour -figures 3600,1800,7200

# 複数の日数を合計して月単位で表示
go run ./cmd/cli/datetime-calculator -operation sum -input-unit day -output-unit month -figures 30,15,45

# 短縮形を使用
go run ./cmd/cli/datetime-calculator -o sum -iu second -ou hour -f 3600,1800,7200
```

### 時間単位変換

```bash
# 時間を分に変換
go run ./cmd/cli/datetime-calculator -operation sum -input-unit hour -output-unit minute -figures 2.5

# 日数を週に変換
go run ./cmd/cli/datetime-calculator -operation sum -input-unit day -output-unit week -figures 14

# 短縮形を使用
go run ./cmd/cli/datetime-calculator -o sum -iu hour -ou minute -f 2.5
```

### 時間抽出

```bash
# ファイルから時間を抽出（デフォルト：分単位）
go run ./cmd/cli/datetime-calculator -operation parse-time -file-path /path/to/file.txt

# テキストから時間を抽出（デフォルト：分単位）
go run ./cmd/cli/datetime-calculator -operation parse-time -text-input "作業は合計30分掛かった。別の作業は合計45分掛かった。"

# 時間単位で出力
go run ./cmd/cli/datetime-calculator -operation parse-time -text-input "合計120分掛かった。" -output-unit hour

# 秒単位で出力
go run ./cmd/cli/datetime-calculator -operation parse-time -file-path /path/to/file.txt -output-unit second

# 短縮形を使用（ファイルから）
go run ./cmd/cli/datetime-calculator -o parse-time -fp /path/to/work_log.md

# 短縮形を使用（テキストから、時間単位で出力）
go run ./cmd/cli/datetime-calculator -o parse-time -ti "合計120分掛かった。" -ou hour
```

### 日次見出し生成

```bash
# 昨日から今日にかけての見出しを生成
go run ./cmd/cli/datetime-calculator -operation generate-daily-heading -day-offset -1

# 今日から明日にかけての見出しを生成（デフォルト）
go run ./cmd/cli/datetime-calculator -operation generate-daily-heading -day-offset 0

# 3日後から4日後にかけての見出しを生成
go run ./cmd/cli/datetime-calculator -operation generate-daily-heading -day-offset 3

# 短縮形を使用
go run ./cmd/cli/datetime-calculator -o generate-daily-heading -do -1
```

## テスト

```bash
# 単体テスト
go test ./internal/datetime_calculator/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/datetime_calculator/...
go tool cover -html=coverage.out -o coverage.html
```

### 関連ツール

このCLIツールと同じサービスを使用するMCPツールも利用可能です：
- `devbox/cmd/mcp/datetime_calculator/mcp.go`

MCPツールを使用することで、Model Context Protocol経由で同じ日時計算機能にアクセスできます。

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
