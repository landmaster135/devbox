# EXIF Modifier

画像ファイルのEXIF・メタデータ情報を参照・表示するためのCLIツールです。

## 機能

- 指定したディレクトリ内の画像ファイルからEXIF・メタデータ情報を抽出
- テーブル形式での結果表示
- 複数の画像形式をサポート（JPEG, TIFF, PNG等）
- PNGファイルの詳細メタデータ読み取り対応
- 特定のプロパティのフィルタリング
- 再帰的なディレクトリ検索
- 詳細ログ出力
- ファイル情報（サイズ、日時、種別）も表示

## インストール

```bash
cd /home/nov/devbox/cmd/cli/exif-modifier
go build -o exif-modifier main.go
```

## 使用方法

### 基本的な使用法

```bash
# カレントディレクトリの画像ファイルの情報を表示
./exif-modifier

# 特定のディレクトリを指定
./exif-modifier -dir ./photos

# PNGファイルのメタデータを表示
./exif-modifier -dir ./screenshots -ext png

# 再帰的にサブディレクトリも検索
./exif-modifier -dir ./photos -r

# 詳細ログを出力
./exif-modifier -dir ./photos -v
```

### 詳細オプション

```bash
# 特定の拡張子のみを対象にする
./exif-modifier -dir ./photos -ext jpg,jpeg

# 特定のプロパティのみを表示
./exif-modifier -dir ./photos -props "File Size,Image Width,Image Height"

# 組み合わせ使用例
./exif-modifier -dir ./photos -ext jpg,tiff,png -props "File Name,Image Size,File Type" -r -v
```

## コマンドラインオプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-dir` | `.` | 画像ファイルを検索するディレクトリ |
| `-ext` | `jpg,jpeg,tiff,tif,png` | 対象の画像拡張子（カンマ区切り） |
| `-props` | （全て表示） | 表示するプロパティ（カンマ区切り） |
| `-r` | `false` | サブディレクトリも再帰的に検索 |
| `-v` | `false` | 詳細なログを出力 |
| `-help` | - | ヘルプメッセージを表示 |
| `-version` | - | バージョン情報を表示 |

## 出力例

### PNG スクリーンショット

```
File Path                          File Name                      File Size  Image Width  Image Height  Image Size  File Type
--------------------------------------------------  --------------------  ----------   -----------  -----------   ---------   ---------
screenshot1.png                    screenshot1.png               341 kB     1746         1316          1746x1316   PNG
screenshot2.png                    screenshot2.png               523 kB     1920         1080          1920x1080   PNG

Summary: 2 files processed, 6 unique properties found
```

### JPEG 写真

```
File Path     DateTime               Make     Model    File Size  Image Size
----------    --------------------   ------   -------  ---------  ----------
photo1.jpg    2024:01:15 10:30:45    Canon    EOS R5   8.2 MB     6000x4000
photo2.jpg    2024:01:15 11:15:22    Sony     A7R IV   12.1 MB    7952x5304

Summary: 2 files processed, 5 unique properties found
```

## サポートされる画像形式

- **JPEG (.jpg, .jpeg)** - 完全なEXIF情報
- **PNG (.png)** - ファイル情報 + PNG メタデータ（チャンク情報）
- **TIFF (.tiff, .tif)** - EXIF情報 + ファイル情報
- **その他のEXIF対応形式**

## 取得できる情報の種類

### 共通情報（全形式）
- `File Name` - ファイル名
- `File Size` - ファイルサイズ
- `File Modification Date/Time` - 更新日時
- `Directory` - ディレクトリパス
- `File Type` - ファイル形式
- `File Type Extension` - 拡張子
- `MIME Type` - MIMEタイプ

### PNG固有の情報
- `Image Width` / `Image Height` - 画像サイズ
- `Image Size` - サイズ（幅x高さ）
- `Megapixels` - メガピクセル数
- `Bit Depth` - ビット深度
- `Color Type` - カラータイプ
- `Compression` - 圧縮方式
- `Filter` - フィルター方式
- `Interlace` - インターレース
- `Pixels Per Unit X/Y` - 解像度
- `Pixel Units` - 解像度の単位
- `Gamma` - ガンマ値
- `SRGB Rendering` - sRGBレンダリング

### JPEG EXIF情報（カメラによる）
- `DateTime` - 撮影日時
- `Make` - カメラメーカー
- `Model` - カメラモデル
- `ExposureTime` - 露光時間
- `FNumber` - F値
- `ISO` - ISO感度
- `FocalLength` - 焦点距離
- `WhiteBalance` - ホワイトバランス
- `Flash` - フラッシュ設定

## エラーハンドリング

- EXIF情報が存在しない画像ファイルでも、ファイル情報は表示されます
- 読み取りエラーが発生したファイルは警告として表示され、処理は継続されます
- 詳細モード（`-v`）でエラーの詳細情報を確認できます
- PNGファイルではEXIFではなくPNGメタデータ（チャンク）を読み取ります

## 依存関係

- `github.com/dsoprea/go-exif/v3` - EXIF情報の読み取り
- `github.com/dsoprea/go-jpeg-image-structure/v2` - JPEG構造の解析
- 標準ライブラリ（flag, fmt, path/filepath, text/tabwriter, image/png等）

## 技術仕様

- PNGファイルはPNGチャンク構造を直接パースしてメタデータを取得
- JPEGファイルは専用ライブラリでEXIF情報を抽出
- その他の画像形式は汎用的なEXIF読み取りを試行
- UTF-8文字コード対応
- クロスプラットフォーム対応（Windows, macOS, Linux）

## ライセンス

このプロジェクトのライセンスに従います。
