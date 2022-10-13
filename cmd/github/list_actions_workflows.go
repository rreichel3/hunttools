package github

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	listActionsWorkflows.MarkFlagRequired("owner")
	GitHubRootCmd.AddCommand(listActionsWorkflows)
}

func getWorkflow(owner string, repo string, runID int64) *github.Workflow {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_TOKEN environment variable.\n")
		return nil
	}

	workflow, _, err := client.Actions.GetWorkflowByID(ctx, owner, repo, runID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return workflow
}

func getContents(owner string, repo string, path string) *github.RepositoryContent {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_TOKEN environment variable.\n")
		return nil
	}

	opt := &github.RepositoryContentGetOptions{
		// Ref: sha,
	}

	contents, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opt)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contents

}

var listActionsWorkflows = &cobra.Command{
	Use:   "list-actions-workflows",
	Short: "Lists all actions workflows for the given organization or repo",
	Long:  `Lists all actions workflows for the given organization or repo`,
	RunE: func(cmd *cobra.Command, args []string) error {

		repos := getInputTargets()

		for _, nwo := range repos {
			fmt.Println("Fetching workflow runs for ", nwo)
			splitRepo := strings.Split(nwo, "/")
			owner := splitRepo[0]
			repo := splitRepo[1]
			workflows := listWorkflowsForNWO(owner, repo)
			if len(workflows) > 0 {
				destinationDir := fmt.Sprintf("workflows/%s", owner)
				if _, err := os.Stat(destinationDir); !os.IsNotExist(err) {
					continue
				}
				ensureDir(destinationDir)
			}
			for i, workflow := range workflows {
				if i < continuationPoint {
					continue
				}
				fmt.Println(i, "/", len(workflows)-1)

				// I need the file contents then I can see if it references `self-hosted`
				workflow := getWorkflow(owner, repo, *workflow.ID)
				if workflow == nil {
					fmt.Println("workflow came back nil :( ")
					break
				}
				// Get file contents
				fileContents := getContents(owner, repo, workflow.GetPath())
				if fileContents == nil {
					fmt.Println("fileContents came back nil :( ")
					continue
				}
				// Now that we have the contents we need to parse them, save them, etc
				rawFileString, err := fileContents.GetContent()
				if err != nil {
					fmt.Println("Error encountered decoding request")
				}
				// uniqueFileID := fmt.Sprintf("%s:%s", workflowRun.GetHeadSHA(), workflow.GetPath())
				// selfHosted := strings.Contains(rawFileString, "self-hosted")

				fullPath := fmt.Sprintf("workflows/%s/%s/%s", owner, repo, workflow.GetPath())
				basePath := ""
				splitPath := strings.Split(fullPath, "/")
				for _, pathSegment := range splitPath[:len(splitPath)-1] {

					basePath = path.Join(basePath, pathSegment)
					ensureDir(basePath)
				}
				os.WriteFile(fullPath, []byte(rawFileString), 0644)
			}

		}

		return nil

	},
}
