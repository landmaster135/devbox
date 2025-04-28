# 画像変換ツール (Image Converter)

複数の画像ファイルを一括で別のフォーマットに変換するコマンドラインツールです。

## 機能

- 複数の画像ファイルを一括変換
- PNG, JPEG, WebP, AVIFなど様々なフォーマット間の変換をサポート
- 非可逆圧縮フォーマットの品質調整
- 複数のCPUコアを活用した並列処理
- 再帰的なディレクトリ走査オプション

## 使い方

```bash
image-converter [オプション]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-src` | `.` | 変換元ディレクトリ |
| `-out` | `./999_converted_images` | 出力先ディレクトリ |
| `-ext` | `png` | 変換先フォーマット (png/jpg/webp/avif) |
| `-q` | `80` | 非可逆圧縮フォーマットの品質 (1-100) |
| `-workers` | CPU数 | 同時実行ワーカー数 |
| `-R` | `false` | サブディレクトリを再帰的に処理 |

### 使用例

#### カレントディレクトリの画像をPNGに変換

```bash
image-converter -ext png
```

#### 指定ディレクトリの画像をWebPに変換（品質90）

```bash
image-converter -src ./photos -ext webp -q 90
```

#### サブディレクトリも含めて全画像をJPEGに変換

```bash
image-converter -src ./photos -ext jpg -R
```

#### 出力先ディレクトリを指定して変換

```bash
image-converter -src ./photos -out ./converted -ext avif
```

## サポートされているフォーマット

### 入力フォーマット
- PNG
- JPEG
- WebP
- AVIF
- その他（usecasesパッケージのコーデックテーブルに依存）

### 出力フォーマット
- PNG
- JPEG (jpg)
- WebP (webp)
- AVIF (avif)

## ビルド方法

リポジトリのルートディレクトリで以下のコマンドを実行します：

```bash
./scripts/build_image_converter.sh
```

または手動でビルドする場合：

```bash
cd devbox
go build -o bin/image-converter ./cmd/cli/image-converter
