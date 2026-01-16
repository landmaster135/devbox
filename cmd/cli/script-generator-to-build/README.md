# Script Generator to Build

このツールは、指定されたGoパッケージのビルドスクリプトを生成するためのコマンドラインツールです。

## 機能

- 指定されたGoパッケージのビルドスクリプトを自動生成
- パッケージが指定されない場合は、利用可能なパッケージの一覧から選択可能
- 生成されたビルドスクリプトは、Linux/AMD64、Windows/AMD64、macOS/ARM64向けのクロスコンパイルを行う
- READMEファイルから使用例を自動抽出して、生成されるビルドスクリプトに含める

### READMEからの使用例抽出

このツールは、対象パッケージのREADMEファイルから使用例を自動的に抽出します：

1. 「## 使用例」セクションを検索
   - 見つからない場合は「## 使用方法」セクションを検索
2. セクション内の最初のコードブロック（```bashで囲まれた部分）を抽出
   - サブセクション（### で始まる見出し）内のコードブロックも検索対象
3. 抽出したコードを生成するビルドスクリプトの使用例として追加

これにより、パッケージの開発者が提供した使用例がそのままビルドスクリプトに反映され、ユーザーは簡単に正しい使い方を確認できます。

## 使用方法

```
script-generator-to-build [パッケージ名]
```

### オプション

- `-h`, `--help`: ヘルプメッセージを表示
- `--base_dir <path>`: CLIやスクリプトの探索元となるベースディレクトリを指定（省略時はカレントディレクトリ）
- `--cli_dir <path>`: ベースディレクトリからの相対パスでCLIディレクトリを指定（デフォルト: `cmd/cli`）
- `--scripts_dir <path>`: ベースディレクトリからの相対パスでスクリプトディレクトリを指定（デフォルト: `scripts`）
- `--output_dir <path>`: 生成したビルドスクリプトの出力先ディレクトリを指定（デフォルト: `./pkg/bin/cli`）
- `--package_name <name>`: 対象パッケージ名をフラグで直接指定

## 出力

指定されたパッケージに対応するビルドスクリプト（`build_<パッケージ名>.sh`）が生成されます。

## 使用例

```bash
# code-analyzerパッケージのビルドスクリプトを生成
go run ./cmd/cli/script-generator-to-build --package_name code-analyzer

# 対話的にパッケージを選択してビルドスクリプトを生成
go run ./cmd/cli/script-generator-to-build

# CLIや出力ディレクトリをカスタマイズ
go run ./cmd/cli/script-generator-to-build --base_dir /path/to/project --cli_dir custom/cmd --output_dir ./artifacts code-analyzer
```
