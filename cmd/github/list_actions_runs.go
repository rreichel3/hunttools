package github

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
)

func init() {
	listSelfHostedRuns.MarkFlagRequired("owner")
	listSelfHostedRuns.Flags().IntVarP(&continuationPoint, "continuation", "c", -1, "Contiuation point for where a specific repo dropped")
	GitHubRootCmd.AddCommand(listSelfHostedRuns)
}

func getRunsForNWO(owner string, repo string) []*github.WorkflowRun {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}
	allRuns := []*github.WorkflowRun{}
	opt := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opt)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		allRuns = append(allRuns, runs.WorkflowRuns...)
		if resp.NextPage == 0 {
			break
		}
		fmt.Println(resp.NextPage)
		opt.Page = resp.NextPage
	}

	return allRuns
}

func getActionsWorkflow(owner string, repo string, runID int64) *github.Workflow {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}

	workflow, _, err := client.Actions.GetWorkflowByID(ctx, owner, repo, runID)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return workflow
}

func getRepoContents(owner string, repo string, path string, sha string) *github.RepositoryContent {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}

	opt := &github.RepositoryContentGetOptions{
		Ref: sha,
	}

	contents, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, opt)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return contents

}

func getJobsForRun(owner string, repo string, runID int64) []*github.WorkflowJob {
	client, err := getGitHubClient()
	if err != nil {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil
	}

	allJobsForWorkflow := []*github.WorkflowJob{}

	opt := &github.ListWorkflowJobsOptions{}

	for {

		jobs, resp, err := client.Actions.ListWorkflowJobs(ctx, owner, repo, runID, opt)
		if err != nil {
			fmt.Println("error encountered listing jobs")
			return nil
		}
		allJobsForWorkflow = append(allJobsForWorkflow, jobs.Jobs...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allJobsForWorkflow
}

var continuationPoint int

var listSelfHostedRuns = &cobra.Command{
	Use:   "list-actions-runs",
	Short: "Lists all actions runs for the given organization or repo",
	Long:  `Lists all actions runs for the given organization or repo`,
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
			repos = append(repos, strings.Split(repo, ",")...)
		}

		allWorkflowRuns := []*github.WorkflowRun{}

		for _, repo := range repos {
			fmt.Println("Fetching workflow runs for ", repo)
			allWorkflowRuns = append(allWorkflowRuns, getRunsForNWO(owner, repo)...)
		}
		workflowOutput := []string{}
		// I need to record the owner, repo, runID, jobID, workflow file, and headref
		// Print csv header
		seenFlows := map[string]bool{}
		var f *os.File = nil
		var er error
		csvOutPath := fmt.Sprintf("workflows/%s/results.csv", owner)
		if continuationPoint != -1 {
			f, er = os.OpenFile(csvOutPath, os.O_APPEND|os.O_WRONLY, 0600)
			if er != nil {
				fmt.Println("Error opening existing file - You likely need to not use a continuation point")
				return er
			}
		} else {
			f, er = os.Create(csvOutPath)
			if er != nil {
				fmt.Println("Error creating file")
				return er
			}
			f.WriteString("owner,repo,runID,,workflow,headref,event,workflowURL,selfHosted\n")
			f.Sync()
		}

		defer f.Close()

		for i, workflowRun := range allWorkflowRuns {
			if i < continuationPoint {
				continue
			}
			fmt.Println(i, "/", len(allWorkflowRuns)-1)

			// I need the file contents then I can see if it references `self-hosted`
			workflow := getWorkflow(owner, *workflowRun.Repository.Name, *workflowRun.WorkflowID)
			if workflow == nil {
				fmt.Println("workflow came back nil :( ")
				break
			}
			// Get file contents
			fileContents := getRepoContents(owner, *workflowRun.Repository.Name, workflow.GetPath(), workflowRun.GetHeadSHA())
			if fileContents == nil {
				fmt.Println("fileContents came back nil :( ")
				break
			}
			// Now that we have the contents we need to parse them, save them, etc
			rawFileString, err := fileContents.GetContent()
			if err != nil {
				fmt.Println("Error encountered decoding request")
			}
			uniqueFileID := fmt.Sprintf("%s:%s", workflowRun.GetHeadSHA(), workflow.GetPath())
			selfHosted := strings.Contains(rawFileString, "self-hosted")

			entry := fmt.Sprintf("%s,%s,%d,%s,%s,%s,%s,%t", owner, *workflowRun.Repository.Name, *workflowRun.ID, workflow.GetPath(), workflowRun.GetHeadSHA(), workflowRun.GetEvent(), *workflowRun.WorkflowURL, selfHosted)
			f.WriteString(entry + "\n")
			f.Sync()
			workflowOutput = append(workflowOutput, entry)
			if !seenFlows[uniqueFileID] {
				fullPath := fmt.Sprintf("workflows/%s/%s/%s/%s", owner, *workflowRun.Repository.Name, workflowRun.GetHeadSHA(), workflow.GetPath())
				basePath := ""
				splitPath := strings.Split(fullPath, "/")
				for _, pathSegment := range splitPath[:len(splitPath)-1] {

					basePath = path.Join(basePath, pathSegment)
					ensureDir(basePath)
				}
				os.WriteFile(fullPath, []byte(rawFileString), 0644)
			}

			seenFlows[uniqueFileID] = true

		}
		return nil

	},
}
