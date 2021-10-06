package github

import (
	"fmt"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	overscopedPages.MarkFlagRequired("owner")

	GitHubRootCmd.AddCommand(overscopedPages)
}

var overscopedPages = &cobra.Command{
	Use:   "overscoped-pages",
	Short: "Finds who can access a given nwo",
	Long:  `Who can access a given nwo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_PAT environment variable.")
			return nil
		}
		opt := &github.RepositoryListByOrgOptions{
			Type: "all",
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
			if *repo.HasPages {
				// We know the repo has pages enabled
				pagesConfig, _, err := client.Repositories.GetPagesInfo(ctx, owner, *repo.Name)
				if err != nil {
					fmt.Printf("%v", err)
					continue
				}
				fmt.Println("")
				if *repo.Visibility == "private" {
					//TODO: This should be filled out once the GitHub Library is updated to support visibility of Pages
				}
				fmt.Printf("%v/%v - %v\n", owner, *repo.Name, *pagesConfig.Source.Path)

			}

		}
		return nil

	},
}
