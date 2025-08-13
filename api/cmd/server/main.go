package main

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo インスタンスを作成
	e := echo.New()

	// ミドルウェア設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// ヘルスチェックエンドポイント
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "construction-mgmt-api",
		})
	})

	// API v1 グループ
	v1 := e.Group("/api/v1")

	// 基本的なルート
	v1.GET("/projects", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "プロジェクト一覧",
			"data":    []interface{}{},
		})
	})

	// ポート設定
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// サーバー起動
	log.Printf("サーバーを起動しています。ポート: %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("サーバーの起動に失敗しました:", err)
	}
}
