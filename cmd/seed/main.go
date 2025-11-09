package main

import (
	"KaungHtetHein116/IVY-backend/config"
	"KaungHtetHein116/IVY-backend/internal/db/seeder"
	"KaungHtetHein116/IVY-backend/internal/entity"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Initialize DB connection
	godotenv.Load(".env.development")
	db := config.ConnectDB()

	db.AutoMigrate(
		&entity.Booking{},
		&entity.Branch{},
		&entity.Category{},
		&entity.Service{},
		&entity.User{},
	)

	// Create and run seeder
	dbSeeder := seeder.NewSeeder(db)
	if err := dbSeeder.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
		os.Exit(1)
	}

	log.Println("Database seeding completed successfully!")
}
