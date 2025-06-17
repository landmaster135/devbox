# EXIF Modifier

画像ファイルのEXIF情報を参照・表示するためのCLIツールです。

## 機能

- 指定したディレクトリ内の画像ファイルからEXIF情報を抽出
- テーブル形式での結果表示
- 複数の画像形式をサポート（JPEG, TIFF等）
- 特定のEXIFプロパティのフィルタリング
- 再帰的なディレクトリ検索
- 詳細ログ出力

## インストール

```bash
cd /home/nov/devbox/cmd/cli/exif-modifier
go build -o exif-modifier main.go
```

## 使用方法

### 基本的な使用法

```bash
# カレントディレクトリの画像ファイルのEXIF情報を表示
./exif-modifier

# 特定のディレクトリを指定
./exif-modifier -dir ./photos

# 再帰的にサブディレクトリも検索
./exif-modifier -dir ./photos -r

# 詳細ログを出力
./exif-modifier -dir ./photos -v
```

### 詳細オプション

```bash
# 特定の拡張子のみを対象にする
./exif-modifier -dir ./photos -ext jpg,jpeg

# 特定のEXIFプロパティのみを表示
./exif-modifier -dir ./photos -props DateTime,Camera,GPS

# 組み合わせ使用例
./exif-modifier -dir ./photos -ext jpg,tiff -props DateTime,Make,Model -r -v
```

## コマンドラインオプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-dir` | `.` | 画像ファイルを検索するディレクトリ |
| `-ext` | `jpg,jpeg,tiff,tif` | 対象の画像拡張子（カンマ区切り） |
| `-props` | （全て表示） | 表示するEXIFプロパティ（カンマ区切り） |
| `-r` | `false` | サブディレクトリも再帰的に検索 |
| `-v` | `false` | 詳細なログを出力 |
| `-help` | - | ヘルプメッセージを表示 |
| `-version` | - | バージョン情報を表示 |

## 出力例

```
File Path                    DateTime              Make     Model           GPS
--------------------------------------------------  --------------------  --------------------  --------------------
./photo1.jpg                 2024:01:15 10:30:45   Canon    EOS R5          40.7128,-74.0060
./photo2.jpg                 2024:01:15 11:15:22   Sony     A7R IV          -
./subfolder/photo3.tiff      2024:01:16 09:45:13   Nikon    D850           41.8781,-87.6298

Summary: 3 files processed, 4 unique EXIF properties found
```

## サポートされる画像形式

- JPEG (.jpg, .jpeg)
- TIFF (.tiff, .tif)
- その他のEXIF情報を含む画像形式

## よく使用されるEXIFプロパティ

- `DateTime` - 撮影日時
- `Make` - カメラメーカー
- `Model` - カメラモデル
- `GPS` - GPS座標情報
- `ExposureTime` - 露光時間
- `FNumber` - F値
- `ISO` - ISO感度
- `FocalLength` - 焦点距離
- `WhiteBalance` - ホワイトバランス
- `Flash` - フラッシュ設定

## エラーハンドリング

- EXIF情報が存在しない画像ファイルはスキップされます
- 読み取りエラーが発生したファイルは警告として表示され、処理は継続されます
- 詳細モード（`-v`）でエラーの詳細情報を確認できます

## 依存関係

- `github.com/dsoprea/go-exif/v3` - EXIF情報の読み取り
- `github.com/dsoprea/go-jpeg-image-structure/v2` - JPEG構造の解析
- 標準ライブラリ（flag, fmt, path/filepath, text/tabwriter等）

## ライセンス

このプロジェクトのライセンスに従います。
