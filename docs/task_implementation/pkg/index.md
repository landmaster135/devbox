# 配布用 Taskfile 実装フロー

対象:
- `$HOME/devbox/pkg/taskfile/Taskfile.yml`
- `$HOME/devbox/pkg/taskfile/taskfiles/*.yml`

## 手順

1. 追加先を決める
  - 実処理タスクは `taskfiles/*.yml` に追加する
  - 呼び出し窓口タスクは `pkg/taskfile/Taskfile.yml` に追加する
2. 分割 Taskfile を更新する
  - 対応する `taskfiles/*.yml` が存在するか確認して、対応するファイル内の既存規則に合わせてタスク定義を追記する
  - 対応先がない場合は新しいカテゴリファイルを `taskfiles/` 配下に作成してタスク定義を追記する
3. 呼び出し窓口タスク `pkg/taskfile/Taskfile.yml` を更新する
  - 新規カテゴリの場合は `includes` 内に追加する
  - 必要に応じて呼び出し窓口タスクを `tasks` に追加する
4. 動作確認する
  - `task -t $HOME/devbox/pkg/taskfile/Taskfile.yml --list-all`
  - `task -t $HOME/devbox/pkg/taskfile/Taskfile.yml <追加したタスク名>`
5. 関連更新を行う
  - 必要に応じて `pkg/taskfile/README.md` を更新する

## 最小テンプレート

```yaml
tasks:
  image-convert:image-converter:to-jpg:
    desc: "Convert PNG files to JPG"
    cmds:
      - ./pkg/bin/cli/linux_amd64/image-converter --operation=to-jpg
```
