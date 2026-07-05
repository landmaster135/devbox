# 配布用 Taskfile 実装フロー

対象:
- `$HOME/devbox/pkg/taskfile/Taskfile.yml`
- `$HOME/devbox/pkg/taskfile/taskfiles/*.yml`

## 手順

1. 追加先を決める
  - 実処理タスクは `taskfiles/*.yml` に追加する
  - 呼び出し窓口タスクは `pkg/taskfile/Taskfile.yml` に追加する
  - 依頼文のタスク名をそのまま採用せず、追加先ごとの既存名前空間に合わせる
2. 分割 Taskfile を更新する
  - 対応する `taskfiles/*.yml` が存在するか確認して、対応するファイル内の既存規則に合わせてタスク定義を追記する
  - 対応先がない場合は新しいカテゴリファイルを `taskfiles/` 配下に作成してタスク定義を追記する
3. 呼び出し窓口タスク `pkg/taskfile/Taskfile.yml` を更新する
  - 新規カテゴリの場合は `includes` 内に追加する
  - 必要に応じて呼び出し窓口タスクを `tasks` に追加する
  - 既存の窓口タスク名を周辺行で確認し、同じカテゴリの名前空間に揃える
4. 動作確認する
  - `task -t $HOME/devbox/pkg/taskfile/Taskfile.yml --list-all`
  - `task -t $HOME/devbox/pkg/taskfile/Taskfile.yml <追加したタスク名>`
5. 関連更新を行う
  - 必要に応じて `pkg/taskfile/README.md` を更新する

## 追加前後のチェック
### 追加前のチェック
`pkg/taskfile/Taskfile.yml` を更新する前に、追加位置の周辺タスクを確認して、窓口タスクの名前空間を確定する。
```bash
sed -n '<追加位置の前後行>p' pkg/taskfile/Taskfile.yml
```

例: 画像変換の窓口タスクが `image:convert:*` で並んでいる場合、追加する窓口タスクも `image:convert:*` にする。

そして、 `taskfiles/*.yml` を更新する前に、対象ファイル内の同カテゴリタスクを確認して、実処理タスクの名前空間を確定する。
```bash
sed -n '<追加位置の前後行>p' pkg/taskfile/taskfiles/<target>.yml
```

例: `pkg/taskfile/taskfiles/image_convert.yml` 内の実処理タスクが `image-converter:*` で並んでいる場合、追加する実処理タスクも `image-converter:*` にする。

### 追加後のチェック
追加後は、窓口タスク、実処理タスク、`cmds[].task` の名前空間がそれぞれ意図どおりか検索する。
```bash
rg -n '<追加したタスク名の主要部分>' pkg/taskfile/Taskfile.yml pkg/taskfile/taskfiles/<target>.yml
```

例:
- `pkg/taskfile/Taskfile.yml`: `image:convert:jpg-to-webp:cwebp`
- `pkg/taskfile/taskfiles/image_convert.yml`: `image-converter:jpg-to-webp:cwebp`
- `cmds[].task`: `image-convert:image-converter:jpg-to-webp:cwebp`

## 最小テンプレート

```yaml
tasks:
  image-convert:image-converter:to-jpg:
    desc: "Convert PNG files to JPG"
    cmds:
      - ./pkg/bin/cli/linux_amd64/image-converter --operation=to-jpg
```
