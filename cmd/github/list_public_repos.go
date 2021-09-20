package github

import (
	"fmt"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listPublicRepos.Flags().StringVarP(&owner, "owner", "o", "", "The owner to list all of the public repos")
	listPublicRepos.MarkFlagRequired("owner")

	GitHubRootCmd.AddCommand(listPublicRepos)
}

var owner string

var listPublicRepos = &cobra.Command{
	Use:   "list-public-repos",
	Short: "Finds who can access a given nwo",
	Long:  `Who can access a given nwo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
			return nil
		}
		opt := &github.RepositoryListByOrgOptions{
			Type: "public",
		}

		allPublicRepos := []*github.Repository{}
		for {
			repos, resp, err := client.Repositories.ListByOrg(ctx, owner, opt)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			allPublicRepos = append(allPublicRepos, repos...)
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
		for _, repo := range allPublicRepos {

			fmt.Printf("%v/%v\n", owner, *repo.Name)

		}
		// fmt.Printf("%d users have access to %s/%s\n", len(allUsers), owner, repo)

		// if root.VerboseOutput {
		// 	for _, user := range allUsers {
		// 		fmt.Println(*user.Login)
		// 		fmt.Println(*user.Permissions)
		// 	}
		// }

		return nil

	},
}
