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
		// Try to load .env.local, but don't fail if it doesn't exist
		// In local, environment variables should be set directly
		if err := godotenv.Load(".env.local"); err != nil {
			log.Println("Warning: .env.local file not found, using environment variables")
		}

		api.StartServer()
	},
}
