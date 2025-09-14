package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, menggunakan environment variables dari sistem")
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		log.Fatal("❌ DATABASE_URL belum di-set")
	}

	database, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi database:", err)
	}

	fmt.Println("✅ Koneksi database berhasil")
	DB = database
}
