# JSON File Merger

JSONファイルを読み込み、APIリクエスト用のリクエストボディを作成するツールです。

## 機能

1. 指定されたディレクトリからJSONファイルを読み込みます
2. 読み込んだJSONデータを指定されたキー名の配列に格納します
3. 以下の形式のリクエストボディを作成します：
   ```json
   {
     "data": {
       "[指定されたキー名]": [
         // JSONファイルから読み込んだデータの配列
       ]
     },
     "description": "By manually bulk request",
     "name": "manual_request"
   }
   ```
4. 作成したリクエストボディをJSONファイルとして保存します（オプション）

## 使用方法

```bash
json-file-merger -dir <JSONディレクトリパス> [-key <キー名>] [-output <出力ファイルパス>]
```

### オプション

- `-dir`: JSONファイルが格納されているディレクトリのパス（必須）
- `-key`: JSONデータの配列が入るキーの名前（デフォルト: "pc_stats"）
- `-output`: 作成したリクエストボディを保存するファイルのパス（省略可）

### 例

```bash
# sample_dataディレクトリ内のJSONファイルを使用して、キー名「pc_stats」のリクエストボディを作成
json-file-merger -dir ../../sample_data -key pc_stats -output output.json
```

## ビルド方法

```bash
cd /path/to/devbox
go build -o bin/json-file-merger ./cmd/json-file-merger
```

## 実行方法

ビルド後：

```bash
./bin/json-file-merger -dir ./sample_data -key pc_stats -output output.json
```

ビルドせずに実行：

```bash
go run ./cmd/json-file-merger/main.go -dir ./sample_data -key pc_stats -output output.json
