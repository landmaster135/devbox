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

## 新規のCLIツールの追加の実装計画の作成
/home/user/devbox/cmd/cli/gcloud-genset-spanner/main.goにCloud Spanner用のgcloudコマンドを出力するCLIツールを実装したい。そのCLIツールでは、--operationとその他のCLIフラグを受け取る。下記の--operationを実装する計画を/home/user/devbox/.agents/draft.mdに書いて。実装の仕方は/home/user/devbox/cmd/cli/gcloud-genset-cloudsql/main.goを参考にして。

### --operation: instance-list
下記のgcloudコマンドを出力する。
```
gcloud spanner instances list
```

### --operation: instance-create
--instance-id、--config、--description、--nodesのCLIフラグも受け取る。
```
gcloud spanner instances create my-instance \
    --config=regional-asia-northeast1 \
    --description="My Spanner Instance" \
    --nodes=1
```

### --operation: db-create
--instance-id、--db-id、--ddl-file-pathのCLIフラグも受け取る。そして下記のgcloudコマンドを出力する。
```
gcloud spanner databases create my-database \
    --instance=my-instance \
    --ddl-file=schema.sql
```

### --operation: db-list
--instance-idのCLIフラグも受け取る。そして下記のgcloudコマンドを出力する。
```
gcloud spanner databases list --instance=my-instance
```

### --operation: db-describe
--instance-id、--db-idのCLIフラグも受け取る。そして下記のgcloudコマンドを出力する。
```
gcloud spanner databases describe my-database --instance=my-instance
```

## 新規のCronワークフローの追加
/home/user/devbox/cmd/cli/cron-workflow/workflow/core.goにDiscordに天気予報を通知するための新規ワークフローを追加して。/home/user/devbox/cmd/cli/weather-notificator/main.go内で呼び出されているservice.HandleWeatherNotification関数を呼び出せば処理出来るはずだ。city: "Tokyo", maxDays: 3を渡して、/home/user/devbox/cmd/cli/cron-workflow/workflow/env.go内にある下記の環境変数にある値も渡す。
- EnvKeyDiscordWebhookURL
- EnvKeyOpenWeatherAPIKey

## 新規のCronワークフローの追加（PC情報取得タスクの場合）
/home/user/devbox/cmd/cli/cron-workflow/workflow/core.goにPC情報を取得するための新規ワークフローを追加して。/home/user/devbox/internal/machine_info/usecases/services.goにあるCollectAndSaveUbuntuInfo関数を呼び出せば処理出来るはずだ。引数networkInterfaceに "eth0"を渡して、/home/user/devbox/cmd/cli/cron-workflow/workflow/env.go内にある下記の環境変数にある値をoutDirとして、引数outputDirにfilepath.Join(c.VolumeDir, outDir)も渡す。cronは10分おきに設定して。
- EnvKeyPCInfoOutputDirectory

## templコンポーネントに対するテストケースの追加
/home/user/devbox/internal/templ_components/core/heading/index_test.goと/home/user/devbox/internal/templ_components/core/hidden_input/input_test.goを参考に/home/user/devbox/internal/templ_components/core/paragraphにテストコードを追加して。
