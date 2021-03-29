package network

import (
	"fmt"

	"github.com/go-ping/ping"
	root "github.com/rreichel3/hunttools/cmd/root_flags"
	"github.com/rreichel3/hunttools/cmd/utils"

	"github.com/spf13/cobra"
)

func init() {
	pingCmd.Flags().StringVarP(&PingDestinationsPath, "infile", "i", "", "Newline delimited file of ping destinations")
	pingCmd.MarkFlagRequired("infile")
	NetworkRootCmd.AddCommand(pingCmd)
}

var PingDestinationsPath string

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Ping a list of hosts",
	Long:  `Runs the ping command against a list of hosts`,
	RunE: func(cmd *cobra.Command, args []string) error {

		var pingDestinations, err = utils.ReadFileToList(PingDestinationsPath)
		if err != nil {
			return err
		}

		for _, address := range pingDestinations {
			alive, err := pingAddress(address)

			if alive {
				if root.VerboseOutput {
					fmt.Printf("%s, UP\n", address)
				} else {
					fmt.Println(address)
				}
			} else {
				if root.VerboseOutput {
					if err != nil {
						fmt.Printf("%s, Error: %s\n", address, err)
					}
					fmt.Printf("%s, DOWN\n", address)
				}
			}
		}
		return nil

	},
}

func pingAddress(addr string) (bool, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return false, err
	}
	pinger.Count = 1
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return false, err
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	if stats.PacketLoss != 0 {
		return false, nil
	}
	return true, nil
}
