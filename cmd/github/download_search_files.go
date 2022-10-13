package github

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

func init() {
	downloadSearchFiles.Flags().StringVarP(&searchQuery, "query", "q", "", "The search query to run")
	downloadSearchFiles.MarkFlagRequired("query")
	downloadSearchFiles.Flags().StringVarP(&outDir, "outdir", "d", "./search-results", "")
	downloadSearchFiles.Flags().IntVarP(&maxCount, "count", "n", -1, "Max count. This may yield more but will stop within 30")
	GitHubRootCmd.AddCommand(downloadSearchFiles)
}

const DELAY_FOR_SEARCH = 25 * time.Second

var searchQuery string
var outDir string
var maxCount int
var downloadSearchFiles = &cobra.Command{
	Use:   "download-search-files",
	Short: "Downloads contents of all search result files",
	Long:  `Finds all of a user's public gists`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getGitHubClient()
		if err != nil {
			fmt.Println("You need to set the GITHUB_TOKEN environment variable.\n")
			return nil
		}

		err = ensureDir(outDir)
		if err != nil {
			fmt.Println("Error encountered trying to create your directory.\n")
			return nil
		}

		allSearchResults := []*github.CodeResult{}
		opt := &github.SearchOptions{TextMatch: false}
		for {
			results, resp, err := client.Search.Code(ctx, searchQuery, opt)
			// fmt.Printf("%v results\n", *results.Total)
			if resp.StatusCode == 403 {
				fmt.Println("Sleeping")
				fmt.Printf("%d found so far\n", len(allSearchResults))
				time.Sleep(DELAY_FOR_SEARCH)
				continue
			}
			if err != nil {
				fmt.Println(err)
				return nil
			}
			allSearchResults = append(allSearchResults, results.CodeResults...)
			if resp.NextPage == 0 {
				fmt.Println("Exiting because next page is 0")
				break
			}
			opt.Page = resp.NextPage

			if maxCount != -1 && len(allSearchResults) > maxCount {
				fmt.Println("Exiting due to max count hit")
				break
			}

		}

		for _, codeResult := range allSearchResults {
			owner := *codeResult.Repository.Owner.Login
			name := *codeResult.Repository.Name
			filepath := *codeResult.Path
			ownerPath := path.Join(outDir, owner)
			namePath := path.Join(ownerPath, name)
			if strings.Contains(filepath, "nextjs") {
				continue
			}
			err = ensureDir(ownerPath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = ensureDir(namePath)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// fmt.Printf("%v\n", *codeResult.Name)
			contentOps := &github.RepositoryContentGetOptions{}
			fileReadCloser, _, err := client.Repositories.DownloadContents(ctx, owner, name, filepath, contentOps)
			if err != nil {
				fmt.Println(err)
				return nil
			}
			defer fileReadCloser.Close()
			newFilename := strings.ReplaceAll(*codeResult.Path, "/", "_")
			f, e := os.Create(path.Join(namePath, newFilename))
			if e != nil {
				panic(e)
			}
			defer f.Close()
			f.ReadFrom(fileReadCloser)
		}

		return nil

	},
}
