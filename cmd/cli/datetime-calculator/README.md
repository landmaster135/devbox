# Datetime Calculator

日時計算を行うCLIツールです。指定された基準日時に対して年月日時分秒の加算・減算を行い、結果を表示します。

## 機能

- **日時加算**: 基準日時に指定された期間を加算
- **日時減算**: 基準日時から指定された期間を減算
- **時間単位計算**: 複数の時間単位値を合計し、異なる時間単位に変換
- **時間単位変換**: 単一の時間単位値を別の時間単位に変換
- **時間抽出**: テキストやファイルから「合計[数値]分掛かった。」パターンを抽出して合計時間を計算
- **柔軟な期間指定**: 年、月、日、時、分、秒を個別に指定可能
- **短縮オプション**: 全てのオプションに短縮形を提供
- **直感的な出力**: 計算式と結果を分かりやすく表示

## インストール

```bash
# プロジェクトルートから
go build -o bin/datetime-calculator ./cmd/cli/datetime-calculator
```

## 使用方法

### 基本的な使用方法

```bash
./bin/datetime-calculator -operation add -year 2025 -month 1 -day 15 -hour 12 -minute 30 -second 0 -duration-year 1 -duration-month 2 -duration-day 10 -duration-hour 5 -duration-minute 30 -duration-second 45
```

出力例:
```
2025-1-15 12:30:0 + 1年2月10日5時間30分45秒 = 2026-03-25 18:00:45
```

### オプション

#### 基本オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-operation` | `-o` | 日時操作 (add, subtract, sum, parse-time) | | `-o add` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

#### 日時加算・減算用オプション

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

#### 時間単位計算・変換用オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-figures` | `-f` | カンマ区切りの数値リスト | | `-f 3600,1800,7200` |
| `-input-unit` | `-iu` | 入力時間単位 (year, month, day, hour, minute, second) | | `-iu second` |
| `-output-unit` | `-ou` | 出力時間単位 (year, month, day, hour, minute, second) | | `-ou hour` |

#### 時間抽出用オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-file-path` | `-fp` | ファイルパス (.mdまたは.txt形式) | | `-fp /path/to/file.txt` |
| `-text-input` | `-ti` | テキスト入力 | | `-ti "合計30分掛かった。"` |
| `-output-unit` | `-ou` | 出力時間単位 (year, month, day, hour, minute, second) | minute | `-ou hour` |

## 使用例

### 日時加算

```bash
# 基本的な加算
./bin/datetime-calculator -operation add -year 2025 -month 1 -day 15 -hour 12 -minute 30 -second 0 -duration-year 1 -duration-month 2 -duration-day 10 -duration-hour 5 -duration-minute 30 -duration-second 45

# 短縮形を使用
./bin/datetime-calculator -o add -y 2025 -m 1 -d 15 -hr 12 -min 30 -s 0 -dy 1 -dm 2 -dd 10 -dh 5 -dmin 30 -ds 45

# 一部の期間のみ指定
./bin/datetime-calculator -o add -y 2025 -m 6 -d 15 -dy 2 -dm 3
```

### 日時減算

```bash
# 基本的な減算
./bin/datetime-calculator -operation subtract -year 2025 -month 3 -day 25 -hour 18 -minute 0 -second 45 -duration-year 0 -duration-month 1 -duration-day 5 -duration-hour 2 -duration-minute 15 -duration-second 30

# 短縮形を使用
./bin/datetime-calculator -o subtract -y 2025 -m 3 -d 25 -hr 18 -min 0 -s 45 -dy 0 -dm 1 -dd 5 -dh 2 -dmin 15 -ds 30

# 日数のみ減算
./bin/datetime-calculator -o subtract -y 2025 -m 12 -d 31 -dd 10
```

### 時間単位計算（合計）

```bash
# 複数の秒数を合計して時間単位で表示
./bin/datetime-calculator -operation sum -input-unit second -output-unit hour -figures 3600,1800,7200

# 複数の日数を合計して月単位で表示
./bin/datetime-calculator -operation sum -input-unit day -output-unit month -figures 30,15,45

# 短縮形を使用
./bin/datetime-calculator -o sum -iu second -ou hour -f 3600,1800,7200
```

### 時間単位変換

```bash
# 時間を分に変換
./bin/datetime-calculator -operation sum -input-unit hour -output-unit minute -figures 2.5

# 日数を週に変換
./bin/datetime-calculator -operation sum -input-unit day -output-unit week -figures 14

# 短縮形を使用
./bin/datetime-calculator -o sum -iu hour -ou minute -f 2.5
```

### 時間抽出

```bash
# ファイルから時間を抽出（デフォルト：分単位）
./bin/datetime-calculator -operation parse-time -file-path /path/to/file.txt

# テキストから時間を抽出（デフォルト：分単位）
./bin/datetime-calculator -operation parse-time -text-input "作業は合計30分掛かった。別の作業は合計45分掛かった。"

# 時間単位で出力
./bin/datetime-calculator -operation parse-time -text-input "合計120分掛かった。" -output-unit hour

# 秒単位で出力
./bin/datetime-calculator -operation parse-time -file-path /path/to/file.txt -output-unit second

# 短縮形を使用（ファイルから）
./bin/datetime-calculator -o parse-time -fp /path/to/work_log.md

# 短縮形を使用（テキストから、時間単位で出力）
./bin/datetime-calculator -o parse-time -ti "合計120分掛かった。" -ou hour
```

## 出力フォーマット

### 加算の出力

```
2025-1-15 12:30:0 + 1年2月10日5時間30分45秒 = 2026-03-25 18:00:45
2025-6-15 0:0:0 + 2年3月 = 2027-09-15 00:00:00
```

