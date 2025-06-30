package server

import (
	"KaungHtetHein116/IVY-backend/api"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var StartLocalCmd = &cobra.Command{
	Use: "local",
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env.local for local environment
		if err := godotenv.Load(".env.local"); err != nil {
			log.Fatal("Error loading .env.local file")
		}

		api.StartServer()
	},
}
