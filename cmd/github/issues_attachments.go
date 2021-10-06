package github

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
)

func init() {
	listIssuesAttachments.MarkFlagRequired("owner")

	GitHubRootCmd.AddCommand(listIssuesAttachments)
}

var fileUploadRegex = regexp.MustCompile(`(https:\/\/github.com\/.*\/files\/.*)\)`)

func getCommentsForIssue(owner string, name string, number int) []*github.IssueComment {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}

	allComments := []*github.IssueComment{}
	opt := &github.IssueListCommentsOptions{}
	for {
		comments, resp, err := client.Issues.ListComments(ctx, owner, name, number, opt)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allComments
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

var listIssuesAttachments = &cobra.Command{
	Use:   "list-issue-attachments",
	Short: "Lists issue attachements for the given organization or repo",
	Long:  `Lists issue attachements for the given organization or repo`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var repos = []string{}
		if repo == "" {
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
			//TODO: List all org repos
			opt := &github.RepositoryListByOrgOptions{
				Type: "all",
			}

			for {
				reposListResponse, resp, err := client.Repositories.ListByOrg(ctx, owner, opt)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				for _, repo := range reposListResponse {
					repos = append(repos, *repo.Name)
				}
				if resp.NextPage == 0 {
					break
				}
				opt.Page = resp.NextPage
			}
		} else {
			repos = append(repos, repo)
		}
		allIssues := []*github.Issue{}

		for _, repo := range repos {
			allIssues = append(allIssues, getIssuesForNWO(owner, repo)...)
		}
		// TODO: Scrape each issue and get the attachments
		allFileUrls := []string{}
		for _, issue := range allIssues {
			fileUrls := fileUploadRegex.FindAll([]byte(issue.GetBody()), -1)
			for _, url := range fileUrls {
				allFileUrls = append(allFileUrls, string(url))
			}
			// We need to split like this because the Issues response doesn't actually return which repository its attached to
			urlSplit := strings.Split(issue.GetRepositoryURL(), "/")
			localRepo := urlSplit[len(urlSplit)-1]
			issueNumber := issue.GetNumber()
			comments := getCommentsForIssue(owner, localRepo, issueNumber)
			for _, comment := range comments {
				fileUrls := fileUploadRegex.FindAll([]byte(comment.GetBody()), -1)
				for _, url := range fileUrls {
					allFileUrls = append(allFileUrls, string(url))
				}
			}
		}
		for _, url := range allFileUrls {
			//TODO: I need to use match groups properly
			fmt.Println(url[:len(url)-1])
		}
		return nil

	},
}
