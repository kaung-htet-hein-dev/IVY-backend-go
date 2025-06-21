package server

import (
	"KaungHtetHein116/IVY-backend/api"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var StartDevCmd = &cobra.Command{
	Use: "dev",
	Run: func(cmd *cobra.Command, args []string) {
		// Try to load .env.development, but don't fail if it doesn't exist
		// In development, environment variables should be set directly
		if err := godotenv.Load(".env.development"); err != nil {
			log.Println("Warning: .env.development file not found, using environment variables")
		}

		api.StartServer()
	},
}
