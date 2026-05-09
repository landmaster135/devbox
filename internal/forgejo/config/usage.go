package config

import (
	"fmt"
	"os"

	flagParser "github.com/landmaster135/devbox/internal/forgejo/infrastructures/flag_parser"
)

const usageTemplate = `Forgejo CLI

使用方法:
  %s -operation "repo list" -forgejo-host https://codeberg.org -forgejo-username user -forgejo-token your_token
  %s -operation "project list" -forgejo-host https://codeberg.org -forgejo-username user -forgejo-token your_token

  位置引数:
  repo list
  project list

オプション:
  -operation          実行する操作 (repo list, project list)
  -forgejo-host       Forgejoのホスト URL (例: https://codeberg.org)
  -forgejo-username   Forgejoのユーザー名
  -forgejo-token      Forgejo APIトークン
  -repos-workers  repo list の同時ワーカー数 (初期値: 4)
  -json               JSON形式で出力
  -help, -h           ヘルプを表示

環境変数:
  .env (またはOS環境変数)
    forgejo-host
    forgejo-username
    forgejo-token

注意:
  カレントディレクトリの .env を読み込みます。

`

// PrintUsage は使用方法を表示します。
func PrintUsage() {
	command := os.Args[0]
	flagParser.PrintUsage(fmt.Sprintf(usageTemplate, command, command))
}
