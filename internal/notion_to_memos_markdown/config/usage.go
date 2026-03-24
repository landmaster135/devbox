package config

import (
	flagParser "github.com/landmaster135/devbox/internal/notion_to_memos_markdown/infrastructures/flag_parser"
)

const usageTemplate = `notion-to-memos-markdown CLI

使用方法:
  %[1]s --operation=distribute-files --page-type=content --src-json-file=./tmp/contents.json --src-body-dir=./tmp/body --out-dir=./tmp/out
  %[1]s --operation=craft-markdown --page-type=content --category=software --skips-no-src-body=false --con_number_start=1 --con_number_end=9999 --src-json-file=./tmp/contents.json --src-body-dir=./tmp/body --out-dir=./tmp/out
  %[1]s --operation=craft-markdown --page-type=artifact --con_number_start=1 --con_number_end=9999 --src-json-path=./tmp/artifacts.json --src-body-dir=./tmp/body --out-dir=./tmp/out
  %[1]s --operation=craft-markdown --page-type=task --con_number_start=1 --con_number_end=9999 --src-json-path=./tmp/tasks.json --src-body-dir=./tmp/body --out-dir=./tmp/out
  %[1]s --operation=check-body-length --src-body-dir=./tmp/body --threshold=1000
  %[1]s --operation=grep-str --src-body-dir=./tmp/body --target-str=TODO
  %[1]s --operation=rename-bodies-by-category-id --page-type=content --con_number_start=1 --con_number_end=9999 --src-json-file=./tmp/contents.json --src-resource-dir=./tmp/resources
  %[1]s --operation=migrate-to-memos --page-type=content --base-url=https://memos.example.com --api-token=token --src-body-dir=./tmp/body --src-resource-dir=./tmp/resources

オプション:
  --operation        操作タイプ（必須: distribute-files, craft-markdown, check-body-length, grep-str, rename-bodies-by-category-id, migrate-to-memos）
  --page-type        ページタイプ（distribute-files/rename-bodies-by-category-id/migrate-to-memos: content, craft-markdown: content|artifact|task）
  --base-url         Memos API のベースURL（migrate-to-memosで必須）
  --api-token        Memos API のトークン（migrate-to-memosで必須）
  --category         対象category（craft-markdownで任意。指定時は一致するContentのみ処理）
  --skips-no-src-body コピー元Markdownなしをスキップするか（craft-markdownで任意。デフォルト:false）
  --con_number_start craft-markdown/rename-bodies-by-category-id時の開始con番号（必須）
  --con_number_end   craft-markdown/rename-bodies-by-category-id時の終了con番号（必須）
  --threshold        文字数の閾値（check-body-lengthで必須、0以上）
  --target-str       検索文字列（grep-strで必須）
  --src-json-file    入力JSONファイルのパス（distribute-files/craft-markdown/rename-bodies-by-category-idで必須）
  --src-body-dir     入力ディレクトリ（distribute-files/craft-markdown/check-body-length/grep-str/migrate-to-memosで必須）
  --src-resource-dir リソース入力ディレクトリ（rename-bodies-by-category-id/migrate-to-memosで必須）
  --out-dir          出力先ルートディレクトリ（distribute-files/craft-markdownで必須）
  -help, -h          このヘルプを表示
`

func PrintUsage() {
	flagParser.PrintUsage(usageTemplate)
}
