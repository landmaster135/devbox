# image-trimmer

画像をトリミングするためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。

## 機能

- 指定したディレクトリ内の画像ファイル（PNG、JPG、JPEG）をトリミング
- 再帰的なサブディレクトリの処理
- 並行処理による高速な処理
- 元ファイルのアーカイブ機能

## インストール

```bash
go install github.com/landmaster135/devbox/cmd/cli/image-trimmer@latest
```

## 使用方法

```bash
image-trimmer [オプション]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-src` | `.` | トリミング元ディレクトリ |
| `-out` | srcと同じ | 出力先ディレクトリ |
| `-arc` | `./5_original_files` | アーカイブ先ディレクトリ |
| `-x1` | `0` | トリミングを開始するX座標 |
| `-y1` | `0` | トリミングを開始するY座標 |
| `-x2` | - | トリミングを終了するX座標（必須） |
| `-y2` | - | トリミングを終了するY座標（必須） |
| `-suffix` | `trimmed` | 出力ファイル名に付加するサフィックス |
| `-move` | `false` | 元ファイルを移動する（コピーではなく） |
| `-r` | `false` | サブディレクトリを再帰的に処理 |
| `-workers` | CPU数 | 同時実行ワーカー数 |

## 使用例

### 基本的な使用方法

カレントディレクトリ内の画像を座標(10,20)から(300,400)までトリミング：

```bash
go run ./cmd/cli/image-trimmer -x1 10 -y1 20 -x2 300 -y2 400
```

### 特定のディレクトリの画像を処理

```bash
go run ./cmd/cli/image-trimmer -src ./images -x1 10 -y1 20 -x2 300 -y2 400
```

### 出力先を指定

```bash
go run ./cmd/cli/image-trimmer -src ./images -out ./trimmed_images -x1 10 -y1 20 -x2 300 -y2 400
```

### サブディレクトリも含めて処理

```bash
go run ./cmd/cli/image-trimmer -src ./images -r -x1 10 -y1 20 -x2 300 -y2 400
```

### 元ファイルを移動

```bash
go run ./cmd/cli/image-trimmer -src ./images -move -x1 10 -y1 20 -x2 300 -y2 400
```

### カスタムサフィックスを指定

```bash
go run ./cmd/cli/image-trimmer -src ./images -suffix cropped -x1 10 -y1 20 -x2 300 -y2 400
```

## 注意事項

- トリミング座標は必ず `x2 > x1` かつ `y2 > y1` となるように指定してください。
- サポートしているファイル形式は PNG、JPG、JPEG のみです。
- 大量のファイルを処理する場合は、`-workers` オプションでワーカー数を調整することで処理速度を最適化できます。
