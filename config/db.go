package config

import (
	"fmt"
	"net"
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
	dsn := getDSN()

	log.Printf("Connecting to database with DSN: %s", dsn)

	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DSN: %v", err)
	}

	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	sqlDB := stdlib.OpenDB(*cfg)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

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

	// Resolve host to IPv4 address to avoid Docker IPv6 issues
	if host != "" {
		addrs, err := net.LookupHost(host)
		if err == nil {
			// Find first IPv4 address
			for _, addr := range addrs {
				if ip := net.ParseIP(addr); ip != nil && ip.To4() != nil {
					log.Printf("Resolved %s to IPv4: %s", host, addr)
					host = addr
					break
				}
			}
		}
	}

	// For Supabase and other cloud databases, ensure proper SSL configuration
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)

	// Log DSN without password for debugging
	log.Printf("DSN (without password): user=%s host=%s port=%s dbname=%s sslmode=%s", user, host, port, dbname, sslmode)

	return dsn
}
