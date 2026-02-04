# スクリーンショットリネームツール (image-renamer-for-screenshot)

このツールは、VLCスナップショットファイル、Windowsスクリーンショットファイル、およびPixel端末で取得したスクリーンショット/録画ファイルを統一された命名規則でリネームするためのコマンドラインユーティリティです。

## 機能
- VLCスナップショットファイルを `Screenshot_YYYYMMDD-HHMMSS.png` 形式にリネーム
  - 対応形式1: `vlcsnap-YYYY-MM-DD-HH-MM-SS.png`
  - 対応形式2: `vlcsnap-YYYY-MM-DD-HHhMMmSSsNNN.png`（例: `vlcsnap-2025-05-06-23h59m44s239.png`）
- Windowsスクリーンショットファイル (`スクリーンショット YYYY-MM-DD HHMMSS.png`) を `Screenshot_YYYYMMDD-HHMMSS.png` 形式にリネーム
- Pixelスクリーンショット/録画ファイル (`screen-YYYYMMDD-HHMMSS.png/mp4`) を `Screenshot_YYYYMMDD-HHMMSS.png/mp4` 形式にリネーム
- Xiaomi端末のスクリーンショットファイル (`Screenshot_YYYY-MM-DD-HH-MM-SS-SSS_package.name.png/jpg` など) を `Screenshot_YYYYMMDD-HHMMSS.png/jpg` 形式にリネーム
- **新機能**: 全てのスクリーンショットファイルを `YYYYMMDDHHMMSS.png/mp4` 形式（日時のみ）にリネーム
- リネーム後のファイル名が既存ファイルと衝突する場合は処理開始前に検知して安全に中止
- 複数のファイルを並行処理
- 再帰的なディレクトリスキャン

## 使用方法
```bash
./image-renamer-for-screenshot [オプション]
```

### オプション

- `-src <ディレクトリ>`: スキャンするソースディレクトリ（デフォルト: カレントディレクトリ）
- `-operation <文字列>`: リネーム対象を指定します。指定できる値は `vlc`、`win`、`pixel`、`xiaomi` のいずれかです（例: `-operation=vlc`）
- `-to-datetime`: 全てのスクリーンショットファイルをYYYYMMDDHHMMSS形式にリネーム
- `-r`: サブディレクトリを再帰的にスキャン
- `-workers <数値>`: 並行ワーカー数（デフォルト: CPUコア数）

**注意**: `-operation` で指定できる値は一つだけです。複数の値を同時に指定することはできません。`-operation` を省略した場合は `-to-datetime` との組み合わせを除き、処理対象がないためエラーになります。値は `vlc` / `win` / `pixel` / `xiaomi` のいずれかを文字列で指定してください（例: `-operation=vlc`）。

### 使用例

VLCスナップショットファイルのリネーム
```bash
go run ./cmd/cli/image-renamer-for-screenshot -operation=vlc -src ./videos/screenshots
```

Windowsスクリーンショットファイルのリネーム
```bash
go run ./cmd/cli/image-renamer-for-screenshot -operation=win -src ./screenshots
```

Pixelスクリーンショット/録画ファイルのリネーム
```bash
go run ./cmd/cli/image-renamer-for-screenshot -operation=pixel -src ./pixel/media
```

Xiaomiスクリーンショットファイルのリネーム
```bash
go run ./cmd/cli/image-renamer-for-screenshot -operation=xiaomi -src ./miui/screenshots
```

再帰的にスキャン
```bash
go run ./cmd/cli/image-renamer-for-screenshot -operation=vlc -r -src ./media
```

全てのスクリーンショットファイルをYYYYMMDDHHMMSS形式にリネーム
```bash
go run ./cmd/cli/image-renamer-for-screenshot -to-datetime -src ./screenshots
```

このオプションは以下のファイル名パターンを自動検出してリネームします：
- `Screenshot_20250430-224038.png` → `20250430224038.png`
- `vlcsnap-2025-06-23-18h10m10s763.png` → `20250623181010.png`
- `スクリーンショット 2025-03-19 113905.png` → `20250319113905.png`
- `screen-20250601-235426.mp4` → `20250601235426.mp4`

## 対応ファイル形式

- `.jpg`
- `.jpeg`
- `.png`
- `.webp`
- `.avif`
- `.mp4`

## ビルド方法

リポジトリのルートディレクトリから以下のコマンドを実行します：

```bash
go build -o bin/image-renamer-for-screenshot ./cmd/cli/image-renamer-for-screenshot
```

または、提供されているビルドスクリプトを使用します：

```bash
./scripts/build_image_renamer_for_screenshot.sh
```
