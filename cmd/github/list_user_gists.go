package github

import (
	"fmt"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listUserGists.Flags().StringVarP(&gistUsername, "owner", "o", "", "The owner to list all of the public repos")
	listUserGists.MarkFlagRequired("owner")

	GitHubRootCmd.AddCommand(listUserGists)
}

var gistUsername string

var listUserGists = &cobra.Command{
	Use:   "list-user-gists",
	Short: "Finds all of a user's public gists",
	Long:  `Finds all of a user's public gists`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
			return nil
		}
		allGists := []*github.Gist{}
		opt := &github.GistListOptions{}
		for {
			gists, resp, err := client.Gists.List(ctx, gistUsername, opt)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			allGists = append(allGists, gists...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
		for _, gist := range allGists {
			fmt.Printf("%s/\n", *gist.HTMLURL)

		}

		return nil

	},
}
