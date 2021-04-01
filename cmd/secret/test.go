package secret

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	secret_utils "github.com/rreichel3/hunttools/cmd/secret/utils"
	"github.com/spf13/cobra"
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
			token, err := getSecret(&client, secretName)
			if err != nil {
				fmt.Println("Unable to get secret from the vault.")
				fmt.Println(err)
			}
			contentType, _ := secret_utils.GetSecretTypeFromString(*secret.ContentType)
			theSecret := secret_utils.Secret{SecretType: contentType, Token: token}
			valid, err := theSecret.IsValid()
			if valid {
				fmt.Println(secretName, "is valid.")
			} else {
				fmt.Println(secretName, "is invalid.")
			}
		}

		return nil
	},
}

func getSecretNameFromID(id string) string {
	return strings.Split(id, "/")[len(strings.Split(id, "/"))-1]
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
