package github

import (
	"context"

	"github.com/spf13/cobra"
)

func init() {
}

var ctx = context.Background()

var GitHubRootCmd = &cobra.Command{
	Use:   "gh",
	Short: "GitHub related commands ",
	Long:  `GitHub related commands`,
}
