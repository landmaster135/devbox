# Image Renamer

画像ファイルを指定したプレフィックスとシリアル番号でリネームするツールです。

## 機能

- 指定したディレクトリ内の画像ファイル（jpg, jpeg, png, webp, avif）を検索
- ファイルを更新日時順またはファイル名順に並べ替え
- 指定したプレフィックスとシリアル番号でファイルをリネーム
- 並行処理による高速なリネーム操作
- 再帰的なディレクトリ走査オプション

## 使用方法

```
./image-renamer [オプション]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| -src       | .            | スキャンするソースディレクトリ |
| -time      | false        | 画像ファイルを更新日時順に並べ替え |
| -name      | false        | 画像ファイルをファイル名順に並べ替え |
| -prefix    | (必須)       | 記事番号のプレフィックス |
| -delimiter | _            | プレフィックスとシリアル番号の間の区切り文字 |
| -digits    | 4            | シリアル番号の桁数 |
| -start     | 1            | リネーム操作の開始番号 |
| -r         | false        | サブディレクトリを再帰的にスキャン |
| -workers   | CPU数        | 並行ワーカー数 |
| -extensions | (未指定)    | リネーム対象拡張子（カンマ区切り）。未指定時は `jpg,jpeg,png,webp,avif` |

**注意**: `-time` または `-name` のいずれかを指定する必要があります。両方指定した場合は `-name` が優先されます。

## 使用例

### 基本的な使用方法

```bash
# カレントディレクトリの画像ファイルを日付順に並べ替えてリネーム
go run ./cmd/cli/image-renamer -prefix "20250507" -time

# 指定したディレクトリの画像ファイルをファイル名順に並べ替えてリネーム
go run ./cmd/cli/image-renamer -src ./photos -prefix "article01" -name

# 3桁のシリアル番号を使用（例: article01_001.jpg）
go run ./cmd/cli/image-renamer -prefix "article01" -digits 3 -time

# 開始番号を10から始める
go run ./cmd/cli/image-renamer -prefix "article01" -start 10 -time

# カスタム区切り文字を使用（例: article01-001.jpg）
go run ./cmd/cli/image-renamer -prefix "article01" -delimiter "-" -time
```

### 再帰的なスキャン

```bash
# サブディレクトリも含めて画像ファイルをリネーム
go run ./cmd/cli/image-renamer -src ./photos -prefix "article01" -time -r
```

### 並行処理の調整

```bash
# ワーカー数を8に設定
go run ./cmd/cli/image-renamer -prefix "article01" -time -workers 8

# HEICとPNGのみを対象にリネーム
go run ./cmd/cli/image-renamer -prefix "article01" -name -extensions ".heic,.png"
```

## ビルド方法

付属のビルドスクリプトを使用してビルドできます：

```bash
./scripts/build_image_renamer.sh
```

これにより、Linux、Windows、macOS向けのバイナリが `pkg/bin` ディレクトリに生成されます。
