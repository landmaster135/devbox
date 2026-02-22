# coverage-badge

Goのカバレッジ結果から `img.shields.io` 形式のバッジMarkdownを生成し、任意のREADMEへ追加/更新するCLIです。

## フラグ一覧

| フラグ | 必須 | デフォルト | 説明 |
| :--- | :---: | :--- | :--- |
| `-operation` | 必須 | なし | `create-badge` または `patch-badge` |
| `-badge-title` | 任意 | `Coverage` | バッジ左側のタイトル |
| `-coverage-file` | 任意 | `coverage.out` | `go tool cover -func` 出力ファイル |
| `-green-threshold` | 任意 | `70` | 緑色へ切り替える閾値 |
| `-yellow-threshold` | 任意 | `30` | 黄色へ切り替える閾値 |
| `-force-color` | 任意 | なし | `green` / `yellow` / `red` を強制指定 |
| `-badge-link` | 任意 | なし | バッジクリック時のリンク先URL |
| `-badge-value` | 任意 | なし | カバレッジ値を直接指定（`58.6` / `58.6%`） |
| `-target-file` | `patch-badge` で実質必須 | `README.md` | 更新対象のMarkdownファイル |
| `-dry-run` | 任意 | `false` | `patch-badge` で書き込みせず更新後内容のみ出力 |
| `-help`, `-h` | 任意 | `false` | 使用方法を表示 |

## 使い方

```bash
go run ./cmd/cli/coverage-badge/main.go -operation=create-badge -coverage-file=coverage.out
```

```bash
go run ./cmd/cli/coverage-badge/main.go \
  -operation=create-badge \
  -badge-value=58.6 \
  -badge-title=Coverage \
  -badge-link=https://github.com/owner/repo/actions/workflows/test_integration.yml
```

```bash
go run ./cmd/cli/coverage-badge/main.go \
  -operation=patch-badge \
  -target-file=README.md \
  -coverage-file=coverage.out
```

```bash
go run ./cmd/cli/coverage-badge/main.go \
  -operation=patch-badge \
  -target-file=README.md \
  -badge-value=72.1 \
  -dry-run
```

## 出力例

成功時（`create-badge`）

```text
![Coverage](https://img.shields.io/badge/Coverage-58.6%25-yellow)
```

成功時（`patch-badge`）

```text
カバレッジバッジを更新しました: README.md
```

エラー時（不正なoperation）

```text
エラー: 設定の初期化に失敗しました: --operation は次のいずれかを指定してください: create-badge, patch-badge
```
