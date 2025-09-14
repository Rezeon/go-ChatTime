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
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Gagal load .env")
	}

	dburl := os.Getenv("DATABASE_URL")
	database, err := gorm.Open(postgres.Open(dburl), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Gagal koneksi database")
	}

	fmt.Println("✅ Koneksi database berhasil")
	DB = database
}
