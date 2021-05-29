package secret_utils

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type SecretType string

const (
	GitHubToken       SecretType = "github_token"
	GitHubSSHKey                 = "github_ssh_key"
	AzureSASToken                = "azure_sas_token"
	SlackToken                   = "slack_token"
	AWSAccessTokens              = "aws_access_tokens"
	HerokuAccessToken            = "heroku_access_token"
	AzureDevOpsPAT               = "azure_devops_pat"
)

func GetSecretTypeFromString(secretType string) (SecretType, error) {
	switch secretType {
	case "github_token":
		return GitHubToken, nil
	case "github_ssh_key":
		return GitHubSSHKey, nil
	case "azure_sas_token":
		return AzureSASToken, nil
	case "slack_token":
		return SlackToken, nil
	case "aws_access_tokens":
		return AWSAccessTokens, nil
	case "heroku_access_token":
		return HerokuAccessToken, nil
	case "azure_devops_pat":
		return AzureDevOpsPAT, nil
	default:
		return "", fmt.Errorf("Unknown secret type: %s", secretType)
	}
}

const GITHUB_TOKEN_TEST_URL = "https://api.github.com/user"
const SLACK_TOKEN_TEST_URL = "https://slack.com/api/auth.test"
const HEROKU_TOKEN_TEST_URL = "https://api.heroku.com/account"

type TestableSecret interface {
	IsValid() (bool, error)
}

type Secret struct {
	SecretType SecretType
	Token      string
}

func (secret Secret) IsValid() (bool, error) {
	switch secret.SecretType {
	case GitHubToken:
		return secret.testGitHubToken()
	case GitHubSSHKey:
		return secret.testGitHubSSHKey()
	case AzureSASToken:
		return secret.testAzureSASToken()
	case SlackToken:
		return secret.testSlackToken()
	case AWSAccessTokens:
		return secret.testAWSAccessTokens()
	case HerokuAccessToken:
		return secret.testHerokuAccessToken()
	case AzureDevOpsPAT:
		return secret.testAzureDevOpsPAT()
	default:
		return false, fmt.Errorf("Unknown secret type: %s", secret.SecretType)
	}
}

func (secret Secret) testGitHubToken() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("GitHub token is empty")
	}
	tokenB64 := b64.StdEncoding.EncodeToString([]byte(secret.Token))
	if testForOkGetResponse(GITHUB_TOKEN_TEST_URL, "Basic", tokenB64, "") {
		return true, nil
	}
	return false, fmt.Errorf("GitHub token is invalid")
}

func (secret Secret) testGitHubSSHKey() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("GitHub SSH Key is empty")
	}
	//TODO: Implement me
	return true, nil
}

func (secret Secret) testAzureSASToken() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("Azure SAS Token is empty")
	}
	//TODO: Implement me
	return true, nil
}

func (secret Secret) testSlackToken() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("Slack Token is empty")
	}
	req, _ := http.NewRequest("GET", SLACK_TOKEN_TEST_URL, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", secret))
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	var responsePayload map[string]string
	json.NewDecoder(res.Body).Decode(responsePayload)
	if responsePayload["ok"] == "true" {
		return true, nil
	} else {
		return false, fmt.Errorf("Slack Token is invalid")
	}
}

func (secret Secret) testAWSAccessTokens() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("AWS Access Tokens is empty")
	}
	// We store AWS Access Tokens separated by a colon
	// Ex: ID:SECRET_KEY
	// We need to split the token and check each one
	tokens := strings.Split(secret.Token, ":")
	aws_access_key_id := tokens[0]
	aws_access_secret_key := tokens[1]

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_access_secret_key, ""),
	})
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	_, err = svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return false, err
		}
	}
	return true, nil
}

func (secret Secret) testHerokuAccessToken() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("Heroku Access Token is empty")
	}
	if testForOkGetResponse(HEROKU_TOKEN_TEST_URL, "Bearer", secret.Token, "application/vnd.heroku+json; version=3") {
		return true, nil
	}
	return false, fmt.Errorf("Heroku Access Token is invalid")
}

func (secret Secret) testAzureDevOpsPAT() (bool, error) {
	if secret.Token == "" {
		return false, fmt.Errorf("Azure DevOps PAT is empty")
	}
	//TODO: Implement me
	return true, nil
}

func testForOkGetResponse(url, authHeaderPrefix, token, accept string) bool {
	req, _ := http.NewRequest("GET", url, nil)

	headerPayload := fmt.Sprintf("%s %s", authHeaderPrefix, token)

	req.Header.Set("Authorization", headerPayload)
	if accept != "" {
		req.Header.Set("Accept", accept)

	}
	res, _ := http.DefaultClient.Do(req)
	if res.StatusCode != 200 {
		return false
	}
	return true
}
