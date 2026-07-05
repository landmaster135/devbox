package config

import (
	"fmt"
	"os"

	flagParser "github.com/landmaster135/devbox/internal/forgejo/infrastructures/flag_parser"
)

const usageTemplate = `Forgejo CLI

使用方法:
  %s -operation "repo list"
  %s -operation "issue list"
  %s repo list
  %s issue list

  位置引数:
  repo list
  issue list

オプション:
  -operation          実行する操作 (repo list, issue list)
  -repos-workers  repo list の同時ワーカー数 (初期値: 4)
  -json               JSON形式で出力
  -help, -h           ヘルプを表示

環境変数:
  .env (またはOS環境変数)
    FORGEJO_HOST
    FORGEJO_USERNAME
    FORGEJO_TOKEN

注意:
  カレントディレクトリの .env を読み込みます。

`

// PrintUsage は使用方法を表示します。
func PrintUsage() {
	command := os.Args[0]
	flagParser.PrintUsage(fmt.Sprintf(usageTemplate, command, command, command, command))
}
