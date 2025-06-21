package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	// Build DSN
	dsn := getDSN()

	// Parse pgx config
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DSN: %v", err)
	}

	// Use simple protocol to avoid prepared statements
	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Open stdlib DB connection
	sqlDB := stdlib.OpenDB(*cfg)

	// Open GORM DB with custom connection
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false, // make sure GORM doesn't prepare statements either
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool settings
	rawDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw DB: %v", err)
	}
	rawDB.SetMaxOpenConns(25)
	rawDB.SetMaxIdleConns(25)
	rawDB.SetConnMaxLifetime(5 * time.Minute)

	DB = db
	return DB
}

func getDSN() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")
	if sslmode == "" {
		sslmode = "require"
	}

	return fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}
