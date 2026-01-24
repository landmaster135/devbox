# Machine Info

Ubuntu系LinuxでPCのハードウェア／ネットワーク情報を収集し、JSONとして保存するCLIツールです。

## 主な機能
- CPU名、コア数、論理プロセッサ数、クロック情報、温度を取得
- メモリ容量と使用量、ホスト名を取得
- 指定インターフェースの送受信スループットを複数サンプリングして平均化
- 取得した情報を標準出力へ表示し、`log_<timestamp>.json` に保存
- `--output-dir` でログの保存先ディレクトリを指定可能
- `--operation` フラグで操作モードを切り替え（現在は`ubuntu`のみ）

## ビルド
```bash
# プロジェクトルートから
./scripts/build_machine_info.sh  # もしくは
go build -o bin/machine-info ./cmd/cli/machine-info
```

## 使い方
### 基本コマンド
```bash
go run ./cmd/cli/machine-info --operation=ubuntu --network-interface=eth0
```

実行後、標準出力に計測結果とJSONが表示され、指定（未指定ならカレントディレクトリ）したフォルダに `log_YYYYmmdd-HHMMSS.json` が作成されます。

### フラグ一覧
| フラグ | 説明 | デフォルト | 備考 |
|--------|------|------------|------|
| `--operation` | 実行する操作モード。現在は`ubuntu`のみサポート | `ubuntu` | 今後の拡張用に追加 |
| `--network-interface` | ネットワーク速度を計測するインターフェース名 | `eth0` |  |
| `--output-dir` | JSONログを保存するディレクトリ | カレント | 不存在の場合は自動作成 |
| `--help`, `-h` | 使い方を表示 | `false` | - |

### Operation一覧
| operation | 説明 |
|-----------|------|
| `ubuntu` | Ubuntu系OSで`lscpu`や`lm-sensors`等の標準ユーティリティを用いて各種情報を取得 |

### 実行例
```
$ go run ./cmd/cli/machine-info --operation=ubuntu --network-interface=enp0s3
マシン情報を取得中...
CPU名: AMD Ryzen 7 7840HS w/ Radeon 780M Graphics
CPUコア数: 8
論理プロセッサ数: 16
CPU最大クロック速度: 5100.00 MHz
CPU現在のクロック速度: 3050.00 MHz
CPU温度: 54.20 °C
メモリ総容量: 32000.00 MB
メモリ使用量: 10240.00 MB
ホスト名: devbox
ネットワーク速度の計測結果:
平均送信速度: 1200.50 Kbps
平均受信速度: 3400.20 Kbps

取得したシステム情報JSON:
{
  "cpu_name": "AMD Ryzen 7 7840HS w/ Radeon 780M Graphics",
  ...
}

ログファイルに保存しました: log_20250101-153000.json
```
※ 実際の値は環境によって異なります。
