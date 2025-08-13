## 1. 使用技術・バージョン

| 領域 | 技術スタック |
|------|--------------|
| フロントエンド | Node.js 20 系 / Vite 6 / React 18 / TypeScript 5 / TailwindCSS 3 / Biome 1.8 |
| バックエンド | Go 1.24.6 (devcontainer では 1.23) / Echo / GORM / Swagger / Docker Compose |
| その他 | MySQL 8.0.28 / MinIO / PHPMyAdmin / AWS SDK / CDK (TypeScript) |

## 2. ルートディレクトリ構成

```
.
├── app            # フロントエンド (React)
├── api            # バックエンド (Go)
├── cdk            # AWS CDK によるインフラ定義
├── portal         # ドキュメントサイト (Docusaurus)
├── e2e            # E2E テスト
├── sampleApp      # フロントエンド実装サンプル
├── sampleApi      # バックエンド実装サンプル
└── .devcontainer  # Dev Container 設定
```

## 3. 主要設定ファイル

以下のファイルを新規プロジェクトにコピーすることで、本プロジェクトと同じ環境を構築できます。

### `.devcontainer/devcontainer.json`
```json
{
  "name": "api",
  "dockerComposeFile": "docker-compose.yml",
  "service": "app",
  "workspaceFolder": "/workspace",
  "remoteUser": "root",
  "updateRemoteUserUID": true,
  "features": {
    "ghcr.io/devcontainers/features/go:1": {
      "version": "1.23",
      "golangciLintVersion": "1.54.1"
    },
    "ghcr.io/devcontainers/features/common-utils:2": {
      "installZsh": "true",
      "configureZshAsDefaultShell": "true",
      "installOhMyZsh": "true",
      "installOhMyZshConfig": "true"
    },
    "ghcr.io/devcontainers/features/docker-outside-of-docker": {
      "version": "latest",
      "dockerSock": "/var/run/docker.sock"
    }
  },
  "customizations": {
    "vscode": {
      "extensions": [
        "golang.Go",
        "emeraldwalk.RunOnSave",
        "ms-azuretools.vscode-docker",
        "streetsidesoftware.code-spell-checker",
        "PKief.material-icon-theme",
        "Arjun.swagger-viewer",
        "shardulm94.trailing-spaces"
      ]
    }
  }
}
```

### `.devcontainer/docker-compose.yml`
```yaml
services:
  app:
    image: mcr.microsoft.com/vscode/devcontainers/base:bookworm
    command: sleep infinity
    environment:
      - ENV=devcontainer
    working_dir: '/workspace/api'
    volumes:
      - ..:/workspace

  api:
    build:
      context: ../api
      dockerfile: ../api/Dockerfile.dev
    volumes:
      - ../api:/app
    ports:
      - 80:80
    tty: true
    env_file: ../api/.env.dev
    command: >
      ash -c "go run main.go"
    depends_on:
      mysql:
        condition: service_healthy

  test-api:
    build:
      context: ../api
      dockerfile: ../api/Dockerfile.dev
    volumes:
      - ../api:/app
    ports:
      - 88:80
    tty: true
    command: >
      ash -c "go run main.go"
    depends_on:
      mysql:
        condition: service_healthy

  mysql:
    image: mysql:8.0.28
    platform: linux/x86_64
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=template
      - MYSQL_USER=docker
      - MYSQL_PASSWORD=password
      - TZ=Asia/Tokyo
    volumes:
      - ../api/mysql/mysql-data:/var/lib/mysql
      - ../api/mysql/my.cnf:/etc/mysql/conf.d/my.cnf
      - ../api/mysql/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
    healthcheck:
      test: mysqladmin ping -h 127.0.0.1 -u$$MYSQL_USER -p$$MYSQL_PASSWORD

  phpmyadmin:
    image: phpmyadmin/phpmyadmin
    platform: linux/x86_64
    environment:
      - PMA_ARBITRARY=1
      - PMA_HOST=mysql
      - PMA_USER=docker
      - PMA_PASSWORD=password
    ports:
      - 4040:80
    depends_on:
      mysql:
        condition: service_healthy

  minio:
    image: minio/minio:RELEASE.2024-07-29T22-14-52Z
    environment:
      - MINIO_ROOT_USER=root
      - MINIO_ROOT_PASSWORD=minio123
    command: server /data --address :9000 --console-address :9001
    volumes:
      - ../api/minio/data:/data

  mc:
    image: minio/mc:RELEASE.2023-12-23T08-47-21Z
    depends_on:
      - minio
    entrypoint: >
      /bin/sh -c "
      mc alias set myminio http://minio:9000 root minio123;
      mc mb myminio/sample;
      "
```

### `api/Dockerfile.dev`
```Dockerfile
FROM --platform=linux/amd64 golang:1.24.6-alpine AS builder

WORKDIR /app
COPY ./go.mod ./go.sum ./Makefile ./
RUN apk update && \
    apk add --update --no-cache git && \
    apk add --no-cache gcc && \
    apk add --no-cache musl-dev && \
    apk add --no-cache make && \
    # タイムゾーンの指定（指定しないとテストが通らないため）
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone && \
    wget https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz && \
    tar -zxvf migrate.linux-amd64.tar.gz && \
    mkdir -p /usr/local/bin && \
    mv migrate /usr/local/bin/migrate && \
    go install github.com/air-verse/air@latest

RUN go mod download

EXPOSE 80
```

