package github

import (
	"errors"
	"fmt"
	"os"

	"github.com/google/go-github/v34/github"
	"golang.org/x/oauth2"
)

func getGitHubClient() (*github.Client, error) {
	auth_token, ok := os.LookupEnv("GITHUB_PAT")
	if !ok {
		fmt.Println("You need to set the GITHUB_PAT environment variable.\n")
		return nil, errors.New("You need to set the GITHUB_PAT environment variable.")
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: auth_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, nil
}

func ensureDir(dirName string) error {
	// fmt.Printf("Ensuring %v exists\n", dirName)
	err := os.Mkdir(dirName, 0777)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
