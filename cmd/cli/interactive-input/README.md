# Interactive Input

ユーザーに1つの質問を投げかけ、標準出力へ機械的に扱いやすいキー付きの値を返すCLIツールです。Windowsバッチの `set /p` や `choice` の代替として設計されており、案内文やリトライ通知は標準エラーに集約されます。

## 主な特徴
- **3種類の入力モード**: `text`(任意文字列)、`choice`(ショートカット選択)、`confirm`(Y/N)
- **キー付き出力**: text/choiceは `--<key>=<value>`、confirmは肯定時のみ `--<key>` を出力
- **リトライ制御**: `--max-attempts` で再入力回数を制御（0で無制限）
- **複数選択肢**: `--choice-option shortcut|output` を複数指定可能

## ビルド

```bash
./scripts/build/build_interactive_input.sh
# もしくは全CLI一括: ./scripts/build.sh
```

ビルド後のバイナリは `pkg/bin/cli/<platform>/interactive-input` に配置されます。

## 共通フラグ

| フラグ | 必須 | 説明 |
| ------ | ---- | ---- |
| `--prompt` | * | 質問文。`\n` で改行可 |
| `--input-type` | * | `text` / `choice` / `confirm` |
| `--key` | * | 出力に使うキー（英数字・`-`・`_`） |
| `--default` | text | Enterのみ時に採用する既定値（空文字可） |
| `--choice-option` | choice | `shortcut|output` 形式。複数指定で選択肢追加 |
| `--max-attempts` | 任意 | バリデーション失敗時の再入力回数。0で無制限（既定3） |
| `--help`, `-h` | 任意 | ヘルプを表示 |

### 出力形式
- text / choice: `--<key>=<value>`
- confirm: YES → `--<key>` / NO → 空文字（標準出力なし）

## 使用例

### テキスト入力
```bash
interactive-input \
  --prompt "Input start number: " \
  --input-type text \
  --key start \
  --max-attempts 0
```
ユーザーが `123` を入力すると `--start=123` が標準出力へ返ります。Enterのみの場合は再入力を促します。`--default` を付ければ Enter だけで既定値が採用されます。

### 選択肢入力
```bash
interactive-input \
  --prompt "Select rename mode [v/w/x/a]: " \
  --input-type choice \
  --key operation \
  --choice-option "v|vlc" \
  --choice-option "w|win" \
  --choice-option "x|xiaomi" \
  --choice-option "a|pixel" \
  --max-attempts 0
```
`v` を入力すると `--operation=vlc` が返ります。ショートカットは大文字/小文字を区別せず、未定義値はリトライ対象になります。

### 確認入力
```bash
interactive-input \
  --prompt "Move originals to archive? " \
  --input-type confirm \
  --key move
```
`y` なら `--move` が返り、`n` なら何も出力されません。`[y] Yes / [n] No` の案内は自動でstderrに追記されます。

## スクリプトでの利用例
```bash
INTERACTIVE_BIN="./pkg/bin/cli/win_amd64/interactive-input.exe"
MOVE_FLAG="$($INTERACTIVE_BIN \
  --prompt 'Move originals? ' \
  --input-type confirm \
  --key move)"

if [ -n "$MOVE_FLAG" ]; then
  echo "will move with flag: $MOVE_FLAG"
fi
```
標準出力のみをコマンド置換で受け取り、stderrのプロンプト表示はそのままユーザーへ見せられます。

## リトライ挙動
- `--max-attempts=3` の場合、バリデーション失敗3回目で終了コード1を返します。
- `--max-attempts=0` は無制限に再入力を促します。
- `Ctrl+D` などで標準入力が閉じられた場合はユーザーキャンセルとみなし終了コード1を返します。

## エラーコード
| コード | 内容 |
| ------ | ---- |
| 0 | 成功 |
| 1 | ユーザーキャンセル、もしくはリトライ上限到達 |
| 2 | フラグの誤りなど予期しないエラー |
