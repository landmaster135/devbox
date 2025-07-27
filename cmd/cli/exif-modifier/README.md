# Exif Modifier

任意のディレクトリ内にある画像・動画ファイルのExifプロパティを編集するCLIツールです。

## 特徴

- **高速並行処理対応**：複数ファイルを同時に処理して大幅な処理時間短縮を実現
- **拡張子指定の簡略化**：ドット（.）なしで拡張子を指定可能（例：`jpg` instead of `.jpg`）
- **ファイル名順ソート**：処理順序が一定でログ確認が容易
- 複数の画像・動画形式をサポート（.jpg, .jpeg, .tiff, .tif, .png, .webp, .mp4, .webm。デフォルト：`jpg`）
- File Modification Date/Time プロパティの編集
- **ファイル名から自動的に日時を抽出する機能**
- **スクリーンショットファイル名から日時を抽出する機能**
- 再帰的なディレクトリ処理
- ドライランモード（実際の変更前の確認）
- 特定の拡張子のみの処理

## 必要な依存関係

このツールはgo-exifライブラリを使用してExifデータの存在確認を行い、ファイルシステムレベルでの時刻更新を実装しています。

## ビルド

```bash
cd ~/devbox/cmd/cli/exif-modifier
go build -o exif-modifier main.go
```

## 使用方法

### 基本的な使用方法（固定日時を設定）

```bash
# 現在のディレクトリの全てのファイルの日時を設定
./exif-modifier --datetime 20240315143000

# 特定のディレクトリを処理
./exif-modifier --dir ./photos --datetime 20240315143000

# 特定の拡張子のみ処理（ドット不要！）
./exif-modifier --dir ./photos --datetime 20240315143000 --ext jpg

# 並行処理のワーカー数を指定（デフォルト: 4）
./exif-modifier --dir ./photos --datetime 20240315143000 --workers 8

# サブディレクトリも再帰的に処理
./exif-modifier --dir ./photos --datetime 20240315143000 --recursive

# ドライランモード（実際には変更せずに確認のみ）
./exif-modifier --dir ./photos --datetime 20240315143000 --dry-run

# 詳細な出力を表示
./exif-modifier --dir ./photos --datetime 20240315143000 --verbose
```

### ファイル名から日時を自動抽出

ファイル名が「yyyyMMddhhmmss」の形式になっている場合、その日時を自動的に抽出してExifに設定できます。

```bash
# ファイル名から日時を抽出して設定（並行処理）
./exif-modifier --from-filename --dir ./photos

# 特定の拡張子のみ対象にする（ドット不要）
./exif-modifier --from-filename --dir ./photos --ext jpg

# 高速処理のためにワーカー数を増やす
./exif-modifier --from-filename --dir ./photos --workers 8

# サブディレクトリも再帰的に処理する
./exif-modifier --from-filename --dir ./photos --recursive

# ドライランで確認
./exif-modifier --from-filename --dir ./photos --dry-run

# 詳細モードで実行
./exif-modifier --from-filename --dir ./photos --verbose
```

### スクリーンショットファイル名から日時を抽出（新機能）

「Screenshot_yyyyMMdd-hhmmss」形式のファイル名から日時を自動抽出できます。

```bash
# スクリーンショットファイルの日時を抽出
./exif-modifier --from-screenshot --dir ./screenshots

# 特定の拡張子のみ対象（並行処理）
./exif-modifier --from-screenshot --dir ./screenshots --ext png --workers 6

# ドライランで確認
./exif-modifier --from-screenshot --dir ./screenshots --dry-run
```

#### 対応ファイル名形式の例

**通常のファイル名形式：**
- `20240315143000.jpg` → 2024年3月15日 14時30分00秒
- `20231225120000.png` → 2023年12月25日 12時00分00秒
- `20220101000000.mp4` → 2022年1月1日 00時00分00秒

**スクリーンショット形式：**
- `Screenshot_20240315-143000.png` → 2024年3月15日 14時30分00秒
- `Screenshot_20231225-120000.jpg` → 2023年12月25日 12時00分00秒

### オプション

