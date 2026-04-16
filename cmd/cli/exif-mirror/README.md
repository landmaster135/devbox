# EXIF Mirror Tool

指定されたソースフォルダ内の画像ファイルから、ターゲットフォルダ内の同名ファイルにEXIFデータをコピーするCLIツールです。

## 機能

- JPEGファイルからJPEG/PNG/WebPファイルへのEXIFデータコピー
- PNGファイルからJPEG/PNG/WebPファイルへのEXIFデータコピー（限定的）
- TIFFファイルからの読み取り（限定的サポート）
- 再帰的なディレクトリ処理
- 並行処理によるパフォーマンス向上
- ドライランモード（実際の変更なしでプレビュー）
- 詳細ログ出力

## 使用方法

### 基本的な使用例

```bash
# JPGファイルからWebPファイルにEXIFをコピー
./exif-mirror --source-dir ./originals --target-dir ./converted --source-ext jpg --target-ext webp

# PNGファイルからWebPファイルにEXIFをコピー（再帰処理）
./exif-mirror --source-dir ./photos --target-dir ./photos --source-ext png --target-ext webp --recursive

# ドライラン（実際の処理は行わず、対象ファイルのみ表示）
./exif-mirror --source-dir ./src --target-dir ./dst --source-ext jpg --target-ext png --dry-run

# 詳細ログ付きで実行
./exif-mirror --source-dir ./src --target-dir ./dst --source-ext jpg --target-ext webp --verbose
```

### コマンドラインオプション

- `--source-dir`: ソース画像ファイルがあるディレクトリのパス（必須）
- `--target-dir`: ターゲット画像ファイルがあるディレクトリのパス（必須）
- `--source-ext`: ソース拡張子（必須）。例: jpg, jpeg, png
- `--target-ext`: ターゲット拡張子（必須）。例: webp, png, jpg
- `--recursive`: サブフォルダも再帰的に処理する
- `--dry-run`: 実際には変更せず、処理対象ファイルのみ表示
- `--verbose`: 詳細な出力を表示
- `--workers`: 並行処理のワーカー数（デフォルト: CPU数）
- `--help`: ヘルプメッセージを表示
- `--version`: バージョン情報を表示

## 動作の仕組み

1. ターゲットディレクトリ内で指定された拡張子のファイルを検索
2. 各ターゲットファイルに対して、対応するソースファイルを検索
   - ファイル名（拡張子除く）が同じファイルを検索
   - ソースディレクトリ内の対応する相対パスで検索
3. ソースファイルからEXIFデータを抽出
4. ターゲットファイルにEXIFデータを書き込み

## ファイル対応関係の例

```
ソースディレクトリ:
├── photo1.jpg
├── photo2.jpg
└── subfolder/
    └── photo3.jpg

ターゲットディレクトリ:
├── photo1.webp    <- photo1.jpgからEXIFをコピー
├── photo2.webp    <- photo2.jpgからEXIFをコピー
└── subfolder/
    └── photo3.webp <- photo3.jpgからEXIFをコピー
```

# EXIF Mirror Tool

PowerShellの`Copy-ExifFromSrc`コマンドレットと同じ機能を提供するGo製CLIツールです。指定されたソースフォルダ内の画像ファイルから、ターゲットフォルダ内の同名ファイルにEXIFデータをコピーします。

## 機能

- ソースファイルからターゲットファイルへのEXIFデータコピー
- exiftoolを使用した信頼性の高いEXIF処理
- 多様な画像フォーマットサポート（JPEG、PNG、WebP、TIFF等）
- 再帰的なディレクトリ処理
- 並行処理によるパフォーマンス向上
- ドライランモード（実際の変更なしでプレビュー）
- 詳細ログ出力

## 前提条件

このツールは**exiftool**を使用してEXIFデータの処理を行います。事前にexiftoolをインストールしてください。

### exiftoolのインストール

