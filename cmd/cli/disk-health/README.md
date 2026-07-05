# Disk Health

`smartctl -a` または dmesg の出力を保存した `.log` / `.txt` ファイルから、ディスクの健康度を評価する CLI ツールです。

## 概要

`disk-health` は SMART 全体ヘルス、主要な故障予兆属性、dmesg に出力されたディスク I/O エラーを解析し、`healthy` / `warning` / `critical` / `unknown` の 4 段階で評価します。SMART では代替待ちセクタ、訂正不能セクタ、代替済みセクタ、CRC エラー、温度を重点的に確認します。dmesg では `critical medium error`、`Medium Error`、`Unrecovered read error`、ディスク文脈の `I/O error` などを確認します。

健康度評価に加えて、取得できる場合は Rotation Rate、電源投入時間、電源投入回数、温度、読み込み量、書き込み量も出力します。この追加サマリーは `--verbose` を指定しない通常出力にも含まれます。

## フラグ一覧

| フラグ | 必須 | デフォルト値 | 説明 |
|---|---|---|---|
| `--operation` | はい | なし | 実行する操作。対応値は `assess-smart`, `assess-dmesg` |
| `--src-file` | はい | なし | `smartctl -a` または dmesg の出力を保存した `.log` / `.txt` ファイル |
| `--json` | いいえ | `false` | JSON 形式で出力 |
| `--verbose` | いいえ | `false` | 判定根拠に使った SMART 属性または dmesg 行を詳細表示 |
| `--help` | いいえ | `false` | ヘルプを表示 |

## 使用方法

```bash
go run cmd/cli/disk-health/main.go --operation=assess-smart --src-file=smart.log
```

dmesg ログを評価する場合:

```bash
go run cmd/cli/disk-health/main.go --operation=assess-dmesg --src-file=dmesg-disk.log --json
```

## 使用例

```bash
smartctl -a /dev/sdb > smart.log
go run cmd/cli/disk-health/main.go --operation=assess-smart --src-file=smart.log
```

```bash
sudo dmesg | grep -E 'error|I/O|sdb' > dmesg-disk.log
go run cmd/cli/disk-health/main.go --operation=assess-dmesg --src-file=dmesg-disk.log
```

## 出力例

### 成功時

```text
status: healthy
score: 100
summary: SMART情報に重大な問題は見つかりませんでした。
model: ST5000LM000-2AN170
serial_number: WCJ925F8
overall_health: PASSED
disk_info:
  rotation_rate_rpm: 5526
  power_on_hours: 87
  power_cycle_count: 159
  temperature_celsius: 30
  total_lbas_written: 3836222920
  total_bytes_written: 1964146135040
  total_lbas_read: 1647232405
  total_bytes_read: 843382991360
```

### 致命的な欠損を検出した場合

```text
status: critical
score: 20
summary: 重大な SMART 欠損指標を検出しました。速やかなバックアップと交換を推奨します。
model: WDC WD50NDZM-11BCXS1
serial_number: WD-WX12D126617N
overall_health: PASSED
findings:
- [critical] Current_Pending_Sector raw=404: 代替待ちセクタを検出しました
- [critical] Reallocated_Sector_Ct raw=0: Current_Pending_Sector=404 による代替セクタへの移し替えに失敗、ドライブが自己修復できていません
```

### JSON 出力

```json
{
  "status": "critical",
  "score": 20,
  "summary": "重大な SMART 欠損指標を検出しました。速やかなバックアップと交換を推奨します。",
  "overall_health": "PASSED",
  "model": "WDC WD50NDZM-11BCXS1",
  "serial_number": "WD-WX12D126617N",
  "disk_info": {
    "rotation_rate_rpm": 5526,
    "power_on_hours": 87,
    "power_cycle_count": 159,
    "temperature_celsius": 30,
    "total_lbas_written": 3836222920,
    "total_bytes_written": 1964146135040,
    "total_lbas_read": 1647232405,
    "total_bytes_read": 843382991360
  },
  "findings": [
    {
      "severity": "critical",
      "attribute_id": 197,
      "attribute_name": "Current_Pending_Sector",
      "raw_value": 404,
      "message": "代替待ちセクタを検出しました"
    },
    {
      "severity": "critical",
      "attribute_id": 5,
      "attribute_name": "Reallocated_Sector_Ct",
      "raw_value": 0,
      "message": "Current_Pending_Sector=404 による代替セクタへの移し替えに失敗、ドライブが自己修復できていません"
    }
  ]
}
```

### dmesg で致命的な欠損を検出した場合

```text
status: critical
score: 20
summary: dmesgログからディスクI/Oの重大エラーを検出しました。速やかなバックアップと交換を推奨します。
findings:
- [critical] sdb: critical medium error を検出しました
- [warning] sdb: FAILED Result を検出しました
- [critical] sdb: Medium Error を検出しました
- [critical] sdb: Unrecovered read error を検出しました
```

### dmesg の JSON 出力

```json
{
  "status": "critical",
  "score": 20,
  "summary": "dmesgログからディスクI/Oの重大エラーを検出しました。速やかなバックアップと交換を推奨します。",
  "findings": [
    {
      "severity": "critical",
      "device": "sdb",
      "message": "critical medium error を検出しました"
    },
    {
      "severity": "warning",
      "device": "sdb",
      "message": "FAILED Result を検出しました"
    },
    {
      "severity": "critical",
      "device": "sdb",
      "message": "Medium Error を検出しました"
    },
    {
      "severity": "critical",
      "device": "sdb",
      "message": "Unrecovered read error を検出しました"
    }
  ]
}
```

### エラー時

```text
エラー: --src-file は必須パラメータです
使用方法: disk-health --operation=assess-smart --src-file=<SMARTログファイル> [オプション]
```

## 判定基準

| ステータス | 条件 |
|---|---|
| `healthy` | SMART 全体ヘルスが `PASSED` で、重大指標が検出されない |
| `warning` | 代替済みセクタ、代替イベント、CRC エラー、50 度以上の温度など注意指標がある |
| `critical` | SMART 全体ヘルス `FAILED`、代替待ちセクタ、訂正不能セクタ、属性しきい値割れ、60 度以上の温度などがある |
| `unknown` | SMART 全体情報または属性表が不足している |

### dmesg 判定基準

| ステータス | 条件 |
|---|---|
| `healthy` | ディスク故障に関連する dmesg イベントが検出されない。ACPI BIOS Error や I/O scheduler など、ディスク故障と無関係な行は対象外 |
| `warning` | ディスクデバイスに紐づく `I/O error`、`FAILED Result`、`failed command`、`timeout` など注意イベントがある |
| `critical` | `critical medium error`、`Medium Error`、`Unrecovered read error` など媒体エラーや回復不能な読み取りエラーがある |
| `unknown` | 入力ログが空で評価できない |

## テスト

```bash
go test ./cmd/cli/disk-health ./internal/disk_health/...
```