- `--dir`: 処理対象のディレクトリパス（デフォルト: 現在のディレクトリ）
- `--datetime`: 設定する日時（yyyyMMddhhmmss形式）
- `--from-filename`: ファイル名から日時を取得してExifに設定する（ファイル名がyyyyMMddhhmmss形式の場合）
- `--from-screenshot`: スクリーンショットファイル名から日時を取得してExifに設定する（Screenshot_yyyyMMdd-hhmmss形式の場合）
- `--ext`: 対象とする拡張子（例: jpg, jpeg, png, webp, mp4）**ドット不要**
- `--workers`: 並行処理のワーカー数（デフォルト: 4）
- `--recursive`: サブディレクトリも再帰的に処理する
- `--dry-run`: 実際には変更せず、処理対象ファイルのみ表示
- `--verbose`: 詳細な出力を表示

**注意**: `--datetime`、`--from-filename`、`--from-screenshot` は同時に指定できません。いずれか一つを選択してください。

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

## 実行例

### 固定日時を設定する場合（並行処理対応）

```bash
# JPEGファイルのみを対象に、2024年3月15日14時30分に設定（並行処理）
./exif-modifier --dir ./photos --datetime 20240315143000 --ext jpg --workers 8 --verbose

# 全てのファイルを再帰的に処理（ドライラン）
./exif-modifier --dir ./photos --datetime 20240315143000 --recursive --dry-run

# 処理結果例
ディレクトリ: ./photos
設定する日時: 2024-03-15 14:30:00
対象拡張子: jpg
再帰処理: false
ドライラン: false
ワーカー数: 8

処理中: ./photos/IMG_001.jpg
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました
処理中: ./photos/video.mp4
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました

処理完了: 2個のファイルを処理しました
```

### ファイル名から日時を自動抽出する場合（並行処理対応）

```bash
# ファイル名から日時を抽出（詳細モード、並行処理）
./exif-modifier --from-filename --dir ./photos --workers 6 --verbose

# 処理結果例
ディレクトリ: ./photos
モード: ファイル名から日時を取得
再帰処理: false
ドライラン: false
詳細モード: true
ワーカー数: 6

2024/06/17 12:00:00 ファイル: 20240315143000 -> 日時: 2024-03-15 14:30:00
処理中: ./photos/20240315143000.jpg
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました

2024/06/17 12:00:01 ファイル名が日時形式ではありません（スキップ）: IMG_001
2024/06/17 12:00:01 ファイル: 20231225120000 -> 日時: 2023-12-25 12:00:00
処理中: ./photos/20231225120000.png
  ✅ File Modification Date/Time を 2023-12-25 12:00:00 に設定しました

処理完了: 2個のファイルを処理しました
```

### スクリーンショットファイルの処理例

```bash
# スクリーンショットファイルから日時を抽出（並行処理）
./exif-modifier --from-screenshot --dir ./screenshots --ext png --workers 4 --verbose

# 処理結果例
ディレクトリ: ./screenshots
モード: スクリーンショットファイル名から日時を取得
対象拡張子: png
再帰処理: false
ドライラン: false
詳細モード: true
ワーカー数: 4

2024/06/17 12:00:00 ファイル: Screenshot_20240315-143000 -> 日時: 2024-03-15 14:30:00
処理中: ./screenshots/Screenshot_20240315-143000.png
  ✅ File Modification Date/Time を 2024-03-15 14:30:00 に設定しました

処理完了: 1個のファイルを処理しました
```

### ドライランモードの例

```bash
# ファイル名から日時を抽出（ドライラン）
./exif-modifier --from-filename --dir ./photos --dry-run

# 処理結果例
ディレクトリ: ./photos
モード: ファイル名から日時を取得
再帰処理: false
ドライラン: true

[DRY RUN] ./photos/20240315143000.jpg -> 2024-03-15 14:30:00
[DRY RUN] ./photos/20231225120000.png -> 2023-12-25 12:00:00

処理完了: 2個のファイルを処理しました
```

## 注意事項

