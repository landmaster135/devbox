---
name: go-version-update
description: Go のバージョン更新手順を、OS と既存導入経路に合わせて安全に案内する。Ubuntu の package 管理運用と手動インストール運用の両方に対応し、PATH 競合や GOTOOLCHAIN 設定を含めて実施・検証・復旧まで提示する。
---

# Go Version Update

Go のバージョンを上げるときに、環境破壊を避けつつ最短で更新するための実務手順。

## このスキルを使う場面

- `go mod tidy` や `go test` が「Go バージョン不足」で失敗する
- 特定バージョン（例: `1.25.8`）に上げたい
- `/usr/bin/go` と `/usr/local/go/bin/go` の競合を解消したい
- システム Go を変えずに必要バージョンだけ使いたい

## 事前確認

必ず最初に現状を確認する。

```bash
which go
readlink -f "$(which go)"
go version
go env GOROOT GOPATH GOTOOLCHAIN
echo "$PATH"
```

## 判断フロー

1. 要件が「今のプロジェクトだけ動かしたい」なら: `GOTOOLCHAIN` を使う
2. 要件が「システム全体で go 自体を更新したい」なら: 既存導入経路に合わせて更新
3. `/usr/bin/go` と `/usr/local/go/bin/go` が混在しているなら: どちらかに一本化

## アプローチA: システムGoを維持して必要バージョンを使う（推奨）

システムの配置を変えず、必要ツールチェーンを使う。

```bash
go env -w GOTOOLCHAIN=go<VERSION>+auto
go version
```

例:

```bash
go env -w GOTOOLCHAIN=go1.25.8+auto
```

この方式は Ubuntu の `/usr/bin/go` 運用と相性がよい。

## アプローチB: UbuntuでシステムGoそのものを更新する

### B-1. apt 管理を維持する

- `apt` で提供されるバージョンに依存する
- 狙ったパッチ版が必要な場合は不向き

```bash
sudo apt update
sudo apt install --only-upgrade golang-go
```

### B-2. 公式 tarball で `/usr/local/go` 運用に切り替える

- 狙ったバージョンを正確に入れたい場合に有効
- 既存 `/usr/local/go` は入れ替え前提

```bash
cd /tmp
wget https://go.dev/dl/go<VERSION>.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go<VERSION>.linux-amd64.tar.gz
```

PATH を先頭に通す。

```bash
echo 'export PATH=/usr/local/go/bin:$PATH' >> ~/.profile
source ~/.profile
hash -r
```

## 競合解消ルール

- `/usr/bin/go` を使う方針なら `/usr/local/go` を削除する
- `/usr/local/go` を使う方針なら PATH で `/usr/local/go/bin` を優先
- 両方を残して順序未管理のままにしない

確認:

```bash
which go
readlink -f "$(which go)"
go version
```

## 実施後の検証

1. バージョン確認
```bash
go version
go env GOTOOLCHAIN
```

2. 対象プロジェクトで再実行
```bash
cd <PROJECT_DIR>
go mod tidy
go test ./...
```

## 失敗時の切り分け

- `requires go >= X.Y.Z`: Goバージョン不足。`GOTOOLCHAIN` か Go 本体更新
- `proxy.golang.org` 到達失敗: ネットワーク/DNS/プロキシ確認
- `read-only file system`（`~/.cache/go-build` 等）: 実行環境の権限制約
- `go env -w` 失敗: `~/.config/go/env` への書き込み権限確認

## 応答テンプレート

ユーザーには次の順で簡潔に返す。

1. 現状（`which go` / `go version`）
2. 原因（バージョン不足 or PATH競合 or 権限）
3. 実行コマンド（コピー可能）
4. 検証コマンド
5. 失敗時の次アクション

## 注意

- 破壊的操作（例: `rm -rf /usr/local/go`）前には必ず意図を確認する
- 既存の運用方針（apt管理 or 手動管理）を尊重する
- バージョンは固定値で断定せず、要件に合わせて `<VERSION>` を埋める
