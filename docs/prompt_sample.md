# コーディングエージェントへのサンプルプロンプト集
過去に使ったプロンプトをまとめました。

## 専用MCPを利用するエージェントの設定ファイルの作成（ペルソナ抽出）
/home/user/devbox/.config/cagent/config_linux/persona_extractor/config.ymlに、実際にpersona_extraction MCPツールを使うペルソナ抽出エージェントの設定ファイルを作って。そのペルソナ抽出エージェントは、ユーザーから文章をプロンプトされたら、その文章内の各登場人物のペルソナを抽出するのである。その文章はプロンプトから提供されても受け付けて、ファイルパスを指定されて提供されても受け付けられるようにしたい。

## cagentワークフローの新規作成
/home/user/devbox/.config/cagent/config_linux/github_pull_request_creator/config.ymlを参考に、/home/user/devbox/.config/cagent/config_linux/github_issue_creator/config.ymlにGithubに新規イシューを起票するためのエージェント用のワークフローを作成して。イシュー本文に書く内容をユーザーに尋ねるようにして。

## Taskfileへのタスクの追加
/home/user/devbox/pkg/win_dos/Taskfile.ymlで、下記のファイルにあるタスクを下記のファイルの代わりに行えるようにして。
- /home/user/devbox/pkg/win_dos/batch_files/Z2-3_image_renamer_for_content_of_date.bat
- /home/user/devbox/pkg/win_dos/batch_files/Z2-3_image_renamer_for_content_of_mackerel.bat

## Taskfileへのタスクの追加（分割パターン）
/home/user/devbox/pkg/win_dos/Taskfile.image_convert.ymlに下記のファイルにあるタスクを下記のファイルの代わりに行えるようにして、/home/user/devbox/pkg/win_dos/Taskfile.ymlから呼び出すようにして。
- /home/user/devbox/pkg/win_dos/batch_files/Z3-1_image_converter_jpg_to_webp.bat
- 
