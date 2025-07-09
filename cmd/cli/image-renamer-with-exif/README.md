# Image Renamer with EXIF

画像ファイルのEXIF CreateDateプロパティまたはファイルの更新時刻を使用して、ファイル名を年月日時分秒の形式（YYYYMMDDHHMMSS.拡張子）にリネームするCLIツールです。

PowerShellの `X0-01_rename_by_exif.ps1` と同等の機能をGoで実装したツールです。

## 機能

- 画像ファイルのEXIF CreateDateまたはDateTimeOriginalからファイル名を生成（JPEG形式のみ）
- EXIF情報がない場合は、ファイルの更新時刻をフォールバックとして使用
- 複数の画像形式をサポート（JPEG, PNG, TIFF, WebP）
- 再帰的なディレクトリ処理
- ドライランモード（実際の変更を行わずに処理対象を表示）
- 並行処理による高速なファイル処理
- 詳細な処理ログ

## サポートしている画像形式

- JPEG (.jpg, .jpeg) - EXIF CreateDateまたはDateTimeOriginalを使用
- PNG (.png) - ファイルの更新時刻を使用（軽量化のため）
- TIFF (.tiff, .tif) - ファイルの更新時刻を使用（軽量化のため）
- WebP (.webp) - ファイルの更新時刻を使用

## インストール

```bash
go build -o image-renamer-with-exif main.go
```

## 使用方法

### 基本的な使用例

```bash
# 現在のディレクトリの画像ファイルをリネーム
./image-renamer-with-exif

# 特定のディレクトリのJPEGファイルのみリネーム
./image-renamer-with-exif --dir ./photos --ext jpg

# サブディレクトリも含めて再帰的にリネーム
./image-renamer-with-exif --dir ./photos --recursive

# ドライランで処理対象ファイルを確認
./image-renamer-with-exif --dir ./photos --dry-run

# ファイルの更新時刻を使用してリネーム
./image-renamer-with-exif --dir ./photos --use-file-modtime

# 詳細出力でリネーム
./image-renamer-with-exif --dir ./photos --verbose
```

### コマンドラインオプション

```
  -dir string
        画像ファイルがあるディレクトリのパス (default ".")
  -ext string
        対象とする拡張子 (例: jpg, jpeg, png, webp, tiff, mp4, webm)
  -recursive
        サブフォルダも再帰的に処理する
  -dry-run
        実際には変更せず、処理対象ファイルのみ表示
  -verbose
        詳細な出力を表示
  -workers int
        並行処理のワーカー数 (default: CPU数)
  -use-file-modtime
        ExifのCreateDateではなくファイルの更新時刻を使用
  -help
        ヘルプメッセージを表示
  -version
        バージョン情報を表示
```

## 動作仕様

1. **EXIF情報の優先順位（JPEG形式のみ）**: 
   - `DateTime`（作成日時）
   - `DateTimeOriginal`（撮影日時）
   - ファイルの更新時刻（フォールバック）

2. **その他形式の処理**:
   - PNG, TIFF, WebP: 依存関係を減らすためファイルの更新時刻を使用

3. **ファイル名形式**: `YYYYMMDDHHMMSS.拡張子`
   - 例: `20240315143025.jpg`

4. **重複ファイル名の処理**: 
   - 同じ名前のファイルが存在する場合、連番を追加
   - 例: `20240315143025_01.jpg`, `20240315143025_02.jpg`

5. **エラーハンドリング**: 
   - 個別ファイルのエラーは全体の処理を中断しない
   - エラーファイルは最終的にカウントして報告

## PowerShellスクリプトとの対応

このツールは以下のPowerShellコマンドと同等の機能を提供します：

```powershell
# PowerShell版
exiftool -FileName<CreateDate -d "%Y%m%d%H%M%S.%%e" ./

# Go版
./image-renamer-with-exif --dir ./
```

## 技術詳細

- **言語**: Go
- **EXIFライブラリ**: `github.com/dsoprea/go-exif/v3`
- **画像構造解析**: 
  - JPEG: `github.com/dsoprea/go-jpeg-image-structure/v2`
  - PNG, TIFF, WebP: ファイルシステムレベルの更新時刻を使用（軽量化）
- **並行処理**: Goroutineとチャネルを使用

## 注意事項

- ファイルのリネームは元に戻すことができません。必要に応じて事前にバックアップを作成してください
- `--dry-run` オプションを使用して、実際の処理前に対象ファイルを確認することを推奨します
- JPEG以外の形式（PNG, TIFF, WebP）では、依存関係を減らすためファイルの更新時刻を使用します

## トラブルシューティング

### EXIF情報が読み取れない場合
- ファイルの更新時刻を使用してリネームされます
- `--verbose` オプションで詳細なログを確認できます

### ファイルアクセス権限エラー
- 対象ディレクトリとファイルに適切な読み書き権限があることを確認してください

### メモリ使用量が多い場合
- `--workers` オプションで並行処理数を調整してください

## 開発者向け情報

### ビルド

```bash
# 通常のビルド
go build -o image-renamer-with-exif main.go

# クロスコンパイル（Windows用）
GOOS=windows GOARCH=amd64 go build -o image-renamer-with-exif.exe main.go

# クロスコンパイル（Linux用）
GOOS=linux GOARCH=amd64 go build -o image-renamer-with-exif main.go
```

### テスト

```bash
go test ./...
```

### 依存関係の更新

```bash
go mod tidy
go mod download
```

## ライセンス

このプロジェクトは、親プロジェクトのライセンスに従います。
