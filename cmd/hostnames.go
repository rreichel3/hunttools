package cmd

import (
	"fmt"
	"net"

	"github.com/rreichel3/hunttools/cmd/utils"

	"github.com/spf13/cobra"
)

func init() {
	hostnamesCommand.Flags().StringVarP(&DestinationsPath, "in", "i", "", "Newline delimited file IPs to fetch hostname")
	hostnamesCommand.Flags().BoolVarP(&HostnamesAliveFlag, "alive", "a", false, "Only return successful hostnames")
	rootCmd.AddCommand(pingCmd)
}

var DestinationsPath string
var HostnamesAliveFlag bool

var hostnamesCommand = &cobra.Command{
	Use:   "hostnames",
	Short: "Get hostnames for a list of IP Addresses",
	Long:  `Gets the hostnames for a list of IP Addresses.`,
	Run: func(cmd *cobra.Command, args []string) {

		var hostnameIPs, err = utils.ReadFileToList(DestinationsPath)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, each_ln := range hostnameIPs {
			addr, err := net.LookupAddr(each_ln)
			if err != nil {
				if !HostnamesAliveFlag {
					fmt.Println(err)
				}
			} else {
				fmt.Println(addr)
			}
		}

	},
}
