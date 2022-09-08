package github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listUserGists.Flags().StringVarP(&userInput, "user", "u", "", "The users to query (csv supported)")

	GitHubRootCmd.AddCommand(listUserGists)
}

var userInput string

var listUserGists = &cobra.Command{
	Use:   "list-user-gists",
	Short: "Finds all of a user's public gists",
	Long:  `Finds all of a user's public gists`,
	RunE: func(cmd *cobra.Command, args []string) error {

		usersInput := []string{}
		if usePipe {
			if user != "" {
				return fmt.Errorf("You can't use the pipe with the --user flag")
			}
			usersInput = append(usersInput, getContentFromStdin()...)
		} else if user != "" {
			usersInput = append(usersInput, strings.Split(user, ",")...)
		} else {
			return fmt.Errorf("You need to provide a user to check")
		}

		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
			return nil
		}

		allGists := []*github.Gist{}
		opt := &github.GistListOptions{}
		for _, username := range usersInput {
			for {
				gists, resp, err := client.Gists.List(ctx, username, opt)
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
		}
		for _, gist := range allGists {
			fmt.Printf("%s/\n", *gist.HTMLURL)

		}

		return nil

	},
}
