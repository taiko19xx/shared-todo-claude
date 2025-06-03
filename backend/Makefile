.PHONY: test test-verbose test-coverage clean build

# テスト実行
test:
	go test ./...

# 詳細出力でテスト実行
test-verbose:
	go test -v ./...

# カバレッジ付きテスト実行
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# カバレッジ表示
show-coverage:
	go tool cover -func=coverage.out

# テストファイルのクリーンアップ
clean:
	rm -f coverage.out coverage.html

# ビルド
build:
	go build -o bin/main .

# 開発用実行
dev:
	air -c .air.toml

# 全テストを実行してカバレッジを生成
test-all: clean test-coverage show-coverage