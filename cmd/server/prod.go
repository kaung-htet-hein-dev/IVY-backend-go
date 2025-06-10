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
		// Load .env.production for production environment
		if err := godotenv.Load(".env.production"); err != nil {
			log.Fatal("Error loading .env.production file")
		}

		api.StartServer()
	},
}
