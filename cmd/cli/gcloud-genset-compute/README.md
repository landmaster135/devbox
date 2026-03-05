# Gcloud Genset Compute CLI

Google Compute Engine 向けの CLI ツールです。  
基本操作は `gcloud` コマンドを生成して表示し、`list-disk-types` / `list-machine-types` は `gcloud` を実行して結果を直接表示します。

## サポートする操作

| operation | 説明 |
|---|---|
| `create-gce-instance` | VM インスタンス作成コマンドを生成 |
| `create-gce-instance-with-startup-script` | VM インスタンス作成 + startup-script 設定コマンドを生成 |
| `create-gce-instance-and-configure` | VM インスタンス作成 + metadata 設定 + startup-script 設定コマンドを生成 |
| `create-gce-router-and-nat` | Cloud Router 作成 + Cloud NAT 作成コマンドを生成 |
| `create-gce-iap-ssh-firewall-rule` | IAP SSH 用 firewall rule 作成コマンドを生成 |
| `create-gce-ingress-ssh-firewall-rule` | VPC 内 SSH 用 firewall rule 作成コマンドを生成 |
| `list-gcloud-instances` | インスタンス一覧取得コマンドを生成 |
| `list-disk-types` | ディスクタイプ一覧を取得して表示（サイズ条件付き絞り込み対応） |
| `list-machine-types` | マシンタイプ一覧を取得して表示（最大永続ディスクサイズ/メモリサイズ条件付き絞り込み対応） |
| `start-gce-instance` | インスタンス起動コマンドを生成 |
| `stop-gce-instance` | インスタンス停止コマンドを生成 |
| `reboot-gce-instance` | インスタンス再起動コマンドを生成 |
| `delete-gce-instance` | インスタンス削除コマンドを生成 |
| `copy-gce-ssh-key` | SSH 秘密鍵を `/tmp` へコピーするコマンドを生成 |
| `scp-dir` | ローカルディレクトリを再帰的にインスタンスへコピーするコマンドを生成 |
| `connect-gce-instance` | IAP トンネル経由の SSH 接続コマンドを生成 |
| `set-gce-instance-metadata-from-yaml` | YAML ファイルをもとに instance metadata 設定コマンドを生成 |
| `add-startup-script-to-gce-instance` | startup-script を instance metadata に設定するコマンドを生成 |
| `setup-gce-firewall-and-ssh` | firewall 作成 + SSH 鍵コピー + SSH 接続の複合コマンドを生成 |

## フラグ一覧

### 共通

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-operation` | 必須 | なし | 実行する操作 |
| `-help` | 任意 | `false` | ヘルプ表示 |

### create-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | インスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-machine-type` | 任意 | `e2-medium` | マシンタイプ |
| `-boot-disk-size` | 任意 | `100GB` | ブートディスクサイズ |
| `-boot-disk-type` | 任意 | `pd-balanced` | ブートディスクタイプ |

### create-gce-instance-with-startup-script

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | インスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-machine-type` | 任意 | `e2-medium` | マシンタイプ |
| `-boot-disk-size` | 任意 | `100GB` | ブートディスクサイズ |
| `-boot-disk-type` | 任意 | `pd-balanced` | ブートディスクタイプ |
| `-startup-script-path` | 任意 | `cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh` | startup-script のファイルパス |

### create-gce-instance-and-configure

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | インスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-machine-type` | 任意 | `e2-medium` | マシンタイプ |
| `-boot-disk-size` | 任意 | `100GB` | ブートディスクサイズ |
| `-boot-disk-type` | 任意 | `pd-balanced` | ブートディスクタイプ |
| `-metadata-yaml-path` | 任意 | `cmd/cli/gcloud-genset-compute/metadata/config/env.yml` | metadata 用 YAML ファイルパス |
| `-startup-script-path` | 任意 | `cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh` | startup-script のファイルパス |

### create-gce-router-and-nat

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-router-name` | 必須 | なし | Cloud Router 名 |
| `-region` | 任意 | `us-central1` | リージョン |
| `-network` | 任意 | `default` | VPC ネットワーク |
| `-nat-name` | 任意 | `nat1` | Cloud NAT 名 |

### create-gce-iap-ssh-firewall-rule

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-rule-name` | 任意 | `allow-ssh-ingress-from-iap` | ルール名 |
| `-direction` | 任意 | `INGRESS` | ルール方向 |
| `-action` | 任意 | `allow` | アクション |
| `-rules` | 任意 | `tcp:22` | 許可ルール |
| `-source-ranges` | 任意 | `35.235.240.0/20` | 送信元 CIDR |

