## CI実行ガイド
このディレクトリは、GitHub Actions とローカルで同じ Docker ベースのテスト環境を使うためのものです。
- `Dockerfile.ci`
  - CIテスト用コンテナ定義
- `ci_test.sh`
  - コンテナをビルドして `go test -v ./... -coverpkg=./... -covermode=count -coverprofile=...` を実行
- `ci_test_with_logging.sh`
  - `ci_test.sh` のログ保存と失敗時サマリー表示を担当

## 前提
- リポジトリルート: `$HOME/devbox`
- ローカル実行には Docker が必要

## ローカルでCI相当テストを実行する
リポジトリルートで実行:
```bash
bash ./scripts/ops/ci/ci_test_with_logging.sh \
  --log-file="tmp/go-test.log" \
  --go-version=1.25 \
  --cov-file=coverage.out \
  --coverage-report-file=tmp/coverage-report.txt \
  --image-tag=devbox-ci-test:local \
  --run-context=local
```

## GitHub Actionsで実行する
現在の `go-test-integration` は `on: push` トリガーです。
1. ブランチへ push する
2. GitHub の Actions で `go-test-integration` を開く
3. 失敗時は `Run Test` ステップを確認する

## Run Testで実際に呼ばれるコマンド
workflow `/.github/workflows/test_integration.yml` では以下の流れです。
```bash
bash ./scripts/ops/ci/ci_test_with_logging.sh \
  --log-file="${GITHUB_WORKSPACE}/tmp/go-test.log" \
  --go-version="${GO_VERSION}" \
  --cov-file="${COV_FILE}" \
  --coverage-report-file="${COV_REPORT_FILE}" \
  --image-tag="devbox-ci-test:${GITHUB_RUN_ID}-${GITHUB_RUN_ATTEMPT}" \
  --run-context="github-actions"
```

`--run-context` は `local` / `github-actions` / `auto` を指定できます。`local` の場合のみ、既知の失敗テスト名をバイネームで許容します。

## 失敗時ログの見方
`ci_test_with_logging.sh` は失敗時に目印を出します。
- `==================== FAILURE SUMMARY START ====================`
- `==================== FAILURE SUMMARY END ====================`
- `==================== FAILURE CONTEXT START ====================`
- `==================== FAILURE CONTEXT END ====================`
- `==================== LOCAL FAILURE FILTER START ====================`
- `==================== LOCAL FAILURE FILTER END ====================`
- `==================== FINAL RESULT START ====================`
- `==================== FINAL RESULT END ====================`

まず `FAILURE SUMMARY` を見て失敗テスト名や `panic` を確認し、次に `FAILURE CONTEXT` で前後行を確認してください。
ローカル実行で既知失敗フィルタが動いた場合は `filter-result` に `PASS local known failures only` または `FAIL unknown failures remain` が表示されます。
最終的な成否は `FINAL RESULT` の `overall-result` で判定できます。
さらに `COVERAGE TOTAL` セクションで、自モジュール `go list -m` 配下のみを対象にした総カバレッジ `coverage-total` が表示されます。  
この値の算出に使った `go tool cover -func` 形式のレポートは `--coverage-report-file` に保存され、coverage badge 更新でも同じファイルを使用します。
コンソール出力では `go test` のパッケージ別カバレッジ行を非表示にして、総カバレッジ表示のみを使うようにしています。

## よくある失敗

- `could not find default credentials`
  - GCP連携テストでADCが必要です。`google-github-actions/auth` 相当の認証情報をローカルにも設定する
- `no lines matched failure pattern`
  - テスト失敗以外の異常終了です。`--log-file` で指定したログ全文を確認する
