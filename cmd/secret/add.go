package secret

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"syscall"
)

func init() {

	addSecretCmd.MarkFlagRequired("name")

	addSecretCmd.Flags().StringVarP(&SecretType, "type", "t", "", "Secret type")
	addSecretCmd.MarkFlagRequired("type")

	SecretRootCmd.AddCommand(addSecretCmd)
}

var SecretType string

var addSecretCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a secret to the given vault. Prompts for the secret",
	Long:  `Adds a secret to the given vault. Prompts for the secret`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vaultUri := fmt.Sprintf("https://%s.vault.azure.net/", VaultName)
		authorizer, err := auth.NewAuthorizerFromCLIWithResource("https://vault.azure.net")
		if err != nil {
			fmt.Println("Unable to get an authorizer. Please run `az login`.")
			fmt.Println(err)

			return err
		}

		client := keyvault.New()
		client.Authorizer = authorizer

		secret, err := getSecretFromUser()
		if err != nil {
			return err
		}

		secretParams := keyvault.SecretSetParameters{
			ContentType: &SecretType,
			Value:       &secret,
		}

		client.SetSecret(ctx, vaultUri, SecretName, secretParams)
		fmt.Println("Secret added successfully.")
		return nil

	},
}

func promptForSecret() (string, error) {

	fmt.Print("Enter Secret Name: ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}

func getSecretFromUser() (string, error) {

	fmt.Print("Enter Secret: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}
