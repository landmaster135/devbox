# Forgejo CLI

Forgejo のリポジトリ/Issue 情報を取得する CLI です。

## 機能

- `repo list`: `repo list` で自分がアクセス可能なリポジトリ一覧を取得
- `issue list`: 指定ユーザーのリポジトリ配下の issue 一覧を取得

認証は `forgejo-token` の API トークンを使用します。  
`.env`（または OS 環境変数）から以下のキーを読み込みます。

- `forgejo-host`
- `forgejo-username`
- `forgejo-token`

## 使い方

```bash
./forgejo -operation "repo list" -forgejo-host https://codeberg.org -forgejo-username YOUR_NAME -forgejo-token YOUR_TOKEN

./forgejo -operation "issue list" -forgejo-host https://codeberg.org -forgejo-username YOUR_NAME -forgejo-token YOUR_TOKEN

./forgejo repo list -forgejo-host https://codeberg.org -forgejo-username YOUR_NAME -forgejo-token YOUR_TOKEN
./forgejo issue list -forgejo-username YOUR_NAME
```

## パラメータ

### 必須

| 名前 | 説明 |
| --- | --- |
| `-operation` | 実行操作 (`repo list` / `issue list`) |
| `forgejo-host` | Forgejo ホスト（例: `https://codeberg.org`） |
| `forgejo-username` | ユーザー名 |
| `forgejo-token` | API トークン |

### 任意

| 名前 | 説明 | デフォルト |
| --- | --- | --- |
| `-json` | 結果を JSON で整形して出力 | `false` |
| `-repos-workers` | repo list の同時ワーカー数 | `4` |
| `-help`, `-h` | このヘルプを表示 | `false` |

## 使用例

```bash
./forgejo -operation "repo list" -forgejo-host https://codeberg.org -forgejo-username yourname -forgejo-token xxx -json
./forgejo -operation "repo list" -forgejo-host https://codeberg.org -forgejo-username yourname -forgejo-token xxx -repos-workers 8

./forgejo issue list -forgejo-host https://codeberg.org -forgejo-username yourname -forgejo-token xxx
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
- `open_issues_count`
- `closed_issues_count`
- `open_pulls_count`
- `closed_pulls_count`
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
    "open_issues_count": 1,
    "closed_issues_count": 0,
    "open_pulls_count": 2,
    "closed_pulls_count": 1,
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

### `issue list`

出力項目:

- `repo_full_name`
- `number`
- `title`
- `state`
- `html_url`
- `author`
- `assignees`
- `labels`
- `comments`
- `is_locked`
- `created_at`
- `updated_at`
- `closed_at`

`-json` 指定時は以下のような JSON 形式で出力されます。

```json
[
  {
    "repo_full_name": "owner/repo",
    "number": 5,
    "title": "Bug fix",
    "state": "open",
    "html_url": "https://forgejo.example.com/owner/repo/issues/5",
    "author": "alice",
    "assignees": [
      "bob"
    ],
    "labels": [
      "bug"
    ],
    "comments": 2,
    "is_locked": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z",
    "closed_at": ""
  }
]
```

## エラー

- 必須パラメータ不足時はエラーメッセージを表示して終了します
- API 呼び出し失敗時は `stderr` にエラー内容を表示します
