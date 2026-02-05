package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Client() *gorm.DB {
	return DB
}

func Connect() {

	getEnv := func(key, fallback string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return fallback
	}

	dbUser := getEnv("DB_USER", "manager")
	dbName := getEnv("DB_NAME", "wallet_db")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	dsnRoot := fmt.Sprintf("host=%s user=%s dbname=postgres port=%s sslmode=disable", dbHost, dbUser, dbPort)
	rootDB, err := gorm.Open(postgres.Open(dsnRoot), &gorm.Config{})

	if err == nil {
		var count int64
		rootDB.Raw("SELECT count(*) FROM pg_database WHERE datname = ?", dbName).Scan(&count)
		if count == 0 {
			log.Printf("Database '%s' not found. Creating it...", dbName)
			sqlDB, _ := rootDB.DB()
			sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		}
		sqlRoot, _ := rootDB.DB()
		sqlRoot.Close()
	}

	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Successfully connected to Postgres")
	DB = db
}
