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
		// Load .env.development for development environment
		if err := godotenv.Load(".env.development"); err != nil {
			log.Fatal("Error loading .env.development file")
		}

		api.StartServer()
	},
}
