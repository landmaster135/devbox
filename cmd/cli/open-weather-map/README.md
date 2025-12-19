# OpenWeather API CLIツール

OpenWeather APIを使用して天気予報を取得するコマンドラインツールです。

## 機能

- 指定した都市の天気予報を取得
- 1-5日間の予報に対応
- 日本語での天気情報表示
- 3時間毎の詳細予報も表示

## 使用方法

### 基本的な使用方法

```bash
# フラグを使用した指定
go run ./cmd/cli/open-weather-map -api-key YOUR_API_KEY -city "Tokyo,JP" -max-days 3
go run ./cmd/cli/open-weather-map -k YOUR_API_KEY -c "Tokyo,JP" -d 5

# 位置引数での指定
go run ./cmd/cli/open-weather-map YOUR_API_KEY "Tokyo,JP" 3
```

### オプション

| オプション | 短縮形 | 説明 | 必須 | デフォルト |
|-----------|--------|------|------|-----------|
| `-api-key` | `-k` | OpenWeather APIキー | * | - |
| `-city` | `-c` | 都市名（例: Tokyo,JP, London,UK） | * | - |
| `-max-days` | `-d` | 取得する最大日数（1-5） | * | 3 |
| `-help` | `-h` | ヘルプを表示 | - | - |

## 使用例

```bash
# 東京の3日間予報
go run ./cmd/cli/open-weather-map -k abc123 -c "Tokyo,JP" -d 3

# ロンドンの5日間予報
go run ./cmd/cli/open-weather-map -api-key abc123 -city "London,UK" -max-days 5

# ニューヨークの2日間予報（位置引数）
go run ./cmd/cli/open-weather-map abc123 "New York,US" 2
```

## 出力例

```
=== Tokyo,JP の3日間天気予報 ===

📅 2024年01月01日 (Mon)
🌡️  気温: 5.2°C ～ 12.8°C
☁️  天気: Clear (晴れ)
💨 湿度: 45% | 気圧: 1013 hPa | 風速: 2.1 m/s
⏰ 3時間毎の詳細:
   00:00: 8.5°C 晴れ
   03:00: 6.2°C 晴れ
   06:00: 5.2°C 晴れ
   09:00: 9.1°C 晴れ
   12:00: 12.8°C 晴れ
   15:00: 11.5°C 晴れ
   18:00: 9.3°C 晴れ
   21:00: 7.8°C 晴れ
--------------------------------------------------
```

## セットアップ

### 1. OpenWeather APIキーの取得

1. [OpenWeatherMap](https://openweathermap.org/)にアカウントを作成
2. APIキーを取得（無料プランで十分）

### 2. ビルド

```bash
# プロジェクトルートから
go build -o open-weather-map ./cmd/cli/open-weather-map
```

### 3. 実行

```bash
go run ./cmd/cli/open-weather-map -api-key YOUR_API_KEY -city "Tokyo,JP" -max-days 3
```

## 注意事項

- APIキーはOpenWeatherMapから取得してください
- 都市名は "都市名,国コード" の形式で指定してください（例: Tokyo,JP, London,UK）
- 無料プランでは最大5日間の予報が取得可能です
- APIレート制限にご注意ください（無料プランは1分間に60回まで）

## エラーハンドリング

以下のエラーが発生する可能性があります：

- **APIキー未指定**: `APIキーが指定されていません`
- **都市名未指定**: `都市名が指定されていません`
- **日数範囲外**: `最大日数は1-5の範囲で指定してください`
- **無効なAPIキー**: `APIエラー: ステータスコード 401`
- **都市名が見つからない**: `APIエラー: ステータスコード 404`
- **ネットワークエラー**: `APIリクエストエラー`

## 技術仕様

- **言語**: Go
- **API**: OpenWeather API 2.5
- **対応OS**: Linux, macOS, Windows
- **依存関係**: 標準ライブラリのみ

## 開発

### テスト実行

```bash
cd internal/open_weather_map
go test -v ./...
```

### カバレッジ確認

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。
