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
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}

		api.StartServer()
	},
}
