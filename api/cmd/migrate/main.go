package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Project プロジェクトテーブル
type Project struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	Status      string `json:"status" gorm:"default:'planning'"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Address     string `json:"address"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// SiteMaterial 現場資材テーブル
type SiteMaterial struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	ProjectID    uint    `json:"project_id" gorm:"not null"`
	Name         string  `json:"name" gorm:"not null"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity"`
	UnitPrice    float64 `json:"unit_price"`
	Location     string  `json:"location"`
	DeliveryDate string  `json:"delivery_date"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Project      Project `json:"project" gorm:"foreignKey:ProjectID"`
}

// PurchaseOrder 発注テーブル
type PurchaseOrder struct {
	ID           uint    `json:"id" gorm:"primaryKey"`
	ProjectID    uint    `json:"project_id" gorm:"not null"`
	SupplierID   uint    `json:"supplier_id" gorm:"not null"`
	OrderNumber  string  `json:"order_number" gorm:"unique;not null"`
	TotalAmount  float64 `json:"total_amount"`
	Status       string  `json:"status" gorm:"default:'pending'"`
	OrderDate    string  `json:"order_date"`
	DeliveryDate string  `json:"delivery_date"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	Project      Project `json:"project" gorm:"foreignKey:ProjectID"`
}

// Supplier 仕入先テーブル
type Supplier struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	ContactName string `json:"contact_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Task 工程テーブル
type Task struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	ProjectID   uint    `json:"project_id" gorm:"not null"`
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Status      string  `json:"status" gorm:"default:'not_started'"`
	Progress    float64 `json:"progress" gorm:"default:0"`
	StartDate   string  `json:"start_date"`
	EndDate     string  `json:"end_date"`
	AssigneeID  uint    `json:"assignee_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Project     Project `json:"project" gorm:"foreignKey:ProjectID"`
}

// User ユーザーテーブル
type User struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"-" gorm:"not null"`
	Role      string `json:"role" gorm:"default:'user'"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func main() {
	// 環境変数から接続情報を取得
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3307")
	dbUser := getEnv("DB_USER", "app")
	dbPassword := getEnv("DB_PASSWORD", "apppass")
	dbName := getEnv("DB_NAME", "construction_mgmt")

	// DSN作成
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// データベース接続
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("データベース接続に失敗しました:", err)
	}

	log.Println("データベースマイグレーションを開始します...")

	// マイグレーション実行
	err = db.AutoMigrate(
		&User{},
		&Project{},
		&Supplier{},
		&SiteMaterial{},
		&PurchaseOrder{},
		&Task{},
	)
	if err != nil {
		log.Fatal("マイグレーションに失敗しました:", err)
	}

	log.Println("データベースマイグレーションが完了しました!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
