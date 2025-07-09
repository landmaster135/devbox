# image-rotator

画像を回転させるためのコマンドラインツールです。複数の画像ファイルを一括処理することができます。

## 機能

- 指定したディレクトリ内の画像ファイル（PNG、JPG、JPEG）を回転
- 再帰的なサブディレクトリの処理
- 並行処理による高速な処理
- 元ファイルのアーカイブ機能

## インストール

```bash
go install github.com/landmaster135/devbox/cmd/cli/image-rotator@latest
```

## 使用方法

```bash
image-rotator [オプション]
```

### オプション

| オプション | デフォルト値 | 説明 |
|------------|--------------|------|
| `-src` | `.` | 回転元ディレクトリ |
| `-out` | srcと同じ | 出力先ディレクトリ |
| `-arc` | `./5_original_files` | アーカイブ先ディレクトリ |
| `-angle` | - | 回転角度（度数法、時計回り、必須） |
| `-suffix` | `rotated` | 出力ファイル名に付加するサフィックス |
| `-move` | `false` | 元ファイルを移動する（コピーではなく） |
| `-r` | `false` | サブディレクトリを再帰的に処理 |
| `-workers` | CPU数 | 同時実行ワーカー数 |

## 使用例

### 基本的な使用方法

カレントディレクトリ内の画像を90度回転：

```bash
image-rotator -angle 90
```

### 特定のディレクトリの画像を処理

```bash
image-rotator -src ./images -angle 90
```

### 出力先を指定

```bash
image-rotator -src ./images -out ./rotated_images -angle 90
```

### サブディレクトリも含めて処理

```bash
image-rotator -src ./images -r -angle 90
```

### 元ファイルを移動

```bash
image-rotator -src ./images -move -angle 90
```

### カスタムサフィックスを指定

```bash
image-rotator -src ./images -suffix turned -angle 90
```

### 反時計回りに回転

```bash
image-rotator -src ./images -angle -90
```

## 注意事項

- 回転角度は必須パラメータで、0以外の値を指定してください。
- サポートしているファイル形式は PNG、JPG、JPEG のみです。
- 大量のファイルを処理する場合は、`-workers` オプションでワーカー数を調整することで処理速度を最適化できます。
