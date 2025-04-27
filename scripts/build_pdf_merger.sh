#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/pdf-merger"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="pdf-merger"
WIN_OUTPUT_NAME="${OUTPUT_NAME}.exe"

echo "Building ${OUTPUT_NAME}..."

# Linux/AMD64向けビルド
echo "Building for Linux/AMD64..."
mkdir -p "${LINUX_AMD64_DIR}"
GOOS=linux GOARCH=amd64 go build -o "${LINUX_AMD64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

# Windows/AMD64向けビルド
echo "Building for Windows/AMD64..."
mkdir -p "${WIN_AMD64_DIR}"
GOOS=windows GOARCH=amd64 go build -o "${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}" "${PACKAGE}"

echo "Build completed successfully!"
echo "Usage as example:"
echo "  ./pdf-merger -dir ./images"
echo "  ./pdf-merger -dir ./vacation_photos -out ./vacation_album.pdf"
