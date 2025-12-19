# Image Renamer for Content

任意のディレクトリにある画像ファイルを、コンテンツIDと連番を基にした命名規則でリネームするCLIツールです。

## 主な機能

- 指定ディレクトリ内の画像（jpg, jpeg, png, webp, avif, bmp, gif）を検出
- ファイル名順または更新日時順にソートしてリネーム順を決定
- `contentId + delimiter + serial + _suffix + extension` 形式でファイル名を生成
- サブディレクトリの再帰処理に対応
- 並行処理により高速にリネーム

## 使い方

```
./image-renamer-for-content [オプション]
```

### オプション一覧

| オプション | デフォルト | 説明 |
|------------|------------|------|
| `-src` | `.` | スキャンするソースディレクトリ |
| `-name` | `false` | ファイル名順で並べ替え |
| `-time` | `false` | 更新日時順で並べ替え |
| `-operation` | (必須) | 実行モード。`mackerel` (MA・4桁), `web_clip` (WC・9桁), `date` (DA・5桁), `wine` (WI・4桁) に対応 |
| `-suffix` | `01` | 連番の後ろに付けるサフィックス（`_<suffix>` が付与されます） |
| `-delimiter` | `` | コンテンツIDと連番の間に挟む任意の区切り文字 |
| `-start` | `1` | 連番の開始番号 |
| `-r` | `false` | サブディレクトリを再帰的に処理 |
| `-workers` | CPU数 − 1 | 並行処理に用いるワーカー数 |

**注意:** `-name` か `-time` のいずれかを必ず指定してください。両方指定した場合は `-name` が優先されます。

## 使用例

```bash
# ファイル名順でリネーム (MA0001_01.jpg 形式)
go run ./cmd/cli/image-renamer-for-content -operation mackerel -name

# Webクリップ向けに9桁連番でリネーム
go run ./cmd/cli/image-renamer-for-content -operation web_clip -name

# 更新日時順で連番を付与し、サフィックスと区切り文字を変更
go run ./cmd/cli/image-renamer-for-content -src ./images -operation mackerel -time -suffix final -delimiter "-"

# 5番から連番を開始し、再帰的に処理
go run ./cmd/cli/image-renamer-for-content -operation date -name -start 5 -r

# ワーカー数を固定して実行
go run ./cmd/cli/image-renamer-for-content -operation mackerel -name -workers 4
```

## ビルド

付属スクリプトからビルドできます。

```bash
./scripts/build_image_renamer_for_content.sh
```

ビルド成果物は `pkg/bin/image-renamer-for-content/<platform>/` に生成されます。
