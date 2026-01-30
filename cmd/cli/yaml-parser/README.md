# YAML Parser

YAMLファイルや直接指定したYAML文字列を解析して、構造化データをJSONとして出力するCLIツールです。設定値の確認やテストデータの検証に利用できます。

## 機能

- **read**: `--file-path` で指定したYAMLファイルを読み込み、構造を解析してJSON形式で表示
- **parse**: `--yaml-content` で受け取ったYAML文字列を解析し、JSON形式で表示
- **複数ドキュメント対応**: `---` 区切りで複数ドキュメントが存在する場合は配列として出力
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

## オプション

| オプション | 必須 | 説明 |
|------------|------|------|
| `--operation` | ✔ | `read` または `parse` を指定 |
| `--file-path` | `read` 操作で必要 | 読み込むYAMLファイルのパス |
| `--yaml-content` | `parse` 操作で必要 | 解析対象のYAML文字列 |

## 例

```bash
# 設定ファイルをJSONで確認
./bin/yaml-parser --operation read --file-path ./configs/app.yaml

# CI のENV YAMLをそのまま解析
./bin/yaml-parser --operation parse --yaml-content $'env:\n  stage: prod'
```
