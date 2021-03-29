package heroku

import (
	"context"
	"github.com/spf13/cobra"
)

func init() {
	// herokuRootCommand.PersistentFlags().StringVarP(&HerokuBearer, "bearer", "b", "localhost", "Bearer token to auth to Heroku")
	// RootCmd.AddCommand(herokuRootCommand)
}

var ctx = context.Background()

var HerokuBearer string
var HerokuRootCmd = &cobra.Command{
	Use:   "heroku",
	Short: "Commands to interact with Heroku ",
	Long:  `Commands to interact with heroku.`,
}
