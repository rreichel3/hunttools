package secret

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	secret_utils "github.com/rreichel3/hunttools/cmd/secret/utils"
	"github.com/spf13/cobra"
)

func init() {

	testCredscanSecretsCmd.Flags().StringVarP(&CredscanResultsDir, "credscan-results-dir", "c", "", "Directory with credscan results")
	addSecretCmd.MarkFlagRequired("credscan-results-dir")

	SecretRootCmd.AddCommand(testCredscanSecretsCmd)
}

func ioReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

var CredscanResultsDir string

// {"type":"ALICLOUD_SECRET_KEY","token":"createRegistrationTokenForRepo","blob":"830ef93df968452499dc41026e116696a6dd8a35","report_url":"https://query.aliyun.com/rest/github.post_tokens","start_line":269,"end_line":269,"start_column":4,"end_column":34}
type CredScanResult struct {
	Type     string
	Token    string
	Blob     string
	Report   string
	Start    int
	End      int
	StartCol int
	EndCol   int
}

var testCredscanSecretsCmd = &cobra.Command{
	Use:   "test-credscan",
	Short: "Tests all secrets from the directory.",
	Long:  `Tests all secrets from the directory.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		// Iterate over every file in the directory
		var files, err = ioReadDir(CredscanResultsDir)
		if err != nil {
			panic(err)
		}
		seenSecrets := make(map[string]bool)
		for idx, filename := range files {
			// Open file and read line by line
			// If filename startswith dsp-testing ignore it
			if strings.HasPrefix(filename, "dsp-testing") {
				continue
			}
			filePath := fmt.Sprintf("%s/%s", CredscanResultsDir, filename)
			println("Opening file: " + filename + "(" + strconv.Itoa(idx) + " of " + strconv.Itoa(len(files)) + ")")
			file, err := os.Open(filePath)
			if err != nil {
				panic(err)
			}

			scanner := bufio.NewScanner(file)
			// optionally, resize scanner's capacity for lines over 64K, see next example

			for scanner.Scan() {
				content := strings.TrimSpace(scanner.Text())
				if len(content) == 0 {
					continue
				}
				var credScanResult CredScanResult
				json.Unmarshal(scanner.Bytes(), &credScanResult)
				if _, ok := seenSecrets[credScanResult.Token]; ok {
					//do something here
					continue
				}
				contentType, _ := secret_utils.GetSecretTypeFromString(credScanResult.Type)
				theSecret := secret_utils.Secret{SecretType: contentType, Token: credScanResult.Token}
				valid, _ := theSecret.IsValid()

				seenSecrets[credScanResult.Token] = valid
				if valid {
					fmt.Println("Wooooooo")
					fmt.Println("VALID: ", content)
				}
				// seenSecrets[credScanResult.Token] = true
			}
			file.Close()
		}
		print(len(seenSecrets))
		return nil
	},
}
