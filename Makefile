# 建築業向け統合管理システム - 開発環境管理

.PHONY: setup start stop clean build test lint

# デフォルトターゲット
all: setup

# 開発環境の初期セットアップ
setup:
	@echo "=== 開発環境セットアップ開始 ==="
	@echo "1. APIディレクトリ構造の作成..."
	mkdir -p api/cmd/server api/cmd/migrate api/cmd/seed api/internal/adapter api/internal/api api/internal/config api/internal/domain api/internal/usecase api/docs
	@echo "2. go.modの初期化..."
	cd api && if [ ! -f go.mod ]; then go mod init construction-mgmt-api; fi
	cd api && if [ ! -f go.sum ]; then go mod tidy; fi
	@echo "3. Docker環境の構築（バックエンド・DB・MinIO）..."
	docker-compose up -d db minio api
	@echo "4. データベースの初期化を待機中..."
	sleep 15
	@echo "5. フロントエンドの依存関係をインストール..."
	cd app && npm install --legacy-peer-deps
	@echo "6. データベースのマイグレーション実行..."
	sleep 5
	docker-compose exec api go run cmd/migrate/main.go
	@echo "7. 初期データの投入..."
	docker-compose exec api go run cmd/seed/main.go
	@echo "=== セットアップ完了 ==="
	@echo "フロントエンド起動: cd app && npm start"

# 開発サーバー起動
start:
	@echo "=== 開発サーバー起動 ==="
	docker-compose up -d
	@echo "バックエンド: http://localhost:8080"
	@echo "MinIO: http://localhost:9010"
	@echo ""
	@echo "フロントエンド起動は別ターミナルで:"
	@echo "cd app && npm start"

# フロントエンド個別起動
start-frontend:
	@echo "=== フロントエンド起動 ==="
	cd app && npm start

# 開発サーバー停止
stop:
	@echo "=== 開発サーバー停止 ==="
	docker-compose down

# 全体クリーンアップ
clean:
	@echo "=== 環境クリーンアップ ==="
	docker-compose down -v
	docker system prune -f
	cd app && rm -rf node_modules

# バックエンドビルド
build-api:
	@echo "=== バックエンドビルド ==="
	cd api && go build -o bin/server cmd/server/main.go

# フロントエンドビルド
build-app:
	@echo "=== フロントエンドビルド ==="
	cd app && npm run build

# 全体ビルド
build: build-api build-app

# テスト実行
test:
	@echo "=== テスト実行 ==="
	docker-compose exec api go test ./...
	cd app && npm test

# リント実行
lint:
	@echo "=== リント実行 ==="
	docker-compose exec api golangci-lint run
	cd app && npm run lint

# データベースリセット
db-reset:
	@echo "=== データベースリセット ==="
	docker-compose exec db mysql -u root -prootpass -e "DROP DATABASE IF EXISTS construction_mgmt; CREATE DATABASE construction_mgmt;"
	docker-compose exec api go run cmd/migrate/main.go
	docker-compose exec api go run cmd/seed/main.go

# 既存MySQLとの競合確認
check-mysql:
	@echo "=== MySQL競合確認 ==="
	@echo "既存MySQLプロセス:"
	@ps aux | grep mysql | grep -v grep || echo "MySQLプロセスなし"
	@echo ""
	@echo "ポート使用状況:"
	@lsof -i :3306 || echo "ポート3306は空いています"
	@lsof -i :3307 || echo "ポート3307は空いています"
	@lsof -i :9000 || echo "ポート9000は空いています"
	@lsof -i :9001 || echo "ポート9001は空いています"
	@lsof -i :9010 || echo "ポート9010は空いています"
	@lsof -i :9011 || echo "ポート9011は空いています"

# API初期化（ディレクトリ作成とgo.mod初期化）
init-api:
	@echo "=== API環境初期化 ==="
	mkdir -p api/cmd/server api/cmd/migrate api/cmd/seed
	mkdir -p api/internal/adapter api/internal/api api/internal/config api/internal/domain api/internal/usecase
	mkdir -p api/docs
	cd api && if [ ! -f go.mod ]; then go mod init construction-mgmt-api; fi
	@echo "go.modファイルを作成しました"

# ログ確認
logs:
	docker-compose logs -f
	@echo "=== 開発環境情報 ==="
	@echo "フロントエンド: http://localhost:3000"
	@echo "バックエンドAPI: http://localhost:8080"
	@echo "Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "MinIO: http://localhost:9010 (admin/password)"
	@echo "MySQL: localhost:3307 (root/rootpass)"
