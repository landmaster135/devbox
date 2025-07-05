# file-maneuver

複数のディレクトリから指定した拡張子のファイルを検索し、単一の宛先ディレクトリに移動またはコピーするCLIツールです。

## 機能

- **複数ソースディレクトリ対応**: カンマ区切りで複数のソースディレクトリを指定可能
- **複数拡張子対応**: カンマ区切りで複数の拡張子を指定可能
- **移動・コピー選択**: ファイルの移動またはコピーを選択可能
- **上書き制御**: 同名ファイルの上書きを許可または禁止
- **再帰的検索**: サブディレクトリも含めた検索が可能
- **並行処理**: マルチワーカーによる高速なファイル処理
- **ドライランモード**: 実際の処理を行わずに動作確認が可能
- **衝突検出**: 同名ファイルの衝突を検出し、安全に処理
- **詳細ログ**: 処理状況の詳細な出力

## 使用方法

### 基本的な使用例

```bash
# 単一ディレクトリから画像ファイルを移動
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg,png,gif" \
  --dest-dir "/path/to/destination"

# 複数ディレクトリから複数拡張子のファイルを移動
./file-maneuver \
  --src-dirs "/path/to/source1,/path/to/source2,/path/to/source3" \
  --extensions "pdf,doc,docx,txt" \
  --dest-dir "/path/to/documents"

# 再帰的検索を有効にして移動
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "mp4,avi,mkv" \
  --dest-dir "/path/to/videos" \
  --recursive

# ドライランモードで動作確認
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg,png" \
  --dest-dir "/path/to/destination" \
  --dry-run

# ワーカー数を指定して高速処理
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg,png" \
  --dest-dir "/path/to/destination" \
  --workers 8

# ファイルをコピー（移動ではなく）
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg,png" \
  --dest-dir "/path/to/destination" \
  --copy

# コピーモードでドライラン
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "pdf,doc" \
  --dest-dir "/path/to/backup" \
  --copy \
  --dry-run

# 上書きモードで移動
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg,png" \
  --dest-dir "/path/to/destination" \
  --overwrite

# 上書きモードでコピー
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "pdf,doc" \
  --dest-dir "/path/to/backup" \
  --copy \
  --overwrite

# 上書きモードでドライラン
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg,png" \
  --dest-dir "/path/to/destination" \
  --overwrite \
  --dry-run
```

## コマンドラインオプション

| オプション | 必須 | 説明 | デフォルト値 |
|-----------|------|------|-------------|
| `--src-dirs` | ✅ | ソースディレクトリ（カンマ区切りで複数指定可能） | - |
| `--extensions` | ✅ | 対象拡張子（カンマ区切りで複数指定可能） | - |
| `--dest-dir` | ✅ | 宛先ディレクトリ（単一のみ） | - |
| `--copy` | ❌ | ファイルを移動ではなくコピーする | `false` |
| `--overwrite` | ❌ | 宛先に同名ファイルが存在する場合に上書きする | `false` |
| `--recursive` | ❌ | 再帰的にサブディレクトリを検索 | `false` |
| `--workers` | ❌ | 並行処理のワーカー数 | CPU数 |
| `--dry-run` | ❌ | 実際の処理を行わずに動作確認のみ | `false` |

## 実行例

### 例1: 写真の整理

```bash
./file-maneuver \
  --src-dirs "/Users/user/Downloads,/Users/user/Desktop" \
  --extensions "jpg,jpeg,png,heic" \
  --dest-dir "/Users/user/Pictures/Organized" \
  --recursive
```

**出力例:**
```
✨ ファイル移動処理を開始します
  ソースディレクトリ: /Users/user/Downloads, /Users/user/Desktop
  対象拡張子: jpg, jpeg, png, heic
  宛先ディレクトリ: /Users/user/Pictures/Organized
  再帰的検索: true
  ワーカー数: 8

🔍 対象ファイルを検索中...
ディレクトリ /Users/user/Downloads から 15 ファイルを発見しました
ディレクトリ /Users/user/Desktop から 3 ファイルを発見しました
合計 18 ファイルが見つかりました

📦 ファイル移動を開始します...
18 ファイルを移動します（0 ファイルをスキップ）
ファイル移動に 8 ワーカーを使用します
移動中: /Users/user/Downloads/IMG_001.jpg -> /Users/user/Pictures/Organized/IMG_001.jpg
移動中: /Users/user/Downloads/IMG_002.png -> /Users/user/Pictures/Organized/IMG_002.png
...

✅ ファイル移動処理が完了しました
  成功: 18 ファイル
```

### 例2: ドキュメントの整理（ドライランモード）

```bash
./file-maneuver \
  --src-dirs "/Users/user/Downloads" \
  --extensions "pdf,doc,docx,txt" \
  --dest-dir "/Users/user/Documents" \
  --dry-run
```

**出力例:**
```
✨ ファイル移動処理を開始します
  ソースディレクトリ: /Users/user/Downloads
  対象拡張子: pdf, doc, docx, txt
  宛先ディレクトリ: /Users/user/Documents
  再帰的検索: false
  ワーカー数: 8
  モード: ドライラン（実際の移動は行いません）

🔍 対象ファイルを検索中...
ディレクトリ /Users/user/Downloads から 5 ファイルを発見しました
合計 5 ファイルが見つかりました

📦 ファイル移動を開始します...
5 ファイルを移動します（0 ファイルをスキップ）
ドライランモード: 実際の移動は行いません
移動予定: /Users/user/Downloads/report.pdf -> /Users/user/Documents/report.pdf
移動予定: /Users/user/Downloads/memo.txt -> /Users/user/Documents/memo.txt
...

✅ ファイル移動処理が完了しました
  成功: 5 ファイル
```

