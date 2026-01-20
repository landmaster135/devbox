# Taskfile CLI

Taskfile を検証するための CLI ツールです。指定した Taskfile が、プロジェクト標準 (`internal/taskfile/usecases/taskfiles/root.yml`) に含まれるフィールドを欠けなく定義しているかを確認します。

## 主な機能
- **inspect**: 参照 Taskfile と照合し、存在しないフィールドを一覧化
- **root Taskfile チェック**: プロジェクトルートの Taskfile を対象に、必要なタスク・desc・cmds が定義されているかを検証

## インストール

```bash
# プロジェクトルートから
go build -o bin/taskfile ./cmd/cli/taskfile
```

## 使用例

### 基本コマンド

```bash
go run ./cmd/cli/taskfile \
  --operation inspect \
  --task-type root \
  --taskfile-path ./Taskfile.yml
```

### オプション

| オプション | 必須 | 説明 | 例 |
|------------|------|------|----|
| `--operation` | * | 実行する操作。現在は `inspect` のみ対応 | `--operation inspect` |
| `--task-type` | * | 参照 Taskfile の種類。現在は `root` のみ対応 | `--task-type root` |
| `--taskfile-path` | * | 検証対象となる Taskfile のパス | `--taskfile-path ./Taskfile.yml` |
| `--help` / `-h` | | ヘルプを表示 | `--help` |

## 出力例

### すべてのフィールドがある場合
```
Taskfileには参照Taskfileのすべてのフィールドが含まれています。
```

### フィールドが不足している場合
```
不足しているフィールドが 3 個見つかりました:
  - tasks.alias
  - tasks.alias.cmds
  - tasks.alias.desc
```

## ワークフロー例
1. Taskfile を編集・生成
2. `go run ./cmd/cli/taskfile --operation inspect --task-type root --taskfile-path ./Taskfile.yml`
3. 不足フィールドが表示されたら追記して再実行

参照 Taskfile は `internal/taskfile/usecases/taskfiles/root.yml` に保守されており、必要なフィールドを更新する際はこのファイルを起点にしてください。
