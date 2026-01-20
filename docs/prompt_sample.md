# コーディングエージェントへのサンプルプロンプト集
過去に使ったプロンプトをまとめました。

## 専用MCPを利用するエージェントの設定ファイルの作成（ペルソナ抽出）
/home/user/devbox/.config/cagent/config_linux/persona_extractor/config.ymlに、実際にpersona_extraction MCPツールを使うペルソナ抽出エージェントの設定ファイルを作って。そのペルソナ抽出エージェントは、ユーザーから文章をプロンプトされたら、その文章内の各登場人物のペルソナを抽出するのである。その文章はプロンプトから提供されても受け付けて、ファイルパスを指定されて提供されても受け付けられるようにしたい。

## cagentワークフローの新規作成
/home/user/devbox/.config/cagent/config_linux/github_pull_request_creator/config.ymlを参考に、/home/user/devbox/.config/cagent/config_linux/github_issue_creator/config.ymlにGithubに新規イシューを起票するためのエージェント用のワークフローを作成して。イシュー本文に書く内容をユーザーに尋ねるようにして。

## Taskfileへのタスクの追加（分割パターン）
/home/user/devbox/pkg/taskfile/taskfiles/exif.ymlに下記のファイルにあるタスクを下記のファイルの代わりに行えるようにして、/home/user/devbox/pkg/taskfile/Taskfile.ymlから呼び出すようにして。
- /home/user/devbox/pkg/win_dos/Z3-4_exif_modifier_from_filename.bat
- 

## 新規のCLIツールの追加
/home/user/devbox/cmd/cli/taskfile/main.goにTaskfileを管理するためのCLIツールを実装して。CLIフラグには"--operation"、"--task-type"、"--taskfile-path"を受け付けるようにして。
- "--operation"には、"inspect"を受け付ける。
- "--task-type"には、"root"を受け付ける。
- --operation: inspect, --task-type: root の場合は、/home/user/devbox/internal/taskfile/usecases/taskfiles/root.ymlにあるフィールドが、"--taskfile-path"で渡されたTaskfileで不足していないかどうかを確認する。