### create-gce-ingress-ssh-firewall-rule

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-rule-name` | 任意 | `allow-ingress-ssh` | ルール名 |
| `-allow-rule` | 任意 | `tcp:22` | `--allow` の値 |
| `-source-ranges` | 任意 | `10.0.0.0/8` | 送信元 CIDR |

### list-gcloud-instances

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-filter` | 任意 | なし | 一覧絞り込み条件 |
| `-format` | 任意 | table形式 | 出力形式 |

### list-disk-types

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-zones` | 任意 | なし | 対象ゾーン（カンマ区切り） |
| `-min-disk-size-gib` | 任意 | `0` | 必要な最小ディスクサイズ（GiB） |
| `-max-disk-size-gib` | 任意 | `0` | 必要な最大ディスクサイズ（GiB） |

`-min-disk-size-gib` と `-max-disk-size-gib` を同時指定した場合、指定レンジ全体を扱える disk type のみを抽出するコマンドを生成します。

### list-machine-types

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-zones` | 任意 | なし | 対象ゾーン（カンマ区切り） |
| `-min-memory-size-mib` | 任意 | `0` | `memoryMb` の下限 |
| `-max-memory-size-mib` | 任意 | `0` | `memoryMb` の上限 |
| `-min-disk-size-gib` | 任意 | `0` | `maximumPersistentDisksSizeGb` の下限 |
| `-max-disk-size-gib` | 任意 | `0` | `maximumPersistentDisksSizeGb` の上限 |

### start-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | 起動するインスタンス名 |
| `-zone` | 必須 | なし | ゾーン |

### stop-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | 停止するインスタンス名 |
| `-zone` | 必須 | なし | ゾーン |

### reboot-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | 再起動するインスタンス名 |
| `-zone` | 必須 | なし | ゾーン |

### delete-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | 削除するインスタンス名 |
| `-zone` | 必須 | なし | ゾーン |

### copy-gce-ssh-key

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | 鍵コピー対象のインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-ssh-key-path` | 任意 | `$HOME/.ssh/google_compute_engine` | コピーする秘密鍵パス |
| `-creates-ssh-key` | 任意 | `false` | `true` のとき SSH 秘密鍵を新規作成してから処理を続行（`ssh-keygen` の対話入力でパスフレーズ設定可） |
| `-forces` | 任意 | `false` | `creates-ssh-key=true` かつ既存鍵がある場合に上書き許可（上書き時は標準エラーへログ出力） |

### scp-dir

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | コピー先のインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-src-dir` | 必須 | なし | ローカルのコピー元ディレクトリ（存在必須） |
| `-dest-dir` | 必須 | なし | インスタンス上のコピー先ディレクトリパス |

### connect-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | SSH 接続対象のインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-ssh-key-path` | 任意 | `$HOME/.ssh/google_compute_engine` | `ssh-add` および鍵生成時に使用する秘密鍵パス |
| `-creates-ssh-key` | 任意 | `false` | `true` のとき SSH 秘密鍵を新規作成してから接続（`ssh-keygen` の対話入力でパスフレーズ設定可） |
| `-forces` | 任意 | `false` | `creates-ssh-key=true` かつ既存鍵がある場合に上書き許可（上書き時は標準エラーへログ出力） |

### set-gce-instance-metadata-from-yaml

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | metadata を設定するインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-metadata-yaml-path` | 任意 | `cmd/cli/gcloud-genset-compute/metadata/config/env.yml` | metadata 用 YAML ファイルパス |

### add-startup-script-to-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | startup-script を設定するインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-startup-script-path` | 任意 | `cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh` | startup-script のファイルパス |

### setup-gce-firewall-and-ssh

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | セットアップ対象のインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |
| `-ssh-key-path` | 任意 | `$HOME/.ssh/google_compute_engine` | コピーする秘密鍵パス |
| `-creates-ssh-key` | 任意 | `false` | `true` のとき SSH 秘密鍵を新規作成してから処理を続行（`ssh-keygen` の対話入力でパスフレーズ設定可） |
| `-forces` | 任意 | `false` | `creates-ssh-key=true` かつ既存鍵がある場合に上書き許可（上書き時は標準エラーへログ出力） |

firewall rule 名には実行時刻のサフィックス `YYMMDD-hhmmss` が付与されます。
`creates-ssh-key=true` で鍵が既存の場合、`forces=false` だとエラーで終了します。

## 使用例

