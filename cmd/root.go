package cmd

import (
	github "github.com/rreichel3/hunttools/cmd/github"
	heroku "github.com/rreichel3/hunttools/cmd/heroku"
	network "github.com/rreichel3/hunttools/cmd/network"
	redis "github.com/rreichel3/hunttools/cmd/redis"

	"github.com/spf13/cobra"

	"github.com/rreichel3/hunttools/cmd/root_flags"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	RootCmd = &cobra.Command{
		Use:   "ht",
		Short: "Tools for the hunt",
		Long: `HuntTools (ht) is a CLI designed to help you perform repetitive operations easily.  
The expected input type is always a newline separated file and the default output type will return the same.   
Using this, you should spend less time trying to remember how to write a loop in bash and more time actually hunting`,
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&root_flags.VerboseOutput, "verbose", "v", false, "Output full results, not newline formatted")
	RootCmd.AddCommand(redis.RedisRootCmd)
	RootCmd.AddCommand(github.GitHubRootCmd)
	RootCmd.AddCommand(heroku.HerokuRootCmd)
	RootCmd.AddCommand(network.NetworkRootCmd)

}
