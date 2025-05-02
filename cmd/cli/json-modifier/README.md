# JSON Modifier ツール使用方法

## 概要

JSON Modifierは、JSONファイルに任意のキーと値を追加したり、既存のキーの値を取得したりするためのコマンドラインツールです。単一のJSONファイルだけでなく、ディレクトリ内の全てのJSONファイルに対しても操作できます。

## インストール

ビルドスクリプトを使用してツールをビルドします：

```bash
cd devbox
./scripts/build_json_modifier.sh
```

ビルドが成功すると、以下の場所に実行ファイルが生成されます：

- Linux用: `pkg/bin/linux_amd64/json-modifier`
- Windows用: `pkg/bin/win_amd64/json-modifier.exe`

## 基本的な使い方

### 単一ファイルの操作

#### キーと値の追加

JSONファイルに新しいキーと値を追加するには、以下のコマンドを使用します：

```bash
./json-modifier -file <JSONファイルパス> -key <キー> -set <値>
```

例：

```bash
./json-modifier -file data.json -key description -set "これはサンプルJSONファイルです"
```

このコマンドは、`data.json`ファイルに`description`キーと値`"これはサンプルJSONファイルです"`を追加します。

- ファイルが存在しない場合は、新しいJSONファイルが作成されます。
- 指定されたキーが既に存在する場合は、値が上書きされます。

#### 値の取得

JSONファイルから特定のキーの値を取得するには、以下のコマンドを使用します：

```bash
./json-modifier -file <JSONファイルパス> -key <キー> -get
```

例：

```bash
./json-modifier -file data.json -key description -get
```

このコマンドは、`data.json`ファイルから`description`キーの値を取得して表示します。

#### すべてのデータの取得

JSONファイルのすべてのキーと値を取得するには、以下のコマンドを使用します：

```bash
./json-modifier -file <JSONファイルパス> -get-all
```

例：

```bash
./json-modifier -file data.json -get-all
```

このコマンドは、`data.json`ファイルのすべてのキーと値をJSON形式で表示します。

### ディレクトリ内の全てのJSONファイルの操作

#### ディレクトリ内の全てのJSONファイルにキーと値を追加

ディレクトリ内の全てのJSONファイルに新しいキーと値を追加するには、以下のコマンドを使用します：

```bash
./json-modifier -dir <ディレクトリパス> -key <キー> -set <値> [-recursive]
```

例：

```bash
./json-modifier -dir ./data -key description -set "これはサンプルJSONファイルです"
```

このコマンドは、`./data`ディレクトリ内の全てのJSONファイルに`description`キーと値`"これはサンプルJSONファイルです"`を追加します。

```bash
./json-modifier -dir ./data -key version -set "1.0.0" -recursive
```

このコマンドは、`./data`ディレクトリとそのサブディレクトリ内の全てのJSONファイルに`version`キーと値`"1.0.0"`を追加します。

- `-recursive`オプションを指定すると、サブディレクトリ内のJSONファイルも処理されます。
- 処理されたファイル数が出力されます。
- ディレクトリモードでは、`-get`と`-get-all`オプションは使用できません。

## オプション

| オプション | 説明 |
|------------|------|
| `-file` | 操作するJSONファイルのパス |
| `-dir` | 操作するJSONファイルが含まれるディレクトリのパス |
| `-recursive` | ディレクトリを再帰的に処理する（`-dir`オプションと共に使用）（デフォルト: `false`） |
| `-key` | 追加または取得するキー（`-get`または値の追加時に必須） |
| `-set` | 追加する値（値の追加時に必須）。文字列、整数、浮動小数点数を指定可能 |
| `-get` | 指定されたキーの値を取得するフラグ（単一ファイルモードのみ） |
| `-get-all` | すべてのキーと値を取得するフラグ（単一ファイルモードのみ） |

**注意**: `-file`と`-dir`オプションは同時に指定できません。どちらか一方を指定してください。

## エラーメッセージ

以下のような場合にエラーメッセージが表示されます：

- ファイルパスが指定されていない場合
- キーが指定されていない場合（`-get`または値の追加時）
- 値が指定されていない場合（値の追加時）
- ファイルが存在しない場合（`-get`または`-get-all`時）
- 指定されたキーが存在しない場合（`-get`時）
- JSONファイルの読み込みや書き込みに失敗した場合

## 使用例

### 単一ファイルの操作

```bash
# 新しいJSONファイルの作成
./json-modifier -file new_data.json -key name -set "サンプル"
./json-modifier -file new_data.json -key version -set "1.0.0"

# 整数値の追加
./json-modifier -file data.json -key count -set 42

# 浮動小数点数の追加
./json-modifier -file data.json -key price -set 19.99

# 既存のキーの値を更新
./json-modifier -file data.json -key version -set "2.0.0"

# キーの値を取得
./json-modifier -file data.json -key name -get

# すべてのデータを取得
./json-modifier -file data.json -get-all
```

### ディレクトリ内の全てのJSONファイルの操作

```bash
# ディレクトリ内の全てのJSONファイルにキーと値を追加
./json-modifier -dir ./data -key author -set "サンプル作成者"

# ディレクトリとサブディレクトリ内の全てのJSONファイルにキーと値を追加
./json-modifier -dir ./data -key version -set "1.0.0" -recursive

# ディレクトリ内の全てのJSONファイルに整数値を追加
./json-modifier -dir ./data -key count -set 42 -recursive

# ディレクトリ内の全てのJSONファイルに浮動小数点数を追加
./json-modifier -dir ./data -key price -set 19.99 -recursive
```

## 注意事項

- JSONファイルは常にUTF-8エンコーディングで処理されます。
- 値は自動的に型変換されます。整数や浮動小数点数は数値として、それ以外は文字列として追加されます。
- ファイルパスとディレクトリパスは相対パスまたは絶対パスで指定できます。
- ディレクトリモードでは、指定したディレクトリ内のJSONファイルのみが処理されます。
- `-recursive`オプションを指定すると、サブディレクトリ内のJSONファイルも処理されます。
- ディレクトリモードでは、`-get`と`-get-all`オプションは使用できません。
- 処理中にエラーが発生した場合でも、可能な限り処理を継続し、最終的に処理されたファイル数が表示されます。
