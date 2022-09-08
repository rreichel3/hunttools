package github

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {

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

var listIssuesAttachments = &cobra.Command{
	Use:   "list-issue-attachments",
	Short: "Lists issue attachements for the given organization or repo",
	Long:  `Lists issue attachements for the given organization or repo. Piping expects the owner/name syntax`,
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
			owner := urlSplit[len(urlSplit)-2]
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