### 一般的なユースケース
```bash
# Create VM instance
go run ./cmd/cli/gcloud-genset-compute --operation=create-gce-instance-and-configure --instance-name=my-vm00 --machine-type=e2-standard-2

# If the startup script is NOT executed. Reboot VM instance (Depends on timing)
go run ./cmd/cli/gcloud-genset-compute --operation=reboot-gce-instance --instance-name=my-vm00 --zone=us-central1-a

# Connect VM instance with NOT overwritten SSH key
go run ./cmd/cli/gcloud-genset-compute --operation=setup-gce-firewall-and-ssh --instance-name=my-vm00 --zone=us-central1-a --creates-ssh-key=false

### Process for cleaning up ###

# Stop VM instance
go run ./cmd/cli/gcloud-genset-compute --operation=stop-gce-instance --instance-name=my-vm00 --zone=us-central1-a

# Delete VM instance
go run ./cmd/cli/gcloud-genset-compute --operation=delete-gce-instance --instance-name=my-vm00 --zone=us-central1-a
```

### インスタンス一覧コマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=list-gcloud-instances \
  -filter='zone:us-central1-a'
```

### ディスクタイプ一覧取得（サイズ条件付き）

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=list-disk-types \
  -zones=asia-southeast3-a \
  -min-disk-size-gib=4 \
  -max-disk-size-gib=65536
```

### マシンタイプ一覧取得（サイズ条件付き）

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=list-machine-types \
  -zones=asia-southeast3-a,asia-southeast3-b \
  -min-memory-size-mib=8192 \
  -max-memory-size-mib=65536 \
  -min-disk-size-gib=1024 \
  -max-disk-size-gib=524288
```

### firewall作成 + SSHセットアップ（鍵を再生成して上書き許可）

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=setup-gce-firewall-and-ssh \
  -instance-name=my-vm \
  -zone=us-central1-a \
  -creates-ssh-key=true \
  -forces=true
```

### SSH接続コマンド生成（IAPトンネル経由）

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=connect-gce-instance \
  -instance-name=my-vm \
  -zone=us-central1-a \
  -creates-ssh-key=false
```

### ディレクトリ再帰コピーコマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=scp-dir \
  -instance-name=my-vm \
  -zone=us-central1-a \
  -src-dir=./local_workspace \
  -dest-dir=remote_workspace
```

## 備考

### startup-script 処理時間

下記は実際の通知履歴 (`--machine-type=e2-medium`) を元に算出した実測値です。目安としてください。
既にstartup_scriptによる構築が終わっている場合、全工程を通して1分ぐらいで完了します。

| 対応処理 | 各工程の所要時間 |
|---|---|
| カスタム metadata 反映 | 0分 |
| startup-script 開始通知 | 0分 |
| デスクトップ環境設定 | 7分 |
| Chrome Remote Desktop 設定完了通知 | 0分 |
| ロケール設定 | 1分 |
| タイムゾーン設定 | 0分 |
| IME 設定 | 1分 |
| 開発リソース設定 | 0分 |
| Docker 設定 | 0分 |
| VSCode 設定 | 1分 |
| startup-script 全体完了通知 | 0分 |

## ビルド

```bash
go build -o bin/gcloud-genset-compute ./cmd/cli/gcloud-genset-compute
```

## 参考
References:
- [gcloud compute instances create](https://docs.cloud.google.com/sdk/gcloud/reference/compute/instances/create)
- [gcloud compute instances add-metadata](https://docs.cloud.google.com/sdk/gcloud/reference/compute/instances/add-metadata)
- [gcloud compute instances list](https://docs.cloud.google.com/sdk/gcloud/reference/compute/instances/list)
- [gcloud compute instances start](https://docs.cloud.google.com/sdk/gcloud/reference/compute/instances/start)
- [gcloud compute instances delete](https://docs.cloud.google.com/sdk/gcloud/reference/compute/instances/delete)
- [gcloud compute firewall-rules create](https://docs.cloud.google.com/sdk/gcloud/reference/compute/firewall-rules/create)
- [gcloud compute scp](https://docs.cloud.google.com/sdk/gcloud/reference/compute/scp)
- [gcloud compute ssh](https://docs.cloud.google.com/sdk/gcloud/reference/compute/ssh)
- [gcloud compute disk-types list](https://docs.cloud.google.com/sdk/gcloud/reference/compute/disk-types/list)
- [gcloud compute machine-types list](https://docs.cloud.google.com/sdk/gcloud/reference/compute/machine-types/list)