### 減算の出力

```
2025-3-25 18:0:45 - 1月5日2時間15分30秒 = 2025-02-20 15:45:15
2025-12-31 0:0:0 - 10日 = 2025-12-21 00:00:00
```

### 時間単位計算の出力

```
sum([3600 1800 7200] second) = 3.500000 hour
sum([30 15 45] day) = 3.000000 month
sum([2.5] hour) = 150.000000 minute
```

### 時間抽出の出力

```
# デフォルト（分単位）
抽出された時間の合計: 75.000000分
抽出された時間の合計: 120.000000分

# 時間単位で出力
抽出された時間の合計: 1.250000時間
抽出された時間の合計: 2.000000時間

# 秒単位で出力
抽出された時間の合計: 4500.000000秒
抽出された時間の合計: 7200.000000秒
```

## エラーハンドリング

### 無効な操作タイプ

```bash
./bin/datetime-calculator -operation invalid
```

```
エラー: 無効な操作タイプです: invalid
```

### 必須パラメータの不足

```bash
./bin/datetime-calculator
```

```
エラー: 操作タイプが指定されていません
日時計算CLIツール

使用方法:
  日時加算:
    ./bin/datetime-calculator -operation add -year 2025 -month 1 -day 15 -hour 12 -minute 30 -second 0 -duration-year 1 -duration-month 2 -duration-day 10 -duration-hour 5 -duration-minute 30 -duration-second 45
...
```

### 無効な数値

```bash
./bin/datetime-calculator -o add -y invalid -m 1
```

```
エラー: 無効な年の値です: invalid
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **TDD**: テスト駆動開発によるテストファースト実装

### ディレクトリ構造

```
internal/datetime_calculator/
├── config/          # 設定管理
│   ├── config.go    # 設定構造体とパーサー
│   ├── flag_parser.go # フラグ解析
│   └── interfaces.go  # インターフェース定義
└── usecases/        # ビジネスロジック
    ├── calculator.go # 日時計算ロジック
    └── services.go  # サービス層
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `time`, `flag`, `fmt`, `strconv`など
- **テストフレームワーク**: Go標準のtestingパッケージ

### 日時計算の仕組み

- **基準日時**: `time.Date()`を使用してtime.Time型を生成
- **年月日の加算/減算**: `time.Time.AddDate()`メソッドを使用
- **時分秒の加算/減算**: `time.Duration`を使用して`time.Time.Add()`メソッドで処理
- **出力フォーマット**: `"2006-01-02 15:04:05"`形式で統一

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/datetime-calculator ./cmd/cli/datetime-calculator

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/datetime-calculator-linux ./cmd/cli/datetime-calculator
GOOS=windows GOARCH=amd64 go build -o bin/datetime-calculator.exe ./cmd/cli/datetime-calculator
GOOS=darwin GOARCH=amd64 go build -o bin/datetime-calculator-mac ./cmd/cli/datetime-calculator
```

### テスト

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

## 実用的な使用例

### プロジェクト管理

```bash
# プロジェクト開始日から3ヶ月後の締切日を計算
./bin/datetime-calculator -o add -y 2025 -m 1 -d 15 -dm 3

# 締切日から2週間前のリマインダー日を計算
./bin/datetime-calculator -o subtract -y 2025 -m 4 -d 15 -dd 14
```

### イベント計画

```bash
# イベント開催日の1ヶ月前の準備開始日を計算
./bin/datetime-calculator -o subtract -y 2025 -m 6 -d 20 -dm 1

# 定期イベントの次回開催日を計算（3ヶ月間隔）
./bin/datetime-calculator -o add -y 2025 -m 3 -d 15 -dm 3
```

### 勤務時間計算

```bash
# 勤務開始時刻に8時間を加算して終業時刻を計算
./bin/datetime-calculator -o add -y 2025 -m 1 -d 20 -hr 9 -min 0 -dh 8

# 残業時間を加算
./bin/datetime-calculator -o add -y 2025 -m 1 -d 20 -hr 17 -min 0 -dh 2 -dmin 30
```

### 作業時間の集計

```bash
# 作業ログファイルから時間を抽出して合計（デフォルト：分単位）
./bin/datetime-calculator -o parse-time -fp work_log.md

# 作業ログファイルから時間を抽出して時間単位で表示
./bin/datetime-calculator -o parse-time -fp work_log.md -ou hour

# 複数の作業報告から時間を抽出（デフォルト：分単位）
./bin/datetime-calculator -o parse-time -ti "タスクAは合計60分掛かった。タスクBは合計90分掛かった。タスクCは合計45分掛かった。"

# 複数の作業報告から時間を抽出して時間単位で表示
./bin/datetime-calculator -o parse-time -ti "タスクAは合計60分掛かった。タスクBは合計90分掛かった。タスクCは合計45分掛かった。" -ou hour
```

### 時間抽出の対応パターン

時間抽出機能は以下のパターンに対応しています：

- **基本パターン**: `合計[数値]分掛かった。`
- **例**: 
  - `合計30分掛かった。`
  - `合計120分掛かった。`
  - `作業は合計45分掛かった。`

**対応ファイル形式**:
- `.md` (Markdown)
- `.txt` (テキスト)

**使用例のファイル内容**:
```
# 作業ログ

## タスクA
詳細な作業内容...
合計60分掛かった。

## タスクB  
別の作業内容...
合計90分掛かった。

## タスクC
さらに別の作業...
合計45分掛かった。
```

上記のファイルを処理すると、`抽出された時間の合計: 195分` が出力されます。

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
