package cmd

import (
	"KaungHtetHein116/IVY-backend/cmd/server"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "server",
	Short:        "Setting Server",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to IVY Backend Server!")
	},
}

func init() {
	rootCmd.AddCommand(server.StartDevCmd)
	rootCmd.AddCommand(server.StartProdCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
