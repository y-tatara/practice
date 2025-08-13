# Nihon Logix Order Management – 開発環境ガイド

本リポジトリと同じ開発環境を別プロジェクトで再現するための最小構成をまとめます。フロントエンドは React + TypeScript、バックエンドは Go、データベースは MySQL を利用します。

## 技術スタック
- フロントエンド: React 18, TypeScript 5, Vite 6, Tailwind CSS 3, Node.js 20
- バックエンド: Go 1.24.6 (devcontainer では 1.23), Echo, GORM
- データベース: MySQL 8.0.28
- 補助ツール: Docker Compose, VS Code Dev Containers, MinIO

## ルートディレクトリ構造
├── app            # React/TypeScript フロントエンド
│   ├── src        # コンポーネント・画面・API クライアント
│   └── public     # 静的ファイル
├── api            # Go 製 Web API
│   ├── adapter    # DB・S3 等のインフラ接続
│   ├── api        # ハンドラ・ルーティング・スキーマ
│   ├── config     # 環境変数・設定値
│   ├── domain     # ドメインモデル
│   └── usecase    # アプリケーションサービス
├── cdk            # AWS CDK によるインフラ定義
├── portal         # Docusaurus で管理するドキュメント
├── e2e            # E2E テスト用（任意）
├── sampleApp      # フロントエンド実装例
├── sampleApi      # バックエンド実装例
└── .devcontainer  # Dev Container 設定
