#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/claude-code-usage"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin/cli"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="claude-code-usage"
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
echo "Built binaries:"
echo "  Linux AMD64:   ${LINUX_AMD64_DIR}/${OUTPUT_NAME}"
echo "  Windows AMD64: ${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}"
echo "  macOS ARM64:   ${MAC_ARM64_DIR}/${OUTPUT_NAME}"
echo ""
echo "Usage examples:"
echo "  # Daily report"
echo "  ${OUTPUT_NAME} daily"
echo ""
echo "  # Session report with JSON output"
echo "  ${OUTPUT_NAME} session --json"
echo ""
echo "  # Filter by date range"
echo "  ${OUTPUT_NAME} daily --since 20250525 --until 20250530"
echo ""
echo "  # Custom Claude data path"
echo "  ${OUTPUT_NAME} daily --path /custom/path/to/.claude"
