package github

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v34/github"
	root "github.com/rreichel3/hunttools/cmd/root_flags"

	"github.com/spf13/cobra"
)

func init() {
	checkRepoUsesActionsCmd.Flags().StringVarP(&continuationNWO, "continuation", "c", "", "Continuation NWO")

	GitHubRootCmd.AddCommand(checkRepoUsesActionsCmd)
}

var continuationNWO string

func listWorkflowsForNWO(owner string, repo string) []*github.Workflow {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}
	allWorkflows := []*github.Workflow{}

	opt := &github.ListOptions{
		PerPage: 100,
	}
	for {
		workflows, resp, err := client.Actions.ListWorkflows(ctx, owner, repo, opt)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		allWorkflows = append(allWorkflows, workflows.Workflows...)
		if resp.NextPage == 0 {
			break
		}
		fmt.Println(resp.NextPage)
		opt.Page = resp.NextPage
	}

	return allWorkflows
}
func nwoHasWorkflows(owner string, repo string) (bool, error) {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return false, err
	}

	opt := &github.ListOptions{
		PerPage: 1,
	}
	workflows, _, err := client.Actions.ListWorkflows(ctx, owner, repo, opt)
	if err != nil {
		return false, err
	}
	if *workflows.TotalCount > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

var checkRepoUsesActionsCmd = &cobra.Command{
	Use:   "uses-actions",
	Short: "Checks if the owner/repo pair exists. If it does, we print the nwo",
	Long:  `Checks if the owner/repo pair exists. If it does, we print the nwo`,
	RunE: func(cmd *cobra.Command, args []string) error {
		nwos := getInputTargets()
		continuationNWOHit := false
		if continuationNWO == "" {
			continuationNWOHit = true
		}
		if root.VerboseOutput {
			// Print how many nwos we found
			fmt.Printf("Found %d nwos\n", len(nwos))
			fmt.Println(nwos)
		}
		for _, nwo := range nwos {
			if nwo == "" {
				continue
			}
			if nwo == continuationNWO {
				continuationNWOHit = true
			}
			if !continuationNWOHit {
				continue
			}

			fmt.Println("Checking if ", nwo, " has actions workflows")
			parts := strings.Split(nwo, "/")
			if len(parts) != 2 {
				fmt.Printf("%s is not a valid nwo\n", nwo)
				continue
			}
			for rateLimited := false; !rateLimited; {
				hasWorkflows, err := nwoHasWorkflows(parts[0], parts[1])
				if err != nil {
					if strings.Contains(err.Error(), "404 Not Found") {
						fmt.Println("Repo does not exist")
						rateLimited = false
					} else if strings.Contains(err.Error(), "403 API rate limit exceeded") {
						// Rate limit was hit
						rateLimited = true
						// Sleep for a minute
						fmt.Println("Rate limit hit, sleeping for a minute")
						time.Sleep(time.Minute)
					} else {
						return err
					}

				} else {
					rateLimited = false
					if hasWorkflows {
						fmt.Printf("%s/%s\n", parts[0], parts[1])
					} else if root.VerboseOutput {
						fmt.Printf("%s/%s has no workflows\n", parts[0], parts[1])
					}
				}

			}

		}

		return nil

	},
}
