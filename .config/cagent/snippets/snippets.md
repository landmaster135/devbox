# cagent 実行スニペット

## Linux 用: `config_linux`
`/home/user` 配下で動かす前提の設定。例として Codespaces や CI の非特権ユーザー環境で利用。

### Gitコミットメッセージ生成 (外部API 利用)
非ステージング差分を解析:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux/git_commit_message_generator/config_prod.yml "/git-commit-message" --yolo
```
ステージング済み差分のみ:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux/git_commit_message_generator/config_prod.yml "git-commit-message-staged" --yolo
```
カレントディレクトリを `git_dir` に指定 (非ステージング):
```bash
cagent run /home/user/devbox/.config/cagent/config_linux/git_commit_message_generator/config_prod.yml "git-commit-message-pwd" --yolo
```
カレントディレクトリ + ステージングのみ:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux/git_commit_message_generator/config_prod.yml "git-commit-message-staged-pwd" --yolo
```

### 逆張り応答エージェント
単発プロンプトで実行:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux_prod/contradict.yml "I think TypeScript is the best."
```
標準入力を使った逐次実行:
```bash
echo "Rust is difficult." | cagent run /home/user/devbox/.config/cagent/config_linux_prod/contradict.yml -
```

### タスク計画エージェント
単発プロンプトで実行:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux/task_planner/config_prod.yml 'plan-from-consult'
```

### 学習支援エージェント Alloy
学習テーマを指定して実行:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux_prod/alloy.yml "Help me understand differential equations."
```
既存セッションに追加メッセージを送る場合:
```bash
cagent run /home/user/devbox/.config/cagent/config_linux_prod/alloy.yml "Can you give me practice problems?"
```

## Windows 用: `config_win`
PowerShell での実行例。事前に `cagent.exe` が PATH にあり、`OLLAMA_API_KEY` を任意のトークン文字列で設定しておく。

### Gitコミットメッセージ生成 (Ollama 利用)
非ステージング差分:
```powershell
cagent.exe run $Env:USERPROFILE\.config\cagent\config_win\git_commit_message_generator_ollama.yml --command git-commit-message
```
ステージングのみ:
```powershell
cagent.exe run $Env:USERPROFILE\.config\cagent\config_win\git_commit_message_generator_ollama.yml --command git-commit-message-staged
```
カレントディレクトリ対象 (非ステージング):
```powershell
cagent.exe run $Env:USERPROFILE\.config\cagent\config_win\git_commit_message_generator_ollama.yml --command git-commit-message-pwd
```
カレントディレクトリ対象 (ステージングのみ):
```powershell
cagent.exe run $Env:USERPROFILE\.config\cagent\config_win\git_commit_message_generator_ollama.yml --command git-commit-message-staged-pwd
```

## メモ
- いずれのスニペットも TUI モードが既定で起動します。非TUIで動かす場合は `--tui=false` を適宜追加してください。
