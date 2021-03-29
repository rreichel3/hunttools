package github

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"github.com/rreichel3/hunttools/cmd/utils"

	"github.com/spf13/cobra"
)

func init() {
	danglingPagesCmd.Flags().StringVarP(&DomainJsonListPath, "infile", "i", "", "Json list of subdomains")
	danglingPagesCmd.MarkFlagRequired("infile")
	danglingPagesCmd.Flags().StringVarP(&RootDomain, "rootdomain", "d", "", "The root domain for the provided subdomain lists")
	danglingPagesCmd.MarkFlagRequired("rootdomain")

	GitHubRootCmd.AddCommand(danglingPagesCmd)
}

var DomainJsonListPath string
var RootDomain string

var danglingPagesCmd = &cobra.Command{
	Use:   "find-dangling-pages",
	Short: "Finds dangling GitHub pages",
	Long:  `Takes a JSON list of subdomains (The format from Azure DNS is what's expected), then iterates over them to discover takeoverable domains`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var addresses, err = utils.LoadJsonList(DomainJsonListPath)
		if err != nil {
			return err
		}

		for _, address := range addresses {
			hostname := fmt.Sprintf("%v.%s", address["name"], RootDomain)
			if root.VerboseOutput {
				fmt.Println("Processing host: " + hostname)
			}
			if isDangling(hostname) {
				fmt.Printf("%s\n", hostname)
			}
		}
		return nil

	},
}

func isPages(addr string) bool {
	cmd := exec.Command("dig", addr)
	out, err := cmd.Output()
	if err != nil {
		if root.VerboseOutput {
			fmt.Println(err)
		}
		return false
	}
	output := string(out)
	return strings.Contains(output, "github.io")
}

func isUnallocated(addr string) bool {
	resp, err := http.Get("https://" + addr)
	if err != nil {
		if root.VerboseOutput {
			fmt.Println(err)
		}
		return false
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if root.VerboseOutput {
			fmt.Println(err)
		}
		return false
	}
	bodyString := string(bodyBytes)
	return resp.StatusCode == http.StatusNotFound && strings.Contains(bodyString, "There isn't a GitHub Pages site here.")
}

func isDangling(addr string) bool {
	if isPages(addr) && isUnallocated(addr) {
		return true
	}
	return false

}
