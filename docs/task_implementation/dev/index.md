# 開発用 Taskfile 実装フロー

対象: `$HOME/devbox/Taskfile.yml`

## 手順

1. タスク名を決める
  - 既存規則に合わせる（例: `build:zip:tools`, `test:all:cov`）
2. 実行内容を決める
  - 短い処理は `cmds` へ直接記述
  - 長い処理は `scripts/ops/*.sh` に切り出す
3. `Taskfile.yml` にタスクを追加する
  - `desc` と `cmds` を最低限設定する
  - 必要なら `vars` を追加する
4. テンプレート Taskfile を同期する
  - `internal/taskfile/usecases/taskfiles/root.yml` に同名タスクを追加する
  - `go run ./cmd/cli/taskfile/main.go --operation inspect --task-type root --taskfile-path $HOME/devbox/Taskfile.yml` で不足フィールドがないことを確認する
5. 動作確認する
  - `task -d $HOME/devbox --list-all`
  - `task -d $HOME/devbox <追加したタスク名>`
6. 関連更新を行う
  - 必要に応じて `scripts/ops` や関連ドキュメントを更新する

## 最小テンプレート

```yaml
tasks:
  build:example:
    desc: "Run example build step"
    cmds:
      - ./scripts/ops/example.sh
```