- 処理前には必ずバックアップを取ることをお勧めします
- `--dry-run` オプションを使用して、事前に処理対象ファイルを確認してください
- `--from-filename` モードでは、ファイル名が「yyyyMMddhhmmss」形式でない場合はスキップされます
- `--from-screenshot` モードでは、ファイル名が「Screenshot_yyyyMMdd-hhmmss」形式でない場合はスキップされます
- `--datetime`、`--from-filename`、`--from-screenshot` は同時に指定できません
- 現在の実装では、ファイルシステムレベルでの更新時刻の変更を行います
- JPEGファイルでExifデータが存在する場合は検出されますが、実際のExifタグの編集は将来の機能拡張で対応予定です

## パフォーマンス

### 並行処理による高速化

- **デフォルト**: 4並列でファイルを処理
- **カスタム**: `--workers` オプションでワーカー数を調整可能（1-16推奨）
- **推奨設定**: CPUコア数と同じかやや多い程度（例：8コアCPUなら`--workers 8`または`--workers 10`）

### 処理速度の目安

```bash
# 例：1000ファイルの処理時間比較
# シングルスレッド: 約30秒
# 4並列（デフォルト）: 約8秒  
# 8並列: 約4秒

# 大量ファイル処理の例
./exif-modifier --from-filename --dir ./large_photo_dir --workers 8 --verbose
```

### メモリ使用量

- 各ワーカーは独立してファイルを処理するため、ワーカー数に比例してメモリ使用量が増加
- 通常のファイルサイズ（数MB程度）であれば、8並列でも問題なし
- 非常に大きなファイル（100MB以上）を大量処理する場合は、ワーカー数を減らすことを推奨

## トラブルシューティング

### 権限エラーが発生する場合
```bash
# ファイルの権限を確認
ls -la /path/to/files/

# 必要に応じて権限を変更
chmod 644 /path/to/files/*
```

### ファイル名から日時が抽出されない場合
```bash
# ファイル名の形式を確認
ls -la /path/to/files/

# 正しい形式の例: 
#   20240315143000.jpg (通常形式)
#   Screenshot_20240315-143000.png (スクリーンショット形式)
# 間違った形式の例: 
#   IMG_20240315_143000.jpg, 2024-03-15-14-30-00.jpg
```

### 大量のファイルを処理する場合
```bash
# まずはドライランで確認
./exif-modifier --from-filename --dir /path/to/files --dry-run

# 少数のファイルでテスト
./exif-modifier --from-filename --dir /path/to/test --verbose

# 並行処理数を調整してパフォーマンステスト
./exif-modifier --from-filename --dir /path/to/files --workers 8 --verbose
```

### パフォーマンスが上がらない場合
```bash
# I/O待機が原因の可能性：SSD推奨、ネットワークドライブは避ける
# ワーカー数を調整してテスト
./exif-modifier --datetime 20240315143000 --workers 4   # 少なめ
./exif-modifier --datetime 20240315143000 --workers 8   # 標準
./exif-modifier --datetime 20240315143000 --workers 16  # 多め
```

## 使用例：カメラのスクリーンショットファイルの処理

多くのカメラアプリケーションは「yyyyMMddhhmmss.jpg」のようなファイル名でスクリーンショットを保存します。このツールを使用することで、これらのファイルのExif情報を自動的に正しい日時に設定できます。

```bash
# カメラのスクリーンショットディレクトリを処理（並行処理）
./exif-modifier --from-filename --dir ./screenshots --ext jpg --workers 6 --verbose

# 結果をドライランで事前確認
./exif-modifier --from-filename --dir ./screenshots --ext jpg --dry-run

# スクリーンショット形式のファイルを処理
./exif-modifier --from-screenshot --dir ./screenshots --ext png --verbose
```

## バージョン履歴

### v1.1.0 (2025-06-18)
- **並行処理対応**: 複数ファイルを同時処理して大幅な処理時間短縮を実現
- **拡張子指定の簡略化**: ドット（.）なしで拡張子を指定可能
- **ファイル名順ソート**: 処理順序が一定でログ確認が容易
- **スクリーンショット形式対応**: Screenshot_yyyyMMdd-hhmmss形式のファイル名に対応
- **ワーカー数調整**: `--workers` オプションで並行処理数を調整可能

### v1.0.0
- 初回リリース
- 基本的なExif修正機能
- ファイル名から日時抽出機能
