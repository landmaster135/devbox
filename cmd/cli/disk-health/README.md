# Disk Health

`smartctl -a` の出力を保存した `.log` / `.txt` ファイルから、ディスクの SMART 健康度を評価する CLI ツールです。

## 概要

`disk-health` は SMART 全体ヘルスと主要な故障予兆属性を解析し、`healthy` / `warning` / `critical` / `unknown` の 4 段階で評価します。特に代替待ちセクタ、訂正不能セクタ、代替済みセクタ、CRC エラー、温度を重点的に確認します。

## フラグ一覧

| フラグ | 必須 | デフォルト値 | 説明 |
|---|---|---|---|
| `-operation` | はい | なし | 実行する操作。対応値は `assess-smart` |
| `-src-file` | はい | なし | `smartctl -a` の出力を保存した `.log` / `.txt` ファイル |
| `-json` | いいえ | `false` | JSON 形式で出力 |
| `-verbose` | いいえ | `false` | 判定根拠に使った SMART 属性を詳細表示 |
| `-help` | いいえ | `false` | ヘルプを表示 |

## 使用方法

```bash
go run cmd/cli/disk-health/main.go -operation=assess-smart -src-file=smart.log
```

JSON で出力する場合:

```bash
go run cmd/cli/disk-health/main.go -operation=assess-smart -src-file=smart.log -json
```

詳細情報を含める場合:

```bash
go run cmd/cli/disk-health/main.go -operation=assess-smart -src-file=smart.log -verbose
```

## 使用例

```bash
smartctl -a /dev/sdb > smart.log
go run cmd/cli/disk-health/main.go -operation=assess-smart -src-file=smart.log
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
  "findings": [
    {
      "severity": "critical",
      "attribute_id": 197,
      "attribute_name": "Current_Pending_Sector",
      "raw_value": 404,
      "message": "代替待ちセクタを検出しました"
    }
  ]
}
```

### エラー時

```text
エラー: --src-file は必須パラメータです
使用方法: disk-health -operation=assess-smart -src-file=<SMARTログファイル> [オプション]
```

## 判定基準

| ステータス | 条件 |
|---|---|
| `healthy` | SMART 全体ヘルスが `PASSED` で、重大指標が検出されない |
| `warning` | 代替済みセクタ、代替イベント、CRC エラー、50 度以上の温度など注意指標がある |
| `critical` | SMART 全体ヘルス `FAILED`、代替待ちセクタ、訂正不能セクタ、属性しきい値割れ、60 度以上の温度などがある |
| `unknown` | SMART 全体情報または属性表が不足している |

## テスト

```bash
go test ./cmd/cli/disk-health ./internal/disk_health/...
```
