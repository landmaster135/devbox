#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/movie-converter-for-webm"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="movie-converter-for-webm"
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
echo "  ./movie-converter-for-webm -input video.mp4 -crf 20 -audio-bitrate 192k"
echo "  ./movie-converter-for-webm -input video.mp4 -crf 15 -audio-codec vorbis -audio-bitrate 256k"
echo "  ./movie-converter-for-webm -input video.mp4 -conversion-mode cbr -video-bitrate 1M"
echo "  ./movie-converter-for-webm -input video.mp4 -conversion-mode cbr -video-bitrate 500k -audio-bitrate 96k"
echo "  ./movie-converter-for-webm -input-dir ./project -input-ext mp4 -output-dir ./webm_output -output-ext webm -recursive -crf 25"
echo "  for ext in mp4 mkv avi; do"
echo "  ./movie-converter-for-webm -input-dir ./videos -input-ext $ext -output-dir ./webm -output-ext webm -crf 28"
echo "  done"
