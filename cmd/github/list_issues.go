package github

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {

	GitHubRootCmd.AddCommand(listIssues)
}

func getIssuesForNWO(owner string, name string) []*github.Issue {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}

	allIssues := []*github.Issue{}
	opt := &github.IssueListByRepoOptions{
		State: "all",
	}
	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, name, opt)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues
}

var listIssues = &cobra.Command{
	Use:   "list-issues",
	Short: "Lists issues for the given organization or repo",
	Long:  `Lists issues for the given organization or repo. Piping expects the owner/name syntax`,
	RunE: func(cmd *cobra.Command, args []string) error {

		repos := getInputTargets()

		allIssues := []*github.Issue{}

		for _, repo := range repos {
			if repo == "" {
				continue
			}
			fmt.Println(repo)
			split := strings.Split(repo, "/")
			owner := split[0]
			name := split[1]
			allIssues = append(allIssues, getIssuesForNWO(owner, name)...)
		}
		allIssueOutput := []string{}
		for _, issue := range allIssues {
			content := issue.GetBody()
			// We need to split like this because the Issues response doesn't actually return which repository its attached to
			urlSplit := strings.Split(issue.GetRepositoryURL(), "/")
			localRepo := urlSplit[len(urlSplit)-1]
			owner := urlSplit[len(urlSplit)-2]
			issueNumber := issue.GetNumber()
			issueOutput := fmt.Sprintf("------------------------\n%s/%s/%d\n\n%s", owner, localRepo, issueNumber, content)
			allIssueOutput = append(allIssueOutput, issueOutput)

		}
		for _, url := range allIssueOutput {
			//TODO: I need to use match groups properly
			fmt.Println(url[:len(url)-1])
		}
		return nil

	},
}
