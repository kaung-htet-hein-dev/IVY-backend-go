package config

import (
	"fmt"
	"os"
	"time"

	"KaungHtetHein116/IVY-backend/internal/entity"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dns := getDSN()

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}

	sqlDB, err := db.DB()

	sqlDB.SetMaxOpenConns(25)                 // Max open connections
	sqlDB.SetMaxIdleConns(25)                 // Max idle connections
	sqlDB.SetConnMaxLifetime(5 * time.Minute) // Connection lifetime

	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Branch{},
		&entity.Category{},
		&entity.Service{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database %v", err)
	}

	DB = db

	return DB
}

func getDSN() string {
	err := godotenv.Load()
	if err != nil {
		log.Error("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return dsn
}
