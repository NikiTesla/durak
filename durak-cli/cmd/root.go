package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "durak",
	Short: "Durak CLI to interact with the Durak card game server",
	Long:  `A command line interface to interact with the Durak card game server.`,
}

func Execute() error {
	return rootCmd.Execute()
}
