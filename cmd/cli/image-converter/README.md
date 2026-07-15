# 画像変換ツール (Image Converter)

複数の画像ファイルを一括で別のフォーマットに変換するコマンドラインツールです。

## 機能

- 複数の画像ファイルを様々なフォーマット間で一括変換
- 非可逆圧縮フォーマットの品質調整

### オプション

| オプション | 必須/任意 | デフォルト値 | 説明 |
|---|---|---|---|
| `--src-dir` | 任意 | `.` | 変換元ディレクトリ |
| `--output-dir` | 必須 | なし | 出力先ディレクトリ |
| `--archive-dir` | 任意 | `""` (空) | アーカイブ先ディレクトリ (空の場合は無効) |
| `--move` | 任意 | `false` | 元ファイルをコピーではなく移動 (`--archive-dir` が指定されている場合のみ有効) |
| `--ext` | 任意 | `png` | 変換先フォーマット (png/jpg/webp/avif) |
| `--q` | 任意 | `80` | 非可逆圧縮フォーマットの品質 (1-100) |
| `--workers` | 任意 | CPU数 | 同時実行ワーカー数 |
| `--R` | 任意 | `false` | サブディレクトリを再帰的に処理 |
| `--lossless` | 任意 | `false` | ロスレス圧縮の有効化 |

## 使用例

カレントディレクトリの画像をPNGに変換
```bash
go run ./cmd/cli/image-converter --output-dir ./converted --ext png
```

指定ディレクトリの画像をWebPに変換 (品質90)
```bash
go run ./cmd/cli/image-converter --src-dir ./photos --output-dir ./converted --ext webp --q 90
```

出力先ディレクトリを指定して変換
```bash
go run ./cmd/cli/image-converter --src-dir ./photos --output-dir ./converted --ext avif
```

変換後に元ファイルをアーカイブディレクトリにコピー
```bash
go run ./cmd/cli/image-converter --src-dir ./photos --output-dir ./converted --ext webp --archive-dir ./originals
```

変換後に元ファイルをアーカイブディレクトリに移動
```bash
go run ./cmd/cli/image-converter --src-dir ./photos --output-dir ./converted --ext webp --archive-dir ./originals --move
```

## サポートされているフォーマット

### 入力フォーマット
- PNG
- JPEG
- WebP
- AVIF
- その他 (usecasesパッケージのコーデックテーブルに依存)

### 出力フォーマット
- PNG
- JPEG (jpg)
- WebP (webp)
- AVIF (avif)
