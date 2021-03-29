package github

import (
	"fmt"
	"github.com/google/go-github/v34/github"
	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"golang.org/x/oauth2"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	whoCanAccessCmd.Flags().StringVarP(&NWO, "nwo", "n", "", "The given nwo (github/github)")
	whoCanAccessCmd.MarkFlagRequired("nwo")

	GitHubRootCmd.AddCommand(whoCanAccessCmd)
}

var NWO string

var whoCanAccessCmd = &cobra.Command{
	Use:   "who-can-access",
	Short: "Finds who can access a given nwo",
	Long:  `Who can access a given nwo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		auth_token, ok := os.LookupEnv("GITHUB_PAT")
		if !ok {
			fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
			return nil
		}
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: auth_token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := github.NewClient(tc)

		components := strings.Split(NWO, "/")
		owner := components[0]
		repo := components[1]
		allUsers := []*github.User{}
		opt := &github.ListCollaboratorsOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}
		for {
			users, resp, err := client.Repositories.ListCollaborators(ctx, owner, repo, opt)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			allUsers = append(allUsers, users...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
		permissionCounts := make(map[string]int)
		for _, user := range allUsers {
			perms := *user.Permissions
			for permission, yes := range perms {
				if val, ok := permissionCounts[permission]; ok {
					if yes {
						permissionCounts[permission] = val + 1
					}
				} else {
					if yes {
						permissionCounts[permission] = 1
					} else {
						permissionCounts[permission] = 0
					}
				}
			}
		}

		fmt.Printf("%d users have access to %s/%s\n", len(allUsers), owner, repo)
		for permission, count := range permissionCounts {
			fmt.Printf("%s: %d\n", permission, count)
		}
		if root.VerboseOutput {
			for _, user := range allUsers {
				fmt.Println(*user.Login)
				fmt.Println(*user.Permissions)
			}
		}

		return nil

	},
}
