# devbox
![Go](https://img.shields.io/badge/Go-1.25-%2300ADD8?logo=go)
![Coverage](https://img.shields.io/badge/Coverage-58.7%25-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

Provides utilities for development.

# Usage
- Go 1.25.5 or later

# Development

## Generate Go packages
```bash
./scripts/create_project_files.sh <PACKAGE_NAME>
```

## Generate shell scripts to build
```bash
cd devbox
./pkg/bin/cli/linux_amd64/script-generator-to-build <TOOL_NAME>
```

# Build

## Common Tools
```bash
./scripts/build.sh
```

Confirm compilable distributions.
```bash
go tool dist list
```

## MCP Tools
```bash
./scripts/build_mcp_tools.sh
```

## Git Hooks
Set Git hooks with built binary files
```bash
cd /path/to/dir
$HOME/devbox/pkg/bash/setup-git-pre-commit-hooks.sh
```

## RESTful API
```bash
go run ./cmd/http/main.go
```

## gRPC API
```bash
go run ./cmd/grpc/main.go
```

# Project Structure

`devbox`は複数のCLIツールを提供する開発ユーティリティ集合体です。Clean Architectureに基づいて設計されており、以下のパッケージ依存関係を持ちます。

```mermaid
graph TD
  A[cmd Layer] --> B[internal/usecases]
  A --> C[internal/interfaces]
  A --> D[internal/domain]
  
  B --> D
  C --> D
  
  E[scripts] --> F[pkg/bin]
  E --> A
  
  F --> H[Cross-platform Binaries]

  classDef cmd fill:#f96,stroke:#333,stroke-width:2px;
  classDef internal fill:#bbf,stroke:#333,stroke-width:1px;
  classDef scripts fill:#bfb,stroke:#333,stroke-width:1px;
  classDef pkg fill:#fbf,stroke:#333,stroke-width:1px;

  class A cmd;
  class B,C,D internal;
  class E scripts;
  class F,H pkg;
```

## Package Overview
- **cmd/**: 各CLIツールのエントリーポイント群。多様なコマンドラインツールが含まれており、それぞれが独立したアプリケーションとして動作します。主要ツール例：
  - `exif-modifier`: Exif処理ユーティリティ
  - `image-converter`: 画像形式変換ツール
  - `json-file-merger`: JSONファイル統合ツール
  - `code-analyzer`: コード解析ツール
  - `depends-visualizer`: 依存関係可視化ツール

- **internal/**: 各ツールのビジネスロジック実装。Clean Architectureに従って以下の層に分離：
  - `domain/`: ドメインモデルとリポジトリインターフェース
  - `usecases/`: ビジネスロジックとユースケース実装
  - `interfaces/`: 外部システムとのインターフェース実装

- **pkg/**: ビルド成果物とデプロイメント用ファイル群：
  - `bin/`: クロスプラットフォーム対応のバイナリファイル（Linux、macOS、Windows）
  - `dos/`: Windows環境向けのバッチファイル群

- **scripts/**: ビルドとデプロイメント自動化スクリプト群。各ツールの個別ビルドスクリプトと統合ビルドスクリプトを提供します。

# Service Implementing Status
[Here](./docs/service_implementation_status.md)

# License
MIT License
