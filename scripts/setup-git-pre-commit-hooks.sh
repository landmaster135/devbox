#!/bin/bash

# Git pre-commit hook用のGo製シークレット検知ツールセットアップスクリプト
# ビルド済みバイナリを使用してGitフックを設定

set -e

# 色付きの出力用
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔧 Setting up Go-based secret detector for Git hooks...${NC}"

# 現在のディレクトリを確認
if [ ! -d ".git" ]; then
  echo -e "${RED}❌ This directory is not a Git repository${NC}"
  exit 1
fi

# ビルド済みバイナリのパス
BINARY_PATH="./pkg/bin/cli/linux_amd64/git-pre-commit-hooks"

# バイナリが存在するかチェック
if [ ! -f "$BINARY_PATH" ]; then
  echo -e "${RED}❌ Binary not found at $BINARY_PATH${NC}"
  echo -e "${YELLOW}Please run './scripts/build_secret_detector.sh' first${NC}"
  exit 1
fi

echo -e "${GREEN}✅ Binary found: $BINARY_PATH${NC}"

# バイナリが実行可能かチェック
if [ ! -x "$BINARY_PATH" ]; then
  echo -e "${YELLOW}⚠️  Making binary executable...${NC}"
  chmod +x "$BINARY_PATH"
fi

# ツールディレクトリを作成
TOOL_DIR=".git/hooks/git-pre-commit-hooks"
mkdir -p "$TOOL_DIR"

# バイナリをコピー
echo -e "${YELLOW}📦 Copying binary to Git hooks directory...${NC}"
cp "$BINARY_PATH" "$TOOL_DIR/git-pre-commit-hooks"
chmod +x "$TOOL_DIR/git-pre-commit-hooks"

# Pre-commitフックを作成
echo -e "${YELLOW}🔗 Creating pre-commit hook...${NC}"

cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

# Execute the secret detector tool implemented with Go
TOOL_PATH="$(dirname "$0")/git-pre-commit-hooks/git-pre-commit-hooks"

if [ -x "$TOOL_PATH" ]; then
  "$TOOL_PATH"
else
  echo "❌ Secret detector tool not found at $TOOL_PATH"
  echo "Please run the setup script again."
  exit 1
fi
EOF

# 実行権限を付与
chmod +x .git/hooks/pre-commit

echo -e "${GREEN}✅ Setup completed successfully!${NC}"
echo ""
echo -e "${BLUE}📝 Usage:${NC}"
echo "  • The pre-commit hook will automatically run on 'git commit'"
echo "  • To bypass the check: git commit --no-verify"
echo "  • To test manually: .git/hooks/git-pre-commit-hooks/git-pre-commit-hooks"
echo ""

# テスト実行
echo -e "${YELLOW}🧪 Testing the tool...${NC}"
if .git/hooks/git-pre-commit-hooks/git-pre-commit-hooks --version; then
  echo -e "${GREEN}✅ Tool is working correctly${NC}"
else
  echo -e "${YELLOW}⚠️  Tool test completed (may show 'no staged files' message)${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Secret detector is now active for this repository!${NC}"
echo ""
echo -e "${BLUE}📋 What was installed:${NC}"
echo "  • Binary: .git/hooks/git-pre-commit-hooks/git-pre-commit-hooks"
echo "  • Hook: .git/hooks/pre-commit"
echo ""
echo -e "${BLUE}🔧 To update the tool:${NC}"
echo "  1. Run: ./scripts/build_secret_detector.sh"
echo "  2. Run: ./scripts/setup-git-pre-commit-hooks-hook.sh"
echo ""
echo -e "${BLUE}🗑️  To remove the hook:${NC}"
echo "  • Delete: .git/hooks/pre-commit"
echo "  • Delete: .git/hooks/git-pre-commit-hooks/"
