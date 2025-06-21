package server

import (
	"KaungHtetHein116/IVY-backend/api"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var StartProdCmd = &cobra.Command{
	Use: "prod",
	Run: func(cmd *cobra.Command, args []string) {
		// Try to load .env.production, but don't fail if it doesn't exist
		// In production (like Render), environment variables should be set directly
		if err := godotenv.Load(".env.production"); err != nil {
			log.Println("Warning: .env.production file not found, using environment variables")
		}

		api.StartServer()
	},
}
