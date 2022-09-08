package github

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	GitHubRootCmd.PersistentFlags().StringVarP(&owner, "owner", "o", "", "The owner to query")
	GitHubRootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "The repo to query")
	GitHubRootCmd.PersistentFlags().BoolVarP(&usePipe, "pipe", "p", false, "Use pipe to read from stdin. Overrides owner and repo flags")
}

var ctx = context.Background()

var owner string
var repo string
var usePipe bool

func getContentFromStdin() []string {
	lines, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Split(string(lines), "\n")
}

func getInputTargets() []string {
	if usePipe {
		return getContentFromStdin()
	} else {
		// Handle one specific owner + repo pair
		var targets []string
		if repo != "" {
			if strings.Contains(owner, ",") {
				fmt.Println("Error encountered. We can't have a comma separated owner along with a specified repo")
				return []string{}
			}
			// This branch handles multiple repos, one owner
			repos := strings.Split(repo, ",")
			for _, r := range repos {
				targets = append(targets, owner+"/"+r)
			}
		} else {
			// Repo wasn't set, lets grab all of the owners (we support csv lists)
			owners := strings.Split(owner, ",")
			for _, o := range owners {
				targets = append(targets, o)
			}
			return targets
		}
		return targets
	}
}

var GitHubRootCmd = &cobra.Command{
	Use:   "gh",
	Short: "GitHub related commands ",
	Long:  `GitHub related commands`,
}
