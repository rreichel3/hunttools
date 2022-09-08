package github

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
)

func init() {
	cloneRepos.Flags().StringVar(&outdir, "outdir", "./repos/", "output directory")
	GitHubRootCmd.AddCommand(cloneRepos)
}

func cloneRepo(owner string, name string, outputDirectory string) error {
	githubToken, err := getGitHubToken()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://github.com/%s/%s.git", owner, name)

	_, err = git.PlainClone(outputDirectory, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: "hunttools",
			Password: githubToken,
		},
	})
	if err != nil {
		return err
	} else {
		return nil
	}
}

var outdir string

var cloneRepos = &cobra.Command{
	Use:   "clone-repos",
	Short: "Clones the repos provided",
	Long:  `Clones the repos provided`,
	RunE: func(cmd *cobra.Command, args []string) error {

		repos := getInputTargets()
		ensureDir(outdir)

		for _, repo := range repos {
			if repo == "" {
				continue
			}
			fmt.Println(repo)
			split := strings.Split(repo, "/")
			owner := split[0]
			name := split[1]
			ensureDir(fmt.Sprintf("%s/%s", outdir, owner))
			destDir := fmt.Sprintf("%s/%s/%s", outdir, owner, name)
			ensureDir(destDir)
			err := cloneRepo(owner, name, destDir)
			if err != nil {
				fmt.Println("Error encountered cloning: " + repo)
				fmt.Println(err)
			}
		}

		return nil

	},
}
