#!/bin/bash

# 共有ToDoリストアプリケーション カバレッジ実行スクリプト

set -e

echo "📊 共有ToDoリストアプリケーション カバレッジ測定"
echo "=================================================="

# カラー定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# バックエンドカバレッジ
echo -e "\n${YELLOW}📋 バックエンドカバレッジ測定中...${NC}"
echo "----------------------------------------"
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
echo -e "${GREEN}✅ バックエンドカバレッジレポート生成: backend/coverage.html${NC}"
go tool cover -func=coverage.out | tail -n 1
cd ..

# フロントエンドカバレッジ
echo -e "\n${YELLOW}🎨 フロントエンドカバレッジ測定中...${NC}"
echo "----------------------------------------"
cd frontend

# 依存関係がインストールされているかチェック
if [ ! -d "node_modules" ]; then
    echo "📦 依存関係をインストール中..."
    npm install
fi

npm run test:coverage -- --run
echo -e "${GREEN}✅ フロントエンドカバレッジレポート生成: frontend/coverage/index.html${NC}"
cd ..

echo -e "\n${GREEN}🎉 カバレッジ測定完了！${NC}"
echo "バックエンドレポート: backend/coverage.html"
echo "フロントエンドレポート: frontend/coverage/index.html"