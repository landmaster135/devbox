# Withings CLI

Withings Public Health Data API の OAuth フローに沿って認可 URL の生成、アクセストークン／リフレッシュトークンの取得、そして日次のヘルスデータ取得を行う CLI ツールです。

## 提供コマンド

| operation        | 説明 | 必須パラメータ |
|------------------|------|----------------|
| `auth-url`       | 認可 URL を生成し、ブラウザで開くべきリンクを出力します。 | `-client-id`, `-redirect-uri`（任意で `-scope`, `-state`, `-mode`） |
| `request-token`  | 認可コードからアクセストークン／リフレッシュトークン／ユーザーIDを取得します。 | `-client-id`, `-client-secret`, `-authorization-code`, `-redirect-uri` |
| `refresh-token`  | リフレッシュトークンを使って新しいアクセストークンを取得します。 | `-client-id`, `-client-secret`, `-refresh-token` |
| `daily-summary` (既定) | 指定期間の測定値と活動サマリを取得して JSON 出力します。 | `-access-token`, `-user-id`, `-start-date`（任意で `-end-date`, `-measure-types`, `-include-activity`） |

`-operation` を省略した場合は `daily-summary` が実行されます。

## OAuth フローに沿った利用手順

1. **認可 URL を生成**
```bash
go run ./cmd/cli/withings \
  -operation auth-url \
  -client-id "YOUR_CLIENT_ID" \
  -redirect-uri "https://yourapp.example/oauth/callback" \
  -scope "user.metrics,user.activity" \
  -state "random_state_value"
```
表示された URL をブラウザで開くと、ユーザーが Withings のログインと許可を行い、設定した `redirect_uri` に `code` と `state` が付加されてリダイレクトされます。`code` はアクセストークンに交換することが可能な認可コードとなります。（リダイレクト先で必ずしもサーバが起動されている必要はありません。）

2. **認可コードをアクセストークンへ交換**
```bash
go run ./cmd/cli/withings \
  -operation request-token \
  -client-id "YOUR_CLIENT_ID" \
  -client-secret "YOUR_CLIENT_SECRET" \
  -authorization-code "CODE_FROM_CALLBACK" \
  -redirect-uri "https://yourapp.example/oauth/callback"
```
標準出力に JSON が表示され、その中に `body.userid`, `body.access_token`, `body.refresh_token` が含まれます。リフレッシュトークンは長期的なアクセスに必要なので安全な場所へ保存してください。

3. **リフレッシュトークンでアクセストークンを更新**（任意、後日再実行する際）
```bash
go run ./cmd/cli/withings \
  -operation refresh-token \
  -client-id "YOUR_CLIENT_ID" \
  -client-secret "YOUR_CLIENT_SECRET" \
  -refresh-token "STORED_REFRESH_TOKEN"
```
新しいアクセストークンと最新のリフレッシュトークンが出力されます。Withings では新しいリフレッシュトークンを受け取った後 8 時間で旧アクセストークンが失効するため、常に最新の値を保管してください。

4. **日次サマリを取得**
```bash
go run ./cmd/cli/withings \
  -operation daily-summary \
  -access-token "ACCESS_TOKEN" \
  -user-id 12345678 \
  -start-date 2025-09-01 \
  -end-date 2025-09-07 \
  -measure-types weight,diastolic,systolic,heart_rate
```
体重・血圧・アクティビティなどの日次サマリが JSON で出力されます。`-include-activity=false` を指定すると活動サマリを省くことも可能です。

### 測定タイプの指定 (`-measure-types`)

- `all` を指定すると Withings Measure API のフィルタを無効化し、取得可能なすべての測定値を返します。
- 個別に指定する場合はカンマ区切りで入力します（例: `weight,diastolic,heart_rate`）。
- 利用可能なエイリアスは `internal/withings/config/aliases.go` に定義されています。数値 ID を直接指定することもできます。

> アクティビティ項目（歩数など）は `user.activity` スコープが付与されたアクセストークンで `-include-activity=true` の場合に `activity` ブロックとして出力されます。

## 主なフラグ

| フラグ | 説明 |
|--------|------|
| `-client-id`, `-client-secret` | Developer Portal で登録した OAuth アプリのクライアント情報。
| `-redirect-uri` | 認可コードを受け取る HTTPS エンドポイント。認可 URL 生成とコード交換の両方で同じ値を指定する必要があります。
| `-scope` | 認可 URL で要求するスコープ（例: `user.metrics,user.activity`）。複数はカンマ区切り。
| `-state` | 認可 URL に付与する任意文字列。CSRF 対策に使用し、コールバックで検証してください。
| `-authorization-code` | 認可コールバックのクエリパラメータ `code`。
| `-refresh-token` | 既に取得済みのリフレッシュトークン。`refresh-token` 操作で新しいアクセストークンを得る際に使用します。
| `-access-token` | API 呼び出しで使用するアクセストークン。`daily-summary` 操作で必須。
| `-user-id` | Withings から払い出される数値のユーザー ID。初回の `request-token` レスポンス内 `body.userid` でも取得できます。
| `-start-date` / `-end-date` | 日次データ取得の範囲（YYYY-MM-DD）。終了日を省略すると開始日のみ取得します。
| `-measure-types` | 取得したい測定タイプの別名または数値 ID（例: `weight`, `diastolic`, `heart_rate`）。
| `-include-activity` | 活動サマリを同時取得するか (`true`/`false`)。既定は `true`。
| `-timeout` | HTTP タイムアウト。既定は 15 秒。

## 代表的なスコープ

- `user.metrics` — 測定値（体重、血圧など）
- `user.activity` — 日次アクティビティ（歩数、消費カロリーなど）
- `user.info` — プロフィール情報

必要なスコープは開発する機能に応じて調整してください。Public Health Data API の詳細は Withings Developer Guide を参照してください。

## 認証フローの参考リンク

- [OAuth Web Flow (Public Health Data API)](https://developer.withings.com/developer-guide/v3/integration-guide/public-health-data-api/get-access/oauth-web-flow)
- [Authorization URL パラメータ詳細](https://developer.withings.com/developer-guide/v3/integration-guide/public-health-data-api/get-access/oauth-authorization-url)
- [Access / Refresh Token 取得](https://developer.withings.com/developer-guide/v3/get-access/access-and-refresh-tokens-no-recover)

## テスト

```bash
GOCACHE=$(pwd)/.gocache go test ./internal/withings/...
```

## 注意事項

- アクセストークンは約 3 時間、リフレッシュトークンは最長 1 年有効ですが、新しいリフレッシュトークンが発行されると旧トークンは 8 時間で失効します。常に最新値を安全に管理してください。
- 認可コードは 30 秒で失効するため、リダイレクトを受け取ったら即座に `request-token` を実行してください。
- `state` パラメータはコールバックで必ず検証し、不正なリクエストを排除してください。
