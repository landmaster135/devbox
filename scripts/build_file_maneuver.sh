#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/file-maneuver"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="file-maneuver"
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
echo "  ./file-maneuver "
echo "  --src-dirs "/path/to/source" "
echo "  --extensions "jpg,png,gif" "
echo "  --dest-dir "/path/to/destination""
echo "  ./file-maneuver "
echo "  --src-dirs "/path/to/source1,/path/to/source2,/path/to/source3" "
echo "  --extensions "pdf,doc,docx,txt" "
echo "  --dest-dir "/path/to/documents""
echo "  ./file-maneuver "
echo "  --src-dirs "/path/to/source" "
echo "  --extensions "mp4,avi,mkv" "
echo "  --dest-dir "/path/to/videos" "
echo "  --recursive"
echo "  ./file-maneuver "
echo "  --src-dirs "/path/to/source" "
echo "  --extensions "jpg,png" "
echo "  --dest-dir "/path/to/destination" "
echo "  --dry-run"
echo "  ./file-maneuver "
echo "  --src-dirs "/path/to/source" "
echo "  --extensions "jpg,png" "
echo "  --dest-dir "/path/to/destination" "
echo "  --workers 8"
