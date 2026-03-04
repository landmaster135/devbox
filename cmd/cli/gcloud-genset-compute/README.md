# Gcloud Genset Compute CLI

Google Compute Engine 向けの `gcloud` コマンドを生成する CLI ツールです。コマンドは実行せず、標準出力へ表示します。

## サポートする操作

| operation | 説明 |
|---|---|
| `create-gce-instance` | VM インスタンス作成コマンドを生成 |
| `create-gce-router-and-nat` | Cloud Router 作成 + Cloud NAT 作成コマンドを生成 |
| `create-gce-iap-ssh-firewall-rule` | IAP SSH 用 firewall rule 作成コマンドを生成 |
| `create-gce-ingress-ssh-firewall-rule` | VPC 内 SSH 用 firewall rule 作成コマンドを生成 |
| `list-gcloud-instances` | インスタンス一覧取得コマンドを生成 |

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
