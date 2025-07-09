#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/exif-modifier"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="exif-modifier"
WIN_OUTPUT_NAME="${OUTPUT_NAME}.exe"

echo "Building ${OUTPUT_NAME}..."

# Linux/AMD64向けビルド
echo "Building for Linux/AMD64..."
mkdir -p "${LINUX_AMD64_DIR}"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${LINUX_AMD64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

# Windows/AMD64向けビルド
echo "Building for Windows/AMD64..."
mkdir -p "${WIN_AMD64_DIR}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}" "${PACKAGE}"

# macOS/ARM64向けビルド
echo "Building for macOS/ARM64..."
mkdir -p "${MAC_ARM64_DIR}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o "${MAC_ARM64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

echo "Build completed successfully!"
echo "Usage as example:"
echo "# 現在のフォルダの全ての画像ファイルの日時を設定"
echo "./exif-modifier --datetime 20240315143000"
echo "# 特定のフォルダを処理"
echo "./exif-modifier --folder /path/to/images --datetime 20240315143000"
echo "# 特定の拡張子のみ処理"
echo "./exif-modifier --folder /path/to/images --datetime 20240315143000 --ext .jpg"
echo "# サブフォルダも再帰的に処理"
echo "./exif-modifier --folder /path/to/images --datetime 20240315143000 --recursive"
echo "# ドライランモード（実際には変更せずに確認のみ）"
echo "./exif-modifier --folder /path/to/images --datetime 20240315143000 --dry-run"
echo "# 詳細な出力を表示"
echo "./exif-modifier --folder /path/to/images --datetime 20240315143000 --verbose"
