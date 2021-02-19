package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "ht",
		Short: "Tools for the hunt",
		Long: `HuntTools (ht) is a CLI designed to help you quickly hunt for things in a network. 
		In the end, you should spend less time trying to remember how to write a loop in bash and more time actually hunting`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(pingCmd)
	rootCmd.AddCommand(hostnamesCommand)
}
