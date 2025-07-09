#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/movie-converter-for-gif"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin/cli"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="movie-converter-for-gif"
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
echo "  ./movie-converter-for-gif -input sample.mp4"
echo "  ./movie-converter-for-gif -input sample.mp4 -fps 15 -width 480 -speed 1 -loop -1"
echo "  ./movie-converter-for-gif -input animation.gif -fps 30"
echo "  ./movie-converter-for-gif -input video.mp4 -output converted_video.gif"
echo "  ./movie-converter-for-gif -input-dir ./videos -input-ext mp4 -output-dir ./gifs -output-ext gif"
echo "  ./movie-converter-for-gif -input-dir ./media -input-ext mp4 -output-dir ./converted -output-ext gif -recursive"
echo "  ./movie-converter-for-gif -input-dir ./videos -input-ext mp4 -output-dir ./gifs -output-ext gif -fps 15 -width 320 -speed 1"
echo "  ./movie-converter-for-gif -input-dir ./animations -input-ext gif -output-dir ./videos -output-ext mp4 -fps 24"
