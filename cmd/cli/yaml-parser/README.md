# YAML Parser

YAMLファイルや直接指定したYAML文字列を解析して、構造化データをJSONとして出力するCLIツールです。設定値の確認やテストデータの検証に利用できます。

## 機能

- **read**: `--file-path` で指定したYAMLファイルを読み込み、構造を解析してJSON形式で表示
- **parse**: `--yaml-content` で受け取ったYAML文字列を解析し、JSON形式で表示
- **edit-file**: `--file-path` と `--key-value-list` で指定したキーを更新し、書き戻した結果をJSONで表示（ドット区切りでネスト指定可）
- **複数ドキュメント対応**: `---` 区切りで複数ドキュメントが存在する場合は配列として出力（`edit-file` は単一ドキュメントのみ対応）
- **キーの正規化**: YAML特有の `map[interface{}]interface{}` をJSONに適した `map[string]interface{}` に変換

## インストール

```bash
# プロジェクトルートから
go build -o bin/yaml-parser ./cmd/cli/yaml-parser
```

## 使用方法

### ファイルを読み取る (`read`)

```bash
go run ./cmd/cli/yaml-parser \
  --operation read \
  --file-path ./configs/app.yaml
```

出力例:

```
{
  "server": {
    "port": 8080,
    "host": "localhost"
  },
  "debug": true
}
```

### 文字列から解析する (`parse`)

```bash
go run ./cmd/cli/yaml-parser \
  --operation parse \
  --yaml-content "name: demo\nfeatures:\n  - api\n  - cli"
```

出力例:

```
{
  "name": "demo",
  "features": [
    "api",
    "cli"
  ]
}
```

### 複数ドキュメント

```bash
go run ./cmd/cli/yaml-parser \
  --operation parse \
  --yaml-content $'---\nname: doc1\n---\nname: doc2'
```

出力例:

```
[
  {
    "name": "doc1"
  },
  {
    "name": "doc2"
  }
]
```

### キーを書き換える (`edit-file`)

`--key-value-list` には `key=value` 形式をカンマまたは改行で並べます。値は YAML として解釈されるため、`true` や `123`、`[1,2,3]` などもそのまま指定できます。ネストは `parent.child` のようにドット記法で表現し、配列要素は `servers.0.port` のようにインデックスで指定できます。

```bash
go run ./cmd/cli/yaml-parser \
  --operation edit-file \
  --file-path ./configs/app.yaml \
  --key-value-list $'server.port=9090\ndebug=true\ninfo.region="ap-northeast-1"'
```

出力例:

```
{
  "debug": true,
  "info": {
    "env": "dev",
    "region": "ap-northeast-1"
  },
  "server": {
    "host": "localhost",
    "port": 9090
  }
}
```

> **Note**: `edit-file` は単一ドキュメント（1つの `---` ブロック）を対象としており、複数ドキュメントが含まれるファイルではエラーになります。

## オプション

| オプション | 必須 | 説明 |
|------------|------|------|
| `--operation` | ✔ | `read` / `parse` / `edit-file` を指定 |
| `--file-path` | `read` / `edit-file` で必要 | 対象YAMLファイルのパス |
| `--yaml-content` | `parse` 操作で必要 | 解析対象のYAML文字列 |
| `--key-value-list` | `edit-file` 操作で必要 | `key=value` をカンマ/改行区切りで列挙（例: `server.port=8081,debug=true`） |

## 例

```bash
# 設定ファイルをJSONで確認
./bin/yaml-parser --operation read --file-path ./configs/app.yaml

# CI のENV YAMLをそのまま解析
./bin/yaml-parser --operation parse --yaml-content $'env:\n  stage: prod'

# 設定ファイルを書き換え（複数キーを改行で指定）
./bin/yaml-parser --operation edit-file --file-path ./configs/app.yaml --key-value-list $'server.port=9090\ndebug=false'
```
