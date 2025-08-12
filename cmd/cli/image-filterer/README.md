# image-filterer

画像の指定領域にフィルター効果を適用するためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。

## 機能

- 指定したディレクトリ内の画像ファイル（PNG、JPG、JPEG、WebP）の特定領域にフィルターを適用
- サポートしているフィルター：
  - ぼかし（Gaussianブラー）
  - グレースケール（カスタムRGB重み対応）
- 再帰的なサブディレクトリの処理
- 並行処理による高速な処理
- 元ファイルのアーカイブ機能

## インストール

```bash
go install github.com/landmaster135/devbox/cmd/cli/image-filterer@latest
```

## 使用方法

```bash
image-filterer [オプション]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-src` | `.` | 処理元ディレクトリ |
| `-out` | srcと同じ | 出力先ディレクトリ |
| `-arc` | `./5_original_files` | アーカイブ先ディレクトリ |
| `-x1` | `0` | フィルター適用領域の左上X座標（全て0の場合は画像全体を処理） |
| `-y1` | `0` | フィルター適用領域の左上Y座標（全て0の場合は画像全体を処理） |
| `-x2` | `0` | フィルター適用領域の右下X座標（全て0の場合は画像全体を処理） |
| `-y2` | `0` | フィルター適用領域の右下Y座標（全て0の場合は画像全体を処理） |
| `-mode` | `blur` | フィルターモード（`blur` または `grayscale`） |
| `-radius` | `10.0` | ぼかしの半径（値が大きいほど強いぼかし効果、blurモード時のみ有効） |
| `-r-weight` | `0.3` | グレースケール変換時の赤チャンネルの重み（0.0-1.0、grayscaleモード時のみ有効） |
| `-g-weight` | `0.6` | グレースケール変換時の緑チャンネルの重み（0.0-1.0、grayscaleモード時のみ有効） |
| `-b-weight` | `0.1` | グレースケール変換時の青チャンネルの重み（0.0-1.0、grayscaleモード時のみ有効） |
| `-suffix` | `filtered` | 出力ファイル名に付加するサフィックス |
| `-move` | `false` | 元ファイルを移動する（コピーではなく） |
| `-r` | `false` | サブディレクトリを再帰的に処理 |
| `-workers` | CPU数 | 同時実行ワーカー数 |

## 使用例

画像全体をぼかし処理（座標指定なし）：
```bash
image-filterer -mode blur
```

画像全体をグレースケール変換：
```bash
image-filterer -mode grayscale
```

カレントディレクトリ内の画像の座標(10,20)から(300,400)までの領域をぼかし処理：
```bash
image-filterer -x1 10 -y1 20 -x2 300 -y2 400
```

特定のディレクトリの画像を処理：
```bash
image-filterer -src ./images -x1 10 -y1 20 -x2 300 -y2 400
```

出力先を指定：
```bash
image-filterer -src ./images -out ./filtered_images -x1 10 -y1 20 -x2 300 -y2 400
```

ぼかしの強さを調整：
```bash
image-filterer -src ./images -x1 10 -y1 20 -x2 300 -y2 400 -radius 5.0
```

サブディレクトリも含めて処理：
```bash
image-filterer -src ./images -r -x1 10 -y1 20 -x2 300 -y2 400
```

元ファイルを移動：
```bash
image-filterer -src ./images -move -x1 10 -y1 20 -x2 300 -y2 400
```

カスタムサフィックスを指定：
```bash
image-filterer -src ./images -suffix blurred -x1 10 -y1 20 -x2 300 -y2 400
```

グレースケール変換（デフォルト重み）：
```bash
image-filterer -src ./images -mode grayscale -x1 10 -y1 20 -x2 300 -y2 400
```

グレースケール変換（カスタム重み）- 青チャンネルを強調：
```bash
image-filterer -src ./images -mode grayscale -r-weight 0.2 -g-weight 0.3 -b-weight 0.5 -x1 10 -y1 20 -x2 300 -y2 400
```

グレースケール変換（セピア調）：
```bash
image-filterer -src ./images -mode grayscale -r-weight 0.4 -g-weight 0.4 -b-weight 0.2 -x1 10 -y1 20 -x2 300 -y2 400
```

## 注意事項

- フィルター適用領域の座標は必ず `x2 > x1` かつ `y2 > y1` となるように指定してください。
- 座標を指定しない場合（全て0のデフォルト値）は、画像全体にフィルターが適用されます。
- サポートしているファイル形式は PNG、JPG、JPEG、WebP です。
- 大量のファイルを処理する場合は、`-workers` オプションでワーカー数を調整することで処理速度を最適化できます。
- ぼかし効果の強さは `-radius` オプションで調整できます。値が大きいほど強いぼかし効果になります。
