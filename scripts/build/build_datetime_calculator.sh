#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/datetime-calculator"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin/cli"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="datetime-calculator"
WIN_OUTPUT_NAME="${OUTPUT_NAME}.exe"

echo "Building ${OUTPUT_NAME}..."

# Linux/AMD64向けビルド
echo "Building for Linux/AMD64..."
mkdir -p "${LINUX_AMD64_DIR}"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${LINUX_AMD64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

# Windows/AMD64向けビルド(Windows Defender誤検知対策で-ldflagsから-sを省く)
echo "Building for Windows/AMD64..."
mkdir -p "${WIN_AMD64_DIR}"
GOOS=windows GOARCH=amd64 go build -ldflags="-w" -trimpath -o "${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}" "${PACKAGE}"

# macOS/ARM64向けビルド
echo "Building for macOS/ARM64..."
mkdir -p "${MAC_ARM64_DIR}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o "${MAC_ARM64_DIR}/${OUTPUT_NAME}" "${PACKAGE}"

echo "Build completed successfully!"
echo "Usage as example:"
echo "  ./bin/datetime-calculator -operation add -year 2025 -month 1 -day 15 -hour 12 -minute 30 -second 0 -duration-year 1 -duration-month 2 -duration-day 10 -duration-hour 5 -duration-minute 30 -duration-second 45"
echo "  ./bin/datetime-calculator -o add -y 2025 -m 1 -d 15 -hr 12 -min 30 -s 0 -dy 1 -dm 2 -dd 10 -dh 5 -dmin 30 -ds 45"
echo "  ./bin/datetime-calculator -o add -y 2025 -m 6 -d 15 -dy 2 -dm 3"
