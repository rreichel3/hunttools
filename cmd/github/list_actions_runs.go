package github

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"

	"github.com/spf13/cobra"
)

func init() {
	listSelfHostedRuns.MarkFlagRequired("owner")

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
		fmt.Println(len(allRuns))
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRuns
}

func getWorkflow(owner string, repo string, runID int64) *github.Workflow {
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
		for i, workflowRun := range allWorkflowRuns {
			fmt.Println(i, "/", len(allWorkflowRuns))
			// Temporary workaround to continue where we stopped
			// if i < 2748 {
			// 	continue
			// }
			fmt.Println("Getting workflow")
			workflow := getWorkflow(owner, *workflowRun.Repository.Name, *workflowRun.WorkflowID)
			actionUniqueID := workflowRun.GetHeadSHA() + ":" + workflow.GetPath()
			// We do this to decrease the number of calls we need to make
			if seenFlows[actionUniqueID] {
				continue
			} else {
				seenFlows[actionUniqueID] = true
			}
			fmt.Println("Getting jobs")
			jobs := getJobsForRun(owner, *workflowRun.Repository.Name, *workflowRun.ID)
			fmt.Println("Got", len(jobs), "jobs")

			for _, job := range jobs {
				groupID := int64(-1)
				if job.RunnerGroupID != nil {
					groupID = *job.RunnerGroupID
				}
				selfHosted := false
				for _, label := rangejob.Labels {
					if label == "self-hosted" {
						selfHosted = true
						break
					}
				}
				entry := fmt.Sprintf("%s,%s,%d,%d,%d,%s,%s,%s,%s", owner, *workflowRun.Repository.Name, *workflowRun.ID, job.GetID(), groupID, workflow.GetPath(), job.GetHeadSHA(), workflowRun.GetEvent(), *workflowRun.WorkflowURL)
				workflowOutput = append(workflowOutput, entry)

			}

		}
		fmt.Println("owner,repo,runID,jobID,runnerGroup,workflow,headref,event,workflowURL")
		for _, workflowEntry := range workflowOutput {
			//TODO: I need to use match groups properly
			fmt.Println(workflowEntry)
		}
		return nil

	},
}
