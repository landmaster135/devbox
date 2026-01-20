# Taskfile CLI

Taskfile を検証・補完・新規作成するための CLI ツールです。指定した Taskfile が、プロジェクト標準 (`internal/taskfile/usecases/taskfiles/root.yml`) に含まれるフィールドを欠けなく定義しているかを確認し、空欄フィールドをテンプレートの値で埋めたり、テンプレートそのものを任意のパスへ複製して新規作成できます。

## 主な機能
- **inspect**: 参照 Taskfile と照合し、存在しないフィールドを一覧化
- **fill**: 空欄または未定義のフィールドを参照 Taskfile の値で補完
- **new**: `internal/taskfile/usecases/taskfiles/root.yml` をベースに `--taskfile-path` で指定した場所へ Taskfile を新規作成
- **root Taskfile チェック**: プロジェクトルートの Taskfile を対象に、必要なタスク・desc・cmds が定義されているかを検証

## インストール

```bash
# プロジェクトルートから
go build -o bin/taskfile ./cmd/cli/taskfile
```

## 使用例

### テンプレートで構成を検閲

```bash
go run ./cmd/cli/taskfile \
  --operation inspect \
  --task-type root \
  --taskfile-path ./Taskfile.yml
```

### テンプレートで空欄を補完

```bash
go run ./cmd/cli/taskfile \
  --operation fill \
  --task-type root \
  --taskfile-path ./Taskfile.yml
```

### テンプレートから新規作成

```bash
go run ./cmd/cli/taskfile \
  --operation new \
  --task-type root \
  --taskfile-path ./Taskfile.yml
```

### オプション

| オプション | 必須 | 説明 | 例 |
|------------|------|------|----|
| `--operation` | * | 実行する操作。`inspect` / `fill` / `new` | `--operation new` |
| `--task-type` | * | 参照 Taskfile の種類。現在は `root` のみ対応 | `--task-type root` |
| `--taskfile-path` | * | 検証/補完/新規作成対象となる Taskfile のパス | `--taskfile-path ./Taskfile.yml` |
| `--help` / `-h` | | ヘルプを表示 | `--help` |

## 出力例

### すべてのフィールドがある場合 (inspect)
```
Taskfileには参照Taskfileのすべてのフィールドが含まれています。
```

### フィールドが不足している場合 (inspect)
```
不足しているフィールドが 3 個見つかりました:
  - tasks.alias
  - tasks.alias.cmds
  - tasks.alias.desc
```

### 空欄フィールドが埋められた場合 (fill)
```
Taskfileの空欄フィールドをテンプレートの値で補完しました。
```

### 補完対象がない場合 (fill)
```
補完対象の空欄フィールドは見つかりませんでした。
```

### 新規作成が完了した場合 (new)
```
Taskfileを新規作成しました: ./Taskfile.yml
```

## ワークフロー例
1. `--operation new` でテンプレート Taskfile を任意のパスへ作成
2. 必要に応じて Taskfile を編集
3. `go run ./cmd/cli/taskfile --operation inspect --task-type root --taskfile-path ./Taskfile.yml`
4. 不足フィールドが表示されたら追記して再実行
