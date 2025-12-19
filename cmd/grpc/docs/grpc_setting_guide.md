# gRPC設定ガイド

このガイドでは、protoファイルからGoコードを生成する方法と、よくある問題の解決方法について説明します。

## 目次

1. [基本的なコード生成](#基本的なコード生成)
2. [必要なツールの確認とインストール](#必要なツールの確認とインストール)
3. [トラブルシューティング](#トラブルシューティング)
4. [高度な設定](#高度な設定)

## 基本的なコード生成

最も基本的なprotoファイルからのコード生成コマンドです。

```bash
protoc --go_out=. --go-grpc_out=. proto/weather_notificator/weather_notificator.proto
```

## 必要なツールの確認とインストール

### バージョン確認

コード生成に失敗した場合は、まず必要なツールのバージョンを確認してください。

```bash
# バージョン確認
protoc-gen-go --version
protoc-gen-go-grpc --version
```

### パッケージのパス確認

必要なツールがインストールされているかを確認します。

```bash
# protoc-gen-goのパス確認
which protoc-gen-go || find $HOME -name "protoc-gen-go" 2>/dev/null || echo $GOPATH/bin/protoc-gen-go

# protoc-gen-go-grpcのパス確認
which protoc-gen-go-grpc || find $HOME -name "protoc-gen-go-grpc" 2>/dev/null || echo $GOPATH/bin/protoc-gen-go-grpc
```

### インストール

必要なツールをインストールします。

```bash
# Protocol Bufferコンパイラのインストール（protoc-gen-goも含まれます）
sudo apt install protobuf-compiler

# 不足している場合は追加でインストール
sudo apt install protoc-gen-go-grpc
```

## トラブルシューティング

### 絶対パスを指定したコード生成

PATHが通っていない等の原因でパッケージが見つからない場合は、絶対パスを指定してコード生成を実行します。

```bash
protoc --plugin=protoc-gen-go=/home/user/go/bin/protoc-gen-go \
       --go_out=. \
       --go-grpc_out=. \
       proto/weather_notificator/weather_notificator.proto
```

> **注意**: `/home/user/go/bin/protoc-gen-go`の部分は、実際の環境に合ったprotoc-gen-goのパスに置き換えてください。

### ディレクトリ構造の問題

Goプロジェクトのパスによっては、非常に深いディレクトリ構造になったり、プロジェクトディレクトリとの齟齬が発生してバグが起こることがあります。

## 高度な設定

### 相対パスを使用したコード生成

ディレクトリ構造の問題を解決するために、相対パスを指定してコード生成を行います。

```bash
# 相対パスを使用した基本的なコード生成
protoc --go_out=paths=source_relative:. \
       --go-grpc_out=paths=source_relative:. \
       proto/weather_notificator/weather_notificator.proto
```

### 絶対パス指定と相対パス出力の組み合わせ

最も確実な方法として、パッケージの絶対パス指定と相対パス出力を組み合わせます。

```bash
# 推奨: 絶対パス指定 + 相対パス出力
protoc --plugin=protoc-gen-go=/home/user/go/bin/protoc-gen-go \
       --go_out=paths=source_relative:. \
       --go-grpc_out=paths=source_relative:. \
       proto/weather_notificator/weather_notificator.proto
```

> **重要**: `/home/user/go/bin/protoc-gen-go`の部分は、実際の環境に合ったprotoc-gen-goのパスに置き換えてください。

このコマンドを使用することで、本来生成したいファイルが正しく出力されるはずです。

## まとめ

- 生成されるGoファイルは、protoファイルで定義されたpackage名に基づいてディレクトリが作成されます
- 問題が発生した場合は、まずバージョン確認とパス確認から始めることをお勧めします
