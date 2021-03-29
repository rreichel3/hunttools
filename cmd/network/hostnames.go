package network

import (
	"fmt"
	"net"

	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"github.com/rreichel3/hunttools/cmd/utils"

	"github.com/spf13/cobra"
)

func init() {
	hostnamesCommand.Flags().StringVarP(&DestinationsPath, "infile", "i", "", "Newline delimited file of IPs for which to fetch their hostname")
	hostnamesCommand.MarkFlagRequired("infile")
	NetworkRootCmd.AddCommand(hostnamesCommand)
}

var DestinationsPath string

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
				if root.VerboseOutput {
					fmt.Println(err)
				}
			} else {
				fmt.Println(addr)
			}
		}

	},
}
