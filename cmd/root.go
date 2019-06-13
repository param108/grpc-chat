package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chat",
	Short: "Central auth provides centralized authorization for opstech products",
}

func init() {
	rootCmd.AddCommand(MigrateCmd, RollbackCmd, ClientCmd, ServerCmd)
}

// Run function lets you run the commands
func Run(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
