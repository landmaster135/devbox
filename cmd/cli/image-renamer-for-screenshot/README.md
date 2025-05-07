# スクリーンショットリネームツール (image-renamer-for-screenshot)

このツールは、VLCスナップショットファイルとWindowsスクリーンショットファイルを統一された命名規則でリネームするためのコマンドラインユーティリティです。

## 機能

- VLCスナップショットファイル (`vlcsnap-YYYY-MM-DD-HH-MM-SS.png`) を `Screenshot_YYYYMMDD-HHMMSS.png` 形式にリネーム
- Windowsスクリーンショットファイル (`スクリーンショット YYYY-MM-DD HHMMSS.png`) を `Screenshot_YYYYMMDD-HHMMSS.png` 形式にリネーム
- 複数のファイルを並行処理
- 再帰的なディレクトリスキャン

## 使用方法

```bash
./image-renamer-for-screenshot [オプション]
```

### オプション

- `-src <ディレクトリ>`: スキャンするソースディレクトリ（デフォルト: カレントディレクトリ）
- `-vlc`: VLCスナップショットファイル (vlcsnap-*.png) をリネーム
- `-win`: Windowsスクリーンショットファイル (スクリーンショット *.png) をリネーム
- `-r`: サブディレクトリを再帰的にスキャン
- `-workers <数値>`: 並行ワーカー数（デフォルト: CPUコア数）

**注意**: `-vlc` と `-win` のフラグは同時に設定できません。どちらか一方のみを指定してください。

### 使用例

#### VLCスナップショットファイルのリネーム

```bash
./image-renamer-for-screenshot -vlc -src ./videos/screenshots
```

#### Windowsスクリーンショットファイルのリネーム

```bash
./image-renamer-for-screenshot -win -src ./screenshots
```

#### 再帰的にスキャン

```bash
./image-renamer-for-screenshot -vlc -r -src ./media
```

## 対応ファイル形式

- `.jpg`
- `.jpeg`
- `.png`
- `.webp`
- `.avif`

## ビルド方法

リポジトリのルートディレクトリから以下のコマンドを実行します：

```bash
go build -o bin/image-renamer-for-screenshot ./cmd/cli/image-renamer-for-screenshot
```

または、提供されているビルドスクリプトを使用します：

```bash
./scripts/build_image_renamer_for_screenshot.sh
```
