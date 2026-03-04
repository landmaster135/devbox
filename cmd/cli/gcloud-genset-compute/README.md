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

### connect-gce-instance

| フラグ | 必須 | デフォルト | 説明 |
|---|---|---|---|
| `-instance-name` | 必須 | なし | SSH 接続対象のインスタンス名 |
| `-zone` | 任意 | `us-central1-a` | ゾーン |

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

firewall rule 名には実行時刻のサフィックス `YYMMDD-hhmmss` が付与されます。

## 使用例

### VM 作成コマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=create-gce-instance \
  -instance-name=my-vm \
  -zone=us-central1-a \
  -machine-type=e2-medium
```

出力例:

```bash
==============================
生成された gcloud コマンド
==============================
gcloud compute instances create 'my-vm' --zone='us-central1-a' --machine-type='e2-medium' --no-address --boot-disk-size='100GB' --boot-disk-type='pd-balanced'
==============================
```

### インスタンス一覧コマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=list-gcloud-instances \
  -filter='zone:us-central1-a'
```

出力例:

```bash
==============================
生成された gcloud コマンド
==============================
gcloud compute instances list --filter='zone:us-central1-a' --format='table(name, zone.basename(), scheduling.preemptible.yesno(yes=true, no='"'"''), networkInterfaces.internal_ip():label=INTERNAL_IP, external_ip():label=EXTERNAL_IP, status)'
==============================
```

### ディスクタイプ一覧取得（サイズ条件付き）

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=list-disk-types \
  -zones=asia-southeast3-a \
  -min-disk-size-gib=4 \
  -max-disk-size-gib=65536
```

出力例:

```text
NAME                  ZONE               VALID_DISK_SIZES
hyperdisk-balanced    asia-southeast3-a  4GB-65536GB
hyperdisk-extreme     asia-southeast3-a  64GB-65536GB
hyperdisk-ml          asia-southeast3-a  4GB-65536GB
hyperdisk-throughput  asia-southeast3-a  2048GB-32768GB
local-ssd             asia-southeast3-a  375GB-375GB
pd-balanced           asia-southeast3-a  10GB-65536GB
pd-extreme            asia-southeast3-a  500GB-65536GB
pd-ssd                asia-southeast3-a  10GB-65536GB
pd-standard           asia-southeast3-a  10GB-65536GB
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

出力例:

```text
NAME         ZONE               GUEST_CPUS  MEMORY_MB  MAX_PERSISTENT_DISKS_SIZE_GB
c3-highmem-4 asia-southeast3-a  4           32768      65536
```

### インスタンス起動コマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=start-gce-instance \
  -instance-name=my-vm \
  -zone=us-central1-a
```

出力例:

```bash
==============================
生成された gcloud コマンド
==============================
gcloud compute instances start 'my-vm' --zone='us-central1-a'
==============================
```

### インスタンス削除コマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=delete-gce-instance \
  -instance-name=my-vm \
  -zone=us-central1-a
```

出力例:

```bash
==============================
生成された gcloud コマンド
==============================
gcloud compute instances delete 'my-vm' --zone='us-central1-a' --quiet
==============================
```

### firewall作成 + SSHセットアップコマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=setup-gce-firewall-and-ssh \
  -instance-name=my-vm \
  -zone=us-central1-a
```

出力例:

```bash
==============================
生成された gcloud コマンド
==============================
if [ -z "${SSH_AUTH_SOCK:-}" ]; then eval "$(ssh-agent -s)" >/dev/null; fi && ssh-add "$HOME/.ssh/google_compute_engine" && \
gcloud compute firewall-rules create 'allow-ssh-ingress-from-iap-260304-131645' --direction='INGRESS' --action='allow' --rules='tcp:22' --source-ranges='35.235.240.0/20' && \
gcloud compute firewall-rules create 'allow-ingress-ssh-260304-131645' --allow='tcp:22' --source-ranges='10.0.0.0/8' && \
gcloud compute scp "$HOME/.ssh/google_compute_engine" 'my-vm:/tmp' --zone='us-central1-a' --tunnel-through-iap && \
gcloud compute ssh 'my-vm' --zone='us-central1-a' --tunnel-through-iap
==============================
```

### インスタンス作成 + metadata + startup-script コマンド生成

```bash
go run ./cmd/cli/gcloud-genset-compute \
  -operation=create-gce-instance-and-configure \
  -instance-name=my-vm
```

出力例:

```bash
==============================
生成された gcloud コマンド
==============================
gcloud compute instances create 'my-vm' --zone='us-central1-a' --machine-type='e2-medium' --no-address --boot-disk-size='100GB' --boot-disk-type='pd-balanced' && \
gcloud compute instances add-metadata 'my-vm' --zone='us-central1-a' --metadata='VSC_PROFILE_URL="",KEYMAP_PROFILE_URL="",GITHUB_ACCOUNT_NAME="",GITHUB_ACCOUNT_EMAIL="",DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD="",GCE_ICON_URL="",GSM_ICON_URL="",GCS_ICON_URL="",GCSCHEDULER_ICON_URL="",GCIAM_ICON_URL="",GCLOUD_RUN_ICON_URL="",GCLOUD_RUN_FUNCTION_ICON_URL="",DEV_HOME=""' && \
gcloud compute instances add-metadata 'my-vm' --zone='us-central1-a' --metadata-from-file startup-script='cmd/cli/gcloud-genset-compute/metadata/setup_scripts/startup-script.sh'
==============================
```

## 備考

### startup-script 処理時間

下記は実際の通知履歴を元に算出した実測値です。目安としてください。既にstartup_scriptによる構築が終わっている場合、全工程を通して1分ぐらいで完了します。

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

## エラー例

```bash
go run ./cmd/cli/gcloud-genset-compute -operation=create-gce-instance
```

```text
エラー: instance-name は必須です
```

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
