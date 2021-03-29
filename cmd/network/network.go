package network

import (
	"github.com/spf13/cobra"
)

func init() {
	// herokuRootCommand.PersistentFlags().StringVarP(&HerokuBearer, "bearer", "b", "localhost", "Bearer token to auth to Heroku")
	// RootCmd.AddCommand(herokuRootCommand)
}

var HerokuBearer string
var NetworkRootCmd = &cobra.Command{
	Use:   "network",
	Short: "Commands to interact with Heroku ",
	Long:  `Commands to interact with heroku.`,
}
