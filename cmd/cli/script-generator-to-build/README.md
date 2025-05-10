# Script Generator to Build

このツールは、指定されたGoパッケージのビルドスクリプトを生成するためのコマンドラインツールです。

## 機能

- 指定されたGoパッケージのビルドスクリプトを自動生成
- パッケージが指定されない場合は、利用可能なパッケージの一覧から選択可能
- 生成されたビルドスクリプトは、Linux/AMD64、Windows/AMD64、macOS/ARM64向けのクロスコンパイルを行う

## 使用方法

```
script-generator-to-build [パッケージ名]
```

### オプション

- `-h`, `--help`: ヘルプメッセージを表示

### 例

```
script-generator-to-build code-analyzer
script-generator-to-build
```

## 出力

指定されたパッケージに対応するビルドスクリプト（`build_<パッケージ名>.sh`）が生成されます。