### `api/.env.example`
```env
DB_USER=docker
DB_PASSWORD=password
DB_HOST=mysql
DB_PORT=3306
DB_NAME=template
ENV=development
AWS_ACCESS_KEY_ID=root
AWS_SECRET_ACCESS_KEY=minio123
AWS_REGION=XXXXX
MINIO_ENDPOINT=http://localhost:9000
FORCE_PATH_STYLE=true
S3_BUCKET=sample
EMAIL_FROM=noreply@dev.digeon.org
POST_CODE_JP_TOKEN=XXXXX
FRONTEND_URL=http://localhost:5173
SIG_KEY=XXXXX
SLOW_QUERY_THRESHOLD_MILLISECOND=1000
QUEUE_NAME=sample-stg-sqs-queue
DISABLE_AUTH=true # 認証を無効にするかどうか
```

### `api/go.mod`
```go
module github.com/digeon-inc/nihon-logix-order-management/api

go 1.24.6

require (
        github.com/aws/aws-lambda-go v1.49.0
        github.com/aws/aws-sdk-go v1.55.7
        github.com/deckarep/golang-set/v2 v2.8.0
        github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
        github.com/go-ozzo/ozzo-validation/v4 v4.3.0
        github.com/go-resty/resty/v2 v2.16.5
        github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
        github.com/godror/godror v0.48.3
        github.com/golang-jwt/jwt/v5 v5.2.2
        github.com/golang/mock v1.6.0
        github.com/google/go-cmp v0.7.0
        github.com/google/uuid v1.6.0
        github.com/labstack/echo/v4 v4.13.4
        github.com/oklog/ulid/v2 v2.1.1
        github.com/shopspring/decimal v1.4.0
        github.com/stretchr/testify v1.10.0
        github.com/swaggo/swag v1.16.4
        github.com/xuri/excelize/v2 v2.9.1
        go.uber.org/zap v1.27.0
        golang.org/x/crypto v0.39.0
        gopkg.in/yaml.v2 v2.4.0
        gorm.io/driver/mysql v1.6.0
        gorm.io/gorm v1.30.0
)

require (
        filippo.io/edwards25519 v1.1.0 // indirect
        github.com/KyleBanks/depth v1.2.1 // indirect
        github.com/PuerkitoBio/purell v1.1.1 // indirect
        github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
        github.com/VictoriaMetrics/easyproto v0.1.4 // indirect
        github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
        github.com/davecgh/go-spew v1.1.1 // indirect
        github.com/go-logfmt/logfmt v0.6.0 // indirect
        github.com/go-openapi/jsonpointer v0.19.5 // indirect
        github.com/go-openapi/jsonreference v0.19.6 // indirect
        github.com/go-openapi/spec v0.20.4 // indirect
        github.com/go-openapi/swag v0.19.15 // indirect
        github.com/go-sql-driver/mysql v1.9.3 // indirect
        github.com/godror/knownpb v0.3.0 // indirect
        github.com/jinzhu/inflection v1.0.0 // indirect
        github.com/jinzhu/now v1.1.5 // indirect
        github.com/jmespath/go-jmespath v0.4.0 // indirect
        github.com/josharian/intern v1.0.0 // indirect
        github.com/labstack/gommon v0.4.2 // indirect
        github.com/mailru/easyjson v0.7.6 // indirect
        github.com/mattn/go-colorable v0.1.14 // indirect
        github.com/mattn/go-isatty v0.0.20 // indirect
        github.com/pmezard/go-difflib v1.0.0 // indirect
        github.com/richardlehane/mscfb v1.0.4 // indirect
        github.com/richardlehane/msoleps v1.0.4 // indirect
        github.com/tiendc/go-deepcopy v1.6.0 // indirect
        github.com/valyala/bytebufferpool v1.0.0 // indirect
        github.com/valyala/fasttemplate v1.2.2 // indirect
        github.com/xuri/efp v0.0.1 // indirect
        github.com/xuri/nfp v0.0.1 // indirect
        go.uber.org/multierr v1.10.0 // indirect
        golang.org/x/exp v0.0.0-20250506013437-ce4c2cf36ca6 // indirect
        golang.org/x/net v0.40.0 // indirect
        golang.org/x/sys v0.33.0 // indirect
        golang.org/x/text v0.26.0 // indirect
        golang.org/x/time v0.11.0 // indirect
        golang.org/x/tools v0.33.0 // indirect
        google.golang.org/protobuf v1.36.6 // indirect
        gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

### `app/package.json`
```json
{
  "name": "app",
  "private": true,
  "version": "0.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc -b && vite build",
    "lint": "biome lint --write ./src",
    "format": "biome format --write ./src",
    "fix": "biome check --write ./src",
    "check-format": "biome check ./src && knip",
    "preview": "vite preview",
    "knip": "knip --fix",
    "test": "vitest",
    "coverage": "vitest run --coverage",
    "check": "tsc -b && npm run fix && npm run knip && npm run coverage",
    "storybook": "storybook dev -p 6006",
    "build-storybook": "storybook build",
    "orval": "orval --config ./orval.config.ts"
  },
  "dependencies": {
    "@dnd-kit/core": "^6.1.0",
    "@dnd-kit/sortable": "^8.0.0",
    "@dnd-kit/utilities": "^3.2.2",
    "@hookform/resolvers": "^3.6.0",
    "@lexical/link": "^0.19.0",
    "@lexical/list": "^0.19.0",
    "@lexical/react": "^0.19.0",
    "@lexical/rich-text": "^0.19.0",
    "@lexical/selection": "^0.19.0",
    "@storybook/manager-api": "^8.3.6",
    "@storybook/theming": "^8.3.6",
    "@storybook/types": "^8.3.6",
    "@tanstack/react-query": "^5.49.2",
    "@tanstack/react-router": "^1.120.13",
    "@tanstack/router-devtools": "^1.43.4",
    "axios": "^1.7.2",
    "date-fns": "^4.1.0",
    "date-holidays": "^3.23.12",
    "js-cookie": "^3.0.5",
    "jsdom": "^25.0.1",
    "lexical": "^0.19.0",
    "msw-storybook-addon": "^2.0.2",
    "react": "^18.3.1",
    "react-dom": "^18.3.1",
    "react-dropzone": "^14.2.3",
    "react-hook-form": "^7.52.0",
    "react-icons": "^5.2.1",
    "react-intersection-observer": "^9.13.1",
    "react-to-print": "^3.1.1",
    "recharts": "^2.15.0",
    "storybook": "^8.1.11",
    "zod": "^3.23.8"
  },
  "devDependencies": {
    "@biomejs/biome": "1.8.3",
    "@faker-js/faker": "^9.5.0",
    "@storybook/addon-essentials": "^8.1.11",
    "@storybook/addon-interactions": "^8.1.11",
    "@storybook/addon-links": "^8.1.11",
    "@storybook/addon-onboarding": "^8.1.11",
    "@storybook/blocks": "^8.1.11",
    "@storybook/react": "^8.1.11",
    "@storybook/react-vite": "^8.1.11",
    "@tanstack/router-plugin": "^1.43.1",
    "@testing-library/react": "^16.0.0",
    "@types/js-cookie": "^3.0.6",
    "@types/node": "^20.14.9",
    "@types/react": "^18.3.3",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react-swc": "^3.5.0",
    "@vitest/coverage-v8": "^3.0.9",
    "autoprefixer": "^10.4.19",
    "knip": "^5.23.2",
    "msw": "^2.3.1",
    "orval": "^7.3.0",
    "postcss": "^8.4.39",
    "tailwindcss": "^3.4.4",
    "typescript": "5.5.2",
    "vite": "^6.2.2",
    "vitest": "^3.0.9"
  },
  "msw": {
    "workerDirectory": [
      "public"
    ]
  }
}
```

### `app/.env.example`
```env
VITE_APP_API_URL=http://localhost/api # APIのURL
VITE_APP_ENABLE_API_MOCKING=false # APIのモックを有効にするかどうか
VITE_APP_DISABLE_AUTH=true # 認証を無効にするかどうか
VITE_APP_DUMMY_USER_ID=00000000000000000000000000 # 認証無効時に使用するダミーユーザーID(26桁)
```

### `package.json` (リポジトリルート)
```json
{
  "scripts": {
    "prepare": "husky"
  },
  "devDependencies": {
    "husky": "^9.1.4"
  }
}
```

## 4. 環境構築手順

1. **リポジトリをクローン**
   ```bash
   git clone <新規リポジトリURL>
   cd <repo>
   ```
2. **依存パッケージのインストール**
   ```bash
   # ルート (husky 用)
   npm install

   # フロントエンド
   cd app
   cp .env.example .env
   npm install
   cd ..

   # バックエンド
   cd api
   cp .env.example .env.dev
   go mod download
   cd ..

   # 必要に応じて
   cd cdk && npm install && cd ..
   cd portal && npm install && cd ..
   cd e2e && npm install && cd ..
   ```
3. **開発用コンテナの起動（バックエンド）**
   ```bash
   cd api
   make run-dev
   ```
4. **フロントエンドの起動**
   ```bash
   cd app
   npm run dev
   ```
5. **API 変更時の連携**
   ```bash
   # バックエンド側で Swagger 生成
   go run create_api.go

   # フロントエンド側で API クライアント生成
   npm run orval
   ```

## 5. 共通開発フロー

- バックエンドのリント/フォーマット
  ```bash
  golangci-lint run ./... --fix
  swag init
  ```
- フロントエンドのチェック
  ```bash
  npm run build
  npm run format
  npm run check
  npm run check-format
  ```

上記の内容を新規プロジェクトに適用することで、本リポジトリと同じ開発環境を再現できます。
