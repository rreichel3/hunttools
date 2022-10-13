package github

import (
	"fmt"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listPublicRepos.MarkFlagRequired("owner")

	GitHubRootCmd.AddCommand(listPublicRepos)
}

func getAllPublicReposForOrganization(client *github.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		Type: "public",
	}

	allPublicRepos := []*github.Repository{}
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		allPublicRepos = append(allPublicRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allPublicRepos, nil
}

var listPublicRepos = &cobra.Command{
	Use:   "list-public-repos",
	Short: "Lists public repos",
	Long:  `Lists public repos`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_TOKEN environment variable.\n")
			return nil
		}

		if owner == "" && !usePipe {
			return fmt.Errorf("owner or piping is required")
		}

		organizations := getInputTargets()
		allPublicRepos := []*github.Repository{}

		for _, organization := range organizations {
			repos, err := getAllPublicReposForOrganization(client, organization)
			if err != nil {
				return err
			}
			allPublicRepos = append(allPublicRepos, repos...)
		}

		for _, repo := range allPublicRepos {

			fmt.Printf("%v/%v\n", repo.GetOwner().GetLogin(), *repo.Name)

		}

		return nil

	},
}
