package github

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v34/github"
	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
)

func init() {
	checkUserExistsCmd.Flags().StringVarP(&user, "user", "u", "", "The users to query (csv supported)")

	GitHubRootCmd.AddCommand(checkUserExistsCmd)
}

var user string

var checkUserExistsCmd = &cobra.Command{
	Use:   "check-user-not-exists",
	Short: "Checks if the user exists. If it doesn't it prints the username",
	Long:  `Who can access a given nwo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		auth_token, ok := os.LookupEnv("GITHUB_TOKEN")
		if !ok {
			fmt.Println("You need to set the GITHUB_TOKEN environment variable.\n")
			return nil
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: auth_token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)
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
		allUsers := []*github.User{}

		for _, username := range usersInput {
			if username == "" {
				continue
			}
			users, resp, err := client.Users.Get(ctx, username)
			if resp.StatusCode == 404 {
				fmt.Println(username)
				continue
			}
			if err != nil {
				fmt.Println(err)
				return nil
			}

			allUsers = append(allUsers, users)
		}
		if root.VerboseOutput {
			fmt.Printf("%d users didn't exist/%s\n", len(allUsers), owner, repo)
		}

		return nil

	},
}