#### Windows
1. [ExifTool公式サイト](https://exiftool.org/)からダウンロード
2. exiftool.exeをPATHの通った場所に配置

#### macOS
```bash
brew install exiftool
```

#### Ubuntu/Debian
```bash
sudo apt-get install exiftool
```

#### CentOS/RHEL
```bash
sudo yum install perl-Image-ExifTool
```

## 使用方法

### 基本的な使用例

```bash
# JPGファイルからWebPファイルにEXIFをコピー
./exif-mirror --source-dir ./originals --target-dir ./converted --source-ext jpg --target-ext webp

# PNGファイルからWebPファイルにEXIFをコピー（再帰処理）
./exif-mirror --source-dir ./photos --target-dir ./photos --source-ext png --target-ext webp --recursive

# ドライラン（実際の処理は行わず、対象ファイルのみ表示）
./exif-mirror --source-dir ./src --target-dir ./dst --source-ext jpg --target-ext png --dry-run

# 詳細ログ付きで実行
./exif-mirror --source-dir ./src --target-dir ./dst --source-ext jpg --target-ext webp --verbose
```

### コマンドラインオプション

- `--source-dir`: ソース画像ファイルがあるディレクトリのパス（必須）
- `--target-dir`: ターゲット画像ファイルがあるディレクトリのパス（必須）
- `--source-ext`: ソース拡張子（必須）。例: jpg, jpeg, png
- `--target-ext`: ターゲット拡張子（必須）。例: webp, png, jpg
- `--recursive`: サブフォルダも再帰的に処理する
- `--dry-run`: 実際には変更せず、処理対象ファイルのみ表示
- `--verbose`: 詳細な出力を表示
- `--workers`: 並行処理のワーカー数（デフォルト: CPU数）
- `--help`: ヘルプメッセージを表示
- `--version`: バージョン情報を表示

## 動作の仕組み

1. ターゲットディレクトリ内で指定された拡張子のファイルを検索
2. 各ターゲットファイルに対して、対応するソースファイルを検索
   - ファイル名（拡張子除く）が同じファイルを検索
   - ソースディレクトリ内の対応する相対パスで検索
3. exiftoolを使用してソースファイルからターゲットファイルにEXIFデータをコピー

## ファイル対応関係の例

```
ソースディレクトリ:
├── photo1.jpg
├── photo2.jpg
└── subfolder/
    └── photo3.jpg

ターゲットディレクトリ:
├── photo1.webp    <- photo1.jpgからEXIFをコピー
├── photo2.webp    <- photo2.jpgからEXIFをコピー
└── subfolder/
    └── photo3.webp <- photo3.jpgからEXIFをコピー
```

## サポートされるフォーマット

exiftoolを使用するため、exiftoolがサポートするすべての画像・動画フォーマットに対応しています：

### 主要なサポート形式
- **JPEG** (.jpg, .jpeg) - 完全サポート
- **PNG** (.png) - 完全サポート
- **WebP** (.webp) - 完全サポート
- **TIFF** (.tiff, .tif) - 完全サポート
- **HEIC/HEIF** (.heic, .heif) - 完全サポート
- **DNG** (.dng) - 完全サポート
- **CR2, NEF, ARW等のRAW形式** - 完全サポート
- **MP4, MOV等の動画形式** - 完全サポート

## ビルド方法

```bash
# リポジトリのルートディレクトリから
go build -o exif-mirror ./cmd/cli/exif-mirror
```

## テスト実行

```bash
# ユニットテストの実行
go test ./internal/exif_mirror/usecases

# カバレッジ付きテスト
go test -cover ./internal/exif_mirror/usecases
```

## エラーハンドリング

- exiftoolが見つからない場合は明確なエラーメッセージを表示
- ソースファイルが見つからない場合はスキップ
- ソースファイルにEXIFデータがない場合は警告を表示
- 並行処理中のエラーは個別に処理され、全体の処理は継続

## PowerShellとの互換性

このツールはPowerShellの`Copy-ExifFromSrc`コマンドレットと同じ機能を提供します：

```powershell
# PowerShell版
Copy-ExifFromSrc $logtxt ".jpg"

# Go版での同等処理
./exif-mirror --source-dir . --target-dir . --source-ext jpg --target-ext webp
```

## パフォーマンス

- 並行処理により複数ファイルを同時処理
- デフォルトでCPU数と同じワーカー数を使用
- `--workers`オプションで並行度を調整可能

## トラブルシューティング

### exiftool not foundエラー
```
Error: EXIF copying requires exiftool to be installed
```
- exiftoolがインストールされていない、またはPATHに含まれていません
- 上記の「前提条件」セクションを参照してexiftoolをインストールしてください

### Permission denied エラー
- ファイルやディレクトリの読み書き権限を確認してください
- Windowsの場合、管理者権限で実行してみてください

## 制限事項

1. **exiftool依存**: exiftoolが必要です（単体では動作しません）
2. **大量ファイル処理**: 数万ファイルの処理では時間がかかる場合があります
3. **ディスク容量**: ファイル書き換え時に一時的に追加容量が必要です

## ライセンス

このプロジェクトは適切なライセンス下で公開されています。

## 貢献

バグレポートや機能改善の提案は歓迎します。プルリクエストを送信する前に、テストが通ることを確認してください。

## ビルド方法

```bash
# リポジトリのルートディレクトリから
go build -o exif-mirror ./cmd/cli/exif-mirror
```

## テスト実行

```bash
# ユニットテストの実行
go test ./internal/exif_mirror/usecases

# カバレッジ付きテスト
go test -cover ./internal/exif_mirror/usecases
```

## 制限事項

1. **WebPサポート**: go-exifライブラリの制限により、WebPファイルへの書き込みは現在サポートされていません
2. **TIFFサポート**: TIFFファイルのEXIF書き込みは複雑で、現在は限定的なサポートです
3. **EXIFデータの完全性**: 一部の特殊なEXIFタグは正しく転送されない場合があります

## エラーハンドリング

- ソースファイルが見つからない場合はスキップ
- ソースファイルにEXIFデータがない場合はスキップ
- サポートされていないフォーマットの場合はエラーメッセージを表示
- 並行処理中のエラーは個別に処理され、全体の処理は継続

## PowerShellとの互換性

このツールはPowerShellの`Copy-ExifFromSrc`コマンドレットと同じ機能を提供します：

```powershell
# PowerShell版
Copy-ExifFromSrc $logtxt ".jpg"

# Go版での同等処理
./exif-mirror --source-dir . --target-dir . --source-ext jpg --target-ext webp
```

## ライセンス

このプロジェクトは適切なライセンス下で公開されています。

## 貢献

バグレポートや機能改善の提案は歓迎します。プルリクエストを送信する前に、テストが通ることを確認してください。
