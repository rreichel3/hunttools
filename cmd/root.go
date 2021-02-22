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
		Long: `HuntTools (ht) is a CLI designed to help you perform repetitive operations easily.  
The expected input type is always a newline separated file and the default output type will return the same.   
Using this, you should spend less time trying to remember how to write a loop in bash and more time actually hunting`,
	}
	VerboseOutput bool
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().BoolVarP(&VerboseOutput, "verbose", "v", false, "Output full results, not newline formatted")
}
