#!/bin/bash

# エラーが発生したら終了
set -e

# ビルド対象のディレクトリ
CMD_DIR="cmd/cli/file-character-replacer"

# 出力先ディレクトリ
OUTPUT_DIR="./pkg/bin/cli"

# プラットフォーム別の出力先
LINUX_AMD64_DIR="${OUTPUT_DIR}/linux_amd64"
WIN_AMD64_DIR="${OUTPUT_DIR}/win_amd64"
MAC_ARM64_DIR="${OUTPUT_DIR}/darwin_arm64"

# ビルド情報
PACKAGE="github.com/landmaster135/devbox/${CMD_DIR}"
OUTPUT_NAME="file-character-replacer"
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
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./test.txt -from="old" -to="new""
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./test.txt -from="old" -to="new" -backup"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go  -target=./test.txt -from="old" -to="new" -backup -backup-dir=./backups"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./src -from="TODO" -to="DONE" -recursive"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./src -from="TODO" -to="DONE" -recursive -backup"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=$HOME/devbox/pkg/dos/test_file.bat -from=".\\pkg\\bin\\" -to=".\\pkg\\bin\cli\\" -encoding=shift_jis"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./data.txt -from="古い" -to="新しい" -encoding=shift_jis"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./legacy.txt -from="旧" -to="新" -encoding=euc-jp"
echo "  go run $HOME/devbox/cmd/cli/file-character-replacer/main.go -target=./project -from="debug" -to="release" -recursive -dry-run"
