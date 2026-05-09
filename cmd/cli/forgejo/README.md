# Forgejo CLI

Forgejo のリポジトリ情報を取得する CLI です。

## 機能

- `repo list`: `repo list` で自分がアクセス可能なリポジトリ一覧を取得
- `project list`: 指定ユーザーのリポジトリ配下のプロジェクト一覧を取得

認証は `forgejo-token` の API トークンを使用します。  
`.env`（または OS 環境変数）から以下のキーを読み込みます。

- `forgejo-host`
- `forgejo-username`
- `forgejo-token`

## 使い方

```bash
./forgejo -operation "repo list" -forgejo-host https://codeberg.org -forgejo-username YOUR_NAME -forgejo-token YOUR_TOKEN

./forgejo -operation "project list" -forgejo-host https://codeberg.org -forgejo-username YOUR_NAME -forgejo-token YOUR_TOKEN

./forgejo repo list -forgejo-host https://codeberg.org -forgejo-username YOUR_NAME -forgejo-token YOUR_TOKEN
./forgejo project list -forgejo-username YOUR_NAME
```

## パラメータ

### 必須

| 名前 | 説明 |
| --- | --- |
| `-operation` | 実行操作 (`repo list` / `project list`) |
| `forgejo-host` | Forgejo ホスト（例: `https://codeberg.org`） |
| `forgejo-username` | ユーザー名 |
| `forgejo-token` | API トークン |

### 任意

| 名前 | 説明 | デフォルト |
| --- | --- | --- |
| `-json` | 結果を JSON で整形して出力 | `false` |
| `-help`, `-h` | このヘルプを表示 | `false` |

## 使用例

```bash
./forgejo -operation "repo list" -forgejo-host https://codeberg.org -forgejo-username yourname -forgejo-token xxx -json

./forgejo project list -forgejo-host https://codeberg.org -forgejo-username yourname -forgejo-token xxx
```

`.env` に設定する例:

```bash
cat <<EOF > .env
forgejo-host=https://codeberg.org
forgejo-username=yourname
forgejo-token=xxx
EOF

./forgejo repo list
```

## 出力

### `repo list`

出力項目:

- `name`
- `description`
- `is_private`
- `http_url`
- `issues_count`
- `pulls_count`
- `forks_count`
- `stargazers_count`
- `subscribers_count`
- `language`
- `languages`
- `size`
- `repo_created_at`
- `repo_updated_at`
- `is_archived`
- `tags`

`-json` 指定時は以下のような JSON 形式で出力されます。

```json
[
  {
    "name": "repo-name",
    "description": "",
    "is_private": false,
    "http_url": "https://forgejo.example.com/owner/repo-name",
    "issues_count": 0,
    "pulls_count": 1,
    "forks_count": 0,
    "stargazers_count": 10,
    "subscribers_count": 0,
    "language": "Go",
    "languages": {
      "Go": 1234.0
    },
    "size": 100,
    "repo_created_at": "2022-10-18T00:00:00Z",
    "repo_updated_at": "2022-10-19T00:00:00Z",
    "is_archived": false,
    "tags": "game,go"
  }
]
```

### `project list`

出力項目:

- `name`
- `description`
- `is_private`
- `is_archived`
- `repo_full_name`
- `created_at`
- `updated_at`

`-json` 指定時は以下のような JSON 形式で出力されます。

```json
[
  {
    "name": "Backend",
    "description": "infra",
    "is_private": false,
    "is_archived": false,
    "repo_full_name": "owner/repo",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
]
```

## エラー

- 必須パラメータ不足時はエラーメッセージを表示して終了します
- API 呼び出し失敗時は `stderr` にエラー内容を表示します
- `project list` でプロジェクト API が未対応の場合は  
  `project list API is not supported on this server`
