# Exif Modifier

任意のディレクトリ内にある画像・動画ファイルのExifプロパティを編集するCLIツールです。

## 特徴

- 複数の画像・動画形式をサポート（.jpg, .jpeg, .tiff, .tif, .png, .webp, .mp4, .webm）
- File Modification Date/Time プロパティの編集
- 再帰的なディレクトリ処理
- ドライランモード（実際の変更前の確認）
- 特定の拡張子のみの処理

## 必要な依存関係

このツールはgo-exifライブラリを使用してExifデータの存在確認を行い、ファイルシステムレベルでの時刻更新を実装しています。

## ビルド

```bash
cd /home/nov/devbox/cmd/cli/exif-modifier
go build -o exif-modifier main.go
```

## 使用方法

### 基本的な使用方法

```bash
# 現在のディレクトリの全てのファイルの日時を設定
./exif-modifier --datetime 20240315143000

# 特定のディレクトリを処理
./exif-modifier --dir ./photos --datetime 20240315143000

# 特定の拡張子のみ処理
./exif-modifier --dir ./photos --datetime 20240315143000 --ext .jpg

# サブディレクトリも再帰的に処理
./exif-modifier --dir ./photos --datetime 20240315143000 --recursive

# ドライランモード（実際には変更せずに確認のみ）
./exif-modifier --dir ./photos --datetime 20240315143000 --dry-run

# 詳細な出力を表示
./exif-modifier --dir ./photos --datetime 20240315143000 --verbose
```

### オプション

- `--dir`: 処理対象のディレクトリパス（デフォルト: 現在のディレクトリ）
- `--datetime`: 設定する日時（yyyyMMddhhmmss形式、必須）
- `--ext`: 対象とする拡張子（例: .jpg, .jpeg, .png, .webp, .mp4）
- `--recursive`: サブディレクトリも再帰的に処理する
- `--dry-run`: 実際には変更せず、処理対象ファイルのみ表示
- `--verbose`: 詳細な出力を表示

### 日時形式

日時は `yyyyMMddhhmmss` 形式で指定してください。

例：
- `20240315143000` → 2024年3月15日 14時30分00秒
- `20231225120000` → 2023年12月25日 12時00分00秒

## サポートしている形式

### 画像形式
- JPEG（.jpg, .jpeg）
- TIFF（.tiff, .tif）
- PNG（.png）
- WebP（.webp）

### 動画形式
- MP4（.mp4）
- WebM（.webm）

## 注意事項

- 処理前には必ずバックアップを取ることをお勧めします
- `--dry-run` オプションを使用して、事前に処理対象ファイルを確認してください
- 現在の実装では、ファイルシステムレベルでの更新時刻の変更を行います
- JPEGファイルでExifデータが存在する場合は検出されますが、実際のExifタグの編集は将来の機能拡張で対応予定です

## 実行例

```bash
# JPEGファイルのみを対象に、2024年3月15日14時30分に設定
./exif-modifier --dir ./photos --datetime 20240315143000 --ext .jpg --verbose

# 全てのファイルを再帰的に処理（ドライラン）
./exif-modifier --dir ./photos --datetime 20240315143000 --recursive --dry-run

# 処理結果例
ディレクトリ: ./photos
設定する日時: 2024-03-15 14:30:00
対象拡張子: .jpg
再帰処理: false
ドライラン: false

処理中: ./photos/IMG_001.jpg
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました
処理中: ./photos/video.mp4
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました

処理完了: 2個のファイルを処理しました
```

## トラブルシューティング

### 権限エラーが発生する場合
```bash
# ファイルの権限を確認
ls -la /path/to/files/

# 必要に応じて権限を変更
chmod 644 /path/to/files/*
```

### 大量のファイルを処理する場合
```bash
# まずはドライランで確認
./exif-modifier --dir /path/to/files --datetime 20240315143000 --dry-run

# 少数のファイルでテスト
./exif-modifier --dir /path/to/test --datetime 20240315143000 --verbose
```
