package github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listRepos.MarkFlagRequired("owner")
	listRepos.Flags().StringVarP(&language, "language", "l", "", "Language to filter on")

	GitHubRootCmd.AddCommand(listRepos)
}

var language string

func listAllRepos(client *github.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{}

	allRepos := []*github.Repository{}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allRepos, nil
}

var listRepos = &cobra.Command{
	Use:   "list-repos",
	Short: "Lists repos",
	Long:  `Lists repos`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
			return nil
		}

		if owner == "" && !usePipe {
			return fmt.Errorf("owner or piping is required")
		}

		organizations := getInputTargets()
		allRepos := []*github.Repository{}

		for _, organization := range organizations {
			repos, err := listAllRepos(client, organization)
			if err != nil {
				return err
			}
			allRepos = append(allRepos, repos...)
		}

		for _, repo := range allRepos {
			if language != "" && strings.ToLower(repo.GetLanguage()) == strings.ToLower(language) {
				fmt.Printf("%v/%v\n", repo.GetOwner().GetLogin(), *repo.Name)
			} else if language == "" {
				fmt.Printf("%v/%v\n", repo.GetOwner().GetLogin(), *repo.Name)
			}

		}

		return nil

	},
}
