package secret

import (
	"context"
	"github.com/spf13/cobra"
)

func init() {
	SecretRootCmd.Flags().StringVarP(&VaultName, "vault-name", "", "red-master-secrets", "Vault Name")
	SecretRootCmd.Flags().StringVarP(&SecretName, "name", "n", "", "Secret name")

}

var VaultName string
var SecretName string

var ctx = context.Background()

var SecretRootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Commands to do magic with secrets and Azure Key Vault",
	Long:  `Commands to do magic with secrets and Azure Key Vault`,
}
