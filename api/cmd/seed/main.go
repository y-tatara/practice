package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 同じ構造体定義（migrate/main.goと同じ）
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Role      string    `json:"role" gorm:"default:'user'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:'planning'"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Supplier struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	ContactName string    `json:"contact_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

	log.Println("初期データの投入を開始します...")

	// ユーザーデータ
	if err := seedUsers(db); err != nil {
		log.Fatal("ユーザーデータの投入に失敗しました:", err)
	}

	// プロジェクトデータ
	if err := seedProjects(db); err != nil {
		log.Fatal("プロジェクトデータの投入に失敗しました:", err)
	}

	// 仕入先データ
	if err := seedSuppliers(db); err != nil {
		log.Fatal("仕入先データの投入に失敗しました:", err)
	}

	log.Println("初期データの投入が完了しました!")
}

func seedUsers(db *gorm.DB) error {
	// パスワードハッシュ化
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	managerPassword, _ := bcrypt.GenerateFromPassword([]byte("manager123"), bcrypt.DefaultCost)

	users := []User{
		{
			Name:     "管理者",
			Email:    "admin@construction.com",
			Password: string(adminPassword),
			Role:     "admin",
		},
		{
			Name:     "施工管理者",
			Email:    "manager@construction.com",
			Password: string(managerPassword),
			Role:     "manager",
		},
		{
			Name:     "現場監督",
			Email:    "supervisor@construction.com",
			Password: string(managerPassword),
			Role:     "supervisor",
		},
	}

	for _, user := range users {
		if err := db.FirstOrCreate(&user, User{Email: user.Email}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedProjects(db *gorm.DB) error {
	projects := []Project{
		{
			Name:        "東京駅前オフィスビル建設",
			Description: "地上20階建てのオフィスビル建設プロジェクト",
			Status:      "planning",
			StartDate:   time.Now().AddDate(0, 1, 0),
			EndDate:     time.Now().AddDate(2, 0, 0),
			Address:     "東京都千代田区丸の内1-1-1",
		},
		{
			Name:        "横浜マンション建設",
			Description: "地上15階建ての分譲マンション建設",
			Status:      "in_progress",
			StartDate:   time.Now().AddDate(0, -3, 0),
			EndDate:     time.Now().AddDate(1, 3, 0),
			Address:     "神奈川県横浜市西区みなとみらい2-2-2",
		},
		{
			Name:        "大阪商業施設リノベーション",
			Description: "既存の商業施設の大規模リノベーション",
			Status:      "completed",
			StartDate:   time.Now().AddDate(-1, 0, 0),
			EndDate:     time.Now().AddDate(0, -1, 0),
			Address:     "大阪府大阪市北区梅田3-3-3",
		},
	}

	for _, project := range projects {
		if err := db.FirstOrCreate(&project, Project{Name: project.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedSuppliers(db *gorm.DB) error {
	suppliers := []Supplier{
		{
			Name:        "東京建材株式会社",
			ContactName: "田中太郎",
			Phone:       "03-1234-5678",
			Email:       "tanaka@tokyo-kenzai.co.jp",
			Address:     "東京都江東区豊洲1-1-1",
		},
		{
			Name:        "関西鉄鋼商事",
			ContactName: "山田花子",
			Phone:       "06-2345-6789",
			Email:       "yamada@kansai-tekko.co.jp",
			Address:     "大阪府大阪市住之江区南港2-2-2",
		},
		{
			Name:        "横浜電設工業",
			ContactName: "佐藤次郎",
			Phone:       "045-3456-7890",
			Email:       "sato@yokohama-densetsu.co.jp",
			Address:     "神奈川県横浜市港北区新横浜3-3-3",
		},
	}

	for _, supplier := range suppliers {
		if err := db.FirstOrCreate(&supplier, Supplier{Name: supplier.Name}).Error; err != nil {
			return err
		}
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
