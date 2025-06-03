#!/bin/bash

# 共有ToDoリストアプリケーション テスト実行スクリプト

echo "🧪 共有ToDoリストアプリケーション テスト実行"
echo "=================================================="

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# テスト結果を格納する変数
BACKEND_TEST_RESULT=0
FRONTEND_TEST_RESULT=0

# バックエンドテスト実行
echo -e "\n${YELLOW}📋 バックエンドテスト実行中...${NC}"
echo "----------------------------------------"
cd backend
if go test -v ./...; then
    echo -e "${GREEN}✅ バックエンドテスト: 成功${NC}"
else
    echo -e "${RED}❌ バックエンドテスト: 失敗${NC}"
    BACKEND_TEST_RESULT=1
fi
cd ..

# フロントエンドテスト実行（依存関係のインストールが必要）
echo -e "\n${YELLOW}🎨 フロントエンドテスト実行中...${NC}"
echo "----------------------------------------"
cd frontend

# 依存関係がインストールされているかチェック
if [ ! -d "node_modules" ]; then
    echo "📦 依存関係をインストール中..."
    npm install
fi

if npm test -- --run; then
    echo -e "${GREEN}✅ フロントエンドテスト: 成功${NC}"
else
    echo -e "${RED}❌ フロントエンドテスト: 失敗${NC}"
    FRONTEND_TEST_RESULT=1
fi
cd ..

# 結果の表示
echo -e "\n=================================================="
echo -e "${YELLOW}🏁 テスト結果サマリー${NC}"
echo "=================================================="

if [ $BACKEND_TEST_RESULT -eq 0 ]; then
    echo -e "バックエンド:  ${GREEN}✅ PASS${NC}"
else
    echo -e "バックエンド:  ${RED}❌ FAIL${NC}"
fi

if [ $FRONTEND_TEST_RESULT -eq 0 ]; then
    echo -e "フロントエンド: ${GREEN}✅ PASS${NC}"
else
    echo -e "フロントエンド: ${RED}❌ FAIL${NC}"
fi

# 全体の結果
if [ $BACKEND_TEST_RESULT -eq 0 ] && [ $FRONTEND_TEST_RESULT -eq 0 ]; then
    echo -e "\n${GREEN}🎉 全てのテストが成功しました！${NC}"
    exit 0
else
    echo -e "\n${RED}💥 一部のテストが失敗しました${NC}"
    exit 1
fi