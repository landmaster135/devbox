# devbox
![Go](https://img.shields.io/badge/Go-1.23-%2300ADD8?logo=go)
![Coverage](https://img.shields.io/badge/Coverage-37.2%25-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

Provides utilities for development.

# Usage
- Go 1.23.5 or later

# Development

## Build
```bash
./scripts/build.sh
```

Confirm compilable distributions.
```bash
go tool dist list
```

## Project Structure

`devbox` is based on Clean Architecture. That package dependencies are the following.

```mermaid
graph TD
  A[cmd/file-processor/main.go] --> B[internal/interfaces/repositories]
  A --> C[internal/usecases/services]
  
  C[internal/usecases/services] --> D[internal/domain/repositories]
  
  B[internal/interfaces/repositories] --> D[internal/domain/repositories]
  B --> E[internal/domain/models]
  
  D[internal/domain/repositories] --> E[internal/domain/models]
  
  F[util] --> G[標準ライブラリ]
  
  A --> F[util]
  C --> F[util]
  B --> F[util]

  %% Style Settings
  classDef main fill:#f96,stroke:#333,stroke-width:2px;
  classDef domain fill:#bbf,stroke:#333,stroke-width:1px;
  classDef interfaces fill:#bfb,stroke:#333,stroke-width:1px;
  classDef usecases fill:#fbf,stroke:#333,stroke-width:1px;
  classDef util fill:#ddd,stroke:#333,stroke-width:1px;
  classDef stdlib fill:#eee,stroke:#333,stroke-width:1px,stroke-dasharray: 5 5;

  class A main;
  class E domain;
  class D domain;
  class B interfaces;
  class C usecases;
  class F util;
  class G stdlib;
```

### Package Overview

- **cmd/file-processor/main.go**: アプリケーションのエントリーポイント。コマンドライン引数の解析と依存関係の注入を行います。
- **internal/domain/models**: ドメインモデル（`FileContent`）の定義。ファイル内容の操作に関するビジネスロジックを実装しています。
- **internal/domain/repositories**: リポジトリインターフェース（`FileRepository`）の定義。ファイルの読み書きを抽象化します。
- **internal/interfaces/repositories**: リポジトリの実装（`FileRepositoryImpl`）。実際のファイルシステムとのやり取りを担当します。
- **internal/usecases/services**: ユースケース（`FileService`）の実装。ドメインモデルとリポジトリを組み合わせてビジネスロジックを実行します。
- **util**: ロギングやユーティリティ機能を提供します。アプリケーション全体で使用される共通機能です。

# License
MIT License
