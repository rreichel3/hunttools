package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	// herokuRootCommand.PersistentFlags().StringVarP(&HerokuBearer, "bearer", "b", "localhost", "Bearer token to auth to Heroku")
	rootCmd.AddCommand(herokuRootCommand)
}

var HerokuBearer string
var herokuRootCommand = &cobra.Command{
	Use:   "heroku",
	Short: "Commands to interact with Heroku ",
	Long:  `Commands to interact with heroku.`,
}
