package main

import (
	"KaungHtetHein116/IVY-backend/config"
	"KaungHtetHein116/IVY-backend/internal/db/seeder"
	"log"
	"os"
)

func main() {
	// Initialize DB connection
	db := config.ConnectDB()

	// Create and run seeder
	dbSeeder := seeder.NewSeeder(db)
	if err := dbSeeder.Seed(); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
		os.Exit(1)
	}

	log.Println("Database seeding completed successfully!")
}
