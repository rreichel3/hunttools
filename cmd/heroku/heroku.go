package heroku

import (
	"context"
	"github.com/spf13/cobra"
)

func init() {
}

var ctx = context.Background()

var HerokuBearer string
var HerokuRootCmd = &cobra.Command{
	Use:   "heroku",
	Short: "Commands to interact with Heroku ",
	Long:  `Commands to interact with heroku.`,
}