### 例3: ファイルのコピー（バックアップ作成）

```bash
./file-maneuver \
  --src-dirs "/Users/user/Documents/Important" \
  --extensions "pdf,doc,docx,txt" \
  --dest-dir "/Users/user/Backup/Documents" \
  --copy \
  --recursive
```

**出力例:**
```
✨ ファイルコピー処理を開始します
  ソースディレクトリ: /Users/user/Documents/Important
  対象拡張子: pdf, doc, docx, txt
  宛先ディレクトリ: /Users/user/Backup/Documents
  再帰的検索: true
  ワーカー数: 8
  モード: コピー

🔍 対象ファイルを検索中...
ディレクトリ /Users/user/Documents/Important から 12 ファイルを発見しました
合計 12 ファイルが見つかりました

📦 ファイルコピーを開始します...
12 ファイルをコピーします（0 ファイルをスキップ）
ファイルコピーに 8 ワーカーを使用します
コピー中: /Users/user/Documents/Important/report.pdf -> /Users/user/Backup/Documents/report.pdf
コピー中: /Users/user/Documents/Important/memo.txt -> /Users/user/Backup/Documents/memo.txt
...

✅ ファイルコピー処理が完了しました
  成功: 12 ファイル
```

### 例4: ファイル衝突の処理

```bash
./file-maneuver \
  --src-dirs "/path/to/source" \
  --extensions "jpg" \
  --dest-dir "/path/to/destination"
```

**出力例（衝突がある場合）:**
```
✨ ファイル移動処理を開始します
...

🔍 対象ファイルを検索中...
合計 3 ファイルが見つかりました

📦 ファイル移動を開始します...
警告: 以下のファイルは宛先に同名ファイルが存在するためスキップされます:
  /path/to/source/photo1.jpg

2 ファイルを移動します（1 ファイルをスキップ）
...

✅ ファイル移動処理が完了しました
  成功: 2 ファイル
```

## 注意事項

### 移動とコピーの違い

- **移動モード（デフォルト）**: ファイルが宛先に移動され、元の場所からは削除されます
- **コピーモード（`--copy`）**: ファイルが宛先にコピーされ、元の場所にも残ります
- コピーモードでは、ファイルの権限とタイムスタンプが保持されます

### ファイル衝突と上書きについて

- **通常モード（デフォルト）**: 宛先ディレクトリに同名ファイルが存在する場合、そのファイルはスキップされます
- **上書きモード（`--overwrite`）**: 宛先ディレクトリに同名ファイルが存在する場合、そのファイルを上書きします
- 通常モードでは元ファイルは削除されず、そのまま残ります
- 衝突したファイルは警告メッセージで表示されます（通常モードのみ）
- 上書きモードでは衝突チェックをスキップし、すべてのファイルを処理対象とします

### 拡張子の処理

- 拡張子は大文字・小文字を区別しません（`JPG` と `jpg` は同じ扱い）
- ドット（`.`）の有無は自動で正規化されます（`jpg` → `.jpg`）

### パフォーマンス

- デフォルトのワーカー数はCPU数と同じです
- 大量のファイルを処理する場合は、`--workers` オプションで調整可能
- 最大ワーカー数は `CPU数 × 2` に制限されます

### 安全性

- ドライランモードで事前に動作確認することを推奨
- 重要なファイルは事前にバックアップを取ることを推奨
- 書き込み権限のないディレクトリは事前にエラーで検出されます

## ビルド方法

```bash
cd /path/to/devbox/cmd/cli/file-maneuver
go build -o file-maneuver .
```

## テスト実行

```bash
cd /path/to/devbox
go test ./internal/file_maneuver/usecases/... -v
```

### カバレッジ確認

```bash
cd /path/to/devbox
go test -coverprofile=coverage.out ./internal/file_maneuver/usecases/...
go tool cover -html=coverage.out -o coverage.html
```

## アーキテクチャ

### ディレクトリ構造

```
devbox/
├── cmd/cli/file-maneuver/
│   ├── main.go              # CLI層（フラグ処理、UI）
│   └── README.md            # このファイル
└── internal/file_maneuver/usecases/
    ├── services.go          # サービス層（ビジネスロジック）
    └── services_test.go     # テストコード
```

### 設計原則

- **関心の分離**: CLI層とビジネスロジック層を明確に分離
- **単一責任**: 各構造体・メソッドは単一の責任を持つ
- **早期バリデーション**: 構造体作成時に全ての検証を完了
- **エラーハンドリング**: 包括的なエラー処理と適切なメッセージ
- **テスタビリティ**: TDD原則に基づく高いテストカバレッジ

### 主要コンポーネント

#### Config構造体
- 設定の保持とバリデーション
- 構造体作成時に全ての検証を実行
- 不正な設定での作成を防止

#### FileManeuverService構造体
- ファイル検索・移動の主要ロジック
- 並行処理による高速化
- 衝突検出と安全な処理

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのライセンスファイルを参照してください。
