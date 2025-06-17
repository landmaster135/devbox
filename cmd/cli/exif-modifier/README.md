# Exif Modifier

任意のフォルダ内にある画像ファイルのExifプロパティを編集するCLIツールです。

## 特徴

- 複数の画像形式をサポート（.jpg, .jpeg, .tiff, .tif, .png, .cr2, .nef, .arw）
- File Modification Date/Time プロパティの編集
- 再帰的なフォルダ処理
- ドライランモード（実際の変更前の確認）
- 特定の拡張子のみの処理

## 必要な依存関係

このツールは `exiftool` を使用します。事前にインストールしてください。

### Ubuntu/Debian
```bash
sudo apt-get install exiftool
```

### macOS
```bash
brew install exiftool
```

### Windows
[ExifTool公式サイト](https://exiftool.org/)からダウンロードしてインストールしてください。

## ビルド

```bash
cd /home/nov/devbox/cmd/cli/exif-modifier
go build -o exif-modifier main.go
```

## 使用方法

### 基本的な使用方法

```bash
# 現在のフォルダの全ての画像ファイルの日時を設定
./exif-modifier --datetime 20240315143000

# 特定のフォルダを処理
./exif-modifier --folder /path/to/images --datetime 20240315143000

# 特定の拡張子のみ処理
./exif-modifier --folder /path/to/images --datetime 20240315143000 --ext .jpg

# サブフォルダも再帰的に処理
./exif-modifier --folder /path/to/images --datetime 20240315143000 --recursive

# ドライランモード（実際には変更せずに確認のみ）
./exif-modifier --folder /path/to/images --datetime 20240315143000 --dry-run

# 詳細な出力を表示
./exif-modifier --folder /path/to/images --datetime 20240315143000 --verbose
```

### オプション

- `--folder`: 処理対象のフォルダパス（デフォルト: 現在のディレクトリ）
- `--datetime`: 設定する日時（yyyyMMddhhmmss形式、必須）
- `--ext`: 対象とする拡張子（例: .jpg, .jpeg, .tiff）
- `--recursive`: サブフォルダも再帰的に処理する
- `--dry-run`: 実際には変更せず、処理対象ファイルのみ表示
- `--verbose`: 詳細な出力を表示

### 日時形式

日時は `yyyyMMddhhmmss` 形式で指定してください。

例：
- `20240315143000` → 2024年3月15日 14時30分00秒
- `20231225120000` → 2023年12月25日 12時00分00秒

## サポートしている画像形式

- JPEG（.jpg, .jpeg）
- TIFF（.tiff, .tif）
- PNG（.png）
- Canon RAW（.cr2）
- Nikon RAW（.nef）
- Sony RAW（.arw）

## 注意事項

- 処理前には必ずバックアップを取ることをお勧めします
- `--dry-run` オプションを使用して、事前に処理対象ファイルを確認してください
- Exifデータが存在しない画像ファイルにはタグが追加されます
- 元のファイルは上書きされます（`exiftool -overwrite_original` を使用）

## 実行例

```bash
# JPEGファイルのみを対象に、2024年3月15日14時30分に設定
./exif-modifier --folder ./photos --datetime 20240315143000 --ext .jpg --verbose

# 全ての画像ファイルを再帰的に処理（ドライラン）
./exif-modifier --folder ./photos --datetime 20240315143000 --recursive --dry-run

# 処理結果例
フォルダ: ./photos
設定する日時: 2024-03-15 14:30:00
対象拡張子: .jpg
再帰処理: false
ドライラン: false

処理中: ./photos/IMG_001.jpg
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました
処理中: ./photos/IMG_002.jpg
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました

処理完了: 2個のファイルを処理しました
```

## トラブルシューティング

### exiftoolが見つからない場合
```bash
# インストール確認
which exiftool

# Ubuntu/Debianでインストール
sudo apt-get update
sudo apt-get install exiftool

# macOSでインストール
brew install exiftool
```

### 権限エラーが発生する場合
```bash
# ファイルの権限を確認
ls -la /path/to/images/

# 必要に応じて権限を変更
chmod 644 /path/to/images/*.jpg
```

### 大量のファイルを処理する場合
```bash
# まずはドライランで確認
./exif-modifier --folder /path/to/images --datetime 20240315143000 --dry-run

# 少数のファイルでテスト
./exif-modifier --folder /path/to/test --datetime 20240315143000 --verbose
```
