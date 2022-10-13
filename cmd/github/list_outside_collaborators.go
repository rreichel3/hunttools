package github

import (
	"fmt"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listRepos.MarkFlagRequired("owner")

	GitHubRootCmd.AddCommand(listOustideCollaborators)
}

func listAllOutsideCollaborators(client *github.Client, org string) ([]*github.User, error) {
	opt := &github.ListOutsideCollaboratorsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	allCollaborators := []*github.User{}

	for {
		users, resp, err := client.Organizations.ListOutsideCollaborators(ctx, org, opt)
		if err != nil {
			fmt.Println(err)
			return nil, nil
		}
		allCollaborators = append(allCollaborators, users...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allCollaborators, nil
}

var listOustideCollaborators = &cobra.Command{
	Use:   "list-outside-collaborators",
	Short: "Lists outside collaborators for org",
	Long:  `Lists outside collaborators for org`,
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
		for _, organization := range organizations {
			collaborators, err := listAllOutsideCollaborators(client, organization)
			if err != nil {
				return err
			}
			if len(collaborators) > 0 {
				fmt.Printf("%s\n", organization)
				for _, collaborator := range collaborators {
					fmt.Printf("\t%s\n", collaborator.GetLogin())
				}
			}
		}

		return nil

	},
}
