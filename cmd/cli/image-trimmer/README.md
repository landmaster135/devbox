# image-trimmer

画像をトリミングするためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。

## 機能

- 指定したディレクトリ内の画像ファイル（PNG、JPG、JPEG）をトリミング

### オプション

| オプション | 必須/任意 | デフォルト値 | 説明 |
|---|---|---|---|
| `-src-dir` | 任意 | `.` | トリミング元ディレクトリ |
| `-output-dir` | 必須 | - | 出力先ディレクトリ |
| `-archive-dir` | 任意 | `./5_original_files` | アーカイブ先ディレクトリ |
| `-x1` | 任意 | `0` | トリミングを開始するX座標 |
| `-y1` | 任意 | `0` | トリミングを開始するY座標 |
| `-x2` | 必須 | - | トリミングを終了するX座標 |
| `-y2` | 必須 | - | トリミングを終了するY座標 |
| `-suffix` | 任意 | `trimmed` | 出力ファイル名に付加するサフィックス |
| `-move` | 任意 | `false` | 元ファイルを移動する（コピーではなく） |
| `-r` | 任意 | `false` | サブディレクトリを再帰的に処理 |
| `-workers` | 任意 | CPU数 | 同時実行ワーカー数 |

## 使用例

カレントディレクトリ内の画像を座標 (10,20) から (300,400) までトリミング：
```bash
go run ./cmd/cli/image-trimmer -output-dir ./trimmed_images -x1 10 -y1 20 -x2 300 -y2 400
```

特定のディレクトリの画像を処理
```bash
go run ./cmd/cli/image-trimmer -src-dir ./images -output-dir ./trimmed_images -x1 10 -y1 20 -x2 300 -y2 400
```

出力先を指定
```bash
go run ./cmd/cli/image-trimmer -src-dir ./images -output-dir ./trimmed_images -x1 10 -y1 20 -x2 300 -y2 400
```

サブディレクトリも含めて処理
```bash
go run ./cmd/cli/image-trimmer -src-dir ./images -output-dir ./trimmed_images -r -x1 10 -y1 20 -x2 300 -y2 400
```

元ファイルを移動
```bash
go run ./cmd/cli/image-trimmer -src-dir ./images -output-dir ./trimmed_images -move -x1 10 -y1 20 -x2 300 -y2 400
```

カスタムサフィックスを指定
```bash
go run ./cmd/cli/image-trimmer -src-dir ./images -output-dir ./trimmed_images -suffix cropped -x1 10 -y1 20 -x2 300 -y2 400
```

## 注意事項

- トリミング座標は必ず `x2 > x1` かつ `y2 > y1` となるように指定してください。
- `-output-dir` は必須です。出力先ディレクトリを明示してください。
- サポートしているファイル形式は PNG、JPG、JPEG のみです。
- 大量のファイルを処理する場合は、`-workers` オプションでワーカー数を調整することで処理速度を最適化できます。
