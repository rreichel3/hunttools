package github

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	listActionsWorkflows.MarkFlagRequired("owner")
	GitHubRootCmd.AddCommand(listActionsSecrets)
}

var listActionsSecrets = &cobra.Command{
	Use:   "list-actions-secrets",
	Short: "Lists all actions secrets for the given organization or repo",
	Long:  `Lists all actions secrets for the given organization or repo`,
	RunE: func(cmd *cobra.Command, args []string) error {

		repos := getInputTargets()

		for _, nwo := range repos {
			// fmt.Println("Fetching secrets runs for ", nwo)
			fmt.Println(nwo)
			splitRepo := strings.Split(nwo, "/")
			owner := splitRepo[0]
			repo := splitRepo[1]

			client, err := getGitHubClient()
			if err != nil {
				return err
			}
			secrets, _, err := client.Actions.ListRepoSecrets(ctx, owner, repo, nil)
			if err != nil {
				return err
			}
			for _, secret := range secrets.Secrets {
				fmt.Printf("%s/%s/%s\n", owner, repo, secret.Name)
			}
		}

		return nil

	},
}
