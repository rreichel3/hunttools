package secret

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

func init() {

	SecretRootCmd.AddCommand(testSecretCmd)
}

var VaultUri string
var testSecretCmd = &cobra.Command{
	Use:   "test-all",
	Short: "Tests all secrets from the vault if possible.",
	Long:  `Tests all secrets from the vault if possible.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		VaultUri = fmt.Sprintf("https://%s.vault.azure.net/", VaultName)
		authorizer, err := auth.NewAuthorizerFromCLIWithResource("https://vault.azure.net")
		if err != nil {
			fmt.Println("Unable to get an authorizer. Please run `az login`.")
			fmt.Println(err)

			return err
		}

		client := keyvault.New()
		client.Authorizer = authorizer
		secretsToFetch := int32(25)
		vaultSecrets, err := client.GetSecretsComplete(ctx, VaultUri, &secretsToFetch)
		if err != nil {
			fmt.Println("Unable to get secrets from the vault.")
			fmt.Println(err)
		}
		for ; vaultSecrets.NotDone(); vaultSecrets.NextWithContext(ctx) {
			secret := vaultSecrets.Value()
			secretName := getSecretNameFromID(*secret.ID)
			switch contentType := secret.ContentType; *contentType {
			case "github_token":
				if testGitHubToken(&client, secretName) {
					fmt.Println(secretName, "is valid")
				} else {
					fmt.Println(secretName, "is invalid")
				}
			case "test":
				fmt.Println("Got a test token")
			default:
				fmt.Println("Content type didn't match an automated testing mechanism:", *contentType)
			}
		}

		return nil
	},
}

func getSecretNameFromID(id string) string {
	return strings.Split(id, "/")[len(strings.Split(id, "/"))-1]
}

func testGitHubToken(client *keyvault.BaseClient, secretName string) bool {

	token, err := getSecret(client, secretName)
	if err != nil {
		return false
	}
	url := "https://api.github.com/user"

	req, _ := http.NewRequest("GET", url, nil)

	tokenB64 := b64.StdEncoding.EncodeToString([]byte(token))

	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", tokenB64))

	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		return false
	}
	return true
}

func getSecret(client *keyvault.BaseClient, secretName string) (string, error) {
	bundle, err := client.GetSecret(ctx, VaultUri, secretName, "")
	if err != nil {
		fmt.Println("Unable to get the GitHub token.")
		fmt.Println(err)
		return "", err
	}
	token := *bundle.Value
	return token, nil
}
