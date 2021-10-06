package github

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
	GitHubRootCmd.PersistentFlags().StringVar(&owner, "owner", "o", "The owner to query")
	GitHubRootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "The repo to query")
}

var ctx = context.Background()

var owner string
var repo string

var GitHubRootCmd = &cobra.Command{
	Use:   "gh",
	Short: "GitHub related commands ",
	Long:  `GitHub related commands`,
}
